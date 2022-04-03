# Main

This example creates a bot that creates an application command and handles it.

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

// Define an event handler.
handler := disgo.EventHandler{
    Event: disgo.EventInteractionCreate,
    Call: func (i disgo.ResourceInteraction) {
        log.Printf("main called by %s", i.User.Username)
    },
}

// Add the event handler to the session.
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

A user creates an [`InteractionCreate`](https://discord.com/developers/docs/topics/gateway#commands-and-events-gateway-events) event by using `/main` in a direct message with the bot.

```
main called by SCB
```