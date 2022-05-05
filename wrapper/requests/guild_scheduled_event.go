package requests

import "github.com/switchupcb/disgo/wrapper/resources"

// List Scheduled Events for Guild
// GET /guilds/{guild.id}/scheduled-events
// https://discord.com/developers/docs/resources/guild-scheduled-event#list-scheduled-events-for-guild
type ListScheduledEventsforGuild struct {
	WithUserCount bool `json:"with_user_count,omitempty"`
}

// Create Guild Scheduled Event
// POST /guilds/{guild.id}/scheduled-events
// https://discord.com/developers/docs/resources/guild-scheduled-event#create-guild-scheduled-event
type CreateGuildScheduledEvent struct {
	ChannelID          *resources.Snowflake                         `json:"channel_id,omitempty"`
	EntityMetadata     *resources.GuildScheduledEventEntityMetadata `json:"entity_metadata,omitempty"`
	Name               *string                                      `json:"name,omitempty"`
	PrivacyLevel       resources.Flag                               `json:"privacy_level,omitempty"`
	ScheduledStartTime resources.Snowflake                          `json:"scheduled_start_time,omitempty"`
	ScheduledEndTime   resources.Snowflake                          `json:"scheduled_end_time,omitempty"`
	Description        *string                                      `json:"description,omitempty"`
	EntityType         *resources.Flag                              `json:"entity_type,omitempty"`
	Image              *string                                      `json:"image,omitempty"`
}

// Get Guild Scheduled Event
// GET /guilds/{guild.id}/scheduled-events/{guild_scheduled_event.id}
// https://discord.com/developers/docs/resources/guild-scheduled-event#get-guild-scheduled-event
type GetGuildScheduledEvent struct {
	GuildID               resources.Snowflake
	GuildScheduledEventID resources.Snowflake
	WithUserCount         bool `json:"with_user_count,omitempty"`
}

// Modify Guild Scheduled Event
// PATCH /guilds/{guild.id}/scheduled-events/{guild_scheduled_event.id}
// https://discord.com/developers/docs/resources/guild-scheduled-event#modify-guild-scheduled-event
type ModifyGuildScheduledEvent struct {
	GuildID               resources.Snowflake
	GuildScheduledEventID resources.Snowflake
	ChannelID             *resources.Snowflake                         `json:"channel_id,omitempty"`
	EntityMetadata        *resources.GuildScheduledEventEntityMetadata `json:"entity_metadata,omitempty"`
	Name                  *string                                      `json:"name,omitempty"`
	PrivacyLevel          resources.Flag                               `json:"privacy_level,omitempty"`
	ScheduledStartTime    resources.Snowflake                          `json:"scheduled_start_time,omitempty"`
	ScheduledEndTime      resources.Snowflake                          `json:"scheduled_end_time,omitempty"`
	Description           *string                                      `json:"description,omitempty"`
	EntityType            *resources.Flag                              `json:"entity_type,omitempty"`
	Image                 *string                                      `json:"image,omitempty"`
	Status                resources.Flag                               `json:"status,omitempty"`
}

// Delete Guild Scheduled Event
// DELETE /guilds/{guild.id}/scheduled-events/{guild_scheduled_event.id}
// https://discord.com/developers/docs/resources/guild-scheduled-event#delete-guild-scheduled-event
type DeleteGuildScheduledEvent struct {
	GuildID               resources.Snowflake
	GuildScheduledEventID resources.Snowflake
}

// Get Guild Scheduled Event Users
// GET /guilds/{guild.id}/scheduled-events/{guild_scheduled_event.id}/users
// https://discord.com/developers/docs/resources/guild-scheduled-event#get-guild-scheduled-event-users
type GetGuildScheduledEventUsers struct {
	GuildID               resources.Snowflake
	GuildScheduledEventID resources.Snowflake
	Limit                 resources.Flag       `json:"limit,omitempty"`
	WithMember            bool                 `json:"with_member,omitempty"`
	Before                *resources.Snowflake `json:"before,omitempty"`
	After                 *resources.Snowflake `json:"after,omitempty"`
}
