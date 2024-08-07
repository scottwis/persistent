package persistent

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestEmptySetExContains(t *testing.T) {
	var x SetEx[Int]
	require.False(t, x.Contains(42))
}

func TestNilSetExContains(t *testing.T) {
	var x *SetEx[Int]
	require.False(t, x.Contains(42))
}

func TestEmptySetExRemove(t *testing.T) {
	var x SetEx[String]
	pSetEx := x.Remove("2")
	require.Equal(t, 0, pSetEx.Size())
}

func TestNilSetExRemove(t *testing.T) {
	var x *SetEx[String]
	pSetEx := x.Remove("2")
	require.Equal(t, 0, pSetEx.Size())
}

func TestEmptySetExAdd(t *testing.T) {
	var x SetEx[String]
	pSetEx := x.Add("2")
	require.Equal(t, 1, pSetEx.Size())
	require.True(t, pSetEx.Contains("2"))
}

func TestNilSetExAdd(t *testing.T) {
	var x *SetEx[String]
	pSetEx := x.Add("2")
	require.Equal(t, 1, pSetEx.Size())
	require.True(t, pSetEx.Contains("2"))
}

func TestSetExAddTwice(t *testing.T) {
	var x *SetEx[Int]
	x = x.Add(42)
	y := x.Add(42)
	if x != y {
		require.Fail(t, "Adding twice modified pointer")
	}
}

func TestSetExAddRemove(t *testing.T) {
	var x *SetEx[Int]
	x = x.Add(4).Add(2).Remove(4)
	require.Equal(t, 1, x.Size())
	require.True(t, x.Contains(2))
	require.False(t, x.Contains(1))
}

func TestEmptySetExLub(t *testing.T) {
	var x SetEx[Int]
	v, ok := x.LeastUpperBound(42)
	require.Equal(t, Int(0), v)
	require.False(t, ok)
}

func TestNilSetExLub(t *testing.T) {
	var x *SetEx[Int]
	v, ok := x.LeastUpperBound(42)
	require.Equal(t, Int(0), v)
	require.False(t, ok)
}

func TestEmptySetExGlb(t *testing.T) {
	var x SetEx[Int]
	v, ok := x.GreatestLowerBound(42)
	require.Equal(t, Int(0), v)
	require.False(t, ok)
}

func TestNilSetExGlb(t *testing.T) {
	var x *SetEx[Int]
	v, ok := x.GreatestLowerBound(42)
	require.Equal(t, Int(0), v)
	require.False(t, ok)
}

func TestSetExLub(t *testing.T) {
	var x *SetEx[Int]
	for i := 0; i < 20; i += 2 {
		x = x.Add(Int(i))
	}
	v, ok := x.LeastUpperBound(3)
	require.Equal(t, Int(4), v)
	require.True(t, ok)
	v, ok = x.LeastUpperBound(22)
	require.Equal(t, Int(0), v)
	require.False(t, ok)
	v, ok = x.LeastUpperBound(-1)
	require.Equal(t, Int(0), v)
	require.True(t, ok)
}

func TestSetExGlb(t *testing.T) {
	var x *SetEx[Int]
	for i := 0; i < 20; i += 2 {
		x = x.Add(Int(i))
	}
	v, ok := x.GreatestLowerBound(3)
	require.Equal(t, Int(2), v)
	require.True(t, ok)
	v, ok = x.GreatestLowerBound(22)
	require.Equal(t, Int(18), v)
	require.True(t, ok)
	v, ok = x.GreatestLowerBound(-1)
	require.Equal(t, Int(0), v)
	require.False(t, ok)
}

func TestNilSetExIsEmpty(t *testing.T) {
	var x *SetEx[Int]
	require.True(t, x.IsEmpty())
}

func TestEmptySetExIsEmpty(t *testing.T) {
	var x SetEx[Int]
	require.True(t, x.IsEmpty())
}

func TestSetExNotEmpty(t *testing.T) {
	var x *SetEx[Int]
	x = x.Add(42)
	require.False(t, x.IsEmpty())
}

func TestSetExIter(t *testing.T) {
	var x *SetEx[Int]
	for i := 0; i < 20; i += 2 {
		x = x.Add(Int(i))
	}
	iter := x.Iter()
	i := 0
	for iter.Next() {
		require.Equal(t, Int(i), iter.Current())
		i += 2
	}
	require.Equal(t, 20, i)
}

func TestNilSetExIter(t *testing.T) {
	var x *SetEx[Int]
	iter := x.Iter()
	i := 0
	for iter.Next() {
		require.Fail(t, "should not enter loop")
	}
	require.Equal(t, 0, i)
}

func TestSetExIterGte(t *testing.T) {
	var x *SetEx[Int]
	for i := 0; i < 20; i += 2 {
		x = x.Add(Int(i))
	}
	iter := x.IterGte(3)
	i := 4
	for iter.Next() {
		require.Equal(t, Int(i), iter.Current())
		i += 2
	}
	require.Equal(t, 20, i)
}

func TestNilSetExIterGte(t *testing.T) {
	var x *SetEx[Int]
	iter := x.IterGte(3)
	i := 0
	for iter.Next() {
		require.Fail(t, "should not enter loop")
	}
	require.Equal(t, 0, i)
}

func TestSetExMarshalJson(t *testing.T) {
	var x *SetEx[Int]
	var expected []int
	for i := 0; i < 20; i += 2 {
		x = x.Add(Int(i))
		expected = append(expected, i)
	}
	serialized, err := json.Marshal(x)
	require.NoError(t, err)
	var actual []int
	err = json.Unmarshal(serialized, &actual)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestStringSetExMarshalJson(t *testing.T) {
	var x *SetEx[String]
	expected := []string{"Hello", "World"}
	x = x.Add("Hello").Add("World")

	serialized, err := json.Marshal(x)
	require.NoError(t, err)
	var actual []string
	err = json.Unmarshal(serialized, &actual)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestSetExGetKth(t *testing.T) {
	var s *SetEx[Int]
	for i := 0; i < 20; i += 2 {
		s = s.Add(Int(i))
	}

	for i := 0; i < 10; i++ {
		e, ok := s.GetKthElement(i)
		require.True(t, ok)
		require.Equal(t, i*2, int(e))
	}

	e, ok := s.GetKthElement(-1)
	require.False(t, ok)
	require.Equal(t, 0, int(e))

	e, ok = s.GetKthElement(11)
	require.False(t, ok)
	require.Equal(t, 0, int(e))
}

func TestSetExGetKeyNil(t *testing.T) {
	var s *SetEx[Int]
	_, ok := s.GetKthElement(0)
	require.False(t, ok)
}
