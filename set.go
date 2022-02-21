package persistent

import "encoding/json"

//Set defines a set interface implemented using a persistent AVL tree.
//
// Persistent AVL trees are immutable. Each mutating operation will return the root of a new tree with the requested
// update applied. The implementation uses structural sharing to make immutability efficient. For any given
// update, at most O(log(n)) nodes will be replaced in the new tree. The implementation is also concurrency-safe
// and non-blocking. A *Tree[K,V] instance may be accessed from multiple go-routines without synchronization. See the
// docs for Iterator[K,V] for notes on current use of iterators.
//
// To create a set use the EmptySet[K]() methods.
type Set[T Ordered[T]] struct {
	tree *Tree[T, bool]
}

//Contains return true if the set contains the given element.
func (s *Set[T]) Contains(elem T) bool {
	var currentRoot *Tree[T, bool]
	if s != nil {
		currentRoot = s.tree
	}
	return !currentRoot.Find(elem).IsEmpty()
}

//Remove returns the root of a new tree with the given element removed.
func (s *Set[T]) Remove(elem T) *Set[T] {
	var currentRoot *Tree[T, bool]
	if s != nil {
		currentRoot = s.tree
	}

	newRoot := currentRoot.Delete(elem)
	if newRoot != currentRoot {
		return &Set[T]{tree: newRoot}
	}
	return s
}

//Add returns the root of a new tree with the given element added.
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

//LeastUpperBound returns the smallest element e in s, such that
//value <= e. If no such element exists, ok will be false and e will be set to the zero value for T.
func (s *Set[T]) LeastUpperBound(value T) (e T, ok bool) {
	if s == nil {
		var ret T
		return ret, false
	}

	n := s.tree.LeastUpperBound(value)
	return n.Key(), !n.IsEmpty()
}

//GreatestLowerBound returns the largest element e in s, such that
// e <= value. If no such element exists, ok will be false and e will be set to the zero value for T.
func (s *Set[T]) GreatestLowerBound(value T) (T, bool) {
	if s == nil {
		var ret T
		return ret, false
	}

	n := s.tree.GreatestLowerBound(value)
	return n.Key(), !n.IsEmpty()
}

//Size returns the number of elements in the set.
func (s *Set[T]) Size() int {
	if s == nil {
		return 0
	}
	return s.tree.Size()
}

//Iter returns an in-order traversal iterator over the elements in the set.
func (s *Set[T]) Iter() Iterator[T] {
	var currentRoot *Tree[T, bool]
	if s != nil {
		currentRoot = s.tree
	}
	return &setIterator[T]{wrapped: currentRoot.Iter()}
}

//IterGte returns an in-order traversal iterator over all elements x in s that are >= e
func (s *Set[T]) IterGte(e T) Iterator[T] {
	var currentRoot *Tree[T, bool]
	if s != nil {
		currentRoot = s.tree
	}
	return &setIterator[T]{wrapped: currentRoot.IterGte(e)}
}

//IsEmpty returns true iif s is empty.
func (s *Set[T]) IsEmpty() bool {
	return s == nil || s.tree.IsEmpty()
}

//MarshalJSON marshals the set s as a json array.a
func (s *Set[T]) MarshalJSON() ([]byte, error) {
	var arr []T
	iter := s.Iter()
	for iter.Next() {
		arr = append(arr, iter.Current())
	}
	return json.Marshal(arr)
}

//UnmarshalJSON unmarshals a json array into s.
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

//EmptySet returns a new EmptySet of type T.
func EmptySet[T Ordered[T]]() *Set[T] {
	return nil
}

type setIterator[T Ordered[T]] struct {
	wrapped Iterator[Pair[T, bool]]
}

func (s *setIterator[T]) Next() bool {
	return s.wrapped.Next()
}

func (s *setIterator[T]) Current() T {
	return s.wrapped.Current().Key()
}
