# What is a Request?

A request is an act of communication. Suppose that we have a conversation:
- I send you a message.
- You receive the message.
- You process the message.
- You send a new message back to me.

This conversation occurs similarly between our computers. Whenever you [enter a URL in your browser](https://github.com/alex/what-happens-when#browser), a client (e.g., web browser) sends a request to a web server that processes the request, then sends another back.

## What is a REST HTTP Request?

To communicate better, humans create protocols that others adhere to during conversation:
- The **HTTP** protocol is used to send and receive resources _(data)_. 
- A **REST API** is a style used to communicate resources.

Discord uses an **HTTP REST API** to transfer information between its servers and your client (e.g., Discord Bot).

### Client vs. Server

In the context of a network, a **client** is a computer that receives information, while a **server** is a computer that serves information. This distinction can be confusing because a computer with a specialized use case is colloquially referred to as a server.

In the context of a Discord Bot, the **client** refers to the server (computer) that your bot application runs on, while the **server** refers to Discord's server(computer).

# Disgo Requests

Disgo provides a simple way to send requests from a Discord Bot.

## How does it work?

Suppose that your Discord Bot sends a request to Discord:
1. Disgo sends a request to Discord's API Server(s). 
2. Discord's API Server processes the sent **request** based on several factors _(e.g., headers, endpoint, URL query string, data)_.
3. Discord's API Server returns a **response** with a status code and other data to the Discord Bot.
4. Disgo handles the response to provide you with the requested data.

_Disgo automatically handles request rate limits so your bot isn't blacklisted from Discord._

## How do I send a Request?

A request is sent using the `Send(bot)` function. 

Disgo is a 1:1 API, meaning the objects defined in the [Discord API Documentation](https://discord.com/developers/docs/intro) are **directly** represented in Disgo. For example, a [`CreateGlobalApplicationCommand`](https://discord.com/developers/docs/interactions/application-commands#create-global-application-command) request can be prepared and sent using the following code:

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

### What is a Request Retry?

A request retry occurs when your request fails to receive a response (from Discord) due to an error. You can set the amount of retries per request by setting the `Client.Config.Request.Retries` field _(default: 1)_.

```go
bot.Config.Request.Retries = 1
```

### What is a Request Timeout?

A request timeout represents the amount of time a request will wait for a response (from Discord). You can set a bot's request timeout from the `Client.Config.Request.Timeout` field _(default: 1s)_. 

```go
bot.Config.Request.Timeout = time.Second * 15
```

_`fasthttp.ErrTimeout`  is returned from timed out requests._

## What is a Rate Limit?

Servers use rate limits to prevent spam, abuse, and service overload. A rate limit defines the speed at which a server can handle requests _(in requests per second)_. 


While there are many rate limit strategies a server may employ, [Google Architecture Rate Limiting Strategies](https://cloud.google.com/architecture/rate-limiting-strategies-techniques#techniques-enforcing-rate-limits) explains the most common cases. Discord enforces multiple rate limit strategies depending on the data sent to the server. 

The Discord rate limit strategies for requests include:
- [Global (Requests)](https://discord.com/developers/docs/topics/rate-limits#global-rate-limit)
- [Per Route (Requests)](https://discord.com/developers/docs/topics/rate-limits#rate-limits)
- Per Resource (Requests)

Disgo makes adhering to Discord's Rate Limits easy by providing a customizable rate limiter:
- Use the builtin [`RateLimit`](/wrapper/ratelimit.go) implementation or develop your own by implementing the [`RateLimiter interface`](/wrapper/ratelimiter.go) _(which stores Buckets)_.
- Set the `Client.Request.RateLimiter` to customize how rate limiting works for HTTP Requests.
- Set entries in the `RateLimitHashFuncs` map to control how a route is rate limited _(per-route, per-resource, etc)_.
- Configure the `RateLimit.DefaultBucket` to control the behavior for requests that are sent without a known rate limit.
  
### Per-Route vs. Per-Resource

Discord maintains two main rate limit strategies: Per-Route and Per-Resource.
 
A **per-route (user) rate limit** refers to a rate limit that applies at the route level and for the user. Specifically, a **route** is a combination of an HTTP Method and Endpoint _(e.g., `GET /guilds/id`)_. 

**Users and Bots must adhere to per-route rate limits.** 

A **per-resource (shared) rate limit** refers to a rate limit that applies at the resource level and for a resource _(e.g., `guild`)_. It's not recommended to adhere to per-resource rate limits.

_BOTH rate limits can be applied to the same route._

### Disclaimer

Discord expects the bot (user) to send a request until it succeeds _(as many times as necessary)_. So the per-route (user) rate limit is used to limit the number of requests _a user_ sends per second, while the per-resource (shared) rate limit is used to limit the usage of a resource _(to control the overall load on Discord's servers)_. 

Per-resource routes depend on factors your bot can **NOT** keep track of. So the bot is only required to adhere to per-route rate limits. In addition, experiencing `429 Status Codes` with the `shared` Rate Limit Scope Header does **NOT** count against you.

Disgo helps the developer implement the above behavior through the `Request.RetryShared` field. When the `Client.Request.RetryShared` field of a bot is set to `true` _(default)_, the bot will send a request — within the per-route rate limit — until one is successful or until one experiences a non-shared 429 status code.

In any other case, the `Request.Retries` field can be set to control the number of times a request may be retried upon any failure. Implementing per-user per-resource route rate limits is possible using the `RateLimitHashFuncs` map _(see example)_, but not recommended.

## What is a Default Bucket?

Discord utilizes a Token Bucket Rate Limit Algorithm for their rate limits. Unfortunately, Discord's specific implementation of this rate limit strategy does **NOT** allow the application to determine the rate limit of a **route** (HTTP Method + Endpoint) until a request with that **route** is sent. This results in a dilemma where you must determine whether to sacrifice performance or safety to send specific requests _(before those requests have ever been sent)_.

A **Default Bucket** is used when a rate limit is **NOT** yet known by the application: In other words, when a request for a route has **NEVER** been sent. 

In Disgo, the `RateLimit.DefaultBucket` field represents the Default Rate Limit Bucket used for requests which operate at the per-route level. However, configuring Default Buckets for per-resource (n) routes is also possible _(see example)_. 

Set the `DefaultBucket.Limit` field-value to control how many requests of a given route can be sent _(per second)_ **BEFORE** the actual Rate Limit Bucket of that route is known.


## Example

Suppose that **Route A** (`POST /A`) is constrained by a Global Rate Limit of 50 requests per second and a Per Route Rate Limit of 25 requests per second. When a bot (application) starts, it opts to send these requests as often as needed. 

No problems occur while the bot is small since it always sends less than 25 requests per second. However, the bot grows and eventually receives `429 Too Many Request` responses upon startup.

_What went wrong?_

Instead of sending <=25 requests upon startup, the bot sends 40. While this adheres to the Global Rate Limit, it does **NOT** adhere to the Per Route Rate Limit.

**There are two valid solutions to the above issue:**

1. Send a request synchronously _(blocking; in-order)_ until the rate limit is known. In other words, require the bot to send one request and wait _until the response is received_ to be able to send requests concurrently _(non-blocking; asynchronous)_.
2. Send as many requests as needed (while adhering to the Global Rate Limit), and resend the requests for responses that receive `429 Too Many Request Status Codes`.

In either case, the bot will eventually successfully send all the required requests. However, the first case will take 3 batches (1 + 25 + 14), while the second case will only take 2 batches (25 + 15); at the cost of 15 `429 Status Codes`.

Employing the second strategy is more efficient but could be costly. In an actual application, there are other implications to failed requests that we haven't even considered.

_Receiving 10,000 `(user) 429 Status Codes` in 10 minutes results in a [Cloudflare Ban](https://discord.com/developers/docs/topics/rate-limits) for approximately one hour._ 

### Solution

Disgo solves the problem described in the above example using **configurable Rate Limits and Default Buckets**:
- When a request's rate limit is unknown, Disgo will only send as many requests as the configured Default Bucket allows _(1 by default)_. 
- Once the request receives a response, the Default Bucket will be discarded and replaced by the request's actual Rate Limit Bucket _(or nil)_.
  
This implementation gives you multiple ways to address the issue described above.

#### Configuring the Route Default Bucket

If you want to ensure that every request at the **route** level initially sends 25 requests per second, you can set the `DefaultBucket` of the `Client.RateLimiter`.

```go
bot.Config.Request.RateLimiter.SetDefaultBucket(&disgo.Bucket{
	Limit:     25,
	Remaining: 25,
})
```

#### Configuring the "Route A" Default Bucket

If you want to ensure that **ONLY** Route `A` initially sends 25 requests per second, you can initialize a `RateLimiter` with that Route ID `Bucket`. 

```go
// create a Client using a Default Configuration.
bot := disgo.Client{
    Config: disgo.DefaultConfig()
    ...
}

// set Route A to a Bucket (ID "temp") in the Client's initialized Request Rate Limiter.
bot.Config.Request.RateLimiter.SetBucketID("A", "temp")
bot.Config.Request.RateLimiter.SetBucketFromID("temp", &disgo.Bucket{
	Limit:     25,
	Remaining: 25,
})

// optional: set other Routes (i.e "B") to Route A's Bucket using the SetBucketID function.
bot.Config.Request.RateLimiter.SetBucketID("B", "temp")
```

_NOTE: `"A"` is used as the ID for Route A in this example. Use the Route ID showcased in [`request_send.go`](/wrapper/request_send.go) for actual requests._

#### Configuring the Parent Default Bucket

When Route `A` refers to a per-resource route, a Default Bucket can be configured by using the `Route ID` of the parent route. As an example, Route `A/Guild2/Channel3` _(Route ID `A/Guild2`, ResourceID `Channel3`)_ uses the Default Bucket at `A/Guild2` _(if it exists)_. This is possible with two steps.

**1.** Configure the hashing function for Route `A/Guild2/Channel3` to use `A/Guild2` as a Route ID. **This step can also be used to change the rate limit algorithm of any route.**

```go
disgo.RateLimitHashFuncs[disgo.RouteIDs["A"]] = func(routeid string, parameters ...string) (string, string) {
    return routeid + parameters[0], parameters[1] // where parameters is [Guild2, Channel3]
}
```

**2.** Set the default bucket for Route `A/Guild2`.

```go
bot.Config.Request.RateLimiter.SetBucketID("AGuild2", "AGuild2BucketID")
bot.Config.Request.RateLimiter.SetBucketFromID("AGuild2BucketID", &disgo.Bucket{
	Limit:     25,
	Remaining: 25,
})
```

This results in the first requests of Route `A/Guild2/Channel3`, Route `A/Guild2/Channel4`, Route  `A/Guild2/Channel...` to initially use a Rate Limit Bucket that allows 25 requests per second.

#### Configuring Both Buckets

When you configure both buckets _(Route, Parent Route, and Route ID)_, Route `A` is **ONLY** assigned to the `Route ID` Default Bucket, since Route `A` already has a "known" bucket. This known bucket will be updated upon receiving a response — that results from a Route `A` request — from Discord.
