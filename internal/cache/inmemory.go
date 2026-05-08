package cache

import (
	"time"

	gocache "github.com/patrickmn/go-cache"
)

type Cache struct {
	c *gocache.Cache
}

func New(defaultTTL, cleanupInterval time.Duration) *Cache {
	return &Cache{c: gocache.New(defaultTTL, cleanupInterval)}
}

func (c *Cache) Set(key string, value any, ttl time.Duration) {
	c.c.Set(key, value, ttl)
}

func (c *Cache) Get(key string) (any, bool) {
	return c.c.Get(key)
}

func (c *Cache) Delete(key string) {
	c.c.Delete(key)
}
