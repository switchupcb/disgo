package wrapper

import (
	"bytes"
	"compress/zlib"
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

// TODO: ensure disconnections upon unrecoverable errors.
// TODO: simplify code (make it more readeable and easier to understand and follow)
// TODO: fix data races
// TODO: handle generation
// TODO: ensure conn, write, read is correct with regards to concurrency.

const (
	module                    = "github.com/switchupcb/disgo"
	gatewayDisconnectMsg      = "Disconnected Session %s from the Discord Gateway with code %d"
	gatewayEncoding           = "json"
	invalidSessionWaitTime    = 3 * time.Second
	maxIdentifyLargeThreshold = 250
)

// Session represents a Discord Gateway WebSocket Session.
type Session struct {
	// ID represents the session ID of the Session.
	ID string

	// Seq represents the last sequence number received by the client.
	//
	// https://discord.com/developers/docs/topics/gateway#heartbeat
	Seq int64

	// Endpoint represents the endpoint that is used to connect to the Gateway.
	Endpoint string

	// Connected represents whether the session is connected to Discord Gateway.
	//
	// Connected is a channel because it's used in a select statement.
	// When Connected == nil, the Session is NOT connected.
	Connected chan bool

	// Context carries request-scoped data for the Discord Gateway Connection.
	Context *context.Context

	// Conn represents a connection to the Discord Gateway.
	Conn *websocket.Conn

	// heartbeat contains the fields required to implement the heartbeat mechanism.
	heartbeat *heartbeat

	// mu represents the session mutex of this session.
	mu *sessionMutex

	cancel context.CancelFunc
}

// heartbeat represents the heartbeat mechanism for a session.
type heartbeat struct {
	// ticker is a timer used to time the interval between each Heartbeat Payload.
	ticker *time.Ticker

	// interval represents the interval of time between each Heartbeat Payload.
	interval time.Duration

	// send represents a channel of heartbeats that will be sent to the Discord Gateway.
	send chan GatewayPayload

	// respond represents a channel of heartbeats in response to an Opcode 1 Heartbeat
	// that will be sent to the Discord Gateway.
	//
	// respond prevents a theoretical race condition where the ticker queues a Heartbeat,
	// while the listen thread (onPayload) queues a Heartbeat (in response to the Discord Gateway),
	// resulting in two Heartbeat Payloads being sent to the Discord Gateway consecutively within nanoseconds,
	// while the first Heartbeat Payload is outdated.
	respond chan GatewayPayload

	// ack represents a channel of times a HeartbeatACK was received.
	ack chan time.Time
}

// sessionMutex represents a Session's mutexes used to prevent race conditions.
type sessionMutex struct {
	// connect ensures that the Connect and Disconnect functionality is only running on one goroutine per Session.
	//
	// Prevents a theoretical race condition where the main thread starts a heartbeat() goroutine in Connect(),
	// but fails to receive an ACK (within the FailedHeartbeatInterval) BEFORE Connect() returns,
	// resulting in the session connecting while also attempting to reconnect.
	//
	// Prevents a theoretical race condition where the main thread attempts to Connect(),
	// while another thread calls s.Disconnect() resulting in undefined behavior while connecting
	// to the Discord Gateway.
	connect sync.Mutex

	// conn represents the mutex for connections and contexts.
	//
	// TODO: It is unclear whether this needs to be used,
	// as the websocket library Disgo currently uses states concurrent features.
	// However, it imports gorilla/websocket which states
	// "All methods may be called concurrently except for Reader and Read."
	conn sync.RWMutex
}

// isConnected returns whether the session is connected.
func (s *Session) isConnected() bool {
	return s.Connected != nil
}

// canReconnect determines whether the session is in a valid state to reconnect.
func (s *Session) canReconnect() bool {
	return s.ID != "" && atomic.LoadInt64(&s.Seq) != 0
}

// Connect creates or reestablishes a session's open connection to the Discord Gateway.
func (bot *Client) Connect(s *Session) error {
	s.mu.connect.Lock()
	defer s.mu.connect.Unlock()

	if s.isConnected() {
		return nil
	}

	// reconnect returns whether the session is in a valid state to reconnect.
	reconnect := s.canReconnect()
	reconnectID := s.ID
	reconnectSeq := atomic.LoadInt64(&s.Seq)

	// request a valid Gateway URL endpoint from the Discord API.
	gateway := GetGateway{}
	response, err := gateway.Send(bot)
	if err != nil {
		return fmt.Errorf("an error occurred getting the Gateway API Endpoint\n%w", err)
	}

	// TODO: zlib compression
	s.Endpoint = response.URL + "?v=" + VersionDiscordAPI + "&encoding=" + gatewayEncoding

	// connect to the Discord Gateway Websocket.
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)

	s.cancel = cancel
	//defer cancel()

	s.Context = &ctx
	s.Conn, _, err = websocket.Dial(*s.Context, s.Endpoint, nil)
	if err != nil {
		return fmt.Errorf("an error occurred while connecting to the Discord Gateway\n%w", err)
	}

	//defer s.Conn.Close(websocket.StatusInternalError, "StatusInternalError")

	// handle the incoming Hello event upon connecting to the Gateway.
	var payload GatewayPayload
	if err = wsjson.Read(*s.Context, s.Conn, &payload); err != nil {
		return ErrorEvent{
			Event:  FlagGatewayEventNameHello,
			Err:    err,
			Action: ErrorEventActionRead,
		}
	}

	var hello Hello
	if err := json.Unmarshal(payload.Data, &hello); err != nil {
		return fmt.Errorf("%v", err)
	}

	// Sending a valid Identify Payload triggers the initial handshake with the Discord Gateway.
	// This will result in the Gateway responding with a Ready event.
	// Add a Ready event handler to the bot prior to sending a Heartbeat and Identify Payload.
	//
	// do NOT add multiple Ready event handlers to the bot.
	if len(bot.Handlers.Ready) == 0 {
		if err := bot.Handle(FlagGatewayEventNameReady, func(r *Ready) {
			s.ID = r.SessionID
			// SHARD: set shard information using r.Shard
			bot.ApplicationID = r.Application.ID
		}); err != nil {
			fmt.Println("disconnecting from bot.Handlers.Ready == 0 : line 176")
			return s.disconnectFromConnect(ErrorEventHandler{
				Event:  FlagGatewayEventNameReady,
				Err:    err,
				Action: ErrorEventHandlerAdd,
			})
		}
	}

	// mark the connection as connected.
	s.Connected = make(chan bool)

	// begin listening for payloads.
	//
	// This is done BEFORE sending the first Heartbeat to ensure that
	// the incoming HeartbeatACK is guaranteed to be is handled.
	go bot.listen(s)

	// begin sending heartbeat payloads every heartbeat_interval ms.
	ms := time.Millisecond * time.Duration(hello.HeartbeatInterval)
	fmt.Println(ms)
	s.heartbeat = &heartbeat{
		ticker:   time.NewTicker(ms),
		interval: ms,
		send:     make(chan GatewayPayload),
		respond:  make(chan GatewayPayload),
		ack:      make(chan time.Time, 1),
	}
	go bot.heartbeat(s)

	// send an Opcode 2 Identify to the Discord Gateway.
	identify, err := json.Marshal(Identify{
		Token: bot.Authentication.Token,
		Properties: IdentifyConnectionProperties{
			OS:      runtime.GOOS,
			Browser: module,
			Device:  module,
		},
		Compress:       true,
		LargeThreshold: maxIdentifyLargeThreshold,
		Shard:          nil, // SHARD: set shard information using s.Shard.
		Presence:       *bot.Config.GatewayPresenceUpdate,
		Intents:        bot.Config.Intents,
	})

	if err != nil {
		fmt.Println("disconnecting from sending Identify : line 213")
		return s.disconnectFromConnect(ErrorEvent{
			Event:  FlagGatewayCommandNameIdentify,
			Err:    err,
			Action: ErrorEventActionMarshal,
		})
	}

	fmt.Println(string(identify))

	if err = wsjson.Write(*s.Context, s.Conn, GatewayPayload{
		Op:   FlagGatewayOpcodeIdentify,
		Data: identify,
	}); err != nil {
		return s.disconnectFromConnect(ErrorEvent{
			Event:  FlagGatewayCommandNameIdentify,
			Err:    err,
			Action: ErrorEventActionWrite,
		})
	}

	// send an Opcode 6 Resume to the Discord Gateway to reconnect the session.
	if reconnect {
		resume, err := json.Marshal(Resume{
			Token:     bot.Authentication.Token,
			SessionID: reconnectID,
			Seq:       reconnectSeq,
		})

		fmt.Println(string(resume))

		if err != nil {
			fmt.Println("disconnecting from reconnect : line 249")
			return s.disconnectFromConnect(ErrorEvent{
				Event:  FlagGatewayCommandNameResume,
				Err:    err,
				Action: ErrorEventActionMarshal,
			})
		}

		if err = wsjson.Write(*s.Context, s.Conn, GatewayPayload{
			Op:   FlagGatewayOpcodeIdentify,
			Data: resume,
		}); err != nil {
			return s.disconnectFromConnect(ErrorEvent{
				Event:  FlagGatewayCommandNameResume,
				Err:    err,
				Action: ErrorEventActionWrite,
			})
		}

		// When a reconnection is successful, the Discord Gateway will respond
		// by replaying all missed events in order, finalized by a Resumed event.
		// However, Resumed events do NOT need to be handled.
	}

	return nil
}

// Disconnect disconnects a session from the Discord Gateway using the given status code.
func (s *Session) Disconnect(code int) error {
	if !s.isConnected() {
		return nil
	}

	defer func() {
		s.Connected = nil
	}()

	log.Println("DISCONNECT")

	if err := s.Conn.Close(websocket.StatusCode(code), fmt.Sprintf(gatewayDisconnectMsg, s.ID, code)); err != nil {
		return ErrorDisconnect{
			SessionID: s.ID,
			Err:       err,
			Action:    nil,
		}
	}

	s.cancel()
	return nil
}

// heartbeat continuously sends Opcode 1 Heartbeats to the Discord Gateway to verify the connection is alive.
func (bot *Client) heartbeat(s *Session) {
	// Upon spawning the goroutine, add a HeartbeatACK to the HeartbeatACK channel
	// to prevent the length of the HeartbeatACK channel from being 0 immediately,
	// which results in an attempt to reconnect.
	s.heartbeat.ack <- time.Now().UTC()

	for {
		select {
		case hb := <-s.heartbeat.respond:
			if err := wsjson.Write(*s.Context, s.Conn, hb); err != nil {
				s.mu.connect.Lock()
				defer s.mu.connect.Unlock()

				log.Println("Closing the connection due to a write error...")
				if disconnectErr := s.Disconnect(FlagClientCloseEventCodeNormal); err != nil {
					log.Println(ErrorDisconnect{
						SessionID: s.ID,
						Err:       disconnectErr,
						Action:    ErrorEvent{Event: FlagGatewayCommandNameHeartbeat, Err: err, Action: ErrorEventActionWrite},
					})
				}

				return
			}

			// reset the ticker (and empty existing ticks).
			s.heartbeat.ticker.Reset(s.heartbeat.interval)
			for len(s.heartbeat.ticker.C) > 0 {
				<-s.heartbeat.ticker.C
			}

			// clear queued (outdated) send heartbeats.
			for len(s.heartbeat.send) > 0 {
				<-s.heartbeat.send
			}

			continue

		case hb := <-s.heartbeat.send:
			if err := wsjson.Write(*s.Context, s.Conn, hb); err != nil {
				s.mu.connect.Lock()
				defer s.mu.connect.Unlock()

				log.Println("Closing the connection due to a write error...")
				if disconnectErr := s.Disconnect(FlagClientCloseEventCodeNormal); err != nil {
					log.Println(ErrorDisconnect{
						SessionID: s.ID,
						Err:       disconnectErr,
						Action:    ErrorEvent{Event: FlagGatewayCommandNameHeartbeat, Err: err, Action: ErrorEventActionWrite},
					})
				}

				return
			}

			continue

		// every Heartbeat Interval...
		case <-s.heartbeat.ticker.C:
			// determine if a HeartbeatACK was received from the last sent Heartbeat.
			//
			// close the connection if the last Heartbeat the bot sent never received a HeartbeatACK.
			if len(s.heartbeat.ack) == 0 {
				// close the active connection with a non-1000 and non-1001 close code.
				s.mu.connect.Lock()

				log.Printf("attempting to reconnect session %s due to no HeartbeatACK\n", s.ID)
				if disconnectErr := s.Disconnect(FlagClientCloseEventCodeReconnect); disconnectErr != nil {
					log.Println(ErrorDisconnect{
						SessionID: s.ID,
						Err:       disconnectErr,
						Action:    fmt.Errorf("No HeartbeatACK"),
					})

					s.mu.connect.Unlock()

					return
				}
				s.mu.connect.Unlock()

				// reconnect to the new Discord Gateway Server.
				err := bot.Connect(s)
				if err != nil {
					log.Printf("could not reconnect to session %s due to error: %v\n", s.ID, err)
				}

				return
			}

			// clear the HeartbeatACK channel.
			for len(s.heartbeat.ack) > 0 {
				<-s.heartbeat.ack
			}

			// otherwise, queue a heartbeat.
			seq := make([]byte, 8)
			binary.LittleEndian.PutUint64(seq, uint64(atomic.LoadInt64(&s.Seq)))

			s.heartbeat.send <- GatewayPayload{
				Op:   FlagGatewayOpcodeHeartbeat,
				Data: seq,
			}

			continue

		case <-s.Connected:
			return
		}
	}
}

// listen listens to the connection for payloads from the Discord Gateway.
func (bot *Client) listen(s *Session) {

	for {
		messageType, message, err := s.Conn.Read(*s.Context)

		if err != nil {
			s.mu.connect.Lock()
			defer s.mu.connect.Unlock()

			eventError := ErrorEvent{
				Event:  "Payload",
				Err:    err,
				Action: ErrorEventActionRead,
			}

			log.Println(eventError)
			log.Println("Closing the connection due to a read error...")

			if disconnectErr := s.Disconnect(FlagClientCloseEventCodeNormal); err != nil {
				log.Println(ErrorDisconnect{
					SessionID: s.ID,
					Err:       disconnectErr,
					Action:    eventError,
				})
			}
		}

		if !(len(message) == 0) {
			var reader io.Reader
			reader = bytes.NewBuffer(message)

			if messageType == websocket.MessageBinary {
				readCloser, uncompressErr := zlib.NewReader(reader)
				if uncompressErr != nil {
					log.Println("Error occurred when uncompressing websocket message", uncompressErr)
					if disconnectErr := s.Disconnect(FlagClientCloseEventCodeNormal); err != nil {
						log.Println(ErrorDisconnect{
							SessionID: s.ID,
							Err:       disconnectErr,
							Action:    uncompressErr,
						})
					}
					return
				}

				defer func() {
					zlibCloseErr := readCloser.Close()

					if zlibCloseErr != nil {
						log.Println("Error occured when closing zlib", zlibCloseErr)
						if disconnectErr := s.Disconnect(FlagClientCloseEventCodeNormal); err != nil {
							log.Println(ErrorDisconnect{
								SessionID: s.ID,
								Err:       disconnectErr,
								Action:    zlibCloseErr,
							})
						}
					}
				}()

				reader = readCloser

			}

			var payload GatewayPayload
			decoder := json.NewDecoder(reader)

			if err = decoder.Decode(&payload); err != nil {
				log.Println("error decoding websocket message", err)
				s.mu.connect.Lock()
				defer s.mu.connect.Unlock()

				eventError := ErrorEvent{
					Event:  "Payload",
					Err:    err,
					Action: ErrorEventActionRead,
				}

				log.Println(eventError)
				log.Println("Closing the connection due to a read error...")

				if disconnectErr := s.Disconnect(FlagClientCloseEventCodeNormal); err != nil {
					log.Println(ErrorDisconnect{
						SessionID: s.ID,
						Err:       disconnectErr,
						Action:    eventError,
					})
				}
				return
			}

			log.Println("\tPAYLOAD:", payload)
			bot.onPayload(s, payload)
			return
		}

	}
}

// onPayload handles an Discord Gateway Payload.
func (bot *Client) onPayload(s *Session, payload GatewayPayload) error {

	fmt.Println("HERE ONPAYLOAD")
	// https://discord.com/developers/docs/topics/opcodes-and-status-codes#gateway-gateway-opcodes
	switch payload.Op {

	// run the bot's event handlers.
	case FlagGatewayOpcodeDispatch:
		atomic.StoreInt64(&s.Seq, *payload.SequenceNumber)
		go bot.handle(*payload.EventName, payload.Data)

	// send an Opcode 1 Heartbeat to the Discord Gateway.
	case FlagGatewayOpcodeHeartbeat:
		var heartbeat Heartbeat
		if err := json.Unmarshal(payload.Data, &heartbeat); err != nil {
			return ErrorEvent{
				Event:  FlagGatewayCommandNameHeartbeat,
				Err:    err,
				Action: ErrorEventActionUnmarshal,
			}
		}
		atomic.StoreInt64(&s.Seq, heartbeat.Data)

		s.heartbeat.respond <- payload

	// handle the successful acknowledgement of the client's last heartbeat.
	case FlagGatewayOpcodeHeartbeatACK:
		s.heartbeat.ack <- time.Now().UTC()

	// occurs when the maximum concurrency limit has been reached while connecting,
	// or when the session does NOT reconnect in time.
	case FlagGatewayOpcodeInvalidSession:
		if s.canReconnect() {
			<-time.NewTimer(invalidSessionWaitTime).C
			return fmt.Errorf("Session %s couldn't reconnect to the Discord Gateway", s.ID)
		}

		return fmt.Errorf("Session %s couldn't connect to the Discord Gateway or has invalidated an active session", s.ID)

	// occurs when the Discord Gateway is shutting down the connection,
	// while signalling the client to reconnect.
	case FlagGatewayOpcodeReconnect:
		// close the active connection with a non-1000 and non-1001 close code.
		s.mu.connect.Lock()
		if err := s.Disconnect(FlagGatewayCloseEventCodeSessionTimed.Code); err != nil {
			s.mu.connect.Unlock()
			return ErrorDisconnect{
				SessionID: s.ID,
				Err:       err,
				Action:    errOpcodeReconnect,
			}
		}
		s.mu.connect.Unlock()

		// reconnect to the new Discord Gateway Server.
		return bot.Connect(s)
	}

	return nil
}

// disconnectFromConnect is a helper function for disconnecting from the Connect() func,
// which does NOT require s.mu.connect mutex calls.
//
// err represents the main error that returns if disconnection is SUCCESSFUL.
// err2 represents the format string that returns with the disconnect error and main error if disconnection FAILS.
func (s *Session) disconnectFromConnect(err error) error {
	if disconnectErr := s.Disconnect(FlagClientCloseEventCodeAway); disconnectErr != nil {
		return ErrorDisconnect{
			SessionID: s.ID,
			Err:       disconnectErr,
			Action:    err,
		}
	}

	return err
}
