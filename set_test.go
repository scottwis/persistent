package persistent

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestEmptySetContains(t *testing.T) {
	var x Set[Key[int]]
	require.False(t, x.Contains(AsKey(42)))
}

func TestNilSetContains(t *testing.T) {
	var x *Set[Key[int]]
	require.False(t, x.Contains(AsKey(42)))
}

func TestEmptySetRemove(t *testing.T) {
	var x Set[Key[string]]
	pSet := x.Remove(AsKey("2"))
	require.Equal(t, 0, pSet.Size())
}

func TestNilSetRemove(t *testing.T) {
	var x *Set[Key[string]]
	pSet := x.Remove(AsKey("2"))
	require.Equal(t, 0, pSet.Size())
}

func TestEmptySetAdd(t *testing.T) {
	var x Set[Key[string]]
	pSet := x.Add(AsKey("2"))
	require.Equal(t, 1, pSet.Size())
	require.True(t, pSet.Contains(AsKey("2")))
}

func TestNilSetAdd(t *testing.T) {
	var x *Set[Key[string]]
	pSet := x.Add(AsKey("2"))
	require.Equal(t, 1, pSet.Size())
	require.True(t, pSet.Contains(AsKey("2")))
}

func TestSetAddTwice(t *testing.T) {
	x := EmptySet[Key[int]]().Add(AsKey(42))
	y := x.Add(AsKey(42))
	if x != y {
		require.Fail(t, "Adding twice modified pointer")
	}
}

func TestSetAddRemove(t *testing.T) {
	x := EmptySet[Key[int]]().Add(AsKey(4)).Add(AsKey(2)).Remove(AsKey(4))
	require.Equal(t, 1, x.Size())
	require.True(t, x.Contains(AsKey(2)))
	require.False(t, x.Contains(AsKey(1)))
}

func TestEmptySetLub(t *testing.T) {
	var x Set[Key[int]]
	v, ok := x.LeastUpperBound(AsKey(42))
	require.Equal(t, 0, v.Key)
	require.False(t, ok)
}

func TestNilSetLub(t *testing.T) {
	var x *Set[Key[int]]
	v, ok := x.LeastUpperBound(AsKey(42))
	require.Equal(t, 0, v.Key)
	require.False(t, ok)
}

func TestEmptySetGlb(t *testing.T) {
	var x Set[Key[int]]
	v, ok := x.GreatestLowerBound(AsKey(42))
	require.Equal(t, 0, v.Key)
	require.False(t, ok)
}

func TestNilSetGlb(t *testing.T) {
	var x *Set[Key[int]]
	v, ok := x.GreatestLowerBound(AsKey(42))
	require.Equal(t, 0, v.Key)
	require.False(t, ok)
}

func TestSetLub(t *testing.T) {
	var x *Set[Key[int]]
	for i := 0; i < 20; i += 2 {
		x = x.Add(AsKey(i))
	}
	v, ok := x.LeastUpperBound(AsKey(3))
	require.Equal(t, 4, v.Key)
	require.True(t, ok)
	v, ok = x.LeastUpperBound(AsKey(22))
	require.Equal(t, 0, v.Key)
	require.False(t, ok)
	v, ok = x.LeastUpperBound(AsKey(-1))
	require.Equal(t, 0, v.Key)
	require.True(t, ok)
}

func TestSetGlb(t *testing.T) {
	var x *Set[Key[int]]
	for i := 0; i < 20; i += 2 {
		x = x.Add(AsKey(i))
	}
	v, ok := x.GreatestLowerBound(AsKey(3))
	require.Equal(t, 2, v.Key)
	require.True(t, ok)
	v, ok = x.GreatestLowerBound(AsKey(22))
	require.Equal(t, 18, v.Key)
	require.True(t, ok)
	v, ok = x.GreatestLowerBound(AsKey(-1))
	require.Equal(t, 0, v.Key)
	require.False(t, ok)
}

func TestNilSetIsEmpty(t *testing.T) {
	var x *Set[Key[int]]
	require.True(t, x.IsEmpty())
}

func TestEmptySetIsEmpty(t *testing.T) {
	var x Set[Key[int]]
	require.True(t, x.IsEmpty())
}

func TestSetNotEmpty(t *testing.T) {
	x := EmptySet[Key[int]]().Add(AsKey(42))
	require.False(t, x.IsEmpty())
}

func TestSetIter(t *testing.T) {
	var x *Set[Key[int]]
	for i := 0; i < 20; i += 2 {
		x = x.Add(AsKey(i))
	}
	iter := x.Iter()
	i := 0
	for iter.Next() {
		require.Equal(t, i, iter.Current().Key)
		i += 2
	}
	require.Equal(t, 20, i)
}

func TestNilSetIter(t *testing.T) {
	var x *Set[Key[int]]
	iter := x.Iter()
	i := 0
	for iter.Next() {
		require.Fail(t, "should not enter loop")
	}
	require.Equal(t, 0, i)
}

func TestSetIterGte(t *testing.T) {
	var x *Set[Key[int]]
	for i := 0; i < 20; i += 2 {
		x = x.Add(AsKey(i))
	}
	iter := x.IterGte(AsKey(3))
	i := 4
	for iter.Next() {
		require.Equal(t, i, iter.Current().Key)
		i += 2
	}
	require.Equal(t, 20, i)
}

func TestNilSetIterGte(t *testing.T) {
	var x *Set[Key[int]]
	iter := x.IterGte(AsKey(3))
	i := 0
	for iter.Next() {
		require.Fail(t, "should not enter loop")
	}
	require.Equal(t, 0, i)
}

func TestSetMarshalJson(t *testing.T) {
	var x *Set[Key[int]]
	var expected []int
	for i := 0; i < 20; i += 2 {
		x = x.Add(AsKey(i))
		expected = append(expected, i)
	}
	serialized, err := json.Marshal(x)
	require.NoError(t, err)
	var actual []int
	err = json.Unmarshal(serialized, &actual)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestStringSetMarshalJson(t *testing.T) {
	var x *Set[Key[string]]
	expected := []string{"Hello", "World"}
	x = x.Add(AsKey("Hello")).Add(AsKey("World"))

	serialized, err := json.Marshal(x)
	require.NoError(t, err)
	var actual []string
	err = json.Unmarshal(serialized, &actual)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}