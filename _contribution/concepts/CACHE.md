# What is a Cache?

A [cache (computing)](https://en.wikipedia.org/wiki/Cache_(computing)) is a component that stores data for quick retrieval.

## Why Cache?

A cache stores data so that future requests for that data can be served immediately. A cache is an alternative to sending redundant requests — which take time to complete — to Discord.

For example, retrieving _the amount of users in a guild_ requires a network request to be sent to Discord, and another network request to be returned from Discord. Without a cache, retrieving _the exact same information_ again requires two more network requests, even when the guild's condition remains unchanged. 

## When to Use a Cache?

Use a cache to stop wasting time executing costly requests and calculations relevant to the application's lifetime. An application cache should **NOT** used for long-term storage because the cache resets when the application restarts.

**Use a database when you need data to persist when your bot restarts.**

## How Does a Cache Work?

_Read [Caching Overview](https://aws.amazon.com/caching) for an in-depth explanation._

A cache is typically stored in-memory, which allows the application to store and retrieve data fast (with minimal latency).

A cache receiving a request — for the number of users in a guild — will store an entry for use later. When a request is made for _the exact same information_, the cache will use the in-memory entry instead of creating a costly network request. 

**Cache Invalidation** describes the process of replacing or removing cache entries. In the above example, we know to invalidate or update the stored value for the _amount of users in a guild_ when a user joins or leaves the server.

_For more information, read [Cache Invalidation](https://en.wikipedia.org/wiki/Cache_invalidation)._


## How Do I Cache?

Disgo provides an **optional** cache along with a **cache interface** for your Discord Bot.

Read [The Disgo Cache](/cache/README.md) for information about its implementation.
