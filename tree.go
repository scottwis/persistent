package persistent

import (
	"encoding/json"
	"golang.org/x/exp/constraints"
)

// Tree implements a persistent AVL tree for keys types that support the < operator. For custom key types see
// TreeEx[K,V].
//
// Note: Both an empty Tree struct and a nil *Tree are valid empty trees.
//
// Persistent AVL trees are immutable. Each mutating operation will return the root of a new tree with the requested
// update applied.The implementation uses structural sharing to make immutability efficient. For any given update
// at most O(log(N)) nodes will be replaced in the new tree. The implementation is concurrency safe and non-blocking.
// A *Tree[K,V] instance may be accessed from multiple go-routines without synchronization. See the docs for Iterator[T]
// for notes on the concurrent use of iterators.
//
// Example:
// var root *Tree[string, int]
// root = root.Update("Hello", 1).Update("World, 2).Update("Apple", 3).Update("Stuff", 4)
// root = root.Remove("Apple")
// iter := root.Iter()
// for iter.MoveNext() {
//     fmt.Printf("%v == %v\n", iter.Current().Key(), iter.Current.Value())
// }
type Tree[K constraints.Ordered, V any] struct {
	left   *Tree[K, V]
	right  *Tree[K, V]
	key    K
	value  V
	size   int
	height int
}

// TreeIterator defines an iterator over a Tree
type TreeIterator[K constraints.Ordered, V any] struct {
	stack   []*Tree[K, V]
	current *Tree[K, V]
}

// Key returns the key associated with the node n. If n is empty, the zero value for K is returned.
func (n *Tree[K, V]) Key() K {
	if n.IsEmpty() {
		var ret K
		return ret
	}
	return n.key
}

// Value returns the value associated with the node n. If n is empty, the zero value for V is returned.
func (n *Tree[K, V]) Value() V {
	if n.IsEmpty() {
		var ret V
		return ret
	}
	return n.value
}

// IsEmpty returns true iif n is empty.
func (n *Tree[K, V]) IsEmpty() bool {
	return n == nil || n.size == 0
}

// Contains returns true if the tree contains the key.
func (n *Tree[K, V]) Contains(key K) bool {
	_, found := n.FindOpt(key)
	return found
}

// Find returns the value associated with key in the tree. Will return a zero value if no such item exists.
func (n *Tree[K, V]) Find(key K) V {
	value, _ := n.FindOpt(key)
	return value
}

// FindOpt returns the value associated with key in the tree. Returns true if found; otherwise false.
// Zero value is returned when found is false.
func (n *Tree[K, V]) FindOpt(key K) (V, bool) {
	if n.IsEmpty() {
		var ret V
		return ret, false
	}

	if n.key < key {
		return n.right.FindOpt(key)
	}

	if key < n.key {
		return n.left.FindOpt(key)
	}

	return n.value, true
}

func newNode[K constraints.Ordered, V any](left *Tree[K, V], right *Tree[K, V], key K, value V) *Tree[K, V] {
	return &Tree[K, V]{
		left:   left,
		right:  right,
		key:    key,
		value:  value,
		size:   left.Size() + right.Size() + 1,
		height: max(left.Height(), right.Height()) + 1,
	}
}

// Update returns the root of a new tree with the value for 'key' set to 'value'.
func (n *Tree[K, V]) Update(key K, value V) *Tree[K, V] {
	if n.IsEmpty() {
		return newNode(nil, nil, key, value)
	}

	if n.key < key {
		return newNode(n.left, n.right.Update(key, value), n.key, n.value).rebalance()
	}

	if key < n.key {
		return newNode(n.left.Update(key, value), n.right, n.key, n.value).rebalance()
	}

	return newNode(n.left, n.right, key, value)
}

// Height returns the height of the tree rooted at node n. Will return 0 if n is empty.
func (n *Tree[K, V]) Height() int {
	if n.IsEmpty() {
		return 0
	}
	return n.height
}

func (n *Tree[K, V]) balanceFactor() int {
	if n.IsEmpty() {
		return 0
	}
	return n.right.Height() - n.left.Height()
}

func (n *Tree[K, V]) rebalance() *Tree[K, V] {
	balance := n.balanceFactor()
	if abs(balance) <= 1 {
		return n
	}

	if balance > 0 {
		rightBalance := n.right.balanceFactor()

		if rightBalance > 0 {
			return n.rotateLeft()
		}

		return n.rotateRightLeft()
	} else {
		leftBalance := n.left.balanceFactor()
		if leftBalance < 0 {
			return n.rotateRight()
		}

		return n.rotateLeftRight()
	}
}

func (n *Tree[K, V]) rotateLeft() *Tree[K, V] {
	return newNode(
		newNode(n.left, n.right.left, n.key, n.value),
		n.right.right,
		n.right.key,
		n.right.value,
	)
}

func (n *Tree[K, V]) rotateRight() *Tree[K, V] {
	return newNode(
		n.left.left,
		newNode(n.left.right, n.right, n.key, n.value),
		n.left.key,
		n.left.value,
	)
}

func (n *Tree[K, V]) rotateRightLeft() *Tree[K, V] {
	return newNode(
		n.left,
		n.right.rotateRight(),
		n.key,
		n.value,
	).rotateLeft()
}

func (n *Tree[K, V]) rotateLeftRight() *Tree[K, V] {
	return newNode(
		n.left.rotateLeft(),
		n.right,
		n.key,
		n.value,
	).rotateRight()
}

// Left returns the left subtree of the current tree. Will never be nil. For empty nodes t.Left() == t.
// If you are using Left() and Right() to traverse a tree, make sure to use IsEmpty() as your termination condition.
func (n *Tree[K, V]) Left() *Tree[K, V] {
	if n.IsEmpty() {
		return n
	}
	return n.left
}

// Right returns the right subtree of the current tree. Will never be nil. For empty nodes t.Left() == t.
// If you are using Left() and Right() to traverse a tree, make sure to use IsEmpty() as your termination condition.
func (n *Tree[K, V]) Right() *Tree[K, V] {
	if n.IsEmpty() {
		return nil
	}
	return n.right
}

// Delete returns the root of a new tree with the entry for 'key' removed.
func (n *Tree[K, V]) Delete(key K) *Tree[K, V] {
	if n.IsEmpty() {
		return nil
	}

	if n.key < key {
		r := n.right.Delete(key)
		if r != n.right {
			return newNode(
				n.left,
				n.right.Delete(key),
				n.key,
				n.value,
			).rebalance()
		}
		return n
	}

	if key < n.key {
		l := n.left.Delete(key)
		if l != n.left {
			return newNode(
				n.left.Delete(key),
				n.right,
				n.key,
				n.value,
			).rebalance()
		}
		return n
	}

	return n.deleteCurrent().rebalance()
}

func (n *Tree[K, V]) deleteCurrent() *Tree[K, V] {
	if n.IsEmpty() {
		return nil
	}
	if n.left.IsEmpty() {
		return n.right
	}

	if n.right.IsEmpty() {
		return n.left
	}

	replacement := n.left.rightMost()

	return newNode(
		n.left.Delete(replacement.key),
		n.right,
		replacement.key,
		replacement.value,
	).rebalance()
}

func (n *Tree[K, V]) rightMost() *Tree[K, V] {
	if n.IsEmpty() {
		return nil
	}
	current := n

	for !current.right.IsEmpty() {
		current = current.right
	}
	return current
}

func (n *Tree[K, V]) Size() int {
	if n.IsEmpty() {
		return 0
	}
	return n.size
}

//LeastUpperBound returns the key-value-pair for the smallest node n such that n.Key() >= key. If there is no such
//node then boolean is false.
func (n *Tree[K, V]) LeastUpperBound(key K) (Pair[K, V], bool) {
	if n.IsEmpty() {
		return n.pair(), false
	}

	if n.key < key {
		return n.right.LeastUpperBound(key)
	}

	if key < n.key {
		ret, found := n.left.LeastUpperBound(key)

		if !found {
			return n.pair(), true
		}
		return ret, true
	}

	return n.pair(), true
}

//GreatestLowerBound returns the key-value-par for the largest node n such that n.Key() <= key. If there is no such
//node then boolean is false.
func (n *Tree[K, V]) GreatestLowerBound(key K) (Pair[K, V], bool) {
	if n.IsEmpty() {
		return n.pair(), false
	}

	if n.key < key {
		ret, found := n.right.GreatestLowerBound(key)
		if !found {
			return n.pair(), true
		}
		return ret, true
	}

	if key < n.key {
		return n.left.GreatestLowerBound(key)
	}

	return n.pair(), true
}

//Iter returns an in-order iterator for the tree.
func (n *Tree[K, V]) Iter() Iterator[Pair[K, V]] {
	ret := TreeIterator[K, V]{
		stack:   nil,
		current: n,
	}

	for !ret.current.IsEmpty() {
		ret.stack = append(ret.stack, ret.current)
		ret.current = ret.current.left
	}
	return &ret
}

//IterGte returns an in-order iterator for the tree for all nodes n such that n.Key() >= glb.
func (n *Tree[K, V]) IterGte(glb K) Iterator[Pair[K, V]] {
	ret := TreeIterator[K, V]{
		stack:   nil,
		current: n,
	}

	for !ret.current.IsEmpty() {
		if ret.current.key < glb {
			ret.current = ret.current.right
		} else if glb < ret.current.key {
			ret.stack = append(ret.stack, ret.current)
			ret.current = ret.current.left
		} else {
			ret.stack = append(ret.stack, ret.current)
			ret.current = nil
		}
	}

	return &ret
}

//Least returns the key-value-pair for the lowest element in the tree. If the tree is empty then boolean
//is false
func (n *Tree[K, V]) Least() (Pair[K, V], bool) {
	if n.IsEmpty() {
		return n.pair(), false
	}

	var ret = n

	for !ret.left.IsEmpty() {
		ret = ret.left
	}

	return ret.pair(), true
}

//Most returns the key-value-pair for the greatest element in the tree. If the tree is empty then boolean is false
func (n *Tree[K, V]) Most() (Pair[K, V], bool) {
	if n.IsEmpty() {
		return n.pair(), false
	}

	var ret = n

	for !ret.right.IsEmpty() {
		ret = ret.right
	}

	return ret.pair(), false
}

func (n *Tree[K, V]) MarshalJSON() ([]byte, error) {
	m := make(map[K]V)
	iter := n.Iter()
	for iter.Next() {
		m[iter.Current().Key] = iter.Current().Value
	}
	return json.Marshal(m)
}

func (n *Tree[K, V]) UnmarshalJSON(data []byte) error {
	var m map[K]V
	err := json.Unmarshal(data, &m)
	if err != nil {
		return err
	}
	tree := &Tree[K, V]{}
	for k, v := range m {
		tree = tree.Update(k, v)
	}
	*n = *tree
	return nil
}

func (i *TreeIterator[K, V]) Next() bool {
	if !i.current.IsEmpty() {
		i.current = i.current.right
		for !i.current.IsEmpty() {
			i.stack = append(i.stack, i.current)
			i.current = i.current.left
		}
	}

	if len(i.stack) != 0 {
		i.current = i.stack[len(i.stack)-1]
		i.stack = i.stack[:len(i.stack)-1]
		return true
	}

	return false
}

func (i *TreeIterator[K, V]) Current() Pair[K, V] {
	if i.current.IsEmpty() {
		panic("invalid iterator position")
	}
	return i.current.pair()
}

func EmptyTree[K constraints.Ordered, V any]() *Tree[K, V] {
	return nil
}

func (n *Tree[K, V]) pair() Pair[K, V] {
	return Pair[K, V]{Key: n.Key(), Value: n.Value()}
}
