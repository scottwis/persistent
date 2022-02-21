package persistent

import "encoding/json"

// Pair defines a key-value pair.
type Pair[K any, V any] interface {
	//Key returns the key stored in the pair. If the pair is empty, will return the 0 value for K.
	Key() K
	//Value returns the value stored in the pair. If the pair is empty, will return the 0 value of V.
	Value() V
	//IsEmpty returns true if the pair is empty.
	IsEmpty() bool
}

// Tree implements a persistent AVL Tree.
//
// Persistent AVL trees are immutable. Each mutating operation will return the root of a new tree with the requested
// update applied. The implementation uses structural sharing to make immutability efficient. For any given
// update, at most O(log(n)) nodes will be replaced in the new tree. The implementation is also concurrency-safe
// and non-blocking. A *Tree[K,V] instance may be accessed from multiple go-routines without synchronization. See the
// docs for Iterator[K,V] for notes on current use of iterators.
//
// To create a tree use the EmptyTree[K,V]() methods.
type Tree[K Ordered[K], V any] struct {
	left   *Tree[K, V]
	right  *Tree[K, V]
	key    K
	value  V
	size   int
	height int
}

type treeIterator[K Ordered[K], V any] struct {
	stack   []*Tree[K, V]
	current *Tree[K, V]
}

func EmptyTree[K Ordered[K], V any]() *Tree[K, V] {
	return nil
}

// Key returns the key stored in the root node of the tree. If the tree is empty, the zero value of K is returned.
func (n *Tree[K, V]) Key() K {
	if n.IsEmpty() {
		var ret K
		return ret
	}
	return n.key
}

// Value returns the value stored in the root node of the tree. If the tree is empty, the zero value of V is returned.
func (n *Tree[K, V]) Value() V {
	if n.IsEmpty() {
		var ret V
		return ret
	}
	return n.value
}

// IsEmpty returns true if the tree is currently empty, false otherwise
func (n *Tree[K, V]) IsEmpty() bool {
	return n == nil || n.size == 0
}

// Find returns the pair associated with key in the tree. Will return an empty pair if no such item exists.
func (n *Tree[K, V]) Find(key K) Pair[K, V] {
	if n.IsEmpty() {
		return n
	}

	if n.key.Equal(key) {
		return n
	}

	if n.key.Less(key) {
		return n.right.Find(key)
	}

	return n.left.Find(key)
}

// Update returns the root of a new tree with the value for 'key' set to 'value'.
func (n *Tree[K, V]) Update(key K, value V) *Tree[K, V] {
	if n.IsEmpty() {
		return newNode(n, n, key, value)
	}

	if n.key.Equal(key) {
		return newNode(n.left, n.right, key, value)
	}

	if n.key.Less(key) {
		return newNode(n.left, n.right.Update(key, value), n.key, n.value).rebalance()
	}

	return newNode(n.left.Update(key, value), n.right, n.key, n.value).rebalance()
}

func newNode[K Ordered[K], V any](left *Tree[K, V], right *Tree[K, V], key K, value V) *Tree[K, V] {
	return &Tree[K, V]{
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

func (n *Tree[K, V]) balanceFactor() int {
	if n.IsEmpty() {
		return 0
	}
	return n.right.Height() - n.left.Height()
}

func (n *Tree[K, V]) Height() int {
	if n.IsEmpty() {
		return 0
	}
	return n.height
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
		return n
	}
	return n.right
}

// Delete returns the root of a new tree with the entry for 'key' removed.
func (n *Tree[K, V]) Delete(key K) *Tree[K, V] {
	if n.IsEmpty() {
		return n
	}

	if n.key.Equal(key) {
		return n.deleteCurrent().rebalance()
	}

	if n.key.Less(key) {
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

func (n *Tree[K, V]) deleteCurrent() *Tree[K, V] {
	if n.IsEmpty() {
		return n
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
		return n
	}
	current := n

	for !current.right.IsEmpty() {
		current = current.right
	}
	return current
}

//Size returns the number of nodes in the tree. Will be 0 for empty trees.
func (n *Tree[K, V]) Size() int {
	if n.IsEmpty() {
		return 0
	}
	return n.size
}

//LeastUpperBound returns the key-value-pair for the smallest node n such that n.Key() >= key. If there is no such
//node then an empty pair is returned.
func (n *Tree[K, V]) LeastUpperBound(key K) Pair[K, V] {
	if n.IsEmpty() {
		return n
	}

	if n.key.Equal(key) {
		return n
	}
	if n.key.Less(key) {
		return n.right.LeastUpperBound(key)
	}

	ret := n.left.LeastUpperBound(key)

	if ret.IsEmpty() {
		return n
	}
	return ret
}

//GreatestLowerBound returns the key-value-par for the largest node n such that n.Key() <= key. If there is no such
//node than an empty pair is returned.
func (n *Tree[K, V]) GreatestLowerBound(key K) Pair[K, V] {
	if n.IsEmpty() {
		return n
	}

	if n.key.Equal(key) {
		return n
	}

	if n.key.Less(key) {
		ret := n.right.GreatestLowerBound(key)

		if ret.IsEmpty() {
			return n
		}
		return ret
	}

	return n.left.GreatestLowerBound(key)
}

//Iter returns an in-order iterator for the tree.
func (n *Tree[K, V]) Iter() Iterator[Pair[K, V]] {
	ret := treeIterator[K, V]{
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
	ret := treeIterator[K, V]{
		stack:   nil,
		current: n,
	}

	for !ret.current.IsEmpty() {
		if ret.current.Key().Equal(glb) {
			ret.stack = append(ret.stack, ret.current)
			ret.current = nil
		} else if ret.current.Key().Less(glb) {
			ret.current = ret.current.right
		} else {
			ret.stack = append(ret.stack, ret.current)
			ret.current = ret.current.left
		}
	}

	return &ret
}

//Least returns the key-value-pair for the lowest element in the tree. If the tree is empty, the returned pair
//is also empty.
func (n *Tree[K, V]) Least() Pair[K, V] {
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
func (n *Tree[K, V]) Most() Pair[K, V] {
	if n.IsEmpty() {
		return n
	}

	var ret = n

	for !ret.right.IsEmpty() {
		ret = ret.right
	}

	return ret
}

func (n *Tree[K, V]) MarshalJSON() ([]byte, error) {
	m := make(map[K]V)
	iter := n.Iter()
	for iter.Next() {
		m[iter.Current().Key()] = iter.Current().Value()
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

func (i *treeIterator[K, V]) Next() bool {
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

func (i *treeIterator[K, V]) Current() Pair[K, V] {
	if i.current.IsEmpty() {
		panic("invalid iterator position")
	}
	return i.current
}
