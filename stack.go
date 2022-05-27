package persistent

type Stack[T any] struct {
	next  *Stack[T]
	value T
	size  int
}

func (s *Stack[T]) Peek() T {
	if s.IsEmpty() {
		var ret T
		return ret
	}
	return s.value
}

func (s *Stack[T]) Pop() *Stack[T] {
	if s.IsEmpty() {
		return nil
	}
	return s.next
}

func (s *Stack[T]) IsEmpty() bool {
	return s == nil || s.size == 0
}

func (s *Stack[T]) Push(value T) *Stack[T] {
	return &Stack[T]{
		next:  s,
		value: value,
		size:  s.Size() + 1,
	}
}

func (s *Stack[T]) Size() int {
	if s == nil {
		return 0
	}
	return s.size
}

func (s *Stack[T]) Reverse() *Stack[T] {
	var ret *Stack[T]
	for ; !s.IsEmpty(); s = s.Pop() {
		ret = ret.Push(s.Peek())
	}
	return ret
}

func EmptyStack[T any]() *Stack[T] {
	var ret *Stack[T]
	return ret
}
