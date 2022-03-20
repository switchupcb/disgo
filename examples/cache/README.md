# Cache

This example creates a bot that uses modified cache settings to create an application command and handle it. For more information about caching, read [What is a Cache?](contribution/concepts/CACHE.md)

## Configuration

Use the client to configure the bot's settings.

```go
bot := disgo.Client{
    Config: disgo.Config{

    },


    // Override the automatic cache.
    Cache: disgo.Cache{
        // Use the client to manage the cache settings.
        Config:  disgo.CacheConfig{ 

        },

        // Populate the cache with information prior to its usage.
        // Struct
    },
}
```