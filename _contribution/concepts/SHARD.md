# What is Discord Sharding?

Disgo provides an **optional** shard manager for your Discord Bot. For information about its implementation, read about [The Disgo Shard Manager](/shard/README.md).

## What is a Session?

A bot connects to Discord using a WebSocket. This **WebSocket Session** provides an ongoing open connection for us to receive or send data to Discord. In contrast, an **HTTP Request** makes a single request for data to or from discord. A **new bot** uses one session and/or multiple requests to handle data _(pertaining to guilds, users, etc)_ from Discord.

## What is Scale?

As the **new bot** is used by more guilds, it handles more and more data. This growth requires the bot to work at-scale _(with a large amount of data)_. 

## What is a Shard?

When the **new bot** handles more data _(without sharding)_, all that data is sent or received through one **Session**. This becomes inefficient, so Discord requires you to handle the data through multiple **Sessions**: This lowers the amount of data each **Session** maintains. A WebSocket **Session** (also known as a Gateway) holds a number of **shards**, and each **shard** contains multiple **guilds**.

### Advanced

A shard only refers to which **Session(s)** a guild's data will be sent to. In other words, a shard represents a _traffic route_ for guild data. This presents the following hierarchy: 
- **One Session manages multiple shards**. 
- **One Shard manages multiple guilds' data**.

Multiple sessions can use the same shard _(routes for a guild's data)_; each session can contain a different number of shards. This allows you to create sessions that handle a different amount traffic _(which is done for a multitude of reasons)_.

## What is Sharding?

Sharding refers to managing data through multiple Shards _(WebSocket Sessions, Gateways, Whatever)_. Discord shards bots by guild: **One Shard manages the data for multiple guilds.**

### How Does It Work?

Read [Sharding on Discord](https://discord.com/developers/docs/topics/gateway#sharding) for a technical explanation. Sharding works by using a **Sharding Formula** to assign _(the data from)_ guilds to a specific shard.

```
shard_id = (guild_id >> 22) % num_shards
```

The `shard_id` is equal to: The `guild_id` (which is an integer-based timestamp) `shifted 22 bits to the right` modulo `the number of shards`. Using a [bit shift calculator](https://bit-calculator.com/bit-shift-calculator), we can see that a `guild_id = 197038439483310086` will result in a `shard_id = (0 % 1) = 0`. This means that all data for the guild in a single-shard connection will go to the shard at index 0. In reality, more than one shard is used while sharding. _Note that all Direct Message data is sent to shard 0._

### How Does It Impact Performance?

Maintaining multiple WebSocket Sessions does **NOT** have any performance implications on its own. However, processing more data among many connections _(within a set period of time)_ warrants more processing power (resulting in higher CPU usage). [Goroutines](https://gobyexample.com/goroutines) allow you to manage asynchronous connections concurrently. This language-specific feature — in addition to other factors — is why Go is the best language for creating Discord Bots.
