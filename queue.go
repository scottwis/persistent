package persistent

import (
	"encoding/json"
)

// Queue implements a persistent queue.
//
// Persistent queues are immutable. Each mutating operation will return a pointer to a new queue with the
// update applied.The implementation uses structural sharing to make immutability efficient. The implementation is
// concurrency safe and non-blocking. A *Queue[T] instance may be accessed from multiple go-routines without
// synchronization. Each mutating Queue[T] operation has amortized O(1) cost.
type Queue[T any] struct {
	e       *Stack[T]
	d       *Stack[T]
	eBottom T // Tracks the first item added to e
}

// IsEmpty returns true iif the queue is empty.
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

// Enqueue returns a new queue with 'value' added to the end. This is O(1).
func (q *Queue[T]) Enqueue(value T) *Queue[T] {
	ret := Queue[T]{
		e: q.enqueueStack().Push(value),
		d: q.dequeueStack(),
	}

	if q.enqueueStack().IsEmpty() {
		ret.eBottom = value
	} else {
		ret.eBottom = q.eBottom
	}

	return &ret
}

// Dequeue removes the top item from the queue, and returns the dequeued value along with a new queue with the value
// removed. If the queue is empty, the returned 'value' will be the 0 value of T.
//
// Dequeue is amortized O(1). This means an individual Dequeue operation may be O(n), but with a series of n calls
// the average runtime is O(1).
func (q *Queue[T]) Dequeue() (value T, queue *Queue[T]) {
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
	retQueue := &Queue[T]{
		d: d,
		e: q.e,
	}

	if !retQueue.e.IsEmpty() {
		retQueue.eBottom = q.eBottom
	}
	return v, retQueue
}

// Top returns the top element in the queue without removing it. This is O(1).
// If the queue is empty, the returned value will be the zero value of T.
func (q *Queue[T]) Top() T {
	if q.IsEmpty() {
		var zv T
		return zv
	}

	// If dequeue stack has elements, the top is there
	if !q.d.IsEmpty() {
		return q.d.Peek()
	}

	// Otherwise, the top is the bottom of enqueue stack
	return q.eBottom
}

// Size returns the number of elements in 'q'.
func (q *Queue[T]) Size() int {
	if q == nil {
		return 0
	}
	return q.d.Size() + q.e.Size()
}

func (q *Queue[T]) MarshalJSON() ([]byte, error) {
	cur := q
	var buf []T
	var item T
	for !cur.IsEmpty() {
		item, cur = cur.Dequeue()
		buf = append(buf, item)
	}
	return json.Marshal(buf)
}

func (q *Queue[T]) UnmarshalJSON(bytes []byte) error {
	var tmp *Queue[T]
	var buf []T
	err := json.Unmarshal(bytes, &buf)
	if err != nil {
		return err
	}
	for _, b := range buf {
		tmp = tmp.Enqueue(b)
	}
	*q = *tmp
	return nil
}

// EmptyQueue returns a new empty queue.
func EmptyQueue[T any]() *Queue[T] {
	return nil
}
