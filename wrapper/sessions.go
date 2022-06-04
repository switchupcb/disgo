package wrapper

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"time"

	json "github.com/goccy/go-json"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

const (
	module                    = "github.com/switchupcb/disgo"
	gatewayEncoding           = "json"
	gatewayCloseStatusCode    = 1000
	maxIdentifyLargeThreshold = 250
)

// Session represents a Discord Gateway WebSocket Session.
type Session struct {
	// Token represents the bot token the session used to connect to the Gateway.
	Token string

	// ID represents the session ID of the Session.
	ID string

	// Seq represents the last sequence number received by the client.
	//
	// https://discord.com/developers/docs/topics/gateway#heartbeat
	Seq int

	// Endpoint represents the endpoint that is used to connect to the Gateway.
	Endpoint string

	// Connected represents whether the session is connected to Discord Gateway.
	Connected bool
}

// Connect creates an open connection to Discord Gateway.
func (bot *Client) Connect(s *Session) error {
	// request a valid Gateway URL endpoint from the API.
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
		return fmt.Errorf("an error occurred connecting to the Discord Gateway\n%w", err)
	}

	defer conn.Close(websocket.StatusInternalError, "StatusInternalError")
	s.Connected = true

	// handle the incoming Hello event (JSON) upon connecting to the Gateway.
	// TODO: Replace with bot.onEvent()
	var hello Hello
	err = wsjson.Read(ctx, conn, &hello)
	if err != nil {
		return fmt.Errorf(ErrEventUnmarshal, "Hello", err)
	}

	// begin sending heartbeat payloads every heartbeat_interval ms
	// go send heartbeat
	// go listen for events
	// TODO: deal with heartbeat by starting to listen for events.

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
	// TODO: Replace with bot.onEvent()
	var ready Ready
	err = wsjson.Read(ctx, conn, &ready)
	if err != nil {
		return fmt.Errorf(ErrEventUnmarshal, "Ready", err)
	}

	return nil
}

// reconnect reestablishes a session's connection to the Discord Gateway.
func (s *Session) reconnect() error {
	if s.Connected {
		return nil
	}

	// connect to Discord Gateway Websocket.
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

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

	// read in the Resumed event to access its SessionID.
	// TODO: handle all of the events in order using bot.onEvent
	// TODO: handle opcode 9 using bot.onEvent
	var resumed Resumed
	err = wsjson.Read(ctx, conn, &resumed)
	if err != nil {
		return fmt.Errorf(ErrEventUnmarshal, "Resumed", err)
	}

	// TODO: determine if this check is necessary given bot.onEvent
	if resumed.Op != FlagGatewayOpcodeReconnect {
		return fmt.Errorf(ErrEventUnmarshal, "Ready", err)
	}

	return nil
}

// Disconnect disconnects a session from the Discord Gateway.
func (s *Session) Disconnect(c *websocket.Conn) error {
	// TODO: Ensure this specific session is disconnected.
	// For example, storing conn in the session to use later.
	c.Close(gatewayCloseStatusCode, fmt.Sprintf("Disconnected Session %s from the Discord Gateway.", s.ID))

	return nil
}

// onEvent handles an event using its JSON data.
func (bot *Client) onEvent(s *Session, data []byte) error {
	var event GatewayPayload
	err := json.Unmarshal(data, &event)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	// run the bot's event handlers based on the Receive Opcode.
	//
	// https://discord.com/developers/docs/topics/opcodes-and-status-codes#gateway-gateway-opcodes
	switch *event.Op {
	case FlagGatewayOpcodeDispatch:
		s.Seq = event.SequenceNumber
		go bot.handle(event.EventName, event.Data)

	case FlagGatewayOpcodeHello:

	case FlagGatewayOpcodeHeartbeat:

	case FlagGatewayOpcodeHeartbeatACK:

	case FlagGatewayOpcodeReconnect:

	case FlagGatewayOpcodeInvalidSession:

	}

	return nil
}

// TODO: Automatically generate the following code using copygen.
// Handlers represents a bot's event handlers.
type Handlers struct {
	Hello          []func(*Hello)
	Ready          []func(*Ready)
	Resumed        []func(*Resumed)
	Reconnect      []func(*Reconnect)
	InvalidSession []func(*InvalidSession)
}

// handle handles an event using its name and JSON data.
func (bot *Client) handle(name string, data json.RawMessage) {
	switch name {
	case FlagGatewayEventNameHello:
		var event *Hello
		err := json.Unmarshal(data, event)
		if err != nil {
			// TODO: fix goroutine error handling semantics.
			log.Panicf("%v", err)
		}

		for _, handler := range bot.Handlers.Hello {
			go handler(event)
		}

	case FlagGatewayEventNameReady:
	case FlagGatewayEventNameResumed:
	case FlagGatewayEventNameReconnect:
	case FlagGatewayEventNameInvalidSession:
	}
}
