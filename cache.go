package velty

import (
	"github.com/viant/velty/est"
	"sync"
)

type (
	cache struct {
		cacheSize int
		cache     map[string]int
		values    []*value
		locker    sync.Locker
	}

	value struct {
		planner *Planner
		compute est.Compute
	}
)

func newCache(cacheSize int) *cache {
	return &cache{
		cache:     map[string]int{},
		locker:    &sync.Mutex{},
		cacheSize: cacheSize,
	}
}

func (c *cache) expression(name string) (*value, bool) {
	val, ok := c.cache[name]
	if !ok {
		return nil, false
	}

	return c.values[val], ok
}

func (c *cache) put(name string, planner *Planner, compute est.Compute) {
	if c.cacheSize == 0 {
		return
	}

	c.locker.Lock()
	defer c.locker.Unlock()
	c.cleanCacheIfNeeded()

	c.values = append(c.values, &value{
		planner: planner,
		compute: compute,
	})
	c.cache[name] = len(c.values) - 1
}

func (c *cache) cleanCacheIfNeeded() {
	if len(c.cache) > c.cacheSize {
		c.cache = map[string]int{}
		c.values = c.values[:0]
	}
}
