package objcache

import (
	`bytes`
	`encoding/json`
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

func (s *simpleTStruct) Key() string {
	return `name:` + s.Name
}

func (s *simpleTStruct) Encode() ([]byte, error) {
	return json.Marshal(s)
}

func (s *simpleTStruct) Decode(raw []byte) error {
	return json.Unmarshal(raw, s)
}

func (s nestedTStruct) Key() string {
	return `alias:` + s.Alias
}

func (s *nestedTStruct) Encode() ([]byte, error) {
	return json.Marshal(s)
}

func (s *nestedTStruct) Decode(raw []byte) error {
	return json.Unmarshal(raw, s)
}

func (s embeddedTStruct) Key() string {
	return `side:` + s.Side
}

func (s *embeddedTStruct) Encode() ([]byte, error) {
	return json.Marshal(s)
}

func (s *embeddedTStruct) Decode(raw []byte) error {
	return json.Unmarshal(raw, s)
}

func (s *objcacheTS) TestSet(c *C) {
	// item: cachable, not cachable
	// conn: have / don't have (how to test)
	// item: simple value, simple struct (1 level), nested struct, embedded struct

	tests := []struct {
		item   CachableItem
		key    string
		actual string
		err    error
	}{
		{
			&simpleTStruct{`yoda`, 200, testTime},
			`name:yoda`,
			`{"Name":"yoda","Age":200,"Birthday":"2014-06-29T08:00:00+07:00"}`,
			nil,
		},
		{
			&nestedTStruct{`old-jedi`, simpleTStruct{`yoda`, 200, testTime}},
			`alias:old-jedi`,
			`{"Alias":"old-jedi","Simple":{"Name":"yoda","Age":200,"Birthday":"2014-06-29T08:00:00+07:00"}}`,
			nil,
		},
		{
			&embeddedTStruct{`light`, &simpleTStruct{`yoda`, 200, testTime}},
			`side:light`,
			`{"Side":"light","Name":"yoda","Age":200,"Birthday":"2014-06-29T08:00:00+07:00"}`,
			nil,
		},
	}

	for _, test := range tests {
		flushTestRedis()
		actualErr := Set(test.item)

		if test.err != nil {
			c.Assert(actualErr, DeepEquals, test.err)
			continue
		}

		c.Assert(actualErr, IsNil)
		data, err := getTestData(test.key)
		// glog.Info(`guess=`, string(toBytes(test.actual)))
		// glog.Info(`read =`, string(data.([]byte)))
		c.Assert(err, IsNil)
		c.Assert(data, DeepEquals, toBytes(test.actual))
	}
}

func (s *objcacheTS) TestGet(c *C) {
	// key: exists, not exists
	// fetch: nil, fail, not fail

	tests := []struct {
		setup     func()
		key       string
		item      CachableItem
		fetchFunc func() (CachableItem, error)
		actual    interface{}
		err       error
	}{
		{
			func() { Set(&embeddedTStruct{`light`, &simpleTStruct{`yoda`, 200, testTime}}) },
			`side:light`,
			&embeddedTStruct{},
			func() (CachableItem, error) { return nil, nil },
			&embeddedTStruct{`light`, &simpleTStruct{`yoda`, 200, testTime}},
			nil,
		},
		{
			func() { Set(&simpleTStruct{`yoda`, 230, testTime}) },
			`name:yoda`,
			&simpleTStruct{},
			func() (CachableItem, error) { return nil, nil },
			&simpleTStruct{`yoda`, 230, testTime},
			nil,
		},
		{
			nil,
			`name:yoda`,
			&simpleTStruct{},
			func() (CachableItem, error) { return nil, nil },
			nil,
			fmt.Errorf(eCannotFetch, `name:yoda`),
		},
		{
			nil,
			`name:anakin`,
			&simpleTStruct{},
			func() (CachableItem, error) { return &simpleTStruct{`anakin`, 20, testTime}, nil },
			&simpleTStruct{`anakin`, 20, testTime},
			nil,
		},
		{
			nil,
			`name:darthvader`,
			&simpleTStruct{},
			func() (CachableItem, error) { return nil, fmt.Errorf(`%v`, `backend problem`) },
			nil,
			fmt.Errorf(`backend problem`),
		},
	}

	for _, test := range tests {
		flushTestRedis()
		if test.setup != nil {
			test.setup()
		}
		err := Get(test.key, test.item, test.fetchFunc)

		if test.err != nil {
			c.Assert(err, DeepEquals, test.err)
			continue
		}

		c.Assert(err, IsNil)
		c.Assert(test.item, DeepEquals, test.actual)
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
