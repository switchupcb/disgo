package main

import (
	"github.com/switchupcb/disgo"
	"github.com/valyala/fasthttp"
)

func main() {
	// set the bot's configuration during the instantiation of a Client.
	bot := &disgo.Client{
		ApplicationID: "",

		// Used for basic Authentication.
		Authentication: &disgo.Authentication{
			Token:     "",
			TokenType: "",
			Header:    "",
		},

		// Used for OAuth2 Authorization.
		Authorization: &disgo.Authorization{
			ClientID:     "",
			ClientSecret: "",
			RedirectURI:  "",
			State:        "",
			Prompt:       "",
			Scopes:       []string{},
		},

		// Used to configure HTTP Request and WebSocket Connection parameters.
		Config: &disgo.Config{
			Request: disgo.Request{

				// RateLimiter is an interface defined in ./wrapper/ratelimiter.go
				RateLimiter: new(disgo.RateLimit),

				// see https://pkg.go.dev/github.com/valyala/fasthttp#Client
				Client:      &fasthttp.Client{},
				Timeout:     0,
				Retries:     0,
				RetryShared: false,
			},

			Gateway: disgo.Gateway{

				// RateLimiter is an interface defined in ./wrapper/ratelimiter.go
				RateLimiter: new(disgo.RateLimit),

				GatewayPresenceUpdate: &disgo.GatewayPresenceUpdate{
					Since:  nil,
					Status: "",
					Game:   []*disgo.Activity{},
					AFK:    false,
				},

				// It's not recommended to modify these fields directly.
				//
				// Instead, use Automatic Gateway Intents, EnableIntent, or DisableIntent
				// described in ./_contribution/EVENTS.md
				IntentSet: make(map[disgo.BitFlag]bool, 0),
				Intents:   0,
			},
		},

		// Handlers controls the bot's event handlers.
		//
		// It's recommended to manage event handlers using Handle() and Remove()
		// described in ./_contribution/REQUESTS.md
		Handlers: new(disgo.Handlers),

		// Sessions controls the bot's WebSocket Sessions (Gateway, Voice).
		Sessions: []*disgo.Session{},
	}

	// set a configuration option during runtime.
	bot.ApplicationID = ""
}
