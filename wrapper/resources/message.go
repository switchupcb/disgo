package resources

import "time"

// Message Object
// https://discord.com/developers/docs/resources/channel#message-object
type Message struct {
	ID                Snowflake         `json:"id,omitempty"`
	ChannelID         Snowflake         `json:"channel_id,omitempty"`
	GuildID           Snowflake         `json:"guild_id,omitempty"`
	Author            *User             `json:"author,omitempty"`
	Member            *GuildMember      `json:"member,omitempty"`
	Content           string            `json:"content,omitempty"`
	Timestamp         time.Time         `json:"timestamp,omitempty"`
	EditedTimestamp   time.Time         `json:"edited_timestamp,omitempty"`
	TTS               bool              `json:"tts,omitempty"`
	MentionEveryone   bool              `json:"mention_everyone,omitempty"`
	Mentions          []*User           `json:"mentions,omitempty"`
	MentionRoles      []Snowflake       `json:"mention_roles,omitempty"`
	MentionChannels   []*ChannelMention `json:"mention_channels,omitempty"`
	Attachments       []*Attachment     `json:"attachments,omitempty"`
	Embeds            []*Embed          `json:"embeds,omitempty"`
	Reactions         []*Reaction       `json:"reactions,omitempty"`
	Nonce             interface{}       `json:"nonce,omitempty"`
	Pinned            bool              `json:"pinned,omitempty"`
	WebhookID         Snowflake         `json:"webhook_id,omitempty"`
	Type              Flag              `json:"type,omitempty"`
	Activity          MessageActivity   `json:"activity,omitempty"`
	Application       *Application      `json:"application,omitempty"`
	MessageReference  *MessageReference `json:"message_reference,omitempty"`
	Flags             CodeFlag          `json:"flags,omitempty"`
	ReferencedMessage *Message          `json:"referenced_message,omitempty"`
	Interaction       *Interaction      `json:"interaction,omitempty"`
	Thread            *Channel          `json:"thread,omitempty"`
	Components        []*Component      `json:"components,omitempty"`
	StickerItems      []*StickerItem    `json:"sticker_items,omitempty"`
}

// Message Types
// https://discord.com/developers/docs/resources/channel#message-object-message-types
const (
	FlagTypesMessageDEFAULT                                      = 0
	FlagTypesMessageRECIPIENT_ADD                                = 1
	FlagTypesMessageRECIPIENT_REMOVE                             = 2
	FlagTypesMessageCALL                                         = 3
	FlagTypesMessageCHANNEL_NAME_CHANGE                          = 4
	FlagTypesMessageCHANNEL_ICON_CHANGE                          = 5
	FlagTypesMessageCHANNEL_PINNED_MESSAGE                       = 6
	FlagTypesMessageGUILD_MEMBER_JOIN                            = 7
	FlagTypesMessageUSER_PREMIUM_GUILD_SUBSCRIPTION              = 8
	FlagTypesMessageUSER_PREMIUM_GUILD_SUBSCRIPTION_TIER_ONE     = 9
	FlagTypesMessageUSER_PREMIUM_GUILD_SUBSCRIPTION_TIER_TWO     = 10
	FlagTypesMessageUSER_PREMIUM_GUILD_SUBSCRIPTION_TIER_THREE   = 11
	FlagTypesMessageCHANNEL_FOLLOW_ADD                           = 12
	FlagTypesMessageGUILD_DISCOVERY_DISQUALIFIED                 = 14
	FlagTypesMessageGUILD_DISCOVERY_REQUALIFIED                  = 15
	FlagTypesMessageGUILD_DISCOVERY_GRACE_PERIOD_INITIAL_WARNING = 16
	FlagTypesMessageGUILD_DISCOVERY_GRACE_PERIOD_FINAL_WARNING   = 17
	FlagTypesMessageTHREAD_CREATED                               = 18
	FlagTypesMessageREPLY                                        = 19
	FlagTypesMessageCHAT_INPUT_COMMAND                           = 20
	FlagTypesMessageTHREAD_STARTER_MESSAGE                       = 21
	FlagTypesMessageGUILD_INVITE_REMINDER                        = 22
	FlagTypesMessageCONTEXT_MENU_COMMAND                         = 23
)

// Message Activity Structure
// https://discord.com/developers/docs/resources/channel#message-object-message-activity-structure
type MessageActivity struct {
	Type    int    `json:"type,omitempty"`
	PartyID string `json:"party_id,omitempty"`
}

// Message Activity Types
// https://discord.com/developers/docs/resources/channel#message-object-message-activity-types
const (
	FlagTypesActivityMessageJOIN         = 1
	FlagTypesActivityMessageSPECTATE     = 2
	FlagTypesActivityMessageLISTEN       = 3
	FlagTypesActivityMessageJOIN_REQUEST = 5
)

// Message Flags
// https://discord.com/developers/docs/resources/channel#message-object-message-flags
const (
	FlagFlagsMessageCROSSPOSTED                            = 1 << 0
	FlagFlagsMessageIS_CROSSPOST                           = 1 << 1
	FlagFlagsMessageSUPPRESS_EMBEDS                        = 1 << 2
	FlagFlagsMessageSOURCE_MESSAGE_DELETED                 = 1 << 3
	FlagFlagsMessageURGENT                                 = 1 << 4
	FlagFlagsMessageHAS_THREAD                             = 1 << 5
	FlagFlagsMessageEPHEMERAL                              = 1 << 6
	FlagFlagsMessageLOADING                                = 1 << 7
	FlagFlagsMessageFAILED_TO_MENTION_SOME_ROLES_IN_THREAD = 1 << 8
)

// Message Reference Object
// https://discord.com/developers/docs/resources/channel#message-reference-object
type MessageReference struct {
	MessageID       Snowflake `json:"message_id,omitempty"`
	ChannelID       Snowflake `json:"channel_id,omitempty"`
	GuildID         Snowflake `json:"guild_id,omitempty"`
	FailIfNotExists bool      `json:"fail_if_not_exists,omitempty"`
}

// Message Attachment Object
// https://discord.com/developers/docs/resources/channel#attachment-object
type Attachment struct {
	ID       Snowflake `json:"id,omitempty"`
	Filename string    `json:"filename,omitempty"`
	Size     uint      `json:"size,omitempty"`
	URL      string    `json:"url,omitempty"`
	ProxyURL string    `json:"proxy_url,omitempty"`
	Height   uint      `json:"height,omitempty"`
	Width    uint      `json:"width,omitempty"`

	SpoilerTag bool `json:"-,omitempty"`
}

// Sticker Structure
// https://discord.com/developers/docs/resources/sticker#sticker-object-sticker-structure
type Sticker struct {
	ID          Snowflake `json:"id,omitempty"`
	PackID      Snowflake `json:"pack_id,omitempty"`
	Name        string    `json:"name,omitempty"`
	Description string    `json:"description,omitempty"`
	Tags        string    `json:"tags,omitempty"`
	Asset       string    `json:"asset,omitempty"`
	Type        Flag      `json:"type,omitempty"`
	FormatType  Flag      `json:"format_type,omitempty"`
	Available   bool      `json:"available,omitempty"`
	GuildID     Snowflake `json:"guild_id,omitempty"`
	User        *User     `json:"user,omitempty"`
	SortValue   int       `json:"sort_value,omitempty"`
}

// Sticker Types
// https://discord.com/developers/docs/resources/sticker#sticker-object-sticker-types
const (
	FlagTypesStickerSTANDARD = 1
	FlagTypesStickerGUILD    = 2
)

// Sticker Format Types
// https://discord.com/developers/docs/resources/sticker#sticker-object-sticker-format-types
const (
	FlagTypesFormatStickerPNG    = 1
	FlagTypesFormatStickerAPNG   = 2
	FlagTypesFormatStickerLOTTIE = 3
)

// Sticker Item Object
// https://discord.com/developers/docs/resources/sticker#sticker-item-object
type StickerItem struct {
	ID         Snowflake `json:"id,omitempty"`
	Name       string    `json:"name,omitempty"`
	FormatType Flag      `json:"format_type,omitempty"`
}

// Sticker Pack Object
// StickerPack represents a pack of standard stickers.
type StickerPack struct {
	ID            Snowflake `json:"id,omitempty"`
	Type          Flag      `json:"type,omitempty"`
	GuildID       Snowflake `json:"guild_id,omitempty"`
	ChannelID     Snowflake `json:"channel_id,omitempty"`
	User          *User     `json:"user,omitempty"`
	Name          string    `json:"name,omitempty"`
	Avatar        string    `json:"avatar,omitempty"`
	Token         string    `json:"token,omitempty"`
	ApplicationID Snowflake `json:"application_id,omitempty"`
	SourceGuild   *Guild    `json:"source_guild,omitempty"`
	SourceChannel *Channel  `json:"source_channel,omitempty"`
	URL           string    `json:"url,omitempty"`
}

// Webhook Object
// https://discord.com/developers/docs/resources/webhook#webhook-object

// Webhook Used to represent a webhook
// https://discord.com/developers/docs/resources/webhook#webhook-object
type Webhook struct {
	ID            Snowflake `json:"id,omitempty"`
	Type          Flag      `json:"type,omitempty"`
	GuildID       Snowflake `json:"guild_id,omitempty"`
	ChannelID     Snowflake `json:"channel_id,omitempty"`
	User          *User     `json:"user,omitempty"`
	Name          string    `json:"name,omitempty"`
	Avatar        string    `json:"avatar,omitempty"`
	Token         string    `json:"token,omitempty"`
	ApplicationID Snowflake `json:"application_id,omitempty"`
	SourceGuild   *Guild    `json:"source_guild,omitempty"`
	SourceChannel *Channel  `json:"source_channel,omitempty"`
	URL           string    `json:"url,omitempty"`
}

// Webhook Types
// https://discord.com/developers/docs/resources/webhook#webhook-object-webhook-types
const (
	FlagTypesWebhookINCOMING        = 1
	FlagTypesWebhookCHANNELFOLLOWER = 2
	FlagTypesWebhookAPPLICATION     = 3
)
