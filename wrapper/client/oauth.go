package client

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/switchupcb/disgo/wrapper/resources"
)

const baseURL = "https://discord.com/api/oauth2/authorize?client_id="
const webhookURL = "https://discord.com/api/oauth2/authorize?response_type="

type BotAuthorization struct {
	ClientID           resources.Snowflake
	Scope              string
	Permissions        int
	GuildID            resources.Snowflake
	DisableGuildSelect bool
}

type WebhookAutorization struct {
	ResponseType string
	ClientID     resources.Snowflake
	Scope        string
	State        string
	RedirectURI  string
}

// GenerateBotURL generates the URL to be sent for connection.
func GenerateBotURL(bot BotAuthorization) string {

	return baseURL + string(rune(bot.ClientID)) + "&scope=bot&permissions=" + string(rune(bot.Permissions))
}

func GenerateWebhookURL(bot WebhookAutorization) string {
	return webhookURL + bot.ResponseType + "&client_id=" + string(rune(bot.ClientID)) + "scope=webhook.incoming&state=" + bot.State + "&redirect_uri=" + bot.RedirectURI
}

// ConnectOauth registers developer application and retrieves the client ID and secret key.
func ConnectOauth(bot Client) {

	BotAuthorization := BotAuthorization{
		ClientID: bot.ApplicationID,
	}

	url := GenerateBotURL(BotAuthorization)

	resp, err := http.Get(url)

	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatalln(err)
	}

	sb := string(body)
	log.Printf(sb)
}

func SendWebhook(bot Client) {
	WebhookAutorization := WebhookAutorization{
		ClientID: bot.ApplicationID,
	}

	url := GenerateWebhookURL(WebhookAutorization)

	resp, err := http.Get(url)

	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatalln(err)
	}

	sb := string(body)
	log.Printf(sb)
}

// OAuth2 URLs
// https://discord.com/api/oauth2/authorize
// https://discord.com/api/oauth2/token
// https://discord.com/api/oauth2/token/revoke

// authorization URL example
// https://discord.com/api/oauth2/authorize?response_type=code&client_id=157730590492196864&scope=identify%20guilds.join&state=15773059ghq9183habn&redirect_uri=https%3A%2F%2Fnicememe.website&prompt=consent

// redirect URL example
// https://nicememe.website/?code=NhhvTDYsFcdgNLnnLijcl7Ku7bEEeee&state=15773059ghq9183habn
