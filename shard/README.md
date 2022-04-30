# Disgo Shard Manager

The Disgo Shard Manager is a Go module that automatically handles sharding for your Discord Bot. The Disgo Shard Manager works by scaling your WebSocket Sessions through a `disgo.Client.Sessions` array. The `Client` calls `GET /gateway/bot` to retrieve the recommended number of shards to use upon joining a `dynamic number` of guilds. These shards are assigned to a **Session** until the [max concurrency](https://discord.com/developers/docs/topics/gateway#sharding-max-concurrency) limit is reached. For information on the concept of a shard, read [What is a Discord Shard?](contribution/concepts/SHARD.md)

## Implementation

### Terminology

#### What is an application?

An application refers to a built binary that is executed on your computer. You run applications through the terminal (i.e `run.exe`) or via visual shortcuts. In the context of this explanation, an application refers to the code your Discord Bot uses to run on a **server**.

#### What is a server?

A server is a computer (with a specialized use-case). You run applications on computers. A discord bot is hosted _(ran)_ on a server.

#### What is a guild?

Guilds in Discord represent an isolated collection of users and channels, and are often referred to as "servers" in the UI. However, these "servers" are **NOT** the same as the **servers** described above. A Discord guild is a concept, while a **server** is a physical machine.

### Method

Sharding on Discord is a two-step process: implementing shard-logic and sharding your infrastructure. The **Disgo Shard Manager** handles the first step for you, so **how do we shard our infrastructure?** Discord only allows you to shard by guild, so - barring a load-balanced architecture - you must handle _every_ important Discord Event in every Discord Bot application. As a result, the only way to shard the infrastructure of your Discord Bot application _— without requiring additional code  —_ is to host multiple copies of it _(each handling a fraction of your total load)_.

#### Explanation

A WebSocket Session contains **shards** that contain **individual guilds**. Ignoring specific events in one shard would ignore those same events from multiple guilds _(since you can't specify a guild's shard)_. Following this logic, ignoring a **session's events** ignores **multiple shards' events** which ignores **multiple guilds' events**. Since Discord requires you to shard by guild, we **CANNOT** shard the infrastructure of our Discord Bot by creating multiple applications that handle a single event. Doing so would result in only receiving specific events from certain guilds _(depending on their session)_.

#### Alternative

A load balancer allows us to "shard" our bot by event. This entails creating **one application** _(the load balancer)_ that accepts every event your bot receives _(and thus every shard)_, and having that **same application** forward those events to micro-applications _(which run on other servers)_. Placing every "shard" in a single application requires that application to maintain every session. As a result, your load balancer's only purpose should be to balance the load by routing events to other applications. When the load balancer _(that handles every session)_ goes down, so does your bot.

### QA

#### When do I need to shard?

Discord requires you to shard once you've reached a certain number of guilds.

#### What do I need to do?

Discord requires that you implement sharding logic in your bot, which is what this module - the **Disgo Shard Manager** - does.

#### Is that all?

Technically, **yes**. The purpose of Discord's sharding requirement is to minimize the amount of data they send per WebSocket Session on their end. There is nothing stopping you from using one server that creates multiple sessions and handles them in one application.

#### What are the implications of using one server?

Servers are computers with **CPU**, **RAM**, and **Storage**. The reason that you typically run one application on a server is because you expect that application to use **all** of the server's resources _(i.e 100% CPU, 100% RAM, etc)_. Placing multiple applications on one server is only useful if your application does **NOT** use all of the server's resources, cores, etc. This would imply that your application is handling a low amount of data or experiencing a bottleneck _(i.e waiting on a network request)_. If a server with two cores — without any form of multithreading — has an application using _<100% CPU_ on one core, then you would be able to add an additional application _(that uses the other core)_ to the server without a performance hit. 

_In practice, it's **NOT** this straightforward._ If you need to shard your bot efficiently, you _probably_ need to use multiple servers with multiple applications that all represent your "Discord Bot". Each servers' application(s) would accept a separate amount of shards and process the shard's data accordingly. These applications **CAN** be built from the same codebase used prior to sharding, but likely require modification if the bot implements cross-guild functionality. However, all you will need to do in most cases is implement this module in your application.