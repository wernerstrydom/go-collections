package collections

// Predicate is a function that returns true or false for a given input. It is
// used to filter collections.
type Predicate[T any] func(T) bool
