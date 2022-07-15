package wrapper

import (
	"context"
	"errors"
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
	invalidSessionWaitTime    = 1 * time.Second
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

	// Context carries request-scoped data for the Discord Gateway Connection.
	//
	// Context is also used as a signal for the Session's goroutines.
	Context context.Context

	// Conn represents a connection to the Discord Gateway.
	Conn *websocket.Conn

	// cancel represents the cancellation signal for context.
	cancel context.CancelFunc

	// heartbeat contains the fields required to implement the heartbeat mechanism.
	heartbeat *heartbeat

	// a mutex is used to protect the Session's variables from data races
	// by providing transactional functionality.
	sync.RWMutex
}

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

// isConnected returns whether the session is connected.
func (s *Session) isConnected() bool {
	if s.Context == nil {
		return false
	}

	select {
	case <-s.Context.Done():
		return false
	default:
		return true
	}
}

// canReconnect determines whether the session is in a valid state to reconnect.
func (s *Session) canReconnect() bool {
	return s.ID != "" && atomic.LoadInt64(&s.Seq) != 0
}

// Connect connects a session to the Discord Gateway (WebSocket Connection).
func (s *Session) Connect(bot *Client) error {
	s.Lock()
	defer s.Unlock()

	log.Printf("connecting Session %q", s.ID)

	return s.connect(bot)
}

// connect connects a session to a WebSocket Connection.
func (s *Session) connect(bot *Client) error {
	if s.isConnected() {
		return fmt.Errorf("Session %q is already connected", s.ID)
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

	var err error

	// connect to the Discord Gateway Websocket.
	s.Context, s.cancel = context.WithCancel(context.Background())
	if s.Conn, _, err = websocket.Dial(s.Context, s.Endpoint, nil); err != nil {
		return fmt.Errorf("an error occurred while connecting to the Discord Gateway\n%w", err)
	}

	// handle the incoming Hello event upon connecting to the Gateway.
	var hello Hello
	if err := readEvent(s, FlagGatewayEventNameHello, &hello); err != nil {
		if disconnectErr := s.disconnect(FlagClientCloseEventCodeNormal); disconnectErr != nil {
			return ErrorDisconnect{
				SessionID: s.ID,
				Err:       disconnectErr,
				Action:    err,
			}
		}

		return err
	}

	// begin listening for payloads.
	//
	// This is done BEFORE sending the first Heartbeat to ensure that
	// any incoming payloads (Ready, HeartbeatACK) are guaranteed to be handled.
	go s.listen(bot)

	// begin sending heartbeat payloads every heartbeat_interval ms.
	ms := time.Millisecond * time.Duration(hello.HeartbeatInterval)
	s.heartbeat = &heartbeat{
		interval: ms,
		ticker:   time.NewTicker(ms),
		send:     make(chan Heartbeat),

		// add a HeartbeatACK to the HeartbeatACK channel to prevent
		// the length of the HeartbeatACK channel from being 0 immediately,
		// which results in an attempt to reconnect.
		acks: 1,
	}
	go s.pulse()
	go s.beat(bot)

	if err := s.initial(bot); err != nil {
		if disconnectErr := s.disconnect(FlagClientCloseEventCodeNormal); disconnectErr != nil {
			return ErrorDisconnect{
				SessionID: s.ID,
				Err:       disconnectErr,
				Action:    err,
			}
		}

		return err
	}

	return nil
}

// initial sends the initial identify or resume packet required to connect to the Gateway.
func (s *Session) initial(bot *Client) error {
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
			return err
		}
	} else {
		// send an Opcode 6 Resume to the Discord Gateway to reconnect the session.
		if err := writeEvent(s, FlagGatewayOpcodeResume, FlagGatewayCommandNameResume,
			Resume{
				Token:     bot.Authentication.Token,
				SessionID: s.ID,
				Seq:       atomic.LoadInt64(&s.Seq),
			}); err != nil {
			return err
		}

		// When a reconnection is successful, the Discord Gateway will respond
		// by replaying all missed events in order, finalized by a Resumed event.
		// However, Resumed events do NOT need to be handled.
	}

	return nil
}

// Disconnect disconnects a session from the Discord Gateway using the given status code.
func (s *Session) Disconnect() error {
	s.Lock()
	defer s.Unlock()

	if err := s.disconnect(FlagClientCloseEventCodeNormal); err != nil {
		return err
	}

	log.Printf("disconnected Session %q with code %d", s.ID, FlagClientCloseEventCodeNormal)

	return nil
}

// disconnect disconnects a session from a WebSocket Connection using the given status code.
func (s *Session) disconnect(code int) error {
	if !s.isConnected() {
		return fmt.Errorf("Session %q is already disconnected", s.ID)
	}

	if err := s.Conn.Close(websocket.StatusCode(code), ""); err != nil {
		return ErrorDisconnect{
			SessionID: s.ID,
			Err:       err,
			Action:    nil,
		}
	}

	// cancel the context to kill the goroutines of the Session.
	s.cancel()

	return nil
}

// disconnectFromRoutine is a helper function for disconnecting from a non-main goroutine.
func (s *Session) disconnectFromRoutine(reason string, err error) {
	log.Println(reason)
	if disconnectErr := s.disconnect(FlagClientCloseEventCodeAway); disconnectErr != nil {
		err = ErrorDisconnect{
			SessionID: s.ID,
			Err:       disconnectErr,
			Action:    err,
		}
	}

	log.Println(err)
}

// Reconnect reconnects an already connected session to the Discord Gateway
// by disconnecting the session, then connecting again.
func (s *Session) Reconnect(bot *Client) error {
	s.Lock()
	defer s.Unlock()

	log.Printf("reconnecting Session %q", s.ID)

	return s.reconnect(bot)
}

// reconnect reconnects an already connected session to a WebSocket Connection.
func (s *Session) reconnect(bot *Client) error {
	// close the active connection with a non-1000 and non-1001 close code.
	if err := s.disconnect(FlagClientCloseEventCodeReconnect); err != nil {
		return ErrorDisconnect{
			SessionID: s.ID,
			Err:       err,
			Action:    errOpcodeReconnect,
		}
	}

	// allow Discord to close the session.
	<-time.After(time.Second)

	// connect to the Discord Gateway again.
	if err := s.connect(bot); err != nil {
		return fmt.Errorf("an error occurred while reconnecting Session %q: %w", s.ID, err)
	}

	return nil
}

// listen listens to the connection for payloads from the Discord Gateway.
func (s *Session) listen(bot *Client) {
	for {
		payload := getPayload()
		if err := socket.Read(s.Context, s.Conn, payload); err != nil {
			s.Lock()
			defer s.Unlock()

			select {
			case <-s.Context.Done():
			default:
				closeErr := new(websocket.CloseError)
				if errors.As(err, closeErr) {
					if gcErr := s.handleGatewayCloseError(bot, closeErr); gcErr == nil {
						return
					}
				}

				s.disconnectFromRoutine("listen: Closing the connection due to a read error...",
					ErrorEvent{
						Event:  "Payload",
						Err:    err,
						Action: ErrorEventActionRead,
					})
			}

			return
		}

		log.Println("PAYLOAD", payload.Op, string(payload.Data))
		if err := s.onPayload(bot, *payload); err != nil {
			s.Lock()
			defer s.Unlock()

			select {
			case <-s.Context.Done():
			default:
				s.disconnectFromRoutine("onPayload error", err)
			}

			return
		}
	}
}

// onPayload handles an Discord Gateway Payload.
func (s *Session) onPayload(bot *Client, payload GatewayPayload) error {
	defer putPayload(&payload)

	// https://discord.com/developers/docs/topics/opcodes-and-status-codes#gateway-gateway-opcodes
	switch payload.Op {
	// run the bot's event handlers.
	case FlagGatewayOpcodeDispatch:
		atomic.StoreInt64(&s.Seq, *payload.SequenceNumber)
		go bot.handle(*payload.EventName, payload.Data)

		// Sending a valid Identify Payload triggers the initial handshake with the Discord Gateway.
		// This will result in the Gateway responding with a Ready event.
		// The handler for this Ready event is located in the onPayload function.
		// This allows developers to manipulate the Handlers.Ready slice without issue.
		if *payload.EventName == FlagGatewayEventNameReady {
			ready := new(Ready)
			if err := json.Unmarshal(payload.Data, ready); err != nil {
				return ErrorEvent{
					Event:  FlagGatewayEventNameReady,
					Err:    err,
					Action: ErrorEventActionUnmarshal,
				}
			}

			s.Lock()
			s.ID = ready.SessionID
			log.Printf("received Ready event for Session %q", s.ID)
			s.Unlock()
			// SHARD: set shard information using r.Shard
			bot.ApplicationID = ready.Application.ID
		}

	// send an Opcode 1 Heartbeat to the Discord Gateway.
	case FlagGatewayOpcodeHeartbeat:
		go s.respond(payload.Data)

	// handle the successful acknowledgement of the client's last heartbeat.
	case FlagGatewayOpcodeHeartbeatACK:
		s.Lock()
		log.Println("ACK")
		atomic.AddUint32(&s.heartbeat.acks, 1)
		s.Unlock()

	// occurs when the maximum concurrency limit has been reached while connecting,
	// or when the session does NOT reconnect in time.
	case FlagGatewayOpcodeInvalidSession:
		// wait for Discord to close the session, then complete a fresh connect.
		<-time.NewTimer(invalidSessionWaitTime).C

		s.Lock()
		defer s.Unlock()

		s.ID = ""
		s.Seq = 0
		if err := s.initial(bot); err != nil {
			return err
		}

		if !s.isConnected() {
			return fmt.Errorf("Session %q couldn't connect to the Discord Gateway or has invalidated an active session", s.ID)
		}

	// occurs when the Discord Gateway is shutting down the connection,
	// while signalling the client to reconnect.
	case FlagGatewayOpcodeReconnect:
		s.Lock()
		defer s.Unlock()

		log.Printf("reconnecting Session %q due to Opcode 7 Reconnect", s.ID)

		return s.reconnect(bot)
	}

	return nil
}

// handleGatewayCloseError handles a Discord Gateway WebSocket CloseError.
//
// returns the given closeErr if a disconnect is warranted.
func (s *Session) handleGatewayCloseError(bot *Client, closeErr *websocket.CloseError) error {
	code, ok := GatewayCloseEventCodes[int(closeErr.Code)]
	switch ok {
	// Gateway Close Event Code is known.
	case true:
		log.Printf(
			"Session %q received Gateway Close Event Code %d %s: %s",
			s.ID, code.Code, code.Description, code.Explanation,
		)

		if code.Reconnect {
			if reconnectErr := s.reconnect(bot); reconnectErr != nil {
				log.Println(reconnectErr)

				return closeErr
			}

			return nil
		}

		return closeErr

	// Gateway Close Event Code is unknown.
	default:

		// when another goroutine calls disconnect(),
		// s.Conn.Close is called before s.cancel which will result in
		// a CloseError with the close code that Disgo uses to reconnect.
		if closeErr.Code == websocket.StatusCode(FlagClientCloseEventCodeReconnect) {
			return nil
		}

		log.Printf(
			"Session %q received unknown Gateway Close Event Code %d with reason %q",
			s.ID, closeErr.Code, closeErr.Reason,
		)

		return closeErr
	}
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

// readEvent is a helper function for reading events from the WebSocket Session.
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

// writeEvent is a helper function for writing events to the WebSocket Session.
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
