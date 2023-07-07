# Disgo Shard Manager

The Disgo Shard Manager is a Go module that automatically handles sharding for your Discord Bot. 

The Disgo Shard Manager works by managing the connection of multiple `disgo.Session` and setting the `Session.Shard` field:
1. The `Client` requests `GET /gateway/bot` to retrieve the [recommended number of shards by Discord](https://discord.com/developers/docs/topics/gateway#get-gateway-bot).
2. These shards _(defining traffic routes of guild event data)_ are assigned to a `disgo.Session`, which is then connected to Discord.

_For more information about the concept of sharding, read [What is a Discord Shard?](/_contribution/concepts/SHARD.md)._

## Implementation

Sharding is a two-step process that involves implementing shard-logic in your application and sharding your infrastructure _(optional)_.

### Import

Get a specific version of `shard` by specifying a tag or branch.

```
go get github.com/switchupcb/disgo/shard@v1.10.1
```

_Disgo branches are referenced by API version (i.e `v10`)._
### Sharding the Discord Bot

Change the instantiated `disgo.Session` variable to a `shard.InstanceShardManager`.

```go
// Change this line.
s := disgo.NewSession()

// To this line.
s := new(shard.InstanceShardManager)
```

**Technically, this change is all that's required to implement sharding.**

Discord's sharding requirement aims to minimize the amount of data that Discord sends per WebSocket Session. Nothing is stopping you from running a Discord Bot that creates multiple sessions and handles them in one instance.

### Sharding the Infrastructure

Discord doesn't let you select which shard a guild is defined on. This has implications on how you shard the infrastructure of a Discord Bot.

Ignoring a shard is equivalent to ignoring all incoming guild event data from that shard. So it's expected that you handle every event from a shard in a Discord Bot instance _(unless a load balancer is involved)_.

These constraints define the most straightforward sharding strategy:
1. Host multiple instances of your Discord Bot _(copies of a single codebase)_; each with the ability to handle all incoming events.
2. Host a central "Shard Manager instance" that each Discord Bot instance communicates to shard.

This sharding strategy is based on **active-active load balancing** and must be implemented using a modified shard manager.

_Read (TODO) "Implementing a Sharding Strategy (Guide)" for more information about alternative sharding strategy implementations._

# Guide

## Terminology

### What is an application?

An application refers to a built binary that is executed on your computer. You run applications through the terminal (i.e `run.exe`) or via visual shortcuts. 

_In the context of this explanation, an application refers to the code that your Discord Bot uses to run on a server._

### What is an instance?

An instance is a running application. 

_Instances can be run on a single server or spread among multiple servers._

### What is a server?

A server is a computer (with a specialized use-case). You run applications on computers. 

_A discord bot application is hosted (ran) on a server._

### What is a guild?

Guilds in Discord represent an isolated collection of users and channels: These are often referred to as "servers" in the User Interface (UI). However, these "servers" are **NOT** the same as the **servers** described above. 

_A Discord guild is a concept, while a server is a physical or virtual machine._

## Alternative Sharding Strategies

### Service Architecture (Event Architecture)

Implementing a service architecture in a Discord Bot is equivalent to handling incoming [events](https://discord.com/developers/docs/topics/gateway#gateway-events) using an instance (codebase) that solely handles specific types of events. 

This architecture involves the creation of "Shard Instances" which function as [API Gateways](https://learn.microsoft.com/en-us/azure/architecture/microservices/design/gateway) to forward incoming event data to respective "Event Handler Instances" (Services) which handle the events.

#### Shard Instance

Ignoring a shard is equivalent to ignoring all incoming guild event data from that shard. So it's expected that you handle every event from a shard in a Discord Bot instance _(unless a load balancer is involved)_.

The **"Shard Instance"** is tasked with setting up **Discord Sessions** and forwarding all incoming events from Discord to an Event Handler instance _(likely over the network)_.

_It's not recommended to create event-specific Shard Instances because Discord has over 70 events: Using event-specific Shard Instances implies that you use 70+ instances PER SHARD, each receiving duplicate guild event data._

### Event Handler Instance (Service)

The **"Event Handler Instance"** is tasked with handling incoming events of a specific type. 

**Service Based**

Create a unique Event Handler Instance per service, leaving you with as many unique Event Handler instances as you have services.

**Microservice Based**

Create a unique Event Handler Instance per event, leaving you with 70+ unique Event Handler instances, each handling a different event type.

## QA

### When do I need to shard?

Discord requires you to shard your Discord Bot once it's in a [certain number](https://discord.com/developers/docs/topics/gateway#sharding) of guilds.

### What are the implications of using one server to shard?

Servers are computers with **CPU**, **RAM**, and **Storage**. You typically run one application on a server because you expect that application to use **all** of the server's resources _(i.e 100% CPU, 100% RAM, etc)_. 

Placing multiple applications on one server is only useful when your application does **NOT** use all of the server's resources, cores, etc. This strategy implies that your application handles a low amount of data, experiences a bottleneck _(e.g., waiting on a network request)_, or maintains a consistent load.

If a server with two cores — without any form of multithreading — has an application using _<100% CPU_ on one core, then you can add an additional application _(that uses the other core)_ to the server without a performance hit.

_In practice, scaling this way is **NOT** this straightforward._

If you need to shard your bot efficiently, you _probably_ need to use multiple servers with multiple applications that all represent your "Discord Bot" as a single entity: This entity — containing multiple servers — is known as a [cluster](https://en.wikipedia.org/wiki/Computer_cluster). 

Each servers' application(s) would accept a different amount of shards and process the shard's data accordingly. Keep in mind that these applications **CAN** be built from the same codebase that was used before sharding, but require modification if the bot implements cross-guild functionality. 

_Otherwise, all most cases require is for you to implement this module in your application._