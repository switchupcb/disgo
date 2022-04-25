package requests

/// .go filetopics\Gateway.md
// Get Gateway
// GET /gateway
// https://discord.com/developers/docs/topics/gateway#get-gateway
type GetGateway struct {
	// TODO
}

// Get Gateway Bot
// GET /gateway/bot
// https://discord.com/developers/docs/topics/gateway#get-gateway-bot
type GetGatewayBot struct {
	URL               string `json:"url,omitempty"`
	Shards            int    `json:"shards,omitempty"`
	SessionStartLimit int    `json:"session_start_limit,omitempty"`
}
