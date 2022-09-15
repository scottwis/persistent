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