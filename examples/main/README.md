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

Create an application command **resource** and a **request** to add an application command.

```go
// Create a global command resource.
newCommand := disgo.ResourceApplicationCommand{
    Name: "main",
    Description: "A basic command",
} 

// Create a global command registration request.
registeredCommand, err := disgo.RequestApplicationCommandAdd(newCommand)
if err != nil {
    log.Println("error: failure sending command to Discord")
}
```

## Handle an Event

Create an **event handler** and add it to a **session**.

```go
// Add a session.
bot.Sessions = append(bot.Sessions, disgo.Session{})

// Add a handler for an event to the session.
bot.Sessions[0].AddHandler(func(e disgo.EventInteractionCreate) {
    log.Println("/main called.")
})
```

### Output

Open the WebSocket **Session** to receive events.

```go
session, err := bot.Sessions[0].Open()
if err != nil {
    log.Println("error: can't open websocket session to Discord")
}
```

A user creates an interaction by using `/main` in a direct message..

[img]

