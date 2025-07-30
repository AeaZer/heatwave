package heatwave

// lruNode represents a node in the LRU doubly linked list
type lruNode[T any] struct {
	item *CacheItem[T]
	prev *lruNode[T]
	next *lruNode[T]
}

// lru implements LRU algorithm using doubly linked list
type lru[T any] struct {
	head    *lruNode[T]
	tail    *lruNode[T]
	size    int
	nodeMap map[*CacheItem[T]]*lruNode[T] // Map from CacheItem to lruNode for O(1) access
}

// newLRUUpdater creates a new lru updater
func newLRUUpdater[T any]() *lru[T] {
	head := &lruNode[T]{}
	tail := &lruNode[T]{}
	head.next = tail
	tail.prev = head
	return &lru[T]{
		head:    head,
		tail:    tail,
		size:    0,
		nodeMap: make(map[*CacheItem[T]]*lruNode[T]),
	}
}

// Add adds a new item to the lru updater
func (l *lru[T]) Add(item *CacheItem[T]) {
	node := &lruNode[T]{item: item}
	l.nodeMap[item] = node
	l.addNodeToHead(node)
}

// Access marks an item as accessed, moving it to head
func (l *lru[T]) Access(item *CacheItem[T]) {
	if node, exists := l.nodeMap[item]; exists {
		l.moveNodeToHead(node)
	}
}

// Remove removes an item from the lru updater
func (l *lru[T]) Remove(item *CacheItem[T]) {
	if node, exists := l.nodeMap[item]; exists {
		l.removeNode(node)
		delete(l.nodeMap, item)
	}
}

// Evict returns the least recently used item for eviction
func (l *lru[T]) Evict() *CacheItem[T] {
	return l.removeTail()
}

// Size returns the current size
func (l *lru[T]) Size() int {
	return l.size
}

// Clear removes all items from the updater
func (l *lru[T]) Clear() {
	head := &lruNode[T]{}
	tail := &lruNode[T]{}
	head.next = tail
	tail.prev = head
	l.head = head
	l.tail = tail
	l.size = 0
	l.nodeMap = make(map[*CacheItem[T]]*lruNode[T])
}

// addNodeToHead adds a node to the head of the list
func (l *lru[T]) addNodeToHead(node *lruNode[T]) {
	node.prev = l.head
	node.next = l.head.next
	l.head.next.prev = node
	l.head.next = node
	l.size++
}

// removeNode removes a node from the list
func (l *lru[T]) removeNode(node *lruNode[T]) {
	node.prev.next = node.next
	node.next.prev = node.prev
	l.size--
}

// removeTail removes the tail node and returns its item
func (l *lru[T]) removeTail() *CacheItem[T] {
	if l.size == 0 {
		return nil
	}
	lastNode := l.tail.prev
	item := lastNode.item
	l.removeNode(lastNode)
	delete(l.nodeMap, item)
	return item
}

// moveNodeToHead moves a node to the head of the list
func (l *lru[T]) moveNodeToHead(node *lruNode[T]) {
	l.removeNode(node)
	l.addNodeToHead(node)
}
