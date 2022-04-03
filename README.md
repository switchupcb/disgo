# Disgo

**This repository is currently a PROOF OF CONCEPT. For more information, read the [roadmap](/contribution/CONTRIBUTING.md#roadmap).** 

Create a Discord Bot in Go using Disgo. This [Discord API](https://discord.com/developers/docs/reference) Wrapper is designed to be flexible, performant, and secure. Disgo aims to provide every feature in the Discord API along with optional caching and shard management.

| Topic                           | Categories                                                                                                                                          |
| :------------------------------ | :-------------------------------------------------------------------------------------------------------------------------------------------------- |
| [Using the API](#using-the-api) | [Breakdown](#using-the-api), [Caching](#caching), [Sharding](#sharding)                                                                             |
| [Examples](#examples)           | [Configuration](#configuration), [Create a Command](#create-a-command), [Handle an Event](#handle-an-event), [Output](#output), [Summary](#Summary) |
| [Features](#features)           | [Why Go?](#why-go), [Comparison](#comparison), [Contributing](#contributing)                                                                        |
| [License](#license)             | [Credits](#credits)                                                                                                                                 |

## Using the API

This breakdown provides you with a **full understanding** on how to use the API. 

| Abstraction  | Usecase                                                                                                                                                            | Example                                                             |
| :----------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------------- | :------------------------------------------------------------------ |
| **Resource** | A [Discord API Resource](https://discord.com/developers/docs/resources/application).                                                                               | Guild Object. User Object.                                          |
| **Event**    | A [Discord API Event](https://discord.com/developers/docs/topics/gateway#commands-and-events-gateway-events).                                                      | A message is created. A user joins a channel.                       |
| **Client**   | The Discord Bot [Application](https://discord.com/developers/docs/resources/application) that you program. One Bot = One Client.                                   | Configure the bot settings. Set the token.                          |
| **Request**  | Uses the Discord HTTPS/REST API to make one-time requests for information _(i.e resources)_. Provides create, read, update, delete, patch endpoints.               | Create a command. Request Guild Info.                               |
| **Session**  | Uses Discord WebSockets [(Gateways)](https://discord.com/developers/docs/topics/gateway) to receive ongoing **events** that contain information _(i.e resources)_. | Send a message when a command used or a user joins a voice channel. |

You create a **Client** that calls for **Resources** using **Requests** and that handles **Events** using **Sessions**.

### Caching

Read [What is a Cache](contribution/concepts/CACHE.md) for a simple yet full understanding of the Disgo Cache. The [Disgo Cache](contribution/concepts/CACHE.md#the-disgo-cache) is **optional**. The **cache interface** allows you to replace the built-in cache with another store _(such as Redis or Memcached)_ and/or provide your own method of caching data.

### Sharding

Read [What is a Discord Shard](contributing/concepts/SHARD.md) for a simple yet full understanding of sharding on Discord. Using the [Shard Manager](contribution/concepts/SHARD.md#the-shard-manager) is **optional**. You can manually implement a shard manager through the `disgo.Client.Sessions` array.

## Examples

Each example has a **README**.

| Example                  | Description                      |
| :----------------------- | :------------------------------- |
| [main](examples/main/)   | The default example.             |
| [shard](examples/shard/) | Uses the shard manager manually. |
| [cache](examples/shard/) | Uses the cache manually.         |

This [example](examples/main) creates a bot that creates an application command and handles it.

## Configuration

**You must create a Discord Application in the [Discord Developer Portal](https://discord.com/developers/applications) to receive your Bot Token.** 

Use the client to configure the bot's settings.
```go
bot := disgo.Client{
    Config: disgo.Config{

    },
}
```

## Create a Command

Create an application command **request** to add an application command.

```go
// Create a global command request.
request := disgo.RequestCreateApplicationCommand{
    Name: "main",
    Description: "A basic command",
} 

// Register the global command by sending the request to Discord.
// returns a disgo.ResourceApplicationCommand
newCommand, err := request.Send()
if err != nil {
    log.Println("error: failure sending command to Discord")
}
```

## Handle an Event

Create an **event handler** and add it to a **session**.

```go
// Add a session.
bot.Sessions = append(bot.Sessions, disgo.Session{})

// Define a handler (or use an anonymous function).
func handler (e disgo.EventInteractionCreate) {
    log.Println("/main called.")
}

// Add the handler for an event to the session.
bot.Sessions[0].Handlers.Add(handler)
```

### Output

Open the WebSocket **Session** to receive events.

```go
session, err := bot.Sessions[0].Open()
if err != nil {
    log.Println("error: can't open websocket session to Discord")
}
```

A user creates an interaction by using `/main` in a direct message.

[img]

### Summary

```go
// Use resources to represent Discord objects in your application.
disgo.Resource<API Resources>

// Use events to represent Discord events in your application.
disgo.Event<API Events>

// Use the client to manage the bot's settings.
disgo.Client.Config<Settings>

// Use requests to exchange data with Discord's REST API.
disgo.Request<Endpoints>

// Use sessions to handle events from Discord's WebSocket Sessions (Gateways).
disgo.Client.Session.<Handler func(disgo.Event){}>

// Use the client to manage the optional cache.
disgo.Client.Cache.<Settings>
disgo.Client.Cache.<Requests>
disgo.Client.Cache.<...>

// Use the client's shard manager to handle sharding automatically or manually.
disgo.Client.Shard.<Settings>
disgo.Client.Shard.<map[Session][]map[Shard][]GuildIDs>
```

## Features

### Why Go?

Go is a statically typed language with a garbage collector. As a result, it performs computationally better compared to most languages that provide [Discord API Wrappers](https://discord.com/developers/docs/topics/community-resources#libraries). Go maintains superior asynchronous handling due to the use of [Goroutines](https://gobyexample.com/goroutines) and [Channels](https://gobyexample.com/channels). This is useful since **a Discord Bot is a server-side software**.

### Comparison

Disgo supports every feature in the Discord API in the Discord API and provides optional caching and shard management. [DiscordGo](https://github.com/bwmarrin/discordgo#documentation) is not feature-complete and [Disgord](https://github.com/andersfylling/disgord#qa) is limiting. The word `disgo` contains 5 letters — while the others have 7+ — saving you precious keyboard strokes. Most important is Disgo's performance, which saves you money by reducing server costs. _Don't believe me?_ Check this out!

#### CPU
Disgo places a priority on performance. For more information, view [`library decisions`](/contribution/libraries/). Sharding is optional.

### Memory
Every struct uses [fieldalignment](https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/fieldalignment) to reduce the memory footprint of your application. Caching is optional.

### Storage
Disgo adds ~<> MB to a compiled binary.

### Contributing

Disgo is the easiest Discord Go API for developers to use and contribute to. You can contribute to this repository by viewing the [Project Structure, Code Specifications, and Roadmap](contribution/CONTRIBUTING.md).

| Library   | Contribution                                                                                                                                                                                                                                                                                               | Lines of Code to Maintain |
| :-------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | :------------------------ |
| Disgo     | [Contribution Guidelines](contribution/CONTRIBUTING.md), [Project Architecture](contribution/CONTRIBUTING.md#project-structure), [Code Generation](contribution/CONTRIBUTING.md#code-generation), [Linting](contribution/CONTRIBUTING.md#static-code-analysis), [Tests](contribution/CONTRIBUTING.md#test) | ?/?K                      |
| DiscordGo | Limited Guidelines, No Architecture, No Linter, Tests                                                                                                                                                                                                                                                      | ?/10K                     |
| Disgord   | Contribution Guidelines, Project Architecture, No Linter, Tests                                                                                                                                                                                                                                            | ?/30K                     |

## License

The [Apache License 2.0](#license) is permissive for commercial use. For more information, read [Apache Licensing FAQ](https://www.apache.org/foundation/license-faq.html).

### Credits

| Name                                 | Contributions                                                                                            |
| :----------------------------------- | :------------------------------------------------------------------------------------------------------- |
| [SwitchUpCB](https://switchupcb.com) | Project Architecture, [Discord API Spec (Disgo Structs)](https://github.com/switchupcb/discord-api-spec) |

_Earn a credit! [Contribute Now](contribution/CONTRIBUTING.md)._