package resources

// Audit Log Object
// https://discord.com/developers/docs/resources/audit-log
type AuditLog struct {
	AuditLogEntries      []*AuditLogEntry       `json:"audit_log_entries,omitempty"`
	GuildScheduledEvents []*GuildScheduledEvent `json:"guild_scheduled_events,omitempty"`
	Integration          []*Integration         `json:"integrations,omitempty"`
	Threads              []*Channel             `json:"threads,omitempty"`
	Users                []*User                `json:"users,omitempty"`
	Webhooks             []*Webhook             `json:"webhooks,omitempty"`
}

// Audit Log Entry Object
// https://discord.com/developers/docs/resources/audit-log#audit-log-object-audit-log-structure
type AuditLogEntry struct {
	TargetID   string            `json:"target_id,omitempty"`
	Changes    []*AuditLogChange `json:"changes,omitempty"`
	UserID     Snowflake         `json:"user_id,omitempty"`
	ID         Snowflake         `json:"id,omitempty"`
	ActionType Flag              `json:"action_type,omitempty"`
	Options    *AuditLogOptions  `json:"options,omitempty"`
	Reason     *string           `json:"reason,omitempty"`
}

// Audit Log Events
// https://discord.com/developers/docs/resources/audit-log#audit-log-entry-object-audit-log-events
const (
	FlagEventsLogAuditGUILD_UPDATE                 = 1
	FlagEventsLogAuditCHANNEL_CREATE               = 10
	FlagEventsLogAuditCHANNEL_UPDATE               = 11
	FlagEventsLogAuditCHANNEL_DELETE               = 12
	FlagEventsLogAuditCHANNEL_OVERWRITE_CREATE     = 13
	FlagEventsLogAuditCHANNEL_OVERWRITE_UPDATE     = 14
	FlagEventsLogAuditCHANNEL_OVERWRITE_DELETE     = 15
	FlagEventsLogAuditMEMBER_KICK                  = 20
	FlagEventsLogAuditMEMBER_PRUNE                 = 21
	FlagEventsLogAuditMEMBER_BAN_ADD               = 22
	FlagEventsLogAuditMEMBER_BAN_REMOVE            = 23
	FlagEventsLogAuditMEMBER_UPDATE                = 24
	FlagEventsLogAuditMEMBER_ROLE_UPDATE           = 25
	FlagEventsLogAuditMEMBER_MOVE                  = 26
	FlagEventsLogAuditMEMBER_DISCONNECT            = 27
	FlagEventsLogAuditBOT_ADD                      = 28
	FlagEventsLogAuditROLE_CREATE                  = 30
	FlagEventsLogAuditROLE_UPDATE                  = 31
	FlagEventsLogAuditROLE_DELETE                  = 32
	FlagEventsLogAuditINVITE_CREATE                = 40
	FlagEventsLogAuditINVITE_UPDATE                = 41
	FlagEventsLogAuditINVITE_DELETE                = 42
	FlagEventsLogAuditWEBHOOK_CREATE               = 50
	FlagEventsLogAuditWEBHOOK_UPDATE               = 51
	FlagEventsLogAuditWEBHOOK_DELETE               = 52
	FlagEventsLogAuditEMOJI_CREATE                 = 60
	FlagEventsLogAuditEMOJI_UPDATE                 = 61
	FlagEventsLogAuditEMOJI_DELETE                 = 62
	FlagEventsLogAuditMESSAGE_DELETE               = 72
	FlagEventsLogAuditMESSAGE_BULK_DELETE          = 73
	FlagEventsLogAuditMESSAGE_PIN                  = 74
	FlagEventsLogAuditMESSAGE_UNPIN                = 75
	FlagEventsLogAuditINTEGRATION_CREATE           = 80
	FlagEventsLogAuditINTEGRATION_UPDATE           = 81
	FlagEventsLogAuditINTEGRATION_DELETE           = 82
	FlagEventsLogAuditSTAGE_INSTANCE_CREATE        = 83
	FlagEventsLogAuditSTAGE_INSTANCE_UPDATE        = 84
	FlagEventsLogAuditSTAGE_INSTANCE_DELETE        = 85
	FlagEventsLogAuditSTICKER_CREATE               = 90
	FlagEventsLogAuditSTICKER_UPDATE               = 91
	FlagEventsLogAuditSTICKER_DELETE               = 92
	FlagEventsLogAuditGUILD_SCHEDULED_EVENT_CREATE = 100
	FlagEventsLogAuditGUILD_SCHEDULED_EVENT_UPDATE = 101
	FlagEventsLogAuditGUILD_SCHEDULED_EVENT_DELETE = 102
	FlagEventsLogAuditTHREAD_CREATE                = 110
	FlagEventsLogAuditTHREAD_UPDATE                = 111
	FlagEventsLogAuditTHREAD_DELETE                = 112
)

// Optional Audit Entry Info
// https://discord.com/developers/docs/resources/audit-log#audit-log-entry-object-optional-audit-entry-info
type AuditLogOptions struct {
	ChannelID        Snowflake `json:"channel_id,omitempty"`
	Count            string    `json:"count,omitempty"`
	DeleteMemberDays string    `json:"delete_member_days,omitempty"`
	ID               Snowflake `json:"id,omitempty"`
	MembersRemoved   string    `json:"members_removed,omitempty"`
	MessageID        Snowflake `json:"message_id,omitempty"`
	RoleName         string    `json:"role_name,omitempty"`
}

// Audit Log Change Object
// https://discord.com/developers/docs/resources/audit-log#audit-log-change-object
type AuditLogChange struct {
	NewValue interface{} `json:"new_value,omitempty"`
	OldValue interface{} `json:"old_value,omitempty"`
	Key      string      `json:"key,omitempty"`
}

// Audit Log Change Key
// https://discord.com/developers/docs/resources/audit-log#audit-log-change-object-audit-log-change-key
const (
	FlagKeyChangeLogAuditafk_channel_id                = "afk_channel_id"
	FlagKeyChangeLogAuditafk_timeout                   = "afk_timeout"
	FlagKeyChangeLogAuditallow                         = "allow"
	FlagKeyChangeLogAuditapplication_id                = "application_id"
	FlagKeyChangeLogAuditarchived                      = "archived"
	FlagKeyChangeLogAuditasset                         = "asset"
	FlagKeyChangeLogAuditauto_archive_duration         = "auto_archive_duration"
	FlagKeyChangeLogAuditavailable                     = "available"
	FlagKeyChangeLogAuditavatar_hash                   = "avatar_hash"
	FlagKeyChangeLogAuditbanner_hash                   = "banner_hash"
	FlagKeyChangeLogAuditbitrate                       = "bitrate"
	FlagKeyChangeLogAuditchannel_id                    = "channel_id"
	FlagKeyChangeLogAuditcode                          = "code"
	FlagKeyChangeLogAuditcolor                         = "color"
	FlagKeyChangeLogAuditcommunication_disabled_until  = "communication_disabled_until"
	FlagKeyChangeLogAuditdeaf                          = "deaf"
	FlagKeyChangeLogAuditdefault_auto_archive_duration = "default_auto_archive_duration"
	FlagKeyChangeLogAuditdefault_message_notifications = "default_message_notifications"
	FlagKeyChangeLogAuditdeny                          = "deny"
	FlagKeyChangeLogAuditdescription                   = "description"
	FlagKeyChangeLogAuditdiscovery_splash_hash         = "discovery_splash_hash"
	FlagKeyChangeLogAuditenable_emoticons              = "enable_emoticons"
	FlagKeyChangeLogAuditentity_type                   = "entity_type"
	FlagKeyChangeLogAuditexpire_behavior               = "expire_behavior"
	FlagKeyChangeLogAuditexpire_grace_period           = "expire_grace_period"
	FlagKeyChangeLogAuditexplicit_content_filter       = "explicit_content_filter"
	FlagKeyChangeLogAuditformat_type                   = "format_type"
	FlagKeyChangeLogAuditguild_id                      = "guild_id"
	FlagKeyChangeLogAudithoist                         = "hoist"
	FlagKeyChangeLogAuditicon_hash                     = "icon_hash"
	FlagKeyChangeLogAuditimage_hash                    = "image_hash"
	FlagKeyChangeLogAuditid                            = "id"
	FlagKeyChangeLogAuditinvitable                     = "invitable"
	FlagKeyChangeLogAuditinviter_id                    = "inviter_id"
	FlagKeyChangeLogAuditlocation                      = "location"
	FlagKeyChangeLogAuditlocked                        = "locked"
	FlagKeyChangeLogAuditmax_age                       = "max_age"
	FlagKeyChangeLogAuditmax_uses                      = "max_uses"
	FlagKeyChangeLogAuditmentionable                   = "mentionable"
	FlagKeyChangeLogAuditmfa_level                     = "mfa_level"
	FlagKeyChangeLogAuditmute                          = "mute"
	FlagKeyChangeLogAuditname                          = "name"
	FlagKeyChangeLogAuditnick                          = "nick"
	FlagKeyChangeLogAuditnsfw                          = "nsfw"
	FlagKeyChangeLogAuditowner_id                      = "owner_id"
	FlagKeyChangeLogAuditpermission_overwrites         = "permission_overwrites"
	FlagKeyChangeLogAuditpermissions                   = "permissions"
	FlagKeyChangeLogAuditposition                      = "position"
	FlagKeyChangeLogAuditpreferred_locale              = "preferred_locale"
	FlagKeyChangeLogAuditprivacy_level                 = "privacy_level"
	FlagKeyChangeLogAuditprune_delete_days             = "prune_delete_days"
	FlagKeyChangeLogAuditpublic_updates_channel_id     = "public_updates_channel_id"
	FlagKeyChangeLogAuditrate_limit_per_user           = "rate_limit_per_user"
	FlagKeyChangeLogAuditregion                        = "region"
	FlagKeyChangeLogAuditrules_channel_id              = "rules_channel_id"
	FlagKeyChangeLogAuditsplash_hash                   = "splash_hash"
	FlagKeyChangeLogAuditstatus                        = "status"
	FlagKeyChangeLogAuditsystem_channel_id             = "system_channel_id"
	FlagKeyChangeLogAudittags                          = "tags"
	FlagKeyChangeLogAudittemporary                     = "temporary"
	FlagKeyChangeLogAudittopic                         = "topic"
	FlagKeyChangeLogAudittype                          = "type"
	FlagKeyChangeLogAuditunicode_emoji                 = "unicode_emoji"
	FlagKeyChangeLogAudituser_limit                    = "user_limit"
	FlagKeyChangeLogAudituses                          = "uses"
	FlagKeyChangeLogAuditvanity_url_code               = "vanity_url_code"
	FlagKeyChangeLogAuditverification_level            = "verification_level"
	FlagKeyChangeLogAuditwidget_channel_id             = "widget_channel_id"
	FlagKeyChangeLogAuditwidget_enabled                = "widget_enabled"
	FlagKeyChangeLogAuditadd                           = "add"
	FlagKeyChangeLogAuditremove                        = "remove"
)
