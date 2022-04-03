package resources

// Audit Log Object
// https://discord.com/developers/docs/resources/audit-log
type AuditLog struct {
	AuditLogEntries      []*AuditLogEntry       `json:"audit_log_entries"`
	GuildScheduledEvents []*GuildScheduledEvent `json:"guild_scheduled_events"`
	Integration          []*Integration         `json:"integrations"`
	Threads              []*Channel             `json:"threads"`
	Users                []*User                `json:"users"`
	Webhooks             []*Webhook             `json:"webhooks"`
}

// Audit Log Entry Object
// https://discord.com/developers/docs/resources/audit-log#audit-log-object-audit-log-structure
type AuditLogEntry struct {
	TargetID   string            `json:"target_id"`
	Changes    []*AuditLogChange `json:"changes"`
	UserID     string            `json:"user_id"`
	ID         string            `json:"id"`
	ActionType uint8             `json:"action_type"`
	Options    *AuditLogOptions  `json:"options"`
	Reason     string            `json:"reason"`
}

// Audit Log Events
// https://discord.com/developers/docs/resources/audit-log#audit-log-entry-object-audit-log-events
const (
	GUILD_UPDATE                 = 1
	CHANNEL_CREATE               = 10
	CHANNEL_UPDATE               = 11
	CHANNEL_DELETE               = 12
	CHANNEL_OVERWRITE_CREATE     = 13
	CHANNEL_OVERWRITE_UPDATE     = 14
	CHANNEL_OVERWRITE_DELETE     = 15
	MEMBER_KICK                  = 20
	MEMBER_PRUNE                 = 21
	MEMBER_BAN_ADD               = 22
	MEMBER_BAN_REMOVE            = 23
	MEMBER_UPDATE                = 24
	MEMBER_ROLE_UPDATE           = 25
	MEMBER_MOVE                  = 26
	MEMBER_DISCONNECT            = 27
	BOT_ADD                      = 28
	ROLE_CREATE                  = 30
	ROLE_UPDATE                  = 31
	ROLE_DELETE                  = 32
	INVITE_CREATE                = 40
	INVITE_UPDATE                = 41
	INVITE_DELETE                = 42
	WEBHOOK_CREATE               = 50
	WEBHOOK_UPDATE               = 51
	WEBHOOK_DELETE               = 52
	EMOJI_CREATE                 = 60
	EMOJI_UPDATE                 = 61
	EMOJI_DELETE                 = 62
	MESSAGE_DELETE               = 72
	MESSAGE_BULK_DELETE          = 73
	MESSAGE_PIN                  = 74
	MESSAGE_UNPIN                = 75
	INTEGRATION_CREATE           = 80
	INTEGRATION_UPDATE           = 81
	INTEGRATION_DELETE           = 82
	STAGE_INSTANCE_CREATE        = 83
	STAGE_INSTANCE_UPDATE        = 84
	STAGE_INSTANCE_DELETE        = 85
	STICKER_CREATE               = 90
	STICKER_UPDATE               = 91
	STICKER_DELETE               = 92
	GUILD_SCHEDULED_EVENT_CREATE = 100
	GUILD_SCHEDULED_EVENT_UPDATE = 101
	GUILD_SCHEDULED_EVENT_DELETE = 102
	THREAD_CREATE                = 110
	THREAD_UPDATE                = 111
	THREAD_DELETE                = 112
)

// Optional Audit Entry Info
// https://discord.com/developers/docs/resources/audit-log#audit-log-entry-object-optional-audit-entry-info
type AuditLogOptions struct {
	ChannelID        string `json:"channel_id"`
	Count            string `json:"count"`
	DeleteMemberDays string `json:"delete_member_days"`
	ID               string `json:"id"`
	MembersRemoved   string `json:"members_removed"`
	MessageID        string `json:"message_id"`
	RoleName         string `json:"role_name"`
}

// Audit Log Change Object
// https://discord.com/developers/docs/resources/audit-log#audit-log-change-object
type AuditLogChange struct {
	NewValue interface{} `json:"new_value"`
	OldValue interface{} `json:"old_value"`
	Key      string      `json:"key"`
}

// Audit Log Change Key
// https://discord.com/developers/docs/resources/audit-log#audit-log-change-object-audit-log-change-key
const (
	KeyAuditLogafk_channel_id                = "afk_channel_id"
	KeyAuditLogafk_timeout                   = "afk_timeout"
	KeyAuditLogallow                         = "allow"
	KeyAuditLogapplication_id                = "application_id"
	KeyAuditLogarchived                      = "archived"
	KeyAuditLogasset                         = "asset"
	KeyAuditLogauto_archive_duration         = "auto_archive_duration"
	KeyAuditLogavailable                     = "available"
	KeyAuditLogavatar_hash                   = "avatar_hash"
	KeyAuditLogbanner_hash                   = "banner_hash"
	KeyAuditLogbitrate                       = "bitrate"
	KeyAuditLogchannel_id                    = "channel_id"
	KeyAuditLogcode                          = "code"
	KeyAuditLogcolor                         = "color"
	KeyAuditLogcommunication_disabled_until  = "communication_disabled_until"
	KeyAuditLogdeaf                          = "deaf"
	KeyAuditLogdefault_auto_archive_duration = "default_auto_archive_duration"
	KeyAuditLogdefault_message_notifications = "default_message_notifications"
	KeyAuditLogdeny                          = "deny"
	KeyAuditLogdescription                   = "description"
	KeyAuditLogdiscovery_splash_hash         = "discovery_splash_hash"
	KeyAuditLogenable_emoticons              = "enable_emoticons"
	KeyAuditLogentity_type                   = "entity_type"
	KeyAuditLogexpire_behavior               = "expire_behavior"
	KeyAuditLogexpire_grace_period           = "expire_grace_period"
	KeyAuditLogexplicit_content_filter       = "explicit_content_filter"
	KeyAuditLogformat_type                   = "format_type"
	KeyAuditLogguild_id                      = "guild_id"
	KeyAuditLoghoist                         = "hoist"
	KeyAuditLogicon_hash                     = "icon_hash"
	KeyAuditLogimage_hash                    = "image_hash"
	KeyAuditLogid                            = "id"
	KeyAuditLoginvitable                     = "invitable"
	KeyAuditLoginviter_id                    = "inviter_id"
	KeyAuditLoglocation                      = "location"
	KeyAuditLoglocked                        = "locked"
	KeyAuditLogmax_age                       = "max_age"
	KeyAuditLogmax_uses                      = "max_uses"
	KeyAuditLogmentionable                   = "mentionable"
	KeyAuditLogmfa_level                     = "mfa_level"
	KeyAuditLogmute                          = "mute"
	KeyAuditLogname                          = "name"
	KeyAuditLognick                          = "nick"
	KeyAuditLognsfw                          = "nsfw"
	KeyAuditLogowner_id                      = "owner_id"
	KeyAuditLogpermission_overwrites         = "permission_overwrites"
	KeyAuditLogpermissions                   = "permissions"
	KeyAuditLogposition                      = "position"
	KeyAuditLogpreferred_locale              = "preferred_locale"
	KeyAuditLogprivacy_level                 = "privacy_level"
	KeyAuditLogprune_delete_days             = "prune_delete_days"
	KeyAuditLogpublic_updates_channel_id     = "public_updates_channel_id"
	KeyAuditLograte_limit_per_user           = "rate_limit_per_user"
	KeyAuditLogregion                        = "region"
	KeyAuditLogrules_channel_id              = "rules_channel_id"
	KeyAuditLogsplash_hash                   = "splash_hash"
	KeyAuditLogstatus                        = "status"
	KeyAuditLogsystem_channel_id             = "system_channel_id"
	KeyAuditLogtags                          = "tags"
	KeyAuditLogtemporary                     = "temporary"
	KeyAuditLogtopic                         = "topic"
	KeyAuditLogtype                          = "type"
	KeyAuditLogunicode_emoji                 = "unicode_emoji"
	KeyAuditLoguser_limit                    = "user_limit"
	KeyAuditLoguses                          = "uses"
	KeyAuditLogvanity_url_code               = "vanity_url_code"
	KeyAuditLogverification_level            = "verification_level"
	KeyAuditLogwidget_channel_id             = "widget_channel_id"
	KeyAuditLogwidget_enabled                = "widget_enabled"
	KeyAuditLogadd                           = "add"
	KeyAuditLogremove                        = "remove"
)
