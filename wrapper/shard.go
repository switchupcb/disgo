package wrapper

import "time"

// ShardManager represents an interface for Shard Management.
//
// ShardManager is an interface which allows developers to use multi-application architectures,
// which run multiple applications on separate processes or servers.
type ShardManager interface {
	// SetNumShards sets the number of shards the shard manager will use.
	//
	// When the Shards = 0, the automatic shard manager is used.
	SetNumShards(shards int)

	// SetLimit sets the ShardLimit of the ShardManager.
	//
	// This limit is determined using the GetGatewayBot request (which provides the Gateway Endpoint).
	// https://discord.com/developers/docs/topics/gateway#get-gateway-bot
	//
	// Called from the session.go connect() function (at L#123 in /wrapper/session.go).
	SetLimit(bot *Client) (gatewayEndpoint string, response *GetGatewayBotResponse, err error)

	// GetSessions gets the connected sessions of the bot (in order of connection).
	GetSessions() []*Session

	// Ready is called when a Session receives a ready event.
	//
	// Called from the session.go initial() function (at L#304 in /wrapper/session.go).
	Ready(bot *Client, session *Session, event *Ready)

	// Connect connects to the Discord Gateway using the Shard Manager.
	Connect(bot *Client) error

	// Disconnect disconnects from the Discord Gateway using the Shard Manager.
	Disconnect() error

	// Reconnect reconnects to the Discord Gateway using the Shard Manager.
	Reconnect(bot *Client) error
}

// ShardLimit contains information about sharding limits.
type ShardLimit struct {
	// Reset represents the time at which the Session Start Rate Limit resets (daily).
	//
	// Discord represents this value from the "reset_after" field of the SessionStartLimit object.
	// https://discord.com/developers/docs/topics/gateway#session-start-limit-object
	Reset time.Time

	// MaxStarts represents the maximum amount of WebSocket Sessions a bot can start per day.
	//
	// This is equivalent to the maximum amount of Shards a bot can create per day.
	//
	// Discord represents this value from the "total" field of the SessionStartLimit object.
	// https://discord.com/developers/docs/topics/gateway#session-start-limit-object
	MaxStarts int

	// RemainingStarts represents the remaining number of "starts" that the bot is allowed
	// until the reset time.
	//
	// This is equivalent to the remaining number of Shards that the bot can create
	// for the rest of the day.
	//
	// Discord represents this value from the "remaining" field of the SessionStartLimit object.
	// https://discord.com/developers/docs/topics/gateway#session-start-limit-object
	RemainingStarts int

	// MaxConcurrency represents the number of Identify SendEvents the bot can send every 5 seconds.
	MaxConcurrency int

	// RecommendedShards represents the number of shards to use when connecting.
	//
	// https://discord.com/developers/docs/topics/gateway#get-gateway-bot
	RecommendedShards int
}
