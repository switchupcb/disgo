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
	statusCode                = 1000
)

// Session represents a Discord Gateway WebSocket Session.
type Session struct {
	// Token represents the Bot token the session used to connect.
	Token string

	// ID represents the session ID of the session.
	ID string

	// Seq represents the last sequence number received by the client.
	//
	// https://discord.com/developers/docs/topics/gateway#heartbeat
	Seq int

	// Endpoint represents the endpoint that will be connected to.
	Endpoint string

	// Connected represented whether or not the session is connected to Discord Gateway.
	Connected bool
}

// Connect creates an open connection to Discord Gateway.
func (bot *Client) Connect(s *Session) error {
	// requesting a valid Gateway URL endpoint from the API.
	gateway := GetGateway{}
	response, err := gateway.Send(bot)
	if err != nil {
		return fmt.Errorf("an error occurred getting the Gateway API URL\n%w", err)
	}

	// TODO: zlib compression
	endpoint := response.URL + "?v=" + VersionDiscordAPI + "&encoding=" + gatewayEncoding

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	// connect to the Discord Gateway Websocket.
	conn, _, err := websocket.Dial(ctx, endpoint, nil)
	if err != nil {
		return fmt.Errorf("an error occurred performing WebSocket handshake\n%w", err)
	}

	defer conn.Close(websocket.StatusInternalError, "StatusInternalError")
	s.Connected = true

	// handle the incoming Hello event (JSON) upon connecting to the Gateway.
	var hello Hello
	err = wsjson.Read(ctx, conn, &hello)
	if err != nil {
		return fmt.Errorf(ErrEventUnmarshal, "Hello", err)
	}

	// begin sending heartbeat payloads every heartbeat_interval ms
	// go send heartbeat
	// TODO: deal with heartbeat when listening for events (2nd paragraph)

	// send an Identify event to the Discord Gateway.
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

	// handle the incoming Ready event upon identification with the socket.
	var ready Ready
	err = wsjson.Read(ctx, conn, &ready)
	if err != nil {
		return fmt.Errorf(ErrEventUnmarshal, "Ready", err)
	}

	// go listen for events
	return nil
}

// Reconn restablishes the connection to the Discord Gateway.
func (s *Session) Reconn() error {

	if s.Connected {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	// connect to Discord Gateway Websocket.
	conn, _, err := websocket.Dial(ctx, s.Endpoint, nil)
	if err != nil {
		return err
	}

	// send a Resume event to the Discord Gateway.
	event := Resume{
		Token:     s.Token,
		SessionID: s.ID,
		Seq:       uint32(s.Seq),
	}

	err = wsjson.Write(ctx, conn, event)
	if err != nil {
		return fmt.Errorf(ErrEventUnmarshal, "Resume", err)
	}

	// read in the Resumed event to access its SessionID
	// TODO: handle all of the events in order
	// TODO: handle opcode 9
	var resumed Resumed
	err = wsjson.Read(ctx, conn, &resumed)
	if err != nil {
		return fmt.Errorf(ErrEventUnmarshal, "Resumed", err)
	}
	if resumed.Op != 7 {
		return fmt.Errorf(ErrEventUnmarshal, "Ready", err)
	}

	return nil
}

// Terminate disconnects from the Discord Gateway by sending a status code 1000.
func Terminate(c *websocket.Conn) error {

	c.Close(statusCode, "Close TCP Connection")

	return nil
}
