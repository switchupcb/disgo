# Requests

The following steps are required to add a request.

1. Define the request object in [dasgo](https://github.com/switchupcb/dasgo).

```go
// Edit Global Application Command
// PATCH /applications/{application.id}/commands/{command.id}
// https://discord.com/developers/docs/interactions/application-commands#edit-global-application-command
type EditGlobalApplicationCommand struct {
	CommandID                resources.Snowflake
	Name                     string                                `json:"name,omitempty"`
	...
}

```

2. Define the endpoint in [dasgo `endpoints.go`](https://github.com/switchupcb/dasgo/blob/main/dasgo/endpoints.go).

3. Define the return values in [`copygen/requests/setup.go`](/wrapper/copygen/requests/setup.go).
```go
// Copygen defines the functions that will be generated.
type Copygen interface {
	Send(*resources.EditGlobalApplicationCommand) (*resources.ApplicationCommand, error)
}
```

4. Generate `Endpoint` and `Send` functions using [`gen -d`](/_gen/README.md).

View the output in [`send.go`](/wrapper/send.go).