package resources

// User Object
// https://discord.com/developers/docs/resources/user#user-object
type User struct {
	ID            int64  `json:"id,omitempty"`
	Username      string `json:"username,omitempty"`
	Discriminator string `json:"discriminator,omitempty"`
	Avatar        string `json:"avatar"`
	Bot           bool   `json:"bot,omitempty"`
	System        bool   `json:"system,omitempty"`
	MFAEnabled    bool   `json:"mfa_enabled,omitempty"`
	Banner        string `json:"banner"`
	AccentColor   int    `json:"accent_color"`
	Locale        string `json:"locale,omitempty"`
	Verified      bool   `json:"verified,omitempty"`
	Email         string `json:"email,omitempty"`
	Flags         int    `json:"flag,omitempty"`
	PremiumType   uint8  `json:"premium_type,omitempty"`
	PublicFlags   uint8  `json:"public_flag,omitempty"`
}

// User Flags
// https://discord.com/developers/docs/resources/user#user-object-user-flags
const (
	UserFlagsNONE                = 0
	STAFF                        = 1 << 0
	PARTNER                      = 1 << 1
	HYPESQUAD                    = 1 << 2
	BUG_HUNTER_LEVEL_1           = 1 << 3
	HYPESQUAD_ONLINE_HOUSE_ONE   = 1 << 6
	HYPESQUAD_ONLINE_HOUSE_TWO   = 1 << 7
	HYPESQUAD_ONLINE_HOUSE_THREE = 1 << 8
	PREMIUM_EARLY_SUPPORTER      = 1 << 9
	TEAM_PSEUDO_USER             = 1 << 10
	BUG_HUNTER_LEVEL_2           = 1 << 14
	VERIFIED_BOT                 = 1 << 16
	VERIFIED_DEVELOPER           = 1 << 17
	CERTIFIED_MODERATOR          = 1 << 18
	BOT_HTTP_INTERACTIONS        = 1 << 19
)

// Premium Types
// https://discord.com/developers/docs/resources/user#user-object-premium-types
const (
	PremiumTypesNONE         = 0
	PremiumTypesNITROCLASSIC = 1
	PremiumTypesNITRO        = 2
)

// User Connection Object
// https://discord.com/developers/docs/resources/user#connection-object-connection-structure
type Connection struct {
	ID           string         `json:"id"`
	Name         string         `json:"name"`
	Type         string         `json:"type"`
	Revoked      bool           `json:"revoked"`
	Integrations []*Integration `json:"integrations"`
	Verified     bool           `json:"verified"`
	FriendSync   bool           `json:"friend_sync"`
	ShowActivity bool           `json:"show_activity"`
	Visibility   int            `json:"visibility"`
}

// Visibility Types
// https://discord.com/developers/docs/resources/user#connection-object-visibility-types
const (
	VisibilityTypesNONE     = 0
	VisibilityTypesEVERYONE = 1
)
