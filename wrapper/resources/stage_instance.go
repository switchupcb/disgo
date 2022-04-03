package resource

// Stage Instance Object
// https://discord.com/developers/docs/resources/stage-instance#stage-instance-object
type StageInstance struct {
	ID                    int64  `json:"id"`
	GuildID               string `json:"guild_id"`
	ChannelID             string `json:"channel_id"`
	Topic                 string `json:"topic"`
	PrivacyLevel          uint8  `json:"privacy_level"`
	DiscoverableDisabled  bool   `json:"discoverable_disabled"`
	GuildScheduledEventID int64  `json:"guild_scheduled_event_id"`
}

// Privacy Level
// https://discord.com/developers/docs/resources/stage-instance#stage-instance-object-privacy-level
const (
	PUBLIC     = 1
	GUILD_ONLY = 2
)
