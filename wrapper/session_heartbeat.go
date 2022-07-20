package wrapper

import (
	"encoding/json"
	"fmt"
	"log"
	"sync/atomic"
	"time"
)

// heartbeat represents the heartbeat mechanism for a Session.
type heartbeat struct {
	// interval represents the interval of time between each Heartbeat Payload.
	interval time.Duration

	// ticker is a timer used to time the interval between each Heartbeat Payload.
	ticker *time.Ticker

	// send represents a channel of heartbeats that will be sent to the Discord Gateway.
	send chan Heartbeat

	// acks represents the amount of times a HeartbeatACK was received since the last Heartbeat.
	acks uint32
}

// Monitor returns the current amount of HeartbeatACKs for a Session's heartbeat.
func (s *Session) Monitor() uint32 {
	s.Lock()
	acks := atomic.LoadUint32(&s.heartbeat.acks)
	s.Unlock()

	return acks
}

// beat listens for pulses to send Opcode 1 Heartbeats to the Discord Gateway (to verify the connection is alive).
func (s *Session) beat(bot *Client) error {
	s.manager.routines.Done()

	// ensure that all pulse routines are closed prior to closing.
	defer func() {
		for {
			select {
			case <-s.heartbeat.send:
			case <-s.Context.Done():
				s.Lock()

				if atomic.LoadInt32(&s.manager.pulses) != 0 {
					s.Unlock()

					break
				}

				s.logClose("heartbeat")
				s.Unlock()

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

				s.reconnect(fmt.Sprintf("attempting to reconnect Session %q due to no HeartbeatACK", s.ID))

				return nil
			}

			// prevent two Heartbeat Payloads being sent to the Discord Gateway consecutively within nanoseconds,
			// when the ticker queues a Heartbeat while the listen thread (onPayload) queues a Heartbeat
			// (in response to the Discord Gateway).
			//
			// clear queued (outdated) heartbeats.
			for len(s.heartbeat.send) > 0 {
				// ensure the latest sequence is sent.
				if h := <-s.heartbeat.send; h.Data > hb.Data {
					hb.Data = h.Data
				}
			}

			// send a Heartbeat to the Discord Gateway (WebSocket Connection).
			if err := writeEvent(s, FlagGatewayOpcodeHeartbeat, FlagGatewayCommandNameHeartbeat, hb); err != nil {
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

			log.Println("sent heartbeat")

			s.Unlock()

		case <-s.Context.Done():
			return nil
		}
	}
}

// pulse generates Opcode 1 Heartbeats for a Session's heartbeat channel.
func (s *Session) pulse() {
	s.manager.routines.Done()
	defer s.decrementPulses()

	// send an Opcode 1 Heartbeat payload after heartbeat_interval * jitter milliseconds
	// (where jitter is a random value between 0 and 1).
	s.Lock()
	s.heartbeat.send <- Heartbeat{Data: s.Seq}
	log.Println("queued jitter heartbeat")
	s.Unlock()

	for {
		select {
		// every Heartbeat Interval...
		case <-s.heartbeat.ticker.C:
			s.Lock()

			// queue a heartbeat.
			s.heartbeat.send <- Heartbeat{Data: s.Seq}

			log.Println("queued heartbeat")

			s.Unlock()

		case <-s.Context.Done():
			s.Lock()
			s.logClose("pulse")
			s.Unlock()

			return
		}
	}
}

// respond responds to Opcode 1 Heartbeats from the Discord Gateway.
func (s *Session) respond(data json.RawMessage) error {
	defer s.decrementPulses()

	var heartbeat Heartbeat
	if err := json.Unmarshal(data, &heartbeat); err != nil {
		return fmt.Errorf("error occurred unmarshalling incoming Heartbeat: %w", err)
	}

	atomic.StoreInt64(&s.Seq, heartbeat.Data)

	s.Lock()

	// heartbeat() checks for the amount of HeartbeatACKs received since the last Heartbeat.
	// There is a possibility for this value to be 0 due to latency rather than a dead connection.
	// For example, when a Heartbeat is queued, sent, responded, and sent.
	//
	// Prevent this possibility by treating this response from Discord as an indication that the
	// connection is still alive.
	atomic.AddUint32(&s.heartbeat.acks, 1)

	// send an Opcode 1 Heartbeat without waiting the remainder of the current interval.
	s.heartbeat.send <- heartbeat

	log.Println("responded to heartbeat")

	s.Unlock()

	return nil
}
