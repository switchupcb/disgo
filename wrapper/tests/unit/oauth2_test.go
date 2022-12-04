package unit_test

import (
	"testing"

	. "github.com/switchupcb/disgo"
)

// testOAuth2 represents parameters used to test GenerateAuthorizationURL.
type testOAuth2 struct {
	name     string
	bot      *Client
	response string
	output   string
}

// testOAuth2Bot represents parameters used to test GenerateBotAuthorizationURL.
type testOAuth2Bot struct {
	name   string
	p      BotAuthParams
	output string
}

// TestGenerateAuthorizationURL tests GenerateAuthorizationURL for a valid authorization URL.
func TestGenerateAuthorizationURL(t *testing.T) {
	tests := []testOAuth2{
		{
			name: "test1",
			bot: &Client{
				Authorization: &Authorization{
					ClientID:     "988406086973444138",
					ClientSecret: "",
					RedirectURI:  "https://localhost",
					State:        "state",
					Prompt:       "prompt",
					Scopes:       []string{"scope1", "scope2"},
				},
			},
			response: "code",
			output:   EndpointAuthorizationURL() + "?responsetype=code&client_id=988406086973444138&scope=scope1%20scope2&redirect_uri=https%3A%2F%2Flocalhost&state=state&prompt=prompt",
		},
		{
			name: "test2",
			bot: &Client{
				Authorization: &Authorization{
					ClientID:     "983406086973444138",
					ClientSecret: "",
					RedirectURI:  "https://localhost",
					State:        "",
					Prompt:       "",
					Scopes:       []string{"bot", "applications.commands"},
				},
			},
			response: "",
			output:   EndpointAuthorizationURL() + "?client_id=983406086973444138&scope=bot%20applications.commands&redirect_uri=https%3A%2F%2Flocalhost",
		},
		{
			name: "test3",
			bot: &Client{
				Authorization: &Authorization{
					ClientID:     "406086983973444138",
					ClientSecret: "",
					RedirectURI:  "https://localhost",
					State:        "",
					Prompt:       "",
					Scopes:       []string{},
				},
			},
			response: "",
			output:   EndpointAuthorizationURL() + "?client_id=406086983973444138&redirect_uri=https%3A%2F%2Flocalhost",
		},
	}

	for _, test := range tests {
		got := GenerateAuthorizationURL(test.bot, test.response)
		if got != test.output {
			t.Errorf("(%v: got %v, wanted %v", test.name, got, test.output)
		}
	}
}

// TestGenerateBotAuthorizationURL tests GenerateBotAuthorizationURL for a valid bot authorization URL.
func TestGenerateBotAuthorizationURL(t *testing.T) {
	tests := []testOAuth2Bot{
		{
			name: "botTest1",
			p: BotAuthParams{
				Bot: &Client{
					Authorization: &Authorization{
						ClientID:     "988406086973444138",
						ClientSecret: "",
						RedirectURI:  "https://github.com",
						State:        "state",
						Prompt:       "prompt",
						Scopes:       []string{"scope1", "scope2"},
					},
				},
				ResponseType:       "",
				Permissions:        1,
				GuildID:            "9862629927296967350",
				DisableGuildSelect: false,
			},
			output: EndpointAuthorizationURL() + "?client_id=988406086973444138&scope=scope1%20scope2&redirect_uri=https%3A%2F%2Fgithub.com&state=state&prompt=prompt&permissions=1&guild_id=9862629927296967350&disable_guild_select=false",
		},
		{
			name: "botTest2",
			p: BotAuthParams{
				Bot: &Client{
					Authorization: &Authorization{
						ClientID:     "983406086973444138",
						ClientSecret: "",
						RedirectURI:  "https://github.com",
						State:        "",
						Prompt:       "",
						Scopes:       []string{"bot", "applications.commands"},
					},
				},
				GuildID:            "2909267986263572999",
				ResponseType:       "code",
				Permissions:        10,
				DisableGuildSelect: true,
			},
			output: EndpointAuthorizationURL() + "?responsetype=code&client_id=983406086973444138&scope=bot%20applications.commands&redirect_uri=https%3A%2F%2Fgithub.com&permissions=10&guild_id=2909267986263572999&disable_guild_select=true",
		},
		{
			name: "botTest3",
			p: BotAuthParams{
				Bot: &Client{
					Authorization: &Authorization{
						ClientID:     "406086983973444138",
						ClientSecret: "",
						RedirectURI:  "https://localhost",
						State:        "",
						Prompt:       "",
						Scopes:       []string{},
					},
				},
				GuildID:            "2909267986263572999",
				ResponseType:       "code",
				Permissions:        10,
				DisableGuildSelect: true,
			},
			output: EndpointAuthorizationURL() + "?responsetype=code&client_id=406086983973444138&redirect_uri=https%3A%2F%2Flocalhost&permissions=10&guild_id=2909267986263572999&disable_guild_select=true",
		},
	}

	for _, test := range tests {
		got := GenerateBotAuthorizationURL(test.p)
		if got != test.output {
			t.Errorf("(%v: got %v, wanted %v", test.name, got, test.output)
		}
	}
}
