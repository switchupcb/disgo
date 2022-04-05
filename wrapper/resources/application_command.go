package resources

// Application Command Structure
// https://discord.com/developers/docs/interactions/application-commands#application-command-object-application-command-structure
type ApplicationCommand struct {
	ID                       int64                       `json:"id"`
	Type                     uint8                       `json:"type"`
	ApplicationID            int64                       `json:"application_id"`
	GuildID                  int64                       `json:"guild_id"`
	Name                     string                      `json:"name"`
	NameLocalizations        map[Flag]string             `json:"name_localizations,omitempty"`
	Description              string                      `json:"description,omitempty"`
	DescriptionLocalizations map[Flag]string             `json:"description_localizations,omitempty"`
	Options                  []*ApplicationCommandOption `json:"options"`
	DefaultPermission        bool                        `json:"default_permission,omitempty"`
	Version                  int64                       `json:"version"`
}

// Application Command Types
// https://discord.com/developers/docs/interactions/application-commands#application-command-object-application-command-types
const (
	ApplicationCommandCHAT_INPUT = 1
	ApplicationCommandUSER       = 2
	ApplicationCommandMESSAGE    = 3
)

// Application Command Option Structure
// https://discord.com/developers/docs/interactions/application-commands#application-command-object-application-command-option-structure
type ApplicationCommandOption struct {
	Type         uint8                             `json:"type"`
	Name         string                            `json:"name"`
	Description  string                            `json:"description"`
	Required     bool                              `json:"required"`
	Choices      []*ApplicationCommandOptionChoice `json:"choices"`
	Options      []*ApplicationCommandOption       `json:"options"`
	ChannelTypes []uint8                           `json:"channel_types"`
	MinValue     float64                           `json:"min_value"`
	MaxValue     float64                           `json:"max_value"`
	Autocomplete bool                              `json:"autocomplete"`
}

// Application Command Option Type
// https://discord.com/developers/docs/interactions/application-commands#application-command-object-application-command-option-type
const (
	ApplicationCommandOptionSUB_COMMAND       = 1
	ApplicationCommandOptionSUB_COMMAND_GROUP = 2
	ApplicationCommandOptionSTRING            = 3
	ApplicationCommandOptionINTEGER           = 4
	ApplicationCommandOptionBOOLEAN           = 5
	ApplicationCommandOptionUSER              = 6
	ApplicationCommandOptionCHANNEL           = 7
	ApplicationCommandOptionROLE              = 8
	ApplicationCommandOptionMENTIONABLE       = 9
	ApplicationCommandOptionNUMBER            = 10
	ApplicationCommandOptionATTACHMENT        = 11
)

// Application Command Option Choice
// https://discord.com/developers/docs/interactions/application-commands#application-command-object-application-command-option-choice-structure
type ApplicationCommandOptionChoice struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

// Application Command Interaction Data Option Structure
// https://discord.com/developers/docs/interactions/application-commands#application-command-object-application-command-interaction-data-option-structure
type ApplicationCommandInteractionDataOption struct {
	Name    string                                     `json:"name"`
	Type    uint8                                      `json:"type"`
	Value   interface{}                                `json:"value,omitempty"`
	Options []*ApplicationCommandInteractionDataOption `json:"options,omitempty"`
	Focused bool                                       `json:"focused,omitempty"`
}

// Guild Application Command Permissions Object
// https://discord.com/developers/docs/interactions/application-commands#application-command-permissions-object-guild-application-command-permissions-structure
type GuildApplicationCommandPermissions struct {
	ID            int64                            `json:"id"`
	ApplicationID int64                            `json:"application_id"`
	GuildID       int64                            `json:"guild_id"`
	Permissions   []*ApplicationCommandPermissions `json:"permissions"`
}

// Application Command Permissions Structure
// https://discord.com/developers/docs/interactions/application-commands#application-command-permissions-object-application-command-permissions-structure
type ApplicationCommandPermissions struct {
	ID         string `json:"id"`
	Type       uint8  `json:"type"`
	Permission bool   `json:"permission"`
}

// Application Command Permission Type
// https://discord.com/developers/docs/interactions/application-commands#application-command-permissions-object-application-command-permission-type
const (
	ApplicationCommandPermissionROLE = 1
	ApplicationCommandPermissionUSER = 2
)
