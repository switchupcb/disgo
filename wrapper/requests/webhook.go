package requests

import "github.com/switchupcb/disgo/wrapper/resources"

// Create Webhook
// POST /channels/{channel.id}/webhooks
// https://discord.com/developers/docs/resources/webhook#create-webhook
type CreateWebhook struct {
	Name   string `json:"name,omitempty"`
	Avatar string `json:"avatar,omitempty"`
}

// Get Channel Webhooks
// GET /channels/{channel.id}/webhooks
// https://discord.com/developers/docs/resources/webhook#get-channel-webhooks
type GetChannelWebhooks struct {
	ChannelID resources.Snowflake
}

// Get Guild Webhooks
// GET /guilds/{guild.id}/webhooks
// https://discord.com/developers/docs/resources/webhook#get-guild-webhooks
type GetGuildWebhooks struct {
	GuildID resources.Snowflake
}

// Get Webhook
// GET /webhooks/{webhook.id}
// https://discord.com/developers/docs/resources/webhook#get-webhook
type GetWebhook struct {
	WebhookID resources.Snowflake
}

// Get Webhook with Token
// GET /webhooks/{webhook.id}/{webhook.token}
// https://discord.com/developers/docs/resources/webhook#get-webhook-with-token
type GetWebhookwithToken struct {
	WebhookID resources.Snowflake
}

// Modify Webhook
// PATCH /webhooks/{webhook.id}
// https://discord.com/developers/docs/resources/webhook#modify-webhook
type ModifyWebhook struct {
	WebhookID resources.Snowflake
	Name      string              `json:"name,omitempty"`
	Avatar    string              `json:"avatar,omitempty"`
	ChannelID resources.Snowflake `json:"channel_id,omitempty"`
}

// Modify Webhook with Token
// PATCH /webhooks/{webhook.id}/{webhook.token}
// https://discord.com/developers/docs/resources/webhook#modify-webhook-with-token
type ModifyWebhookwithToken struct {
	WebhookID resources.Snowflake
}

// Delete Webhook
// DELETE /webhooks/{webhook.id}
// https://discord.com/developers/docs/resources/webhook#delete-webhook
type DeleteWebhook struct {
	WebhookID resources.Snowflake
}

// Delete Webhook with Token
// DELETE /webhooks/{webhook.id}/{webhook.token}
// https://discord.com/developers/docs/resources/webhook#delete-webhook-with-token
type DeleteWebhookwithToken struct {
	WebhookID resources.Snowflake
}

// Execute Webhook
// POST /webhooks/{webhook.id}/{webhook.token}
// https://discord.com/developers/docs/resources/webhook#execute-webhook
type ExecuteWebhook struct {
	Wait            bool                       `json:"wait,omitempty"`
	ThreadID        resources.Snowflake        `json:"thread_id,omitempty"`
	Content         string                     `json:"content,omitempty"`
	Username        string                     `json:"username,omitempty"`
	AvatarURL       string                     `json:"avatar_url,omitempty"`
	TTS             bool                       `json:"tts,omitempty"`
	Files           []byte                     `disgo:"files"`
	Components      []resources.Component      `json:"components,omitempty"`
	Embeds          []*resources.Embed         `json:"embeds,omitempty"`
	AllowedMentions *resources.AllowedMentions `json:"allowed_mentions,omitempty"`
	PayloadJSON     string                     `json:"payload_json,omitempty"`
	Attachments     []*resources.Attachment    `json:"attachments,omitempty"`
}

// Execute Slack-Compatible Webhook
// POST /webhooks/{webhook.id}/{webhook.token}/slack
// https://discord.com/developers/docs/resources/webhook#execute-slackcompatible-webhook
type ExecuteSlackCompatibleWebhook struct {
	ThreadID resources.Snowflake
	Wait     bool `json:"wait,omitempty"`
}

// Execute GitHub-Compatible Webhook
// POST /webhooks/{webhook.id}/{webhook.token}/github
// https://discord.com/developers/docs/resources/webhook#execute-githubcompatible-webhook
type ExecuteGitHubCompatibleWebhook struct {
	ThreadID resources.Snowflake
	Wait     bool `json:"wait,omitempty"`
}

// Get Webhook Message
// GET /webhooks/{webhook.id}/{webhook.token}/messages/{message.id}
// https://discord.com/developers/docs/resources/webhook#get-webhook-message
type GetWebhookMessage struct {
	ThreadID resources.Snowflake
}

// Edit Webhook Message
// PATCH /webhooks/{webhook.id}/{webhook.token}/messages/{message.id}
// https://discord.com/developers/docs/resources/webhook#edit-webhook-message
type EditWebhookMessage struct {
	WebhookID       resources.Snowflake
	ThreadID        resources.Snowflake        `json:"thread_id,omitempty"`
	Content         *string                    `json:"content,omitempty"`
	Components      []*resources.Component     `json:"components,omitempty"`
	Embeds          []*resources.Embed         `json:"embeds,omitempty"`
	Files           []byte                     `disgo:"files"`
	AllowedMentions *resources.AllowedMentions `json:"allowed_mentions,omitempty"`
	PayloadJSON     string                     `json:"payload_json,omitempty"`
	Attachments     []*resources.Attachment    `json:"attachments,omitempty"`
}

// Delete Webhook Message
// DELETE /webhooks/{webhook.id}/{webhook.token}/messages/{message.id}
// https://discord.com/developers/docs/resources/webhook#delete-webhook-message
type DeleteWebhookMessage struct {
	ThreadID resources.Snowflake
}
