package resources

// Stage Instance Object
// https://discord.com/developers/docs/resources/stage-instance#stage-instance-object
type StageInstance struct {
	ID                    Snowflake `json:"id,omitempty"`
	GuildID               Snowflake `json:"guild_id,omitempty"`
	ChannelID             Snowflake `json:"channel_id,omitempty"`
	Topic                 string    `json:"topic,omitempty"`
	PrivacyLevel          Flag      `json:"privacy_level,omitempty"`
	DiscoverableDisabled  bool      `json:"discoverable_disabled,omitempty"`
	GuildScheduledEventID Snowflake `json:"guild_scheduled_event_id,omitempty"`
}

// Privacy Level
// https://discord.com/developers/docs/resources/stage-instance#stage-instance-object-privacy-level
const (
	FlagPrivacyLevelPUBLIC     = 1
	FlagPrivacyLevelGUILD_ONLY = 2
)
