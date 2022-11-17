# Example: Create a Subcommand

This example creates an application command with a subcommand group, subcommands, and options. Once the command is created, it listens for an interaction _(while connected to the Discord Gateway)_, then responds to the interaction _(with a calculation based on user input)_. The bot waits until the program receives a signal to terminate, then disconnects from the Discord Gateway. For more information about subcommands, read the [Discord API Documentation](https://discord.com/developers/docs/interactions/application-commands#subcommands-and-subcommand-groups).

_Use this example to create subcommands users can interact with._

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

Use `go build` to build the executable binary. Use `subcommand` to run it from the command line.

```
> subcommand
Program is started.
Creating an application command...
Adding an event handler.
Connecting to the Discord Gateway...
Successfully connected to the Discord Gateway. Waiting for an interaction...
```

### /calculate

The following output results from using a command in any channel with the bot (includes direct messages).

#### add int

```
/calculate add int
    There was nothing to do.

/calculate add int 1
    1

/calculate add int 2 4
    6
```

#### add string

```
/calculate add string
    There was nothing to do.

/calculate add string 1
    1

/calculate add int 2 4
    24
```

#### subtract

```
/calculate subtract 5.3 2
    3.300
```

_User input of both options is required to submit `/calculate subtract`._

### SIGINT

Use `ctrl + C` or `cmd + C` in the terminal.

```
Exiting program due to signal...
Program exited successfully.
```
