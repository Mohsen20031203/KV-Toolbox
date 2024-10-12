package iterleveldb

import (
	"github.com/syndtr/goleveldb/leveldb/iterator"
)

type Model struct {
	Iter iterator.Iterator
}

func (m *Model) Next() bool {
	result := m.Iter.Next()
	return result
}

func (m *Model) Key() string {
	text := m.Iter.Key()
	return string(text)
}

func (m *Model) First() bool {
	result := m.Iter.First()
	return result
}

func (m *Model) Value() string {
	result := m.Iter.Value()
	return string(result)
}

func (m *Model) Prev() bool {
	result := m.Iter.Prev()
	return result
}
