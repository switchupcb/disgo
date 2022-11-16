# Example: Create an Application Command

This example creates an application command, listens for an interaction _(while connected to the Discord Gateway)_, then responds to the interaction. Following this response, the bot deletes the application command, then disconnects from the Discord Gateway.

_Use this example to create commands users can interact with._

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

### /main

```
> command
2022/11/15 17:55:39 Program is started.
2022/11/15 17:55:39 Creating an application command...
2022/11/15 17:55:40 Adding an event handler.
2022/11/15 17:55:40 Connecting to the Discord Gateway...
2022/11/15 17:55:40 Successfully connected to the Discord Gateway. Waiting for an interaction...
2022/11/15 17:55:44 main called by SCB.
2022/11/15 17:55:44 Creating a response to the interaction...
2022/11/15 17:55:45 Deleting the application command...
2022/11/15 17:55:45 Disconnecting from the Discord Gateway...
2022/11/15 17:55:45 Program executed successfully.
```

### SIGINT

Use `ctrl + C` or `cmd + C` in the terminal.

```
2022/11/15 17:53:11 Program is started.
2022/11/15 17:53:11 Creating an application command...
2022/11/15 17:53:12 Adding an event handler.
2022/11/15 17:53:12 Connecting to the Discord Gateway...
2022/11/15 17:53:12 Successfully connected to the Discord Gateway. Waiting for an interaction...
2022/11/15 17:53:14 Exiting program due to signal...
2022/11/15 17:53:14 Program exited successfully.
```
