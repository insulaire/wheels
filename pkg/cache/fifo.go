package cache

import (
	"container/list"
	"sync"
)

type fifo struct {
	maxBytes  int
	onEvicted func(key string, value interface{})
	usedBytes int
	ll        *list.List
	cache     map[string]*list.Element
	lock      sync.RWMutex
}

func New(maxBytes int, onEvicted func(key string, value interface{})) Cache {
	return &fifo{
		maxBytes:  maxBytes,
		onEvicted: onEvicted,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		lock:      sync.RWMutex{},
	}
}

func (f *fifo) Set(key string, value interface{}) {
	f.lock.Lock()
	defer f.lock.Unlock()
	if e, ok := f.cache[key]; ok {
		f.ll.MoveToBack(e)
		en := e.Value.(*entry)
		f.usedBytes = f.usedBytes - CalcLen(en.value) + CalcLen(value)
		en.value = value
		return
	}
	en := &entry{key, value}
	e := f.ll.PushBack(en)
	f.cache[key] = e
	f.usedBytes += en.Len()
	if f.maxBytes > 0 && f.usedBytes > f.usedBytes {
		f.DelOldst()
	}
}

func (f *fifo) Get(key string) interface{} {
	f.lock.RLock()
	defer f.lock.RUnlock()
	if e, ok := f.cache[key]; ok {
		return e.Value.(*entry).value
	}
	return nil
}

func (f *fifo) Del(key string) {
	f.lock.Lock()
	defer f.lock.Unlock()
	if e, ok := f.cache[key]; ok {
		f.removeElement(e)
	}
}

func (f *fifo) DelOldst() {
	f.lock.Lock()
	f.removeElement(f.ll.Front())
	f.lock.Unlock()
}

func (f *fifo) Len() int {
	return f.ll.Len()
}

func (f *fifo) removeElement(e *list.Element) {
	if e == nil {
		return
	}
	f.ll.Remove(e)
	en := e.Value.(*entry)
	f.usedBytes -= en.Len()
	delete(f.cache, en.key)
	if f.onEvicted != nil {
		f.onEvicted(en.key, en.value)
	}
}
