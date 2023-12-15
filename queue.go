package collections

import (
	"fmt"
	"sync"
)

// Queue implements a FIFO data structure. It is not thread-safe.
type Queue[T any] struct {
	items []T
}

// NewQueue returns a new queue with the given initial items.
func NewQueue[T any](values ...T) *Queue[T] {
	return &Queue[T]{items: values}
}

// Enqueue adds an item to the end of the queue.
func (q *Queue[T]) Enqueue(item T) {
	q.items = append(q.items, item)
}

// Dequeue removes and returns the item at the front of the queue. If the queue
// is empty, an error is returned.
func (q *Queue[T]) Dequeue() (T, error) {
	var zero T
	if len(q.items) == 0 {
		return zero, ErrEmptyQueue
	}

	item := q.items[0]
	q.items = q.items[1:]
	return item, nil
}

// Peek returns the item at the front of the queue without removing it. If the
// queue is empty, an error is returned.
func (q *Queue[T]) Peek() (T, error) {
	var zero T
	if len(q.items) == 0 {
		return zero, ErrEmptyQueue
	}
	return q.items[0], nil
}

// IsEmpty returns true if the queue is empty.
func (q *Queue[T]) IsEmpty() bool {
	return len(q.items) == 0
}

// Size returns the number of items in the queue.
func (q *Queue[T]) Size() int {
	return len(q.items)
}

// String returns a string representation of the queue.
func (q *Queue[T]) String() string {
	return fmt.Sprintf("%v", q.items)
}

// Clear removes all items from the queue.
func (q *Queue[T]) Clear() {
	q.items = []T{}
}

// CopyTo copies the items in the queue to the given slice, starting at the
// given index. If the index is out of range, an error is returned. If the
// slice is not large enough to hold all the items, an error is returned.
func (q *Queue[T]) CopyTo(items []T, index int) error {
	if index < 0 || index > len(items) {
		return ErrIndexOutOfRange
	}

	if len(items)-index < len(q.items) {
		return ErrIndexOutOfRange
	}

	copy(items[index:], q.items)

	return nil
}

// ConcurrentQueue implements a FIFO data structure. It is thread-safe.
type ConcurrentQueue[T any] struct {
	items []T
	mutex sync.RWMutex
}

// NewConcurrentQueue returns a new queue with the given initial items.
func NewConcurrentQueue[T any](values ...T) *ConcurrentQueue[T] {
	return &ConcurrentQueue[T]{items: values}
}

// Enqueue adds an item to the end of the queue.
func (q *ConcurrentQueue[T]) Enqueue(item T) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	q.items = append(q.items, item)
}

// Dequeue removes and returns the item at the front of the queue. If the queue
// is empty, an error is returned.
func (q *ConcurrentQueue[T]) Dequeue() (T, error) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	var zero T
	if len(q.items) == 0 {
		return zero, ErrEmptyQueue
	}

	item := q.items[0]
	q.items = q.items[1:]
	return item, nil
}

// Peek returns the item at the front of the queue without removing it. If the
// queue is empty, an error is returned.
func (q *ConcurrentQueue[T]) Peek() (T, error) {
	q.mutex.RLock()
	defer q.mutex.RUnlock()

	var zero T
	if len(q.items) == 0 {
		return zero, ErrEmptyQueue
	}
	return q.items[0], nil
}

// IsEmpty returns true if the queue is empty.
func (q *ConcurrentQueue[T]) IsEmpty() bool {
	q.mutex.RLock()
	defer q.mutex.RUnlock()

	return len(q.items) == 0
}

// Size returns the number of items in the queue.
func (q *ConcurrentQueue[T]) Size() int {
	q.mutex.RLock()
	defer q.mutex.RUnlock()

	return len(q.items)
}

// String returns a string representation of the queue.
func (q *ConcurrentQueue[T]) String() string {
	q.mutex.RLock()
	defer q.mutex.RUnlock()

	return fmt.Sprintf("%v", q.items)
}

// Clear removes all items from the queue.
func (q *ConcurrentQueue[T]) Clear() {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	q.items = []T{}
}

// CopyTo copies the items in the queue to the given slice, starting at the
// given index. If the index is out of range, an error is returned. If the
// slice is not large enough to hold all the items, an error is returned.
func (q *ConcurrentQueue[T]) CopyTo(items []T, index int) error {
	q.mutex.RLock()
	defer q.mutex.RUnlock()

	if index < 0 || index > len(items) {
		return ErrIndexOutOfRange
	}

	if len(items)-index < len(q.items) {
		return ErrIndexOutOfRange
	}

	copy(items[index:], q.items)

	return nil
}
