package persistent

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestEnqueue(t *testing.T) {
	q := EmptyQueue[int]().Enqueue(1).Enqueue(2).Enqueue(3)
	require.Equal(t, 3, q.Size())
	require.False(t, q.IsEmpty())

	var v int
	for expected := 1; expected <= 3; expected++ {
		v, q = q.Dequeue()
		require.Equal(t, expected, v)
	}
	require.True(t, q.IsEmpty())
}

func TestDequeueEmpty(t *testing.T) {
	v, q := EmptyQueue[int]().Dequeue()
	require.Equal(t, 0, v)
	require.True(t, q.IsEmpty())
}

func TestQueueSize(t *testing.T) {
	q := EmptyQueue[int]().Enqueue(1).Enqueue(2).Enqueue(3)
	_, q = q.Dequeue()
	q = q.Enqueue(4).Enqueue(5)
	require.Equal(t, 4, q.Size())
}

func TestEmptyQueueSize(t *testing.T) {
	require.Equal(t, 0, EmptyQueue[int]().Size())
}

func TestIsEmpty(t *testing.T) {
	q := EmptyQueue[int]().Enqueue(1).Enqueue(2)
	require.False(t, q.IsEmpty())
	_, q = q.Dequeue()
	require.False(t, q.IsEmpty())
	_, q = q.Dequeue()
	require.True(t, q.IsEmpty())
}

func TestJsonMarshalQueue(t *testing.T) {
	q := EmptyQueue[int]().Enqueue(1).Enqueue(2).Enqueue(3)
	bytes, err := json.Marshal(q)
	require.NoError(t, err)

	var actual []int
	err = json.Unmarshal(bytes, &actual)
	require.NoError(t, err)

	require.Equal(t, []int{1, 2, 3}, actual)
}

func TestJsonUnMarshalQueue(t *testing.T) {
	var items = []int{0, 1, 2}
	bytes, err := json.Marshal(items)
	require.NoError(t, err)

	var q *Queue[int]
	err = json.Unmarshal(bytes, &q)
	require.NoError(t, err)
	require.False(t, q.IsEmpty())
	require.Equal(t, 3, q.Size())

	var x int
	for i := 0; i < 3; i++ {
		x, q = q.Dequeue()
		require.Equal(t, i, x)
	}

}

func TestTopEmpty(t *testing.T) {
	q := EmptyQueue[int]()
	require.Equal(t, 0, q.Top())
	// Verify queue is still empty
	require.True(t, q.IsEmpty())
}

func TestTopWithElementsInDequeueStack(t *testing.T) {
	// Create a queue and dequeue once to ensure elements are in the dequeue stack
	q := EmptyQueue[int]().Enqueue(1).Enqueue(2).Enqueue(3)
	val, q := q.Dequeue()
	require.Equal(t, 1, val)

	// Now the remaining elements should be in the dequeue stack
	// Top should return the next element without modifying the queue
	require.Equal(t, 2, q.Top())
	require.Equal(t, 2, q.Size()) // Size should still be 2

	// Verify queue wasn't modified by calling Top()
	v, q := q.Dequeue()
	require.Equal(t, 2, v)
	require.Equal(t, 1, q.Size())
}

func TestTopWithElementsInEnqueueStackOnly(t *testing.T) {
	// In a freshly created queue, elements are only in enqueue stack
	q := EmptyQueue[int]().Enqueue(10).Enqueue(20).Enqueue(30)

	// Top should return the first element added (10) without dequeuing
	require.Equal(t, 10, q.Top())
	require.Equal(t, 3, q.Size()) // Size should still be 3

	// Verify the queue state by dequeuing all elements
	v1, q := q.Dequeue()
	require.Equal(t, 10, v1)

	v2, q := q.Dequeue()
	require.Equal(t, 20, v2)

	v3, q := q.Dequeue()
	require.Equal(t, 30, v3)

	require.True(t, q.IsEmpty())
}

func TestTopNonDestructive(t *testing.T) {
	q := EmptyQueue[string]().Enqueue("a").Enqueue("b").Enqueue("c")

	// Call Top() multiple times and verify it returns the same value
	require.Equal(t, "a", q.Top())
	require.Equal(t, "a", q.Top())

	// Verify queue size hasn't changed
	require.Equal(t, 3, q.Size())

	// Verify elements can still be dequeued properly in the right order
	v1, q := q.Dequeue()
	require.Equal(t, "a", v1)

	v2, q := q.Dequeue()
	require.Equal(t, "b", v2)

	v3, q := q.Dequeue()
	require.Equal(t, "c", v3)

	require.True(t, q.IsEmpty())
}
