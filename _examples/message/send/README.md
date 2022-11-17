# Example: Send a Message

This example sends a message and/or files to a channel.

_Use this example to send content to any given channel. Do **NOT** use it to respond to [interactions](/_examples/command/)._

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
> send -h
Usage of send:
  -c string
        Set the channel (ID) the message will be sent to using -c.
  -f string
        Set the location (filepath or URL) of the file using -f.
  -m string
        Set the text content of the message using -m.
```

## Usage

Use `go build` to build the executable binary. Use `send` to run it from the command line.

_NOTE: Get the Channel ID by enabling **Developer Mode** from the settings of your account, then right clicking any channel._

### Text

```
send -c 1041179872518737990 -m "This is a message."
```

### Emoji

```
send -c 1041179872518737990 -m ":smile:"
```

### File

```
send -c 1041179872518737990 -f file.png
```

### URL

```
send -c 1041179872518737990 -f https://assets-global.website-files.com/6257adef93867e50d84d30e2/636e0b52aa9e99b832574a53_full_logo_blurple_RGB.png
```

### Both

```
send -c 1041179872518737990 -m "This is a message. :smile:" -f file.png
```

