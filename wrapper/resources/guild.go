package resource

import "time"

// Guild Object
// https://discord.com/developers/docs/resources/guild#guild-object
type Guild struct {
	ID                          int64                  `json:"id"`
	Name                        string                 `json:"name"`
	Icon                        string                 `json:"icon"`
	Splash                      string                 `json:"splash"`
	DiscoverySplash             string                 `json:"discovery_splash,omitempty"`
	Owner                       bool                   `json:"owner,omitempty"`
	OwnerID                     int64                  `json:"owner_id"`
	Permissions                 uint8                  `json:"permissions,omitempty"`
	Region                      string                 `json:"region"`
	AfkChannelID                int64                  `json:"afk_channel_id"`
	AfkTimeout                  uint                   `json:"afk_timeout"`
	WidgetEnabled               bool                   `json:"widget_enabled,omit_empty"`
	WidgetChannelID             int64                  `json:"widget_channel_id,omit_empty"`
	VerificationLevel           uint8                  `json:"verification_level"`
	DefaultMessageNotifications uint8                  `json:"default_message_notifications"`
	ExplicitContentFilter       uint8                  `json:"explicit_content_filter"`
	Roles                       []*Role                `json:"roles"`
	Emojis                      []*Emoji               `json:"emojis"`
	Features                    []string               `json:"features"`
	MFALevel                    uint8                  `json:"mfa_level"`
	ApplicationID               int64                  `json:"application_id"`
	SystemChannelID             int64                  `json:"system_channel_id,omitempty"`
	SystemChannelFlags          int                    `json:"system_channel_flags,omitempty"`
	RulesChannelID              int64                  `json:"rules_channel_id,omitempty"`
	JoinedAt                    *time.Time             `json:"joined_at,omitempty"`
	Large                       bool                   `json:"large,omitempty"`
	Unavailable                 bool                   `json:"unavailable"`
	MemberCount                 uint                   `json:"member_count,omitempty"`
	VoiceStates                 []*VoiceState          `json:"voice_states,omitempty"`
	Members                     []*GuildMember         `json:"members,omitempty"`
	Channels                    []*Channel             `json:"channels,omitempty"`
	Threads                     []*Channel             `json:"threads,omitempty"`
	Presences                   []*PresenceUpdate      `json:"presences,omitempty"`
	MaxPresences                int                    `json:"max_presences"`
	MaxMembers                  int                    `json:"max_members"`
	VanityUrl                   string                 `json:"vanity_url_code,omitempty"`
	Description                 string                 `json:"description,omitempty"`
	Banner                      string                 `json:"banner,omitempty"`
	PremiumTier                 uint8                  `json:"premium_tier"`
	PremiumSubscriptionCount    uint                   `json:"premium_subscription_count,omitempty"`
	PreferredLocale             string                 `json:"preferred_locale"`
	PublicUpdatesChannelID      int64                  `json:"public_updates_channel_id"`
	MaxVideoChannelUsers        int                    `json:"max_video_channel_users"`
	ApproximateMemberCount      int                    `json:"approximate_member_count"`
	ApproximatePresenceCount    int                    `json:"approximate_presence_count"`
	WelcomeScreen               *WelcomeScreen         `json:"welcome_screen"`
	NSFWLevel                   int                    `json:"nsfw_level"`
	StageInstances              []*StageInstance       `json:"stage_instances"`
	Stickers                    []*Sticker             `json:"stickers"`
	GuildScheduledEvents        []*GuildScheduledEvent `json:"guild_scheduled_events"`
	PremiumProgressBarEnabled   bool                   `json:"premium_progress_bar_enabled"`
}

// Default Message Notification Level
// https://discord.com/developers/docs/resources/guild#guild-object-default-message-notification-level
const (
	ALL_MESSAGES  = 0
	ONLY_MENTIONS = 1
)

// Explicit Content Filter Level
// https://discord.com/developers/docs/resources/guild#guild-object-explicit-content-filter-level
const (
	DISABLED              = 0
	MEMBERS_WITHOUT_ROLES = 1
	ALL_MEMBERS           = 2
)

// MFA Level
// https://discord.com/developers/docs/resources/guild#guild-object-mfa-level
const (
	MFALevelNONE     = 0
	MFALevelELEVATED = 1
)

// Verification Level
// https://discord.com/developers/docs/resources/guild#guild-object-verification-level
const (
	VerificationLevelNONE      = 0
	VerificationLevelLOW       = 1
	VerificationLevelMEDIUM    = 2
	VerificationLevelHIGH      = 3
	VerificationLevelVERY_HIGH = 4
)

// Guild NSFW Level
// https://discord.com/developers/docs/resources/guild#guild-object-guild-nsfw-level
const (
	GuildNSFWLevelDEFAULT        = 0
	GuildNSFWLevelEXPLICIT       = 1
	GuildNSFWLevelSAFE           = 2
	GuildNSFWLevelAGE_RESTRICTED = 3
)

// Premium Tier
// https://discord.com/developers/docs/resources/guild#guild-object-premium-tier
const (
	PremiumTierNONE  = 0
	PremiumTierONE   = 1
	PremiumTierTWO   = 2
	PremiumTierTHREE = 3
)

// System Channel Flags
// https://discord.com/developers/docs/resources/guild#guild-object-system-channel-flags
const (
	SUPPRESS_JOIN_NOTIFICATIONS           = 1 << 0
	SUPPRESS_PREMIUM_SUBSCRIPTIONS        = 1 << 1
	SUPPRESS_GUILD_REMINDER_NOTIFICATIONS = 1 << 2
	SUPPRESS_JOIN_NOTIFICATION_REPLIES    = 1 << 3
)

// Guild Features
// https://discord.com/developers/docs/resources/guild#guild-object-guild-features
// TODO

// Guild Preview Object
// https://discord.com/developers/docs/resources/guild#guild-preview-object-guild-preview-structure

type GuildPreview struct {
	ID                       string     `json:"id"`
	Name                     string     `json:"name"`
	Icon                     string     `json:"icon"`
	Splash                   string     `json:"splash"`
	DiscoverySplash          string     `json:"discovery_splash"`
	Emojis                   []*Emoji   `json:"emojis"`
	Features                 []string   `json:"features"`
	ApproximateMemberCount   int        `json:"approximate_member_count"`
	ApproximatePresenceCount int        `json:"approximate_presence_count"`
	Description              string     `json:"description"`
	Stickers                 []*Sticker `json:"stickers"`
}

// Guild Widget Settings Object
// https://discord.com/developers/docs/resources/guild#guild-widget-settings-object
type GuildWidget struct {
	Enabled   bool  `json:"enabled"`
	ChannelID int64 `json:"channel_id"`
}

// Guild Ban Object
// https://discord.com/developers/docs/resources/guild#ban-object
type Ban struct {
	Reason string `json:"reason"`
	User   *User  `json:"user"`
}

// Guild Scheduled Event Object
// https://discord.com/developers/docs/resources/guild-scheduled-event#guild-scheduled-event-object-guild-scheduled-event-structure
type GuildScheduledEvent struct {
	ID                 int64                             `json:"id"`
	GuildID            int64                             `json:"guild_id"`
	ChannelID          int64                             `json:"channel_id"`
	CreatorID          int64                             `json:"creator_id"`
	Name               string                            `json:"name"`
	Description        string                            `json:"description"`
	ScheduledStartTime time.Time                         `json:"scheduled_start_time"`
	ScheduledEndTime   time.Time                         `json:"scheduled_end_time"`
	PrivacyLevel       int                               `json:"privacy_level"`
	Status             int                               `json:"status"`
	EntityType         int                               `json:"entity_type"`
	EntityID           int64                             `json:"entity_id"`
	EntityMetadata     GuildScheduledEventEntityMetadata `json:"entity_metadata"`
	Creator            *User                             `json:"creator"`
	UserCount          int                               `json:"user_count"`
	Image              string                            `json:"image"`
}

// Guild Scheduled Event Privacy Level
// https://discord.com/developers/docs/resources/guild-scheduled-event#guild-scheduled-event-object-guild-scheduled-event-privacy-level
const (
	GuildScheduledEventPrivacyLevelGUILD_ONLY = 2
)

// Guild Scheduled Event Entity Types
// https://discord.com/developers/docs/resources/guild-scheduled-event#guild-scheduled-event-object-guild-scheduled-event-entity-types
const (
	STAGE_INSTANCE = 1
	VOICE          = 2
	EXTERNAL       = 3
)

// Guild Scheduled Event Status
// https://discord.com/developers/docs/resources/guild-scheduled-event#guild-scheduled-event-object-guild-scheduled-event-status
const (
	SCHEDULED = 1
	ACTIVE    = 2
	COMPLETED = 3
	CANCELED  = 4
)

// Guild Scheduled Event Entity Metadata
// https://discord.com/developers/docs/resources/guild-scheduled-event#guild-scheduled-event-object-guild-scheduled-event-entity-metadata
type GuildScheduledEventEntityMetadata struct {
	// location of the event (1-100 characters)
	// required for events with 'entity_type': EXTERNAL
	Location string `json:"location"`
}

// Guild Scheduled Event User Object
// https://discord.com/developers/docs/resources/guild-scheduled-event#guild-scheduled-event-user-object-guild-scheduled-event-user-structure
type GuildScheduledEventUser struct {
	GuildScheduledEventID string       `json:"guild_scheduled_event_id"`
	User                  *User        `json:"user"`
	Member                *GuildMember `json:"member"`
}

// Guild Template Object
// https://discord.com/developers/docs/resources/guild-template#guild-template-object
type GuildTemplate struct {
	Code                  string    `json:"code"`
	Name                  string    `json:"name"`
	Description           string    `json:"description"`
	UsageCount            string    `json:"usage_count"`
	CreatorID             string    `json:"creator_id"`
	Creator               *User     `json:"creator"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
	SourceGuildID         string    `json:"source_guild_id"`
	SerializedSourceGuild *Guild    `json:"serialized_source_guild"`
	IsDirty               bool      `json:"is_dirty"`
}

// Invite Object
// https://discord.com/developers/docs/resources/invite#invite-object
type Invite struct {
	Code                     string              `json:"code"`
	Guild                    *Guild              `json:"guild"`
	Channel                  *Channel            `json:"channel"`
	Inviter                  *User               `json:"inviter"`
	TargetUser               *User               `json:"target_user"`
	TargetType               uint8               `json:"target_type"`
	TargetApplication        *Application        `json:"target_application"`
	ApproximatePresenceCount int                 `json:"approximate_presence_count"`
	ApproximateMemberCount   int                 `json:"approximate_member_count"`
	ExpiresAt                time.Time           `json:"expires_at"`
	StageInstance            StageInstance       `json:"stage_instance"`
	GuildScheduledEvent      GuildScheduledEvent `json:"guild_scheduled_event"`
}

// Invite Target Types
// https://discord.com/developers/docs/resources/invite#invite-object-invite-target-types
const (
	InviteTargetTypesSTREAM = 1
	EMBEDDED_APPLICATION    = 2
)

// Invite Metadata Object
// https://discord.com/developers/docs/resources/invite#invite-metadata-object-invite-metadata-structure
type InviteMetadata struct {
	Uses      int       `json:"uses"`
	MaxUses   int       `json:"max_uses"`
	MaxAge    int       `json:"max_age"`
	Temporary bool      `json:"temporary"`
	CreatedAt time.Time `json:"created_at"`
}

// Invite Stage Instance Object
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
