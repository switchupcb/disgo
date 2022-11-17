# Example: Create a Followup Response

This example creates an application command, listens for an interaction _(while connected to the Discord Gateway)_, then responds to the interaction with multiple responses. First, an original response is sent and edited. Then, a followup message to that response is sent and edited. The bot waits until the program receives a signal to terminate, then disconnects from the Discord Gateway and deletes the created command. For more information about interaction responses, read the [Discord API Documentation](https://discord.com/developers/docs/interactions/receiving-and-responding#responding-to-an-interaction).

_Use this example to create responses to interactions from users._

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

Use `go build` to build the executable binary. Use `followup` to run it from the command line.

```
> followup
Program is started.
Creating an application command...
Adding an event handler.
Connecting to the Discord Gateway...
Successfully connected to the Discord Gateway. Waiting for an interaction...
```

### /followup

Upon receiving an interaction from the user, the bot will respond with a regular interaction response. Once this has been sent, it waits two seconds and edits that response. Once edited, it waits two seconds and sends a followup message. Once that is sent, it waits two more seconds to edit the followup message.

```
0:00 followup called by Flame.
0:00 Creating a response to the interaction...
0:00 Sent original interaction response.
0:02 Editing the original response to the interaction...
0:04 Edited original interaction response.
0:06 Sending a followup message to the interaction...
0:06 Sent a followup message to the interaction.
0:08 Editing the followup message to the interaction...
0:08 Edited the followup message to the interaction.
```

### SIGINT

Use `ctrl + C` or `cmd + C` in the terminal.

```
Deleting the application command...
Program executed successfully.
```
