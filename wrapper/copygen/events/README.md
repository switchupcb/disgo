# Events

The following steps are required to add an event.

1. Define the event object in [dasgo](https://github.com/switchupcb/dasgo).

```go
// Hello Structure
// https://discord.com/developers/docs/topics/gateway#hello-hello-structure
type Hello struct {
	HeartbeatInterval time.Duration `json:"heartbeat_interval,omitempty"`
}
```

2. Define the event name in [dasgo `events.go`](https://github.com/switchupcb/dasgo/blob/main/dasgo/events.go).

3. Define the return values in [`copygen/events/setup.go`](/wrapper/copygen/events/setup.go).
```go
// Copygen defines the functions that will be generated.
type Copygen interface {
    Hello(*disgo.Hello)
}
```

1. Generate the `Handlers` struct,  `Handle` and `handle` functions using [`gen -d`](/_gen/README.md).

View the output in [`handle.go`](/wrapper/handle.go).