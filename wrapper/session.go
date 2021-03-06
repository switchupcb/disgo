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
	"github.com/switchupcb/websocket"
	"golang.org/x/sync/errgroup"
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

	// heartbeat contains the fields required to implement the heartbeat mechanism.
	heartbeat *heartbeat

	// manager represents a manager of a Session's goroutines.
	manager *manager

	// RWMutex is used to protect the Session's variables from data races
	// by providing transactional functionality.
	sync.RWMutex
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

		// set the maximum allowed (Identify) concurrency rate limit.
		//
		// https://discord.com/developers/docs/topics/gateway#rate-limiting
		bot.Config.Gateway.RateLimiter.StartTx()

		identifyBucket := bot.Config.Gateway.RateLimiter.GetBucketFromHash(FlagGatewayCommandNameIdentify)
		if identifyBucket == nil {
			identifyBucket = getBucket()
			bot.Config.Gateway.RateLimiter.SetBucketFromHash(FlagGatewayCommandNameIdentify, identifyBucket)
		}

		identifyBucket.Limit = int16(response.SessionStartLimit.MaxConcurrency) + 1
		identifyBucket.Remaining = identifyBucket.Limit
		identifyBucket.Expiry = time.Now().Add(FlagGlobalRateLimitIdentifyInterval)

		bot.Config.Gateway.RateLimiter.EndTx()
	}

	var err error

	// connect to the Discord Gateway Websocket.
	s.manager = new(manager)
	s.Context, s.manager.cancel = context.WithCancel(context.Background())
	if s.Conn, _, err = websocket.Dial(s.Context, s.Endpoint, nil); err != nil {
		return fmt.Errorf("an error occurred while connecting to the Discord Gateway\n%w", err)
	}

	// handle the incoming Hello event upon connecting to the Gateway.
	hello := new(Hello)
	if err := readEvent(s, FlagGatewayEventNameHello, hello); err != nil {
		if disconnectErr := s.disconnect(FlagClientCloseEventCodeNormal); disconnectErr != nil {
			return DisconnectError{
				SessionID: s.ID,
				Err:       disconnectErr,
				Action:    err,
			}
		}

		return err
	}

	for _, handler := range bot.Handlers.Hello {
		go handler(hello)
	}

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

	// create a goroutine group for the Session.
	s.manager.Group, s.manager.signal = errgroup.WithContext(s.Context)
	s.manager.err = make(chan error, 1)

	// spawn the heartbeat pulse goroutine.
	s.manager.routines.Add(1)
	atomic.AddInt32(&s.manager.pulses, 1)
	s.manager.Go(func() error {
		s.pulse()
		return nil
	})

	// spawn the heartbeat beat goroutine.
	s.manager.routines.Add(1)
	s.manager.Go(func() error {
		if err := s.beat(bot); err != nil {
			return fmt.Errorf("heartbeat: %w", err)
		}

		return nil
	})

	// send the initial Identify or Resumed packet.
	if err := s.initial(bot, 0); err != nil {
		if disconnectErr := s.disconnect(FlagClientCloseEventCodeNormal); disconnectErr != nil {
			return DisconnectError{
				SessionID: s.ID,
				Err:       disconnectErr,
				Action:    err,
			}
		}

		return err
	}

	// spawn the event listener listen goroutine.
	s.manager.routines.Add(1)
	s.manager.Go(func() error {
		if err := s.listen(bot); err != nil {
			return fmt.Errorf("listen: %w", err)
		}

		return nil
	})

	// spawn the manager goroutine.
	s.manager.routines.Add(1)
	go s.manage(bot)

	// ensure that the Session's goroutines are spawned.
	s.manager.routines.Wait()

	return nil
}

// initial sends the initial Identify or Resume packet required to connect to the Gateway,
// then handles the incoming Ready or Resumed packet that indicates a successful connection.
func (s *Session) initial(bot *Client, attempt int) error {
	if !s.canReconnect() {
		// send an Opcode 2 Identify to the Discord Gateway.
		identify := Identify{
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
		}

		if err := identify.Command(bot, s); err != nil {
			return err
		}
	} else {
		// send an Opcode 6 Resume to the Discord Gateway to reconnect the session.
		resume := Resume{
			Token:     bot.Authentication.Token,
			SessionID: s.ID,
			Seq:       atomic.LoadInt64(&s.Seq),
		}

		if err := resume.Command(bot, s); err != nil {
			return err
		}
	}

	// handle the incoming Ready, Resumed or Replayed event (or Opcode 9 Invalid Session).
	payload := new(GatewayPayload)
	if err := socket.Read(s.Context, s.Conn, payload); err != nil {
		return fmt.Errorf("error occurred while reading initial payload: %w", err)
	}

	log.Println("INITIAL PAYLOAD", payload.Op, string(payload.Data))

	switch payload.Op {
	case FlagGatewayOpcodeDispatch:
		switch {
		// When a connection is successful, the Discord Gateway will respond with a Ready event.
		case *payload.EventName == FlagGatewayEventNameReady:
			ready := new(Ready)
			if err := json.Unmarshal(payload.Data, ready); err != nil {
				return fmt.Errorf("%w", err)
			}

			s.ID = ready.SessionID
			s.Seq = 0
			// SHARD: set shard information using r.Shard
			bot.ApplicationID = ready.Application.ID

			log.Printf("received Ready event for Session %q", s.ID)

			for _, handler := range bot.Handlers.Ready {
				go handler(ready)
			}

		// When a reconnection is successful, the Discord Gateway will respond
		// by replaying all missed events in order, finalized by a Resumed event.
		case *payload.EventName == FlagGatewayEventNameResumed:
			log.Printf("received Resumed event for Session %q", s.ID)

			for _, handler := range bot.Handlers.Resumed {
				go handler(&Resumed{})
			}

		// When a reconnection is successful, the Discord Gateway will respond
		// by replaying all missed events in order, finalized by a Resumed event.
		default:
			// handle the initial payload(s) until a Resumed event is encountered.
			go bot.handle(*payload.EventName, payload.Data)

			for {
				replayed := new(GatewayPayload)
				if err := socket.Read(s.Context, s.Conn, replayed); err != nil {
					return fmt.Errorf("error occurred while replaying events: %w", err)
				}

				if replayed.Op == FlagGatewayOpcodeDispatch && *replayed.EventName == FlagGatewayEventNameResumed {
					log.Printf("received Resumed event for Session %q", s.ID)

					for _, handler := range bot.Handlers.Resumed {
						go handler(&Resumed{})
					}

					return nil
				}

				go bot.handle(*payload.EventName, payload.Data)
			}
		}

	// When the maximum concurrency limit has been reached while connecting, or when
	// the session does NOT reconnect in time, the Discord Gateway send an Opcode 9 Invalid Session.
	case FlagGatewayOpcodeInvalidSession:
		if attempt < 1 {
			// wait for Discord to close the session, then complete a fresh connect.
			<-time.NewTimer(invalidSessionWaitTime).C

			s.ID = ""
			s.Seq = 0
			if err := s.initial(bot, attempt+1); err != nil {
				return err
			}

			return nil
		}

		return fmt.Errorf("Session %q couldn't connect to the Discord Gateway or has invalidated an active session", s.ID)
	default:
		return fmt.Errorf("Session %q received payload %d during connection which is unexpected", s.ID, payload.Op)
	}

	return nil
}

// Disconnect disconnects a session from the Discord Gateway using the given status code.
func (s *Session) Disconnect() error {
	s.Lock()

	if !s.isConnected() {
		s.Unlock()

		return fmt.Errorf("Session %q is already disconnected", s.ID)
	}

	id := s.ID
	log.Printf("disconnecting Session %q with code %d", id, FlagClientCloseEventCodeNormal)

	s.manager.signal = context.WithValue(s.manager.signal, keySignal, signalDisconnect)

	if err := s.disconnect(FlagClientCloseEventCodeNormal); err != nil {
		s.Unlock()

		return err
	}

	s.Unlock()

	if err := <-s.manager.err; err != nil {
		return err
	}

	putSession(s)

	log.Printf("disconnected Session %q with code %d", id, FlagClientCloseEventCodeNormal)

	return nil
}

// disconnect disconnects a session from a WebSocket Connection using the given status code.
func (s *Session) disconnect(code int) error {
	// cancel the context to kill the goroutines of the Session.
	defer s.manager.cancel()

	if err := s.Conn.Close(websocket.StatusCode(code), ""); err != nil {
		return DisconnectError{
			SessionID: s.ID,
			Err:       err,
			Action:    nil,
		}
	}

	return nil
}

// Reconnect reconnects an already connected session to the Discord Gateway
// by disconnecting the session, then connecting again.
func (s *Session) Reconnect(bot *Client) error {
	s.reconnect(fmt.Sprintf("reconnecting Session %q", s.ID))

	if err := <-s.manager.err; err != nil {
		return err
	}

	// connect to the Discord Gateway again.
	if err := s.Connect(bot); err != nil {
		return fmt.Errorf("an error occurred while reconnecting Session %q: %w", s.ID, err)
	}

	return nil
}

// readEvent is a helper function for reading events from the WebSocket Session.
func readEvent(s *Session, name string, dst any) error {
	payload := new(GatewayPayload)
	if err := socket.Read(s.Context, s.Conn, payload); err != nil {
		return fmt.Errorf("readEvent: %w", err)
	}

	if err := json.Unmarshal(payload.Data, dst); err != nil {
		return fmt.Errorf("readEvent: %w", err)
	}

	return nil
}

// writeEvent is a helper function for writing events to the WebSocket Session.
func writeEvent(bot *Client, s *Session, op int, name string, dst any) error {
RATELIMIT:
	// a single command is PROCESSED at any point in time.
	bot.Config.Gateway.RateLimiter.Lock()

	for {
		bot.Config.Gateway.RateLimiter.StartTx()

		globalBucket := bot.Config.Gateway.RateLimiter.GetBucket(0)

		// stop waiting when the Global Rate Limit Bucket is NOT empty.
		if isNotEmpty(globalBucket) { //nolint:nestif
			switch op {
			case FlagGatewayOpcodeIdentify:
				identifyBucket := bot.Config.Gateway.RateLimiter.GetBucketFromHash(FlagGatewayCommandNameIdentify)

				if isNotEmpty(identifyBucket) {
					if globalBucket != nil {
						globalBucket.Remaining--
					}

					if identifyBucket != nil {
						identifyBucket.Remaining--
					}

					bot.Config.Gateway.RateLimiter.EndTx()

					goto SEND
				}

				if isExpired(identifyBucket) {
					if globalBucket != nil {
						globalBucket.Remaining--
					}

					if identifyBucket != nil {
						identifyBucket.Reset(time.Now().Add(FlagGlobalRateLimitIdentifyInterval))
						identifyBucket.Remaining--
					}

					bot.Config.Gateway.RateLimiter.EndTx()

					goto SEND
				}

				var wait time.Time
				if identifyBucket != nil {
					wait = identifyBucket.Expiry
				}

				// do NOT block other requests due to a Command Rate Limit.
				bot.Config.Gateway.RateLimiter.EndTx()
				bot.Config.Gateway.RateLimiter.Unlock()

				// reduce CPU usage by blocking the current goroutine
				// until it's eligible for action.
				if identifyBucket != nil {
					<-time.After(time.Until(wait))
				}

				goto RATELIMIT

			default:
				if globalBucket != nil {
					globalBucket.Remaining--
				}

				bot.Config.Gateway.RateLimiter.EndTx()

				goto SEND
			}
		}

		// reset the Global Rate Limit Bucket when the current Bucket has passed its expiry.
		if isExpired(globalBucket) {
			globalBucket.Reset(time.Now().Add(time.Minute))
		}

		bot.Config.Gateway.RateLimiter.EndTx()
	}

SEND:
	bot.Config.Gateway.RateLimiter.Unlock()

	// write the event to the WebSocket Connection.
	event, err := json.Marshal(dst)
	if err != nil {
		return fmt.Errorf("writeEvent: %w", err)
	}

	if err = socket.Write(s.Context, s.Conn, websocket.MessageBinary,
		GatewayPayload{ //nolint:exhaustruct
			Op:   op,
			Data: event,
		}); err != nil {
		return fmt.Errorf("writeEvent: %w", err)
	}

	return nil
}
