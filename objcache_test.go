package objcache

import (
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
