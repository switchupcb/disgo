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
	// Connected is a channel because it needs to be used in a select statement
	// when it is nil it is false.
	Connected chan bool

	// Context is used to create the Gateway Connection.
	Context context.Context

	// Conn represents a connection to the Discord Gateway at the specified Endpoint.
	Conn *websocket.Conn

	// HeartbeatInterval represents the interval of time between Heartbeat Payloads
	// being sent to the Discord Gateway and should be converted to miliseconds.
	HeartbeatInterval time.Duration

	// Ticker represents a timer and is used when tasks must be completed within specified intervals.
	// Ticker is used to properly time the intervals in which Heartbeat Payloads are sent.
	Ticker *time.Ticker

	// LastHeartbeatACK represents the time that the last ACK was received by the Gateway and is used to
	// determine whether the Connection has zombied.
	LastHeartbeatACK time.Time
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
	s.Connected = make(chan bool)

	// handle the incoming Hello event (JSON) upon connecting to the Gateway.
	var hello Hello
	err = wsjson.Read(ctx, conn, &hello)
	if err != nil {
		return fmt.Errorf(ErrEventUnmarshal, "Hello", err)
	}

	// begin sending heartbeat payloads every heartbeat_interval ms.
	go bot.heartbeat(s)
	go bot.listen(s)

	// send an Identify event to the Discord Gateway (Opcode 2).
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
	// TODO: Replace with bot.onPayload()
	var ready Ready
	err = wsjson.Read(ctx, conn, &ready)
	if err != nil {
		return fmt.Errorf(ErrEventUnmarshal, "Ready", err)
	}

	return nil
}

// reconnect reestablishes a session's connection to the Discord Gateway.
func (s *Session) reconnect() error {

	// TODO: fix select statement so below code is not unreacheable
	// need to check if the session is already connected.
	select {
	case <-s.Connected:
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
	// TODO: handle all of the events in order using bot.onPayload

	var resumed Resumed
	err = wsjson.Read(ctx, conn, &resumed)
	if err != nil {
		return fmt.Errorf(ErrEventUnmarshal, "Resumed", err)
	}

	return nil
}

// Disconnect disconnects a session from the Discord Gateway.
func (s *Session) Disconnect(c *websocket.Conn) error {
	// TODO: Ensure this specific session is disconnected.
	// For example, storing conn in the session to use later.
	c.Close(gatewayCloseStatusCode, fmt.Sprintf("Disconnected Session %s from the Discord Gateway.", s.ID))
	s.Connected = nil

	return nil
}

// onPayload handles an Discord Gateway Payload.
func (bot *Client) onPayload(s *Session, data []byte) error {
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

	case FlagGatewayOpcodeHeartbeat:
		// Sequence number must be updated for when Opcode 1 Heartbeat is sent to Gateway.
		s.Seq = event.SequenceNumber

		hb := Heartbeat{
			Op:   FlagGatewayOpcodeHeartbeat,
			Data: s.Seq,
		}
		err := wsjson.Write(s.Context, s.Conn, hb)
		if err != nil {
			log.Printf("%v", err)
		}
		s.Ticker.Reset(s.HeartbeatInterval * time.Millisecond)

	case FlagGatewayOpcodeHeartbeatACK:
		// receive Hearbeat ACK from Gateway.
		// TODO: deal with not receiving ACK after first heartbeat
		s.LastHeartbeatACK = time.Now().UTC()

	case FlagGatewayOpcodeReconnect:
		s.reconnect()

	case FlagGatewayOpcodeInvalidSession:
		// when received Opcode 9 Invalid Session from gateway.
		// wait random amount of time between 1-5 seconds.
		const waitTime = 3
		timer := time.NewTimer(waitTime * time.Second)
		<-timer.C

		// send an Identify event to the Discord Gateway (Opcode 2).
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

		wsjson.Write(s.Context, s.Conn, event)
		if err != nil {
			log.Printf("%v", err)
		}

	}

	return nil
}

// heatbeat sends the payload to the Discord Gateway to verify the connection is alive.
func (bot *Client) heartbeat(s *Session) {

	//TODO: figure out when to use Mutex

	// Heartbeat is what the payload sends to the Gateway every HeartbeatInterval miliseconds.
	var hb Heartbeat
	hb.Op = FlagGatewayOpcodeHeartbeat

	s.Ticker = time.NewTicker(s.HeartbeatInterval * time.Millisecond)
	// Begin sending Opcode 1 Heartbeat Payloads
	for {
		hb.Data = s.Seq
		err := wsjson.Write(s.Context, s.Conn, hb)
		if err != nil {
			log.Printf("%v", err)
			return
		}

		select {
		case <-s.Ticker.C:
		case <-s.Connected:
			return
		}
	}
}

// listen listens to the connection for payloads from the Gateway.
func (bot *Client) listen(s *Session) {

	for {
		var payload []byte
		err := wsjson.Read(s.Context, s.Conn, payload)
		if err != nil {
			log.Printf("%v", err)
			return
		}

		select {
		case <-s.Connected:
			return
		default:
			bot.onPayload(s, payload)
		}
	}
}
