package persistent

type Queue[T any] struct {
	e *Stack[T]
	d *Stack[T]
}

func (q *Queue[T]) IsEmpty() bool {
	return q == nil || (q.e.IsEmpty() && q.d.IsEmpty())
}

func (q *Queue[T]) enqueueStack() *Stack[T] {
	if q == nil {
		return nil
	}
	return q.e
}

func (q *Queue[T]) dequeueStack() *Stack[T] {
	if q == nil {
		return nil
	}
	return q.d
}

func (q *Queue[T]) Enqueue(value T) *Queue[T] {
	return &Queue[T]{
		e: q.enqueueStack().Push(value),
		d: q.dequeueStack(),
	}
}

func (q *Queue[T]) Dequeue() (T, *Queue[T]) {
	if q.IsEmpty() {
		var zv T
		return zv, nil
	}
	if q.d.IsEmpty() {
		reversed := q.e.Reverse()
		return reversed.Peek(), &Queue[T]{
			e: nil,
			d: reversed.Pop(),
		}
	}
	v := q.d.Peek()
	d := q.d.Pop()
	if d.IsEmpty() && q.e.IsEmpty() {
		return v, nil
	}
	return v, &Queue[T]{
		d: d,
		e: q.e,
	}
}

func (q *Queue[T]) Size() int {
	if q == nil {
		return 0
	}
	return q.d.Size() + q.e.Size()
}

func EmptyQueue[T any]() *Queue[T] {
	return nil
}
