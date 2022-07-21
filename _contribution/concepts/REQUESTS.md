# What is a Request?

A request is an act of communication. Whenever we have a conversation, I tell you something, you receive it and process it, then tell me something back. This conversation occurs in a similar manner between our computers. Whenever you [enter a URL in your browser](https://github.com/alex/what-happens-when#browser), a request is sent to a server that processes it and sends you a request back.

## What is a REST HTTP Request?

In order to communicate better, we create protocols that we must adhere to during conversation. The HTTP protocol is used to send and receive resources _(data)_ while REST is a style of communicating those same resources. Discord uses an HTTP REST API to transfer information between its servers and your bot.

### Client vs. Server

In the context of a network, a **client** is a computer that _receives_ information while a **server** is a computer that _serves_ information. This can be confusing because a computer with a specialized use-case is colloquially referred to as a server. In the case of a Discord Bot, a client refers to the server (computer) that your bot runs on, while the server refers to Discord's servers.

# Disgo Requests

Disgo provides a simple way to send requests using a Discord Bot. 

## How It Works

An HTTP library is used to send Discord requests to a Discord API Server. This server will process the request based on a number of factors _(headers, endpoint, url query string, data, etc)_ and return with a status code and respective data. Disgo handles this information accordingly to provide you with the data you requested. In addition, Disgo automatically handles request rate limits so that your bot isn't blacklisted from Discord.

### How do I send a request?

Disgo is a 1:1 API which means that the objects defined in the [Discord API Documentation](https://discord.com/developers/docs/intro) are **directly** represented in Disgo. A request is sent using the `Send(bot)` function. For example, the `CreateGlobalApplicationCommand` can be prepared and sent using the following code.

```go
// Create a global command request.
request := disgo.CreateGlobalApplicationCommand{
    Name: "main",
    Description: "A basic command",
} 

// Send the global command request to Discord using the bot.
newCommand, err := request.Send(bot)
if err != nil {
    log.Println("error: failure sending command to Discord")
}
```