package resources

// Application Command Structure
// https://discord.com/developers/docs/interactions/application-commands#application-command-object-application-command-structure
type ApplicationCommand struct {
	ID                Snowflake                   `json:"id,omitempty"`
	Type              Flag                        `json:"type,omitempty"`
	ApplicationID     Snowflake                   `json:"application_id,omitempty"`
	GuildID           Snowflake                   `json:"guild_id,omitempty"`
	Name              string                      `json:"name,omitempty"`
	Description       string                      `json:"description,omitempty"`
	Options           []*ApplicationCommandOption `json:"options,omitempty"`
	DefaultPermission bool                        `json:"default_permission,omitempty"`
	Version           Snowflake                   `json:"version,omitempty"`
}

// Application Command Types
// https://discord.com/developers/docs/interactions/application-commands#application-command-object-application-command-types
const (
	FlagApplicationCommandTypesCHAT_INPUT = 1
	FlagApplicationCommandTypesUSER       = 2
	FlagApplicationCommandTypesMESSAGE    = 3
)

// Application Command Option Structure
// https://discord.com/developers/docs/interactions/application-commands#application-command-object-application-command-option-structure
type ApplicationCommandOption struct {
	Type         Flag                              `json:"type,omitempty"`
	Name         string                            `json:"name,omitempty"`
	Description  string                            `json:"description,omitempty"`
	Required     bool                              `json:"required,omitempty"`
	Choices      []*ApplicationCommandOptionChoice `json:"choices,omitempty"`
	Options      []*ApplicationCommandOption       `json:"options,omitempty"`
	ChannelTypes []Flag                            `json:"channel_types,omitempty"`
	MinValue     float64                           `json:"min_value,omitempty"`
	MaxValue     float64                           `json:"max_value,omitempty"`
	Autocomplete bool                              `json:"autocomplete,omitempty"`
}

// Application Command Option Type
// https://discord.com/developers/docs/interactions/application-commands#application-command-object-application-command-option-type
const (
	FlagApplicationCommandOptionTypeSUB_COMMAND       = 1
	FlagApplicationCommandOptionTypeSUB_COMMAND_GROUP = 2
	FlagApplicationCommandOptionTypeSTRING            = 3
	FlagApplicationCommandOptionTypeINTEGER           = 4
	FlagApplicationCommandOptionTypeBOOLEAN           = 5
	FlagApplicationCommandOptionTypeUSER              = 6
	FlagApplicationCommandOptionTypeCHANNEL           = 7
	FlagApplicationCommandOptionTypeROLE              = 8
	FlagApplicationCommandOptionTypeMENTIONABLE       = 9
	FlagApplicationCommandOptionTypeNUMBER            = 10
	FlagApplicationCommandOptionTypeATTACHMENT        = 11
)

// Application Command Option Choice
// https://discord.com/developers/docs/interactions/application-commands#application-command-object-application-command-option-choice-structure
type ApplicationCommandOptionChoice struct {
	Name  string      `json:"name,omitempty"`
	Value interface{} `json:"value,omitempty"`
}

// Application Command Interaction Data Option Structure
// https://discord.com/developers/docs/interactions/application-commands#application-command-object-application-command-interaction-data-option-structure
type ApplicationCommandInteractionDataOption struct {
	Name    string                                     `json:"name,omitempty"`
	Type    Flag                                       `json:"type,omitempty"`
	Value   interface{}                                `json:"value,omitempty"`
	Options []*ApplicationCommandInteractionDataOption `json:"options,omitempty"`
	Focused bool                                       `json:"focused,omitempty"`
}

// Guild Application Command Permissions Object
// https://discord.com/developers/docs/interactions/application-commands#application-command-permissions-object-guild-application-command-permissions-structure
type GuildApplicationCommandPermissions struct {
	ID            Snowflake                        `json:"id,omitempty"`
	ApplicationID Snowflake                        `json:"application_id,omitempty"`
	GuildID       Snowflake                        `json:"guild_id,omitempty"`
	Permissions   []*ApplicationCommandPermissions `json:"permissions,omitempty"`
}

// Application Command Permissions Structure
// https://discord.com/developers/docs/interactions/application-commands#application-command-permissions-object-application-command-permissions-structure
type ApplicationCommandPermissions struct {
	ID         Snowflake `json:"id,omitempty"`
	Type       Flag      `json:"type,omitempty"`
	Permission bool      `json:"permission,omitempty"`
}

// Application Command Permission Type
// https://discord.com/developers/docs/interactions/application-commands#application-command-permissions-object-application-command-permission-type
const (
	FlagApplicationCommandPermissionTypeROLE = 1
	FlagApplicationCommandPermissionTypeUSER = 2
)
