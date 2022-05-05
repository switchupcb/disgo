package requests

import "github.com/switchupcb/disgo/wrapper/resources"

// Create Interaction Response
// POST /interactions/{interaction.id}/{interaction.token}/callback
// https://discord.com/developers/docs/interactions/receiving-and-responding#create-interaction-response
type CreateInteractionResponse struct {
	InteractionToken string `json:"token,omitempty"`
	InteractionID    resources.Snowflake
}

// Get Original Interaction Response
// GET /webhooks/{application.id}/{interaction.token}/messages/@original
// https://discord.com/developers/docs/interactions/receiving-and-responding#get-original-interaction-response
type GetOriginalInteractionResponse struct {
	ApplicationID    resources.Snowflake
	InteractionToken string `json:"token,omitempty"`
}

// Edit Original Interaction Response
// PATCH /webhooks/{application.id}/{interaction.token}/messages/@original
// https://discord.com/developers/docs/interactions/receiving-and-responding#edit-original-interaction-response
type EditOriginalInteractionResponse struct {
	ApplicationID    resources.Snowflake
	InteractionToken string `json:"token,omitempty"`
}

// Delete Original Interaction Response
// DELETE /webhooks/{application.id}/{interaction.token}/messages/@original
// https://discord.com/developers/docs/interactions/receiving-and-responding#delete-original-interaction-response
type DeleteOriginalInteractionResponse struct {
	ApplicationID    resources.Snowflake
	InteractionToken string `json:"token,omitempty"`
}

// Create Followup Message
// POST /webhooks/{application.id}/{interaction.token}
// https://discord.com/developers/docs/interactions/receiving-and-responding#create-followup-message
type CreateFollowupMessage struct {
	InteractionToken string `json:"token,omitempty"`
	ApplicationID    resources.Snowflake
}

// Get Followup Message
// GET /webhooks/{application.id}/{interaction.token}/messages/{message.id}
// https://discord.com/developers/docs/interactions/receiving-and-responding#get-followup-message
type GetFollowupMessage struct {
	ApplicationID    resources.Snowflake
	MessageID        resources.Snowflake
	InteractionToken string `json:"token,omitempty"`
}

// Edit Followup Message
// PATCH /webhooks/{application.id}/{interaction.token}/messages/{message.id}
// https://discord.com/developers/docs/interactions/receiving-and-responding#edit-followup-message
type EditFollowupMessage struct {
	MessageID        resources.Snowflake
	ApplicationID    resources.Snowflake
	InteractionToken string `json:"token,omitempty"`
}

// Delete Followup Message
// DELETE /webhooks/{application.id}/{interaction.token}/messages/{message.id}
// https://discord.com/developers/docs/interactions/receiving-and-responding#delete-followup-message
type DeleteFollowupMessage struct {
	MessageID        resources.Snowflake
	ApplicationID    resources.Snowflake
	InteractionToken string `json:"token,omitempty"`
}
