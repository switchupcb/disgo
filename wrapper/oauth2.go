package wrapper

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ContentTypeURL represents an HTTP header that indicates a encoded URL.
var ContentTypeURL = []byte("application/x-www-form-urlencoded")

// AuthorizationCodeGrant performs an OAuth2 authorization code grant.
//https://discord.com/developers/docs/topics/oauth2#authorization-code-grant
func AuthorizationCodeGrant(bot *Client) error {
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

	// TODO
	// set header := ContentTypeURL and data where necessary.
	// fix endpoints for send()

	// retrieve an access code.
	authurl := AuthorizationURL{
		ResponseType: "code",
		ClientID:     bot.Authorization.ClientID,
		Scope:        scope.String(),
		State:        bot.Authorization.State,
		RedirectURI:  bot.Authorization.RedirectURI,
		Prompt:       bot.Authorization.Prompt,
	}

	redirecturl, err := authurl.SendAuthorizationCodeGrantURL(bot)
	if err != nil {
		return err
	}

	// exchange the access code for a user's access token.
	accesstokenexchange := AccessTokenExchange{
		ClientID:     bot.Authorization.ClientID,
		ClientSecret: bot.Authorization.ClientSecret,
		GrantType:    "authorization_code",
		Code:         redirecturl.Code,
		RedirectURI:  bot.Authorization.RedirectURI,
	}

	// set bot AccessToken information.
	bot.AccessToken, err = accesstokenexchange.SendAccessTokenExchange(bot)
	if err != nil {
		return err
	}

	return nil
}

// AuthorizationURL generates an authorization URL from a given client
func EndAuthorizationURL(bot *Client) string {
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

// SendAuthorizationCodeGrantURL sends a AuthorizationURL to Discord and returns a RedirectURL.
func (r *AuthorizationURL) SendAuthorizationCodeGrantURL(bot *Client) (*RedirectURL, error) {
	var result *RedirectURL
	body, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf(ErrSendMarshal, "AuthorizationURL", err)
	}

	err = SendRequest(result, bot.client, TODO, EndpointAuthorizationURL(), body)
	if err != nil {
		return nil, fmt.Errorf(ErrSendRequest, "AuthorizationURL", err)
	}

	return result, nil
}

// SendAccessTokenExchange sends a AccessTokenExchange to Discord and returns a AccessTokenResponse.
func (r *AccessTokenExchange) SendAccessTokenExchange(bot *Client) (*AccessTokenResponse, error) {
	var result *AccessTokenResponse
	body, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf(ErrSendMarshal, "AccessTokenExchange", err)
	}

	err = SendRequest(result, bot.client, TODO, EndpointTokenURL(), body)
	if err != nil {
		return nil, fmt.Errorf(ErrSendRequest, "AccessTokenExchange", err)
	}

	return result, nil
}

// SendRefreshTokenExchange sends a RefreshTokenExchange to Discord and returns a AccessTokenResponse.
func (r *RefreshTokenExchange) SendRefreshTokenExchange(bot *Client) (*AccessTokenResponse, error) {
	var result *AccessTokenResponse
	body, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf(ErrSendMarshal, "RefreshTokenExchange", err)
	}

	err = SendRequest(result, bot.client, TODO, EndpointTokenURL(), body)
	if err != nil {
		return nil, fmt.Errorf(ErrSendRequest, "RefreshTokenExchange", err)
	}

	return result, nil
}

// SendImplicitAuthorizationURL sends a AuthorizationURL to Discord and returns a RedirectURI.
func (r *AuthorizationURL) SendImplicitAuthorizationURL(bot *Client) (*RedirectURI, error) {
	var result *RedirectURI
	body, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf(ErrSendMarshal, "ImplicitAuthorizationURL", err)
	}

	err = SendRequest(result, bot.client, TODO, EndpointAuthorizationURL(), body)
	if err != nil {
		return nil, fmt.Errorf(ErrSendRequest, "ImplicitAuthorizationURL", err)
	}

	return result, nil
}

// SendBotAuth sends a BotAuth to Discord and returns a error.
func (r *BotAuth) SendBotAuth(bot *Client) error {
	var result error
	body, err := json.Marshal(r)
	if err != nil {
		return fmt.Errorf(ErrSendMarshal, "BotAuth", err)
	}

	err = SendRequest(result, bot.client, TODO, EndpointAuthorizationURL(), body)
	if err != nil {
		return fmt.Errorf(ErrSendRequest, "BotAuth", err)
	}

	return nil
}

// SendClientCredentialsTokenRequest sends a ClientCredentialsTokenRequest to Discord and returns a ClientCredentialsTokenRequest.
func (r *ClientCredentialsTokenRequest) SendClientCredentialsTokenRequest(bot *Client) (*ClientCredentialsTokenRequest, error) {
	var result *ClientCredentialsTokenRequest
	body, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf(ErrSendMarshal, "ClientCredentialsTokenRequest", err)
	}

	err = SendRequest(result, bot.client, TODO, EndpointTokenURL(), body)
	if err != nil {
		return nil, fmt.Errorf(ErrSendRequest, "ClientCredentialsTokenRequest", err)
	}

	return result, nil
}

// SendAdvancedBotAuth sends a AuthorizationURL to Discord and returns a ExtendedBotAuthorizationAccessTokenResponse.
func (r *AuthorizationURL) SendAdvancedBotAuth(bot *Client) (*ExtendedBotAuthorizationAccessTokenResponse, error) {
	var result *ExtendedBotAuthorizationAccessTokenResponse
	body, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf(ErrSendMarshal, "AdvancedBotAuth", err)
	}

	err = SendRequest(result, bot.client, TODO, EndpointAuthorizationURL(), body)
	if err != nil {
		return nil, fmt.Errorf(ErrSendRequest, "AdvancedBotAuth", err)
	}

	return result, nil
}

// SendWebhookAuth sends a AuthorizationURL to Discord and returns a WebhookTokenResponse.
func (r *AuthorizationURL) SendWebhookAuth(bot *Client) (*WebhookTokenResponse, error) {
	var result *WebhookTokenResponse
	body, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf(ErrSendMarshal, "WebhookAuth", err)
	}

	err = SendRequest(result, bot.client, TODO, EndpointAuthorizationURL(), body)
	if err != nil {
		return nil, fmt.Errorf(ErrSendRequest, "WebhookAuth", err)
	}

	return result, nil
}
