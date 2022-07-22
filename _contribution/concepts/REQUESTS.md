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
// Create a Create Global Application Command request.
request := disgo.CreateGlobalApplicationCommand{
    Name: "main",
    Description: "A basic command",
}

// Send the Global Application Command to Discord using the bot.
newCommand, err := request.Send(bot)
if err != nil {
    log.Printf("failure sending command to Discord: %v", err)
}
```

### What is a Rate Limit?

Rate limits are used by servers to prevent spam, abuse, and service overload. A rate limit defines the speed at which a server can handle requests _(in requests per second)_. While there are many rate limit strategies a server may employ, [Google Architecture Rate Limiting Strategies](https://cloud.google.com/architecture/rate-limiting-strategies-techniques#techniques-enforcing-rate-limits) provides an explanation for the most common ones. Discord enforces multiple rate limit strategies dependent on the data that is sent to the server. This includes Global (Requests), Per Route (Requests), Per Resource (Requests), Per Resource Per Routes (Emoji), Global (Gateway), and Identify (Gateway) rate limits.

Disgo makes adhering to Discord's Rate Limits easy by providing a customizable rate limiter. A [`Bucket`](../../wrapper/ratelimiter.go) represents a single rate limit. Use the [`RateLimiter interface`](../../wrapper/ratelimiter.go) to provide your own Rate Limiter implementation _(which stores Buckets)_. Set the `Client.Request.RateLimiter` or `Client.Gateway.RateLimiter` to customize how rate limiting works for HTTP Requests and Gateway Events. Configure the entries in the `DefaultBuckets` map to control the behavior for requests sent without a known rate limit.

#### What is a Default Bucket?

Discord utilizes a Token Bucket Rate Limit Algorithm for Per Route, Per Resource, and Per Resource per Route (Emoji) requests. Unfortunately, Discord's specific implementation of this rate limit strategy does **NOT** allow the application to determine the rate limit of certain **Routes** (HTTP METHOD + Endpoint) until a request with that route is sent. In other cases, the rate limit cannot be known until a request for a specific **Resource** _(i.e `guild`,`channel`,`webhook`)_ on a route is sent. This results in a dilemma where one must determine whether to sacrifice performance or safety to send certain requests _(before those requests have ever been sent)_.

A Default Bucket is used when a Rate Limit is **NOT** yet known by the application. In other words, when a request for a route has **NEVER** been sent _(since the start of the application)_. In Disgo, the `DefaultBuckets` global variable represents a map of the Default Rate Limit Buckets used for each rate limit strategy. Set the Default `Bucket.Limit` field-value to control how many requests of a given route can be sent _(per second)_ **BEFORE** the actual rate limit Bucket of the route is known.

##### Example

In the following example, **Route A** (`POST /A`) is constrained by a Global Rate Limit of 50 requests per second, and a Per Route Rate Limit of 25 requests per second. When an application _(bot)_ is started, it opts to send this request as many times as it needs. No problems occur while the bot is small, since it always sends less than 25 requests per second. However, the bot eventually begins receiving `429 Too Many Request` responses upon startup. **What is going wrong?** Instead of sending <=25 requests upon startup, the bot is sending 40. While this adheres to the Global Rate Limit, it does **NOT** adhere to the Per Route Rate Limit.

There are two valid solutions to the above issue. 

1. Send a request synchronously _(blocking; in-order)_ until the rate limit is known. In other words, require the bot to send one request and wait _until the response is received_ to be able to send requests concurrently _(non-blocking; asynchronous)_.
2. Send as many requests as needed (while adhering to the Global Rate Limit), and resend the requests for responses that receive `429 Too Many Request Status Codes`.

In either case, the bot will eventually end up successfully sending all the requests it requires. However, the first case will take 3 batches (1 + 25 + 14) while the second case will only take 2 batches (25 + 15); at the cost of 15 `429 Status Codes`. **Receiving 10,000 `429 Status Codes` in 10 minutes results in a [Cloudflare Ban](https://discord.com/developers/docs/topics/rate-limits)** for approximately one hour. As a result, employing the second strategy is more efficient, but could be costly. In an actual application, there are other implications to failed requests that we haven't even considered.

##### Solution

Disgo solves this problem through the use of a configurable Default Bucket. When a request's rate limit is unknown, Disgo will only send as many requests as the configured Default Bucket allows _(which is 1 by default)_. Once the request receives a respective response, the Default Bucket will be discarded and replaced by the request's actual rate limit _(if applicable)_. This implementation gives you two ways to address the issue described above.

**Configuring The Route A Bucket**

If you want to ensure that **ONLY** Route **A** sends 25 requests per second initially, you can initialize a  `RateLimiter` with that `Bucket`, then assign the initialized rate limiter to the `Client`. 

```go
// create a Client using a Default Configuration.
bot := disgo.Client{
    Config: disgo.DefaultConfig()
    ...
}

// add a Bucket to Route A in the Client's initialized Request Rate Limiter.
bot.Config.Request.RateLimiter.SetBucket("A", &Bucket{Limit: 25})
```

_NOTE: `"A"` is used to represent Route A in this example. Use the Route ID showcased in [`request_send.go`](../../wrapper/request_send.go) for actual requests._

**Configuring The Route Default Bucket**

If you want to ensure that every request with a **Route** Rate Limit sends 25 requests per second initially, you can set the `Route` key of the `DefaultBuckets` map.

```go
DefaultBuckets[DefaultBucketKeyRoute] = &Bucket{Limit: 25}
```

**Configuring Both**

If you do both, Route A will never be assigned a Default Bucket, since it already has a "known" bucket. This bucket will be updated upon receiving a response _(respective to a Route A request)_ from Discord.
