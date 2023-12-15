package collections

import (
	"errors"
)

// ErrIndexOutOfRange is returned when the index is out of range.
var ErrIndexOutOfRange = errors.New("index out of range")

// ErrEmptyStack is returned when the stack is empty.
var ErrEmptyStack = errors.New("stack is empty")

// ErrEmptyQueue is returned when the queue is empty.
var ErrEmptyQueue = errors.New("queue is empty")
