package common

import (
	"sync"
	"time"

	"github.com/jellydator/ttlcache/v3"
)

// pendingCacheRequest represents an ongoing request for a specific cache key
type pendingCacheRequest[V any] struct {
	waitGroup sync.WaitGroup
	result    V
	err       error
}

type QueuedCache[K comparable, V any] struct {
	cache *ttlcache.Cache[K, V]

	pendingMutex sync.RWMutex
	pending      map[K]*pendingCacheRequest[V]
}

func NewQueuedCache[K comparable, V any](ttl time.Duration) *QueuedCache[K, V] {
	cache := ttlcache.New(ttlcache.WithTTL[K, V](ttl), ttlcache.WithDisableTouchOnHit[K, V]())
	go cache.Start()

	return &QueuedCache[K, V]{
		cache:        cache,
		pendingMutex: sync.RWMutex{},
		pending:      make(map[K]*pendingCacheRequest[V]),
	}
}

func (c *QueuedCache[K, V]) GetOrSet(key K, fn func() (V, error)) (V, error) {
	var zero V

	// Check cache first
	if cached := c.cache.Get(key); cached != nil {
		return cached.Value(), nil
	}

	c.pendingMutex.Lock()
	if pending, exists := c.pending[key]; exists {
		c.pendingMutex.Unlock()
		pending.waitGroup.Wait()
		if pending.err != nil {
			return zero, pending.err
		}
		return pending.result, nil
	}

	// Create new pending request
	pending := &pendingCacheRequest[V]{}
	pending.waitGroup.Add(1)
	c.pending[key] = pending
	c.pendingMutex.Unlock()

	// Perform the actual request
	result, err := fn()

	// Clean up pending request and set result
	c.pendingMutex.Lock()
	delete(c.pending, key)
	c.pendingMutex.Unlock()

	pending.result = result
	pending.err = err
	pending.waitGroup.Done()

	if err != nil {
		return zero, err
	}

	c.cache.Set(key, result, 0)
	return result, nil
}

func (c *QueuedCache[K, V]) Get(key K) (V, bool) {
	var zero V

	if cached := c.cache.Get(key); cached != nil {
		return cached.Value(), true
	}

	return zero, false
}

func (c *QueuedCache[K, V]) Set(key K, value V) {
	c.cache.Set(key, value, 0)
}

func (c *QueuedCache[K, V]) Delete(key K) {
	c.cache.Delete(key)
}
