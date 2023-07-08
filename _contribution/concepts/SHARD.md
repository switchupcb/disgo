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

A traditional **WebSocket Session** refers to the unique period of time that a client and server are connected. However, the terms WebSocket Connection and WebSocket Session have different meanings in the Discord ecosystem.

![Discord Sharding Diagram outlining difference between a WebSocket Connection, Discord Session, and Discord Shard.](/_contribution/concepts/_disgo-shard-diagram-min.jpg)

<br>

Suppose that a Discord Bot connects to the Discord Gateway using a WebSocket:
- **\[1\]** A WebSocket Connection is created when the client dials in to the server.
- **\[5\]** A Discord WebSocket Session is created when the client receives a Ready event from the server _(with a unique Session ID)_.
- **\[6\]** The Discord WebSocket Connection ends when the client is disconnected from the server.
- **\[6\]** The Discord WebSocket Session remains alive for 5 seconds when the client is disconnected from the server.
- **\[7\]** A **new** Discord WebSocket Connection is created when the client reconnects to the server.
- **\[7\]** The **old** Discord WebSocket Session is still alive when the client receives a Resumed event from the server _(within 5 seconds of the disconnection)_.
- **\[1\]** A **new** Discord WebSocket Connection is crated when the client reconnects to the server; regardless of how much time has passed since disconnection.
- **\[5\]** A **new** Discord WebSocket Session is created when the client receives a **new** Ready event from the server 5 seconds after disconnection.

The **Discord WebSocket Connection** refers to the state of communication between the Discord Bot and the Discord Gateway. 

In contrast, the **Discord WebSocket Session** refers to the unique period of time from \[1\] the initial Ready event to \[2\] the time that occurs 5 seconds after the Discord WebSocket Connection has been disconnected.

_The Discord Gateway Rate Limit is applied per Discord WebSocket Connection._

### When is a WebSocket Used?

An [**HTTP Request**](/_contribution/concepts/REQUESTS.md) is used when data is served upon a request from the client (e.g., webpage, Discord Bot). In contrast, a **WebSocket Connection** is used when the server (e.g., Discord Gateway) expects to send data (e.g., guild events) to the client (e.g., Discord Bot) without a request.

_Discord uses WebSocket Sessions to send real-time event data to Discord Bots (without requiring those bots to make multiple requests to Discord)._

## What is a Discord Shard?

A **Discord Shard** is an abstract concept that defines how guild event data is routed to a Discord WebSocket Session.

Suppose that you develop a Discord Bot:
- Without sharding, one Session communicates with the Discord Gateway for the event data of every guild the bot is in.
- You create a **single shard** that manages the event data of every guild the bot is in. This shard is routed to a single Session, so the session handles the event data of every guild the bot is in.

Suppose that your bot is added to 4,000 guilds:
- You create **two shards**, each managing the event data of 2000 guilds. However, a shard can only be routed to a single Session, so the bot ignores the event data of 2,000 guilds.
  
  **NOT GOOD!**

Suppose that your bot is added to another 2,000 guilds:
- You create a **third shard**. So shard 1 manages 2,000 guilds, shard 2 manages 2,000 guilds, and shard 3 manages 2,000 guilds. However, a shard can only be routed to a single Session, so the bot ignores the event data of shard 2 and 3 containing 4,000 guilds' event data.
  
  **NOT GOOD!**

- You create a **second session** and route the second shard to it. So the first session manages the event data of 2,000 guilds, and the second session manages the event data of 2,000 guilds. Therefore, the bot ignores the event data of 2,000 guilds.
  
  **NOT GOOD!**

- You create a **third session** and route the third shard to it. So the first session manages the event data of 2,000 guilds, and the second session manages the event data of 2,000 guilds, and the third session manages the event data of 2,000 guilds. 

  The bot handles the event data of every guild the bot is in.

  **GREAT!**

- You create a **fourth session** and route the first shard to it.

  The first session manages the event data of 2,000 guilds, the second session manages the event data of 2,000 guilds, the third session manages the event data of 2,000 guilds, and the fourth session _(spawned on an updated bot instance)_ manages the event data of 2,000 guilds; which is the same as the event data in shard 1.
  
  Your bot is still only in 5,000 guilds, but receives the event data from 2,000 guilds in the first shard twice: This lets you shut down your first session _(located on an instance containing outdated code)_ while still handling the first session's event data from your fourth session _(located on an instance containing updated code)_.

**These examples illustrate the following hierarchy between a WebSocket Connection, Discord Session, and Discord Shard:**
- `Discord Session` 
  -  involves `WebSocket Connection(s)` and Disconnection.
  -  can be tied to a `Discord Shard`
        - that defines the incoming guild event data (from multiple guilds) sent to the `Discord Session`.
  
At a certain point, handling all incoming guild event data on a single Session becomes impossible. So Discord requires you to handle this load with multiple **Discord Sessions**, each specifying a **Discord Shard**.

### How Does It Work?

_For a technical explanation, read [Sharding on Discord](https://discord.com/developers/docs/topics/gateway#sharding)._


Discord Sharding is an implementation of a sharding strategy. 

The **Discord Sharding Formula** is used to determine which shard a guild is scoped to.

```
shard_id = (guild_id >> 22) % num_shards
```

The `shard_id` is equal to the following operation: The `guild_id` (which is an integer-based timestamp) `shifted 22 bits to the right` modulo `the number of shards`. 

Suppose that you are calculate the `shard_id` for `guild_id = 197038439483310086` for a single shard connection:
1. Using a [bit shift calculator](https://bit-calculator.com/bit-shift-calculator), `(197038439483310086 >> 22) = 46977624770`.
2. `(46977624770) % 1 = 0`
3. `shard_id = 0`.

A `number % 1` is always `0`: So every guild in a single-shard connection is routed to the shard at index 0. However, in practice, more than one shard is used to shard a Discord Bot.

_Direct Message data is always sent to shard 0._

### How Does It Impact Performance?

Maintaining multiple WebSocket Sessions does **NOT** have any performance implications alone. However, processing more data among many connections (within a set time period) warrants more processing power (i.e., higher CPU usage).

### How Do You Implement Discord Sharding?

Disgo makes implementing Sharding easy by providing a customizable shard manager. Use the [**Disgo Shard Manager (module)**](/shard/README.md) or develop your own by implementing the [`ShardManager interface`](/wrapper/shard.go).
