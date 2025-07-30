# 🔥 Heatwave[热浪] - 高性能泛型内存缓存

[![Go Version](https://img.shields.io/badge/Go-1.18+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Generic](https://img.shields.io/badge/Generic-Type%20Safe-brightgreen.svg)](https://go.dev/doc/tutorial/generics)

Heatwave 是一个高性能、类型安全的 Go 内存缓存系统，**完全支持泛型**。它具有可插拔的淘汰策略（LRU、FIFO、自定义）、自动过期和线程安全操作 - 全部都具有编译时类型安全性！

## ✨ 核心特性

- 🎯 **完整的泛型类型支持** - 支持任意类型 `T`，编译时安全保证
- 🚀 **高性能** - 基于哈希表 + 双向链表的 O(1) 操作
- 🔄 **可插拔淘汰策略** - LRU（默认）、FIFO、随机或自定义策略
- ⏰ **自动过期** - TTL 支持和后台清理
- 🔒 **线程安全** - 使用 RWMutex 支持并发读写
- 🎛️ **高度可配置** - 大小限制、清理间隔、自定义策略
- 📦 **零依赖** - 纯 Go 实现
- 💡 **开发者友好** - 直观的 API 和优秀的类型推断

## 🚀 快速开始

### 安装

```bash
go get github.com/AeaZer/heatwave
```

### 基本用法

```go
package main

import (
    "fmt"
    "time"
    "github.com/AeaZer/heatwave"
)

func main() {
    // 创建字符串缓存
    cache := heatwave.NewBucket[string]()
    defer cache.Close()

    // 存储数据
    cache.Nail("user:123", "Alice Johnson")
    
    // 获取数据
    if name, found := cache.Bring("user:123"); found {
        fmt.Printf("你好，%s！\n", name) // 你好，Alice Johnson！
    }
}
```

## 🎯 类型安全实践

### 编译时类型检查

```go
// ✅ 类型安全操作
stringCache := heatwave.NewBucket[string]()
stringCache.Nail("key", "Hello World") // ✅ 接受字符串类型

intCache := heatwave.NewBucket[int]()
intCache.Nail("count", 42) // ✅ 接受整数类型

// ❌ 这样写不会编译通过 - 类型不匹配！
// stringCache.Nail("key", 123) // ❌ 编译错误
```

### 无需类型断言

```go
// 之前（传统的 interface{} 缓存）
value := cache.Get("key")
if str, ok := value.(string); ok {  // 运行时类型断言
    fmt.Println(strings.ToUpper(str))
}

// 现在（Heatwave 泛型）
if value, found := cache.Bring("key"); found {
    fmt.Println(strings.ToUpper(value)) // 直接使用，类型有保证！
}
```

## 🏗️ 核心概念

| 概念 | 描述 |
|------|------|
| **Nail** | 将数据"钉"在内存中（存储操作） |
| **Bring** | 从缓存中"取出"数据（获取操作） |
| **Bucket** | 管理类型化项目的泛型缓存容器 |
| **Updater** | 可插拔的淘汰策略接口 |

## 📊 支持的类型

### 基本类型
```go
// 字符串缓存
stringCache := heatwave.NewBucket[string]()
stringCache.Nail("name", "Alice")

// 数值缓存
intCache := heatwave.NewBucket[int]()
floatCache := heatwave.NewBucket[float64]()
boolCache := heatwave.NewBucket[bool]()
```

### 复杂类型
```go
// 切片缓存
sliceCache := heatwave.NewBucket[[]byte]()
sliceCache.Nail("data", []byte("二进制数据"))

// 映射缓存
mapCache := heatwave.NewBucket[map[string]int]()
mapCache.Nail("scores", map[string]int{"alice": 100, "bob": 85})

// 自定义结构体缓存
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

### 混合数据的接口类型
```go
// 使用 interface{} 的混合类型缓存
mixedCache := heatwave.NewBucket[interface{}]()
mixedCache.Nail("string", "Hello")
mixedCache.Nail("number", 42)
mixedCache.Nail("user", User{ID: 1, Name: "Alice"})

// interface{} 值仍需要类型断言
if value, found := mixedCache.Bring("string"); found {
    if str, ok := value.(string); ok {
        fmt.Println(str)
    }
}
```

## ⚙️ 配置

### 基本配置

```go
cache := heatwave.NewBucket[string](
    heatwave.WithBucketName[string]("user-sessions"),
    heatwave.WithMaxSize[string](10000),                    // 最大 1万 项目
    heatwave.WithBucketOutdated[string](time.Hour),         // 1小时 TTL
    heatwave.WithCleanupInterval[string](time.Minute * 5),  // 每5分钟清理
)
```

### 使用自定义策略的高级配置

```go
// 自定义淘汰策略（参见自定义策略章节）
customUpdater := NewMyCustomUpdater[string]()

cache := heatwave.NewBucket[string](
    heatwave.WithBucketName[string]("high-priority-cache"),
    heatwave.WithMaxSize[string](5000),
    heatwave.WithUpdater[string](customUpdater),
    heatwave.WithBucketOutdated[string](time.Minute * 30),
)
```

## 🔄 淘汰策略

### 内置 LRU（默认）

```go
// LRU 是默认策略
cache := heatwave.NewBucket[string]()

// 显式指定 LRU（与默认相同）
lruCache := heatwave.NewBucket[string](
    heatwave.WithUpdater[string](heatwave.newLRUUpdater[string]()),
)
```

### 自定义策略

实现 `Updater[T]` 接口：

```go
type MyFIFOStrategy[T any] struct {
    items []*heatwave.CacheItem[T]
}

func (f *MyFIFOStrategy[T]) Add(item *heatwave.CacheItem[T]) {
    f.items = append(f.items, item)
}

func (f *MyFIFOStrategy[T]) Access(item *heatwave.CacheItem[T]) {
    // FIFO 访问时不重新排序
}

func (f *MyFIFOStrategy[T]) Remove(item *heatwave.CacheItem[T]) {
    for i, it := range f.items {
        if it == item {
            f.items = append(f.items[:i], f.items[i+1:]...)
            break
        }
    }
}

func (f *MyFIFOStrategy[T]) Evict() *heatwave.CacheItem[T] {
    if len(f.items) == 0 {
        return nil
    }
    item := f.items[0]
    f.items = f.items[1:]
    return item
}

func (f *MyFIFOStrategy[T]) Size() int {
    return len(f.items)
}

func (f *MyFIFOStrategy[T]) Clear() {
    f.items = f.items[:0]
}

// 使用方法
fifoCache := heatwave.NewBucket[string](
    heatwave.WithUpdater[string](&MyFIFOStrategy[string]{}),
)
```

## 📖 完整 API 参考

### Bucket[T] 方法

| 方法 | 签名 | 描述 |
|------|------|------|
| `Nail` | `(id string, data T) error` | 使用键存储数据 |
| `Bring` | `(id string) (T, bool)` | 通过键获取数据 |
| `Size` | `() int` | 当前缓存大小 |
| `Clear` | `()` | 移除所有项目 |
| `Close` | `()` | 停止后台清理 |

### 配置选项

| 选项 | 类型 | 描述 |
|------|------|------|
| `WithBucketName[T]` | `string` | 设置缓存名称 |
| `WithMaxSize[T]` | `int` | 最大缓存大小 |
| `WithBucketOutdated[T]` | `time.Duration` | 项目 TTL |
| `WithCleanupInterval[T]` | `time.Duration` | 清理频率 |
| `WithUpdater[T]` | `Updater[T]` | 自定义淘汰策略 |

### Updater[T] 接口

```go
type Updater[T any] interface {
    Add(item *CacheItem[T])     // 添加新项目
    Access(item *CacheItem[T])  // 标记项目被访问
    Remove(item *CacheItem[T])  // 移除特定项目
    Evict() *CacheItem[T]       // 淘汰项目（策略相关）
    Size() int                  // 当前大小
    Clear()                     // 清除所有项目
}
```

## 🔄 迁移指南

### 从非泛型版本迁移

```go
// 对于字节数据
byteCache := heatwave.NewBucket[[]byte]()
byteCache.Nail("key", []byte("data"))
if data, found := byteCache.Bring("key"); found {
    fmt.Printf("Data: %s\n", string(data))
}

// 更好的方式：直接使用字符串缓存
stringCache := heatwave.NewBucket[string]()
stringCache.Nail("key", "data")
if data, found := stringCache.Bring("key"); found {
    fmt.Printf("Data: %s\n", data) // 无需转换！
}
```

## 🎯 性能与基准测试

### 时间复杂度
- **存储 (Nail)**: O(1)
- **获取 (Bring)**: O(1)  
- **淘汰**: LRU 为 O(1)
- **空间**: O(n)，其中 n = 缓存大小

### 并发性
- **线程安全**: 使用 `sync.RWMutex`
- **多读取器**: 支持并发读取
- **单写入器**: 写入是独占的
- **后台清理**: 非阻塞 goroutine

## 🛡️ 线程安全

所有操作都是线程安全的：

```go
cache := heatwave.NewBucket[string]()

// 安全的并发访问
go func() {
    cache.Nail("key1", "value1")
    cache.Nail("key2", "value2")
}()

go func() {
    if val, found := cache.Bring("key1"); found {
        fmt.Println("找到:", val)
    }
}()

go func() {
    fmt.Println("缓存大小:", cache.Size())
}()
```

## 📋 默认配置

| 设置 | 默认值 | 描述 |
|------|--------|------|
| **最大大小** | 1,000 项目 | 最大缓存容量 |
| **TTL** | 5 分钟 | 项目过期时间 |
| **清理间隔** | 1 分钟 | 后台清理频率 |
| **策略** | LRU | 默认淘汰策略 |

## 🔧 系统要求

- **Go 版本**: 1.18+（需要泛型支持）
- **依赖**: 无（纯 Go）

## 🧪 测试

```bash
# 运行测试
go test -v

# 运行基准测试
go test -bench=.

# 使用竞态检测运行
go test -race -v
```

## 📄 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件。

### 开发设置

```bash
git clone https://github.com/AeaZer/heatwave.git
cd heatwave
go mod download
go test -v
```

---

<div align="center">
  <sub>用 ❤️ 为 Go 社区构建</sub>
</div> 