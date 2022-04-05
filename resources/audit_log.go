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
	Reason     string            `json:"reason,omitempty"`
}

// Audit Log Events
// https://discord.com/developers/docs/resources/audit-log#audit-log-entry-object-audit-log-events
const (
	FlagAuditLogEventsGUILD_UPDATE                 = 1
	FlagAuditLogEventsCHANNEL_CREATE               = 10
	FlagAuditLogEventsCHANNEL_UPDATE               = 11
	FlagAuditLogEventsCHANNEL_DELETE               = 12
	FlagAuditLogEventsCHANNEL_OVERWRITE_CREATE     = 13
	FlagAuditLogEventsCHANNEL_OVERWRITE_UPDATE     = 14
	FlagAuditLogEventsCHANNEL_OVERWRITE_DELETE     = 15
	FlagAuditLogEventsMEMBER_KICK                  = 20
	FlagAuditLogEventsMEMBER_PRUNE                 = 21
	FlagAuditLogEventsMEMBER_BAN_ADD               = 22
	FlagAuditLogEventsMEMBER_BAN_REMOVE            = 23
	FlagAuditLogEventsMEMBER_UPDATE                = 24
	FlagAuditLogEventsMEMBER_ROLE_UPDATE           = 25
	FlagAuditLogEventsMEMBER_MOVE                  = 26
	FlagAuditLogEventsMEMBER_DISCONNECT            = 27
	FlagAuditLogEventsBOT_ADD                      = 28
	FlagAuditLogEventsROLE_CREATE                  = 30
	FlagAuditLogEventsROLE_UPDATE                  = 31
	FlagAuditLogEventsROLE_DELETE                  = 32
	FlagAuditLogEventsINVITE_CREATE                = 40
	FlagAuditLogEventsINVITE_UPDATE                = 41
	FlagAuditLogEventsINVITE_DELETE                = 42
	FlagAuditLogEventsWEBHOOK_CREATE               = 50
	FlagAuditLogEventsWEBHOOK_UPDATE               = 51
	FlagAuditLogEventsWEBHOOK_DELETE               = 52
	FlagAuditLogEventsEMOJI_CREATE                 = 60
	FlagAuditLogEventsEMOJI_UPDATE                 = 61
	FlagAuditLogEventsEMOJI_DELETE                 = 62
	FlagAuditLogEventsMESSAGE_DELETE               = 72
	FlagAuditLogEventsMESSAGE_BULK_DELETE          = 73
	FlagAuditLogEventsMESSAGE_PIN                  = 74
	FlagAuditLogEventsMESSAGE_UNPIN                = 75
	FlagAuditLogEventsINTEGRATION_CREATE           = 80
	FlagAuditLogEventsINTEGRATION_UPDATE           = 81
	FlagAuditLogEventsINTEGRATION_DELETE           = 82
	FlagAuditLogEventsSTAGE_INSTANCE_CREATE        = 83
	FlagAuditLogEventsSTAGE_INSTANCE_UPDATE        = 84
	FlagAuditLogEventsSTAGE_INSTANCE_DELETE        = 85
	FlagAuditLogEventsSTICKER_CREATE               = 90
	FlagAuditLogEventsSTICKER_UPDATE               = 91
	FlagAuditLogEventsSTICKER_DELETE               = 92
	FlagAuditLogEventsGUILD_SCHEDULED_EVENT_CREATE = 100
	FlagAuditLogEventsGUILD_SCHEDULED_EVENT_UPDATE = 101
	FlagAuditLogEventsGUILD_SCHEDULED_EVENT_DELETE = 102
	FlagAuditLogEventsTHREAD_CREATE                = 110
	FlagAuditLogEventsTHREAD_UPDATE                = 111
	FlagAuditLogEventsTHREAD_DELETE                = 112
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
	FlagAuditLogChangeKeyKeyAuditLogafk_channel_id                = "afk_channel_id"
	FlagAuditLogChangeKeyKeyAuditLogafk_timeout                   = "afk_timeout"
	FlagAuditLogChangeKeyKeyAuditLogallow                         = "allow"
	FlagAuditLogChangeKeyKeyAuditLogapplication_id                = "application_id"
	FlagAuditLogChangeKeyKeyAuditLogarchived                      = "archived"
	FlagAuditLogChangeKeyKeyAuditLogasset                         = "asset"
	FlagAuditLogChangeKeyKeyAuditLogauto_archive_duration         = "auto_archive_duration"
	FlagAuditLogChangeKeyKeyAuditLogavailable                     = "available"
	FlagAuditLogChangeKeyKeyAuditLogavatar_hash                   = "avatar_hash"
	FlagAuditLogChangeKeyKeyAuditLogbanner_hash                   = "banner_hash"
	FlagAuditLogChangeKeyKeyAuditLogbitrate                       = "bitrate"
	FlagAuditLogChangeKeyKeyAuditLogchannel_id                    = "channel_id"
	FlagAuditLogChangeKeyKeyAuditLogcode                          = "code"
	FlagAuditLogChangeKeyKeyAuditLogcolor                         = "color"
	FlagAuditLogChangeKeyKeyAuditLogcommunication_disabled_until  = "communication_disabled_until"
	FlagAuditLogChangeKeyKeyAuditLogdeaf                          = "deaf"
	FlagAuditLogChangeKeyKeyAuditLogdefault_auto_archive_duration = "default_auto_archive_duration"
	FlagAuditLogChangeKeyKeyAuditLogdefault_message_notifications = "default_message_notifications"
	FlagAuditLogChangeKeyKeyAuditLogdeny                          = "deny"
	FlagAuditLogChangeKeyKeyAuditLogdescription                   = "description"
	FlagAuditLogChangeKeyKeyAuditLogdiscovery_splash_hash         = "discovery_splash_hash"
	FlagAuditLogChangeKeyKeyAuditLogenable_emoticons              = "enable_emoticons"
	FlagAuditLogChangeKeyKeyAuditLogentity_type                   = "entity_type"
	FlagAuditLogChangeKeyKeyAuditLogexpire_behavior               = "expire_behavior"
	FlagAuditLogChangeKeyKeyAuditLogexpire_grace_period           = "expire_grace_period"
	FlagAuditLogChangeKeyKeyAuditLogexplicit_content_filter       = "explicit_content_filter"
	FlagAuditLogChangeKeyKeyAuditLogformat_type                   = "format_type"
	FlagAuditLogChangeKeyKeyAuditLogguild_id                      = "guild_id"
	FlagAuditLogChangeKeyKeyAuditLoghoist                         = "hoist"
	FlagAuditLogChangeKeyKeyAuditLogicon_hash                     = "icon_hash"
	FlagAuditLogChangeKeyKeyAuditLogimage_hash                    = "image_hash"
	FlagAuditLogChangeKeyKeyAuditLogid                            = "id"
	FlagAuditLogChangeKeyKeyAuditLoginvitable                     = "invitable"
	FlagAuditLogChangeKeyKeyAuditLoginviter_id                    = "inviter_id"
	FlagAuditLogChangeKeyKeyAuditLoglocation                      = "location"
	FlagAuditLogChangeKeyKeyAuditLoglocked                        = "locked"
	FlagAuditLogChangeKeyKeyAuditLogmax_age                       = "max_age"
	FlagAuditLogChangeKeyKeyAuditLogmax_uses                      = "max_uses"
	FlagAuditLogChangeKeyKeyAuditLogmentionable                   = "mentionable"
	FlagAuditLogChangeKeyKeyAuditLogmfa_level                     = "mfa_level"
	FlagAuditLogChangeKeyKeyAuditLogmute                          = "mute"
	FlagAuditLogChangeKeyKeyAuditLogname                          = "name"
	FlagAuditLogChangeKeyKeyAuditLognick                          = "nick"
	FlagAuditLogChangeKeyKeyAuditLognsfw                          = "nsfw"
	FlagAuditLogChangeKeyKeyAuditLogowner_id                      = "owner_id"
	FlagAuditLogChangeKeyKeyAuditLogpermission_overwrites         = "permission_overwrites"
	FlagAuditLogChangeKeyKeyAuditLogpermissions                   = "permissions"
	FlagAuditLogChangeKeyKeyAuditLogposition                      = "position"
	FlagAuditLogChangeKeyKeyAuditLogpreferred_locale              = "preferred_locale"
	FlagAuditLogChangeKeyKeyAuditLogprivacy_level                 = "privacy_level"
	FlagAuditLogChangeKeyKeyAuditLogprune_delete_days             = "prune_delete_days"
	FlagAuditLogChangeKeyKeyAuditLogpublic_updates_channel_id     = "public_updates_channel_id"
	FlagAuditLogChangeKeyKeyAuditLograte_limit_per_user           = "rate_limit_per_user"
	FlagAuditLogChangeKeyKeyAuditLogregion                        = "region"
	FlagAuditLogChangeKeyKeyAuditLogrules_channel_id              = "rules_channel_id"
	FlagAuditLogChangeKeyKeyAuditLogsplash_hash                   = "splash_hash"
	FlagAuditLogChangeKeyKeyAuditLogstatus                        = "status"
	FlagAuditLogChangeKeyKeyAuditLogsystem_channel_id             = "system_channel_id"
	FlagAuditLogChangeKeyKeyAuditLogtags                          = "tags"
	FlagAuditLogChangeKeyKeyAuditLogtemporary                     = "temporary"
	FlagAuditLogChangeKeyKeyAuditLogtopic                         = "topic"
	FlagAuditLogChangeKeyKeyAuditLogtype                          = "type"
	FlagAuditLogChangeKeyKeyAuditLogunicode_emoji                 = "unicode_emoji"
	FlagAuditLogChangeKeyKeyAuditLoguser_limit                    = "user_limit"
	FlagAuditLogChangeKeyKeyAuditLoguses                          = "uses"
	FlagAuditLogChangeKeyKeyAuditLogvanity_url_code               = "vanity_url_code"
	FlagAuditLogChangeKeyKeyAuditLogverification_level            = "verification_level"
	FlagAuditLogChangeKeyKeyAuditLogwidget_channel_id             = "widget_channel_id"
	FlagAuditLogChangeKeyKeyAuditLogwidget_enabled                = "widget_enabled"
	FlagAuditLogChangeKeyKeyAuditLogadd                           = "add"
	FlagAuditLogChangeKeyKeyAuditLogremove                        = "remove"
)
