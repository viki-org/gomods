// Package objcache implements an object cache which is stored in Redis
// Motivation: to cache objects (not simple value) in Redis and use it as the read cache instead of hitting db
// Why not use X: where X is
//   - golang/groupcache: similar to memcache, but not as familar to the team as Redis. Favor maintainability
//   - beego/cache: works with simple value, but don't support deserializing structs. Still have to write that
//
// Goal for the objcache
// 	 - Cache struct into Redis
//   - Thread-safe
//	 - Could contain stale db data, but should always become consistent eventually
//   - [Future] Can utilize a Redis master/slave setup or Redis cluster
//   - Do serializationa and deserialization for client users (require all struct fields to be exported). [Future] can work with unexported data(?)
package objcache

import (
	`github.com/exklamationmark/glog`
	`github.com/garyburd/redigo/redis`
	`sync`
)

// Config stores configurations for objcache
type Configuration struct {
	WriterURL, ReaderURL string // separate reader/writer to prepare for a master/slave Redis setup
}

var (
	Config Configuration

	// writerPool and readerPool are redis connection pools to be used for writing and reading
	// no of connections in the  pool scale up under load, but get garbarge-collected when it's not
	// no max no of connection as Redis can take ~10k clients easily // future self, watch out
	writerPool, readerPool *sync.Pool
)

// Configure setup the objcache client. It must be called in the caller's init()
func Configure(config Configuration) {
	Config = config
	writerPool = &sync.Pool{
		New: newRedisConn(Config.WriterURL),
	}
	readerPool = &sync.Pool{
		New: newRedisConn(Config.ReaderURL),
	}
}

// returns a redis connection or nil (when cannot connect)
func newRedisConn(url string) func() interface{} {
	return func() interface{} {
		conn, err := redis.Dial(`tcp`, url)
		if err != nil {
			glog.Error(`could not connect to redis: err=`, err)
			return nil
		}
		return conn
	}
}
