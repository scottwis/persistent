package persistent

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPush(t *testing.T) {
	s := EmptyStack[int]().Push(2)
	require.Equal(t, 2, s.Peek())
	require.True(t, s.Pop().IsEmpty())
}

func TestPeekEmpty(t *testing.T) {
	require.Equal(t, 0, EmptyStack[int]().Peek())
}

func TestPopEmpty(t *testing.T) {
	require.True(t, EmptyStack[int]().Pop().Pop().IsEmpty())
}

func TestEmptyZerosStack(t *testing.T) {
	var s Stack[int]
	require.True(t, s.IsEmpty())
}

func TestStackSize(t *testing.T) {
	require.Equal(t, 3, EmptyStack[int]().Push(1).Push(2).Push(3).Size())
}

func TestReverse(t *testing.T) {
	reversed := EmptyStack[int]().Push(1).Push(2).Reverse()
	require.Equal(t, 2, reversed.Size())
	require.Equal(t, 1, reversed.Peek())
	require.Equal(t, 2, reversed.Pop().Peek())
	require.True(t, reversed.Pop().Pop().IsEmpty())
}
