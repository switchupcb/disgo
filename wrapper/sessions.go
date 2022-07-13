package wrapper

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/goccy/go-json"
	"github.com/switchupcb/disgo/wrapper/internal/socket"

	"nhooyr.io/websocket"
)

const (
	gatewayEndpointParams          = "?v=" + VersionDiscordAPI + "&encoding=json"
	gatewayDisconnectMsg           = "Disconnected Session %q from the Discord Gateway with code %d"
	invalidSessionWaitTime         = 3 * time.Second
	maxIdentifyLargeThreshold      = 250
	FlagClientCloseEventUnexpected = 1011
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
	// Connected is a channel due to its use in a select statements.
	//
	// When Connected is closed (or nil), the Session is NOT connected.
	Connected chan bool

	// Context carries request-scoped data for the Discord Gateway Connection.
	Context context.Context

	// Conn represents a connection to the Discord Gateway.
	Conn *websocket.Conn

	// heartbeat contains the fields required to implement the heartbeat mechanism.
	heartbeat *heartbeat

	// mu represents a mutex that is used to protect the Session's variables from data-races.
	mu sync.RWMutex

	// muConnect ensures that the Connect and Disconnect functionality is only running on one goroutine per Session.
	//
	// Prevents a theoretical race condition where the main thread starts a heartbeat() goroutine in Connect(),
	// but fails to receive an ACK (within the FailedHeartbeatInterval) BEFORE Connect() returns,
	// resulting in the session connecting while also attempting to reconnect.
	//
	// Prevents a theoretical race condition where the main thread attempts to Connect(),
	// while another thread calls s.Disconnect() resulting in undefined behavior while connecting
	// to the Discord Gateway.
	muConnect sync.Mutex

	// muHeartbeat represents the mutex for heartbeat functionality.
	//
	// Prevents race conditions where a HeartbeatACK is written while a HeartbeatACK is being read or cleared.
	//
	// Used to process heartbeat operations as a transaction.
	muHeartbeat sync.Mutex
}

// heartbeat represents the heartbeat mechanism for a session.
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

// isConnected returns whether the session is connected.
func (s *Session) isConnected() bool {
	select {
	case <-s.Connected:
		return false
	default:
		return s.Connected != nil
	}
}

// canReconnect determines whether the session is in a valid state to reconnect.
func (s *Session) canReconnect() bool {
	return s.ID != "" && atomic.LoadInt64(&s.Seq) != 0
}

// Connect creates or reestablishes a session's open connection to the Discord Gateway.
func (bot *Client) Connect(s *Session) error {
	s.muConnect.Lock()
	defer s.muConnect.Unlock()

	if s.isConnected() {
		return nil
	}

	// request a valid Gateway URL endpoint from the Discord API.
	if s.Endpoint == "" || !s.canReconnect() {
		gateway := GetGateway{}
		response, err := gateway.Send(bot)
		if err != nil {
			return fmt.Errorf("an error occurred getting the Gateway API Endpoint\n%w", err)
		}

		s.Endpoint = response.URL + gatewayEndpointParams
	}

	// connect to the Discord Gateway Websocket.
	var err error
	s.Context = context.Background()
	if s.Conn, _, err = websocket.Dial(s.Context, s.Endpoint, nil); err != nil {
		return fmt.Errorf("an error occurred while connecting to the Discord Gateway\n%w", err)
	}

	// handle the incoming Hello event upon connecting to the Gateway.
	var hello Hello
	if err := readEvent(s, FlagGatewayEventNameHello, &hello); err != nil {
		return s.disconnectFromConnect(err)
	}

	// Sending a valid Identify Payload triggers the initial handshake with the Discord Gateway.
	// This will result in the Gateway responding with a Ready event.
	// Add a Ready event handler to the bot prior to sending a Heartbeat and Identify Payload.
	//
	// do NOT add multiple Ready event handlers to the bot.
	if len(bot.Handlers.Ready) == 0 {
		if err := bot.Handle(FlagGatewayEventNameReady, func(r *Ready) {
			s.mu.Lock()
			s.ID = r.SessionID
			s.mu.Unlock()
			// SHARD: set shard information using r.Shard
			bot.ApplicationID = r.Application.ID
		}); err != nil {
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
	// the any incoming HeartbeatACK is guaranteed to be is handled.
	go bot.listen(s)

	// begin sending heartbeat payloads every heartbeat_interval ms.
	ms := time.Millisecond * time.Duration(hello.HeartbeatInterval)
	s.muHeartbeat.Lock()
	s.heartbeat = &heartbeat{
		interval: ms,
		ticker:   time.NewTicker(ms),
		send:     make(chan Heartbeat),

		// add a HeartbeatACK to the HeartbeatACK channel to prevent
		// the length of the HeartbeatACK channel from being 0 immediately,
		// which results in an attempt to reconnect.
		acks: 1,
	}
	s.muHeartbeat.Unlock()
	go bot.pulse(s)
	go bot.heartbeat(s)

	return bot.initial(s)
}

// initial sends the initial identify or resume packet required to connect to the Gateway.
func (bot *Client) initial(s *Session) error {
	if !s.canReconnect() {
		// send an Opcode 2 Identify to the Discord Gateway.
		if err := writeEvent(s, FlagGatewayOpcodeIdentify, FlagGatewayCommandNameIdentify,
			Identify{
				Token: bot.Authentication.Token,
				Properties: IdentifyConnectionProperties{
					OS:      runtime.GOOS,
					Browser: module,
					Device:  module,
				},
				Compress:       true,
				LargeThreshold: maxIdentifyLargeThreshold,
				Shard:          nil, // SHARD: set shard information using s.Shard.
				Presence:       *bot.Config.Gateway.GatewayPresenceUpdate,
				Intents:        bot.Config.Gateway.Intents,
			}); err != nil {
			return s.disconnectFromConnect(err)
		}
	} else {
		// send an Opcode 6 Resume to the Discord Gateway to reconnect the session.
		if err := writeEvent(s, FlagGatewayOpcodeResume, FlagGatewayCommandNameResume,
			Resume{
				Token:     bot.Authentication.Token,
				SessionID: s.ID,
				Seq:       atomic.LoadInt64(&s.Seq),
			}); err != nil {
			return s.disconnectFromConnect(err)
		}

		// When a reconnection is successful, the Discord Gateway will respond
		// by replaying all missed events in order, finalized by a Resumed event.
		// However, Resumed events do NOT need to be handled.
	}

	return nil
}

// Reconnect reconnects an already connected session to the Discord Gateway
// by disconnecting the session, then connecting again.
func (bot *Client) Reconnect(s *Session) error {
	if err := s.Disconnect(FlagClientCloseEventUnexpected); err != nil {
		return ErrorDisconnect{
			SessionID: s.ID,
			Err:       err,
			Action:    errOpcodeReconnect,
		}
	}

	return bot.Connect(s)
}

// Disconnect disconnects a session from the Discord Gateway using the given status code.
func (s *Session) Disconnect(code int) error {
	s.muConnect.Lock()
	defer s.muConnect.Unlock()

	if !s.isConnected() {
		return nil
	}

	// close the Connected channel to kill the goroutines of the Session.
	close(s.Connected)

	// close the connection.
	if err := s.Conn.Close(websocket.StatusCode(code), fmt.Sprintf(gatewayDisconnectMsg, s.ID, code)); err != nil {
		return ErrorDisconnect{
			SessionID: s.ID,
			Err:       err,
			Action:    nil,
		}
	}

	return nil
}

// disconnectFromConnect is a helper function for disconnecting from the Connect() func,
// which does NOT require s.mu.connect mutex calls.
//
// err represents the main error that returns if disconnection is SUCCESSFUL.
func (s *Session) disconnectFromConnect(err error) error {
	if disconnectErr := s.Disconnect(FlagClientCloseEventCodeNormal); disconnectErr != nil {
		return ErrorDisconnect{
			SessionID: s.ID,
			Err:       disconnectErr,
			Action:    err,
		}
	}

	return err
}

// disconnectFromRoutine is a helper function for disconnecting from a non-main goroutine,
// which requires s.mu.connect mutex calls and logging.
//
// err represents the main error that returns if disconnection is SUCCESSFUL.
func (s *Session) disconnectFromRoutine(reason string, err error) {
	log.Println(reason)
	if disconnectErr := s.Disconnect(FlagClientCloseEventCodeNormal); disconnectErr != nil {
		err = ErrorDisconnect{
			SessionID: s.ID,
			Err:       disconnectErr,
			Action:    err,
		}
	}

	log.Println(err)
}

// listen listens to the connection for payloads from the Discord Gateway.
func (bot *Client) listen(s *Session) {
	for {
		payload := getPayload()
		if err := socket.Read(s.Context, s.Conn, payload); err != nil {
			select {
			case <-s.Connected:
				return
			default:
				s.disconnectFromRoutine("Closing the connection due to a read error...", ErrorEvent{
					Event:  "Payload",
					Err:    err,
					Action: ErrorEventActionRead,
				})
			}
		}

		log.Println("PAYLOAD", payload.Op, string(payload.Data))
		if err := bot.onPayload(s, *payload); err != nil {
			s.disconnectFromRoutine("onPayload error", err)

			return
		}
	}
}

// onPayload handles an Discord Gateway Payload.
func (bot *Client) onPayload(s *Session, payload GatewayPayload) error {
	defer putPayload(&payload)

	// https://discord.com/developers/docs/topics/opcodes-and-status-codes#gateway-gateway-opcodes
	switch payload.Op {
	// run the bot's event handlers.
	case FlagGatewayOpcodeDispatch:
		atomic.StoreInt64(&s.Seq, *payload.SequenceNumber)
		go bot.handle(*payload.EventName, payload.Data)

	// send an Opcode 1 Heartbeat to the Discord Gateway.
	case FlagGatewayOpcodeHeartbeat:
		go bot.respond(s, payload.Data)

	// handle the successful acknowledgement of the client's last heartbeat.
	case FlagGatewayOpcodeHeartbeatACK:
		s.muHeartbeat.Lock()
		atomic.AddUint32(&s.heartbeat.acks, 1)
		s.muHeartbeat.Unlock()

	// occurs when the maximum concurrency limit has been reached while connecting,
	// or when the session does NOT reconnect in time.
	case FlagGatewayOpcodeInvalidSession:
		// wait for Discord to close the session, then complete a fresh connect.
		<-time.NewTimer(invalidSessionWaitTime).C

		s.ID = ""
		s.Seq = 0
		if err := bot.initial(s); err != nil {
			return err
		}

		if !s.isConnected() {
			return fmt.Errorf("Session %q couldn't connect to the Discord Gateway or has invalidated an active session", s.ID)
		}

	// occurs when the Discord Gateway is shutting down the connection,
	// while signalling the client to reconnect.
	case FlagGatewayOpcodeReconnect:
		if err := bot.Reconnect(s); err != nil {
			return ErrorDisconnect{
				SessionID: s.ID,
				Err:       err,
				Action:    errOpcodeReconnect,
			}
		}
	}

	return nil
}

// heartbeat listens for pulses to send Opcode 1 Heartbeats to the Discord Gateway (to verify the connection is alive).
func (bot *Client) heartbeat(s *Session) {
	for {
		select {
		case hb := <-s.heartbeat.send:
			s.muHeartbeat.Lock()

			// close the connection if the last sent Heartbeat never received a HeartbeatACK.
			if atomic.LoadUint32(&s.heartbeat.acks) == 0 {
				s.muHeartbeat.Unlock()

				log.Printf("attempting to reconnect session %q due to no HeartbeatACK", s.ID)
				if err := bot.Reconnect(s); err != nil {
					log.Println(ErrorDisconnect{
						SessionID: s.ID,
						Err:       err,
						Action:    fmt.Errorf("no HeartbeatACK"),
					})
				}

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
				s.disconnectFromRoutine("Closing the connection due to a write error...", err)
				s.muHeartbeat.Unlock()

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

			s.muHeartbeat.Unlock()

		case <-s.Connected:
			return
		}
	}
}

// respond responds to Opcode 1 Heartbeats from the Discord Gateway.
func (bot *Client) respond(s *Session, data json.RawMessage) {
	var heartbeat Heartbeat
	if err := json.Unmarshal(data, &heartbeat); err != nil {
		s.disconnectFromRoutine("Closing the connection due to an unmarshal error...", ErrorEvent{
			Event:  FlagGatewayCommandNameHeartbeat,
			Err:    err,
			Action: ErrorEventActionUnmarshal,
		})
	}

	atomic.StoreInt64(&s.Seq, heartbeat.Data)

	s.muHeartbeat.Lock()

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

	s.muHeartbeat.Unlock()
}

// pulse generates Opcode 1 Heartbeats for a Session's heartbeat channel.
func (bot *Client) pulse(s *Session) {
	for {
		s.muHeartbeat.Lock()

		select {
		default:
			s.muHeartbeat.Unlock()

		// every Heartbeat Interval...
		case <-s.heartbeat.ticker.C:

			// queue a heartbeat.
			s.heartbeat.send <- Heartbeat{Data: s.Seq}

			log.Println("queued heartbeat")

			s.muHeartbeat.Unlock()

		case <-s.Connected:
			s.muHeartbeat.Unlock()

			return
		}
	}
}

// readEvent is a helper function for reading events from the websocket session.
func readEvent(s *Session, name string, dst any) error {
	payload := new(GatewayPayload)
	if err := socket.Read(s.Context, s.Conn, payload); err != nil {
		return ErrorEvent{
			Event:  name,
			Err:    err,
			Action: ErrorEventActionRead,
		}
	}

	if err := json.Unmarshal(payload.Data, dst); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

// writeEvent is a helper function for writing events to the websocket session.
// returns an ErrorEvent.
func writeEvent(s *Session, op int, name string, dst any) error {
	event, err := json.Marshal(dst)
	if err != nil {
		return ErrorEvent{
			Event:  name,
			Err:    err,
			Action: ErrorEventActionMarshal,
		}
	}

	if err = socket.Write(s.Context, s.Conn, websocket.MessageBinary,
		GatewayPayload{ //nolint:exhaustruct
			Op:   op,
			Data: event,
		}); err != nil {
		return ErrorEvent{
			Event:  name,
			Err:    err,
			Action: ErrorEventActionWrite,
		}
	}

	return nil
}
