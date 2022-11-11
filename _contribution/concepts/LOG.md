# What is a Log?

A **log** is a record that contains information about an application's runtime. Logs are useful for the development and debugging of a program. For more information, read [Logging (Software)](https://en.wikipedia.org/wiki/Logging_(software)).

## What is an Unstructured Log?

An **unstructured log** — in Software Development — is a log _(record)_ without structure. For example, the Go standard library [`log`](https://pkg.go.dev/log) implements unstructured logging.

```go
log.Println("This creates an unstructured log.")
```

```
2009/11/10 23:00:00 This creates an unstructured log.
```

## What is a Structured Log?

A **structured log** — in Software Development — is a log _(record)_ with a specified structure _(i.e CBOR, JSON)_. This allows developers to programmatically analyze and process logs.

```json
{"time":1516134303,"level":"info","message":"This is structured log in the JSON format."}
```

# Disgo Logger

Disgo uses [`rs/zerolog`](https://github.com/rs/zerolog) to provide customizable zero-allocation structured logging.

## Usage

Disgo provides leveled logging of the API Wrapper via the `Logger` global variable. As a result, this variable is accessible via `disgo.Logger`. For more information on the usage of the `zerolog.Logger`, check out its [features](https://github.com/rs/zerolog#features).

## Configuration

Disgo imports the `zerolog` module which allows you to use it in your own program. As a result, the `zerolog` package can be configured from your application. For more information on the configuration of `zerolog`, check out its [Global Settings](https://github.com/rs/zerolog#global-settings).
