package resource

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
	ActionType *AuditLogAction   `json:"action_type"`
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

// AuditLogAction is the Action of the AuditLog (see AuditLogAction* consts)
// https://discord.com/developers/docs/resources/audit-log#audit-log-entry-object-audit-log-events
type AuditLogAction int

// Audit Log Change Object
// https://discord.com/developers/docs/resources/audit-log#audit-log-change-object
type AuditLogChange struct {
	NewValue interface{} `json:"new_value"`
	OldValue interface{} `json:"old_value"`
	Key      string      `json:"key"`
}

// TODO: Audit Log Change Key
// https://discord.com/developers/docs/resources/audit-log#audit-log-change-object-audit-log-change-key
type AuditLogChangeKey string
