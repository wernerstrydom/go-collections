package collections

import (
	"fmt"
	"sync"
)



// Stack implements a LIFO data structure. It is not thread-safe.
type Stack[T any] struct {
	items []T
}

// NewStack returns a new stack with the given initial items.
func NewStack[T any](values ...T) *Stack[T] {
	return &Stack[T]{items: values}
}

// Push adds an item to the top of the stack.
func (s *Stack[T]) Push(item T) {
	s.items = append(s.items, item)
}

// Pop removes and returns the item at the top of the stack. If the stack is
// empty, an error is returned.
func (s *Stack[T]) Pop() (T, error) {
	var zero T
	if len(s.items) == 0 {
		return zero, ErrEmptyStack
	}

	item := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return item, nil
}

// Peek returns the item at the top of the stack without removing it. If the
// stack is empty, an error is returned.
func (s *Stack[T]) Peek() (T, error) {
	var zero T
	if len(s.items) == 0 {
		return zero, ErrEmptyStack
	}
	return s.items[len(s.items)-1], nil
}

// IsEmpty returns true if the stack is empty.
func (s *Stack[T]) IsEmpty() bool {
	return len(s.items) == 0
}

// Size returns the number of items in the stack.
func (s *Stack[T]) Size() int {
	return len(s.items)
}

// String returns a string representation of the stack.
func (s *Stack[T]) String() string {
	return fmt.Sprintf("%v", s.items)
}

// ConcurrentStack implements a LIFO data structure. It is thread-safe.
type ConcurrentStack[T any] struct {
	items []T
	lock  sync.RWMutex
}

// NewConcurrentStack returns a new stack with the given initial items.
func NewConcurrentStack[T any](values ...T) *ConcurrentStack[T] {
	return &ConcurrentStack[T]{items: values}
}

// Push adds an item to the top of the stack.
func (s *ConcurrentStack[T]) Push(item T) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.items = append(s.items, item)
}

// Pop removes and returns the item at the top of the stack. If the stack is
func (s *ConcurrentStack[T]) Pop() (T, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	var zero T
	if len(s.items) == 0 {
		return zero, ErrEmptyStack
	}

	item := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return item, nil
}

// Peek returns the item at the top of the stack without removing it. If the
func (s *ConcurrentStack[T]) Peek() (T, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	var zero T
	if len(s.items) == 0 {
		return zero, ErrEmptyStack
	}
	return s.items[len(s.items)-1], nil
}

// IsEmpty returns true if the stack is empty.
func (s *ConcurrentStack[T]) IsEmpty() bool {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return len(s.items) == 0
}

// Size returns the number of items in the stack.
func (s *ConcurrentStack[T]) Size() int {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return len(s.items)
}

// String returns a string representation of the stack.
func (s *ConcurrentStack[T]) String() string {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return fmt.Sprintf("%v", s.items)
}
