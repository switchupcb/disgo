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
	ApplicationID resources.Snowflake
	CommandID     resources.Snowflake
}

// Edit Global Application Command
// PATCH/applications/{application.id}/commands/{command.id}
// https://discord.com/developers/docs/interactions/application-commands#edit-global-application-command
type EditGlobalApplicationCommand struct {
	CommandID                resources.Snowflake
	ApplicationID            resources.Snowflake
	Name                     string                                `json:"name,omitempty"`
	NameLocalizations        map[resources.Flag]string             `json:"name_localizations,omitempty"`
	Description              string                                `json:"description,omitempty"`
	DescriptionLocalizations map[resources.Flag]string             `json:"description_localizations,omitempty"`
	Options                  []*resources.ApplicationCommandOption `json:"options,omitempty"`
	DefaultPermission        bool                                  `json:"default_permission,omitempty"`
}

// Delete Global Application Command
// DELETE /applications/{application.id}/commands/{command.id}
// https://discord.com/developers/docs/interactions/application-commands#delete-global-application-command
type DeleteGlobalApplicationCommand struct {
	CommandID     resources.Snowflake
	ApplicationID resources.Snowflake
	GuildID       resources.Snowflake
}

// Bulk Overwrite Global Application Commands
// PUT /applications/{application.id}/commands
// https://discord.com/developers/docs/interactions/application-commands#bulk-overwrite-global-application-commands
type BulkOverwriteGlobalApplicationCommands struct {
	ApplicationCommands []*resources.ApplicationCommand
}

// Get Guild Application Commands
// GET /applications/{application.id}/guilds/{guild.id}/commands
// https://discord.com/developers/docs/interactions/application-commands#get-guild-application-commands
type GetGuildApplicationCommands struct {
	GuildID           resources.Snowflake
	ApplicationID     resources.Snowflake
	WithLocalizations bool `json:"with_localizations,omitempty"`
}

// Create Guild Application Command
// POST /applications/{application.id}/guilds/{guild.id}/commands
// https://discord.com/developers/docs/interactions/application-commands#create-guild-application-command
type CreateGuildApplicationCommand struct {
	GuildID                  resources.Snowflake
	Name                     string                                `json:"name,omitempty"`
	NameLocalizations        map[resources.Flag]string             `json:"name_localizations,omitempty"`
	Description              string                                `json:"description,omitempty"`
	DescriptionLocalizations map[resources.Flag]string             `json:"description_localizations,omitempty"`
	Options                  []*resources.ApplicationCommandOption `json:"options,omitempty"`
	DefaultPermission        bool                                  `json:"default_permission,omitempty"`
	Type                     resources.Flag                        `json:"type,omitempty"`
}

// Get Guild Application Command
// GET /applications/{application.id}/guilds/{guild.id}/commands/{command.id}
// https://discord.com/developers/docs/interactions/application-commands#get-guild-application-command
type GetGuildApplicationCommand struct {
	ApplicationID resources.Snowflake
	GuildID       resources.Snowflake
	CommandID     resources.Snowflake
}

// Edit Guild Application Command
// PATCH /applications/{application.id}/guilds/{guild.id}/commands/{command.id}
// https://discord.com/developers/docs/interactions/application-commands#edit-guild-application-command
type EditGuildApplicationCommand struct {
	ApplicationID            resources.Snowflake
	GuildID                  resources.Snowflake
	CommandID                resources.Snowflake
	Name                     string                                `json:"name,omitempty"`
	NameLocalizations        map[resources.Flag]string             `json:"name_localizations,omitempty"`
	Description              string                                `json:"description,omitempty"`
	DescriptionLocalizations map[resources.Flag]string             `json:"description_localizations,omitempty"`
	Options                  []*resources.ApplicationCommandOption `json:"options,omitempty"`
	DefaultPermission        bool                                  `json:"default_permission,omitempty"`
}

// Delete Guild Application Command
// DELETE /applications/{application.id}/guilds/{guild.id}/commands/{command.id}
// https://discord.com/developers/docs/interactions/application-commands#delete-guild-application-command
type DeleteGuildApplicationCommand struct {
	GuildID resources.Snowflake
}

// Bulk Overwrite Guild Application Commands
// PUT /applications/{application.id}/guilds/{guild.id}/commands
// https://discord.com/developers/docs/interactions/application-commands#bulk-overwrite-guild-application-commands
type BulkOverwriteGuildApplicationCommands struct {
	ApplicationID            resources.Snowflake
	GuildID                  resources.Snowflake
	Name                     string                                `json:"name,omitempty"`
	NameLocalizations        map[resources.Flag]string             `json:"name_localizations,omitempty"`
	Description              string                                `json:"description,omitempty"`
	DescriptionLocalizations map[resources.Flag]string             `json:"description_localizations,omitempty"`
	Options                  []*resources.ApplicationCommandOption `json:"options,omitempty"`
	DefaultPermission        bool                                  `json:"default_permission,omitempty"`
	Type                     resources.Flag                        `json:"type,omitempty"`
}

// Get Guild Application Command Permissions
// GET /applications/{application.id}/guilds/{guild.id}/commands/permissions
// https://discord.com/developers/docs/interactions/application-commands#get-guild-application-command-permissions
type GetGuildApplicationCommandPermissions struct {
	GuildID resources.Snowflake
}

// Get Application Command Permissions
// GET /applications/{application.id}/guilds/{guild.id}/commands/{command.id}/permissions
// https://discord.com/developers/docs/interactions/application-commands#get-application-command-permissions
type GetApplicationCommandPermissions struct {
	ApplicationID resources.Snowflake
	GuildID       resources.Snowflake
	CommandID     resources.Snowflake
}

// Edit Application Command Permissions
// PUT /applications/{application.id}/guilds/{guild.id}/commands/{command.id}/permissions
// https://discord.com/developers/docs/interactions/application-commands#edit-application-command-permissions
type EditApplicationCommandPermissions struct {
	ApplicationID resources.Snowflake
	GuildID       resources.Snowflake
	CommandID     resources.Snowflake
	Permissions   []*resources.ApplicationCommandPermissions `json:"permissions,omitempty"`
}

// Batch Edit Application Command Permissions
// PUT /applications/{application.id}/guilds/{guild.id}/commands/permissions
// https://discord.com/developers/docs/interactions/application-commands#batch-edit-application-command-permissions
type BatchEditApplicationCommandPermissions struct {
	GuildID resources.Snowflake
}
