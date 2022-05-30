package wrapper

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"github.com/valyala/fasthttp"
)

const (
	grantTypeAuthorizationCodeGrant = "authorization_code"
	grantTypeRefreshToken           = "refresh_token"
	grantTypeClientCredentials      = "client_credentials"
)

// GenerateAuthorizationURL generates an authorization URL from a given client and response type.
func GenerateAuthorizationURL(bot *Client, response string) string {
	params := make([]string, 0, 6)

	// response_type is the type of response the redirect will return.
	if response != "" {
		params = append(params, "responsetype="+response)
	}

	// client_id is the application client id.
	params = append(params, "client_id="+bot.Authorization.ClientID)

	// scope is a list of OAuth2 scopes separated by url encoded spaces (%20).
	scope := urlQueryStringScope(bot.Authorization.Scopes)
	if scope != "" {
		params = append(params, scope)
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

	return EndpointAuthorizationURL() + "?" + strings.Join(params, "&")
}

// BotAuthParams represents parameters used to generate a bot authorization URL.
type BotAuthParams struct {
	// Bot provides the client_id and scopes parameters.
	Bot *Client

	// Permissions represents the permissions the bot is requesting.
	Permissions BitFlag

	// GuildID pre-selects a guild in the authorization prompt.
	GuildID string

	// DisableGuildSelect disables the ability to select other guilds
	// in the authorization prompt (when GuildID is provided).
	DisableGuildSelect bool

	// ResponseType provides the type of response the OAuth2 flow will return.
	//
	// In the context of bot authorization, response_type is only provided when
	// a scope outside of `bot` and `applications.commands` is requested.
	ResponseType string
}

// GenerateBotAuthorizationURL generates a bot authorization URL using the given BotAuthParams.
//
// Bot.Scopes must include "bot" to enable the OAuth2 Bot Flow.
func GenerateBotAuthorizationURL(p BotAuthParams) string {
	params := make([]string, 0, 3)

	// permissions is permissions the bot is requesting.
	params = append(params, "permissions="+strconv.FormatUint(uint64(p.Permissions), 10))

	// guild_id is the Guild ID of the guild that is pre-selected in the authorization prompt.
	if p.GuildID != "" {
		params = append(params, "guild_id="+p.GuildID)
	}

	// disable_guild_select determines whether the user will be allowed to select a guild
	// other than the guild_id.
	params = append(params, "disable_guild_select="+strconv.FormatBool(p.DisableGuildSelect))

	return GenerateAuthorizationURL(p.Bot, p.ResponseType) + strings.Join(params, "&")
}

// AuthorizationCodeGrant performs an OAuth2 authorization code grant.
//
// Send the user a valid Authorization URL, which can be generated using
// GenerateAuthorizationURL(bot, "code").
//
// When the user visits the Authorization URL, they will be prompted for authorization.
// If the user accepts the prompt, they will be redirected to the `redirect_uri`.
// This issues a GET request to the `redirect_uri` web server which YOU MUST HANDLE
// by parsing the request's URL Query String into a disgo.RedirectURL object.
//
// Retrieve the user's access token by calling THIS FUNCTION (with the disgo.RedirectURL parameter),
// which performs an Access Token Exchange.
//
// Refresh the token by using RefreshAuthorizationCodeGrant(bot, token).
//
// For more information, read https://discord.com/developers/docs/topics/oauth2#authorization-code-grant
func AuthorizationCodeGrant(bot *Client, ru *RedirectURL) (*AccessTokenResponse, error) {
	exchange := &AccessTokenExchange{
		ClientID:     bot.Authorization.ClientID,
		ClientSecret: bot.Authorization.ClientSecret,
		GrantType:    grantTypeAuthorizationCodeGrant,
		Code:         ru.Code,
		RedirectURI:  bot.Authorization.RedirectURI,
	}

	return exchange.Send(bot)
}

// RefreshAuthorizationCodeGrant refreshes an Access Token from an OAuth2 authorization code grant.
func RefreshAuthorizationCodeGrant(bot *Client, token *AccessTokenResponse) (*AccessTokenResponse, error) {
	exchange := &RefreshTokenExchange{
		ClientID:     bot.Authorization.ClientID,
		ClientSecret: bot.Authorization.ClientSecret,
		GrantType:    grantTypeRefreshToken,
		RefreshToken: token.RefreshToken,
	}

	return exchange.Send(bot)
}

// ImplicitGrant converts a RedirectURI (from a simplified OAuth2 grant) to an AccessTokenResponse.
//
// Send the user a valid Authorization URL, which can be generated using
// GenerateAuthorizationURL(bot, "token").
//
// When the user visits the Authorization URL, they will be prompted for authorization.
// If the user accepts the prompt, they will be redirected to the `redirect_uri`.
// This issues a GET request to the `redirect_uri` web server which YOU MUST HANDLE
// by parsing the request's URI Fragments into a disgo.RedirectURI object.
//
// A disgo.RedirectURI object is equivalent to a disgo.AccessTokenResponse,
// but it does NOT contain a refresh token.
//
// For more information, read https://discord.com/developers/docs/topics/oauth2#implicit-grant
func ImplicitGrant(ru *RedirectURI) *AccessTokenResponse {
	return &AccessTokenResponse{
		AccessToken:  ru.AccessToken,
		TokenType:    ru.TokenType,
		ExpiresIn:    ru.ExpiresIn,
		RefreshToken: "",
		Scope:        ru.Scope,
	}
}

// ClientCredentialsGrant performs a client credential OAuth2 grant for TESTING PURPOSES.
//
// The bot client's Authentication Header will be set to a Basic Authentication Header that
// uses the bot's ClientID as a username and ClientSecret as a password.
//
// A request will be made for a Client Credential grant which returns a disgo.AccessTokenResponse
// that does NOT contain a refresh token.
//
// For more information, read https://discord.com/developers/docs/topics/oauth2#client-credentials-grant
func ClientCredentialsGrant(bot *Client) (*AccessTokenResponse, error) {
	bot.Authentication.Header = "Basic " +
		base64.StdEncoding.EncodeToString([]byte(bot.Authorization.ClientID+":"+bot.Authorization.ClientSecret))

	grant := &ClientCredentialsTokenRequest{
		GrantType: grantTypeClientCredentials,
		Scope:     urlQueryStringScope(bot.Authorization.Scopes),
	}

	return grant.Send(bot)
}

// BotAuthorization performs a specialized OAuth2 flow for users to add bots to guilds.
//
// Send the user a valid Bot Authorization URL, which can be generated using
// GenerateBotAuthorizationURL(disgo.BotAuthParams{...}).
//
// When the user visits the Bot Authorization URL, they will be prompted for authorization.
// If the user accepts the prompt (with a guild), the bot will be added to the selected guild.
//
// For more information read, https://discord.com/developers/docs/topics/oauth2#bot-authorization-flow
func BotAuthorization() {}

// AdvancedBotAuthorization performs a specialized OAuth2 flow for users to add bots to guilds.
//
// Send the user a valid Bot Authorization URL, which can be generated using
// GenerateBotAuthorizationURL(disgo.BotAuthParams{...}).
//
// If the user accepts the prompt (with a guild), they will be redirected to the `redirect_uri`.
// This issues a GET request to the `redirect_uri` web server which YOU MUST HANDLE
// by parsing the request's URL Query String into a disgo.RedirectURL object.
//
// Retrieve the user's access token by calling THIS FUNCTION (with the disgo.RedirectURL parameter),
// which performs an Access Token Exchange.
//
// Refresh the token by using RefreshAuthorizationCodeGrant(bot, token).
//
// For more information read, https://discord.com/developers/docs/topics/oauth2#advanced-bot-authorization
func AdvancedBotAuthorization(bot *Client, ru *RedirectURL) (*AccessTokenResponse, error) {
	return AuthorizationCodeGrant(bot, ru)
}

// WebhookAuthorization performs a specialized OAuth2 authorization code grant.
//
// Send the user a valid Authorization URL, which can be generated using
// GenerateAuthorizationURL(bot, "code") when bot.Scopes is set to `webhook.incoming`.
//
// When the user visits the Authorization URL, they will be prompted for authorization.
// If the user accepts the prompt (with a channel), they will be redirected to the `redirect_uri`.
// This issues a GET request to the `redirect_uri` web server which YOU MUST HANDLE
// by parsing the request's URL Query String into a disgo.RedirectURL object.
//
// Retrieve the user's access token by calling THIS FUNCTION (with the disgo.RedirectURL parameter),
// which performs an Access Token Exchange.
//
// Refresh the token by using RefreshAuthorizationCodeGrant(bot, token).
//
// For more information read, https://discord.com/developers/docs/topics/oauth2#webhooks
func WebhookAuthorization(bot *Client, ru *RedirectURL) (*AccessTokenResponse, *Webhook, error) {
	exchange := &AccessTokenExchange{
		ClientID:     bot.Authorization.ClientID,
		ClientSecret: bot.Authorization.ClientSecret,
		GrantType:    grantTypeAuthorizationCodeGrant,
		Code:         ru.Code,
		RedirectURI:  bot.Authorization.RedirectURI,
	}

	query, err := EndpointQueryString(exchange)
	if err != nil {
		return nil, nil, fmt.Errorf(ErrQueryString, "WebhookAuthorization", err)
	}

	var result *WebhookTokenResponse
	err = SendRequest(bot, fasthttp.MethodPost, EndpointTokenURL()+"?"+query, contentTypeURL, nil, result)
	if err != nil {
		return nil, nil, fmt.Errorf(ErrSendRequest, "WebhookAuthorization", err)
	}

	// convert the webhook token response to an access token response (and webhook).
	token := &AccessTokenResponse{
		AccessToken:  result.AccessToken,
		TokenType:    result.TokenType,
		ExpiresIn:    result.ExpiresIn,
		RefreshToken: result.RefreshToken,
		Scope:        result.Scope,
	}

	return token, result.Webhook, nil
}

// Send sends an AccessTokenExchange request to Discord and returns an AccessTokenResponse.
func (r *AccessTokenExchange) Send(bot *Client) (*AccessTokenResponse, error) {
	query, err := EndpointQueryString(r)
	if err != nil {
		return nil, fmt.Errorf(ErrQueryString, "AccessTokenExchange", err)
	}

	var result *AccessTokenResponse
	err = SendRequest(bot, fasthttp.MethodPost, EndpointTokenURL()+"?"+query, contentTypeURL, nil, result)
	if err != nil {
		return nil, fmt.Errorf(ErrSendRequest, "AccessTokenExchange", err)
	}

	return result, nil
}

// Send sends a RefreshTokenExchange request to Discord and returns an AccessTokenResponse.
//
// Uses the RefreshTokenExchange ClientID and ClientSecret.
func (r *RefreshTokenExchange) Send(bot *Client) (*AccessTokenResponse, error) {
	query, err := EndpointQueryString(r)
	if err != nil {
		return nil, fmt.Errorf(ErrQueryString, "RefreshTokenExchange", err)
	}

	var result *AccessTokenResponse
	err = SendRequest(bot, fasthttp.MethodPost, EndpointTokenURL()+"?"+query, contentTypeURL, nil, result)
	if err != nil {
		return nil, fmt.Errorf(ErrSendRequest, "RefreshTokenExchange", err)
	}

	return result, nil
}

// Send sends a ClientCredentialsTokenRequest to Discord and returns a ClientCredentialsTokenRequest.
func (r *ClientCredentialsTokenRequest) Send(bot *Client) (*AccessTokenResponse, error) {
	query, err := EndpointQueryString(r)
	if err != nil {
		return nil, fmt.Errorf(ErrQueryString, "ClientCredentialsTokenRequest", err)
	}

	var result *AccessTokenResponse
	err = SendRequest(bot, fasthttp.MethodPost, EndpointTokenURL()+"?"+query, contentTypeURL, nil, result)
	if err != nil {
		return nil, fmt.Errorf(ErrSendRequest, "ClientCredentialsTokenRequest", err)
	}

	return result, nil
}

// urlQueryStringScope parses a given slice of scopes to generate a valid URL Query String.
func urlQueryStringScope(scopes []string) string {
	if len(scopes) > 0 {
		var scope strings.Builder
		scope.WriteString("scope=")

		for i, s := range scopes {
			if i > 0 {
				scope.WriteString("%20")
			}

			scope.WriteString(s)
		}

		return scope.String()
	}

	return ""
}
