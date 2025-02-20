package wrapper

import (
	"context"
	"fmt"
	"runtime"
	"sync/atomic"
	"time"

	json "github.com/goccy/go-json"
	"github.com/switchupcb/disgo/wrapper/socket"
	"github.com/switchupcb/websocket"
)

const (
	gatewayEndpointParams     = "?v=" + VersionDiscordAPI + "&encoding=json"
	invalidSessionWaitTime    = 1 * time.Second
	maxIdentifyLargeThreshold = 250
)

// connect connects a session to a WebSocket Connection.
func (s *Session) connect(bot *Client) error {
	var err error

	// request a valid Gateway URL endpoint and response from the Discord API.
	gatewayEndpoint := s.Endpoint
	var response *GetGatewayBotResponse

	if bot.Config.Gateway.ShardManager != nil {
		if response, err = bot.Config.Gateway.ShardManager.SetLimit(bot); err != nil {
			return fmt.Errorf("shardmanager: %w", err)
		}
	} else {
		if gatewayEndpoint == "" || !s.canReconnect() {
			gateway := GetGatewayBot{}
			response, err = gateway.Send(bot)
			if err != nil {
				return fmt.Errorf("error getting the Gateway API Endpoint: %w", err)
			}

			gatewayEndpoint = response.URL
		}
	}

	// set the maximum allowed (Identify) concurrency rate limit for the bot.
	//
	// https://discord.com/developers/docs/topics/gateway#rate-limiting
	if response != nil {
		bot.Config.Gateway.RateLimiter.StartTx()

		identifyBucket := bot.Config.Gateway.RateLimiter.GetBucketFromID(FlagGatewaySendEventNameIdentify)
		if identifyBucket == nil {
			identifyBucket = getBucket()
			bot.Config.Gateway.RateLimiter.SetBucketFromID(FlagGatewaySendEventNameIdentify, identifyBucket)
		}

		identifyBucket.Limit = int16(response.SessionStartLimit.MaxConcurrency) //nolint:gosec // disable G115

		if identifyBucket.Expiry.IsZero() {
			identifyBucket.Remaining = identifyBucket.Limit
			identifyBucket.Expiry = time.Now().Add(FlagGlobalRateLimitIdentifyInterval)
		}

		bot.Config.Gateway.RateLimiter.EndTx()
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

	// connect to the Discord Gateway Websocket.
	s.Context, s.manager.cancel = context.WithCancel(context.Background())
	if s.Conn, _, err = websocket.Dial(s.Context, gatewayEndpoint+gatewayEndpointParams, nil); err != nil {
		return fmt.Errorf("error connecting to the Discord Gateway: %w", err)
	}

	s.setState(SessionStateConnectingWebsocket)

	// handle the incoming Hello event upon connecting to the Gateway.
	hello := new(Hello)
	if err := readEvent(s, hello); err != nil {
		return fmt.Errorf("error reading initial Hello event: %w", err)
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

	// spawn the heartbeat pulse goroutine.
	s.manager.routines.Add(1)
	s.manager.Go(func() error {
		atomic.AddInt32(&s.manager.pulses, 1)
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
		return fmt.Errorf("initial: %w", err)
	}

	// spawn the event listener listen goroutine.
	s.manager.routines.Add(1)
	s.manager.Go(func() error {
		if err := s.listen(bot); err != nil {
			return fmt.Errorf("listen: %w", err)
		}

		return nil
	})

	// confirm the Session's goroutines are spawned.
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
	payload := getPayload()
	defer putPayload(payload)
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
		s.client_manager.RemoveGatewaySession(s.ID)

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
		return fmt.Errorf("session %q received unexpected payload %d during connection", s.ID, payload.Op)
	}

	return nil
}

// disconnect disconnects a session from a WebSocket Connection using the given status code.
func (s *Session) disconnect(code int) error {
	// cancel the context to kill the goroutines of the Session.
	defer s.manager.cancel()

	if err := s.Conn.Close(websocket.StatusCode(code), ""); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

// reconnect spawns a goroutine for reconnection which prompts the manager
// to reconnect upon a disconnection.
func (s *Session) reconnect(reason string) {
	LogSession(Logger.Info(), s.ID).Msg(reason)

	s.manager.signals <- sessionSignalReconnect
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
