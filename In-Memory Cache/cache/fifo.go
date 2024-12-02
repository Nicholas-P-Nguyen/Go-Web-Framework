package cache

type FIFONode struct {
	key     string
	value   []byte
	storage int
	next    *FIFONode
	prev    *FIFONode
}

// An FIFO is a fixed-size in-memory cache with first-in first-out eviction
type FIFO struct {
	limit       int
	currentSize int
	stats       Stats
	cache       map[string]*FIFONode
	dummyLeft   FIFONode
	dummyRight  FIFONode
}

// NewFIFO returns a pointer to a new FIFO with a capacity to store limit bytes
func NewFifo(limit int) *FIFO {
	fifo := &FIFO{
		limit:       limit,
		currentSize: 0,
		stats: Stats{
			Hits:   0,
			Misses: 0,
		},
		cache: map[string]*FIFONode{},
		// dummyLeft.next will point to the oldest item that was inserted to the cache
		dummyLeft: FIFONode{
			key:     "",
			value:   nil,
			storage: 0,
			next:    nil,
			prev:    nil,
		},
		// dummyRight.prev will point to newest item that was inserted to the cache
		dummyRight: FIFONode{
			key:     "",
			value:   nil,
			storage: 0,
			next:    nil,
			prev:    nil,
		},
	}
	fifo.dummyLeft.next, fifo.dummyRight.prev = &fifo.dummyRight, &fifo.dummyLeft
	return fifo
}

// MaxStorage returns the maximum number of bytes this FIFO can store
func (fifo *FIFO) MaxStorage() int {
	return fifo.limit
}

// RemainingStorage returns the number of unused bytes available in this FIFO
func (fifo *FIFO) RemainingStorage() int {
	return fifo.limit - fifo.currentSize
}

// Get returns the value associated with the given key, if it exists.
// ok is true if a value was found and false otherwise.
func (fifo *FIFO) Get(key string) (value []byte, ok bool) {
	node, ok := fifo.cache[key]
	if ok {
		fifo.stats.Hits += 1
		return node.value, ok
	}
	fifo.stats.Misses += 1
	return value, ok
}

func (fifo *FIFO) RemoveNode(node *FIFONode) {
	prevNode := node.prev
	nextNode := node.next

	prevNode.next, nextNode.prev = nextNode, prevNode
}

func (fifo *FIFO) InsertNode(node *FIFONode) {
	tail := fifo.dummyRight.prev
	node.next, fifo.dummyRight.prev = &fifo.dummyRight, node
	node.prev, tail.next = tail, node
}

// Remove removes and returns the value associated with the given key, if it exists.
// ok is true if a value was found and false otherwise
func (fifo *FIFO) Remove(key string) (value []byte, ok bool) {
	node, ok := fifo.cache[key]
	if ok {
		value = node.value
		fifo.currentSize -= node.storage
		fifo.RemoveNode(node)
		delete(fifo.cache, key)
	}
	return value, ok
}

// Set associates the given value with the given key, possibly evicting values
// to make room. Returns true if the binding was added successfully, else false.
func (fifo *FIFO) Set(key string, value []byte) bool {
	bindingSize := len(key) + len(value)
	if bindingSize > fifo.limit {
		return false
	}
	// Checking for existing binding, if ok, update the existing binding
	node, ok := fifo.cache[key]
	if ok {
		fifo.currentSize -= node.storage
		node.value = value
		node.storage = bindingSize
		fifo.currentSize += bindingSize
		return true
	}

	newNode := &FIFONode{
		key:     key,
		value:   value,
		storage: bindingSize,
	}
	fifo.cache[key] = newNode
	fifo.currentSize += bindingSize

	// Evicting oldest node if cache if full
	for fifo.currentSize > fifo.limit {
		oldest := fifo.dummyLeft.next
		fifo.currentSize -= oldest.storage
		fifo.RemoveNode(oldest)
		delete(fifo.cache, oldest.key)
	}

	fifo.InsertNode(newNode)
	return true
}

// Len returns the number of bindings in the FIFO.
func (fifo *FIFO) Len() int {
	return len(fifo.cache)
}

// Stats returns statistics about how many search hits and misses have occurred.
func (fifo *FIFO) Stats() *Stats {
	return &fifo.stats
}
