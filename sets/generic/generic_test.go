package generic_test

import (
	"testing"

	"github.com/killa-beez/gopkgs/sets/generic"
	"github.com/stretchr/testify/assert"
)

func exampleItems() []generic.Item {
	return []generic.Item{new(generic.Item), new(generic.Item), new(generic.Item)}
}

func NewTestItemSet(val ...generic.Item) *generic.ItemSet {
	s := generic.NewItemSet(len(val))
	if len(val) > 0 {
		s.Add(val...)
	}
	return s
}

func TestItemSet_Add(t *testing.T) {
	t.Run("add nothing", func(t *testing.T) {
		s := NewTestItemSet()
		s.Add()
		assert.Equal(t, 0, s.Len())
	})

	t.Run("add values", func(t *testing.T) {
		exItems := exampleItems()
		s := NewTestItemSet()
		s.Add(exItems[0:2]...)
		assert.Equal(t, exItems[0:2], s.Values())
		s.Add(exItems[0], exItems[2])
		assert.Equal(t, exItems, s.Values())
	})
}

func TestItemSet_Values(t *testing.T) {
	exItems := exampleItems()
	s := NewTestItemSet(exItems...)
	assert.ElementsMatch(t, s.Values(), exItems)
}

func TestItemSet_Len(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		s := NewTestItemSet()
		assert.Equal(t, 0, s.Len())
	})

	t.Run("with items", func(t *testing.T) {
		s := NewTestItemSet(exampleItems()...)
		assert.Equal(t, 3, s.Len())
	})
}

func TestItemSet_Del(t *testing.T) {
	exItems := exampleItems()
	s := NewTestItemSet(exItems...)
	s.Del(exItems[0])
	assert.ElementsMatch(t, exItems[1:], s.Values())
	s.Del(exItems[0], exItems[1])
	assert.ElementsMatch(t, exItems[2:], s.Values())
	s.Del(exItems[0])
	assert.ElementsMatch(t, exItems[2:], s.Values())
}

func TestItemSet_Contains(t *testing.T) {
	exItems := exampleItems()
	s := NewTestItemSet(exItems...)
	assert.True(t, s.Contains(exItems[0]))
	s.Del(exItems[0])
	assert.False(t, s.Contains(exItems[0]))
	assert.True(t, s.Contains(exItems[1], exItems[2]))
	assert.False(t, s.Contains(exItems[1], exItems[2], exItems[0]))
}

func TestItemSet_Clone(t *testing.T) {
	exItems := exampleItems()
	s := NewTestItemSet(exItems...)
	s2 := s.Clone()
	s2.Del(exItems[0])
	assert.ElementsMatch(t, exItems[1:], s2.Values())
	assert.ElementsMatch(t, exItems, s.Values())
}

func TestItemSet_Union(t *testing.T) {
	exItems := exampleItems()
	s1 := NewTestItemSet(exItems[:1]...)
	s2 := NewTestItemSet(exItems[1:]...)
	assert.ElementsMatch(t, exItems, s1.Union(s2).Values())
	assert.ElementsMatch(t, exItems[:1], s1.Values())
	assert.ElementsMatch(t, exItems[1:], s2.Values())
}

func TestItemSet_Intersection(t *testing.T) {
	exItems := exampleItems()
	s1 := NewTestItemSet(exItems[:2]...)
	s2 := NewTestItemSet(exItems[1:]...)
	assert.ElementsMatch(t, []generic.Item{exItems[1]}, s1.Intersection(s2).Values())
	assert.ElementsMatch(t, exItems[:2], s1.Values())
	assert.ElementsMatch(t, exItems[1:], s2.Values())
}

func TestItemSet_Diff(t *testing.T) {
	exItems := exampleItems()
	s1 := NewTestItemSet(exItems[:2]...)
	s2 := NewTestItemSet(exItems[1:]...)
	assert.ElementsMatch(t, exItems[:1], s1.Diff(s2).Values())
	assert.ElementsMatch(t, exItems[1:2], s2.Diff(s1).Values())
	assert.ElementsMatch(t, exItems[:2], s1.Values())
	assert.ElementsMatch(t, exItems[1:], s2.Values())
}
