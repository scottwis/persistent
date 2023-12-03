package persistent

import (
	"encoding/json"
	"golang.org/x/exp/constraints"
)

// Set defines a set interface implemented using a persistent AVL tree. It works for all types that support the <
// operator. For custom types see the SetEx method.
//
// Persistent AVL trees are immutable. Each mutating operation will return the root of a new tree with the requested
// update applied. The implementation uses structural sharing to make immutability efficient. For any given
// update, at most O(log(n)) nodes will be replaced in the new tree. The implementation is also concurrency-safe
// and non-blocking. A *Set[T] instance may be accessed from multiple go-routines without synchronization. See the
// docs on Iterator[T] for notes on concurrent use of iterators.
type Set[T constraints.Ordered] struct {
	tree *Tree[T, bool]
}

// GetKthElement returns the k'th smallest element in a set.
// If no such element exists, ok will be false.
func (s *Set[T]) GetKthElement(k int) (e T, ok bool) {
	p, ok := s.tree.GetKthElement(k)
	return p.Key, ok
}

// Contains return true if the set contains the given element.
func (s *Set[T]) Contains(elem T) bool {
	var currentRoot *Tree[T, bool]
	if s != nil {
		currentRoot = s.tree
	}
	return currentRoot.Contains(elem)
}

// Remove returns the root of a new tree with the given element removed.
func (s *Set[T]) Remove(elem T) *Set[T] {
	var currentRoot *Tree[T, bool]
	if s != nil {
		currentRoot = s.tree
	}

	newRoot := currentRoot.Delete(elem)
	if newRoot == nil {
		return nil
	}

	if newRoot != currentRoot {
		return &Set[T]{tree: newRoot}
	}
	return s
}

// Add returns the root of a new tree with the given element added.
func (s *Set[T]) Add(elem T) *Set[T] {
	if s.Contains(elem) {
		return s
	}
	var currentRoot *Tree[T, bool]
	if s != nil {
		currentRoot = s.tree
	}

	return &Set[T]{
		tree: currentRoot.Update(elem, true),
	}
}

// LeastUpperBound returns the smallest element e in s, such that
// value <= e. If no such element exists, ok will be false and e will be set to the zero value for T.
func (s *Set[T]) LeastUpperBound(value T) (e T, ok bool) {
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
func (s *Set[T]) GreatestLowerBound(value T) (T, bool) {
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
func (s *Set[T]) Size() int {
	if s == nil {
		return 0
	}
	return s.tree.Size()
}

// Iter returns an in-order traversal iterator over the elements in the set.
func (s *Set[T]) Iter() Iterator[T] {
	var currentRoot *Tree[T, bool]
	if s != nil {
		currentRoot = s.tree
	}
	return &SetIterator[T]{wrapped: currentRoot.Iter()}
}

// IterGte returns an in-order traversal iterator over all elements x in s that are >= e
func (s *Set[T]) IterGte(e T) Iterator[T] {
	var currentRoot *Tree[T, bool]
	if s != nil {
		currentRoot = s.tree
	}
	return &SetIterator[T]{wrapped: currentRoot.IterGte(e)}
}

// IsEmpty returns true iif s is empty.
func (s *Set[T]) IsEmpty() bool {
	return s == nil || s.tree.IsEmpty()
}

// MarshalJSON marshals the set s as a json array.a
func (s *Set[T]) MarshalJSON() ([]byte, error) {
	var arr []T
	iter := s.Iter()
	for iter.Next() {
		arr = append(arr, iter.Current())
	}
	return json.Marshal(arr)
}

// UnmarshalJSON unmarshals a json array into s.
func (s *Set[T]) UnmarshalJSON(data []byte) error {
	var arr []T
	err := json.Unmarshal(data, &arr)
	if err != nil {
		return err
	}
	ret := &Set[T]{}
	for _, e := range arr {
		ret = ret.Add(e)
	}
	*s = *ret
	return nil
}

type SetIterator[T constraints.Ordered] struct {
	wrapped Iterator[Pair[T, bool]]
}

func (s *SetIterator[T]) Next() bool {
	return s.wrapped.Next()
}

func (s *SetIterator[T]) Current() T {
	return s.wrapped.Current().Key
}

func EmptySet[T constraints.Ordered]() *Set[T] {
	return nil
}
