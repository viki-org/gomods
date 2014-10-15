package objcache

import (
	`fmt`
	`github.com/garyburd/redigo/redis`
)

const (
	eCannotGetConnection = `cannot connect to Redis`
	eCannotMarshal       = `cannot marshal %v into json; err=%v`
	eCannotSet           = `cannot run SET %s %v; err=%v`
	eCannotGet           = `cannot run GET %s; err=%v`
	eCannotRead          = `cannot read the data %v; err=%v`
	eCannotFetch         = `cannot fetch data from backend; key=%s`
)

// CachableItem is anything that can be put into the cache
type CachableItem interface {
	// The key in redis. The CachableItem has to manage duplication by itself
	Key() string
	// Decode raw bytes and update the item
	Decode(raw []byte) error //
	// Encode the item into bytes
	Encode() ([]byte, error)
}

// Set method serializes & writes an item implementing CachableItem interface into Redis
// Set can only serialize exported fields of a given item
func Set(item CachableItem) error {
	buffer, err := item.Encode()
	if err != nil {
		return fmt.Errorf(eCannotMarshal, item, err)
	}

	conn, ok := writerPool.Get().(redis.Conn)
	if !ok {
		// both conn == nil and can't cast case are here
		return fmt.Errorf(eCannotGetConnection)
	}
	defer writerPool.Put(conn)

	_, err = conn.Do(`SET`, item.Key(), buffer)
	if err != nil {
		return fmt.Errorf(eCannotSet, item.Key(), buffer, err)
	}

	return nil
}

// Get reads the cached data for a key and store into the given CachableItem
// If key is not in cache, it will use the fetch function to get it from other places
func Get(key string, item CachableItem, fetch func() (CachableItem, error)) error {
	// get connection
	conn, ok := readerPool.Get().(redis.Conn)
	if !ok {
		return fmt.Errorf(eCannotGetConnection)
	}
	defer readerPool.Put(conn)

	// read bytes
	data, err := conn.Do(`GET`, key)
	if err != nil {
		return fmt.Errorf(eCannotGet, key, err)
	}

	// return if in cache
	if data != nil {
		switch data.(type) {
		case []uint8:
			err := item.Decode(data.([]byte))
			if err != nil {
				return fmt.Errorf(eCannotRead, data, `failed to decode`)
			}
			return nil
		default:
			return fmt.Errorf(eCannotRead, data, `redis did not return []unit8`)
		}
	}

	// if not in cache, fetch and set
	fetched, err := fetch()
	if err != nil {
		return err
	}
	if fetched == nil {
		return fmt.Errorf(eCannotFetch, key)
	}
	err = Set(fetched)
	if err != nil {
		return err
	}
	// return the fetched data
	// complicated because item (addresses) might not be assignable
	buffer, err := fetched.Encode()
	if err != nil {
		return err
	}
	return item.Decode(buffer)
}
