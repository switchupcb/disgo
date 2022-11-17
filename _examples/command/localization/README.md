# Example: Create a Localized Command

This example creates an application command with a name and description based on the locale of the user interacting with the bot. Once the command is created, it listens for an interaction _(while connected to the Discord Gateway)_, then responds to the interaction accordingly. The bot waits until the program receives a signal to terminate, then disconnects from the Discord Gateway and deletes the created command. For more information about localization, read the [Discord API Documentation](https://discord.com/developers/docs/interactions/application-commands#localization).

_Use this example to create localized commands for users._

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

Use `go build` to build the executable binary. Use `localization` to run it from the command line.

```
> localization
Program is started.
Creating an application command...
Adding an event handler.
Connecting to the Discord Gateway...
Successfully connected to the Discord Gateway. Waiting for an interaction...
```

### /hello

Upon receiving an interaction from the user, the bot will respond with a message in the language of the user's locale. When the user's locale is not supported, the bot will state that `The current locale is not supported by this command.`

### SIGINT

Use `ctrl + C` or `cmd + C` in the terminal.

```
Deleting the application command...
Program executed successfully.
```
