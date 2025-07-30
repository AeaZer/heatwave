package main

import (
	"fmt"
	"time"

	"github.com/AeaZer/heatwave"
)

// RandomStrategy - also simpler now
type randomStrategy[T any] struct {
	items []*heatwave.CacheItem[T]
}

func newRandomStrategy[T any]() *randomStrategy[T] {
	return &randomStrategy[T]{
		items: make([]*heatwave.CacheItem[T], 0),
	}
}

func (r *randomStrategy[T]) Add(item *heatwave.CacheItem[T]) {
	r.items = append(r.items, item)
}

func (r *randomStrategy[T]) Access(item *heatwave.CacheItem[T]) {
	// Random strategy doesn't care about access patterns
}

func (r *randomStrategy[T]) Remove(item *heatwave.CacheItem[T]) {
	for i, it := range r.items {
		if it == item {
			r.items = append(r.items[:i], r.items[i+1:]...)
			break
		}
	}
}

func (r *randomStrategy[T]) Evict() *heatwave.CacheItem[T] {
	if len(r.items) == 0 {
		return nil
	}
	// Simple random selection (take last for simplicity)
	idx := len(r.items) - 1
	item := r.items[idx]
	r.items = r.items[:idx]
	return item
}

func (r *randomStrategy[T]) Size() int {
	return len(r.items)
}

func (r *randomStrategy[T]) Clear() {
	r.items = r.items[:0]
}

// FrequencyBasedStrategy - tracks access frequency
type frequencyStrategy[T any] struct {
	items     []*heatwave.CacheItem[T]
	frequency map[*heatwave.CacheItem[T]]int
}

func newFrequencyStrategy[T any]() *frequencyStrategy[T] {
	return &frequencyStrategy[T]{
		items:     make([]*heatwave.CacheItem[T], 0),
		frequency: make(map[*heatwave.CacheItem[T]]int),
	}
}

func (f *frequencyStrategy[T]) Add(item *heatwave.CacheItem[T]) {
	f.items = append(f.items, item)
	f.frequency[item] = 0
}

func (f *frequencyStrategy[T]) Access(item *heatwave.CacheItem[T]) {
	if _, exists := f.frequency[item]; exists {
		f.frequency[item]++
	}
}

func (f *frequencyStrategy[T]) Remove(item *heatwave.CacheItem[T]) {
	for i, it := range f.items {
		if it == item {
			f.items = append(f.items[:i], f.items[i+1:]...)
			break
		}
	}
	delete(f.frequency, item)
}

func (f *frequencyStrategy[T]) Evict() *heatwave.CacheItem[T] {
	if len(f.items) == 0 {
		return nil
	}

	// Find item with lowest frequency
	var leastFreqItem *heatwave.CacheItem[T]
	minFreq := int(^uint(0) >> 1) // Max int
	leastFreqIndex := -1

	for i, item := range f.items {
		if freq, exists := f.frequency[item]; exists {
			if freq < minFreq {
				minFreq = freq
				leastFreqItem = item
				leastFreqIndex = i
			}
		}
	}

	if leastFreqIndex >= 0 {
		f.items = append(f.items[:leastFreqIndex], f.items[leastFreqIndex+1:]...)
		delete(f.frequency, leastFreqItem)
		return leastFreqItem
	}

	return nil
}

func (f *frequencyStrategy[T]) Size() int {
	return len(f.items)
}

func (f *frequencyStrategy[T]) Clear() {
	f.items = f.items[:0]
	f.frequency = make(map[*heatwave.CacheItem[T]]int)
}

func demonstrateCustomStrategies() {
	fmt.Println("=== Custom Strategy Demo ===")

	// Example 1: FIFO Strategy
	fmt.Println("\n1. FIFO Strategy Demo:")
	fifoCache := heatwave.NewBucket[string](
		heatwave.WithBucketName[string]("custom-fifo"),
		heatwave.WithMaxSize[string](3),
		heatwave.WithFIFOUpdater[string](),
		heatwave.WithBucketOutdated[string](time.Second*10),
	)
	defer fifoCache.Close()

	// Add items
	fifoCache.Nail("first", "value1")
	fifoCache.Nail("second", "value2")
	fifoCache.Nail("third", "value3")
	fmt.Printf("FIFO cache size: %d\n", fifoCache.Size())

	// Access first item (FIFO doesn't change order)
	if data, ok := fifoCache.Bring("first"); ok {
		fmt.Printf("Accessed first item: %s\n", data)
	}

	// Add fourth item, should evict "first" (oldest)
	fifoCache.Nail("fourth", "value4")
	fmt.Printf("After adding fourth item, size: %d\n", fifoCache.Size())

	if _, ok := fifoCache.Bring("first"); !ok {
		fmt.Println("✓ First item was evicted (FIFO behavior)")
	}

	// Example 2: Frequency-based Strategy
	fmt.Println("\n2. Frequency-based Strategy Demo:")
	freqCache := heatwave.NewBucket[int](
		heatwave.WithBucketName[int]("frequency-cache"),
		heatwave.WithMaxSize[int](3),
		heatwave.WithUpdater[int](newFrequencyStrategy[int]()),
		heatwave.WithBucketOutdated[int](time.Second*10),
	)
	defer freqCache.Close()

	// Add items
	freqCache.Nail("item1", 100)
	freqCache.Nail("item2", 200)
	freqCache.Nail("item3", 300)

	// Access item1 multiple times
	freqCache.Bring("item1")
	freqCache.Bring("item1")
	freqCache.Bring("item1")

	// Access item2 once
	freqCache.Bring("item2")

	// item3 is not accessed

	// Add new item, should evict item3 (least frequent)
	freqCache.Nail("item4", 400)

	if _, ok := freqCache.Bring("item3"); !ok {
		fmt.Println("✓ item3 was evicted (least frequent)")
	}
	if _, ok := freqCache.Bring("item1"); ok {
		fmt.Println("✓ item1 still exists (most frequent)")
	}

	fmt.Println("\n=== Custom Strategy Demo Complete ===")
}

func init() {
	// Uncomment to run the demo
	// demonstrateCustomStrategies()
}
