package sync

import "sync"

type Slice[V any] struct {
	mu      sync.RWMutex
	content []V
}

func NewSlice[V any]() Slice[V] {
	return Slice[V]{
		sync.RWMutex{},
		[]V{},
	}
}

func (self *Slice[V]) Push(value V) {
	self.mu.Lock()
	self.content = append(self.content, value)
	self.mu.Unlock()
}

func (self *Slice[V]) Pop() V {
	self.mu.Lock()
	value := self.content[len(self.content)-1]
	self.content = self.content[:len(self.content)-1]
	self.mu.Unlock()
	return value
}

func (self *Slice[V]) Len() int {
	self.mu.RLock()
	l := len(self.content)
	self.mu.RUnlock()
	return l
}
