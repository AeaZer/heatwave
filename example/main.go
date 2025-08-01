package main

import (
	"fmt"
	"time"

	"github.com/AeaZer/heatwave"
)

func main() {
	fmt.Println("=== Heatwave Cache System Demo ===")

	// Example 1: Default lru Strategy with string values
	fmt.Println("\n1. Default lru Strategy:")
	lruCache := heatwave.NewBucket[string](
		heatwave.WithBucketName[string]("lru-cache"),
		heatwave.WithMaxSize[string](3),
		heatwave.WithBucketExpire[string](time.Second*10),
	)
	defer lruCache.Close()

	// Add items
	lruCache.Nail("user:1", `{"name": "Alice"}`)
	lruCache.Nail("user:2", `{"name": "Bob"}`)
	lruCache.Nail("user:3", `{"name": "Charlie"}`)
	fmt.Printf("lru Cache size after adding 3 items: %d\n", lruCache.Size())

	// Access user:1 to make it recently used
	if data, ok := lruCache.Bring("user:1"); ok {
		fmt.Printf("Accessed user:1: %s\n", data)
	}

	// Add new item, should evict user:2 (least recently used)
	lruCache.Nail("user:4", `{"name": "David"}`)
	fmt.Printf("lru Cache size after adding user:4: %d\n", lruCache.Size())

	// Check which item was evicted
	if _, ok := lruCache.Bring("user:2"); !ok {
		fmt.Println("✓ user:2 was evicted by lru strategy")
	}
	if data, ok := lruCache.Bring("user:1"); ok {
		fmt.Printf("✓ user:1 is still in cache (recently used): %s\n", data)
	}

	// Example 2: Working with different types - integers
	fmt.Println("\n2. Integer Cache Example:")
	intCache := heatwave.NewBucket[int](
		heatwave.WithBucketName[int]("int-cache"),
		heatwave.WithMaxSize[int](3),
		heatwave.WithBucketExpire[int](time.Second*10),
	)
	defer intCache.Close()

	// Add integer values
	intCache.Nail("score:1", 100)
	intCache.Nail("score:2", 85)
	intCache.Nail("score:3", 92)
	fmt.Printf("Integer cache size: %d\n", intCache.Size())

	if score, ok := intCache.Bring("score:1"); ok {
		fmt.Printf("Score 1: %d\n", score)
	}

	// Example 3: Working with struct types
	fmt.Println("\n3. Struct Cache Example:")
	type User struct {
		Name  string
		Age   int
		Email string
	}

	userCache := heatwave.NewBucket[User](
		heatwave.WithBucketName[User]("user-cache"),
		heatwave.WithMaxSize[User](2),
		heatwave.WithBucketExpire[User](time.Second*10),
	)
	defer userCache.Close()

	// Add struct values
	userCache.Nail("user:alice", User{Name: "Alice", Age: 25, Email: "alice@example.com"})
	userCache.Nail("user:bob", User{Name: "Bob", Age: 30, Email: "bob@example.com"})

	if user, ok := userCache.Bring("user:alice"); ok {
		fmt.Printf("User Alice: %+v\n", user)
	}

	// Example 4: Byte slice cache (similar to original)
	fmt.Println("\n4. Byte Slice Cache Example:")
	byteCache := heatwave.NewBucket[[]byte](
		heatwave.WithBucketName[[]byte]("byte-cache"),
		heatwave.WithMaxSize[[]byte](2),
		heatwave.WithBucketExpire[[]byte](time.Second*10),
	)
	defer byteCache.Close()

	byteCache.Nail("data:1", []byte("Hello, World!"))
	byteCache.Nail("data:2", []byte("Go Generics!"))

	if data, ok := byteCache.Bring("data:1"); ok {
		fmt.Printf("Byte data: %s\n", string(data))
	}

	// Example 5: Expiration Demo
	fmt.Println("\n5. Expiration Demo:")
	expCache := heatwave.NewBucket[string](
		heatwave.WithBucketName[string]("exp-cache"),
		heatwave.WithBucketExpire[string](time.Second*2), // 2 seconds TTL
	)
	defer expCache.Close()

	expCache.Nail("temp", "temporary data")
	if data, ok := expCache.Bring("temp"); ok {
		fmt.Printf("Before expiration: %s\n", data)
	}

	fmt.Println("Waiting 3 seconds for expiration...")
	time.Sleep(3 * time.Second)

	if _, ok := expCache.Bring("temp"); !ok {
		fmt.Println("✓ Data expired and was removed")
	}

	// Example 6: Interface{} cache for mixed types
	fmt.Println("\n6. Mixed Types Cache Example:")
	mixedCache := heatwave.NewBucket[interface{}](
		heatwave.WithBucketName[interface{}]("mixed-cache"),
		heatwave.WithMaxSize[interface{}](3),
	)
	defer mixedCache.Close()

	mixedCache.Nail("string", "Hello")
	mixedCache.Nail("number", 42)
	mixedCache.Nail("bool", true)

	if str, ok := mixedCache.Bring("string"); ok {
		fmt.Printf("String value: %v\n", str)
	}
	if num, ok := mixedCache.Bring("number"); ok {
		fmt.Printf("Number value: %v\n", num)
	}
	if b, ok := mixedCache.Bring("bool"); ok {
		fmt.Printf("Boolean value: %v\n", b)
	}

	fmt.Println("\n=== Demo Complete ===")
}
