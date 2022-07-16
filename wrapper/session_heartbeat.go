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
func (s *Session) beat(bot *Client) {
	s.routines.Done()

	for {
		select {
		case hb := <-s.heartbeat.send:
			s.Lock()

			// close the connection if the last sent Heartbeat never received a HeartbeatACK.
			if atomic.LoadUint32(&s.heartbeat.acks) == 0 {
				log.Printf("attempting to reconnect Session %q due to no HeartbeatACK", s.ID)
				if err := s.reconnect(bot); err != nil {
					log.Println(ErrorDisconnect{
						SessionID: s.ID,
						Err:       err,
						Action:    fmt.Errorf("no HeartbeatACK"),
					})
				}

				s.Unlock()

				return
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
				s.disconnectFromRoutine("heartbeat: Closing the connection due to a write error...", err)

				s.Unlock()

				return
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
			return
		}
	}
}

// pulse generates Opcode 1 Heartbeats for a Session's heartbeat channel.
func (s *Session) pulse() {
	s.routines.Done()

	// send an Opcode 1 Heartbeat payload after heartbeat_interval * jitter milliseconds
	// (where jitter is a random value between 0 and 1).
	s.Lock()
	s.heartbeat.send <- Heartbeat{Data: s.Seq}
	log.Println("queued jitter heartbeat")
	s.Unlock()

	for {
		s.Lock()

		select {
		default:
			s.Unlock()

		// every Heartbeat Interval...
		case <-s.heartbeat.ticker.C:

			// queue a heartbeat.
			s.heartbeat.send <- Heartbeat{Data: s.Seq}

			log.Println("queued heartbeat")

			s.Unlock()

		case <-s.Context.Done():
			s.Unlock()

			return
		}
	}
}

// respond responds to Opcode 1 Heartbeats from the Discord Gateway.
func (s *Session) respond(data json.RawMessage) {
	var heartbeat Heartbeat
	if err := json.Unmarshal(data, &heartbeat); err != nil {
		s.Lock()
		defer s.Unlock()

		s.disconnectFromRoutine("respond: Closing the connection due to an unmarshal error...",
			ErrorEvent{
				Event:  FlagGatewayCommandNameHeartbeat,
				Err:    err,
				Action: ErrorEventActionUnmarshal,
			})
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
}
