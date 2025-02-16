# What is an Event?

An event is an action.

A Discord Bot receives a Discord **Event** by connecting to a Discord WebSocket (Gateway). 

_For example, when a User creates a message, Discord sends that (Message Create) event to the Bot._

## What is an Event Handler?

Event handling is an asynchronous operation where an **event listener** waits for an event to occur while an **event handler** handles the respective incoming event. 

_For example, you can implement an event handler to determine what happens when a User creates a message._

# The Disgo Event Handler

You can handle handle events in a Discord Bot with Disgo. 

## How It Works

1. Opening a connection to a Discord WebSocket Session (Gateway) allows Discord to send the bot [**Events**](https://discord.com/developers/docs/topics/threads#gateway-events).
2. An event is sent to a bot's session
   1.  The bot's **event listener** passes the incoming event to the bot's `Client.Handlers`.
   2.  Each **event handler** is called on a goroutine _(separate thread)_, which prevents your bot from being blocked while receiving more events.

### No Reflection

The performance of your Discord Bot increases when it does not use reflection or type assertion. 

 **Disgo does NOT use reflection or type assertion to handle events.** Instead, payloads from the Discord Gateway are marshalled into their respective structs directly.

Other API Wrappers use reflection and type assertion to convert the _Payload Data_ sent by the Discord Gateway into _Go event objects_. These operations negatively effect the performance of the entire application. 


### Selective Event Processing

The performance of your Discord Bot increases when it does not process unhandled events.

**Disgo does NOT process events your bot doesn't handle**.  So, your bot only processes handled events.

The Discord Gateway sends every event that is applicable to your bot. Other API Wrappers process these events _even when_ the bot does **NOT** need them. These operations negatively effect the performance of the entire application. 

### What is a Gateway Intent?

[Gateway Intents](https://discord.com/developers/docs/topics/gateway#gateway-intents) are required to receive certain events.

#### How do you manage Gateway Intents with Disgo?

Disgo **automatically** manages a Bot's Gateway Intents by setting the `Client.Config.Gateway.Intents` when an event handler is added to the Bot using the `Handle(event, handler)` function. The Bot's current `Intents` value will be used to identify which events to receive when a Bot's Session connects to the Discord Gateway.

You can also **manually* manage a bot's intents from the `Client.Config.Gateway` using the `Gateway.EnableIntent(intent)` and `Gateway.Disable(intent)` functions. 

#### How do Gateway Intents work?

Intents are added using a [Bitwise OR operation](https://en.wikipedia.org/wiki/Bitwise_operation), which is a DESTRUCTIVE operation. So, an intent that is added to a Bot can't be removed using `Gateway.Disable(intent)`. Instead, the `Gateway.Intents` value must be reset.

**_[Privileged Intents](https://discord.com/developers/docs/topics/gateway#privileged-intents) must be added using `EnableIntent` or `EnableIntentsPrivileged`._**

### When should you add or remove an event handler?

Disgo Event handlers can be added or removed from the application at any time. 

Event handlers are **NOT** maintained by Discord unlike application commands. So, event handlers do **NOT** persist when your bot restarts.

Event handlers are invoked when a respective event is received from a connected WebSocket Session. The information your bot receives from each event depends on the `Client.Config.Gateway.Intents` value at the start of the connection.

In order to add an event handler, use the `Client.Handle(event, handler)` function.

```go
// Add an event handler to the bot.
err := bot.Handle(disgo.FlagGatewayEventNameInteractionCreate, func(i disgo.InteractionCreate) {
	log.Printf("InteractionCreate event from %s", i.User.Username)
})

// You are recommended to check the error of the Event Handler functions Handle() and Remove():
// 		Handle() fails when the (eventname, function) parameters are not configured correctly.
// 		Remove() fails when there is no event handler to remove at the given index.
if err != nil {
	log.Printf("Failed to add event handler to bot: %v", err)
}
```
Use the `Client.Handlers.Remove(event, index)` function to remove an event handler, 

```go
// Remove the first InteractionCreate event handler from the bot.
bot.Handlers.Remove(disgo.FlagGatewayEventNameInteractionCreate, 0)
```
