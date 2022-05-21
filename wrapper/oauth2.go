package wrapper

import (
	"fmt"
	"strings"
	"time"
)

// OAuth2 Scopes
// https://discord.com/developers/docs/topics/oauth2#shared-resources-oauth2-scopes
const (
	FlagOAuth2ScopeActivitiesRead                        = "activities.read"
	FlagOAuth2ScopeActivitiesWrite                       = "activities.write"
	FlagOAuth2ScopeApplicationsBuildsRead                = "applications.builds.read"
	FlagOAuth2ScopeApplicationsBuildsUpload              = "applications.builds.upload"
	FlagOAuth2ScopeApplicationsCommands                  = "applications.commands"
	FlagOAuth2ScopeApplicationsCommandsUpdate            = "applications.commands.update"
	FlagOAuth2ScopeApplicationsCommandsPermissionsUpdate = "applications.commands.permissions.update"
	FlagOAuth2ScopeApplicationsEntitlements              = "applications.entitlements"
	FlagOAuth2ScopeApplicationsStoreUpdate               = "applications.store.update"
	FlagOAuth2ScopeBot                                   = "bot"
	FlagOAuth2ScopeConnections                           = "connections"
	FlagOAuth2ScopeDM_channelsRead                       = "dm_channels.read"
	FlagOAuth2ScopeEmail                                 = "email"
	FlagOAuth2ScopeGDMJoin                               = "gdm.join"
	FlagOAuth2ScopeGuilds                                = "guilds"
	FlagOAuth2ScopeGuildsJoin                            = "guilds.join"
	FlagOAuth2ScopeGuildsMembersRead                     = "guilds.members.read"
	FlagOAuth2ScopeIdentify                              = "identify"
	FlagOAuth2ScopeMessagesRead                          = "messages.read"
	FlagOAuth2ScopeRelationshipsRead                     = "relationships.read"
	FlagOAuth2ScopeRPC                                   = "rpc"
	FlagOAuth2ScopeRPCActivitiesWrite                    = "rpc.activities.write"
	FlagOAuth2ScopeRPCNotificationsRead                  = "rpc.notifications.read"
	FlagOAuth2ScopeRPCVoiceRead                          = "rpc.voice.read"
	FlagOAuth2ScopeRPCVoiceWrite                         = "rpc.voice.write"
	FlagOAuth2ScopeVoice                                 = "voice"
	FlagOAuth2ScopeWebhookIncoming                       = "webhook.incoming"
)

// OAuth2 URLs
// https://discord.com/developers/docs/topics/oauth2#shared-resources-oauth2-urls
// TODO: swap to dasgo endpoint function.
const (
	EndpointBaseAuthorizationURL   = "https://discord.com/api/oauth2/authorize"
	EndpointBaseTokenURL           = "https://discord.com/api/oauth2/token"
	EndpointBaseTokenRevocationURL = "https://discord.com/api/oauth2/token/revoke"
)

// ContentTypeURL represents an HTTP header that indicates a encoded URL.
var ContentTypeURL = []byte("application/x-www-form-urlencoded")

// AuthorizationCodeGrant performs an OAuth2 authorization code grant.
//https://discord.com/developers/docs/topics/oauth2#authorization-code-grant
func AuthorizationCodeGrant(bot *Client) string {
	// retrieve an access code.
	authorizationURL := AuthorizationURL(bot)
	// get request to authorizationURL for code.

	// exchange the access code for a user's access token.
	/*
		data := map[string]string{
			"client_id":     bot.Authorization.ClientID,
			"client_secret": bot.Authorization.ClientSecret,
			"grant_type":    "authorization_code",
			"code":          code,
			"redirect_uri":  bot.Authorization.RedirectURI,
		}
		header := ContentTypeURL
	*/

	// post request to access token url
	// marshal response to AccessTokenResponse
	// set bot AccessToken information.

	return authorizationURL
}

// AuthorizationURL generates an authorization URL from a given client
func AuthorizationURL(bot *Client) string {
	params := make([]string, 0, 5)

	// client_id is the application client id.
	params = append(params, "client_id="+bot.Authorization.ClientID)

	// scope is a list of OAuth2 scopes separated by url encoded spaces (%20).
	if len(bot.Authorization.Scopes) > 0 {
		var scope strings.Builder
		scope.WriteString("scope=")

		for i, s := range bot.Authorization.Scopes {
			if i > 0 {
				scope.WriteString("%20")
			}

			scope.WriteString(s)
		}
		params = append(params, scope.String())
	}

	// redirect_uri is the URL registered while creating the application.
	if bot.Authorization.RedirectURI != "" {
		params = append(params, "redirect_uri="+bot.Authorization.RedirectURI)
	}

	// state is the unique string mentioned in State and Security.
	if bot.Authorization.State != "" {
		params = append(params, "state="+bot.Authorization.State)
	}

	// prompt controls how the authorization flow handles existing authorizations.
	if bot.Authorization.Prompt != "" {
		params = append(params, "prompt="+bot.Authorization.Prompt)
	}

	return EndpointBaseAuthorizationURL + "?response_type=code&" + strings.Join(params, "&")
}

// Access Token Response
// https://discord.com/developers/docs/topics/oauth2#authorization-code-grant-access-token-response
type AccessTokenResponse struct {
	AccessToken  string        `json:"access_token,omitempty"`
	TokenType    string        `json:"token_type,omitempty"`
	ExpiresIn    time.Duration `json:"expires_in,omitempty"`
	RefreshToken string        `json:"refresh_token,omitempty"`
	Scope        string        `json:"scope,omitempty"`
}

// RefreshAccessTokenExchange refreshes the given client's access token.
func RefreshAccessTokenExchange(bot *Client) error {
	if bot.AccessToken == nil {
		return fmt.Errorf("cannot refresh access token without access token information.")
	}

	/*
		// post request to access token url
		data = {
			"client_id": bot.Authorization.ClientID,
			"client_secret": bot.Authorization.ClientSecret,
			"grant_type": "refresh_token",
			"refresh_token": bot.AccessToken.RefreshToken
		  }
		  headers = {
			'Content-Type': 'application/x-www-form-urlencoded'
		  }
	*/

	return nil
}
