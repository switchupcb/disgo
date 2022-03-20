# Contributing

## License

You agree to license any contribution to this library under the [Apache License 2.0](#license).

## Pull Requests

Pull requests must follow the [code specification](#code-specification) and work with all [test cases](#test).

## Domain

The domain of Disgo in providing an API for HTTP/WebSocket requests. The program uses provided structures (from the [Discord API](https://discord.com/developers/docs/reference)) to provide simple abstractions for end users (developers).

## Project Structure

The repository consists of a detailed [README](/README.md), [examples](/examples/), and [**API Wrapper**](/disgo/).

### Disgo

The API Wrapper consists of multiple packages. A [bundler (TBA)]() is used to package the API into a single-usable package.

| Package   | Description                                   |
| :-------- | :-------------------------------------------- |
| wrapper   | Contains the wrapper bundling functionality.  |
| client    | The Disgo Client Abstraction.                 |
| events    | Discord API Events.                           |
| pkg       | Utility functionality for specific libraries. |
| requests  | Discord API Requests.                         |
| resources | Discord API Resources.                        |
| sessions  | Discord API WebSocket Sessions (Gateways).    |

### Code Specification

#### Comments

Comments follow [Effective Go](https://golang.org/doc/effective_go#commentary) and explain why more than what _(unless the "what" isn't intuitive)_.

#### Static Code Analysis

Disgo uses [golangci-lint](https://github.com/golangci/golangci-lint) in order to statically analyze code. You can install golangci-lint with `go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.45.0` and run it using `golangci-lint run`. If you receive a `diff` error, you must add a `diff` tool in your PATH. There is one located in the `Git` bin.

If you receive `File is not ... with -...`, use `golangci-lint run --disable-all --no-config -Egofmt --fix`.

### Test

#### Unit Tests

Unit tests are used to test logic.

#### Integration Tests

Integration tests are used to ensure functionality between the API and Discord.

# Roadmap
