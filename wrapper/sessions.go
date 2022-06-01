package wrapper

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

const (
	gatewayEncoding           = "json"
	module                    = "github.com/switchupcb/disgo"
	maxIdentifyLargeThreshold = 250
	ErrEventUnmarshal         = "an error occurred while unmarshalling a %v Event:\n%w"
)

// Session represents a session.
type Session struct {
	// Token represents the Bot token the session used to connect.
	Token string

	// ID represents the session ID of the session.
	ID string

	// Seq represents the last sequence number received by the client.
	//
	// https://discord.com/developers/docs/topics/gateway#heartbeat
	Seq int
}

// Connect creates open connection to Discord Gateway.
func (bot *Client) Connect(s *Session) error {
	// get the Gateway URL.
	gateway := GetGateway{}
	response, err := gateway.Send(bot)
	if err != nil {
		return fmt.Errorf("an error occurred getting the Gateway API URL\n%w", err)
	}

	// TODO: zlib compression
	endpoint := response.URL + "?v=" + VersionDiscordAPI + "&encoding=" + gatewayEncoding

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	// connect to Discord Gateway Websocket.
	conn, _, err := websocket.Dial(ctx, endpoint, nil)
	if err != nil {
		return fmt.Errorf("an error occurred performing WebSocket handshake\n%w", err)
	}

	defer conn.Close(websocket.StatusInternalError, "StatusInternalError")

	// unmarshall json file bytes into hello event object
	// handle the Hello event
	var hello Hello
	err = wsjson.Read(ctx, conn, &hello)
	if err != nil {
		return fmt.Errorf(ErrEventUnmarshal, "Hello", err)
	}

	// begin sending heartbeat payloads every heartbeat_interval ms
	// go send heartbeat
	// TODO: deal with heartbeat when listening for events (2nd paragraph)

	// deal with opcode 2
	event := Identify{
		Token: bot.Authentication.Token,
		Properties: IdentifyConnectionProperties{
			OS:      runtime.GOOS,
			Browser: module,
			Device:  module,
		},
		Compress:       true, // TODO: account for compression
		LargeThreshold: maxIdentifyLargeThreshold,
		Shard:          nil, // TODO: sharding
		Presence:       *bot.Config.GatewayPresenceUpdate,
		Intents:        bot.Config.Intents,
	}

	wsjson.Write(ctx, conn, event)

	// Reading in Ready event.
	var ready Ready
	err = wsjson.Read(ctx, conn, &ready)
	if err != nil {
		return fmt.Errorf(ErrEventUnmarshal, "Ready", err)
	}

	// go listen for events
	return nil
}

func (s *Session) Reconn(s *Session) {

	event := Resume{
		Token:     s.Token,
		SessionID: s.ID,
		Seq:       uint32(s.Seq),
	}
}

// find out how to terminate connection
