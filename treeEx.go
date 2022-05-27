package persistent

import (
	"encoding/json"
)

// TreeEx implements a persistent AVL tree for keys implementing the Ordered[K] interface. For built-in
// ordered keys (types supporting <) see Tree[K,V],
//
// Persistent AVL trees are immutable. Each mutating operation will return the root of a new tree with the requested
// update applied.The implementation uses structural sharing to make immutability efficient. For any given update
// at most O(log(N)) nodes will be replaced in the new tree. The implementation is concurrency safe and non-blocking.
// A *Tree[K,V] instance may be accessed from multiple go-routines without synchronization. See the docs for Iterator[T]
// for notes on the concurrent use of iterators.
type TreeEx[K Ordered[K], V any] struct {
	left   *TreeEx[K, V]
	right  *TreeEx[K, V]
	key    K
	value  V
	size   int
	height int
}

// TreeExIterator defines an iterator over a TreeEx
type TreeExIterator[K Ordered[K], V any] struct {
	stack   []*TreeEx[K, V]
	current *TreeEx[K, V]
}

// Key returns the key associated with the node n. If n is empty, the zero value for K is returned.
func (n *TreeEx[K, V]) Key() K {
	if n.IsEmpty() {
		var ret K
		return ret
	}
	return n.key
}

// Value returns the value associated with the node n. If n is empty, the zero value for V is returned.
func (n *TreeEx[K, V]) Value() V {
	if n.IsEmpty() {
		var ret V
		return ret
	}
	return n.value
}

// IsEmpty returns true iif n is empty.
func (n *TreeEx[K, V]) IsEmpty() bool {
	return n == nil || n.size == 0
}

// Find returns the pair associated with key in the tree. Will return an empty pair if no such item exists.
func (n *TreeEx[K, V]) Find(key K) Pair[K, V] {
	if n.IsEmpty() {
		return n
	}

	if n.key.Less(key) {
		return n.right.Find(key)
	}

	if key.Less(n.key) {
		return n.left.Find(key)
	}

	return n
}

func newExNode[K Ordered[K], V any](left *TreeEx[K, V], right *TreeEx[K, V], key K, value V) *TreeEx[K, V] {
	return &TreeEx[K, V]{
		left:   left,
		right:  right,
		key:    key,
		value:  value,
		size:   left.Size() + right.Size() + 1,
		height: max(left.Height(), right.Height()) + 1,
	}
}

func max(i, j int) int {
	if i > j {
		return i
	}
	return j
}

func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

// Update returns the root of a new tree with the value for 'key' set to 'value'.
func (n *TreeEx[K, V]) Update(key K, value V) *TreeEx[K, V] {
	if n.IsEmpty() {
		return newExNode(nil, nil, key, value)
	}

	if n.key.Less(key) {
		return newExNode(n.left, n.right.Update(key, value), n.key, n.value).rebalance()
	}

	if key.Less(n.key) {
		return newExNode(n.left.Update(key, value), n.right, n.key, n.value).rebalance()
	}

	return newExNode(n.left, n.right, key, value)
}

// Height returns the height of the tree rooted at node n. Will return 0 if n is empty.
func (n *TreeEx[K, V]) Height() int {
	if n.IsEmpty() {
		return 0
	}
	return n.height
}

func (n *TreeEx[K, V]) balanceFactor() int {
	if n.IsEmpty() {
		return 0
	}
	return n.right.Height() - n.left.Height()
}

func (n *TreeEx[K, V]) rebalance() *TreeEx[K, V] {
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

func (n *TreeEx[K, V]) rotateLeft() *TreeEx[K, V] {
	return newExNode(
		newExNode(n.left, n.right.left, n.key, n.value),
		n.right.right,
		n.right.key,
		n.right.value,
	)
}

func (n *TreeEx[K, V]) rotateRight() *TreeEx[K, V] {
	return newExNode(
		n.left.left,
		newExNode(n.left.right, n.right, n.key, n.value),
		n.left.key,
		n.left.value,
	)
}

func (n *TreeEx[K, V]) rotateRightLeft() *TreeEx[K, V] {
	return newExNode(
		n.left,
		n.right.rotateRight(),
		n.key,
		n.value,
	).rotateLeft()
}

func (n *TreeEx[K, V]) rotateLeftRight() *TreeEx[K, V] {
	return newExNode(
		n.left.rotateLeft(),
		n.right,
		n.key,
		n.value,
	).rotateRight()
}

// Left returns the left subtree of the current tree. Will never be nil. For empty nodes t.Left() == t.
// If you are using Left() and Right() to traverse a tree, make sure to use IsEmpty() as your termination condition.
func (n *TreeEx[K, V]) Left() *TreeEx[K, V] {
	if n.IsEmpty() {
		return nil
	}
	return n.left
}

// Right returns the right subtree of the current tree. Will never be nil. For empty nodes t.Left() == t.
// If you are using Left() and Right() to traverse a tree, make sure to use IsEmpty() as your termination condition.
func (n *TreeEx[K, V]) Right() *TreeEx[K, V] {
	if n.IsEmpty() {
		return nil
	}
	return n.right
}

// Delete returns the root of a new tree with the entry for 'key' removed.
func (n *TreeEx[K, V]) Delete(key K) *TreeEx[K, V] {
	if n.IsEmpty() {
		return nil
	}

	if n.key.Less(key) {
		r := n.right.Delete(key)
		if r != n.right {
			return newExNode(
				n.left,
				n.right.Delete(key),
				n.key,
				n.value,
			).rebalance()
		}
		return n
	}

	if key.Less(n.key) {
		l := n.left.Delete(key)
		if l != n.left {
			return newExNode(
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

func (n *TreeEx[K, V]) deleteCurrent() *TreeEx[K, V] {
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

	return newExNode(
		n.left.Delete(replacement.key),
		n.right,
		replacement.key,
		replacement.value,
	).rebalance()
}

func (n *TreeEx[K, V]) rightMost() *TreeEx[K, V] {
	if n.IsEmpty() {
		return nil
	}
	current := n

	for !current.right.IsEmpty() {
		current = current.right
	}
	return current
}

func (n *TreeEx[K, V]) Size() int {
	if n.IsEmpty() {
		return 0
	}
	return n.size
}

//LeastUpperBound returns the key-value-pair for the smallest node n such that n.Key() >= key. If there is no such
//node then an empty pair is returned.
func (n *TreeEx[K, V]) LeastUpperBound(key K) Pair[K, V] {
	if n.IsEmpty() {
		return n
	}

	if n.key.Less(key) {
		return n.right.LeastUpperBound(key)
	}

	if key.Less(n.key) {
		ret := n.left.LeastUpperBound(key)

		if ret.IsEmpty() {
			return n
		}
		return ret
	}

	return n
}

//GreatestLowerBound returns the key-value-par for the largest node n such that n.Key() <= key. If there is no such
//node than an empty pair is returned.
func (n *TreeEx[K, V]) GreatestLowerBound(key K) Pair[K, V] {
	if n.IsEmpty() {
		return n
	}

	if n.key.Less(key) {
		ret := n.right.GreatestLowerBound(key)

		if ret.IsEmpty() {
			return n
		}
		return ret
	}

	if key.Less(n.key) {
		return n.left.GreatestLowerBound(key)
	}

	return n
}

//Iter returns an in-order iterator for the tree.
func (n *TreeEx[K, V]) Iter() Iterator[Pair[K, V]] {
	ret := TreeExIterator[K, V]{
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
func (n *TreeEx[K, V]) IterGte(glb K) Iterator[Pair[K, V]] {
	ret := TreeExIterator[K, V]{
		stack:   nil,
		current: n,
	}

	for !ret.current.IsEmpty() {
		if ret.current.Key().Less(glb) {
			ret.current = ret.current.right
		} else if glb.Less(ret.current.Key()) {
			ret.stack = append(ret.stack, ret.current)
			ret.current = ret.current.left
		} else {
			ret.stack = append(ret.stack, ret.current)
			ret.current = nil
		}
	}

	return &ret
}

//Least returns the key-value-pair for the lowest element in the tree. If the tree is empty, the returned pair
//is also empty.
func (n *TreeEx[K, V]) Least() Pair[K, V] {
	if n.IsEmpty() {
		return n
	}

	var ret = n

	for !ret.left.IsEmpty() {
		ret = ret.left
	}

	return ret
}

//Most returns the key-value-pair for the greatest element in the tree. If the tree is empty, the returned pair
//is also empty
func (n *TreeEx[K, V]) Most() Pair[K, V] {
	if n.IsEmpty() {
		return n
	}

	var ret = n

	for !ret.right.IsEmpty() {
		ret = ret.right
	}

	return ret
}

func (n *TreeEx[K, V]) MarshalJSON() ([]byte, error) {
	m := make(map[K]V)
	iter := n.Iter()
	for iter.Next() {
		m[iter.Current().Key()] = iter.Current().Value()
	}
	return json.Marshal(m)
}

func (n *TreeEx[K, V]) UnmarshalJSON(data []byte) error {
	var m map[K]V
	err := json.Unmarshal(data, &m)
	if err != nil {
		return err
	}
	tree := &TreeEx[K, V]{}
	for k, v := range m {
		tree = tree.Update(k, v)
	}
	*n = *tree
	return nil
}

func (i *TreeExIterator[K, V]) Next() bool {
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

func (i *TreeExIterator[K, V]) Current() Pair[K, V] {
	if i.current.IsEmpty() {
		panic("invalid iterator position")
	}
	return i.current
}

func EmptyTreeEx[K Ordered[K], V any]() *TreeEx[K, V] {
	return nil
}
