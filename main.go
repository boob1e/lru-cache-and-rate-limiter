package main

func main() {

	// cache := NewLRUCache[string](5)
}

type Node[T any] struct {
	CacheKey string
	Data     T
	Next     *Node[T]
	Previous *Node[T]
}

type LinkedList[T any] struct {
	Head *Node[T]
	Tail *Node[T]
}

type LRUCache[T any] struct {
	Cache       map[string]*Node[T]
	LinkedList  *LinkedList[T]
	CurrentSize int
	Size        int
}

func NewLRUCache[T any](size int) *LRUCache[T] {
	return &LRUCache[T]{
		Cache:       make(map[string]*Node[T], size),
		LinkedList:  &LinkedList[T]{},
		Size:        size,
		CurrentSize: 0,
	}
}

func (c *LRUCache[T]) Lookup(key string) (T, bool) {
	value, exist := c.Cache[key]
	if !exist {
		var zero T
		return zero, exist
	}
	if c.LinkedList.Head != value && c.CurrentSize > 1 {
		c.moveValueToHead(value)
	}
	return value.Data, exist
}

func (c *LRUCache[T]) Add(key string, value T) {
	if c.CurrentSize >= c.Size {
		// remove oldest value first
		oldTail := c.trimFromCurrentTail()
		delete(c.Cache, oldTail.CacheKey)
	}

	val, exists := c.Cache[key]
	if exists {
		c.Cache[key].Data = value
		c.moveValueToHead(val)
		return
	}

	node := &Node[T]{
		CacheKey: key,
		Data:     value,
		Next:     nil,
		Previous: nil,
	}
	c.moveValueToHead(node)
	c.Cache[key] = node
	c.CurrentSize++

	if c.LinkedList.Tail == nil {
		c.LinkedList.Tail = node
	}
}

func (c *LRUCache[T]) remove(key string) {

}

func (c *LRUCache[T]) moveValueToHead(value *Node[T]) {

	// if exists, move it to the front (head)
	prevNode := value.Previous
	nextNode := value.Next
	needsNewTail := value == c.LinkedList.Tail

	if needsNewTail {
		c.LinkedList.Tail = value.Previous
	}

	// if its not the tail
	if nextNode != nil {
		// remove references to current value
		nextNode.Previous = prevNode
	}

	// if its not the head
	if prevNode != nil {
		prevNode.Next = nextNode
	}

	// new head doesnt have a previous value
	value.Previous = nil

	// make old head the next value
	oldHead := c.LinkedList.Head
	if oldHead != nil {
		oldHead.Previous = value
	}
	value.Next = oldHead
	c.LinkedList.Head = value
}

func (c *LRUCache[T]) trimFromCurrentTail() *Node[T] {
	tail := c.LinkedList.Tail
	newTail := tail.Previous

	if newTail != nil {
		newTail.Next = nil
		c.CurrentSize--
	}

	if c.CurrentSize == 0 {
		c.LinkedList.Head = nil
	}

	return tail
}
