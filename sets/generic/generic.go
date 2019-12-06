package generic

import (
	"sync"

	"github.com/cheekybits/genny/generic"
)

// Item the type of the Set
type Item generic.Type

// ItemSet the set of Items
type ItemSet struct {
	mp map[Item]struct{}
	sync.RWMutex
}

// NewItemSet returns a new *ItemSet with the initial size set to size
func NewItemSet(size int) *ItemSet {
	return &ItemSet{
		mp: make(map[Item]struct{}, size),
	}
}

func (s *ItemSet) lock(others ...*ItemSet) func() {
	s.Lock()
	for _, lck := range others {
		lck.Lock()
	}
	return func() {
		s.Unlock()
		for _, lck := range others {
			lck.Unlock()
		}
	}
}

func (s *ItemSet) readLock(others ...*ItemSet) func() {
	s.RLock()
	for _, lck := range others {
		lck.RLock()
	}
	return func() {
		s.RUnlock()
		for _, lck := range others {
			lck.RUnlock()
		}
	}
}

//Values returns all the values in the set in unspecified order
func (s *ItemSet) Values() []Item {
	defer s.readLock()()
	vals := make([]Item, 0, len(s.mp))
	for k := range s.mp {
		vals = append(vals, k)
	}
	return vals
}

// Add adds new values to the set
func (s *ItemSet) Add(val ...Item) {
	defer s.lock()()
	for _, v := range val {
		s.mp[v] = struct{}{}
	}
}

// Del deletes values from the set
func (s *ItemSet) Del(val ...Item) {
	defer s.readLock()()
	for _, v := range val {
		delete(s.mp, v)
	}
}

//Len returns the number of items in the set
func (s *ItemSet) Len() int {
	defer s.readLock()()
	return len(s.mp)
}

//Contains returns true if the set contains all given values
func (s *ItemSet) Contains(val ...Item) bool {
	defer s.readLock()()
	for _, v := range val {
		_, ok := s.mp[v]
		if !ok {
			return false
		}
	}
	return true
}

//Clone returns a new *ItemSet with the same elements
func (s *ItemSet) Clone() *ItemSet {
	defer s.readLock()()
	newSet := NewItemSet(s.Len())
	newSet.Add(s.Values()...)
	return newSet
}

//Union returns a new *ItemSet with elements from both s and s2
func (s *ItemSet) Union(s2 *ItemSet) *ItemSet {
	defer s.readLock(s2)()
	n := s.Clone()
	n.Add(s2.Values()...)
	return n
}

//Intersection returns a new *ItemSet with elements that s and s2 have in common
func (s *ItemSet) Intersection(s2 *ItemSet) *ItemSet {
	defer s.readLock(s2)()
	n := s.Clone()
	for _, v := range n.Values() {
		if !s2.Contains(v) {
			n.Del(v)
		}
	}
	return n
}

//Diff returns a new *ItemSet with the elements that exist in s but not s2
func (s *ItemSet) Diff(s2 *ItemSet) *ItemSet {
	defer s.readLock(s2)()
	n := s.Clone()
	n.Del(s2.Values()...)
	return n
}
