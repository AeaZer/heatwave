package heatwave

// SimpleFIFO strategy example - now much cleaner without prev/next coupling
type fifo[T any] struct {
	items []*CacheItem[T]
}

func newFIFO[T any]() *fifo[T] {
	return &fifo[T]{
		items: make([]*CacheItem[T], 0),
	}
}

func (f *fifo[T]) Add(item *CacheItem[T]) {
	f.items = append(f.items, item)
}

func (f *fifo[T]) Access(item *CacheItem[T]) {
	// FIFO doesn't reorder on access
}

func (f *fifo[T]) Remove(item *CacheItem[T]) {
	for i, it := range f.items {
		if it == item {
			f.items = append(f.items[:i], f.items[i+1:]...)
			break
		}
	}
}

func (f *fifo[T]) Evict() *CacheItem[T] {
	if len(f.items) == 0 {
		return nil
	}
	item := f.items[0]
	f.items = f.items[1:]
	return item
}

func (f *fifo[T]) Size() int {
	return len(f.items)
}

func (f *fifo[T]) Clear() {
	f.items = f.items[:0]
}
