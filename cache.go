package gocache

//对lru进行封装，加锁使其并发安全
import (
	"gocache/lru"
	"sync"
)

type cache struct {
	mu         sync.Mutex
	lru        *lru.Cache
	cacheBytes int
}

//使用了延迟初始化来创建lru实例
func (c *cache) add(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		c.lru = lru.New(c.cacheBytes, nil)
	}
	c.lru.Add(key, value)
}

func (c *cache) get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		return
	}
	if v, ok := c.lru.Get(key); ok {
		return v.(ByteView), ok
	}
	return
}
