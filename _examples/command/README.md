# Example: Create an Application Command

This example creates an application command, listens for an interaction _(while connected to the Discord Gateway)_, then responds to the interaction. Following this response, the bot deletes the application command, then disconnects from the Discord Gateway. For more information about Application Commands, read the [Discord API Documentation](https://discord.com/developers/docs/interactions/application-commands#application-commands).

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
Program is started.
Creating an application command...
Adding an event handler.
Connecting to the Discord Gateway...
Successfully connected to the Discord Gateway. Waiting for an interaction...
main called by SCB.
Creating a response to the interaction...
Deleting the application command...
Disconnecting from the Discord Gateway...
Program executed successfully.
```

### SIGINT

Use `ctrl + C` or `cmd + C` in the terminal.

```
Program is started.
Creating an application command...
Adding an event handler.
Connecting to the Discord Gateway...
Successfully connected to the Discord Gateway. Waiting for an interaction...
Exiting program due to signal...
Program exited successfully.
```

# Read More

| Example                                        | Description                                                                     |
| :--------------------------------------------- | :------------------------------------------------------------------------------ |
| [`subcommand`](/_examples/command/subcommand/) | Create an application command with subcommand groups, subcommands, and options. |
