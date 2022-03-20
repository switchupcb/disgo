# Shard

This example creates a bot that uses modified a shard manager to create an application command and handle it. For more information about sharding, read [What is a Discord Shard?](contributing/concepts/SHARD.md)

## Configuration

Use the client to configure the bot's settings.

```go
bot := disgo.Client{
    Config: disgo.Config{

    },


    // Override the automatic shard manager.
    Shard: disgo.Shard{
        // Use the client to manage the shard settings.
        Config:  disgo.ShardConfig{ 

        },

        // Instantiate two sessions with 4 shards.
        // map
    },
}
```