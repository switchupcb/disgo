package resources

// Application Command Structure
// https://discord.com/developers/docs/interactions/application-commands#application-command-object-application-command-structure
type ApplicationCommand struct {
	ID                       Snowflake                   `json:"id,omitempty"`
	Type                     Flag                        `json:"type,omitempty"`
	ApplicationID            Snowflake                   `json:"application_id,omitempty"`
	GuildID                  Snowflake                   `json:"guild_id,omitempty"`
	Name                     string                      `json:"name,omitempty"`
	NameLocalizations        map[Flag]string             `json:"name_localizations,omitempty"`
	Description              string                      `json:"description,omitempty"`
	DescriptionLocalizations map[Flag]string             `json:"description_localizations,omitempty"`
	Options                  []*ApplicationCommandOption `json:"options,omitempty"`
	DefaultPermission        bool                        `json:"default_permission,omitempty"`
	Version                  Snowflake                   `json:"version,omitempty"`
}

// Application Command Types
// https://discord.com/developers/docs/interactions/application-commands#application-command-object-application-command-types
const (
	FlagTypesCommandApplicationCHAT_INPUT = 1
	FlagTypesCommandApplicationUSER       = 2
	FlagTypesCommandApplicationMESSAGE    = 3
)

// Application Command Option Structure
// https://discord.com/developers/docs/interactions/application-commands#application-command-object-application-command-option-structure
type ApplicationCommandOption struct {
	Type                     Flag                              `json:"type,omitempty"`
	Name                     string                            `json:"name,omitempty"`
	NameLocalizations        map[Flag]string                   `json:"name_localizations,omitempty"`
	Description              string                            `json:"description,omitempty"`
	DescriptionLocalizations map[Flag]string                   `json:"description_localizations,omitempty"`
	Required                 bool                              `json:"required,omitempty"`
	Choices                  []*ApplicationCommandOptionChoice `json:"choices,omitempty"`
	Options                  []*ApplicationCommandOption       `json:"options,omitempty"`
	ChannelTypes             []*Flag                           `json:"channel_types,omitempty"`
	MinValue                 float64                           `json:"min_value,omitempty"`
	MaxValue                 float64                           `json:"max_value,omitempty"`
	Autocomplete             bool                              `json:"autocomplete,omitempty"`
}

// Application Command Option Type
// https://discord.com/developers/docs/interactions/application-commands#application-command-object-application-command-option-type
const (
	FlagTypeOptionCommandApplicationSUB_COMMAND       = 1
	FlagTypeOptionCommandApplicationSUB_COMMAND_GROUP = 2
	FlagTypeOptionCommandApplicationSTRING            = 3
	FlagTypeOptionCommandApplicationINTEGER           = 4
	FlagTypeOptionCommandApplicationBOOLEAN           = 5
	FlagTypeOptionCommandApplicationUSER              = 6
	FlagTypeOptionCommandApplicationCHANNEL           = 7
	FlagTypeOptionCommandApplicationROLE              = 8
	FlagTypeOptionCommandApplicationMENTIONABLE       = 9
	FlagTypeOptionCommandApplicationNUMBER            = 10
	FlagTypeOptionCommandApplicationATTACHMENT        = 11
)

// Application Command Option Choice
// https://discord.com/developers/docs/interactions/application-commands#application-command-object-application-command-option-choice-structure
type ApplicationCommandOptionChoice struct {
	Name              string          `json:"name,omitempty"`
	NameLocalizations map[Flag]string `json:"name_localizations,omitempty"`
	Value             interface{}     `json:"value,omitempty"`
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
	FlagTypePermissionCommandApplicationROLE = 1
	FlagTypePermissionCommandApplicationUSER = 2
)
