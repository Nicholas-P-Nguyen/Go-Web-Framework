package cache

// Node is a doubly linked list that contains a reference to the next and prev node
type Node struct {
	key     string
	value   []byte
	storage int
	next    *Node
	prev    *Node
}

// An LRU is a fixed-size in-memory cache with least-recently-used eviction
type LRU struct {
	limit       int
	currentSize int
	stats       Stats
	cache       map[string]*Node
	dummyLeft   Node
	dummyRight  Node
}

// NewLru returns a pointer to a new LRU with a capacity to store limit bytes
func NewLru(limit int) *LRU {
	lru := &LRU{
		limit:       limit,
		currentSize: 0,
		stats: Stats{
			Hits:   0,
			Misses: 0,
		},
		cache: map[string]*Node{},
		// dummyLeft.next will point to the item that is least recently used
		dummyLeft: Node{
			key:     "",
			value:   nil,
			storage: 0,
			next:    nil,
			prev:    nil,
		},
		// dummyRight.prev will point to the item that was most recently used
		dummyRight: Node{
			key:     "",
			value:   nil,
			storage: 0,
			next:    nil,
			prev:    nil,
		},
	}

	lru.dummyLeft.next = &lru.dummyRight
	lru.dummyRight.prev = &lru.dummyLeft

	return lru
}

// MaxStorage returns the maximum number of bytes this LRU can store
func (lru *LRU) MaxStorage() int {
	return lru.limit
}

// RemainingStorage returns the number of unused bytes available in this LRU
func (lru *LRU) RemainingStorage() int {
	return lru.limit - lru.currentSize
}

// InsertNode helper function that inserts node to the end of the doubly linked list
func (lru *LRU) InsertNode(node *Node) {
	mruNode := lru.dummyRight.prev

	node.prev, mruNode.next = mruNode, node
	node.next, lru.dummyRight.prev = &lru.dummyRight, node
}

// RemoveNode helper function that removes the specified node
func (lru *LRU) RemoveNode(node *Node) {
	prevNode := node.prev
	nextNode := node.next
	prevNode.next, nextNode.prev = nextNode, prevNode
}

// Remove removes and returns the value associated with the given key, if it exists.
// ok is true if a value was found and false otherwise
func (lru *LRU) Remove(key string) (value []byte, ok bool) {
	node, ok := lru.cache[key]
	if ok {
		value = node.value
		lru.currentSize -= node.storage
		lru.RemoveNode(node)
		delete(lru.cache, node.key)
		return value, ok
	}
	return value, ok
}

// Get returns the value associated with the given key, if it exists.
// This operation counts as a "use" for that key-value pair
// ok is true if a value was found and false otherwise.
func (lru *LRU) Get(key string) (value []byte, ok bool) {
	node, ok := lru.cache[key]
	if ok {
		lru.RemoveNode(node)
		lru.InsertNode(node)
		value = lru.cache[key].value
		lru.stats.Hits += 1
		return value, ok
	}
	lru.stats.Misses += 1
	return value, ok
}

// Set associates the given value with the given key, possibly evicting values
// to make room. Returns true if the binding was added successfully, else false.
func (lru *LRU) Set(key string, value []byte) bool {
	// Reject binding if it's too large for the LRU cache
	bindingSize := len(key) + len(value)
	if bindingSize > lru.limit {
		return false
	}

	// Checking for existing binding, if ok, remove it so we can update its last used
	node, ok := lru.cache[key]
	if ok {
		lru.currentSize -= node.storage
		lru.RemoveNode(node)
	}

	newNode := &Node{
		key:     key,
		value:   value,
		storage: bindingSize,
	}
	lru.cache[key] = newNode
	lru.currentSize += lru.cache[key].storage

	// Evicting LRU node if cache is full
	for lru.currentSize > lru.limit {
		lruNode := lru.dummyLeft.next
		lru.currentSize -= lruNode.storage
		lru.RemoveNode(lruNode)
		delete(lru.cache, lruNode.key)
	}

	node, ok = lru.cache[key]
	if ok {
		lru.InsertNode(node)
		return ok
	}
	return ok
}

// Len returns the number of bindings in the LRU.
func (lru *LRU) Len() int {
	return len(lru.cache)

}

// Stats returns statistics about how many search hits and misses have occurred.
func (lru *LRU) Stats() *Stats {
	return &lru.stats
}
