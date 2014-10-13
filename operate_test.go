package objcache

import (
	`bytes`
	`fmt`
	`github.com/exklamationmark/glog`
	`github.com/garyburd/redigo/redis`
	. `gopkg.in/check.v1`
	`time`
)

type simpleTStruct struct {
	Name     string
	Age      int
	Birthday time.Time
}

type nestedTStruct struct {
	Alias  string
	Simple simpleTStruct
}

type embeddedTStruct struct {
	Side string
	*simpleTStruct
}

var testTime, _ = time.Parse(`2006-01-02 15:04:05-07`, `2014-06-29 08:00:00+07`)

func (s simpleTStruct) Key() string {
	return `name:` + s.Name
}

func (s nestedTStruct) Key() string {
	return `alias:` + s.Alias
}

func (s embeddedTStruct) Key() string {
	return `side:` + s.Side
}

func (s *objcacheTS) TestSet(c *C) {
	// item: cachable, not cachable
	// conn: have / don't have (how to test)
	// item: simple value, simple struct (1 level), nested struct, embedded struct

	tests := []struct {
		item   interface{}
		key    string
		actual string
		err    error
	}{
		{
			20,
			``,
			``,
			fmt.Errorf(eItemNotCachable, 20),
		},
		{
			simpleTStruct{`yoda`, 200, testTime},
			`name:yoda`,
			`{"Name":"yoda","Age":200,"Birthday":"2014-06-29T08:00:00+07:00"}`,
			nil,
		},
		{
			nestedTStruct{`old-jedi`, simpleTStruct{`yoda`, 200, testTime}},
			`alias:old-jedi`,
			`{"Alias":"old-jedi","Simple":{"Name":"yoda","Age":200,"Birthday":"2014-06-29T08:00:00+07:00"}}`,
			nil,
		},
		{
			embeddedTStruct{`light`, &simpleTStruct{`yoda`, 200, testTime}},
			`side:light`,
			`{"Side":"light","Name":"yoda","Age":200,"Birthday":"2014-06-29T08:00:00+07:00"}`,
			nil,
		},
	}

	for _, test := range tests {
		actualErr := Set(test.item)

		if test.err != nil {
			c.Assert(actualErr, DeepEquals, test.err)
		} else {
			c.Assert(actualErr, IsNil)
			data, err := getTestData(test.key)
			// glog.Info(`guess=`, string(toBytes(test.actual)))
			// glog.Info(`read =`, string(data.([]byte)))
			c.Assert(err, IsNil)
			c.Assert(data, DeepEquals, toBytes(test.actual))
		}
	}
}

func getTestData(key string) (interface{}, error) {
	conn, ok := testReaderPool.Get().(redis.Conn)
	defer testReaderPool.Put(conn)
	if !ok {
		glog.Fatal(`cannot connect to redis`)
	}
	return conn.Do(`GET`, key)
}

func toBytes(data interface{}) []byte {
	buffer := bytes.NewBuffer(nil)
	fmt.Fprint(buffer, data)
	return buffer.Bytes()
}
