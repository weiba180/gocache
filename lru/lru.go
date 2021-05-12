package lru

import (
	"container/list"
)

//实现一个基础的lru缓存，非并发安全
type Cache struct {
	maxBytes  int                           //最大内存
	nowBytes  int                           //已使用内存
	ll        *list.List                    //双向链表
	cache     map[string]*list.Element      //键值对，指向数据节点
	onEvicted func(key string, value Value) //数据未命中时的回调函数，可以为nil
}

//键值对
type entry struct {
	key   string
	value Value //接口 需实现Len()方法
}

type Value interface {
	Len() int
}

//Cache的构造函数
func New(maxBytes int, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		onEvicted: onEvicted,
	}
}

//添加节点，若节点存在则直接修改，否则新建，并添加到双向列表的头部
//先更新已使用内存，若溢出则先执行淘汰算法再写入节点数据
func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nowBytes += value.Len() - kv.value.Len()
		c.Eliminate()
		kv.value = value
	} else {
		ele := c.ll.PushFront(&entry{key, value})
		c.nowBytes += len(key) + value.Len()
		c.Eliminate()
		c.cache[key] = ele

	}

}

//get 数据节点，并将节点移至头部
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

//内存淘汰
func (c *Cache) Eliminate() {
	for c.maxBytes != 0 && c.maxBytes < c.nowBytes {
		c.RemoveOldest()
	}
}
func (c *Cache) RemoveOldest() {
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.nowBytes -= len(kv.key) + kv.value.Len()
		if c.onEvicted != nil {
			c.onEvicted(kv.key, kv.value)
		}
	}
}

func (c *Cache) Len() int {
	return c.ll.Len()
}
