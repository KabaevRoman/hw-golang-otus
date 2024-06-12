package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
	mu       sync.Mutex
}

type ValueStruct struct {
	Key Key
	Val interface{}
}

func (l *lruCache) Set(key Key, value interface{}) bool {
	val, ok := l.items[key]
	l.mu.Lock()
	defer l.mu.Unlock()
	if ok {
		val.Value = ValueStruct{key, value}
		l.queue.MoveToFront(val)
		return ok
	}
	if l.queue.Len() == l.capacity {
		backElem := l.queue.Back()
		delete(l.items, backElem.Value.(ValueStruct).Key)
		l.queue.Remove(backElem)
	}
	l.items[key] = l.queue.PushFront(ValueStruct{Val: value, Key: key})
	return ok
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	val, ok := l.items[key]
	if ok {
		l.queue.MoveToFront(val)
		return val.Value.(ValueStruct).Val, ok
	}
	return nil, ok
}

func (l *lruCache) Clear() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.queue = NewList()
	l.items = make(map[Key]*ListItem, l.capacity)
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
