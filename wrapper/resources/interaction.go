package resource

// Interaction Object
// https://discord.com/developers/docs/interactions/receiving-and-responding#interaction-object-interaction-structure
type Interaction struct {
	ID            string          `json:"id"`
	ApplicationID int64           `json:"application_id"`
	Type          uint8           `json:"type"`
	Data          InteractionData `json:"data"`
	GuildID       int64           `json:"guild_id"`
	ChannelID     int64           `json:"channel_id"`
	Member        *GuildMember    `json:"member"`
	User          *User           `json:"user"`
	Token         string          `json:"token"`
	Version       int             `json:"version"`
	Message       *Message        `json:"message"`
	Locale        string          `json:"locale"`
	GuildLocale   string          `json:"guild_locale"`
}

// Interaction Type
// https://discord.com/developers/docs/interactions/receiving-and-responding#interaction-object-interaction-type
const (
	PING                             = 1
	APPLICATION_COMMAND              = 2
	MESSAGE_COMPONENT                = 3
	APPLICATION_COMMAND_AUTOCOMPLETE = 4
	MODAL_SUBMIT                     = 5
)

// Interaction Data Structure
// https://discord.com/developers/docs/interactions/receiving-and-responding#interaction-object-interaction-data-structure
type InteractionData struct {
	ID            int64                                      `json:"id"`
	Name          string                                     `json:"name"`
	Type          uint8                                      `json:"type"`
	Resolved      *ResolvedData                              `json:"resolved"`
	Options       []*ApplicationCommandInteractionDataOption `json:"options"`
	CustomID      string                                     `json:"custom_id"`
	ComponentType uint8                                      `json:"component_type"`
	Values        []string                                   `json:"values"`
	TargetID      int64                                      `json:"target_id"`
	Components    []Component                                `json:"components"`
}

// Resolved Data Structure
// https://discord.com/developers/docs/interactions/receiving-and-responding#interaction-object-resolved-data-structure
type ResolvedData struct {
	Users       map[int64]*User        `json:"users"`
	Members     map[int64]*GuildMember `json:"members"`
	Roles       map[int64]*Role        `json:"roles"`
	Channels    map[int64]*Channel     `json:"channels"`
	Messages    map[int64]*Message     `json:"messages"`
	Attachments map[int64]*Attachment  `json:"attachments"`
}

// Message Interaction Structure
// https://discord.com/developers/docs/interactions/receiving-and-responding#message-interaction-object-message-interaction-structure
type MessageInteraction struct {
	ID     int64        `json:"id"`
	Type   uint8        `json:"type"`
	Name   string       `json:"name"`
	User   *User        `json:"user"`
	Member *GuildMember `json:"member"`
}

// Interaction Response Structure
// https://discord.com/developers/docs/interactions/receiving-and-responding#interaction-response-object-interaction-response-structure
type InteractionResponse struct {
	Type uint8                    `json:"type,omitempty"`
	Data *InteractionCallbackData `json:"data,omitempty"`
}

// Interaction Callback Type
// https://discord.com/developers/docs/interactions/receiving-and-responding#interaction-response-object-interaction-callback-type
const (
	PONG                                    = 1
	CHANNEL_MESSAGE_WITH_SOURCE             = 4
	DEFERRED_CHANNEL_MESSAGE_WITH_SOURCE    = 5
	DEFERRED_UPDATE_MESSAGE                 = 6
	UPDATE_MESSAGE                          = 7
	APPLICATION_COMMAND_AUTOCOMPLETE_RESULT = 8
	MODAL                                   = 9
)

// Interaction Callback Data Structure
// https://discord.com/developers/docs/interactions/receiving-and-responding#interaction-response-object-interaction-callback-data-structure
type InteractionCallbackData interface{}

// Messages
// https://discord.com/developers/docs/interactions/receiving-and-responding#interaction-response-object-messages
type Messages struct {
	TTS             bool             `json:"tts"`
	Content         string           `json:"content"`
	Embeds          []*Embed         `json:"embeds"`
	Components      []Component      `json:"components"`
	AllowedMentions *AllowedMentions `json:"allowed_mentions,omitempty"`
	Flags           uint64           `json:"flags,omitempty"`
	Attachments     []*Attachment    `json:"attachments,omitempty"`
}

// Autocomplete
// https://discord.com/developers/docs/interactions/receiving-and-responding#interaction-response-object-autocomplete
type Autocomplete struct {
	Choices []ApplicationCommandOptionChoice `json:"choices"`
}

// Modal
// https://discord.com/developers/docs/interactions/receiving-and-responding#interaction-response-object-modal
type ModalSubmitInteractionData struct {
	CustomID   string      `json:"custom_id"`
	Title      string      `json:"title"`
	Components []Component `json:"components"`
}
