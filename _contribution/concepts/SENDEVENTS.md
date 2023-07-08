# What is a Send Event?

A [send event](https://discord.com/developers/docs/topics/gateway-events#send-events) is an [event](/_contribution/concepts/EVENTS.md) that your Discord Bot can send to the Discord Gateway.

# Disgo Send Events

Disgo provides a simple way to send events from a Discord Bot.

## How does it work?

Suppose that your Discord Bot sends a request to Discord:
1. Disgo sends an event to Discord's Gateway. 
2. Discord's Gateway Server processes the sent **event** based on several factors _(e.g., opcode, payload data)_.
3. Discord's Gateway Server returns a **payload** with an opcode and payload data to the Discord Bot.
4. Disgo handles the incoming event payload to provide you with the requested data.

_Disgo automatically handles send event rate limits so your bot isn't blacklisted from Discord._

## How do I send a Send Event?

A send event is sent using the `SendEvent(bot)` function.

Disgo is a 1:1 API, meaning the objects defined in the [Discord API Documentation](https://discord.com/developers/docs/intro) are **directly** represented in Disgo. For example, a [`RequestGuildMembers`](https://discord.com/developers/docs/topics/gateway-events#request-guild-members) can be prepared and sent using the following code:

```go
// Create a Request Guild Members send event.
sendevent := disgo.RequestGuildMembers{
    GuildID: "GUILD",
    Query:   disgo.Pointer(""),
    Limit:   disgo.Pointer(0),
}

// Send the Request Guild Members event to Discord using the bot and a connected session or shard manager.
if err := sendevent.SendEvent(bot, session); err != nil {
    log.Printf("failure sending SendEvent to Discord: %v", err)
}
```

## What is a Rate Limit?

_Read ["Requests: What is a Rate Limit?"](/_contribution/concepts/REQUESTS.md#what-is-a-rate-limit) for in-depth information about rate limits._

Servers use rate limits to prevent spam, abuse, and service overload. A rate limit defines the speed at which a server can handle events _(in requests per second)_.

The Discord rate limit strategies for send events include:
- [Global (Gateway)](https://discord.com/developers/docs/topics/gateway#rate-limiting) \[Per [WebSocket Connection](/_contribution/concepts/SHARD.md#what-is-a-websocket)\]
- [Identify (Gateway)](https://discord.com/developers/docs/topics/gateway#identifying)

Disgo makes adhering to Discord's Rate Limits easy by providing a customizable rate limiter:
- Use the builtin [`RateLimit`](/wrapper/ratelimit.go) implementation or develop your own by implementing the [`RateLimiter interface`](/wrapper/ratelimiter.go) _(which stores Buckets)_.
- Set the `Client.Gateway.RateLimiter` or `Session.RateLimiter` to customize how rate limiting works for Gateway Send Events.

### Global

The global rate limiter for the Discord Gateway is automatically configured when a `disgo.Session` connects to the Discord Gateway.

### Identify

The [`Identify`](https://discord.com/developers/docs/topics/gateway-events#identify) rate limit operates under the following conditions:
- Identify Send Event Rate Limit is applied per application (bot token).
- Identify Send Events count towards the Global Rate Limit for the Discord Gateway.
- Identify Send Event Rate Limit Buckets are reset every 5 seconds.

The bot rate limiter for Identify payloads is configured in the `Client.Gateway.RateLimiter`.

