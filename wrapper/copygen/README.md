# Requests

The following steps are required to add a request.

1. Define the request object in [dasgo](https://github.com/switchupcb/dasgo).

```go
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

```

2. Define the return values in `copygen/setup.go`.
```go
// Copygen defines the functions that will be generated.
type Copygen interface {
	Send(*resources.EditGlobalApplicationCommand) (*resources.ApplicationCommand, error)
}
```

3. Generate `Send` functions using a [copygen](/wrapper/requests/copygen/template/generate.go) template.

```
copygen -yml wrapper/requests/copygen/setup/setup.yml
```
_Current Working Directory: `/disgo`_

View the output in [`send.go`](send.go).

## Endpoints

Endpoints are automatically generated using a [generator](../../gen/README.md).

