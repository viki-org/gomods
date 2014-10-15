objcache is the client for a simple object cache layer in Redis

# Why?
Our endpoints are hitting scaling problem, since multiple requests are coming in, asking for things in database. Most of these request, however, is read queries, and could be better served from a read cache.

This will help scale our endpoints performance, at the cost of getting stale data once in a while

# Why not use X: where X is

* __golang/groupcache__: somewhat similar to memcache, but not as familar to the team as Redis. (biased, but we want something familiar)
* __beego/cache__: its redis cache works with simple value, but don't support more complicated data (still need to serialize them)

Hence we created __objcache__

# Features

* Cache struct into Redis
* Thread-safe
* Could contain stale db data, but should always become consistent with db eventually
* Serialization/deserialization for library users (right now through interfaces)
* [Future] Can utilize a Redis master/slave setup or Redis cluster
* [Future] Queue listener to update the cache when other services change them

# Design

<img src="https://docs.google.com/drawings/d/1PovsHTMWC4f-CqyFkO-viUHEdKzqK09Fix5feCYLfMU/pub?w=808&amp;h=410">
