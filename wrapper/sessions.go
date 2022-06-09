package wrapper

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

// TODO: fix error messages
// TODO: fix opcode include in json.
// TODO: ensure context is correct with regards to Mutex and Resource Contention.
// TODO: change heartbeat thread algorithm to use channels instead of mutex.

const (
	module                    = "github.com/switchupcb/disgo"
	gatewayEncoding           = "json"
	maxIdentifyLargeThreshold = 250
	gatewayDisconnectCode     = 1000
	gatewayDisconnectMsg      = "Disconnected Session %s from the Discord Gateway with code %d"

	// allowedFailedHeartbeats represents the allowed amount of failed heartbeats before a connection
	// is considered disconnected.
	//
	// Discord documentation advises clients to sent Opcode 1 Heartbeat Payloads AFTER the HeartbeatInterval * [0,1],
	// AFTER receiving the Opcode 10 Hello. Disgo sends the HeartbeatInterval as soon as possible once Hello has been
	// received (limit of HeartbeatInterval * 0), and once every HeartbeatInterval after.
	//
	// As a result, time.Now() - LastHeartbeatACK (which represents the duration since the last HeartbeatACK)
	// should always be less than or equal to (<=) a single HeartbeatInterval.
	allowedFailedHeartbeats = 2
	invalidSessionWaitTime  = 3 * time.Second
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

	// Ticker is a timer used to time the interval between each Heartbeat Payload.
	Ticker *time.Ticker

	// HeartbeatInterval represents the interval of time between each Heartbeat Payload.
	HeartbeatInterval time.Duration

	// FailedHeartbeatInterval represents the duration of time before the Session is considered
	// disconnected due to a lack of HeartbeatACK Payloads (prompted by Heartbeat Payloads).
	FailedHeartbeatInterval time.Duration

	// LastHeartbeatACK represents the time when the last HeartbeatACK was received.
	LastHeartbeatACK time.Time

	// mu represents the session mutex of this session.
	mu *sessionMutex
}

// sessionMutex represents a Session's mutexes used to prevent race conditions while
// reading and writing from a Session's fields.
type sessionMutex struct {
	// connect ensures that the Connect() function is only running on one goroutine per Session.
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

	// heartbeat ensures that an Opcode 1 Heartbeat is only being sent once per Heartbeat Interval.
	//
	// Prevents a theoretical race condition where a heartbeat thread is scheduled to send a Heartbeat (via Ticker),
	// while a listen thread (onPayload) responds to Discord sending an Opcode 1 Heartbeat to the client,
	// resulting in the session writing two Heartbeat Payloads to the Discord Gateway at the same time.
	//
	// It is currently still possible for two heartbeats to be sent consecutively within nanoseconds.
	heartbeat sync.RWMutex
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
		return fmt.Errorf("an error occurred getting the Gateway API URL\n%w", err)
	}

	// TODO: zlib compression
	s.Endpoint = response.URL + "?v=" + VersionDiscordAPI + "&encoding=" + gatewayEncoding

	// connect to the Discord Gateway Websocket.
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	s.Context = &ctx
	s.Conn, _, err = websocket.Dial(*s.Context, s.Endpoint, nil)
	if err != nil {
		return fmt.Errorf("an error occurred connecting to the Discord Gateway\n%w", err)
	}

	defer s.Conn.Close(websocket.StatusInternalError, "StatusInternalError")

	// handle the incoming Hello event upon connecting to the Gateway.
	var hello Hello
	err = wsjson.Read(ctx, s.Conn, &hello)
	if err != nil {
		return fmt.Errorf(ErrEventRead, FlagGatewayEventNameHello, err)
	}

	// mark the connection as connected.
	s.Connected = make(chan bool)

	// begin sending heartbeat payloads every heartbeat_interval ms.
	s.Ticker = time.NewTicker(s.HeartbeatInterval)
	s.HeartbeatInterval = hello.HeartbeatInterval * time.Millisecond
	s.FailedHeartbeatInterval = hello.HeartbeatInterval * time.Millisecond * allowedFailedHeartbeats
	s.LastHeartbeatACK = time.Now().UTC()
	go bot.heartbeat(s)

	// begin listening for events.
	go bot.listen(s)

	// Sending a valid Identify Payload triggers the initial handshake with the Discord Gateway.
	// This will result in the Gateway responding with a Ready event.
	// Add a Ready event handler to the bot prior to sending the Identify Payload.
	//
	// do NOT add multiple Ready event handlers to the bot.
	if len(bot.Handlers.Ready) == 0 {
		bot.Handle(FlagGatewayEventNameReady, func(r *Ready) {
			s.ID = r.SessionID
			// SHARD: set shard information using r.Shard
			bot.ApplicationID = r.Application.ID
		})
	}

	// send an Opcode 2 Identify to the Discord Gateway.
	identify := Identify{
		Token: bot.Authentication.Token,
		Properties: IdentifyConnectionProperties{
			OS:      runtime.GOOS,
			Browser: module,
			Device:  module,
		},
		Compress:       true, // TODO: account for compression
		LargeThreshold: maxIdentifyLargeThreshold,
		Shard:          nil, // SHARD: set shard information using s.Shard.
		Presence:       *bot.Config.GatewayPresenceUpdate,
		Intents:        bot.Config.Intents,
	}

	err = wsjson.Write(*s.Context, s.Conn, identify)
	if err != nil {
		s.Disconnect(gatewayDisconnectCode)
		return fmt.Errorf(ErrEventWrite, "Identify", err)
	}

	// send an Opcode 6 Resume to the Discord Gateway to reconnect the session.
	if reconnect {
		resume := Resume{
			Token:     bot.Authentication.Token,
			SessionID: reconnectID,
			Seq:       reconnectSeq,
		}

		err = wsjson.Write(*s.Context, s.Conn, resume)
		if err != nil {
			s.Disconnect(gatewayDisconnectCode)
			return fmt.Errorf(ErrEventWrite, "Resume", err)
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

	err := s.Conn.Close(websocket.StatusCode(code), fmt.Sprintf(gatewayDisconnectMsg, s.ID, code))
	if err != nil {
		return fmt.Errorf(ErrDisconnecting, s.ID)
	}

	return nil
}

// heartbeat continuously sends Opcode 1 Heartbeats to the Discord Gateway to verify the connection is alive.
func (bot *Client) heartbeat(s *Session) {
	var hb Heartbeat
	hb.Op = FlagGatewayOpcodeHeartbeat

	for {
		// close the connection if the last two heartbeats were NOT acknowledged.
		if time.Now().UTC().Sub(s.LastHeartbeatACK) > s.FailedHeartbeatInterval {

			// close the active connection with a non-1000 and non-1001 close code.
			s.mu.connect.Lock()
			log.Printf("attempting to reconnect session %s", s.ID)
			s.Disconnect(FlagGatewayCloseEventCodeSessionTimed.Code)
			s.mu.connect.Unlock()

			// reconnect to the new Discord Gateway Server.
			err := bot.Connect(s)
			if err != nil {
				log.Printf("could not reconnect to session %s due to error: %v", s.ID, err)
			}

			return
		}

		// send an Opcode 1 Heartbeat Payload.
		s.mu.heartbeat.Lock()
		*hb.Data = atomic.LoadInt64(&s.Seq)
		err := wsjson.Write(*s.Context, s.Conn, hb)
		s.mu.heartbeat.Unlock()
		if err != nil {
			s.mu.connect.Lock()
			log.Printf("an error occurred writing a heartbeat: %v\nclosing the connection...", err)
			s.Disconnect(gatewayDisconnectCode)
			s.mu.connect.Unlock()
			return
		}

		select {
		case <-s.Ticker.C:
		case <-s.Connected:
			return
		}
	}
}

// listen listens to the connection for payloads from the Discord Gateway.
func (bot *Client) listen(s *Session) {
	for {
		var payload []byte
		err := wsjson.Read(*s.Context, s.Conn, payload)
		if err != nil {
			log.Printf("%v", err)
			return
		}

		select {
		case <-s.Connected:
			return
		default:
			bot.onPayload(s, payload)
		}
	}
}

// onPayload handles an Discord Gateway Payload.
func (bot *Client) onPayload(s *Session, data []byte) error {
	var event GatewayPayload
	err := json.Unmarshal(data, &event)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	// https://discord.com/developers/docs/topics/opcodes-and-status-codes#gateway-gateway-opcodes
	switch *event.Op {

	// run the bot's event handlers.
	case FlagGatewayOpcodeDispatch:
		atomic.StoreInt64(&s.Seq, event.SequenceNumber)
		go bot.handle(event.EventName, event.Data)

	// send an Opcode 1 Heartbeat to the Discord Gateway.
	case FlagGatewayOpcodeHeartbeat:
		s.mu.heartbeat.Lock()
		defer s.mu.heartbeat.Unlock()

		var heartbeat Heartbeat
		err := json.Unmarshal(event.Data, &heartbeat)
		if err != nil {
			return fmt.Errorf("%w", err)
		}

		// send the heartbeat back.
		heartbeat.Op = FlagGatewayOpcodeHeartbeat
		err = wsjson.Write(*s.Context, s.Conn, heartbeat)
		if err != nil {
			log.Printf("%v", err)
		}

		s.Ticker.Reset(s.HeartbeatInterval)

	// handle the successful acknowledgement of the client's last heartbeat.
	case FlagGatewayOpcodeHeartbeatACK:
		s.LastHeartbeatACK = time.Now().UTC()

	// occurs when the maximum concurrency limit has been reached while connecting,
	// or when the session does NOT reconnect in time.
	case FlagGatewayOpcodeInvalidSession:
		if s.canReconnect() {
			<-time.NewTimer(invalidSessionWaitTime).C
			return fmt.Errorf("the session could not reconnect to the Discord Gateway")
		}

		return fmt.Errorf("the session could not connect to the Discord Gateway " +
			"or has invalidated an active session")

	// occurs when the Discord Gateway is shutting down the connection,
	// while signalling the client to reconnect.
	case FlagGatewayOpcodeReconnect:
		// close the active connection with a non-1000 and non-1001 close code.
		s.mu.connect.Lock()
		s.Disconnect(FlagGatewayCloseEventCodeSessionTimed.Code)
		s.mu.connect.Unlock()

		// reconnect to the new Discord Gateway Server.
		return bot.Connect(s)
	}

	return nil
}
