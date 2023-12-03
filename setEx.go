package persistent

import (
	"encoding/json"
)

// SetEx defines a set interface implemented using a persistent AVL tree. It works for custom types that implement
// the Ordered[T] interface. See Set[T] for types that support the < interface.
//
// Persistent AVL trees are immutable. Each mutating operation will return the root of a new tree with the requested
// update applied. The implementation uses structural sharing to make immutability efficient. For any given
// update, at most O(log(n)) nodes will be replaced in the new tree. The implementation is also concurrency-safe
// and non-blocking. A *Set[T] instance may be accessed from multiple go-routines without synchronization. See the
// docs for Iterator[T] for notes on concurrent use of iterators.
type SetEx[T Ordered[T]] struct {
	tree *TreeEx[T, bool]
}

// GetKthElement returns the k'th smallest element in a set.
// If no such element exists, ok will be false.
func (s *SetEx[T]) GetKthElement(k int) (e T, ok bool) {
	p, ok := s.tree.GetKthElement(k)
	return p.Key, ok
}

// Contains return true if the set contains the given element.
func (s *SetEx[T]) Contains(elem T) bool {
	var currentRoot *TreeEx[T, bool]
	if s != nil {
		currentRoot = s.tree
	}
	return currentRoot.Contains(elem)
}

// Remove returns the root of a new tree with the given element removed.
func (s *SetEx[T]) Remove(elem T) *SetEx[T] {
	var currentRoot *TreeEx[T, bool]
	if s != nil {
		currentRoot = s.tree
	}

	newRoot := currentRoot.Delete(elem)
	if newRoot == nil {
		return nil
	}

	if newRoot != currentRoot {
		return &SetEx[T]{tree: newRoot}
	}
	return s
}

// Add returns the root of a new tree with the given element added.
func (s *SetEx[T]) Add(elem T) *SetEx[T] {
	if s.Contains(elem) {
		return s
	}
	var currentRoot *TreeEx[T, bool]
	if s != nil {
		currentRoot = s.tree
	}

	return &SetEx[T]{
		tree: currentRoot.Update(elem, true),
	}
}

// LeastUpperBound returns the smallest element e in s, such that
// value <= e. If no such element exists, ok will be false and e will be set to the zero value for T.
func (s *SetEx[T]) LeastUpperBound(value T) (e T, ok bool) {
	if s == nil {
		var ret T
		return ret, false
	}

	kv, found := s.tree.LeastUpperBound(value)
	if !found {
		var ret T
		return ret, false
	}
	return kv.Key, true
}

// GreatestLowerBound returns the largest element e in s, such that
// e <= value. If no such element exists, ok will be false and e will be set to the zero value for T.
func (s *SetEx[T]) GreatestLowerBound(value T) (T, bool) {
	if s == nil {
		var ret T
		return ret, false
	}

	kv, found := s.tree.GreatestLowerBound(value)
	if !found {
		var ret T
		return ret, false
	}
	return kv.Key, true
}

// Size returns the number of elements in the set.
func (s *SetEx[T]) Size() int {
	if s == nil {
		return 0
	}
	return s.tree.Size()
}

// Iter returns an in-order traversal iterator over the elements in the set.
func (s *SetEx[T]) Iter() Iterator[T] {
	var currentRoot *TreeEx[T, bool]
	if s != nil {
		currentRoot = s.tree
	}
	return &SetExIterator[T]{wrapped: currentRoot.Iter()}
}

// IterGte returns an in-order traversal iterator over all elements x in s that are >= e
func (s *SetEx[T]) IterGte(e T) Iterator[T] {
	var currentRoot *TreeEx[T, bool]
	if s != nil {
		currentRoot = s.tree
	}
	return &SetExIterator[T]{wrapped: currentRoot.IterGte(e)}
}

// IsEmpty returns true iif s is empty.
func (s *SetEx[T]) IsEmpty() bool {
	return s == nil || s.tree.IsEmpty()
}

// MarshalJSON marshals the set s as a json array.a
func (s *SetEx[T]) MarshalJSON() ([]byte, error) {
	var arr []T
	iter := s.Iter()
	for iter.Next() {
		arr = append(arr, iter.Current())
	}
	return json.Marshal(arr)
}

// UnmarshalJSON unmarshals a json array into s.
func (s *SetEx[T]) UnmarshalJSON(data []byte) error {
	var arr []T
	err := json.Unmarshal(data, &arr)
	if err != nil {
		return err
	}
	ret := &SetEx[T]{}
	for _, e := range arr {
		ret = ret.Add(e)
	}
	*s = *ret
	return nil
}

type SetExIterator[T Ordered[T]] struct {
	wrapped Iterator[Pair[T, bool]]
}

func (s *SetExIterator[T]) Next() bool {
	return s.wrapped.Next()
}

func (s *SetExIterator[T]) Current() T {
	return s.wrapped.Current().Key
}

func EmptySetEx[T Ordered[T]]() *SetEx[T] {
	return nil
}
