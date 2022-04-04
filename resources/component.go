package resources

// Component Types
// https://discord.com/developers/docs/interactions/message-components#component-object-component-types
const (
	FlagComponentTypesActionRow  = 1
	FlagComponentTypesButton     = 2
	FlagComponentTypesSelectMenu = 3
	FlagComponentTypesTextInput  = 4
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
	Style    Flag   `json:"style,omitempty"`
	Label    string `json:"label,omitempty"`
	Emoji    *Emoji `json:"emoji,omitempty"`
	CustomID string `json:"custom_id,omitempty"`
	URL      string `json:"url,omitempty"`
	Disabled bool   `json:"disabled,omitempty"`
}

// Button Styles
// https://discord.com/developers/docs/interactions/message-components#button-object-button-styles
const (
	FlagButtonStylesPRIMARY   = 1
	FlagButtonStylesBLURPLE   = 1
	FlagButtonStylesSecondary = 2
	FlagButtonStylesGREY      = 2
	FlagButtonStylesSuccess   = 3
	FlagButtonStylesGREEN     = 3
	FlagButtonStylesDanger    = 4
	FlagButtonStylesRED       = 4
	FlagButtonStylesLINK      = 5
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
	Label       string `json:"label,omitempty"`
	Value       string `json:"value,omitempty"`
	Description string `json:"description,omitempty"`
	Emoji       Emoji  `json:"emoji,omitempty"`
	Default     bool   `json:"default,omitempty"`
}

// Text Input Structure
// https://discord.com/developers/docs/interactions/message-components#text-inputs-text-input-structure
type TextInput struct {
	CustomID    string   `json:"custom_id,omitempty"`
	Style       Flag     `json:"style,omitempty"`
	Label       string   `json:"label,omitempty"`
	MinLength   CodeFlag `json:"min_length,omitempty"`
	MaxLength   CodeFlag `json:"max_length,omitempty"`
	Required    bool     `json:"required,omitempty"`
	Value       string   `json:"value,omitempty"`
	Placeholder string   `json:"placeholder,omitempty"`
}

// TextInputStyle
// https://discord.com/developers/docs/interactions/message-components#text-inputs-text-input-styles
const (
	FlagTextInputStyleShort     = 1
	FlagTextInputStyleParagraph = 2
)
