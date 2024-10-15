package iterPebble

import (
	"github.com/cockroachdb/pebble"
)

type PebbleIter struct {
	Iter *pebble.Iterator
}

func (m *PebbleIter) Next() bool {
	result := m.Iter.Next()
	return result
}

func (m *PebbleIter) Key() string {
	text := m.Iter.Key()
	return string(text)
}

func (m *PebbleIter) First() bool {
	result := m.Iter.First()
	return result
}

func (m *PebbleIter) Value() string {
	result := m.Iter.Value()
	return string(result)
}

func (m *PebbleIter) Prev() bool {
	result := m.Iter.Prev()
	return result
}

func (m *PebbleIter) Close() bool {
	result := m.Close()
	return result
}

func (m *PebbleIter) Seek(key string) bool {
	result := m.Iter.SeekGE([]byte(key))
	return result
}

func (m *PebbleIter) Valid() bool {
	return m.Iter.Valid()
}
