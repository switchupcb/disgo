# What is an Event?

An event is an action that happens. A bot receives a Discord **Event** by connecting to a Discord WebSocket (Gateway). For example, when a message is created by a User, Discord sends that (Message Create) event to the bot.

## What is an Event Handler?

An event can occur at any moment. Event handling is an asynchronous operation where an **event listener** waits for an event to occur, while an **event handler** handles the respective incoming event. For example, you can implement an event handler to determine what happens when a message is created by a User.

# The Disgo Event Handler

Disgo provides a simple way to handle events in a Discord Bot. 

## How It Works

Opening a connection to a Discord WebSocket Session (Gateway) allows Discord to send the bot [**Events**](https://discord.com/developers/docs/topics/threads#gateway-events). When an event is sent to a bot's Session, Disgo's **event listener** passes the incoming event to the bot's `Client.Handlers`. Each **event handler** is called on a goroutine _(separate thread)_ which prevents your bot from being blocked while receiving more events.

### What is a Gateway Intent?

[Gateway Intents](https://discord.com/developers/docs/topics/gateway#gateway-intents) are required to receive certain events. Disgo makes managing a bot's Gateway Intents easy by **automatically** setting the `Client.Config.Gateway.Intents` when an event handler is added to the bot using the `Handle(event, handler)` function. When a bot's Session connects to the Discord Gateway, the bot's current `Intents` value will be used to identify which events to receive.

As a reminder, Disgo already provides **Automatic Intent Calculation**. However, intents can be managed from the `Client.Config.Gateway` using the `Gateway.EnableIntent(intent)` and `Gateway.Disable(intent)` functions. Intents are added using a [Bitwise OR operation](https://en.wikipedia.org/wiki/Bitwise_operation) which is a DESTRUCTIVE operation. As a result, an intent that is added to a bot can't be removed using `Gateway.Disable(intent)`. Instead, the `Gateway.Intents` value must be reset.

**_[Privileged Intents](https://discord.com/developers/docs/topics/gateway#privileged-intents) must be added using `EnableIntent` or `EnableIntentsPrivileged`._**

### When should I add or remove my event handler?

Event handlers can be added or removed from the application at any time. In contrast to application commands, event handlers are **NOT** maintained by Discord. However, this also means that event handlers do **NOT** persist when your bot restarts. Event handlers are invoked when a respective event is received from a connected Websocket Session. Keep in mind that the events your bot receives are dependent on its Intents at the start of the connection.

In order to add an event handler, use the `Client.Handle(event, handler)` function. 

```go
// Add an event handler to the bot.
err := bot.Handle(disgo.FlagGatewayEventNameInteractionCreate, func(i disgo.InteractionCreate) {
	log.Printf("InteractionCreate event from %s", i.User.Username)
})

// It's recommended to check the error of the Event Handler functions Handle() and Remove().
// 	Handle() will fail when the (eventname, function) parameters are not configured correctly.
// 	Remove() will fail when there is no event handler to remove at the given index.
if err != nil {
	log.Printf("Failed to add event handler to bot: %v", err)
}
```

In order to remove an event handler, use the `Client.Handlers.Remove(event, index)` function.

```go
// Remove the first InteractionCreate event handler from the bot.
bot.Handlers.Remove(disgo.FlagGatewayEventNameInteractionCreate, 0)
```
