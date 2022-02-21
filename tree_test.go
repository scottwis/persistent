package persistent

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNilEmpty(t *testing.T) {
	var tree *Tree[Key[int], int]
	require.True(t, tree.IsEmpty())
}

func TestDefaultEmpty(t *testing.T) {
	var tree Tree[Key[int], int]
	require.True(t, tree.IsEmpty())
}

func TestNilHeight(t *testing.T) {
	var tree *Tree[Key[int], int]
	require.Equal(t, 0, tree.Height())
}

func TestDefaultHeight(t *testing.T) {
	var tree Tree[Key[int], int]
	require.Equal(t, 0, tree.Height())
}

func TestNilSize(t *testing.T) {
	var tree *Tree[Key[int], int]
	require.Equal(t, 0, tree.Size())
}

func TestDefaultSize(t *testing.T) {
	var tree Tree[Key[int], int]
	require.Equal(t, 0, tree.Size())
}

func TestNilFind(t *testing.T) {
	var tree *Tree[Key[int], int]
	p := tree.Find(AsKey(2))
	if p == nil {
		require.Fail(t, "wtf?")
	}
	require.Equal(t, 0, p.Key().Key)
	require.Equal(t, 0, p.Value())
	require.True(t, p.IsEmpty())
}

func TestEmptyFind(t *testing.T) {
	var tree Tree[Key[int], int]
	p := tree.Find(AsKey(2))
	if p == nil {
		require.Fail(t, "wtf?")
	}
	require.Equal(t, 0, p.Key().Key)
	require.Equal(t, 0, p.Value())
	require.True(t, p.IsEmpty())
}

func TestNilDelete(t *testing.T) {
	var tree *Tree[Key[int], int]
	tree2 := tree.Delete(AsKey(2))
	require.True(t, tree2.IsEmpty())
}

func TestEmptyDelete(t *testing.T) {
	var tree Tree[Key[int], int]
	tree2 := tree.Delete(AsKey(2))
	require.True(t, tree2.IsEmpty())
}

func TestNilGLB(t *testing.T) {
	var tree *Tree[Key[int], int]
	p := tree.GreatestLowerBound(AsKey(2))
	require.True(t, p.IsEmpty())
	require.Equal(t, 0, p.Key().Key)
	require.Equal(t, 0, p.Value())
}

func TestEmptyGLB(t *testing.T) {
	var tree Tree[Key[int], int]
	p := tree.GreatestLowerBound(AsKey(2))
	require.True(t, p.IsEmpty())
	require.Equal(t, 0, p.Key().Key)
	require.Equal(t, 0, p.Value())
}

func TestNilLUB(t *testing.T) {
	var tree *Tree[Key[int], int]
	p := tree.LeastUpperBound(AsKey(2))
	require.True(t, p.IsEmpty())
	require.Equal(t, 0, p.Key().Key)
	require.Equal(t, 0, p.Value())
}

func TestEmptyLUB(t *testing.T) {
	var tree Tree[Key[int], int]
	p := tree.LeastUpperBound(AsKey(2))
	require.True(t, p.IsEmpty())
	require.Equal(t, 0, p.Key().Key)
	require.Equal(t, 0, p.Value())
}

func TestNilIter(t *testing.T) {
	var tree *Tree[Key[int], int]
	i := tree.Iter()
	require.False(t, i.Next())
}

func TestEmptyIter(t *testing.T) {
	var tree Tree[Key[int], int]
	i := tree.Iter()
	require.False(t, i.Next())
}

func TestNilIterGTE(t *testing.T) {
	var tree *Tree[Key[int], int]
	i := tree.IterGte(AsKey(5))
	require.False(t, i.Next())
}

func TestEmptyIterGTE(t *testing.T) {
	var tree Tree[Key[int], int]
	i := tree.IterGte(AsKey(5))
	require.False(t, i.Next())
}

func TestNilLeft(t *testing.T) {
	var tree *Tree[Key[int], int]
	require.True(t, tree.Left().IsEmpty())
}

func TestEmptyLeft(t *testing.T) {
	var tree Tree[Key[int], int]
	require.True(t, tree.Left().IsEmpty())
}

func TestNilRight(t *testing.T) {
	var tree *Tree[Key[int], int]
	require.True(t, tree.Right().IsEmpty())
}

func TestEmptyRight(t *testing.T) {
	var tree Tree[Key[int], int]
	require.True(t, tree.Right().IsEmpty())
}

func TestNilKey(t *testing.T) {
	var tree *Tree[Key[int], int]
	require.Equal(t, 0, tree.Key().Key)
}

func TestEmptyKey(t *testing.T) {
	var tree Tree[Key[int], int]
	require.Equal(t, 0, tree.Key().Key)
}

func TestNilValue(t *testing.T) {
	var tree *Tree[Key[int], int]
	require.Equal(t, 0, tree.Value())
}

func TestEmptyValue(t *testing.T) {
	var tree Tree[Key[int], int]
	require.Equal(t, 0, tree.Value())
}

func TestNilLeast(t *testing.T) {
	var tree *Tree[Key[int], int]
	p := tree.Least()
	if p == nil {
		require.Fail(t, "wtf?")
	}
	require.Equal(t, 0, p.Key().Key)
	require.Equal(t, 0, p.Value())
	require.True(t, p.IsEmpty())
}

func TestEmptyLeast(t *testing.T) {
	var tree Tree[Key[int], int]
	p := tree.Least()
	if p == nil {
		require.Fail(t, "wtf?")
	}
	require.Equal(t, 0, p.Key().Key)
	require.Equal(t, 0, p.Value())
	require.True(t, p.IsEmpty())
}

func TestNilMost(t *testing.T) {
	var tree *Tree[Key[int], int]
	p := tree.Most()
	if p == nil {
		require.Fail(t, "wtf?")
	}
	require.Equal(t, 0, p.Key().Key)
	require.Equal(t, 0, p.Value())
	require.True(t, p.IsEmpty())
}

func TestEmptyMost(t *testing.T) {
	var tree Tree[Key[int], int]
	p := tree.Most()
	if p == nil {
		require.Fail(t, "wtf?")
	}
	require.Equal(t, 0, p.Key().Key)
	require.Equal(t, 0, p.Value())
	require.True(t, p.IsEmpty())
}

func TestNilUpdate(t *testing.T) {
	var tree *Tree[Key[int], int]
	tree2 := tree.Update(AsKey(2), 3)
	require.NotEqual(t, tree, tree2)
	require.False(t, tree2.IsEmpty())
	require.Equal(t, 1, tree2.Size())
	require.Equal(t, 1, tree2.Height())
	require.Equal(t, 2, tree2.Key().Key)
	require.Equal(t, 3, tree2.Value())
	p := tree2.Find(AsKey(4))
	if p == nil {
		require.Fail(t, "wtf")
	}
	require.True(t, p.IsEmpty())
	p = tree2.Find(AsKey(2))
	require.False(t, p.IsEmpty())
	require.Equal(t, 2, p.Key().Key)
	require.Equal(t, 3, p.Value())
	p = tree2.Find(AsKey(1))
	require.True(t, p.IsEmpty())
	i := tree2.Iter()
	require.True(t, i.Next())
	require.Equal(t, 2, i.Current().Key().Key)
	require.Equal(t, 3, i.Current().Value())
	require.False(t, i.Next())
}

func TestEmptyUpdate(t *testing.T) {
	var tree *Tree[Key[int], int]
	tree2 := tree.Update(AsKey(2), 3)
	require.NotEqual(t, tree, tree2)
	require.False(t, tree2.IsEmpty())
	require.Equal(t, 1, tree2.Size())
	require.Equal(t, 1, tree2.Height())
	require.Equal(t, 2, tree2.Key().Key)
	require.Equal(t, 3, tree2.Value())
	p := tree2.Find(AsKey(4))
	if p == nil {
		require.Fail(t, "wtf")
	}
	require.True(t, p.IsEmpty())
	p = tree2.Find(AsKey(2))
	require.False(t, p.IsEmpty())
	require.Equal(t, 2, p.Key().Key)
	require.Equal(t, 3, p.Value())
	p = tree2.Find(AsKey(1))
	require.True(t, p.IsEmpty())
	i := tree2.Iter()
	require.True(t, i.Next())
	require.Equal(t, 2, i.Current().Key().Key)
	require.Equal(t, 3, i.Current().Value())
	require.False(t, i.Next())
}

func TestUpdateReplace(t *testing.T) {
	var tree *Tree[Key[int], int]
	tree2 := tree.Update(AsKey(2), 3).Update(AsKey(0), 4).Update(AsKey(4), 5)
	tree3 := tree2.Update(AsKey(0), 10)
	require.NotEqual(t, tree, tree2)
	require.NotEqual(t, tree2, tree3)
	p := tree3.Find(AsKey(0))
	require.False(t, p.IsEmpty())
	require.Equal(t, 0, p.Key().Key)
	require.Equal(t, 10, p.Value())
	require.False(t, tree3.Left().IsEmpty())
	require.False(t, tree3.Right().IsEmpty())
	require.True(t, tree3.Left().Left().IsEmpty())
	require.True(t, tree3.Left().Right().IsEmpty())
	require.True(t, tree3.Right().Left().IsEmpty())
	require.True(t, tree3.Right().Right().IsEmpty())
}

func TestRotateLeft(t *testing.T) {
	var tree *Tree[Key[int], int]
	tree2 := tree.Update(AsKey(1), 1).Update(AsKey(2), 2).Update(AsKey(3), 3)
	require.False(t, tree2.IsEmpty())
	require.Equal(t, 2, tree2.Key().Key)
	require.Equal(t, 1, tree2.Left().Key().Key)
	require.Equal(t, 3, tree2.Right().Key().Key)
}

func TestRotateRightLeft(t *testing.T) {
	var tree *Tree[Key[int], int]
	items := []int{10, 9, 15, 13, 16, 14}
	for _, i := range items {
		tree = tree.Update(AsKey(i), i)
	}
	require.Equal(t, 13, tree.Value())
	require.Equal(t, 10, tree.Left().Value())
	require.Equal(t, 15, tree.Right().Value())
	require.Equal(t, 9, tree.Left().Left().Value())
	require.True(t, tree.Left().Right().IsEmpty())
	require.Equal(t, 14, tree.Right().Left().Value())
	require.Equal(t, 16, tree.Right().Right().Value())
}

func TestRotateRight(t *testing.T) {
	var tree *Tree[Key[int], int]
	tree2 := tree.Update(AsKey(2), 2).Update(AsKey(1), 1).Update(AsKey(0), 0)
	require.Equal(t, 1, tree2.Value())
	require.Equal(t, 0, tree2.Left().Value())
	require.Equal(t, 2, tree2.Right().Value())
}

func TestRotateLeftRight(t *testing.T) {
	var tree *Tree[Key[int], int]
	items := []int{20, 21, 15, 17, 14, 19}
	for _, i := range items {
		tree = tree.Update(AsKey(i), i)
	}
	require.Equal(t, 17, tree.Value())
	require.Equal(t, 15, tree.Left().Value())
	require.Equal(t, 20, tree.Right().Value())
	require.Equal(t, 14, tree.Left().Left().Value())
	require.True(t, tree.Left().Right().IsEmpty())
	require.Equal(t, 19, tree.Right().Left().Value())
	require.Equal(t, 21, tree.Right().Right().Value())
}

func TestLeastUppperBound(t *testing.T) {
	tree := EmptyTree[Key[int], int]()
	for i := 0; i < 20; i += 2 {
		tree = tree.Update(AsKey(i), i)
	}
	p := tree.LeastUpperBound(AsKey(4))
	require.Equal(t, 4, p.Value())
	p = tree.LeastUpperBound(AsKey(5))
	require.Equal(t, 6, p.Value())
	p = tree.LeastUpperBound(AsKey(22))
	require.True(t, p.IsEmpty())
	p = tree.LeastUpperBound(AsKey(-1))
	require.Equal(t, 0, p.Value())
}

func TestGreatestLowerBound(t *testing.T) {
	tree := EmptyTree[Key[int], int]()
	for i := 0; i < 20; i += 2 {
		tree = tree.Update(AsKey(i), i)
	}
	p := tree.GreatestLowerBound(AsKey(4))
	require.Equal(t, 4, p.Value())
	p = tree.GreatestLowerBound(AsKey(5))
	require.Equal(t, 4, p.Value())
	p = tree.GreatestLowerBound(AsKey(22))
	require.Equal(t, 18, p.Value())
	p = tree.GreatestLowerBound(AsKey(-1))
	require.True(t, p.IsEmpty())
}

func TestIter(t *testing.T) {
	tree := EmptyTree[Key[int], int]()
	for i := 0; i < 20; i += 2 {
		tree = tree.Update(AsKey(i), i)
	}
	iter := tree.Iter()
	j := 0
	for iter.Next() {
		require.Equal(t, j, iter.Current().Value())
		j += 2
	}
	require.Equal(t, 20, j)
}

func TestIterGte(t *testing.T) {
	tree := EmptyTree[Key[int], int]()
	for i := 0; i < 20; i += 2 {
		tree = tree.Update(AsKey(i), i)
	}
	iter := tree.IterGte(AsKey(7))
	j := 8
	for iter.Next() {
		require.Equal(t, j, iter.Current().Value())
		j += 2
	}
	require.Equal(t, 20, j)
}

func TestDelete(t *testing.T) {
	tree := EmptyTree[Key[int], int]()
	for i := 0; i < 20; i += 2 {
		tree = tree.Update(AsKey(i), i)
	}
	t2 := tree.Delete(AsKey(7))
	if tree != t2 {
		require.Fail(t, "delete missing changed tree")
	}
	t3 := tree.Delete(AsKey(4)).Delete(AsKey(18)).Delete(AsKey(0)).Delete(AsKey(-1)).Delete(AsKey(22))
	require.False(t, t3.IsEmpty())
	require.Equal(t, 7, t3.Size())
	for i := 0; i < 20; i += 2 {
		tree = tree.Delete(AsKey(i))
	}

	require.True(t, tree.IsEmpty())
}

func TestMarshalJson(t *testing.T) {
	tree := EmptyTree[Key[int], int]()
	expected := make(map[int]int)
	for i := 0; i < 20; i += 2 {
		tree = tree.Update(AsKey(i), i)
		expected[i] = i
	}
	serialized, err := json.Marshal(tree)
	require.NoError(t, err)
	var actual map[int]int
	err = json.Unmarshal(serialized, &actual)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestUnmarshalJson(t *testing.T) {
	m := make(map[int]int)
	for i := 0; i < 20; i += 2 {
		m[i] = i
	}
	serialized, err := json.Marshal(m)
	require.NoError(t, err)
	var tree *Tree[Key[int], int]
	err = json.Unmarshal(serialized, &tree)
	require.NoError(t, err)
	require.Equal(t, 10, tree.Size())
	iter := tree.Iter()
	i := 0
	for iter.Next() {
		require.Equal(t, i, iter.Current().Key().Key)
		i += 2
	}
	require.Equal(t, 20, i)
}

func TestUnMarshalJsonStringKey(t *testing.T) {
	expected := map[string]int{
		"Hello": 2,
		"World": 4,
	}
	serialized, err := json.Marshal(expected)
	require.NoError(t, err)

	var x *Tree[Key[string], int]
	err = json.Unmarshal(serialized, &x)
	require.NoError(t, err)
	require.Equal(t, 2, x.Find(AsKey("Hello")).Value())
	require.Equal(t, 4, x.Find(AsKey("World")).Value())
}
