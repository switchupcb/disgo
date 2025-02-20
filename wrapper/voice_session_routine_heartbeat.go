package wrapper

import (
	"sync/atomic"
	"time"
)

// voice_heartbeat represents the heartbeat mechanism for a Voice Session.
type voice_heartbeat struct {
	// interval represents the interval of time between each Heartbeat Payload.
	interval time.Duration

	// ticker is a timer used to time the interval between each Heartbeat Payload.
	ticker *time.Ticker

	// send represents a channel of heartbeats that will be sent to the Discord Gateway.
	send chan VoiceHeartbeat

	// acks represents the amount of times a HeartbeatACK was received since the last Heartbeat.
	acks uint32
}

// Monitor returns the current amount of HeartbeatACKs for a Voice Session's heartbeat.
func (s *VoiceSession) Monitor() uint32 {
	s.Lock()
	acks := atomic.LoadUint32(&s.heartbeat.acks)
	s.Unlock()

	return acks
}

// beat listens for pulses to send Opcode 1 Heartbeats to the Discord Voice Server (to verify the connection is alive).
func (s *VoiceSession) beat() error {
	s.manager.routines.Done()

	// confirm all pulse routines are closed prior to closing.
	defer func() {
		for {
			if s.heartbeat == nil {
				s.logClose("heartbeat")

				return
			}

			select {
			case <-s.heartbeat.send:
			case <-s.Context.Done():
				if atomic.LoadInt32(&s.manager.pulses) != 0 {
					break
				}

				s.logClose("heartbeat")

				return
			}
		}
	}()

	for {
		select {
		case hb := <-s.heartbeat.send:
			s.Lock()

			// close the connection if the last sent Heartbeat never received a HeartbeatACK.
			if atomic.LoadUint32(&s.heartbeat.acks) == 0 {
				s.Unlock()

				s.reconnect("attempting to reconnect voice session due to no HeartbeatACK")

				return nil
			}

			// prevent two Heartbeat Payloads being sent to the Discord Voice Server consecutively within nanoseconds,
			// when the ticker queues a Heartbeat while the listen thread (onPayload) queues a Heartbeat
			// (in response to the Discord Voice Server).
			//
			// clear queued (outdated) heartbeats.
			for len(s.heartbeat.send) > 0 {
				// ensure the latest sequence is sent.
				if h := <-s.heartbeat.send; h.Data > hb.Data {
					hb.Data = h.Data
				}
			}

			// send a Heartbeat to the Discord Voice Server (WebSocket Connection).
			if err := hb.SendEvent(s); err != nil {
				s.Unlock()

				return err
			}

			// reset the ticker (and empty existing ticks).
			s.heartbeat.ticker.Reset(s.heartbeat.interval)
			for len(s.heartbeat.ticker.C) > 0 {
				<-s.heartbeat.ticker.C
			}

			// reset the amount of HeartbeatACKs since the last heartbeat.
			atomic.StoreUint32(&s.heartbeat.acks, 0)

			LogSession(Logger.Info(), s.ID).Msg("sent heartbeat")

			s.Unlock()

		case <-s.Context.Done():
			return nil
		}
	}
}

// pulse generates Opcode 3 Heartbeats for a Voice Session's heartbeat channel.
func (s *VoiceSession) pulse() {
	s.manager.routines.Done()

	// send an Opcode 3 Heartbeat payload after heartbeat_interval * jitter milliseconds
	// (where jitter is a random value between 0 and 1).
	s.Lock()
	s.heartbeat.send <- VoiceHeartbeat{Data: atomic.LoadInt64(&s.Nonce)}
	LogSession(Logger.Info(), s.ID).Msg("queued jitter voice heartbeat")
	s.Unlock()

	for {
		select {
		// every Heartbeat Interval...
		case <-s.heartbeat.ticker.C:
			s.Lock()

			// queue a heartbeat.
			s.heartbeat.send <- VoiceHeartbeat{Data: atomic.LoadInt64(&s.Nonce)}

			LogSession(Logger.Info(), s.ID).Msg("queued heartbeat")

			s.Unlock()

		case <-s.Context.Done():
			s.Lock()
			s.logClose("pulse")
			s.Unlock()

			return
		}
	}
}
