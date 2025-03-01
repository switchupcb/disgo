package tools

import (
	"errors"
	"strings"
	"time"

	"github.com/switchupcb/disgo"
	"github.com/switchupcb/websocket"
)

var (
	// ReconnectOnUnexpectedDisconnectionHandlers represents default handlers the user can pass to
	// ReconnectOnUnexpectedDisconnection() to reconnect the bot.
	//
	// Read https://github.com/switchupcb/disgo/discussions/82 for more information about each handler.
	ReconnectOnUnexpectedDisconnectionHandlers = []func(error) bool{
		// failed to get reader: failed to read frame header: EOF
		func(err error) bool {
			return strings.Contains(err.Error(), "failed to get reader: failed to read frame header: EOF")
		},
		// received unknown Gateway Close Event Code 1001 with reason CloudFlare WebSocket proxy restarting
		func(err error) bool {
			closeErr := new(websocket.CloseError)

			return errors.As(err, closeErr) && closeErr.Code == 1001
		},
	}

	// ReconnectOnUnexpectedDisconnectionWait represents the amount of time to wait before reconnecting
	// to the Discord Gateway in ReconnectOnUnexpectedDisconnection() [default: time.Second * 5].
	ReconnectOnUnexpectedDisconnectionWait = time.Second * 5
)

// ReconnectOnUnexpectedDisconnection reconnects a session when it disconnects unexpectedly.
//
// Use this function with `tools.ReconnectOnUnexpectedDisconnectionHandlers` or custom handlers
// (`func(error) bool`): The bot is reconnected when a single handler returns true.
//
// WARNING: You are recommended to run this function on a goroutine.
//
// Calling this function blocks the thread until the Session exits the connection loop.
func ReconnectOnUnexpectedDisconnection(bot *disgo.Client, s *disgo.Session, reconnectOnErr ...func(error) bool) {
	s.Lock()
	sessionid := s.ID
	s.Unlock()

	disgo.LogSession(disgo.Logger.Info(), sessionid).Msg("ReconnectOnUnexpectedDisconnection: reconnection loop started")

	defer func() {
		disgo.LogSession(disgo.Logger.Info(), sessionid).Msg("ReconnectOnUnexpectedDisconnection: reconnection loop closed")
	}()

	for {
		reconnect := false
		if _, err := s.Wait(); err != nil {
			for _, handler := range reconnectOnErr {
				if handler(err) {
					reconnect = true

					break
				}
			} // handlers

			if reconnect {
				disgo.LogSession(disgo.Logger.Info(), sessionid).Msgf("ReconnectOnUnexpectedDisconnection: Reconnecting in 5 seconds due to handled error: %v", err) //nolint:lll

				<-time.After(ReconnectOnUnexpectedDisconnectionWait)

				if sErr := s.Connect(bot); sErr != nil {
					disgo.LogSession(disgo.Logger.Info(), sessionid).Err(sErr).Msgf("ReconnectOnUnexpectedDisconnection: could not connect to Discord Gateway") //nolint:lll

					return
				}

				s.Lock()
				sessionid = s.ID
				s.Unlock()

				continue
			} // on reconnect

			return
		} // s.Wait()

		break
	}
}
