package lru

import "container/list"

// Cache uses LRU algorithm to control the elimination of kv-pair
// Cache uses queue to control its data, when we access a key,
// move the corresponding element to the front, and when delete we
// directly delete the rear of the queue
type Cache struct {
	// the maximum memory
	maxBytes int64
	// memory currently used
	nbytes int64
	ll     *list.List
	cache  map[string]*list.Element
	// a callback function, optional and executed when an entry is purged
	onEvicted func(key string, value Value)
}

// entry is a key-value pair, which is the datatype of double linked list
type entry struct {
	key   string
	value Value
}

// Value uses Len to count how many bytes it takes to ensure the
// compatibility, we let Value can be any type realizing the interface
type Value interface {
	// Len return the memory usage
	Len() int
}

// New is the constructor of Cache
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		onEvicted: onEvicted,
	}
}

// Get looks up a key's value
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		// if the element exists, move it to the rear
		c.ll.MoveToFront(ele)
		// get the list element's kv pair
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

// RemoveOldest removes the oldest item
func (c *Cache) RemoveOldest() {
	ele := c.ll.Back()
	if ele != nil {
		// remove the queue's rear element
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		// delete the mapping relationship in the dictionary
		delete(c.cache, kv.key)
		// update current memory
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.onEvicted != nil {
			c.onEvicted(kv.key, kv.value)
		}
	}
}

func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		// if the key already exists, update the corresponding node's value
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		ele := c.ll.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.nbytes += int64(len(key)) + int64(value.Len())
	}

	// if current memory exceeds, remove the least recently used node
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOldest()
	}
}

func (c *Cache) Len() int{
	return c.ll.Len()
}