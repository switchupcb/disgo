package requests

import "github.com/switchupcb/disgo/wrapper/resources"

// Create Stage Instance
// POST /stage-instances
// https://discord.com/developers/docs/resources/stage-instance#create-stage-instance
type CreateStageInstance struct {
	ChannelID             resources.Snowflake
	Topic                 string         `json:"topic,omitempty"`
	PrivacyLevel          resources.Flag `json:"privacy_level,omitempty"`
	SendStartNotification bool           `json:"send_start_notification,omitempty"`
}

// Get Stage Instance
// GET /stage-instances/{channel.id}
// https://discord.com/developers/docs/resources/stage-instance#get-stage-instance
type GetStageInstance struct {
	ChannelID resources.Snowflake
}

// Modify Stage Instance
// PATCH /stage-instances/{channel.id}
// https://discord.com/developers/docs/resources/stage-instance#modify-stage-instance
type ModifyStageInstance struct {
	ChannelID    resources.Snowflake
	Topic        string         `json:"topic,omitempty"`
	PrivacyLevel resources.Flag `json:"privacy_level,omitempty"`
}

// Delete Stage Instance
// DELETE /stage-instances/{channel.id}
// https://discord.com/developers/docs/resources/stage-instance#delete-stage-instance
type DeleteStageInstance struct {
	ChannelID resources.Snowflake
}
