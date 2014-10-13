package objcache

import (
	`encoding/json`
	`fmt`
	`github.com/garyburd/redigo/redis`
)

const (
	eItemNotCachable     = `item %v does not implement CachableItem interface`
	eCannotGetConnection = `cannot connect to Redis`
	eCannotMarshal       = `cannot marshal %v into json; err=%v`
	eCannotSet           = `cannot run SET %s %v; err=%v`
)

type CachableItem interface {
	Key() string // generate a key for the object
}

// Set method serializes & writes an item implementing CachableItem interface into Redis
// Set can only serialize exported fields of a given item
func Set(item interface{}) error {
	cachable, ok := item.(CachableItem)
	if !ok {
		return fmt.Errorf(eItemNotCachable, item)
	}

	buffer, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf(eCannotMarshal, item, err)
	}

	conn, ok := writerPool.Get().(redis.Conn)
	if !ok { // covers both when Get return nil & different type
		return fmt.Errorf(eCannotGetConnection)
	}
	defer writerPool.Put(conn)

	_, err = conn.Do(`SET`, cachable.Key(), buffer)
	if err != nil {
		return fmt.Errorf(eCannotSet, cachable.Key(), buffer, err)
	}

	return nil
}
