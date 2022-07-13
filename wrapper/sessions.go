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
	gatewayEndpointParams     = "?v=" + VersionDiscordAPI + "&encoding=json"
	gatewayDisconnectMsg      = "Disconnected Session %s from the Discord Gateway with code %d"
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
	Context context.Context

	// Conn represents a connection to the Discord Gateway.
	Conn *websocket.Conn

	// heartbeat contains the fields required to implement the heartbeat mechanism.
	heartbeat *heartbeat

	// mu represents the session mutex of this session.
	mu *sessionMutex
}

// heartbeat represents the heartbeat mechanism for a session.
type heartbeat struct {
	// ticker is a timer used to time the interval between each Heartbeat Payload.
	ticker *time.Ticker

	// interval represents the interval of time between each Heartbeat Payload.
	interval time.Duration

	// send represents a channel of heartbeats that will be sent to the Discord Gateway.
	send chan Heartbeat

	// respond represents a channel of heartbeats in response to an Opcode 1 Heartbeat
	// that will be sent to the Discord Gateway.
	//
	// respond prevents a theoretical race condition where the ticker queues a Heartbeat,
	// while the listen thread (onPayload) queues a Heartbeat (in response to the Discord Gateway),
	// resulting in two Heartbeat Payloads being sent to the Discord Gateway consecutively within nanoseconds,
	// while the first Heartbeat Payload is outdated.
	respond chan Heartbeat

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

	// resource represents the mutex for session resources.
	resource sync.RWMutex

	// heartbeat represents the mutex for heartbeat functionality.
	//
	// Prevents race conditions where a HeartbeatACK is written while a HeartbeatACK is being read or cleared.
	//
	// Used to process heartbeat operations as a transaction.
	heartbeat sync.Mutex
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

	s.Endpoint = response.URL + gatewayEndpointParams

	// connect to the Discord Gateway Websocket.
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
			s.mu.resource.Lock()
			s.ID = r.SessionID
			s.mu.resource.Unlock()
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
	s.mu.heartbeat.Lock()
	s.heartbeat = &heartbeat{
		ticker:   time.NewTicker(ms),
		interval: ms,
		send:     make(chan Heartbeat, 1),
		respond:  make(chan Heartbeat, 1),
		ack:      make(chan time.Time, 1),
	}
	s.mu.heartbeat.Unlock()
	go bot.heartbeat(s)

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

	// send an Opcode 6 Resume to the Discord Gateway to reconnect the session.
	if reconnect {
		if err := writeEvent(s, FlagGatewayOpcodeResume, FlagGatewayCommandNameResume,
			Resume{
				Token:     bot.Authentication.Token,
				SessionID: reconnectID,
				Seq:       reconnectSeq,
			}); err != nil {
			return s.disconnectFromConnect(err)
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
		close(s.Connected)
	}()

	if err := s.Conn.Close(websocket.StatusCode(code), fmt.Sprintf(gatewayDisconnectMsg, s.ID, code)); err != nil {
		return ErrorDisconnect{
			SessionID: s.ID,
			Err:       err,
			Action:    nil,
		}
	}

	return nil
}

// heartbeat continuously sends Opcode 1 Heartbeats to the Discord Gateway to verify the connection is alive.
func (bot *Client) heartbeat(s *Session) {
	// Upon spawning the goroutine, add a HeartbeatACK to the HeartbeatACK channel
	// to prevent the length of the HeartbeatACK channel from being 0 immediately,
	// which results in an attempt to reconnect.
	s.heartbeat.ack <- time.Now().UTC()

	for {
		s.mu.heartbeat.Lock()

		select {
		default:
			s.mu.heartbeat.Unlock()

		case hb := <-s.heartbeat.respond:
			if err := writeEvent(s, FlagGatewayOpcodeHeartbeat, FlagGatewayCommandNameHeartbeat, hb); err != nil {
				s.disconnectFromRoutine("Closing the connection due to a write error...", err)

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

			log.Println("responded to heartbeat")

			s.mu.heartbeat.Unlock()

			continue

		case hb := <-s.heartbeat.send:
			if err := writeEvent(s, FlagGatewayOpcodeHeartbeat, FlagGatewayCommandNameHeartbeat, hb); err != nil {
				s.disconnectFromRoutine("Closing the connection due to a write error...", err)

				return
			}

			log.Println("sent heartbeat")

			s.mu.heartbeat.Unlock()

			continue

		// every Heartbeat Interval...
		case <-s.heartbeat.ticker.C:
			// determine if a HeartbeatACK was received from the last sent Heartbeat.
			//
			// close the connection if the last Heartbeat the bot sent never received a HeartbeatACK.
			if len(s.heartbeat.ack) == 0 {
				// close the active connection with a non-1000 and non-1001 close code.
				s.mu.connect.Lock()

				log.Printf("attempting to reconnect session %s due to no HeartbeatACK", s.ID)
				if disconnectErr := s.Disconnect(FlagClientCloseEventCodeReconnect); disconnectErr != nil {
					log.Println(ErrorDisconnect{
						SessionID: s.ID,
						Err:       disconnectErr,
						Action:    fmt.Errorf("no HeartbeatACK"),
					})

					s.mu.heartbeat.Unlock()
					s.mu.connect.Unlock()

					return
				}

				s.mu.heartbeat.Unlock()
				s.mu.connect.Unlock()

				// reconnect to the new Discord Gateway Server.
				if err := bot.Connect(s); err != nil {
					log.Printf("could not reconnect to session %s due to error: %v", s.ID, err)
				}

				return
			}

			// clear the HeartbeatACK channel.
			for len(s.heartbeat.ack) > 0 {
				<-s.heartbeat.ack
			}

			// queue a heartbeat.
			s.heartbeat.send <- Heartbeat{Data: s.Seq}

			s.mu.heartbeat.Unlock()

			continue

		case <-s.Connected:
			s.mu.heartbeat.Unlock()

			return
		}
	}
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
				eventError := ErrorEvent{
					Event:  "Payload",
					Err:    err,
					Action: ErrorEventActionRead,
				}

				s.disconnectFromRoutine("Closing the connection due to a read error...", eventError)
			}
		}

		fmt.Println("PAYLOAD", payload.Op, string(payload.Data))
		if err := bot.onPayload(s, *payload); err != nil {
			s.disconnectFromRoutine("onPayload error", err)
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
		s.mu.heartbeat.Lock()

		var heartbeat Heartbeat
		if err := json.Unmarshal(payload.Data, &heartbeat); err != nil {
			return ErrorEvent{
				Event:  FlagGatewayCommandNameHeartbeat,
				Err:    err,
				Action: ErrorEventActionUnmarshal,
			}
		}

		atomic.StoreInt64(&s.Seq, heartbeat.Data)

		s.heartbeat.respond <- heartbeat

		s.mu.heartbeat.Unlock()

	// handle the successful acknowledgement of the client's last heartbeat.
	case FlagGatewayOpcodeHeartbeatACK:
		s.mu.heartbeat.Lock()
		s.heartbeat.ack <- time.Now().UTC()
		s.mu.heartbeat.Unlock()

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

var (
	// gpool represents a synchronized Gateway Payload pool.
	gpool sync.Pool
)

// getPayload gets a Gateway Payload from the pool.
func getPayload() *GatewayPayload {
	if g := gpool.Get(); g != nil {
		return g.(*GatewayPayload)
	}

	return new(GatewayPayload)
}

// putPayload puts a Gateway Payload into the pool.
func putPayload(g *GatewayPayload) {
	// reset the Gateway Payload.
	g.Op = 0
	g.Data = nil
	g.SequenceNumber = nil
	g.EventName = nil
	gpool.Put(g)
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
		return fmt.Errorf("%v", err)
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
		GatewayPayload{
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

// disconnectFromConnect is a helper function for disconnecting from the Connect() func,
// which does NOT require s.mu.connect mutex calls.
//
// err represents the main error that returns if disconnection is SUCCESSFUL.
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

// disconnectFromRoutine is a helper function for disconnecting from a non-main goroutine,
// which does requires s.mu.connect mutex calls and logging.
//
// err represents the main error that returns if disconnection is SUCCESSFUL.
func (s *Session) disconnectFromRoutine(msg string, err error) {
	s.mu.connect.Lock()
	defer s.mu.connect.Unlock()

	log.Println(msg)
	if disconnectErr := s.Disconnect(FlagClientCloseEventCodeNormal); disconnectErr != nil {
		log.Println(ErrorDisconnect{
			SessionID: s.ID,
			Err:       disconnectErr,
			Action:    err,
		}.Error())
	} else {
		log.Println(err.Error())
	}
}
