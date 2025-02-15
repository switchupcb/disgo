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

3. Define the return values in [`copygen/requests/setup.go`](/_gen/tools/copygen/requests/setup.go).
```go
// Copygen defines the functions that will be generated.
type Copygen interface {
	Send(*resources.EditGlobalApplicationCommand) (*resources.ApplicationCommand, error)
}
```

Use [`gen -d`](/_gen/README.md) once to perform the following actions.

1. Generate `Endpoint` functions, `Send` functions, and `RouteIDs`. View the output in [`request_send.go`](/wrapper/request_send.go).

2. Generate the rate limit algorithm map `RateLimitHashFuncs`. View the output in [`ratelimit_algorithm_map.go`](/wrapper/ratelimit_algorithm_map.go).

3. Generate the full endpoints map in the [coverage test order generator](/_gen/coverage).View the output in [`coverage_endpoint_graph.go`](/_gen/coverage/coverage_endpoint_map.go).


