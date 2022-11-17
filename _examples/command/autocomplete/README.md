# Example: Autocomplete User Input

This example creates an application command with two required options, one of which uses required choices, and the other which uses autocompletion. Once the command is created, the bot listens for interactions _(while connected to the Discord Gateway)_, then responds to the interaction accordingly. The bot waits until the program receives a signal to terminate, then disconnects from the Discord Gateway and deletes the created command. For more information about autocompletion, read the [Discord API Documentation](https://discord.com/developers/docs/interactions/application-commands#autocomplete).

_Use this example to provide option validation and autocompletion to users._

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

Use `go build` to build the executable binary. Use `autocomplete` to run it from the command line.

```
> autocomplete
Program is started.
Creating an application command...
Adding an event handler.
Connecting to the Discord Gateway...
Successfully connected to the Discord Gateway. Waiting for an interaction...
```

### /autocomplete

The bot sends a response upon receiving an interaction from the user. This response is dependent on whether autocompletion data or application command data has been provided.

#### Autocomplete

The `freewill` option has autocomplete disabled, since autocompletion may **NOT** be enabled (by Discord) when required option choices are provided. The `confirm` option has autocomplete enabled, but this does **NOT** constrain the user to any specific choice. In other words, autocomplete choices will be provided, but the user can still input any value for the `confirm` option.

Due to the constraints above, all autocompletion data is triggered by a focused non-nil `confirm` option. In any case, server-side validation is provided for each option. When no options are focused, autocompletion data is **NOT** sent. When `freewill` and `confirm` are both EMPTY _(i.e `""`)_, both choices will be sent to the user. When `confirm` is selected (by the user), the choice opposite of the selected `freewill` choice will be sent.

| `freewill` | `confirm` | Focused    | Choice |
| :--------- | :-------- | :--------- | :----- |
| N/A        | N/A       | N/A        | N/A    |
| any        | N/A       | `freewill` | N/A    |
| N/A        | any       | `confirm`  | Both   |
| Yes        | any       | `confirm`  | `No`   |
| No         | any       | `confirm`  | `Yes`  |

The hidden objective for the user is to select and confirm an answer to the question, **"Do you have free will?"** This requires effort from the user because the program will **NOT** let the command be sent unless `freewill` is filled _(with `Yes` or `No`)_, and will always recommended the opposite value of that choice during `confirm` autocompletion. The user can `confirm` their `freewill` choice by manually entering their provided `confirm` option value. 

#### Command

When application command data is received, the bot will respond with a message.

| `freewill` | `confirm` | `==`    | Message                               |
| :--------- | :-------- | :------ | :------------------------------------ |
| `Yes`      | any       | `false` | `Hmmm. I guess you aren't sure...`    |
| `No`       | any       | `false` | `Hmmm. I guess you aren't sure...`    |
| `Yes`      | `Yes`     | `true`  | `Where there's a will there's a way.` |
| `No`       | `No`      | `true`  | `Fate awaits you.`                    |


### SIGINT

Use `ctrl + C` or `cmd + C` in the terminal.

```
Deleting the application command...
Program executed successfully.
```