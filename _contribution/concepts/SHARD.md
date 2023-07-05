# What is Discord Sharding?

[Sharding (computing)](https://en.wikipedia.org/wiki/Shard_(database_architecture)) is the process of splitting data into segments (for scalability purposes).

**Discord Sharding** refers to the management of incoming event data using shards (traffic routes of guild event data) to scale your Discord Bot.

_Here is an in-depth explanation of how Discord Sharding works._

## What is Scale?

[Scale (computing)](https://en.wikipedia.org/wiki/Scalability) refers to the amount of work a system can handle. As your bot is added to new guilds, it must handle an increasing amount of data. 

_**Discord Sharding** is one way to handle this increasing load._

## What is a WebSocket?

A [WebSocket](https://en.wikipedia.org/wiki/WebSocket) is a communication protocol that lets clients and servers communicate in real time. The name references a physical [socket](https://en.wikipedia.org/wiki/Socket), which is a hole that an object is placed in.

### What is a WebSocket Session?

A **WebSocket Connection** refers to the state of communication (i.e., connected, disconnected) between a client (e.g., Discord Bot) and server (e.g., Discord Gateway). When a WebSocket Connection starts, data is transferred between the connected client and server.

A **WebSocket Session** refers to the period of time that a client and server are connected. The terms WebSocket Connection and WebSocket Session have different meanings in the Discord ecosystem.

Suppose that a Discord Bot connects to the Discord Gateway using a WebSocket:
- A Discord WebSocket Session starts when the client and server are connected.
- A Discord WebSocket Session ends when the client and server are disconnected.
- A **new** Discord WebSocket Session starts when the client and server are connected again _(resumed)_.

In contrast, the **Discord WebSocket Connection** refers to the state of communication between the Discord Bot's Sessions and the Discord Gateway.

_The Discord Gateway Rate Limit applies to the Discord WebSocket Connection. So the Discord Gateway Rate Limit is applied per Discord Bot token._

### When is a WebSocket Used?

An [**HTTP Request**](/_contribution/concepts/REQUESTS.md) is used when data is served upon a request from the client (e.g., webpage, Discord Bot). In contrast, a **WebSocket Connection** is used when the server (e.g., Discord Gateway) expects to send data (e.g., guild events) to the client (e.g., Discord Bot) without a request.

_Discord uses WebSocket Sessions to send real-time event data to Discord Bots (without requiring those bots to make multiple requests to Discord)._

## What is a Discord Shard?

A **Discord Shard** is an abstract concept that defines how guild event data is routed to a WebSocket Session.

Suppose that you create a Discord Bot:
- Without sharding, one Session communicates with the Discord Gateway for the event data of every guild the bot is in.
- You create a single shard that manages the event data of every guild the bot is in. This shard is routed to a single Session, so the session handles the event data of every guild the bot is in.

Suppose that your bot is added to 4,000 guilds:
- You create two shards, each managing the event data of 2000 guilds. These shards are routed to a single Session, so the session still handles the event data of every guild the bot is in.

Suppose that your bot is added to another 1,000 guilds:
- You create a third shard but use a passive sharding strategy instead of splitting the load among each shard equally. So shard 1 manages 2,000 guilds, shard 2 manages 2,000 guilds, and shard 3 manages 1,000 guilds. However, you route each shard to a single Session, so the session still handles the event data of every guild the bot is in.
- You create another **session**, route the third shard to it, then disconnect the third shard from the first session. So the first session manages the event data of 4,000 guilds, and the second session manages the event data of 1,000 guilds.
- You decide to route the third shard to the first session again. So the first session manages the event data of 5,000 guilds, and the second session manages the event data of 1,000 guilds. That said, your bot is only in 5,000 guilds.

**These examples illustrate the following hierarchy between a WebSocket Connection, WebSocket Session, and Discord Shard:**
- WebSocket Connection (Discord Bot) manages
  -  WebSocket `Session(s)` manages
     -  Discord `Shard(s)` defines groups of
        -  `Guild(s)` Event Data

At a certain point, handling all incoming Guild Event Data on one Session becomes impossible, so Discord requires you to handle the data through multiple **Sessions** using **Shards**.

### How Does It Work?

_For a technical explanation, read [Sharding on Discord](https://discord.com/developers/docs/topics/gateway#sharding)._


Discord Sharding is an implementation of a sharding strategy. The **Discord Sharding Formula** is used to determine which shard a guild is scoped to.

```
shard_id = (guild_id >> 22) % num_shards
```

The `shard_id` is equal to the following operation: The `guild_id` (which is an integer-based timestamp) `shifted 22 bits to the right` modulo `the number of shards`. 

Suppose that you are calculate the `shard_id` for `guild_id = 197038439483310086` for a single shard connection:
1. Using a [bit shift calculator](https://bit-calculator.com/bit-shift-calculator), `(197038439483310086 >> 22) = 46977624770`.
2. `(46977624770) % 1 = 0`
3. `shard_id = 0`.

A `number % 1` is always `0`: So every guild in a single-shard connection is routed to the shard at index 0. In practice, more than one shard is used to shard a Discord Bot.

_Direct Message data is always sent to shard 0._

### How Does It Impact Performance?

Maintaining multiple WebSocket Sessions does **NOT** have any performance implications on its own. However, processing more data among many connections (within a set period of time) warrants more processing power (i.e., higher CPU usage).

[Goroutines](https://gobyexample.com/goroutines) allow you to manage asynchronous connections concurrently (without blocking). This language-specific feature, in addition to other factors is why you should use Go to create a Discord Bot.

### How Do You Implement Discord Sharding?

Disgo makes implementing Sharding easy by providing a customizable shard manager. Use the [**Disgo Shard Manager (module)**](/shard/README.md) or develop your own by implementing the [`ShardManager interface`](/wrapper/shard.go).
