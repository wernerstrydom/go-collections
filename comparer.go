package collections

// Comparer is a function that returns -1, 0, or 1 for a given input. It is
// used to compare items in a collection. It is used to sort collections.
type Comparer[T any] func(T, T) int

// EqualityComparer is a function that returns true or false when the two
// given inputs are equal. It is used to compare items in a collection.
type EqualityComparer[T any] func(T, T) bool

// DefaultEqualityComparer is the default equality comparer for a given type.
func DefaultEqualityComparer[T comparable](a, b T) bool {
	return a == b
}
