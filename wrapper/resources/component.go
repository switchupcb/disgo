package resources

// Component Types
// https://discord.com/developers/docs/interactions/message-components#component-object-component-types
const (
	FlagTypesComponentActionRow  = 1
	FlagTypesComponentButton     = 2
	FlagTypesComponentSelectMenu = 3
	FlagTypesComponentTextInput  = 4
)

// Component Object
type Component interface{}

// https://discord.com/developers/docs/interactions/message-components#component-object
type ActionsRow struct {
	Components []Component `json:"components,omitempty"`
}

// Button Object
// https://discord.com/developers/docs/interactions/message-components#button-object
type Button struct {
	Style    Flag    `json:"style,omitempty"`
	Label    *string `json:"label,omitempty"`
	Emoji    *Emoji  `json:"emoji,omitempty"`
	CustomID string  `json:"custom_id,omitempty"`
	URL      string  `json:"url,omitempty"`
	Disabled bool    `json:"disabled,omitempty"`
}

// Button Styles
// https://discord.com/developers/docs/interactions/message-components#button-object-button-styles
const (
	FlagStylesbuttonPRIMARY   = 1
	FlagStylesbuttonBLURPLE   = 1
	FlagStylesbuttonSecondary = 2
	FlagStylesbuttonGREY      = 2
	FlagStylesbuttonSuccess   = 3
	FlagStylesbuttonGREEN     = 3
	FlagStylesbuttonDanger    = 4
	FlagStylesbuttonRED       = 4
	FlagStylesbuttonLINK      = 5
)

// Select Menu Structure
// https://discord.com/developers/docs/interactions/message-components#select-menu-object-select-menu-structure
type SelectMenu struct {
	CustomID    string             `json:"custom_id,omitempty"`
	Options     []SelectMenuOption `json:"options,omitempty"`
	Placeholder string             `json:"placeholder,omitempty"`
	MinValues   *Flag              `json:"min_values,omitempty"`
	MaxValues   Flag               `json:"max_values,omitempty"`
	Disabled    bool               `json:"disabled,omitempty"`
}

// Select Menu Option Structure
// https://discord.com/developers/docs/interactions/message-components#select-menu-object-select-option-structure
type SelectMenuOption struct {
	Label       *string `json:"label,omitempty"`
	Value       *string `json:"value,omitempty"`
	Description *string `json:"description,omitempty"`
	Emoji       Emoji   `json:"emoji,omitempty"`
	Default     bool    `json:"default,omitempty"`
}

// Text Input Structure
// https://discord.com/developers/docs/interactions/message-components#text-inputs-text-input-structure
type TextInput struct {
	CustomID    string    `json:"custom_id,omitempty"`
	Style       Flag      `json:"style,omitempty"`
	Label       *string   `json:"label,omitempty"`
	MinLength   *CodeFlag `json:"min_length,omitempty"`
	MaxLength   CodeFlag  `json:"max_length,omitempty"`
	Required    bool      `json:"required,omitempty"`
	Value       string    `json:"value,omitempty"`
	Placeholder *string   `json:"placeholder,omitempty"`
}

// TextInputStyle
// https://discord.com/developers/docs/interactions/message-components#text-inputs-text-input-styles
const (
	FlagStyleInputTextShort     = 1
	FlagStyleInputTextParagraph = 2
)
