# What is a Cache?

## Why Do We Cache?

Using **Sessions** and **Requests** allow us to retrieve data from Discord, but it becomes redundant to receive the same data over and over again. For example, retrieving _the amount of users in a guild_ requires a network request to be sent to and returned from Discord. Without a cache, retrieving the _exact same information_ requires another network request; even when the amount of users in the guild remains unchanged. A cache is an alternative to sending redundant requests — that take time to complete — to Discord. A cache stores data so that future requests for that data can be served immediately (without latency).

## When to Use a Cache?
Caches are useful for storing costly requests or calculations relevant to the lifetime of the application. In other words, A cache is not meant to be used for long-term storage. **If you need data to persist when your bot restarts, use a database.**

## How Does a Cache Work?

Read [Caching Overview](https://aws.amazon.com/caching) for an in-depth explanation. A cache is typically stored in-memory which allows the application to store and retrieve data fast (without a network request). A cache that receives a request — for the amount of users in a guild — will store it for use later. When a request is made for the same information, the cache will use the in-memory value instead of creating a network request. **Cache Invalidation** describes the process of replacing or removing cache entries. In the example above, we know to invalidate the stored value for the _amount of users in a guild_ when a user joins or leaves the server. For more information, read [Cache Invalidation](https://en.wikipedia.org/wiki/Cache_invalidation).

# The Disgo Cache
Disgo provides an **optional** cache along with a **cache interface**. The cache interface allows you to replace the built-in cache with another store _(such as Redis or Memcached)_. Use the cache interface to provide your own method of caching data.

## How It Works

The Disgo Cache is used when the `disgo.Client` creates a **Request** or receives a **Session** event. 

### Caching Resources

### Caching Events