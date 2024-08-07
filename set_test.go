package persistent

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestEmptySetContains(t *testing.T) {
	var x Set[int]
	require.False(t, x.Contains(42))
}

func TestNilSetContains(t *testing.T) {
	var x *Set[int]
	require.False(t, x.Contains(42))
}

func TestEmptySetRemove(t *testing.T) {
	var x Set[string]
	pSet := x.Remove("2")
	require.Equal(t, 0, pSet.Size())
}

func TestNilSetRemove(t *testing.T) {
	var x *Set[string]
	pSet := x.Remove("2")
	require.Equal(t, 0, pSet.Size())
}

func TestEmptySetAdd(t *testing.T) {
	var x Set[string]
	pSet := x.Add("2")
	require.Equal(t, 1, pSet.Size())
	require.True(t, pSet.Contains("2"))
}

func TestNilSetAdd(t *testing.T) {
	var x *Set[string]
	pSet := x.Add("2")
	require.Equal(t, 1, pSet.Size())
	require.True(t, pSet.Contains("2"))
}

func TestSetAddTwice(t *testing.T) {
	var x *Set[int]
	x = x.Add(42)
	y := x.Add(42)
	if x != y {
		require.Fail(t, "Adding twice modified pointer")
	}
}

func TestSetAddRemove(t *testing.T) {
	var x *Set[int]
	x = x.Add(4).Add(2).Remove(4)
	require.Equal(t, 1, x.Size())
	require.True(t, x.Contains(2))
	require.False(t, x.Contains(1))
}

func TestEmptySetLub(t *testing.T) {
	var x Set[int]
	v, ok := x.LeastUpperBound(42)
	require.Equal(t, 0, v)
	require.False(t, ok)
}

func TestNilSetLub(t *testing.T) {
	var x *Set[int]
	v, ok := x.LeastUpperBound(42)
	require.Equal(t, 0, v)
	require.False(t, ok)
}

func TestEmptySetGlb(t *testing.T) {
	var x Set[int]
	v, ok := x.GreatestLowerBound(42)
	require.Equal(t, 0, v)
	require.False(t, ok)
}

func TestNilSetGlb(t *testing.T) {
	var x *Set[int]
	v, ok := x.GreatestLowerBound(42)
	require.Equal(t, 0, v)
	require.False(t, ok)
}

func TestSetLub(t *testing.T) {
	var x *Set[int]
	for i := 0; i < 20; i += 2 {
		x = x.Add(i)
	}
	v, ok := x.LeastUpperBound(3)
	require.Equal(t, 4, v)
	require.True(t, ok)
	v, ok = x.LeastUpperBound(22)
	require.Equal(t, 0, v)
	require.False(t, ok)
	v, ok = x.LeastUpperBound(-1)
	require.Equal(t, 0, v)
	require.True(t, ok)
}

func TestSetGlb(t *testing.T) {
	var x *Set[int]
	for i := 0; i < 20; i += 2 {
		x = x.Add(i)
	}
	v, ok := x.GreatestLowerBound(3)
	require.Equal(t, 2, v)
	require.True(t, ok)
	v, ok = x.GreatestLowerBound(22)
	require.Equal(t, 18, v)
	require.True(t, ok)
	v, ok = x.GreatestLowerBound(-1)
	require.Equal(t, 0, v)
	require.False(t, ok)
}

func TestNilSetIsEmpty(t *testing.T) {
	var x *Set[int]
	require.True(t, x.IsEmpty())
}

func TestEmptySetIsEmpty(t *testing.T) {
	var x Set[int]
	require.True(t, x.IsEmpty())
}

func TestSetNotEmpty(t *testing.T) {
	var x *Set[int]
	x = x.Add(42)
	require.False(t, x.IsEmpty())
}

func TestSetIter(t *testing.T) {
	var x *Set[int]
	for i := 0; i < 20; i += 2 {
		x = x.Add(i)
	}
	iter := x.Iter()
	i := 0
	for iter.Next() {
		require.Equal(t, i, iter.Current())
		i += 2
	}
	require.Equal(t, 20, i)
}

func TestNilSetIter(t *testing.T) {
	var x *Set[int]
	iter := x.Iter()
	i := 0
	for iter.Next() {
		require.Fail(t, "should not enter loop")
	}
	require.Equal(t, 0, i)
}

func TestSetIterGte(t *testing.T) {
	var x *Set[int]
	for i := 0; i < 20; i += 2 {
		x = x.Add(i)
	}
	iter := x.IterGte(3)
	i := 4
	for iter.Next() {
		require.Equal(t, i, iter.Current())
		i += 2
	}
	require.Equal(t, 20, i)
}

func TestNilSetIterGte(t *testing.T) {
	var x *Set[int]
	iter := x.IterGte(3)
	i := 0
	for iter.Next() {
		require.Fail(t, "should not enter loop")
	}
	require.Equal(t, 0, i)
}

func TestSetMarshalJson(t *testing.T) {
	var x *Set[int]
	var expected []int
	for i := 0; i < 20; i += 2 {
		x = x.Add(i)
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
	var x *Set[string]
	expected := []string{"Hello", "World"}
	x = x.Add("Hello").Add("World")

	serialized, err := json.Marshal(x)
	require.NoError(t, err)
	var actual []string
	err = json.Unmarshal(serialized, &actual)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestSetGetKth(t *testing.T) {
	var s *Set[int]
	for i := 0; i < 20; i += 2 {
		s = s.Add(i)
	}

	for i := 0; i < 10; i++ {
		e, ok := s.GetKthElement(i)
		require.True(t, ok)
		require.Equal(t, i*2, e)
	}

	e, ok := s.GetKthElement(-1)
	require.False(t, ok)
	require.Equal(t, 0, e)

	e, ok = s.GetKthElement(11)
	require.False(t, ok)
	require.Equal(t, 0, e)
}

func TestSetGetKthNil(t *testing.T) {
	var s *Set[int]
	_, ok := s.GetKthElement(0)
	require.False(t, ok)
}
