# Example: Send a Message Component

This example sends message components to a channel.

_This example is UNFINISHED: User input to message components are handled using [interactions](/_examples/command/)._

## Setup

**You must create a Discord Application in the [Discord Developer Portal](https://discord.com/developers/docs/getting-started#creating-an-app) to receive your Bot Token.** 

### Environment Variables

Assign an environment variable in the command line you will be running the program from.

#### Windows

```
set TOKEN=value
```

#### Mac/Linux

```
export TOKEN=value
``` 

**NEVER SHOW YOUR TOKEN TO THE PUBLIC.**

### Flags

This example uses [command line flags](https://pkg.go.dev/flag) to input the channel and message content.

```
> components -h
Usage of send:
  -c string
        Set the channel (ID) the message components will be sent to using -c.
```

## Usage

Use `go build` to build the executable binary. Use `components` to run it from the command line.

_NOTE: Get the Channel ID by enabling **Developer Mode** from the settings of your account, then right clicking any channel._

```
components -c 1041179872518737990
```
