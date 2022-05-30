package wrapper

import (
	"fmt"
	"strings"

	"github.com/valyala/fasthttp"
)

const (
	grantTypeAuthorizationCodeGrant = "authorization_code"
	grantTypeRefreshToken           = "refresh_token"
)

// AuthorizationCodeGrant performs an OAuth2 authorization code grant.
//
// Send the user a valid Authorization URL, which can be generated using GenerateAuthorizationURL(bot).
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
func RefreshAuthorizationCodeGrant(bot *Client, token AccessTokenResponse) (*AccessTokenResponse, error) {
	exchange := &RefreshTokenExchange{
		ClientID:     bot.Authorization.ClientID,
		ClientSecret: bot.Authorization.ClientSecret,
		GrantType:    grantTypeRefreshToken,
		RefreshToken: token.RefreshToken,
	}

	return exchange.Send(bot)
}

// GenerateAuthorizationURL generates an authorization URL from a given client.
func GenerateAuthorizationURL(bot *Client) string {
	params := make([]string, 0, 5)

	// client_id is the application client id.
	params = append(params, "client_id="+bot.Authorization.ClientID)

	// scope is a list of OAuth2 scopes separated by url encoded spaces (%20).
	var scope strings.Builder
	if len(bot.Authorization.Scopes) > 0 {
		scope.WriteString("scope=")

		for i, s := range bot.Authorization.Scopes {
			if i > 0 {
				scope.WriteString("%20")
			}

			scope.WriteString(s)
		}
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

	return EndpointAuthorizationURL() + "?response_type=code&" + strings.Join(params, "&")
}

// Send sends an AccessTokenExchange request to Discord and returns an AccessTokenResponse.
func (r *AccessTokenExchange) Send(bot *Client) (*AccessTokenResponse, error) {
	query, err := EndpointQueryString(r)
	if err != nil {
		return nil, fmt.Errorf(ErrQueryString, "AccessTokenExchange", err)
	}

	var result *AccessTokenResponse
	err = SendRequest(bot.client, fasthttp.MethodPost, EndpointTokenURL()+"?"+query, contentTypeURL, nil, result)
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
	err = SendRequest(bot.client, fasthttp.MethodPost, EndpointTokenURL()+"?"+query, contentTypeURL, nil, result)
	if err != nil {
		return nil, fmt.Errorf(ErrSendRequest, "RefreshTokenExchange", err)
	}

	return result, nil
}

// SendImplicitAuthorizationURL sends a AuthorizationURL to Discord and returns a RedirectURI.
func (r *AuthorizationURL) SendImplicitAuthorizationURL(bot *Client) (*RedirectURI, error) {
	var result *RedirectURI
	err := SendRequest(bot.client, fasthttp.MethodGet, GenerateAuthorizationURL(bot), nil, nil, result)
	if err != nil {
		return nil, fmt.Errorf(ErrSendRequest, "ImplicitAuthorizationURL", err)
	}

	//	fmt.Println(u.Fragment)

	return result, nil
}

// SendBotAuth sends a BotAuth to Discord and returns a error.
func (r *BotAuth) SendBotAuth(bot *Client) error {
	var result error
	query, err := EndpointQueryString(r)
	if err != nil {
		return fmt.Errorf(ErrQueryString, "BotAuth", err)
	}

	err = SendRequest(bot.client, fasthttp.MethodGet, EndpointAuthorizationURL()+"?"+query, nil, nil, result)
	if err != nil {
		return fmt.Errorf(ErrSendRequest, "BotAuth", err)
	}

	return nil
}

// SendClientCredentialsTokenRequest sends a ClientCredentialsTokenRequest to Discord and returns a ClientCredentialsTokenRequest.
func (r *ClientCredentialsTokenRequest) SendClientCredentialsTokenRequest(bot *Client) (*ClientCredentialsTokenRequest, error) {
	var result *ClientCredentialsTokenRequest
	query, err := EndpointQueryString(r)
	if err != nil {
		return nil, fmt.Errorf(ErrQueryString, "ClientCredentialsTokenRequest", err)
	}

	err = SendRequest(bot.client, fasthttp.MethodPost, EndpointTokenURL()+"?"+query, contentTypeURL, nil, result)
	if err != nil {
		return nil, fmt.Errorf(ErrSendRequest, "ClientCredentialsTokenRequest", err)
	}

	return result, nil
}

// SendAdvancedBotAuth sends a AuthorizationURL to Discord and returns a ExtendedBotAuthorizationAccessTokenResponse.
func (r *AuthorizationURL) SendAdvancedBotAuth(bot *Client) (*ExtendedBotAuthorizationAccessTokenResponse, error) {
	var result *ExtendedBotAuthorizationAccessTokenResponse
	err := SendRequest(bot.client, fasthttp.MethodGet, GenerateAuthorizationURL(bot), nil, nil, result)
	if err != nil {
		return nil, fmt.Errorf(ErrSendRequest, "AdvancedBotAuth", err)
	}

	return result, nil
}

// SendWebhookAuth sends a AuthorizationURL to Discord and returns a WebhookTokenResponse.
func (r *AuthorizationURL) SendWebhookAuth(bot *Client) (*WebhookTokenResponse, error) {
	var result *WebhookTokenResponse
	err := SendRequest(bot.client, fasthttp.MethodGet, GenerateAuthorizationURL(bot), nil, nil, result)
	if err != nil {
		return nil, fmt.Errorf(ErrSendRequest, "WebhookAuth", err)
	}

	return result, nil
}
