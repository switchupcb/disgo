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

| Package | Description       |
| :------ | :---------------- |
| wrapper | API Wrapper.      |
| cache   | Cache.            |
| shard   | Sharding Manager. |

_A [bundler](https://pkg.go.dev/golang.org/x/tools/cmd/bundle) is used to package the API into a `disgo` package (`disgo.go`)_.

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

Disgo uses [golangci-lint](https://github.com/golangci/golangci-lint) in order to statically analyze code. You can install golangci-lint with `go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.46.2`. If you receive a `diff` error _(while running)_, you must add a `diff` tool in your PATH. There is one located in the `Git` bin.

If you receive `File is not ... with -...`, use `golangci-lint run --disable-all --no-config -Egofmt --fix`.

 | Directory | Command                       | Description                                        |
 | :-------- | :---------------------------- | :------------------------------------------------- |
 | `disgo`   | `golangci-lint run ./wrapper` | Perform static code analysis on the API Wrapper.   |
 | `_gen`    | `golangci-lint run ./_gen`    | Perform static code analysis on the generator.     |
 | `cache`   | `golangci-lint run ./cache`   | Perform static code analysis on the Disgo Cache.   |
 | `shard`   | `golangci-lint run ./shard`   | Perform static code analysis on the Shard Manager. |

#### Fieldalignment

Disgo [fieldaligns](https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/fieldalignment) bundled code to save memory.

### Test

#### Unit Tests

Unit tests are used to test logic.

#### Integration Tests

Integration tests are used to ensure functionality between the API and Discord.

#### Running Tests

Use `go test` to run the tests in the current directory. Use `go test ./<dir>` to run tests in a given directory (from the current directory). Use [Github Action Workflow Files](/.github/workflows/) to find the correct test command for a module.

# Roadmap

Disgo is currently in DEVELOPMENT. Here are the steps required in order to complete it:

1. Merge Rate Limits (HTTP, Gateway)
2. **Implement Testing** _[usable at this stage]_.
3. Implement Sharding.
4. Implement Cache (which is likely where most effort lies; caching is difficult).

In addition, we must make [decisions](/_contribution/libraries/) for the following:
1. UDP connections (Voice)
2. [Audio Processing using Opus](https://discord.com/developers/docs/topics/voice-connections#encrypting-and-sending-voice)
