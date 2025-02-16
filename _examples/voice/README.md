# Example: Record a Sound and Speak

This example joins a voice channel, records a sound, then speaks the recorded sound on request. The bot waits until the program receives a signal to terminate, then disconnects from the Discord Gateway.

For more information about Discord Voice, read the [Discord API Documentation](https://discord.com/developers/docs/topics/voice-connections).

_Use this example to implement voice channel and audio interaction._

## Setup

**You must create a Discord Application in the [Discord Developer Portal](https://discord.com/developers/docs/getting-started#creating-an-app) to receive your Bot Token.** 

### Environment Variables

Assign an environment variable in the command line you will be running the program from.

#### Windows

```
set TOKEN=value
set APPID=value
```

#### Mac/Linux

```
export TOKEN=value
export APPID=value
``` 

**NEVER SHOW YOUR TOKEN TO THE PUBLIC.**

_NOTE: Get the Application ID by enabling **Developer Mode** from the settings of your account, then right clicking your bot._

## Usage

Use `go build` to build the executable binary. Use `command` to run it from the command line.

```
> command
Program is started.
Creating an application command...
Adding an event handler.
Connecting to a voice channel...
Successfully connected to a voice channel. Waiting for an interaction...
```

### /record

```
record called by SCB.
Recording voice channel for 5 seconds...
Recording completed.
```

### /speak

```
speak called by SCB.
Speaking the recorded sound to the voice channel...
Speaking completed.
```

### SIGINT

Use `ctrl + C` or `cmd + C` in the terminal.

```
Exiting program due to signal...
Disconnected from the Discord Gateway.
Deleting the application command...
Program executed successfully.
```