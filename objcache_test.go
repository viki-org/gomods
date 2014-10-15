package objcache

import (
	`github.com/exklamationmark/glog`
	`github.com/garyburd/redigo/redis`
	. `gopkg.in/check.v1`
	`sync`
	`testing`
)

var (
	testReaderPool *sync.Pool
)

func Test(t *testing.T) {
	TestingT(t)
}

type objcacheTS struct{}

func init() {
	Suite(&objcacheTS{})
	Configure(Configuration{
		WriterURL: `:6379`,
		ReaderURL: `:6379`,
	})

	testReaderPool = &sync.Pool{
		New: newRedisConn(Config.ReaderURL),
	}
}

func (s *objcacheTS) SetUpTest(c *C) {
	flushTestRedis()
}

func flushTestRedis() {
	conn, ok := testReaderPool.Get().(redis.Conn)
	if !ok {
		glog.Fatal(`can't connect to Redis`)
	}
	defer testReaderPool.Put(conn)

	conn.Do(`FLUSHDB`)
}
