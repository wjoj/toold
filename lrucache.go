package toold

import (
	"sync"
	"time"

	"github.com/golang/groupcache/lru"
)

//CacheLru CacheLru
type CacheLru struct {
	cache     *lru.Cache
	times     map[string]*TickerObj
	timesLock sync.RWMutex
}

//NewLruCache NewLruCache
func NewLruCache(size int) *lru.Cache {
	return lru.New(size)
}

//NewCacheLru NewCacheLru
func NewCacheLru(size int) *CacheLru {
	return &CacheLru{
		cache: NewLruCache(size),
		times: make(map[string]*TickerObj),
	}
}

//Set Set
func (c *CacheLru) Set(key string, val interface{}) {
	c.cache.Add(key, val)
}

func (c *CacheLru) getTimeer(key string) *TickerObj {
	c.timesLock.RLock()
	defer c.timesLock.RUnlock()
	return c.times[key]
}

func (c *CacheLru) setTimeer(key string, t *TickerObj) {
	c.timesLock.Lock()
	defer c.timesLock.Unlock()
	c.times[key] = t
}

func (c *CacheLru) delTimeer(key string) {
	c.timesLock.Lock()
	defer c.timesLock.Unlock()
	delete(c.times, key)
}

//SetDoExpire SetDoExpire
func (c *CacheLru) SetDoExpire(key string, val interface{}, expire time.Duration) {
	c.cache.Add(key, val)
	t := c.getTimeer(key)
	if t != nil {
		t.Stop()
		t = nil
	}
	t = timer(expire, func(t *TickerObj) bool {
		// t.Stop()
		// t = nil
		c.cache.Remove(key)
		c.delTimeer(key)
		return true
	})
	c.setTimeer(key, t)
}

//GetKey GetKey
func (c *CacheLru) GetKey(key string) interface{} {
	val, _ := c.cache.Get(key)
	return val
}

//GetInt GetInt
func (c *CacheLru) GetInt(key string) int {
	val, s := c.cache.Get(key)
	if !s {
		return 0
	}
	return ConversionToInt(val)
}

//DeleteKey DeleteKey
func (c *CacheLru) DeleteKey(key string) {
	t := c.getTimeer(key)
	if t != nil {
		t.Stop()
		t = nil
		c.delTimeer(key)
	}
	c.cache.Remove(key)
}
