package client

// Top level resources:
// channel_id, guild_id, webhook_id
// endpoints with two different top-levle resources can have independent rate limits

// client config has max retries
// when you receive rate limit bucket
// create a new thread/go routine
// on that thread, we wait
// send once finished waiting
