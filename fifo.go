package heatwave

// fifo implements FIFO (First-In-First-Out) algorithm
type fifo[T any] struct {
	items []*CacheItem[T]
}

// newFIFO creates a new FIFO updater
func newFIFO[T any]() *fifo[T] {
	return &fifo[T]{
		items: make([]*CacheItem[T], 0),
	}
}

// Add adds a new item to the FIFO updater
func (f *fifo[T]) Add(item *CacheItem[T]) {
	f.items = append(f.items, item)
}

// Access does nothing in FIFO strategy (no reordering on access)
func (f *fifo[T]) Access(item *CacheItem[T]) {
	// FIFO doesn't reorder on access
}

// Remove removes an item from the FIFO updater
func (f *fifo[T]) Remove(item *CacheItem[T]) {
	for i, it := range f.items {
		if it == item {
			f.items = append(f.items[:i], f.items[i+1:]...)
			break
		}
	}
}

// Evict returns the first item (oldest) for eviction
func (f *fifo[T]) Evict() *CacheItem[T] {
	if len(f.items) == 0 {
		return nil
	}
	item := f.items[0]
	f.items = f.items[1:]
	return item
}

// Size returns the current size
func (f *fifo[T]) Size() int {
	return len(f.items)
}

// Clear removes all items from the updater
func (f *fifo[T]) Clear() {
	f.items = f.items[:0]
}
