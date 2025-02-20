package wrapper

import (
	"fmt"
	"time"

	json "github.com/goccy/go-json"
	"github.com/switchupcb/disgo/wrapper/socket"
	"github.com/switchupcb/websocket"
)

// writeEvent is a helper function for writing events to the WebSocket Session.
func writeEvent(bot *Client, s *Session, op int, name string, dst any) error {
RATELIMIT:
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
						if globalBucket.Remaining == FlagGlobalRateLimitGateway {
							globalBucket.Reset(time.Now().Add(time.Minute))
						}

						globalBucket.Remaining--
					}

					if identifyBucket != nil {
						identifyBucket.Remaining--
					}

					bot.Config.Gateway.RateLimiter.EndTx()
					s.RateLimiter.EndTx()

					goto SEND
				}

				if isExpired(identifyBucket) {
					if globalBucket != nil {
						if globalBucket.Remaining == FlagGlobalRateLimitGateway {
							globalBucket.Reset(time.Now().Add(time.Minute))
						}

						globalBucket.Remaining--
					}

					if identifyBucket != nil {
						identifyBucket.Reset(time.Now().Add(FlagGlobalRateLimitIdentifyInterval))
						identifyBucket.Remaining--
					}

					bot.Config.Gateway.RateLimiter.EndTx()
					s.RateLimiter.EndTx()

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

				// reduce CPU usage by blocking the current goroutine
				// until it's eligible for action.
				if identifyBucket != nil {
					<-time.After(time.Until(wait))
				}

				goto RATELIMIT

			default:
				if globalBucket != nil {
					if globalBucket.Remaining == FlagGlobalRateLimitGateway {
						globalBucket.Reset(time.Now().Add(time.Minute))
					}

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
