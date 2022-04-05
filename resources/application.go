package resources

// Application Object
// https://discord.com/developers/docs/resources/application
type Application struct {
	ID                  Snowflake      `json:"id,omitempty"`
	Name                string         `json:"name,omitempty"`
	Icon                string         `json:"icon,omitempty"`
	Description         string         `json:"description,omitempty"`
	RPCOrigins          []string       `json:"rpc_origins,omitempty"`
	BotPublic           bool           `json:"bot_public,omitempty"`
	BotRequireCodeGrant bool           `json:"bot_require_code_grant,omitempty"`
	TermsOfServiceURL   string         `json:"terms_of_service_url,omitempty"`
	PrivacyProxyURL     string         `json:"privacy_policy_url,omitempty"`
	Owner               *User          `json:"owner,omitempty"`
	VerifyKey           string         `json:"verify_key,omitempty"`
	Team                *Team          `json:"team,omitempty"`
	GuildID             Snowflake      `json:"guild_id,omitempty"`
	PrimarySKUID        Snowflake      `json:"primary_sku_id,omitempty"`
	Slug                string         `json:"slug,omitempty"`
	CoverImage          string         `json:"cover_image,omitempty"`
	Flags               Flag           `json:"flags,omitempty"`
	Summary             string         `json:"summary,omitempty"`
	InstallParams       *InstallParams `json:"install_params,omitempty"`
	CustomInstallURL    string         `json:"custom_install_url,omitempty"`
}

// Application Flags
// https://discord.com/developers/docs/resources/application#application-object-application-flags
const (
	FlagApplicationFlagsGATEWAY_PRESENCE                 = 1 << 12
	FlagApplicationFlagsGATEWAY_PRESENCE_LIMITED         = 1 << 13
	FlagApplicationFlagsGATEWAY_GUILD_MEMBERS            = 1 << 14
	FlagApplicationFlagsGATEWAY_GUILD_MEMBERS_LIMITED    = 1 << 15
	FlagApplicationFlagsVERIFICATION_PENDING_GUILD_LIMIT = 1 << 16
	FlagApplicationFlagsEMBEDDED                         = 1 << 17
	FlagApplicationFlagsGATEWAY_MESSAGE_CONTENT          = 1 << 18
	FlagApplicationFlagsGATEWAY_MESSAGE_CONTENT_LIMITED  = 1 << 19
)

// Install Params Object
// https://discord.com/developers/docs/resources/application#install-params-object
type InstallParams struct {
	Scopes      []string `json:"scopes,omitempty"`
	Permissions string   `json:"permissions,omitempty"`
}
