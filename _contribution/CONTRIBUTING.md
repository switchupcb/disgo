# Contributing

## License

You agree to license any contribution to this library under the [Apache License 2.0](#license).

## Pull Requests

Pull requests must follow the [code specification](#code-specification) and work with all [test cases](#test).

## Domain

The domain of Disgo in providing an API for HTTP/WebSocket requests. The program uses provided structures (from the [Discord API](https://discord.com/developers/docs/reference)) to provide simple abstractions for end users _(developers)_.

## Project Structure

The repository consists of a detailed [README](/README.md), [examples](/_examples/), [**API Wrapper**](/wrapper), [**Cache**](/cache), and [**Shard Manager**](/shard/).

### Disgo

The API Wrapper consists of multiple packages. A [bundler](https://pkg.go.dev/golang.org/x/tools/cmd/bundle) is used to package the API into a single-usable package _(`disgo.go`)_.

| Package   | Description                                   |
| :-------- | :-------------------------------------------- |
| wrapper   | Contains the wrapper bundling functionality.  |
| client    | The Disgo Client Abstraction.                 |
| events    | Discord API Events.                           |
| pkg       | Utility functionality for specific libraries. |
| requests  | Discord API Requests.                         |
| resources | Discord API Resources.                        |
| sessions  | Discord API WebSocket Sessions (Gateways).    |

#### Structs

Structs are sourced from [Dasgo](https://github.com/switchupcb/dasgo) and refactored into the correct name scheme for end users. The abstraction _(i.e Resource)_ is pre-pended to the resource name _(i.e User)_ to speed up development. Modern IDE's will show the developer a list of resources when `disgo.R` is typed; rather than a bunch of irrelevant resources, functions, and variables.

```go
disgo.ResourceUser
disgo.RequestGetUser
```

#### Requests

Resource GET, DELETE, POST, PUT, BULK _(GET, ...)_ `Send()` functions are generated from the respective requests object. For more information, read the requests [README](/wrapper/requests/README.md).

### Code Specification

#### Code Generation

Disgo uses generators to easily update and maintain over 8000 lines of code. For more information, read [gen](/_gen/README.md).

#### Comments

Comments follow [Effective Go](https://golang.org/doc/effective_go#commentary) and explain why more than what _(unless the "what" isn't intuitive)_.

#### Static Code Analysis

Disgo uses [golangci-lint](https://github.com/golangci/golangci-lint) in order to statically analyze code. You can install golangci-lint with `go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.46.0`. If you receive a `diff` error _(while running)_, you must add a `diff` tool in your PATH. There is one located in the `Git` bin.

If you receive `File is not ... with -...`, use `golangci-lint run --disable-all --no-config -Egofmt --fix`.

 | Directory | Command                       | Description                                      |
 | :-------- | :---------------------------- | :----------------------------------------------- |
 | `disgo`   | `golangci-lint run ./wrapper` | Perform static code analysis on the API Wrapper. |
 | `_gen`    | `golangci-lint run`           | Perform static code analysis on the generator.   |

#### Fieldalignment

Disgo [fieldaligns](https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/fieldalignment) bundled code to save memory.

### Test

#### Unit Tests

Unit tests are used to test logic.

#### Integration Tests

Integration tests are used to ensure functionality between the API and Discord.

# Roadmap

Disgo is currently a PROOF OF CONCEPT. Here are the steps required in order to complete it:

1. Implement Client (OAuth).
2. Implement Gateway (WebSocket, Events).
3. Implement Rate Limits (HTTP, Gateway)
4. Implement Retries.
5. **Implement Testing** _[usable at this stage]_.
6. Implement Sharding.
7. Implement Cache (which is likely where most effort lies; caching is difficult).

In addition, we must make [decisions](/_contribution/libraries/) for the following:
1. UDP connections (Voice)
2. [Audio Processing using Opus](https://discord.com/developers/docs/topics/voice-connections#encrypting-and-sending-voice)
