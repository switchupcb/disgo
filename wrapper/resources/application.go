package resources

// Application Object
// https://discord.com/developers/docs/resources/application
type Application struct {
	ID                  string         `json:"id,omitempty"`
	Name                string         `json:"name"`
	Icon                string         `json:"icon,omitempty"`
	Description         string         `json:"description,omitempty"`
	RPCOrigins          []string       `json:"rpc_origins,omitempty"`
	BotPublic           bool           `json:"bot_public,omitempty"`
	BotRequireCodeGrant bool           `json:"bot_require_code_grant,omitempty"`
	TermsOfServiceURL   string         `json:"terms_of_service_url"`
	PrivacyProxyURL     string         `json:"privacy_policy_url"`
	Owner               *User          `json:"owner"`
	VerifyKey           string         `json:"verify_key"`
	Team                *Team          `json:"team"`
	GuildID             string         `json:"guild_id"`
	PrimarySKUID        string         `json:"primary_sku_id"`
	Slug                string         `json:"slug"`
	CoverImage          string         `json:"cover_image"`
	Flags               int            `json:"flags,omitempty"`
	Summary             string         `json:"summary"`
	InstallParams       *InstallParams `json:"install_params"`
	CustomInstallURL    string         `json:"custom_install_url"`
}

// Application Flags
// https://discord.com/developers/docs/resources/application#application-object-application-flags
const (
	GATEWAY_PRESENCE                 = 1 << 12
	GATEWAY_PRESENCE_LIMITED         = 1 << 13
	GATEWAY_GUILD_MEMBERS            = 1 << 14
	GATEWAY_GUILD_MEMBERS_LIMITED    = 1 << 15
	VERIFICATION_PENDING_GUILD_LIMIT = 1 << 16
	EMBEDDED                         = 1 << 17
	GATEWAY_MESSAGE_CONTENT          = 1 << 18
	GATEWAY_MESSAGE_CONTENT_LIMITED  = 1 << 19
)

// Install Params Object
// https://discord.com/developers/docs/resources/application#install-params-object
type InstallParams struct {
	Scopes      []string `json:"scopes"`
	Permissions string   `json:"permissions"`
}
