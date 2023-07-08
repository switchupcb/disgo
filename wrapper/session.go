package wrapper

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	json "github.com/goccy/go-json"
	"github.com/switchupcb/disgo/wrapper/socket"
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

	// Endpoint represents the endpoint that is used to reconnect to the Gateway.
	Endpoint string

	// Shard represents the [shard_id, num_shards] for the Session.
	//
	// https://discord.com/developers/docs/topics/gateway#sharding
	Shard *[2]int

	// Context carries request-scoped data for the Discord Gateway Connection.
	//
	// Context is also used as a signal for the Session's goroutines.
	Context context.Context

	// Conn represents a WebSocket Connection to the Discord Gateway.
	Conn *websocket.Conn

	// heartbeat contains the fields required to implement the heartbeat mechanism.
	heartbeat *heartbeat

	// manager represents a manager of a Session's goroutines.
	manager *manager

	// client_manager represents the *Client Session Manager of the Session.
	client_manager *SessionManager

	// RateLimiter represents an object that provides rate limit functionality.
	RateLimiter RateLimiter

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
	return s.ID != "" && s.Endpoint != "" && atomic.LoadInt64(&s.Seq) != 0
}

// Connect connects a session to the Discord Gateway (WebSocket Connection).
func (s *Session) Connect(bot *Client) error {
	s.Lock()
	defer s.Unlock()

	LogSession(Logger.Info(), s.ID).Str(LogCtxClient, bot.ApplicationID).Msg("connecting session")

	return s.connect(bot)
}

// connect connects a session to a WebSocket Connection.
func (s *Session) connect(bot *Client) error {
	if bot.Sessions == nil {
		return fmt.Errorf(errNoSessionManager) //lint:ignore ST1005 format help message.
	}

	s.client_manager = bot.Sessions

	if s.isConnected() {
		return fmt.Errorf("session %q is already connected", s.ID)
	}

	// request a valid Gateway URL endpoint from the Discord API.
	gatewayEndpoint := s.Endpoint
	if gatewayEndpoint == "" || !s.canReconnect() {
		gateway := GetGatewayBot{}
		response, err := gateway.Send(bot)
		if err != nil {
			return fmt.Errorf("error getting the Gateway API Endpoint: %w", err)
		}

		gatewayEndpoint = response.URL + gatewayEndpointParams

		// set the maximum allowed (Identify) concurrency rate limit.
		//
		// https://discord.com/developers/docs/topics/gateway#rate-limiting
		bot.Config.Gateway.RateLimiter.StartTx()

		identifyBucket := bot.Config.Gateway.RateLimiter.GetBucketFromID(FlagGatewaySendEventNameIdentify)
		if identifyBucket == nil {
			identifyBucket = getBucket()
			bot.Config.Gateway.RateLimiter.SetBucketFromID(FlagGatewaySendEventNameIdentify, identifyBucket)
		}

		if bot.Config.Gateway.ShardManager != nil {
			bot.Config.Gateway.ShardManager.SetLimit(
				ShardLimit{
					Reset:             time.Now().Add(time.Millisecond*time.Duration(response.SessionStartLimit.ResetAfter) + 1),
					MaxStarts:         response.SessionStartLimit.Total,
					RemainingStarts:   response.SessionStartLimit.Remaining,
					MaxConcurrency:    response.SessionStartLimit.MaxConcurrency,
					RecommendedShards: response.Shards,
				},
			)
		}

		identifyBucket.Limit = int16(response.SessionStartLimit.MaxConcurrency)

		if identifyBucket.Expiry.IsZero() {
			identifyBucket.Remaining = identifyBucket.Limit
			identifyBucket.Expiry = time.Now().Add(FlagGlobalRateLimitIdentifyInterval)
		}

		bot.Config.Gateway.RateLimiter.EndTx()
	}

	var err error

	// connect to the Discord Gateway Websocket.
	s.manager = new(manager)
	s.Context, s.manager.cancel = context.WithCancel(context.Background())
	if s.Conn, _, err = websocket.Dial(s.Context, gatewayEndpoint, nil); err != nil {
		return fmt.Errorf("error connecting to the Discord Gateway: %w", err)
	}

	// set up the Session's Rate Limiter (applied per WebSocket Connection).
	// https://discord.com/developers/docs/topics/gateway#rate-limiting
	s.RateLimiter = &RateLimit{ //nolint:exhaustruct
		ids:     make(map[string]string, totalGatewayBucketsPerConnection),
		buckets: make(map[string]*Bucket, totalGatewayBucketsPerConnection),
	}

	s.RateLimiter.SetBucket(
		GlobalRateLimitRouteID, &Bucket{ //nolint:exhaustruct
			Limit:     FlagGlobalRateLimitGateway,
			Remaining: FlagGlobalRateLimitGateway,
			Expiry:    time.Now().Add(FlagGlobalRateLimitGatewayInterval),
		},
	)

	// handle the incoming Hello event upon connecting to the Gateway.
	hello := new(Hello)
	if err := readEvent(s, hello); err != nil {
		err = fmt.Errorf("error reading initial Hello event: %w", err)
		sessionErr := ErrorSession{SessionID: s.ID, Err: err}
		if disconnectErr := s.disconnect(FlagClientCloseEventCodeNormal); disconnectErr != nil {
			sessionErr.Err = ErrorDisconnect{
				Action:     err,
				Err:        disconnectErr,
				Connection: ErrConnectionSession,
			}
		}

		return sessionErr
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
			return ErrorSession{
				SessionID: s.ID,
				Err:       fmt.Errorf("heartbeat: %w", err),
			}
		}

		return nil
	})

	// send the initial Identify or Resumed packet.
	if err := s.initial(bot, 0); err != nil {
		sessionErr := ErrorSession{SessionID: s.ID, Err: err}
		if disconnectErr := s.disconnect(FlagClientCloseEventCodeNormal); disconnectErr != nil {
			sessionErr.Err = ErrorDisconnect{
				Action:     err,
				Err:        disconnectErr,
				Connection: ErrConnectionSession,
			}
		}

		return sessionErr
	}

	// spawn the event listener listen goroutine.
	s.manager.routines.Add(1)
	s.manager.Go(func() error {
		if err := s.listen(bot); err != nil {
			return ErrorSession{
				SessionID: s.ID,
				Err:       fmt.Errorf("listen: %w", err),
			}
		}

		return nil
	})

	// spawn the manager goroutine.
	s.manager.routines.Add(1)
	go s.manage()

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
			Compress:       Pointer(true),
			LargeThreshold: Pointer(maxIdentifyLargeThreshold),
			Shard:          s.Shard,
			Presence:       bot.Config.Gateway.GatewayPresenceUpdate,
			Intents:        bot.Config.Gateway.Intents,
		}

		if err := identify.SendEvent(bot, s); err != nil {
			return err
		}
	} else {
		// send an Opcode 6 Resume to the Discord Gateway to reconnect the session.
		resume := Resume{
			Token:     bot.Authentication.Token,
			SessionID: s.ID,
			Seq:       atomic.LoadInt64(&s.Seq),
		}

		if err := resume.SendEvent(bot, s); err != nil {
			return err
		}
	}

	// handle the incoming Ready, Resumed or Replayed event (or Opcode 9 Invalid Session).
	payload := new(GatewayPayload)
	if err := socket.Read(s.Context, s.Conn, payload); err != nil {
		return fmt.Errorf("error reading initial payload: %w", err)
	}

	LogPayload(LogSession(Logger.Info(), s.ID), payload.Op, payload.Data).Msg("received initial payload")

	switch payload.Op {
	case FlagGatewayOpcodeDispatch:
		switch {
		// When a connection is successful, the Discord Gateway will respond with a Ready event.
		case *payload.EventName == FlagGatewayEventNameReady:
			ready := new(Ready)
			if err := json.Unmarshal(payload.Data, ready); err != nil {
				return fmt.Errorf("error reading ready event: %w", err)
			}

			LogSession(Logger.Info(), ready.SessionID).Msg("received Ready event")

			// Configure the session.
			s.ID = ready.SessionID
			atomic.StoreInt64(&s.Seq, 0)
			s.Endpoint = ready.ResumeGatewayURL
			bot.ApplicationID = ready.Application.ID

			// Store the session in the session manager.
			s.client_manager.Gateway.Store(s.ID, s)

			if bot.Config.Gateway.ShardManager != nil {
				bot.Config.Gateway.ShardManager.Ready(bot, s, ready)
			}

			for _, handler := range bot.Handlers.Ready {
				go handler(ready)
			}

		// When a reconnection is successful, the Discord Gateway will respond
		// by replaying all missed events in order, finalized by a Resumed event.
		case *payload.EventName == FlagGatewayEventNameResumed:
			LogSession(Logger.Info(), s.ID).Msg("received Resumed event")

			// Store the session in the session manager.
			s.client_manager.Gateway.Store(s.ID, s)

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
					return fmt.Errorf("error replaying events: %w", err)
				}

				if replayed.Op == FlagGatewayOpcodeDispatch && *replayed.EventName == FlagGatewayEventNameResumed {
					LogSession(Logger.Info(), s.ID).Msg("received Resumed event")

					// Store the session in the session manager.
					s.client_manager.Gateway.Store(s.ID, s)

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
		// Remove the session from the session manager.
		s.client_manager.Gateway.Store(s.ID, nil)

		if attempt < 1 {
			// wait for Discord to close the session, then complete a fresh connect.
			<-time.NewTimer(invalidSessionWaitTime).C

			s.ID = ""
			atomic.StoreInt64(&s.Seq, 0)
			if err := s.initial(bot, attempt+1); err != nil {
				return err
			}

			return nil
		}

		return fmt.Errorf("session %q couldn't connect to the Discord Gateway or has invalidated an active session", s.ID)
	default:
		return fmt.Errorf("session %q received payload %d during connection which is unexpected", s.ID, payload.Op)
	}

	return nil
}

// Disconnect disconnects a session from the Discord Gateway using the given status code.
func (s *Session) Disconnect() error {
	s.Lock()

	if !s.isConnected() {
		s.Unlock()

		return fmt.Errorf("session %q is already disconnected", s.ID)
	}

	id := s.ID
	LogSession(Logger.Info(), id).Msgf("disconnecting session with code %d", FlagClientCloseEventCodeNormal)

	s.manager.signal = context.WithValue(s.manager.signal, keySignal, signalDisconnect)

	if err := s.disconnect(FlagClientCloseEventCodeNormal); err != nil {
		s.Unlock()

		return ErrorDisconnect{
			Connection: ErrConnectionSession,
			Action:     nil,
			Err:        err,
		}
	}

	s.Unlock()

	if err := <-s.manager.err; err != nil {
		return err
	}

	putSession(s)

	LogSession(Logger.Info(), id).Msgf("disconnected session with code %d", FlagClientCloseEventCodeNormal)

	return nil
}

// disconnect disconnects a session from a WebSocket Connection using the given status code.
func (s *Session) disconnect(code int) error {
	// cancel the context to kill the goroutines of the Session.
	defer s.manager.cancel()

	// Remove the session from the session manager.
	s.client_manager.Gateway.Store(s.ID, nil)

	if err := s.Conn.Close(websocket.StatusCode(code), ""); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

// Reconnect reconnects an already connected session to the Discord Gateway
// by disconnecting the session, then connecting again.
func (s *Session) Reconnect(bot *Client) error {
	s.reconnect("reconnecting")

	if err := <-s.manager.err; err != nil {
		return err
	}

	// connect to the Discord Gateway again.
	if err := s.Connect(bot); err != nil {
		return fmt.Errorf("error reconnecting session %q: %w", s.ID, err)
	}

	return nil
}

// readEvent is a helper function for reading events from the WebSocket Session.
func readEvent(s *Session, dst any) error {
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
	s.RLock()

	// make sure the Session isn't disconnected while sending an event.
	if !s.isConnected() {
		s.RUnlock()

		return fmt.Errorf("writeEvent: session is disconnected")
	}

	// a single send event is PROCESSED at any point in time.
	s.RateLimiter.Lock()

	LogCommand(LogSession(Logger.Trace(), s.ID), bot.ApplicationID, op, name).Msg("processing gateway command")

	for {
		s.RateLimiter.StartTx()

		globalBucket := s.RateLimiter.GetBucket(GlobalRateLimitRouteID, "")

		// reset the Global Rate Limit Bucket when the current Bucket has passed its expiry.
		if isExpired(globalBucket) {
			globalBucket.Reset(time.Now().Add(time.Minute))
		}

		// stop waiting when the Global Rate Limit Bucket is NOT empty.
		if isNotEmpty(globalBucket) {
			switch op {
			// Identify is also bound by the max_concurrency rate limit.
			case FlagGatewayOpcodeIdentify:
				bot.Config.Gateway.RateLimiter.StartTx()

				identifyBucket := bot.Config.Gateway.RateLimiter.GetBucketFromID(FlagGatewaySendEventNameIdentify)

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

				// do NOT block other send events due to a Send Event Rate Limit.
				bot.Config.Gateway.RateLimiter.EndTx()
				s.RateLimiter.EndTx()
				s.RateLimiter.Unlock()
				s.RUnlock()

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

				s.RateLimiter.EndTx()

				goto SEND
			}
		}

		s.RateLimiter.EndTx()
	}

SEND:
	s.RateLimiter.Unlock()
	defer s.RUnlock()

	LogCommand(LogSession(Logger.Trace(), s.ID), bot.ApplicationID, op, name).Msg("sending gateway command")

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

	LogCommand(LogSession(Logger.Trace(), s.ID), bot.ApplicationID, op, name).Msg("sent gateway command")

	return nil
}
