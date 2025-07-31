package heatwave

import (
	"errors"
	"sync"
	"time"
)

const (
	defaultMaxSize         = 1000
	defaultOutdated        = time.Minute * 5
	defaultCleanupInterval = time.Minute
)

var (
	ErrBucketClosed = errors.New("bucket is closed")
)

// CacheItem represents an item in the cache with generic value type
// This structure is now decoupled from any specific update strategy
type CacheItem[T any] struct {
	key       string
	value     T
	expiredAt time.Time
}

type NewBucketOption[T any] func(b *Bucket[T])

type Bucket[T any] struct {
	name     string         // Name of the bucket
	maxSize  int            // Maximum number of items in cache
	outdated *time.Duration // TTL for cache items

	cleanupInterval time.Duration            // Interval for background cleanup
	cache           map[string]*CacheItem[T] // Hash map for O(1) access
	updater         Updater[T]               // Update strategy interface
	mutex           sync.RWMutex             // Read-write mutex for thread safety
	stopCleanup     chan struct{}            // Channel to stop cleanup goroutine
	closed          bool                     // Flag to track if bucket is closed
	closeMutex      sync.Mutex               // Mutex to protect close operation
}

func NewBucket[T any](opts ...NewBucketOption[T]) *Bucket[T] {
	od := defaultOutdated
	b := &Bucket[T]{
		maxSize:         defaultMaxSize,
		outdated:        &od,
		cache:           make(map[string]*CacheItem[T]),
		updater:         newLRUUpdater[T](),
		cleanupInterval: defaultCleanupInterval,
		stopCleanup:     make(chan struct{}, 1), // Buffered channel to prevent blocking
		closed:          false,
	}

	for _, opt := range opts {
		opt(b)
	}

	// Start background cleanup goroutine
	go b.startCleanup()

	return b
}

// Nail stores data in memory (like nailing it to memory)
func (b *Bucket[T]) Nail(id string, data T) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	// Check if bucket is closed
	if b.isClosed() {
		return ErrBucketClosed
	}

	expiredAt := time.Now().Add(*b.outdated)

	// If key already exists, update it
	if existingItem, exists := b.cache[id]; exists {
		existingItem.value = data
		existingItem.expiredAt = expiredAt
		b.updater.Access(existingItem)
		return nil
	}

	// If cache is full, remove least recently used item
	if b.updater.Size() >= b.maxSize {
		evictedItem := b.updater.Evict()
		if evictedItem != nil {
			delete(b.cache, evictedItem.key)
		}
	}

	// Create new cache item
	newItem := &CacheItem[T]{
		key:       id,
		value:     data,
		expiredAt: expiredAt,
	}

	b.cache[id] = newItem
	b.updater.Add(newItem)

	return nil
}

// Bring retrieves data from the bucket
func (b *Bucket[T]) Bring(id string) (T, bool) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	var zero T

	// Check if bucket is closed
	if b.isClosed() {
		return zero, false
	}

	item, exists := b.cache[id]
	if !exists {
		return zero, false
	}

	// Check if expired
	if time.Now().After(item.expiredAt) {
		b.updater.Remove(item)
		delete(b.cache, id)
		return zero, false
	}

	// Mark as accessed
	b.updater.Access(item)

	return item.value, true
}

// startCleanup starts the background goroutine for cleaning up expired items
func (b *Bucket[T]) startCleanup() {
	ticker := time.NewTicker(b.cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if b.isClosed() {
				return
			}
			b.cleanupExpired()
		case <-b.stopCleanup:
			return
		}
	}
}

// cleanupExpired removes expired cache items
func (b *Bucket[T]) cleanupExpired() {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	// Double-check if closed after acquiring lock
	if b.closed {
		return
	}

	now := time.Now()
	expiredKeys := make([]string, 0)

	// Collect expired keys
	for key, item := range b.cache {
		if now.After(item.expiredAt) {
			expiredKeys = append(expiredKeys, key)
		}
	}

	// Delete expired items
	for _, key := range expiredKeys {
		if item, exists := b.cache[key]; exists {
			b.updater.Remove(item)
			delete(b.cache, key)
		}
	}
}

// Close closes the bucket and stops the cleanup goroutine
// It's safe to call Close multiple times
func (b *Bucket[T]) Close() error {
	b.closeMutex.Lock()
	defer b.closeMutex.Unlock()

	// Check if already closed
	if b.closed {
		return nil // Already closed, no error
	}

	// Mark as closed
	b.closed = true

	// Stop the cleanup goroutine
	select {
	case b.stopCleanup <- struct{}{}:
		// Signal sent successfully
	default:
		// Channel is full or closed, that's fine
	}

	// Close the channel
	close(b.stopCleanup)

	// Clear all data from the bucket
	b.mutex.Lock()
	b.cache = make(map[string]*CacheItem[T])
	b.updater.Clear()
	b.mutex.Unlock()

	return nil
}

// isClosed checks if the bucket is closed (must be called with appropriate locking)
func (b *Bucket[T]) isClosed() bool {
	b.closeMutex.Lock()
	defer b.closeMutex.Unlock()
	return b.closed
}

// IsClosed returns whether the bucket is closed (public method)
func (b *Bucket[T]) IsClosed() bool {
	return b.isClosed()
}

// Size returns the current cache size
func (b *Bucket[T]) Size() int {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	if b.isClosed() {
		return 0
	}

	return b.updater.Size()
}

// Clear removes all cache items
func (b *Bucket[T]) Clear() {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if b.isClosed() {
		return
	}

	b.cache = make(map[string]*CacheItem[T])
	b.updater.Clear()
}

func WithBucketName[T any](name string) NewBucketOption[T] {
	return func(b *Bucket[T]) {
		b.name = name
	}
}

func WithBucketOutdated[T any](outdated time.Duration) NewBucketOption[T] {
	return func(b *Bucket[T]) {
		b.outdated = &outdated
	}
}

func WithMaxSize[T any](maxSize int) NewBucketOption[T] {
	return func(b *Bucket[T]) {
		b.maxSize = maxSize
	}
}

func WithCleanupInterval[T any](interval time.Duration) NewBucketOption[T] {
	return func(b *Bucket[T]) {
		b.cleanupInterval = interval
	}
}

// WithUpdater sets a custom update strategy
func WithUpdater[T any](updater Updater[T]) NewBucketOption[T] {
	return func(b *Bucket[T]) {
		b.updater = updater
	}
}

// WithFIFOUpdater sets FIFO update strategy
func WithFIFOUpdater[T any]() NewBucketOption[T] {
	return func(b *Bucket[T]) {
		b.updater = newFIFO[T]()
	}
}
