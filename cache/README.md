# Disgo Cache

The Disgo Cache is a Go module that automatically handles caching for your Discord Bot. The Disgo Cache implements a **cache interface** that allows you to replace the built-in cache with another store _(such as Redis or Memcached)_ or provide your own method of caching data. For information on the concept of caching, read [What is a Cache?](contribution/concepts/CACHE.md)

## How It Works

The Disgo Cache is used when the `disgo.Client` creates a **Request** or receives a **Session** event. 

### Caching Resources

### Caching Events