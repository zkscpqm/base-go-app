package common

import (
	"fmt"
	"runtime"
	"sync/atomic"
)

type Cache[KeyT comparable, ValT any] map[KeyT]ValT

func (m Cache[KeyT, ValT]) Keys() []KeyT {
	keys := make([]KeyT, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func (m Cache[KeyT, ValT]) Values() []ValT {
	values := make([]ValT, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}

func (m Cache[KeyT, ValT]) Merge(other map[KeyT]ValT) {
	for k, v := range other {
		m[k] = v
	}
}

type AsyncSafeCache[KeyT comparable, ValT any] struct {
	cache   Cache[KeyT, ValT]
	writing int32
	readers int32
}

func (c *AsyncSafeCache[KeyT, ValT]) write(fn func()) {
	for {
		if atomic.LoadInt32(&c.readers) == 0 && atomic.CompareAndSwapInt32(&c.writing, 0, 1) {
			break
		}
		runtime.Gosched()
	}

	defer func() {
		atomic.StoreInt32(&c.writing, 0)

		if r := recover(); r != nil {
			stack := make([]byte, 4096)
			stack = stack[:runtime.Stack(stack, false)]

			fmt.Printf("Recovered panic in AsyncSafeCache.write: %v\nStack:\n%s\n", r, stack)
		}
	}()
	fn()
}

func (c *AsyncSafeCache[KeyT, ValT]) read(fn func()) {
	for {
		if atomic.LoadInt32(&c.writing) == 0 {
			atomic.AddInt32(&c.readers, 1)

			if atomic.LoadInt32(&c.writing) == 0 {
				break
			}
			atomic.AddInt32(&c.readers, -1)
		}
		runtime.Gosched()
	}

	defer func() {
		atomic.AddInt32(&c.readers, -1)

		if r := recover(); r != nil {
			stack := make([]byte, 4096)
			stack = stack[:runtime.Stack(stack, false)]

			fmt.Printf("Recovered panic in AsyncSafeCache.read: %v\nStack:\n%s\n", r, stack)
		}
	}()
	fn()
}

func (c *AsyncSafeCache[KeyT, ValT]) Get(k KeyT) (v ValT, ok bool) {
	c.read(func() {
		v, ok = c.cache[k]
	})

	return
}

func (c *AsyncSafeCache[KeyT, ValT]) Set(k KeyT, v ValT) {
	c.write(func() {
		c.cache[k] = v
	})
}

func (c *AsyncSafeCache[KeyT, ValT]) Delete(k KeyT) {
	c.write(func() {
		delete(c.cache, k)
	})
}

func (c *AsyncSafeCache[KeyT, ValT]) SmartLookup(entries []KeyT) (result Cache[KeyT, ValT], missing []KeyT) {
	result = make(Cache[KeyT, ValT], len(entries))
	missing = make([]KeyT, 0, len(entries))

	c.read(func() {
		for _, e := range entries {
			if v, ok := c.cache[e]; ok {
				result[e] = v
			} else {
				missing = append(missing, e)
			}
		}
	})

	return
}

func (c *AsyncSafeCache[KeyT, ValT]) Merge(other map[KeyT]ValT) (newEntries []ValT) {
	newEntries = make([]ValT, 0, len(other))
	c.write(func() {
		for k, v := range other {
			if _, ok := c.cache[k]; !ok {
				c.cache[k] = v
			} else {
				newEntries = append(newEntries, v)
			}
		}
	})
	return
}

type CacheManager struct {
	SomeCache AsyncSafeCache[string, int64]
}

func NewCacheManager() CacheManager {
	return CacheManager{
		SomeCache: AsyncSafeCache[string, int64]{
			cache: make(Cache[string, int64]),
		},
	}
}
