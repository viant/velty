package plan

import (
	"github.com/viant/velty/est"
	"sync"
)

type (
	Cache struct {
		cacheSize int
		cache     map[string]*Value
		locker    sync.Locker
	}

	Value struct {
		Planner *Planner
		Compute est.Compute
	}
)

func NewCache(cacheSize int) *Cache {
	return &Cache{
		cache:     map[string]*Value{},
		locker:    &sync.Mutex{},
		cacheSize: cacheSize,
	}
}

func (c *Cache) Expression(name string) (*Value, bool) {
	val, ok := c.cache[name]
	return val, ok
}

func (c *Cache) Put(name string, planner *Planner, compute est.Compute) {
	if c.cacheSize == 0 {
		return
	}

	c.locker.Lock()
	defer c.locker.Unlock()
	c.cleanCacheIfNeeded()

	c.cache[name] = &Value{
		Planner: planner,
		Compute: compute,
	}
}

func (c *Cache) cleanCacheIfNeeded() {
	if len(c.cache) > c.cacheSize {
		c.cache = map[string]*Value{}
	}
}
