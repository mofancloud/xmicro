package memcache

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/mofancloud/xmicro/cache"
)

// Cache Memcache adapter.
type Cache struct {
	conn     *memcache.Client
	conninfo []string
}

// NewMemCache create new memcache adapter.
func NewMemCache() cache.Cache {
	return &Cache{}
}

// Get get value from memcache.
func (rc *Cache) Get(key string) interface{} {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return err
		}
	}
	if item, err := rc.conn.Get(key); err == nil {
		return item.Value
	}
	return nil
}

// GetMulti get value from memcache.
func (rc *Cache) GetMulti(keys []string) []interface{} {
	size := len(keys)
	var rv []interface{}
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			for i := 0; i < size; i++ {
				rv = append(rv, err)
			}
			return rv
		}
	}
	mv, err := rc.conn.GetMulti(keys)
	if err == nil {
		for _, v := range mv {
			rv = append(rv, v.Value)
		}
		return rv
	}
	for i := 0; i < size; i++ {
		rv = append(rv, err)
	}
	return rv
}

// Put put value to memcache.
func (rc *Cache) Put(key string, val interface{}, timeout time.Duration) error {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return err
		}
	}
	item := memcache.Item{Key: key, Expiration: int32(timeout / time.Second)}
	if v, ok := val.([]byte); ok {
		item.Value = v
	} else if str, ok := val.(string); ok {
		item.Value = []byte(str)
	} else {
		return errors.New("val only support string and []byte")
	}
	return rc.conn.Set(&item)
}

// Delete delete value in memcache.
func (rc *Cache) Delete(key string) error {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return err
		}
	}
	return rc.conn.Delete(key)
}

// Incr increase counter.
func (rc *Cache) Incr(key string) error {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return err
		}
	}
	_, err := rc.conn.Increment(key, 1)
	return err
}

// Decr decrease counter.
func (rc *Cache) Decr(key string) error {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return err
		}
	}
	_, err := rc.conn.Decrement(key, 1)
	return err
}

// IsExist check value exists in memcache.
func (rc *Cache) IsExist(key string) bool {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return false
		}
	}
	_, err := rc.conn.Get(key)
	return !(err != nil)
}

// ClearAll clear all cached in memcache.
func (rc *Cache) ClearAll() error {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return err
		}
	}
	return rc.conn.FlushAll()
}

// StartAndGC start memcache adapter.
// config string is like {"conn":"connection info"}.
// if connecting error, return.
func (rc *Cache) StartAndGC(config string) error {
	var cf map[string]string
	json.Unmarshal([]byte(config), &cf)
	if _, ok := cf["conn"]; !ok {
		return errors.New("config has no conn key")
	}
	rc.conninfo = strings.Split(cf["conn"], ";")
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return err
		}
	}
	return nil
}

// connect to memcache and keep the connection.
func (rc *Cache) connectInit() error {
	rc.conn = memcache.New(rc.conninfo...)
	return nil
}

func init() {
	cache.Register("memcache", NewMemCache)
}
