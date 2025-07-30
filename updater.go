package heatwave

// Updater interface defines the update strategy for cache items
type Updater[T any] interface {
	// Add adds a new item to the update strategy
	Add(item *CacheItem[T])
	// Access marks an item as accessed, updating its position in the strategy
	Access(item *CacheItem[T])
	// Remove removes an item from the update strategy
	Remove(item *CacheItem[T])
	// Evict returns the item that should be evicted according to the strategy
	Evict() *CacheItem[T]
	// Size returns the current size of items managed by the updater
	Size() int
	// Clear removes all items from the updater
	Clear()
}
