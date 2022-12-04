# Create a Discord Bot in Go

[![Go Doc](https://img.shields.io/badge/godoc-reference-5272B4.svg?style=for-the-badge&logo=appveyor&logo=appveyor)](https://pkg.go.dev/github.com/switchupcb/disgo)
[![License](https://img.shields.io/github/license/switchupcb/disgo.svg?style=for-the-badge)](https://github.com/switchupcb/disgo/blob/main/LICENSE)

**Disgo** is a [Discord API](https://discord.com/developers/docs/reference) Wrapper designed to be flexible, performant, secure, and thread-safe. Disgo aims to provide every feature in the Discord API along with optional rate limiting, structured logging, shard management, and caching. Use the only Go module to provide a **100% one-to-one implementation** of the Discord API.

_This repository is STABLE. For more information, read the [roadmap](/_contribution/CONTRIBUTING.md#roadmap)._

## A Next Generation Discord API Wrapper

High quality code merits easy development. Disgo uses developer operations to stay up-to-date with the ever-changing Discord API. Code generation is used to provide a clean implementation for every request and event. Data race detection is used with _an integration test that covers the entire Discord API_ in order to ensure that Disgo is safe for concurrent usage. In addition, **Disgo provides the following exclusive features**.

- [EVERY Rate Limit (Global, Per Route, Per Resource, Custom, Gateway)](_contribution/concepts/REQUESTS.md#what-is-a-rate-limit) 
- [Automatic Gateway Intent Calculation](_contribution/concepts/EVENTS.md#what-is-a-gateway-intent)
- [Selective Event Processing](_contribution/concepts/EVENTS.md#selective-event-processing)

_Disgo uses [NO reflection or type assertion](_contribution/concepts/EVENTS.md#how-it-works)._

## Table of Contents

| Topic                           | Categories                                                                                                                                                             |
| :------------------------------ | :--------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [Using the API](#using-the-api) | [Breakdown](#using-the-api), [Logging](#logging), [Sharding](#sharding), [Caching](#caching)                                                                           |
| [Examples](#examples)           | [Import](#import), [Configuration](#configuration), [Create a Command](#create-a-command), [Handle an Event](#handle-an-event), [Output](#output), [Summary](#Summary) |
| [Features](#features)           | [Why Go?](#why-go), [Comparison](#comparison), [Contributing](#contributing)                                                                                           |
| [Ecosystem](#ecosystem)         | [License](#license), [Libraries](#libraries), [Credits](#credits)                                                                                                      |

## Using the API

This breakdown provides you with a **full understanding** on how to use the API.

| Abstraction  | Usecase                                                                                                                                           | Example                                                             |
| :----------- | :------------------------------------------------------------------------------------------------------------------------------------------------ | :------------------------------------------------------------------ |
| **Resource** | A [Discord API Resource](https://discord.com/developers/docs/resources/application).                                                              | Guild Object. User Object.                                          |
| **Event**    | A [Discord API Event](https://discord.com/developers/docs/topics/gateway#commands-and-events-gateway-events).                                     | A message is created. A user joins a channel.                       |
| **Client**   | The Discord Bot [Application](https://discord.com/developers/docs/resources/application) that you program. One Bot = One Client.                  | Configure the bot settings. Set the token.                          |
| **Request**  | Uses the Discord HTTP REST API to make one-time _requests_ for information.                                                                       | Create an application command. Request guild information.           |
| **Session**  | Uses a Discord WebSocket Connection [(Gateway)](https://discord.com/developers/docs/topics/gateway) to receive _events_ that contain information. | Send a message when a command used or a user joins a voice channel. |

You create a **Client** that calls for **Resources** using **Requests** and handles **Events** from **Sessions** using event handlers. For more information, please read [What is a Request?](/_contribution/concepts/REQUESTS.md) and [What is an Event?](/_contribution/concepts/EVENTS.md)

### Flags

A flag is a [flag](https://discord.com/developers/docs/resources/application#application-object-application-flags), [type](https://discord.com/developers/docs/resources/channel#embed-object-embed-types), [key](https://discord.com/developers/docs/resources/audit-log#audit-log-change-object-audit-log-change-key), [level](https://discord.com/developers/docs/resources/guild#guild-object-verification-level) or any other option that Discord provides. All flags are denoted by `Flag` in disgo: For example, `disgo.FlagUserSTAFF`, `disgo.FlagVerificationLevelHIGH`, `disgo.FlagPremiumTierNONE`, etc.

### Logging

Read [What is a Log](/_contribution/concepts/LOG.md) for a simple yet full understanding of logging. Disgo provides structured, leveled logging of the API Wrapper via the `disgo.Logger` global variable _(disabled by default)_. Enable the logger using `zerolog.SetGlobalLevel(zerolog.LEVEL)`.

### Sharding

Read [What is a Discord Shard](/_contribution/concepts/SHARD.md) for a simple yet full understanding of sharding on Discord. Using the [Shard Manager](/_contribution/concepts/SHARD.md#the-shard-manager) is **optional**. You can manually implement a shard manager through the `disgo.Client.Sessions` array.

### Caching

Read [What is a Cache](/_contribution/concepts/CACHE.md) for a simple yet full understanding of the Disgo Cache. The [Disgo Cache](/_contribution/concepts/CACHE.md#the-disgo-cache) is **optional**. The **cache interface** allows you to replace the built-in cache with another store _(such as Redis or Memcached)_ and/or provide your own method of caching data.

## Examples

| Example                        | Description                                                |
| :----------------------------- | :--------------------------------------------------------- |
| main                           | Learn how to use `disgo`.                                  |
| [command](/_examples/command/) | Create an application command and respond to interactions. |
| [message](/_examples/message/) | Send a message with text, emojis, files and/or components. |
| [image](/_examples/image/)     | Set the bot's avatar using an image.                       |

_Check out the [examples](/_examples/) directory for more._

### Import

Get a specific version of `disgo` by specifying a tag or branch.

```
go get github.com/switchupcb/disgo@v0.10.1
```

_Disgo branches are referenced by API version (i.e `v10`)._

_DISCLAIMER: `v0.10.1` is a pre-release version. For more information, read the [State of Disgo (v0.10.1)](https://github.com/switchupcb/disgo/discussions/40)._

### Configuration

**You must create a Discord Application in the [Discord Developer Portal](https://discord.com/developers/docs/getting-started#creating-an-app) to receive your Bot Token.** 

Use the client to configure the bot's settings.
```go
bot := &disgo.Client{
    ApplicationID:  "APPID", // optional
    Authentication: disgo.BotToken("TOKEN"), // or BearerToken("TOKEN")
    Authorization:  &disgo.Authorization{ ... },
    Config:         disgo.DefaultConfig(),
    Handlers:       new(disgo.Handlers),
    Sessions:       []*disgo.Session{disgo.NewSession()},
}
```

_Need more information? Read the [bot example](/_examples/bot)._

### Create a Command

Create an application command **request** to add an application command.

```go
// Create a Create Global Application Command request.
request := disgo.CreateGlobalApplicationCommand{
    Name:        "main",
    Description: disgo.Pointer("A basic command."),
}

// Register the new command by sending the request to Discord using the bot.
//
// returns a disgo.ApplicationCommand
newCommand, err := request.Send(bot)
if err != nil {
    log.Printf("failure sending command to Discord: %v", err)

    return
}
```

### Handle an Event

Create an **event handler** and add it to the **bot**.

```go
// Add an event handler to the bot.
bot.Handle(disgo.FlagGatewayEventNameInteractionCreate, func(i *disgo.InteractionCreate) {
	log.Printf("main called by %s", i.User.Username)
})
```

_Disgo provides automatic [Gateway Intent](https://discord.com/developers/docs/topics/gateway#gateway-intents) calculation._

### Output

Open a WebSocket **Session** to receive events.

```go
// Connect the session to the Discord Gateway (WebSocket Connection).
if err := bot.Sessions[0].Connect(bot); err != nil {
    log.Printf("can't open websocket session to Discord Gateway: %v", err)

	return
}
```

The following message will be logged when a user creates an `InteractionCreate` event by using `/main` in a Direct Message with the Bot on Discord.

```
main called by SCB.
```

### Summary

```go
// Use flags to specify options from Discord.
disgo.Flag<Option><Name>

// Use resources to represent Discord objects in your application.
disgo.<API Resources>

// Use events to represent Discord events in your application.
disgo.<API Events>

// Use requests to exchange data with Discord's REST API.
disgo.<Endpoint>.Send()

// Use sessions to connect to the Discord Gateway.
disgo.Session.Connect()
disgo.Session.Disconnect()

// Use event handlers to handle events from Discord's Gateway.
disgo.Client.Handle(<event>, <handler>)
disgo.Client.Remove(<event>, <index>)
disgo.Client.Handlers.<Handler>

// Use the client to manage the bot's settings.
disgo.Client.ApplicationID
disgo.Client.Authentication.<Settings>
disgo.Client.Authorization.<Settings>
disgo.Client.Config.Request.<Settings>
disgo.Client.Config.Gateway.<Settings>

// Use the client's shard manager to handle sharding automatically or manually.
disgo.Client.Shard.<Settings>
disgo.Client.Shard.<map[Session][]map[Shard][]GuildIDs>

// Use the client to manage the optional cache.
disgo.Client.Cache.<Settings>
disgo.Client.Cache.<Requests>
disgo.Client.Cache.<...>
```

## Features

### Why Go?

Go is a statically typed, compiled programming language _(with a garbage collector)_. As a result, it performs computationally better compared to _most_ languages that provide [Discord API Wrappers](https://discord.com/developers/docs/topics/community-resources#libraries). Go maintains superior asynchronous handling due to the use of [Goroutines](https://gobyexample.com/goroutines) and [Channels](https://gobyexample.com/channels). This is useful since **a Discord Bot is a server-side software**.

### Comparison

Disgo supports every feature in the Discord API and is **the most customizable Discord API Wrapper** due to its optional caching, shard management, rate limiting, and logging. **DiscordGo** is not feature-complete and **Disgord** is limiting. Look no further than the name. The word `disgo` contains 5 letters — while the others have 7+ — saving you precious keyboard strokes. Most important is Disgo's performance, which saves you money by reducing server costs. _Don't believe me?_ Check this out!

#### CPU

Disgo places a priority on performance. For more information, view [`library decisions`](/_contribution/libraries/).

### Memory

Every struct uses [fieldalignment](https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/fieldalignment) to reduce the memory footprint of your application.

### Storage

Disgo adds **5 MB** to a compiled binary.

### Contributing

Disgo is the easiest Discord Go API for developers to use and contribute to. You can contribute to this repository by viewing the [Project Structure, Code Specifications, and Roadmap](/_contribution/CONTRIBUTING.md).

## Ecosystem

### License

The [Apache License 2.0](#license) is permissive for commercial use. For more information, read [Apache Licensing FAQ](https://www.apache.org/foundation/license-faq.html).

### Libraries

| Library                                               | Description                              |
| :---------------------------------------------------- | :--------------------------------------- |
| [Copygen](https://github.com/switchupcb/copygen)      | Generate custom type-based code.         |
| [Dasgo](https://github.com/switchupcb/dasgo)          | Go Type Definitions for the Discord API. |
| [Ecosystem](https://github.com/switchupcb/disgo/wiki) | View projects that use Disgo.            |

### Credits

| Name                                      | Contributions                                                         |
| :---------------------------------------- | :-------------------------------------------------------------------- |
| [SwitchUpCB](https://switchupcb.com)      | Project Architecture, Generators, Dasgo, Requests, WebSockets, Events |
| [Thomas Rogers](https://github.com/t-rog) | Dasgo                                                                 |
| [Josh Dawe](https://github.com/joshdawe)  | Dasgo                                                                 |

_Earn a credit! [Contribute Now](_contribution/CONTRIBUTING.md)._
