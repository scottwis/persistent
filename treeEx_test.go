package persistent

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
)

type Int int

func (x Int) Less(y Int) bool {
	return x < y
}

func (x *Int) UnmarshalText(text []byte) error {
	value, err := strconv.Atoi(string(text))
	if err == nil {
		*x = Int(value)
	}
	return err
}

type String string

func (x String) Less(y String) bool {
	return x < y
}

func TestExNilEmpty(t *testing.T) {
	var tree *TreeEx[Int, int]
	require.True(t, tree.IsEmpty())
}

func TestExDefaultEmpty(t *testing.T) {
	var tree TreeEx[Int, int]
	require.True(t, tree.IsEmpty())
}

func TestExNilHeight(t *testing.T) {
	var tree *TreeEx[Int, int]
	require.Equal(t, 0, tree.Height())
}

func TestExDefaultHeight(t *testing.T) {
	var tree TreeEx[Int, int]
	require.Equal(t, 0, tree.Height())
}

func TestExNilSize(t *testing.T) {
	var tree *TreeEx[Int, int]
	require.Equal(t, 0, tree.Size())
}

func TestExDefaultSize(t *testing.T) {
	var tree TreeEx[Int, int]
	require.Equal(t, 0, tree.Size())
}

func TestExNilFind(t *testing.T) {
	var tree *TreeEx[Int, int]
	p := tree.Find(2)
	if p == nil {
		require.Fail(t, "wtf?")
	}
	require.Equal(t, Int(0), p.Key())
	require.Equal(t, 0, p.Value())
	require.True(t, p.IsEmpty())
}

func TestExEmptyFind(t *testing.T) {
	var tree TreeEx[Int, int]
	p := tree.Find(2)
	if p == nil {
		require.Fail(t, "wtf?")
	}
	require.Equal(t, Int(0), p.Key())
	require.Equal(t, 0, p.Value())
	require.True(t, p.IsEmpty())
}

func TestExNilDelete(t *testing.T) {
	var tree *TreeEx[Int, int]
	tree2 := tree.Delete(2)
	require.True(t, tree2.IsEmpty())
}

func TestExEmptyDelete(t *testing.T) {
	var tree TreeEx[Int, int]
	tree2 := tree.Delete(2)
	require.True(t, tree2.IsEmpty())
}

func TestExNilGLB(t *testing.T) {
	var tree *TreeEx[Int, int]
	p := tree.GreatestLowerBound(2)
	require.True(t, p.IsEmpty())
	require.Equal(t, Int(0), p.Key())
	require.Equal(t, 0, p.Value())
}

func TestExEmptyGLB(t *testing.T) {
	var tree TreeEx[Int, int]
	p := tree.GreatestLowerBound(2)
	require.True(t, p.IsEmpty())
	require.Equal(t, Int(0), p.Key())
	require.Equal(t, 0, p.Value())
}

func TestExNilLUB(t *testing.T) {
	var tree *TreeEx[Int, int]
	p := tree.LeastUpperBound(2)
	require.True(t, p.IsEmpty())
	require.Equal(t, Int(0), p.Key())
	require.Equal(t, 0, p.Value())
}

func TestExEmptyLUB(t *testing.T) {
	var tree TreeEx[Int, int]
	p := tree.LeastUpperBound(2)
	require.True(t, p.IsEmpty())
	require.Equal(t, Int(0), p.Key())
	require.Equal(t, 0, p.Value())
}

func TestExNilIter(t *testing.T) {
	var tree *TreeEx[Int, int]
	i := tree.Iter()
	require.False(t, i.Next())
}

func TestExEmptyIter(t *testing.T) {
	var tree TreeEx[Int, int]
	i := tree.Iter()
	require.False(t, i.Next())
}

func TestExNilIterGTE(t *testing.T) {
	var tree *TreeEx[Int, int]
	i := tree.IterGte(5)
	require.False(t, i.Next())
}

func TestExEmptyIterGTE(t *testing.T) {
	var tree TreeEx[Int, int]
	i := tree.IterGte(5)
	require.False(t, i.Next())
}

func TestExNilLeft(t *testing.T) {
	var tree *TreeEx[Int, int]
	require.True(t, tree.Left().IsEmpty())
}

func TestExEmptyLeft(t *testing.T) {
	var tree TreeEx[Int, int]
	require.True(t, tree.Left().IsEmpty())
}

func TestExNilRight(t *testing.T) {
	var tree *TreeEx[Int, int]
	require.True(t, tree.Right().IsEmpty())
}

func TestExEmptyRight(t *testing.T) {
	var tree TreeEx[Int, int]
	require.True(t, tree.Right().IsEmpty())
}

func TestExNilKey(t *testing.T) {
	var tree *TreeEx[Int, int]
	require.Equal(t, Int(0), tree.Key())
}

func TestExEmptyKey(t *testing.T) {
	var tree TreeEx[Int, int]
	require.Equal(t, Int(0), tree.Key())
}

func TestExNilValue(t *testing.T) {
	var tree *TreeEx[Int, int]
	require.Equal(t, 0, tree.Value())
}

func TestExEmptyValue(t *testing.T) {
	var tree TreeEx[Int, int]
	require.Equal(t, 0, tree.Value())
}

func TestExNilLeast(t *testing.T) {
	var tree *TreeEx[Int, int]
	p := tree.Least()
	if p == nil {
		require.Fail(t, "wtf?")
	}
	require.Equal(t, Int(0), p.Key())
	require.Equal(t, 0, p.Value())
	require.True(t, p.IsEmpty())
}

func TestExEmptyLeast(t *testing.T) {
	var tree TreeEx[Int, int]
	p := tree.Least()
	if p == nil {
		require.Fail(t, "wtf?")
	}
	require.Equal(t, Int(0), p.Key())
	require.Equal(t, 0, p.Value())
	require.True(t, p.IsEmpty())
}

func TestExNilMost(t *testing.T) {
	var tree *TreeEx[Int, int]
	p := tree.Most()
	if p == nil {
		require.Fail(t, "wtf?")
	}
	require.Equal(t, Int(0), p.Key())
	require.Equal(t, 0, p.Value())
	require.True(t, p.IsEmpty())
}

func TestExEmptyMost(t *testing.T) {
	var tree TreeEx[Int, int]
	p := tree.Most()
	if p == nil {
		require.Fail(t, "wtf?")
	}
	require.Equal(t, Int(0), p.Key())
	require.Equal(t, 0, p.Value())
	require.True(t, p.IsEmpty())
}

func TestExNilUpdate(t *testing.T) {
	var tree *TreeEx[Int, int]
	tree2 := tree.Update(2, 3)
	require.NotEqual(t, tree, tree2)
	require.False(t, tree2.IsEmpty())
	require.Equal(t, 1, tree2.Size())
	require.Equal(t, 1, tree2.Height())
	require.Equal(t, Int(2), tree2.Key())
	require.Equal(t, 3, tree2.Value())
	p := tree2.Find(4)
	if p == nil {
		require.Fail(t, "wtf")
	}
	require.True(t, p.IsEmpty())
	p = tree2.Find(2)
	require.False(t, p.IsEmpty())
	require.Equal(t, Int(2), p.Key())
	require.Equal(t, 3, p.Value())
	p = tree2.Find(1)
	require.True(t, p.IsEmpty())
	i := tree2.Iter()
	require.True(t, i.Next())
	require.Equal(t, Int(2), i.Current().Key())
	require.Equal(t, 3, i.Current().Value())
	require.False(t, i.Next())
}

func TestExEmptyUpdate(t *testing.T) {
	var tree *TreeEx[Int, int]
	tree2 := tree.Update(2, 3)
	require.NotEqual(t, tree, tree2)
	require.False(t, tree2.IsEmpty())
	require.Equal(t, 1, tree2.Size())
	require.Equal(t, 1, tree2.Height())
	require.Equal(t, Int(2), tree2.Key())
	require.Equal(t, 3, tree2.Value())
	p := tree2.Find(4)
	if p == nil {
		require.Fail(t, "wtf")
	}
	require.True(t, p.IsEmpty())
	p = tree2.Find(2)
	require.False(t, p.IsEmpty())
	require.Equal(t, Int(2), p.Key())
	require.Equal(t, 3, p.Value())
	p = tree2.Find(1)
	require.True(t, p.IsEmpty())
	i := tree2.Iter()
	require.True(t, i.Next())
	require.Equal(t, Int(2), i.Current().Key())
	require.Equal(t, 3, i.Current().Value())
	require.False(t, i.Next())
}

func TestExUpdateReplace(t *testing.T) {
	var tree *TreeEx[Int, int]
	tree2 := tree.Update(2, 3).Update(0, 4).Update(4, 5)
	tree3 := tree2.Update(0, 10)
	require.NotEqual(t, tree, tree2)
	require.NotEqual(t, tree2, tree3)
	p := tree3.Find(0)
	require.False(t, p.IsEmpty())
	require.Equal(t, Int(0), p.Key())
	require.Equal(t, 10, p.Value())
	require.False(t, tree3.Left().IsEmpty())
	require.False(t, tree3.Right().IsEmpty())
	require.True(t, tree3.Left().Left().IsEmpty())
	require.True(t, tree3.Left().Right().IsEmpty())
	require.True(t, tree3.Right().Left().IsEmpty())
	require.True(t, tree3.Right().Right().IsEmpty())
}

func TestExRotateLeft(t *testing.T) {
	var tree *TreeEx[Int, int]
	tree2 := tree.Update(1, 1).Update(2, 2).Update(3, 3)
	require.False(t, tree2.IsEmpty())
	require.Equal(t, Int(2), tree2.Key())
	require.Equal(t, Int(1), tree2.Left().Key())
	require.Equal(t, Int(3), tree2.Right().Key())
}

func TestExRotateRightLeft(t *testing.T) {
	var tree *TreeEx[Int, int]
	items := []int{10, 9, 15, 13, 16, 14}
	for _, i := range items {
		tree = tree.Update(Int(i), i)
	}
	require.Equal(t, 13, tree.Value())
	require.Equal(t, 10, tree.Left().Value())
	require.Equal(t, 15, tree.Right().Value())
	require.Equal(t, 9, tree.Left().Left().Value())
	require.True(t, tree.Left().Right().IsEmpty())
	require.Equal(t, 14, tree.Right().Left().Value())
	require.Equal(t, 16, tree.Right().Right().Value())
}

func TestExRotateRight(t *testing.T) {
	var tree *TreeEx[Int, int]
	tree2 := tree.Update(2, 2).Update(1, 1).Update(0, 0)
	require.Equal(t, 1, tree2.Value())
	require.Equal(t, 0, tree2.Left().Value())
	require.Equal(t, 2, tree2.Right().Value())
}

func TestExRotateLeftRight(t *testing.T) {
	var tree *TreeEx[Int, int]
	items := []int{20, 21, 15, 17, 14, 19}
	for _, i := range items {
		tree = tree.Update(Int(i), i)
	}
	require.Equal(t, 17, tree.Value())
	require.Equal(t, 15, tree.Left().Value())
	require.Equal(t, 20, tree.Right().Value())
	require.Equal(t, 14, tree.Left().Left().Value())
	require.True(t, tree.Left().Right().IsEmpty())
	require.Equal(t, 19, tree.Right().Left().Value())
	require.Equal(t, 21, tree.Right().Right().Value())
}

func TestExLeastUpperBound(t *testing.T) {
	var tree *TreeEx[Int, int]
	for i := 0; i < 20; i += 2 {
		tree = tree.Update(Int(i), i)
	}
	p := tree.LeastUpperBound(4)
	require.Equal(t, 4, p.Value())
	p = tree.LeastUpperBound(5)
	require.Equal(t, 6, p.Value())
	p = tree.LeastUpperBound(22)
	require.True(t, p.IsEmpty())
	p = tree.LeastUpperBound(-1)
	require.Equal(t, 0, p.Value())
}

func TestExGreatestLowerBound(t *testing.T) {
	var tree *TreeEx[Int, int]
	for i := 0; i < 20; i += 2 {
		tree = tree.Update(Int(i), i)
	}
	p := tree.GreatestLowerBound(4)
	require.Equal(t, 4, p.Value())
	p = tree.GreatestLowerBound(5)
	require.Equal(t, 4, p.Value())
	p = tree.GreatestLowerBound(22)
	require.Equal(t, 18, p.Value())
	p = tree.GreatestLowerBound(-1)
	require.True(t, p.IsEmpty())
}

func TestExIter(t *testing.T) {
	var tree *TreeEx[Int, int]
	for i := 0; i < 20; i += 2 {
		tree = tree.Update(Int(i), i)
	}
	iter := tree.Iter()
	j := 0
	for iter.Next() {
		require.Equal(t, j, iter.Current().Value())
		j += 2
	}
	require.Equal(t, 20, j)
}

func TestExIterGte(t *testing.T) {
	var tree *TreeEx[Int, int]
	for i := 0; i < 20; i += 2 {
		tree = tree.Update(Int(i), i)
	}
	iter := tree.IterGte(7)
	j := 8
	for iter.Next() {
		require.Equal(t, j, iter.Current().Value())
		j += 2
	}
	require.Equal(t, 20, j)
}

func TestExDelete(t *testing.T) {
	var tree *TreeEx[Int, int]
	for i := 0; i < 20; i += 2 {
		tree = tree.Update(Int(i), i)
	}
	t2 := tree.Delete(7)
	if tree != t2 {
		require.Fail(t, "delete missing changed tree")
	}
	t3 := tree.Delete(4).Delete(18).Delete(0).Delete(-1).Delete(22)
	require.False(t, t3.IsEmpty())
	require.Equal(t, 7, t3.Size())
	for i := 0; i < 20; i += 2 {
		tree = tree.Delete(Int(i))
	}

	require.True(t, tree.IsEmpty())
}

func TestExMarshalJson(t *testing.T) {
	var tree *TreeEx[Int, int]
	expected := make(map[int]int)
	for i := 0; i < 20; i += 2 {
		tree = tree.Update(Int(i), i)
		expected[i] = i
	}
	serialized, err := json.Marshal(tree)
	require.NoError(t, err)
	var actual map[int]int
	err = json.Unmarshal(serialized, &actual)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestExUnmarshalJson(t *testing.T) {
	m := make(map[int]int)
	for i := 0; i < 20; i += 2 {
		m[i] = i
	}
	serialized, err := json.Marshal(m)
	require.NoError(t, err)
	var tree *TreeEx[Int, int]
	err = json.Unmarshal(serialized, &tree)
	require.NoError(t, err)
	require.Equal(t, 10, tree.Size())
	iter := tree.Iter()
	i := 0
	for iter.Next() {
		require.Equal(t, Int(i), iter.Current().Key())
		i += 2
	}
	require.Equal(t, 20, i)
}

func TestExUnMarshalJsonStringKey(t *testing.T) {
	expected := map[string]int{
		"Hello": 2,
		"World": 4,
	}
	serialized, err := json.Marshal(expected)
	require.NoError(t, err)

	var x *TreeEx[String, int]
	err = json.Unmarshal(serialized, &x)
	require.NoError(t, err)
	require.Equal(t, 2, x.Find("Hello").Value())
	require.Equal(t, 4, x.Find("World").Value())
}
