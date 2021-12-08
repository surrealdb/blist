// Copyright © SurrealDB Ltd
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

import "github.com/google/btree"

// Item represents an item in an in-memory btree.
type Item struct {
	ver  uint64
	val  []byte
	list *List
}

// Less determines whether an item precedes another item in the list.
func (i *Item) Less(than btree.Item) bool {
	return i.ver < than.(*Item).ver
}

// Ver returns the version of this item in the containing list.
func (i *Item) Ver() uint64 {
	return i.ver
}

// Val returns the value of this item in the containing list.
func (i *Item) Val() []byte {
	return i.val
}

// Set updates the value of this item in the containing list.
func (i *Item) Set(val []byte) *Item {
	i.val = val
	return i
}

// Del deletes the item from any containing list and returns it.
func (i *Item) Del() *Item {

	if i.list != nil {

		i.list.lock.Lock()
		defer i.list.lock.Unlock()

		i.list.tree.Delete(i)

		i.list = nil

	}

	return i

}

// Prev returns the previous item to this item in the list.
func (i *Item) Prev() *Item {

	if i.list != nil {

		i.list.lock.RLock()
		defer i.list.lock.RUnlock()

		return i.list.find(i.ver, Prev)

	}

	return nil

}

// Next returns the next item to this item in the list.
func (i *Item) Next() *Item {

	if i.list != nil {

		i.list.lock.RLock()
		defer i.list.lock.RUnlock()

		return i.list.find(i.ver, Next)

	}

	return nil

}
