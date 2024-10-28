package cache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) (interface{}, bool)
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
	mutex    sync.Mutex
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
		mutex:    sync.Mutex{},
	}
}

func (c *lruCache) Set(key Key, value interface{}) (interface{}, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	var deletedValue interface{}
	item, ok := c.items[key]
	if ok {
		c.queue.MoveToFront(item)
		item.Value = value
		return nil, true
	}
	item = c.queue.PushFront(value)
	c.items[key] = item
	if len(c.items) > c.capacity {
		tail := c.queue.Back()
		deletedValue = tail.Value
		for k, listItem := range c.items {
			if listItem == tail {
				delete(c.items, k)
				break
			}
		}
		c.queue.Remove(tail)
	}
	return deletedValue, false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	item, ok := c.items[key]
	if ok {
		c.queue.MoveToFront(item)
		return item.Value, true
	}
	return nil, false
}

func (c *lruCache) Clear() {
	c.mutex.Lock()
	c.items = make(map[Key]*ListItem, c.capacity)
	c.queue = NewList()
	c.mutex.Unlock()
}
