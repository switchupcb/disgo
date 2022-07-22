# What is a Cache?

Disgo provides an **optional** cache along with a **cache interface** for your Discord Bot. Read [The Disgo Cache](/cache/README.md) for information about its implementation.

## Why Do We Cache?

Using **Requests** and **Sessions** allow us to retrieve data from Discord, but it becomes redundant to receive the same data over and over again. For example, retrieving _the amount of users in a guild_ requires a network request to be sent _to_ and returned _from_ Discord. Without a cache, retrieving _the exact same information_ requires yet another network request; even when guild's condition remains unchanged. A cache is an alternative to sending redundant requests — that take time to complete — to Discord. A cache stores data so that future requests for that data can be served immediately.

## When to Use a Cache?

Caches are useful for storing costly requests or calculations relevant to the lifetime of the application. In other words, A cache is not meant to be used for long-term storage. **If you need data to persist when your bot restarts, use a database.**

## How Does a Cache Work?

Read [Caching Overview](https://aws.amazon.com/caching) for an in-depth explanation. A cache is typically stored in-memory which allows the application to store and retrieve data fast (with minimal latency). A cache that receives a request — for the amount of users in a guild — will store an entry for use later. When a request is made for _the exact same information_, the cache will use the in-memory entry instead of creating a costly network request. **Cache Invalidation** describes the process of replacing or removing cache entries. In the example above, we know to invalidate or update the stored value for the _amount of users in a guild_ when a user joins or leaves the server. For more information, read [Cache Invalidation](https://en.wikipedia.org/wiki/Cache_invalidation).