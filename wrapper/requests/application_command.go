package requests

import "github.com/switchupcb/disgo/wrapper/resources"

// Get Global Application Commands
// GET/applications/{application.id}/commands
// https://discord.com/developers/docs/interactions/application-commands#get-global-application-commands
type GetGlobalApplicationCommands struct {
	WithLocalizations bool `json:"with_localizations,omitempty"`
}

// Create Global Application Command
// POST/applications/{application.id}/commands
// https://discord.com/developers/docs/interactions/application-commands#create-global-application-command
type CreateGlobalApplicationCommand struct {
	Name                     string                                `json:"name,omitempty"`
	NameLocalizations        map[resources.Flag]string             `json:"name_localizations,omitempty"`
	Description              string                                `json:"description,omitempty"`
	DescriptionLocalizations map[resources.Flag]string             `json:"description_localizations,omitempty"`
	Options                  []*resources.ApplicationCommandOption `json:"options,omitempty"`
	DefaultPermission        bool                                  `json:"default_permission,omitempty"`
	Type                     resources.Flag                        `json:"type,omitempty"`
}

// Get Global Application Command
// GET/applications/{application.id}/commands/{command.id}
// https://discord.com/developers/docs/interactions/application-commands#get-global-application-command
type GetGlobalApplicationCommand struct {
	CommandID resources.Snowflake
}

// Edit Global Application Command
// PATCH/applications/{application.id}/commands/{command.id}
// https://discord.com/developers/docs/interactions/application-commands#edit-global-application-command
type EditGlobalApplicationCommand struct {
	CommandID                resources.Snowflake
	Name                     string                                `json:"name,omitempty"`
	NameLocalizations        map[resources.Flag]string             `json:"name_localizations,omitempty"`
	Description              string                                `json:"description,omitempty"`
	DescriptionLocalizations map[resources.Flag]string             `json:"description_localizations,omitempty"`
	Options                  []*resources.ApplicationCommandOption `json:"options,omitempty"`
	DefaultPermission        bool                                  `json:"default_permission,omitempty"`
}
