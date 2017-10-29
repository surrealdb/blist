// Copyright © 2016 Abcum Ltd
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package blist

import "sync"
import "github.com/google/btree"

// List represents an in-memory binary list.
type List struct {
	tree *btree.BTree
	lock sync.RWMutex
}

// Find determines which method is used to seek items in the list.
type Find int8

const (
	// Exact returns an item at a specific version from the list. If the exact
	// item does not exist in the list, then a nil value is returned.
	Exact Find = iota
	// Prev returns the nearest item in the list, where the version number is
	// less than the given version. In a time-series list, this can be used
	// to get the version that was valid before a specified time.
	Prev
	// Next returns the nearest item in the list, where the version number is
	// greater than the given version. In a time-series list, this can be used
	// to get the version that was changed after a specified time.
	Next
	// Upto returns the nearest item in the list, where the version number is
	// less than or equal to the given version. In a time-series list, this can
	// be used to get the version that was current at the specified time.
	Upto
	// Nearest returns an item nearest a specific version in the list. If there
	// is a previous version to the given version, then it will be returned,
	// otherwise it will return the next available version.
	Nearest
)

// New creates a new list
func New() *List {
	return &List{tree: btree.New(2)}
}

// Clr clears all of the items from the list.
func (l *List) Clr() {

	l.lock.Lock()
	defer l.lock.Unlock()

	l.tree = btree.New(2)

}

// Put inserts a new item into the list, ensuring that the list is sorted
// after insertion. If an item with the same version already exists in the
// list, then the value is updated.
func (l *List) Put(ver int64, val []byte) *Item {

	l.lock.Lock()
	defer l.lock.Unlock()

	i := &Item{ver: ver, val: val, list: l}

	l.tree.ReplaceOrInsert(i)

	return i

}

// Del deletes a specific item from the list, returning the previous item
// if it existed. If it did not exist, a nil value is returned.
func (l *List) Del(ver int64, meth Find) *Item {

	l.lock.Lock()
	defer l.lock.Unlock()

	i := l.find(ver, meth)

	if i != nil {

		l.tree.Delete(i)

		i.list = nil

	}

	return i

}

// Exp expunges all items in the list, upto and including the specified
// version, returning the latest version, or a nil value if not found.
func (l *List) Exp(ver int64, meth Find) *Item {

	l.lock.Lock()
	defer l.lock.Unlock()

	i := l.find(ver, meth)

	if i != nil {

		l.tree.DescendLessOrEqual(i, func(v btree.Item) bool {
			l.tree.Delete(v.(*Item))
			return true
		})

	}

	return i

}

// Get gets a specific item from the list. If the exact item does not
// exist in the list, then a nil value is returned.
func (l *List) Get(ver int64, meth Find) *Item {

	l.lock.RLock()
	defer l.lock.RUnlock()

	return l.find(ver, meth)

}

// Len returns the total number of items in the list.
func (l *List) Len() int {

	l.lock.RLock()
	defer l.lock.RUnlock()

	return l.tree.Len()

}

// Min returns the first item in the list. In a time-series list this can be
// used to get the initial version.
func (l *List) Min() *Item {

	l.lock.RLock()
	defer l.lock.RUnlock()

	i := l.tree.Min()

	if i != nil {
		return i.(*Item)
	}

	return nil

}

// Max returns the last item in the list. In a time-series list this can be
// used to get the latest version.
func (l *List) Max() *Item {

	l.lock.RLock()
	defer l.lock.RUnlock()

	i := l.tree.Max()

	if i != nil {
		return i.(*Item)
	}

	return nil

}

// Walk iterates over the list starting at the first version, and continuing
// until the walk function returns true.
func (l *List) Walk(fn func(*Item) bool) {

	l.lock.RLock()
	defer l.lock.RUnlock()

	l.tree.Ascend(func(i btree.Item) bool {
		return !fn(i.(*Item))
	})

}

// ---------------------------------------------------------------------------

func (l *List) find(ver int64, what Find) (i *Item) {

	switch what {

	case Prev: // Get the item below the specified version

		l.tree.DescendLessOrEqual(&Item{ver: ver}, func(v btree.Item) bool {
			if v.(*Item).ver != ver {
				i = v.(*Item)
				return false
			}
			return true
		})

	case Next: // Get the item above the specified version

		l.tree.AscendGreaterOrEqual(&Item{ver: ver}, func(v btree.Item) bool {
			if v.(*Item).ver != ver {
				i = v.(*Item)
				return false
			}
			return true
		})

	case Upto: // Get the item up to the specified version

		l.tree.DescendLessOrEqual(&Item{ver: ver}, func(v btree.Item) bool {
			i = v.(*Item)
			return false
		})

	case Exact: // Get the exact specified version

		if v := l.tree.Get(&Item{ver: ver}); v != nil {
			i = v.(*Item)
		}

	case Nearest: // Get the item nearest the specified version

		if i = l.find(ver, Upto); i == nil {
			i = l.find(ver, Next)
		}

	}

	return

}
