package plan

import (
	"github.com/viant/velty/est"
	"sync"
)

type (
	Cache struct {
		cacheSize int
		cache     map[string]int
		values    []*Value
		locker    sync.Locker
	}

	Value struct {
		Planner *Planner
		Compute est.Compute
	}
)

func NewCache(cacheSize int) *Cache {
	return &Cache{
		cache:     map[string]int{},
		locker:    &sync.Mutex{},
		cacheSize: cacheSize,
	}
}

func (c *Cache) Expression(name string) (*Value, bool) {
	val, ok := c.cache[name]
	if !ok {
		return nil, false
	}

	return c.values[val], ok
}

func (c *Cache) Put(name string, planner *Planner, compute est.Compute) {
	if c.cacheSize == 0 {
		return
	}

	c.locker.Lock()
	defer c.locker.Unlock()
	c.cleanCacheIfNeeded()

	c.values = append(c.values, &Value{
		Planner: planner,
		Compute: compute,
	})
	c.cache[name] = len(c.values) - 1
}

func (c *Cache) cleanCacheIfNeeded() {
	if len(c.cache) > c.cacheSize {
		c.cache = map[string]int{}
		c.values = c.values[:0]
	}
}
