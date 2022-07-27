# Disgo

[![Go Doc](https://img.shields.io/badge/godoc-reference-5272B4.svg?style=for-the-badge&logo=appveyor&logo=appveyor)](https://pkg.go.dev/github.com/switchupcb/disgo)
[![License](https://img.shields.io/github/license/switchupcb/disgo.svg?style=for-the-badge)](https://github.com/switchupcb/disgo/blob/main/LICENSE)

**This repository is currently in DEVELOPMENT. For more information, read the [roadmap](/_contribution/CONTRIBUTING.md#roadmap).**

Create a Discord Bot in Go using Disgo. This [Discord API](https://discord.com/developers/docs/reference) Wrapper is designed to be flexible, performant, secure, and thread-safe. Disgo aims to provide every feature in the Discord API along with optional caching, shard management, rate limiting, and logging. Use the only Go module to provide a **100% one-to-one implementation** of the Discord API.

**A Next Generation Discord API Wrapper**

High quality code merits easy development. Disgo uses developer operations to stay up-to-date with the ever-changing Discord API. Code generation is used to provide a clean implementation for every request and event. Data race detection is used with _an integration test that covers nearly 100% of the Discord API_ in order to ensure that Disgo is safe for concurrent usage.

**Don't Miss Out On These Exclusive Features**

- EVERY Rate Limit (Global, Per Route, Per Resource, Emoji, Gateway)
- Automatic Intent Calculation (Gateway)

## Table of Contents

| Topic                           | Categories                                                                                                                                          |
| :------------------------------ | :-------------------------------------------------------------------------------------------------------------------------------------------------- |
| [Using the API](#using-the-api) | [Breakdown](#using-the-api), [Sharding](#sharding), [Caching](#caching)                                                                             |
| [Examples](#examples)           | [Configuration](#configuration), [Create a Command](#create-a-command), [Handle an Event](#handle-an-event), [Output](#output), [Summary](#Summary) |
| [Features](#features)           | [Why Go?](#why-go), [Comparison](#comparison), [Contributing](#contributing)                                                                        |
| [Ecosystem](#ecosystem)         | [License](#license), [Libraries](#libraries), [Credits](#credits)                                                                                   |

## Using the API

This breakdown provides you with a **full understanding** on how to use the API.

| Abstraction  | Usecase                                                                                                                                                            | Example                                                             |
| :----------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------------- | :------------------------------------------------------------------ |
| **Resource** | A [Discord API Resource](https://discord.com/developers/docs/resources/application).                                                                               | Guild Object. User Object.                                          |
| **Event**    | A [Discord API Event](https://discord.com/developers/docs/topics/gateway#commands-and-events-gateway-events).                                                      | A message is created. A user joins a channel.                       |
| **Client**   | The Discord Bot [Application](https://discord.com/developers/docs/resources/application) that you program. One Bot = One Client.                                   | Configure the bot settings. Set the token.                          |
| **Request**  | Uses the Discord HTTP REST API to make one-time requests for information _(i.e resources)_. Provides create, read, update, delete, patch endpoints.                | Create a command. Request Guild Info.                               |
| **Session**  | Uses Discord WebSockets [(Gateways)](https://discord.com/developers/docs/topics/gateway) to receive ongoing **events** that contain information _(i.e resources)_. | Send a message when a command used or a user joins a voice channel. |

You create a **Client** that calls for **Resources** using **Requests** and handles **Events** from **Sessions** using event handlers. For more information, please read [What is a Request?](/_contribution/concepts/REQUESTS.md) and [What is an Event?](/_contribution/concepts/EVENTS.md)

### Flags

A flag is a [flag](https://discord.com/developers/docs/resources/application#application-object-application-flags), [type](https://discord.com/developers/docs/resources/channel#embed-object-embed-types), [key](https://discord.com/developers/docs/resources/audit-log#audit-log-change-object-audit-log-change-key), [level](https://discord.com/developers/docs/resources/guild#guild-object-verification-level) or any other option that Discord provides. All flags are denoted by `Flag` in disgo: For example, `disgo.FlagUserSTAFF`, `disgo.FlagVerificationLevelHIGH`, `disgo.FlagPremiumTierNONE`, etc.

### Sharding

Read [What is a Discord Shard](/_contribution/concepts/SHARD.md) for a simple yet full understanding of sharding on Discord. Using the [Shard Manager](/_contribution/concepts/SHARD.md#the-shard-manager) is **optional**. You can manually implement a shard manager through the `disgo.Client.Sessions` array.

### Caching

Read [What is a Cache](/_contribution/concepts/CACHE.md) for a simple yet full understanding of the Disgo Cache. The [Disgo Cache](/_contribution/concepts/CACHE.md#the-disgo-cache) is **optional**. The **cache interface** allows you to replace the built-in cache with another store _(such as Redis or Memcached)_ and/or provide your own method of caching data.

## Examples

The **main example** creates a bot that creates an application command and handles it. Check out the [examples](/_examples/) directory for more examples.

## Configuration

**You must create a Discord Application in the [Discord Developer Portal](https://discord.com/developers/docs/getting-started#creating-an-app) to receive your Bot Token.** 

Use the client to configure the bot's settings.
```go
bot := &disgo.Client{
    ApplicationID: "APPID", // optional
    Authentication: disgo.BotToken("TOKEN"), // or BearerToken("TOKEN")
    Authorization: &disgo.Authorization{ ... },
    Config: disgo.DefaultConfig(),
    Handlers: new(Handlers),
    Sessions: new(Sessions)
}
```

## Create a Command

Create an application command **request** to add an application command.

```go
// Create a Create Global Application Command request.
request := disgo.CreateGlobalApplicationCommand{
    Name: "main",
    Description: "A basic command",
} 

// Register the new command by sending the request to Discord using the bot.
// returns a disgo.ApplicationCommand
newCommand, err := request.Send(bot)
if err != nil {
    log.Printf("failure sending command to Discord: %v", err)
}
```

## Handle an Event

Create an **event handler** and add it to the **bot**.

```go
// Add an event handler to the bot.
bot.Handle(disgo.FlagGatewayEventNameInteractionCreate, func(i disgo.InteractionCreate) {
	log.Printf("main called by %s", i.User.Username)
})
```

_Disgo provides automatic intent calculation._

### Output

Open a WebSocket **Session** to receive events.

```go
// Connect the session to the Discord Gateway (WebSocket Connection).
if err := bot.Connect(disgo.NewSession()); err != nil {
    log.Printf("can't open websocket session to Discord: %v", err)
}
```

The following message will be logged when a user creates an [`InteractionCreate`](https://discord.com/developers/docs/topics/gateway#commands-and-events-gateway-events) event by using `/main` in a Direct Message with the bot on Discord.

```
main called by SCB
```

### Summary

```go
// Use resources to represent Discord objects in your application.
disgo.<API Resources>

// Use events to represent Discord events in your application.
disgo.<API Events>

// Use the client to manage the bot's settings.
disgo.Client.Config.Request.<Settings>
disgo.Client.Config.Gateway.<Settings>
disgo.Client.Authentication.<Settings>
disgo.Client.Authorization.<Settings>

// Use requests to exchange data with Discord's REST API.
disgo.<Endpoint>.Send()

// Use sessions to handle events from Discord's WebSocket Sessions (Gateways).
disgo.Client.Handle(<event>, <handler>)
disgo.Client.Remove(<event>, <index>)

// Use flags to specify options.
disgo.Flag<Option><Name>

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

Go is a statically typed language with a garbage collector. As a result, it performs computationally better compared to most languages that provide [Discord API Wrappers](https://discord.com/developers/docs/topics/community-resources#libraries). Go maintains superior asynchronous handling due to the use of [Goroutines](https://gobyexample.com/goroutines) and [Channels](https://gobyexample.com/channels). This is useful since **a Discord Bot is a server-side software**.

### Comparison

Disgo supports every feature in the Discord API and is **the most customizable Discord API Wrapper** due to its optional caching, shard management, rate limiting, and logging. **DiscordGo** is not feature-complete and **Disgord** is limiting. Look no further than the name. The word `disgo` contains 5 letters — while the others have 7+ — saving you precious keyboard strokes. Most important is Disgo's performance, which saves you money by reducing server costs. _Don't believe me?_ Check this out!

#### CPU

Disgo places a priority on performance. For more information, view [`library decisions`](/_contribution/libraries/).

### Memory

Every struct uses [fieldalignment](https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/fieldalignment) to reduce the memory footprint of your application.

### Storage

Disgo adds **3.5 MB** to a compiled binary.

### Contributing

Disgo is the easiest Discord Go API for developers to use and contribute to. You can contribute to this repository by viewing the [Project Structure, Code Specifications, and Roadmap](/_contribution/CONTRIBUTING.md).

| Library   | Contribution                                                                                                                                                                                                                              | Lines of Code to Maintain |
| :-------- | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | :------------------------ |
| Disgo     | [Contribution Guidelines](contribution/CONTRIBUTING.md), [Project Architecture](contribution/CONTRIBUTING.md#project-structure), [Linting](contribution/CONTRIBUTING.md#static-code-analysis), [Tests](contribution/CONTRIBUTING.md#test) | 6K/15K                    |
| DiscordGo | No Guidelines, No Architecture, No Linter, Not Feature Complete                                                                                                                                                                           | 12K/12K                   |
| Disgord   | Contribution Guidelines, No Linter, ORM, Not Feature Complete                                                                                                                                                                             | ?/30K                     |

## Ecosystem

### License

The [Apache License 2.0](#license) is permissive for commercial use. For more information, read [Apache Licensing FAQ](https://www.apache.org/foundation/license-faq.html).

### Libraries

| Library                                          | Description                                             |
| :----------------------------------------------- | :------------------------------------------------------ |
| [Copygen](https://github.com/switchupcb/copygen) | Generate custom type-based code.                        |
| [Dasgo](https://github.com/switchupcb/dasgo)     | Go Type Definitions for the Discord API.                |
| Disgo Template                                   | Get started on a Discord Bot with this Disgo Framework. |

### Credits

| Name                                      | Contributions                                             |
| :---------------------------------------- | :-------------------------------------------------------- |
| [SwitchUpCB](https://switchupcb.com)      | Project Architecture, Dasgo, Requests, WebSockets, Events |
| [Thomas Rogers](https://github.com/t-rog) | Dasgo, WebSockets                                         |
| [Josh Dawe](https://github.com/joshdawe)  | Dasgo                                                     |

_Earn a credit! [Contribute Now](_contribution/CONTRIBUTING.md)._