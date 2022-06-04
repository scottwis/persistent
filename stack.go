package persistent

// Stack implements a persistent stack.
//
// Persistent stacks are immutable. Each mutating operation will return a pointer to a new stack with the
// update applied.The implementation uses structural sharing to make immutability efficient. The implementation is
// concurrency safe and non-blocking. A *Stack[T] instance may be accessed from multiple go-routines without
// synchronization. Each mutating Stack[T] operation is O(1).
type Stack[T any] struct {
	next  *Stack[T]
	value T
	size  int
}

// Peek returns the top-most element of the stack. If the stack is empty, it will return the zero value for T.
func (s *Stack[T]) Peek() T {
	if s.IsEmpty() {
		var ret T
		return ret
	}
	return s.value
}

// Pop returns a new stack with the top element removed. If the stack s is empty, s.Pop() is also empty.
func (s *Stack[T]) Pop() *Stack[T] {
	if s.IsEmpty() {
		return nil
	}
	return s.next
}

// IsEmpty returns true iif s is empty.
func (s *Stack[T]) IsEmpty() bool {
	return s == nil || s.size == 0
}

// Push returns a new stack with 'value' added to the top.
func (s *Stack[T]) Push(value T) *Stack[T] {
	return &Stack[T]{
		next:  s,
		value: value,
		size:  s.Size() + 1,
	}
}

// Size returns the number of elements in s.
func (s *Stack[T]) Size() int {
	if s == nil {
		return 0
	}
	return s.size
}

// Reverse returns the stack s in reverse order. This is mostly used in the implementation of Queue[T].
func (s *Stack[T]) Reverse() *Stack[T] {
	var ret *Stack[T]
	for ; !s.IsEmpty(); s = s.Pop() {
		ret = ret.Push(s.Peek())
	}
	return ret
}

// EmptyStack returns a new empty Stack[T].
func EmptyStack[T any]() *Stack[T] {
	return nil
}
