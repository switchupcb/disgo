package resource

// Component Types
// https://discord.com/developers/docs/interactions/message-components#component-object-component-types
const (
	ComponentActionRow  = 1
	ComponentButton     = 2
	ComponentSelectMenu = 3
	ComponentTextInput  = 4
)

// Component Object
type Component interface{}

// https://discord.com/developers/docs/interactions/message-components#component-object
type ActionsRow struct {
	Components []Component `json:"components"`
}

// Button Object
// https://discord.com/developers/docs/interactions/message-components#button-object
type Button struct {
	Style    uint8  `json:"style"`
	Label    string `json:"label"`
	Emoji    *Emoji `json:"emoji"`
	CustomID string `json:"custom_id,omitempty"`
	URL      string `json:"url,omitempty"`
	Disabled bool   `json:"disabled"`
}

// Button Styles
// https://discord.com/developers/docs/interactions/message-components#button-object-button-styles
const (
	PRIMARY   = 1
	BLURPLE   = 1
	Secondary = 2
	GREY      = 2
	Success   = 3
	GREEN     = 3
	Danger    = 4
	RED       = 4
	LINK      = 5
)

// Select Menu Structure
// https://discord.com/developers/docs/interactions/message-components#select-menu-object-select-menu-structure
type SelectMenu struct {
	CustomID    string             `json:"custom_id,omitempty"`
	Options     []SelectMenuOption `json:"options"`
	Placeholder string             `json:"placeholder"`
	MinValues   *int               `json:"min_values,omitempty"`
	MaxValues   int                `json:"max_values,omitempty"`
	Disabled    bool               `json:"disabled"`
}

// Select Menu Option Structure
// https://discord.com/developers/docs/interactions/message-components#select-menu-object-select-option-structure
type SelectMenuOption struct {
	Label       string `json:"label,omitempty"`
	Value       string `json:"value"`
	Description string `json:"description"`
	Emoji       Emoji  `json:"emoji"`
	Default     bool   `json:"default"`
}

// Text Input Structure
// https://discord.com/developers/docs/interactions/message-components#text-inputs-text-input-structure
type TextInput struct {
	CustomID    string `json:"custom_id"`
	Style       uint8  `json:"style"`
	Label       string `json:"label"`
	MinLength   int    `json:"min_length,omitempty"`
	MaxLength   int    `json:"max_length,omitempty"`
	Required    bool   `json:"required"`
	Value       string `json:"value,omitempty"`
	Placeholder string `json:"placeholder,omitempty"`
}

// TextInputStyle
// https://discord.com/developers/docs/interactions/message-components#text-inputs-text-input-styles
const (
	Short     = 1
	Paragraph = 2
)
