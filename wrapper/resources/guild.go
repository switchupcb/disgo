package resources

import "time"

// Guild Object
// https://discord.com/developers/docs/resources/guild#guild-object
type Guild struct {
	ID                          Snowflake              `json:"id,omitempty"`
	Name                        string                 `json:"name,omitempty"`
	Icon                        string                 `json:"icon,omitempty"`
	Splash                      string                 `json:"splash,omitempty"`
	DiscoverySplash             string                 `json:"discovery_splash,omitempty"`
	Owner                       bool                   `json:"owner,omitempty"`
	OwnerID                     Snowflake              `json:"owner_id,omitempty"`
	Permissions                 *string                `json:"permissions,omitempty"`
	Region                      string                 `json:"region,omitempty"`
	AfkChannelID                Snowflake              `json:"afk_channel_id,omitempty"`
	AfkTimeout                  *uint                  `json:"afk_timeout,omitempty"`
	WidgetEnabled               bool                   `json:"widget_enabled,omitempty"`
	WidgetChannelID             Snowflake              `json:"widget_channel_id,omitempty"`
	VerificationLevel           *Flag                  `json:"verification_level,omitempty"`
	DefaultMessageNotifications *Flag                  `json:"default_message_notifications,omitempty"`
	ExplicitContentFilter       *Flag                  `json:"explicit_content_filter,omitempty"`
	Roles                       []*Role                `json:"roles,omitempty"`
	Emojis                      []*Emoji               `json:"emojis,omitempty"`
	Features                    []*string              `json:"features,omitempty"`
	MFALevel                    *Flag                  `json:"mfa_level,omitempty"`
	ApplicationID               Snowflake              `json:"application_id,omitempty"`
	SystemChannelID             Snowflake              `json:"system_channel_id,omitempty"`
	SystemChannelFlags          BitFlag                `json:"system_channel_flags,omitempty"`
	RulesChannelID              Snowflake              `json:"rules_channel_id,omitempty"`
	JoinedAt                    time.Time              `json:"joined_at,omitempty"`
	Large                       bool                   `json:"large,omitempty"`
	Unavailable                 bool                   `json:"unavailable,omitempty"`
	MemberCount                 uint                   `json:"member_count,omitempty"`
	VoiceStates                 []*VoiceState          `json:"voice_states,omitempty"`
	Members                     []*GuildMember         `json:"members,omitempty"`
	Channels                    []*Channel             `json:"channels,omitempty"`
	Threads                     []*Channel             `json:"threads,omitempty"`
	Presences                   []*PresenceUpdate      `json:"presences,omitempty"`
	MaxPresences                CodeFlag               `json:"max_presences,omitempty"`
	MaxMembers                  int                    `json:"max_members,omitempty"`
	VanityUrl                   *string                `json:"vanity_url_code,omitempty"`
	Description                 *string                `json:"description,omitempty"`
	Banner                      string                 `json:"banner,omitempty"`
	PremiumTier                 *Flag                  `json:"premium_tier,omitempty"`
	PremiumSubscriptionCount    *CodeFlag              `json:"premium_subscription_count,omitempty"`
	PreferredLocale             string                 `json:"preferred_locale,omitempty"`
	PublicUpdatesChannelID      Snowflake              `json:"public_updates_channel_id,omitempty"`
	MaxVideoChannelUsers        int                    `json:"max_video_channel_users,omitempty"`
	ApproximateMemberCount      int                    `json:"approximate_member_count,omitempty"`
	ApproximatePresenceCount    int                    `json:"approximate_presence_count,omitempty"`
	WelcomeScreen               *WelcomeScreen         `json:"welcome_screen,omitempty"`
	NSFWLevel                   *Flag                  `json:"nsfw_level,omitempty"`
	StageInstances              []*StageInstance       `json:"stage_instances,omitempty"`
	Stickers                    []*Sticker             `json:"stickers,omitempty"`
	GuildScheduledEvents        []*GuildScheduledEvent `json:"guild_scheduled_events,omitempty"`
	PremiumProgressBarEnabled   bool                   `json:"premium_progress_bar_enabled,omitempty"`
}

// Default Message Notification Level
// https://discord.com/developers/docs/resources/guild#guild-object-default-message-notification-level
const (
	FlagLevelNotificationMessageDefaultALL_MESSAGES  = 0
	FlagLevelNotificationMessageDefaultONLY_MENTIONS = 1
)

// Explicit Content Filter Level
// https://discord.com/developers/docs/resources/guild#guild-object-explicit-content-filter-level
const (
	FlagLevelFilterContentExplicitDISABLED              = 0
	FlagLevelFilterContentExplicitMEMBERS_WITHOUT_ROLES = 1
	FlagLevelFilterContentExplicitALL_MEMBERS           = 2
)

// MFA Level
// https://discord.com/developers/docs/resources/guild#guild-object-mfa-level
const (
	FlagLevelMFANONE     = 0
	FlagLevelMFAELEVATED = 1
)

// Verification Level
// https://discord.com/developers/docs/resources/guild#guild-object-verification-level
const (
	FlagLevelVerificationNONE      = 0
	FlagLevelVerificationLOW       = 1
	FlagLevelVerificationMEDIUM    = 2
	FlagLevelVerificationHIGH      = 3
	FlagLevelVerificationVERY_HIGH = 4
)

// Guild NSFW Level
// https://discord.com/developers/docs/resources/guild#guild-object-guild-nsfw-level
const (
	FlagLevelNSFWGuildDEFAULT        = 0
	FlagLevelNSFWGuildEXPLICIT       = 1
	FlagLevelNSFWGuildSAFE           = 2
	FlagLevelNSFWGuildAGE_RESTRICTED = 3
)

// Premium Tier
// https://discord.com/developers/docs/resources/guild#guild-object-premium-tier
const (
	FlagTierPremiumNONE  = 0
	FlagTierPremiumONE   = 1
	FlagTierPremiumTWO   = 2
	FlagTierPremiumTHREE = 3
)

// System Channel Flags
// https://discord.com/developers/docs/resources/guild#guild-object-system-channel-flags
const (
	FlagFlagsChannelSystemSUPPRESS_JOIN_NOTIFICATIONS           = 1 << 0
	FlagFlagsChannelSystemSUPPRESS_PREMIUM_SUBSCRIPTIONS        = 1 << 1
	FlagFlagsChannelSystemSUPPRESS_GUILD_REMINDER_NOTIFICATIONS = 1 << 2
	FlagFlagsChannelSystemSUPPRESS_JOIN_NOTIFICATION_REPLIES    = 1 << 3
)

// Guild Features
// https://discord.com/developers/docs/resources/guild#guild-object-guild-features
var (
	GuildFeatures = map[string]string{
		"ANIMATED_BANNER":                  "guild has access to set an animated guild banner image",
		"ANIMATED_ICON":                    "guild has access to set an animated guild icon",
		"BANNER":                           "guild has access to set a guild banner image",
		"COMMERCE":                         "guild has access to use commerce features (i.e. create store channels)",
		"COMMUNITY":                        "guild can enable welcome screen, Membership Screening, stage channels and discovery, and receives community updates",
		"DISCOVERABLE":                     "guild is able to be discovered in the directory",
		"FEATURABLE":                       "guild is able to be featured in the directory",
		"INVITE_SPLASH":                    "guild has access to set an invite splash background",
		"MEMBER_VERIFICATION_GATE_ENABLED": "guild has enabled Membership Screening",
		"MONETIZATION_ENABLED":             "guild has enabled monetization",
		"MORE_STICKERS":                    "guild has increased custom sticker slots",
		"NEWS":                             "guild has access to create news channels",
		"PARTNERED":                        "guild is partnered",
		"PREVIEW_ENABLED":                  "guild can be previewed before joining via Membership Screening or the directory",
		"PRIVATE_THREADS":                  "guild has access to create private threads",
		"ROLE_ICONS":                       "guild is able to set role icons",
		"SEVEN_DAY_THREAD_ARCHIVE":         "guild has access to the seven day archive time for threads",
		"THREE_DAY_THREAD_ARCHIVE":         "guild has access to the three day archive time for threads",
		"TICKETED_EVENTS_ENABLED":          "guild has enabled ticketed events",
		"VANITY_URL":                       "guild has access to set a vanity URL",
		"VERIFIED":                         "guild is verified",
		"VIP_REGIONS":                      "guild has access to set 384kbps bitrate in voice (previously VIP voice servers)",
		"WELCOME_SCREEN_ENABLED":           "guild has enabled the welcome screen",
	}
)

// Guild Preview Object
// https://discord.com/developers/docs/resources/guild#guild-preview-object-guild-preview-structure
type GuildPreview struct {
	ID                       string     `json:"id,omitempty"`
	Name                     string     `json:"name,omitempty"`
	Icon                     string     `json:"icon,omitempty"`
	Splash                   string     `json:"splash,omitempty"`
	DiscoverySplash          string     `json:"discovery_splash,omitempty"`
	Emojis                   []*Emoji   `json:"emojis,omitempty"`
	Features                 []*string  `json:"features,omitempty"`
	ApproximateMemberCount   int        `json:"approximate_member_count,omitempty"`
	ApproximatePresenceCount int        `json:"approximate_presence_count,omitempty"`
	Description              *string    `json:"description,omitempty"`
	Stickers                 []*Sticker `json:"stickers,omitempty"`
}

// Guild Widget Settings Object
// https://discord.com/developers/docs/resources/guild#guild-widget-settings-object
type GuildWidget struct {
	Enabled   bool      `json:"enabled,omitempty"`
	ChannelID Snowflake `json:"channel_id,omitempty"`
}

// Guild Member Object
// https://discord.com/developers/docs/resources/guild#guild-member-object
type GuildMember struct {
	User                       *User        `json:"user,omitempty"`
	Nick                       *string      `json:"nick,omitempty"`
	Avatar                     string       `json:"avatar,omitempty"`
	Roles                      []*Snowflake `json:"roles,omitempty"`
	GuildID                    Snowflake    `json:"guild_id,omitempty"`
	JoinedAt                   time.Time    `json:"joined_at,omitempty"`
	PremiumSince               time.Time    `json:"premium_since,omitempty"`
	Deaf                       bool         `json:"deaf,omitempty"`
	Mute                       bool         `json:"mute,omitempty"`
	Pending                    bool         `json:"pending,omitempty"`
	CommunicationDisabledUntil time.Time    `json:"communication_disabled_until,omitempty"`
	Permissions                *string      `json:"permissions,omitempty"`
}

// Guild Ban Object
// https://discord.com/developers/docs/resources/guild#ban-object
type Ban struct {
	Reason *string `json:"reason,omitempty"`
	User   *User   `json:"user,omitempty"`
}

// Guild Scheduled Event Object
// https://discord.com/developers/docs/resources/guild-scheduled-event#guild-scheduled-event-object-guild-scheduled-event-structure
type GuildScheduledEvent struct {
	ID                 Snowflake                         `json:"id,omitempty"`
	GuildID            Snowflake                         `json:"guild_id,omitempty"`
	ChannelID          Snowflake                         `json:"channel_id,omitempty"`
	CreatorID          Snowflake                         `json:"creator_id,omitempty"`
	Name               string                            `json:"name,omitempty"`
	Description        string                            `json:"description,omitempty"`
	ScheduledStartTime time.Time                         `json:"scheduled_start_time,omitempty"`
	ScheduledEndTime   time.Time                         `json:"scheduled_end_time,omitempty"`
	PrivacyLevel       Flag                              `json:"privacy_level,omitempty"`
	Status             Flag                              `json:"status,omitempty"`
	EntityType         Flag                              `json:"entity_type,omitempty"`
	EntityID           Snowflake                         `json:"entity_id,omitempty"`
	EntityMetadata     GuildScheduledEventEntityMetadata `json:"entity_metadata,omitempty"`
	Creator            *User                             `json:"creator,omitempty"`
	UserCount          CodeFlag                          `json:"user_count,omitempty"`
	Image              string                            `json:"image,omitempty"`
}

// Guild Scheduled Event Privacy Level
// https://discord.com/developers/docs/resources/guild-scheduled-event#guild-scheduled-event-object-guild-scheduled-event-privacy-level
const (
	FlagGuildScheduledEventPrivacyLevelGUILD_ONLY = 2
)

// Guild Scheduled Event Entity Types
// https://discord.com/developers/docs/resources/guild-scheduled-event#guild-scheduled-event-object-guild-scheduled-event-entity-types
const (
	FlagTypesEntityEventScheduledGuildSTAGE_INSTANCE = 1
	FlagTypesEntityEventScheduledGuildVOICE          = 2
	FlagTypesEntityEventScheduledGuildEXTERNAL       = 3
)

// Guild Scheduled Event Status
// https://discord.com/developers/docs/resources/guild-scheduled-event#guild-scheduled-event-object-guild-scheduled-event-status
const (
	FlagStatusEventScheduledGuildSCHEDULED = 1
	FlagStatusEventScheduledGuildACTIVE    = 2
	FlagStatusEventScheduledGuildCOMPLETED = 3
	FlagStatusEventScheduledGuildCANCELED  = 4
)

// Guild Scheduled Event Entity Metadata
// https://discord.com/developers/docs/resources/guild-scheduled-event#guild-scheduled-event-object-guild-scheduled-event-entity-metadata
type GuildScheduledEventEntityMetadata struct {
	// location of the event (1-100 characters)
	// required for events with 'entity_type': EXTERNAL
	Location string `json:"location,omitempty"`
}

// Guild Scheduled Event User Object
// https://discord.com/developers/docs/resources/guild-scheduled-event#guild-scheduled-event-user-object-guild-scheduled-event-user-structure
type GuildScheduledEventUser struct {
	GuildScheduledEventID Snowflake    `json:"guild_scheduled_event_id,omitempty"`
	User                  *User        `json:"user,omitempty"`
	Member                *GuildMember `json:"member,omitempty"`
}

// Guild Template Object
// https://discord.com/developers/docs/resources/guild-template#guild-template-object
type GuildTemplate struct {
	Code                  string    `json:"code,omitempty"`
	Name                  string    `json:"name,omitempty"`
	Description           *string   `json:"description,omitempty"`
	UsageCount            *int      `json:"usage_count,omitempty"`
	CreatorID             Snowflake `json:"creator_id,omitempty"`
	Creator               *User     `json:"creator,omitempty"`
	CreatedAt             time.Time `json:"created_at,omitempty"`
	UpdatedAt             time.Time `json:"updated_at,omitempty"`
	SourceGuildID         Snowflake `json:"source_guild_id,omitempty"`
	SerializedSourceGuild *Guild    `json:"serialized_source_guild,omitempty"`
	IsDirty               bool      `json:"is_dirty,omitempty"`
}

// Welcome Screen Object
// https://discord.com/developers/docs/resources/guild#welcome-screen-object-welcome-screen-structure
type WelcomeScreen struct {
	Description           *string                 `json:"description,omitempty"`
	WelcomeScreenChannels []*WelcomeScreenChannel `json:"welcome_channels,omitempty"`
}

// Welcome Screen Channel Structure
// https://discord.com/developers/docs/resources/guild#welcome-screen-object-welcome-screen-channel-structure
type WelcomeScreenChannel struct {
	ChannelID   Snowflake  `json:"channel_id,omitempty"`
	Description *string    `json:"description,omitempty"`
	EmojiID     *Snowflake `json:"emoji_id,omitempty"`
	EmojiName   *string    `json:"emoji_name,omitempty"`
}
