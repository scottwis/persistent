package persistent

// Pair defines an interface for a Key / Value pair.
type Pair[K any, V any] interface {
	// Key Returns the key associated with the pair. Will return the zero value for K if IsEmpty is true.
	Key() K

	// Value Returns the value associated with the pair. Will return the zero value for V if IsEmpty is true.
	Value() V

	// IsEmpty returns true iif the pair is empty.
	IsEmpty() bool
}

// Iterator defines an interface for an iterator over a persistent data structure.
//
// It is safe to use multiple iterators over the same data structure concurrently. A single
// Iterator instance, however, is not concurrency safe. Sharing an iterator instance between go-routines requires
// explicit synchronization.
//
// Example:
//    var iter = tree.Iterator()
//    for iter.Next() {
//        fmt.Println(iter.Current())
//    }
type Iterator[T any] interface {
	//Next advances the iterator by one position, and returns true iif the new position is valid.
	Next() bool

	//Current returns the value at the current iterator position. Will panic if the iterator position is invalid.
	Current() T
}

// Ordered defines an interface for a custom key type. Any type K defining a Less method can be used with
// TreeEx[K, V] or OrderedSet[K]. For primitive key types see Tree[K, V] and Set[K]
type Ordered[T any] interface {
	comparable
	Less(rhs T) bool
}
