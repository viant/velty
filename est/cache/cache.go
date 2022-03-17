package cache

import (
	"github.com/viant/velty/est"
	"sync"
)

type Cache struct {
	cache  map[string]est.Compute
	locker sync.Locker
}

func NewCache() *Cache {
	return &Cache{
		cache:  map[string]est.Compute{},
		locker: &sync.Mutex{},
	}
}

func (c *Cache) Expression(name string) (est.Compute, bool) {
	val, found := c.cache[name]
	return val, found
}

func (c *Cache) Put(name string, compute est.Compute) {
	c.locker.Lock()
	defer c.locker.Unlock()

	c.cache[name] = compute
}
