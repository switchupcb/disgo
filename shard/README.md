# Disgo Shard Manager

The Disgo Shard Manager is a Go module that automatically handles sharding for your Discord Bot. The Disgo Shard Manager works by scaling your WebSocket Sessions through the `disgo.Client.Sessions` field. The `Client` calls `GET /gateway/bot` to retrieve the recommended number of shards to use upon joining a `variable` amount of guilds. These shards are assigned to a **Session** using the recommended number of shards by [Discord](https://discord.com/developers/docs/topics/gateway#get-gateway-bot). For more information about the concept of sharding, read [What is a Discord Shard?](/_contribution/concepts/SHARD.md).

## Implementation

Sharding is a two-step process that involves implementing shard-logic in your application and sharding your infrastructure _(optional)_.

### Sharding the Bot

Implement shard-logic in your application by importing the **Disgo Shard Manager** module, then setting the `Client.Config.Gateway.ShardManager = new(ActiveShardManager)`. This manager will set the `Shard` field of the `Identify` payloads that are sent to Discord, which indicates that your bot is accepting data from multiple shards. **Technically, this is all that's required to implement sharding.** The purpose of Discord's sharding requirement is to minimize the amount of data that Discord sends per WebSocket Session. There is nothing stopping you from running one server that creates multiple sessions and handles them in one application.

### Sharding the Infrastructure

Discord only allows you to shard by guild, so — barring a load-balanced architecture — you must handle _every_ important Discord Event in every Discord Bot application within a single codebase. As a result, the only way to shard the infrastructure of your Discord Bot application _— without requiring additional code  —_ is to host multiple copies of it _(each handling a fraction of your total load)_. This is known as **active-active load balancing**.

_The Disgo Shard Manager implements this approach. Read the following guide for an alternative._

# Guide

## Terminology

### What is an application?

An application refers to a built binary that is executed on your computer. You run applications through the terminal (i.e `run.exe`) or via visual shortcuts. In the context of this explanation, an application refers to the code that your Discord Bot uses to run on a **server**.

### What is a server?

A server is a computer (with a specialized use-case). You run applications on computers. A discord bot application is hosted _(ran)_ on a server.

### What is a guild?

Guilds in Discord represent an isolated collection of users and channels: These are often referred to as "servers" in the User Interface (UI). However, these "servers" are **NOT** the same as the **servers** described above. A Discord guild is a concept, while a **server** is a physical machine.

## Approach

A WebSocket Session contains **shards** that manage **multiple guilds**. Ignoring specific events across one shard would ignore those same events from _multiple guilds_ the shard manages: A guild's shard isn't able to be specified directly. Therefore, ignoring a **session's events** ignores **multiple shards' events** which ignores **multiple guilds' events**. Since Discord requires you to shard by guild, you **CANNOT** shard the infrastructure of a Discord Bot by creating multiple applications that handle a single event _(without an alternative infrastructure)_.

### Alternative

A **load balancer** allows you to "shard" your bot by event. This entails creating **one application** _(the load balancer)_ that accepts every event your bot receives _(and thus every shard)_, then having that **same application** forward those events to micro-applications _(which run on other servers)_. This strategy is also known as a [microservice architecture](https://en.wikipedia.org/wiki/Microservices). Placing every "shard" in a single application requires _that same application_ to maintain every session. As a result, your load balancer's only purpose — in this alternative infrastructure — is to balance the load by routing events to other applications. When the load balancer _(that handles every session)_ goes down, so does your bot.

## QA

### When do I need to shard?

Discord requires you to shard once you've reached a [certain number](https://discord.com/developers/docs/topics/gateway#sharding) of guilds.

### What are the implications of using one server?

Servers are computers with **CPU**, **RAM**, and **Storage**. You typically run one application on a server because you expect that application to use **all** of the server's resources _(i.e 100% CPU, 100% RAM, etc)_. Placing multiple applications on one server is only useful when your application does **NOT** use all of the server's resources, cores, etc. This implies that your application handles a low amount of data, experiences a bottleneck _(i.e waiting on a network request)_, and/or maintains a consistent load. If a server with two cores — without any form of multithreading — has an application using _<100% CPU_ on one core, then you can add an additional application _(that uses the other core)_ to the server without a performance hit. 

_In practice, it's **NOT** this straightforward._

If you need to shard your bot efficiently, you _probably_ need to use multiple servers with multiple applications that all represent your "Discord Bot" as a single entity: This entity — containing multiple servers — is known as a [cluster](https://en.wikipedia.org/wiki/Computer_cluster). Each servers' application(s) would accept a separate amount of shards and process the shard's data accordingly. These applications **CAN** be built from the same codebase that was used prior to sharding, but will require modification if the bot implements cross-guild functionality. Otherwise, in most cases, all you need to do is implement this module in your application.