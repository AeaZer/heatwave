# 🔥 Heatwave - High-Performance Generic Memory Cache

[![Go Version](https://img.shields.io/badge/Go-1.18+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Generic](https://img.shields.io/badge/Generic-Type%20Safe-brightgreen.svg)](https://go.dev/doc/tutorial/generics)

> 🌐 **Language**: [English](README.md) | [中文](README_zh.md)

Heatwave is a blazing-fast, type-safe Go memory cache system with **full generic support**. It features pluggable eviction strategies (LRU, FIFO, Custom), automatic expiration, and thread-safe operations - all with compile-time type safety!

## ✨ Key Features

- 🎯 **Full Generic Type Support** - Work with any type `T` with compile-time safety
- 🚀 **High Performance** - O(1) operations with hash table + doubly linked list
- 🔄 **Pluggable Eviction Strategies** - LRU (default), FIFO, Random, or custom
- ⏰ **Auto Expiration** - TTL support with background cleanup
- 🔒 **Thread Safe** - Concurrent read/write with RWMutex
- 🎛️ **Highly Configurable** - Size limits, cleanup intervals, custom strategies
- 📦 **Zero Dependencies** - Pure Go implementation
- 💡 **Developer Friendly** - Intuitive API with excellent type inference

## 🚀 Quick Start

### Installation

```bash
go get github.com/AeaZer/heatwave
```

### Basic Usage

```go
package main

import (
    "fmt"
    "time"
    "github.com/AeaZer/heatwave"
)

func main() {
    // Create a string cache
    cache := heatwave.NewBucket[string]()
    defer cache.Close()

    // Store data
    cache.Nail("user:123", "Alice Johnson")
    
    // Retrieve data
    if name, found := cache.Bring("user:123"); found {
        fmt.Printf("Hello, %s!\n", name) // Hello, Alice Johnson!
    }
}
```

## 🎯 Type Safety in Action

### Compile-Time Type Checking

```go
// ✅ Type-safe operations
stringCache := heatwave.NewBucket[string]()
stringCache.Nail("key", "Hello World") // ✅ string accepted

intCache := heatwave.NewBucket[int]()
intCache.Nail("count", 42) // ✅ int accepted

// ❌ This won't compile - type mismatch!
// stringCache.Nail("key", 123) // ❌ Compile error
```

### No Type Assertions Needed

```go
// Before (traditional interface{} cache)
value := cache.Get("key")
if str, ok := value.(string); ok {  // Runtime type assertion
    fmt.Println(strings.ToUpper(str))
}

// After (Heatwave generics)
if value, found := cache.Bring("key"); found {
    fmt.Println(strings.ToUpper(value)) // Direct usage, type guaranteed!
}
```

## 🏗️ Core Concepts

| Concept | Description |
|---------|-------------|
| **Nail** | "Nail" data into memory (store operation) |
| **Bring** | "Bring" data from cache (retrieve operation) |
| **Bucket** | Generic cache container managing typed items |
| **Updater** | Pluggable eviction strategy interface |

## 📊 Supported Types

### Primitive Types
```go
// String cache
stringCache := heatwave.NewBucket[string]()
stringCache.Nail("name", "Alice")

// Numeric caches
intCache := heatwave.NewBucket[int]()
floatCache := heatwave.NewBucket[float64]()
boolCache := heatwave.NewBucket[bool]()
```

### Complex Types
```go
// Slice cache
sliceCache := heatwave.NewBucket[[]byte]()
sliceCache.Nail("data", []byte("binary data"))

// Map cache
mapCache := heatwave.NewBucket[map[string]int]()
mapCache.Nail("scores", map[string]int{"alice": 100, "bob": 85})

// Custom struct cache
type User struct {
    ID       int       `json:"id"`
    Name     string    `json:"name"`
    Email    string    `json:"email"`
    Created  time.Time `json:"created"`
}

userCache := heatwave.NewBucket[User]()
userCache.Nail("user:123", User{
    ID:      123,
    Name:    "Alice Johnson",
    Email:   "alice@example.com",
    Created: time.Now(),
})
```

### Interface Types for Mixed Data
```go
// Mixed type cache using interface{}
mixedCache := heatwave.NewBucket[interface{}]()
mixedCache.Nail("string", "Hello")
mixedCache.Nail("number", 42)
mixedCache.Nail("user", User{ID: 1, Name: "Alice"})

// Type assertion still needed for interface{} values
if value, found := mixedCache.Bring("string"); found {
    if str, ok := value.(string); ok {
        fmt.Println(str)
    }
}
```

## ⚙️ Configuration

### Basic Configuration

```go
cache := heatwave.NewBucket[string](
    heatwave.WithBucketName[string]("user-sessions"),
    heatwave.WithMaxSize[string](10000),                    // Max 10K items
    heatwave.WithBucketOutdated[string](time.Hour),         // 1 hour TTL
    heatwave.WithCleanupInterval[string](time.Minute * 5),  // Clean every 5min
)
```

### Advanced Configuration with Custom Strategy

```go
// Custom eviction strategy (see Custom Strategies section)
customUpdater := NewMyCustomUpdater[string]()

cache := heatwave.NewBucket[string](
    heatwave.WithBucketName[string]("high-priority-cache"),
    heatwave.WithMaxSize[string](5000),
    heatwave.WithUpdater[string](customUpdater),
    heatwave.WithBucketOutdated[string](time.Minute * 30),
)
```

## 🔄 Eviction Strategies

### Built-in Strategies

#### LRU (Least Recently Used) - Default

```go
// LRU is the default strategy
cache := heatwave.NewBucket[string]()

// Explicit LRU (same as default)
lruCache := heatwave.NewBucket[string](
    heatwave.WithUpdater[string](heatwave.newLRUUpdater[string]()),
)
```

#### FIFO (First In, First Out) - Built-in

```go
// Use built-in FIFO strategy
fifoCache := heatwave.NewBucket[string](
    heatwave.WithFIFOUpdater[string](),
)

// Or explicitly specify
fifoCache := heatwave.NewBucket[string](
    heatwave.WithUpdater[string](heatwave.newFIFO[string]()),
)
```

### Custom Strategies

Implement the `Updater[T]` interface:

```go
type MyCustomStrategy[T any] struct {
    items []*heatwave.CacheItem[T]
}

func (f *MyCustomStrategy[T]) Add(item *heatwave.CacheItem[T]) {
    f.items = append(f.items, item)
}

func (f *MyCustomStrategy[T]) Access(item *heatwave.CacheItem[T]) {
    // Custom access logic
}

func (f *MyCustomStrategy[T]) Remove(item *heatwave.CacheItem[T]) {
    for i, it := range f.items {
        if it == item {
            f.items = append(f.items[:i], f.items[i+1:]...)
            break
        }
    }
}

func (f *MyCustomStrategy[T]) Evict() *heatwave.CacheItem[T] {
    if len(f.items) == 0 {
        return nil
    }
    // Custom eviction logic
    item := f.items[0]
    f.items = f.items[1:]
    return item
}

func (f *MyCustomStrategy[T]) Size() int {
    return len(f.items)
}

func (f *MyCustomStrategy[T]) Clear() {
    f.items = f.items[:0]
}

// Usage
customCache := heatwave.NewBucket[string](
    heatwave.WithUpdater[string](&MyCustomStrategy[string]{}),
)
```

### Strategy Comparison

| Strategy | Eviction Rule | Use Cases | Time Complexity |
|----------|---------------|-----------|----------------|
| **LRU** | Least recently used | High locality access patterns | O(1) |
| **FIFO** | First in, first out | Time-series data, fair eviction | O(1) |
| **Custom** | User-defined logic | Special business requirements | Depends on implementation |

## 📖 Complete API Reference

### Bucket[T] Methods

| Method | Signature | Description |
|--------|-----------|-------------|
| `Nail` | `(id string, data T) error` | Store data with key |
| `Bring` | `(id string) (T, bool)` | Retrieve data by key |
| `Size` | `() int` | Current cache size |
| `Clear` | `()` | Remove all items |
| `Close` | `()` | Stop background cleanup |

### Configuration Options

| Option | Type | Description |
|--------|------|-------------|
| `WithBucketName[T]` | `string` | Set cache name |
| `WithMaxSize[T]` | `int` | Maximum cache size |
| `WithBucketOutdated[T]` | `time.Duration` | TTL for items |
| `WithCleanupInterval[T]` | `time.Duration` | Cleanup frequency |
| `WithUpdater[T]` | `Updater[T]` | Custom eviction strategy |
| `WithFIFOUpdater[T]` | `none` | Use built-in FIFO strategy |

### Updater[T] Interface

```go
type Updater[T any] interface {
    Add(item *CacheItem[T])     // Add new item
    Access(item *CacheItem[T])  // Mark item as accessed
    Remove(item *CacheItem[T])  // Remove specific item
    Evict() *CacheItem[T]       // Evict item (strategy-dependent)
    Size() int                  // Current size
    Clear()                     // Clear all items
}
```

## 🔄 Migration Guide

### From Non-Generic Version
```go
// For byte data
byteCache := heatwave.NewBucket[[]byte]()
byteCache.Nail("key", []byte("data"))
if data, found := byteCache.Bring("key"); found {
    fmt.Printf("Data: %s\n", string(data))
}

// Better: use string cache directly
stringCache := heatwave.NewBucket[string]()
stringCache.Nail("key", "data")
if data, found := stringCache.Bring("key"); found {
    fmt.Printf("Data: %s\n", data) // No conversion needed!
}
```

## 🎯 Performance & Benchmarks

### Time Complexity
- **Storage (Nail)**: O(1)
- **Retrieval (Bring)**: O(1)  
- **Eviction**: O(1) for LRU
- **Space**: O(n) where n = cache size

### Concurrency
- **Thread-Safe**: Uses `sync.RWMutex`
- **Multiple Readers**: Concurrent reads supported
- **Single Writer**: Writes are exclusive
- **Background Cleanup**: Non-blocking goroutine

## 🛡️ Thread Safety

All operations are thread-safe:

```go
cache := heatwave.NewBucket[string]()

// Safe concurrent access
go func() {
    cache.Nail("key1", "value1")
    cache.Nail("key2", "value2")
}()

go func() {
    if val, found := cache.Bring("key1"); found {
        fmt.Println("Found:", val)
    }
}()

go func() {
    fmt.Println("Cache size:", cache.Size())
}()
```

## 📋 Default Configuration

| Setting | Default Value | Description |
|---------|---------------|-------------|
| **Max Size** | 1,000 items | Maximum cache capacity |
| **TTL** | 5 minutes | Item expiration time |
| **Cleanup Interval** | 1 minute | Background cleanup frequency |
| **Strategy** | LRU | Default eviction strategy |

## 🔧 Requirements

- **Go Version**: 1.18+ (generics support required)
- **Dependencies**: None (pure Go)

## 🧪 Testing

```bash
# Run tests
go test -v

# Run benchmarks  
go test -bench=.

# Run with race detection
go test -race -v
```

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

### Development Setup

```bash
git clone https://github.com/AeaZer/heatwave.git
cd heatwave
go mod download
go test -v
```

---

<div align="center">
  <sub>Built with ❤️ for the Go community</sub>
</div> 