# What is an Event Handler?

## What is an Event?

An event is an action that happens. A bot receives a Discord **Event** by connecting to a Discord WebSocket (Gateway). For example, when a message is created by a User, Discord sends that (message create) event to the bot.

## What is an Event Handler?

An event can occur at any moment. An **event listener** waits for an event to occur, while an **event handler** handles the respective incoming event. As a result, event handling is an asynchronous operation; run by the application.

# The Disgo Event Handler

Disgo provides a simple way to handle events in a Discord Bot.  

### How It Works

Opening a connection to a Discord WebSocket (Gateway) **Session** allows Discord to send the bot [**Events**](https://discord.com/developers/docs/topics/threads#gateway-events). When an event is added to a Session, Disgo makes a request for the event's required [Gateway Intents](https://discord.com/developers/docs/topics/gateway#gateway-intents); _if they're not already granted_. When an event is received on an `Open()` session, it's handled using a `Handler.Call` function. The `Add(handler)` function adds an event handler to a connection while `Remove(handler)` removes it.

### When should I add or remove my event handler?

Event handlers operate on your application and can be added or removed at any time. Unlike application commands, event handlers are **NOT** maintained by Discord. This means that they do **NOT** persist when your bot restarts. You must also keep in mind that a session only receives events when it's `Open()`. As a result, your handler will only run on an open **Session**.