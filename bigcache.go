package toold

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/allegro/bigcache"
)

//Cache Cache
type Cache struct {
	Cache *bigcache.BigCache
}

//NewCache NewCache
func NewCache(times time.Duration, maxCacheSize int, maxEntries int, funs func(key string, entry []byte, reason bigcache.RemoveReason)) (*Cache, error) {
	if maxEntries == 0 {
		maxEntries = 10000
	} else {
		maxEntries = maxEntries * 10000
	}
	cache, err := bigcache.NewBigCache(bigcache.Config{
		Shards:     1024 * 2,
		LifeWindow: times,
		// CleanWindow:        1,
		MaxEntriesInWindow: maxEntries,
		MaxEntrySize:       1024,
		Verbose:            true,
		HardMaxCacheSize:   1024 * maxCacheSize,
		Logger:             bigcache.DefaultLogger(),
		OnRemoveWithReason: funs,
	})
	if err != nil {
		return nil, err
	}
	return &Cache{
		Cache: cache,
	}, nil
}

//SetString SetString
func (c *Cache) SetString(key string, val string) error {
	return c.Cache.Set(key, []byte(val))
}

//Set Set
func (c *Cache) Set(key string, val interface{}) error {
	return c.Cache.Set(key, []byte(fmt.Sprintf("%v", val)))
}

//SetObj
func (c *Cache) SetObj(key string, val interface{}) error {
	body, err := json.Marshal(val)
	if err != nil {
		return err
	}
	return c.Cache.Set(key, body)
}

//GetString GetString
func (c *Cache) GetString(key string) (string, error) {
	by, er := c.Cache.Get(key)
	if er != nil {
		return "", er
	}
	return string(by), nil
}

//GetInt GetInt
func (c *Cache) GetInt(key string) (int, error) {
	by, er := c.Cache.Get(key)
	if er != nil {
		return 0, er
	}
	return ConversionToInt(string(by)), nil
}

//GetInt GetInt
func (c *Cache) GetInt64(key string) (int64, error) {
	by, er := c.Cache.Get(key)
	if er != nil {
		return 0, er
	}
	return ConversionToInt64(string(by)), nil
}

func (c *Cache) GetObj(key string, v interface{}) error {
	by, err := c.Cache.Get(key)
	if err != nil {
		return err
	}
	return json.Unmarshal(by, v)
}

func (c *Cache) Delete(key string) error {
	return c.Cache.Delete(key)
}
