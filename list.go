package collections

import (
	"fmt"
	"sync"
)

// List implements a list data structure. It is not thread-safe.
type List[T any] struct {
	items    []T
	comparer EqualityComparer[T]
}

// NewListWithEqualityComparer returns a new list with the given initial items.
func NewListWithEqualityComparer[T any](comparer EqualityComparer[T], values ...T) *List[T] {
	return &List[T]{items: values, comparer: comparer}
}

// NewList returns a new list with the given initial items.
func NewList[T comparable](values ...T) *List[T] {
	return &List[T]{items: values, comparer: DefaultEqualityComparer[T]}
}

// Add adds an item to the end of the list.
func (l *List[T]) Add(item T) {
	l.items = append(l.items, item)
}

// AddRange adds the given items to the end of the list.
func (l *List[T]) AddRange(items ...T) {
	l.items = append(l.items, items...)
}

// Insert inserts an item at the given index.
func (l *List[T]) Insert(index int, item T) error {
	if index < 0 || index > len(l.items) {
		return ErrIndexOutOfRange
	}

	l.items = append(l.items[:index], append([]T{item}, l.items[index:]...)...)
	return nil
}

// Remove removes the given item from the list. If the item is not found, false
// is returned, otherwise true is returned.
func (l *List[T]) Remove(item T) bool {
	index := l.IndexOf(item)
	if index == -1 {
		return false
	}

	err := l.RemoveAt(index)
	if err != nil {
		return false
	}

	return true
}

// RemoveAt removes the item at the given index.
func (l *List[T]) RemoveAt(index int) error {
	if index < 0 || index >= len(l.items) {
		return ErrIndexOutOfRange
	}

	l.items = append(l.items[:index], l.items[index+1:]...)
	return nil
}

// Clear removes all items from the list.
func (l *List[T]) Clear() {
	l.items = []T{}
}

// Contains returns true if the list contains the given item.
func (l *List[T]) Contains(item T) bool {
	return l.IndexOf(item) != -1
}

// IndexOf returns the index of the given item. If the item is not found, -1 is
// returned.
func (l *List[T]) IndexOf(item T) int {
	for i, v := range l.items {
		if l.comparer(v, item) {
			return i
		}
	}
	return -1
}

// LastIndexOf returns the last index of the given item. If the item is not
// found, -1 is returned.
func (l *List[T]) LastIndexOf(item T) int {
	for i := len(l.items) - 1; i >= 0; i-- {
		if l.comparer(l.items[i], item) {
			return i
		}
	}
	return -1
}

// CopyTo copies the items in the list to the given slice, starting at the given
// index.
func (l *List[T]) CopyTo(array []T, arrayIndex int) error {
	if arrayIndex < 0 || arrayIndex >= len(array) {
		return ErrIndexOutOfRange
	}

	if arrayIndex+len(l.items) > len(array) {
		return ErrIndexOutOfRange
	}

	copy(array[arrayIndex:], l.items)
	return nil
}

// Get returns the item at the given index.
func (l *List[T]) Get(index int) (T, error) {

	var zero T
	if index < 0 || index >= len(l.items) {
		return zero, ErrIndexOutOfRange
	}

	return l.items[index], nil
}

// Set sets the item at the given index.
func (l *List[T]) Set(index int, item T) error {
	if index < 0 || index >= len(l.items) {
		return ErrIndexOutOfRange
	}

	l.items[index] = item
	return nil
}

func (l *List[T]) String() string {
	return fmt.Sprintf("%v", l.items)
}

type ConcurrentList[T any] struct {
	items    []T
	comparer EqualityComparer[T]
	mutex    sync.RWMutex
}

// NewConcurrentListWithComparer returns a new list with the given initial items.
func NewConcurrentListWithComparer[T any](comparer EqualityComparer[T], values ...T) *ConcurrentList[T] {
	return &ConcurrentList[T]{items: values, comparer: comparer}
}

// NewConcurrentList returns a new list with the given initial items.
func NewConcurrentList[T comparable](values ...T) *ConcurrentList[T] {
	return &ConcurrentList[T]{items: values, comparer: DefaultEqualityComparer[T]}
}

// Add adds an item to the end of the list.
func (l *ConcurrentList[T]) Add(item T) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.items = append(l.items, item)
}

// AddRange adds the given items to the end of the list.
func (l *ConcurrentList[T]) AddRange(items ...T) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.items = append(l.items, items...)
}

// Insert inserts an item at the given index.
func (l *ConcurrentList[T]) Insert(index int, item T) error {
	if index < 0 || index > len(l.items) {
		return ErrIndexOutOfRange
	}

	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.items = append(l.items[:index], append([]T{item}, l.items[index:]...)...)
	return nil
}

// Remove removes the given item from the list. If the item is not found, false
// is returned, otherwise true is returned.
func (l *ConcurrentList[T]) Remove(item T) bool {
	index := l.IndexOf(item)
	if index == -1 {
		return false
	}

	err := l.RemoveAt(index)
	if err != nil {
		return false
	}

	return true
}

// RemoveAt removes the item at the given index.
func (l *ConcurrentList[T]) RemoveAt(index int) error {
	if index < 0 || index >= len(l.items) {
		return ErrIndexOutOfRange
	}

	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.items = append(l.items[:index], l.items[index+1:]...)
	return nil
}

// Clear removes all items from the list.
func (l *ConcurrentList[T]) Clear() {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.items = []T{}
}

// Contains returns true if the list contains the given item.
func (l *ConcurrentList[T]) Contains(item T) bool {
	return l.IndexOf(item) != -1
}

// IndexOf returns the index of the given item. If the item is not found, -1 is
// returned.
func (l *ConcurrentList[T]) IndexOf(item T) int {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	for i, v := range l.items {
		if l.comparer(v, item) {
			return i
		}
	}
	return -1
}

// LastIndexOf returns the last index of the given item. If the item is not
// found, -1 is returned.
func (l *ConcurrentList[T]) LastIndexOf(item T) int {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	for i := len(l.items) - 1; i >= 0; i-- {
		if l.comparer(l.items[i], item) {
			return i
		}
	}
	return -1
}

// CopyTo copies the items in the list to the given slice, starting at the given
// index.
func (l *ConcurrentList[T]) CopyTo(array []T, arrayIndex int) error {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	if arrayIndex < 0 || arrayIndex >= len(array) {
		return ErrIndexOutOfRange
	}

	if arrayIndex+len(l.items) > len(array) {
		return ErrIndexOutOfRange
	}

	copy(array[arrayIndex:], l.items)
	return nil
}

// Get returns the item at the given index.
func (l *ConcurrentList[T]) Get(index int) (T, error) {

	var zero T
	if index < 0 || index >= len(l.items) {
		return zero, ErrIndexOutOfRange
	}

	l.mutex.RLock()
	defer l.mutex.RUnlock()
	return l.items[index], nil
}

// Set sets the item at the given index.
func (l *ConcurrentList[T]) Set(index int, item T) error {
	if index < 0 || index >= len(l.items) {
		return ErrIndexOutOfRange
	}

	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.items[index] = item
	return nil
}

func (l *ConcurrentList[T]) String() string {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	return fmt.Sprintf("%v", l.items)
}
