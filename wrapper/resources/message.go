package resource

import "time"

// Message Object
// https://discord.com/developers/docs/resources/channel#message-object
type Message struct {
	ID                int64             `json:"id"`
	ChannelID         int64             `json:"channel_id"`
	GuildID           int64             `json:"guild_id"`
	Author            *User             `json:"author"`
	Member            *GuildMember      `json:"member"`
	Content           string            `json:"content"`
	Timestamp         time.Time         `json:"timestamp"`
	EditedTimestamp   time.Time         `json:"edited_timestamp"`
	TTS               bool              `json:"tts"`
	MentionEveryone   bool              `json:"mention_everyone"`
	Mentions          []*User           `json:"mentions"`
	MentionRoles      []int64           `json:"mention_roles"`
	MentionChannels   []*ChannelMention `json:"mention_channels"`
	Attachments       []*Attachment     `json:"attachments"`
	Embeds            []*Embed          `json:"embeds"`
	Reactions         []*Reaction       `json:"reactions"`
	Nonce             interface{}       `json:"nonce"`
	Pinned            bool              `json:"pinned"`
	WebhookID         int64             `json:"webhook_id"`
	Type              uint8             `json:"type"`
	Activity          MessageActivity   `json:"activity"`
	Application       *Application      `json:"application"`
	MessageReference  *MessageReference `json:"message_reference"`
	Flags             uint8             `json:"flags"`
	ReferencedMessage *Message          `json:"referenced_message"`
	Interaction       *Interaction      `json:"interaction"`
	Thread            *Channel          `json:"thread"`
	Components        []*Component      `json:"components"`
	StickerItems      []*StickerItem    `json:"sticker_items"`
}

// Message Types
// https://discord.com/developers/docs/resources/channel#message-object-message-types
const (
	DEFAULT                                      = 0
	RECIPIENT_ADD                                = 1
	RECIPIENT_REMOVE                             = 2
	CALL                                         = 3
	CHANNEL_NAME_CHANGE                          = 4
	CHANNEL_ICON_CHANGE                          = 5
	CHANNEL_PINNED_MESSAGE                       = 6
	GUILD_MEMBER_JOIN                            = 7
	USER_PREMIUM_GUILD_SUBSCRIPTION              = 8
	USER_PREMIUM_GUILD_SUBSCRIPTION_TIER_ONE     = 9
	USER_PREMIUM_GUILD_SUBSCRIPTION_TIER_TWO     = 10
	USER_PREMIUM_GUILD_SUBSCRIPTION_TIER_THREE   = 11
	CHANNEL_FOLLOW_ADD                           = 12
	GUILD_DISCOVERY_DISQUALIFIED                 = 14
	GUILD_DISCOVERY_REQUALIFIED                  = 15
	GUILD_DISCOVERY_GRACE_PERIOD_INITIAL_WARNING = 16
	GUILD_DISCOVERY_GRACE_PERIOD_FINAL_WARNING   = 17
	THREAD_CREATED                               = 18
	REPLY                                        = 19
	CHAT_INPUT_COMMAND                           = 20
	THREAD_STARTER_MESSAGE                       = 21
	GUILD_INVITE_REMINDER                        = 22
	CONTEXT_MENU_COMMAND                         = 23
)

// Message Activity Structure
// https://discord.com/developers/docs/resources/channel#message-object-message-activity-structure
type MessageActivity struct {
	Type    int    `json:"type"`
	PartyID string `json:"party_id"`
}

// Message Activity Types
// https://discord.com/developers/docs/resources/channel#message-object-message-activity-types
const (
	JOIN         = 1
	SPECTATE     = 2
	LISTEN       = 3
	JOIN_REQUEST = 5
)

// Message Flags
// https://discord.com/developers/docs/resources/channel#message-object-message-flags
const (
	CROSSPOSTED                            = 1 << 0
	IS_CROSSPOST                           = 1 << 1
	SUPPRESS_EMBEDS                        = 1 << 2
	SOURCE_MESSAGE_DELETED                 = 1 << 3
	URGENT                                 = 1 << 4
	HAS_THREAD                             = 1 << 5
	EPHEMERAL                              = 1 << 6
	LOADING                                = 1 << 7
	FAILED_TO_MENTION_SOME_ROLES_IN_THREAD = 1 << 8
)

// Message Reference Object
// https://discord.com/developers/docs/resources/channel#message-reference-object
type MessageReference struct {
	MessageID       int64 `json:"message_id"`
	ChannelID       int64 `json:"channel_id"`
	GuildID         int64 `json:"guild_id"`
	FailIfNotExists bool  `json:"fail_if_not_exists"`
}

// Message Attachment Object
// https://discord.com/developers/docs/resources/channel#attachment-object
type Attachment struct {
	ID       int64  `json:"id"`
	Filename string `json:"filename"`
	Size     uint   `json:"size"`
	URL      string `json:"url"`
	ProxyURL string `json:"proxy_url"`
	Height   uint   `json:"height"`
	Width    uint   `json:"width"`

	SpoilerTag bool `json:"-"`
}

// Sticker Structure
// https://discord.com/developers/docs/resources/sticker#sticker-object-sticker-structure
type Sticker struct {
	ID          string `json:"id"`
	PackID      string `json:"pack_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Tags        string `json:"tags"`
	Asset       string `json:"asset"`
	Type        uint8  `json:"type"`
	FormatType  uint8  `json:"format_type"`
	Available   bool   `json:"available"`
	GuildID     int64  `json:"guild_id"`
	User        *User  `json:"user"`
	SortValue   int    `json:"sort_value"`
}

// Sticker Types
// https://discord.com/developers/docs/resources/sticker#sticker-object-sticker-types
const (
	STANDARD = 1
	GUILD    = 2
)

// Sticker Format Types
// https://discord.com/developers/docs/resources/sticker#sticker-object-sticker-format-types
const (
	PNG    = 1
	APNG   = 2
	LOTTIE = 3
)

// Sticker Item Object
// https://discord.com/developers/docs/resources/sticker#sticker-item-object
type StickerItem struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	FormatType uint8  `json:"format_type"`
}

// Sticker Pack Object
// StickerPack represents a pack of standard stickers.
type StickerPack struct {
	ID             string     `json:"id"`
	Stickers       []*Sticker `json:"stickers"`
	Name           string     `json:"name"`
	SKUID          string     `json:"sku_id"`
	CoverStickerID string     `json:"cover_sticker_id"`
	Description    string     `json:"description"`
	BannerAssetID  int64      `json:"banner_asset_id"`
}

// Webhook Object
// https://discord.com/developers/docs/resources/webhook#webhook-object

// Webhook Used to represent a webhook
// https://discord.com/developers/docs/resources/webhook#webhook-object
type Webhook struct {
	ID            string   `json:"id"`
	Type          uint8    `json:"type"`
	GuildID       string   `json:"guild_id"`
	ChannelID     string   `json:"channel_id"`
	User          *User    `json:"user"`
	Name          string   `json:"name"`
	Avatar        string   `json:"avatar"`
	Token         string   `json:"token"`
	ApplicationID string   `json:"application_id,omitempty"`
	SourceGuild   *Guild   `json:"source_guild,omitempty"`
	SourceChannel *Channel `json:"source_channel,omitempty"`
	URL           string   `json:"url"`
}

// Webhook Types
// https://discord.com/developers/docs/resources/webhook#webhook-object-webhook-types
const (
	WebhookINCOMING        = 1
	WebhookCHANNELFOLLOWER = 2
	WebhookAPPLICATION     = 3
)
