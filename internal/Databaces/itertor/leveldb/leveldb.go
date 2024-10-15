package iterleveldb

import (
	"github.com/syndtr/goleveldb/leveldb/iterator"
)

type LeveldbModel struct {
	Iter iterator.Iterator
}

func (l *LeveldbModel) Next() bool {
	result := l.Iter.Next()
	return result
}

func (l *LeveldbModel) Key() string {
	text := l.Iter.Key()
	return string(text)
}

func (l *LeveldbModel) First() bool {
	result := l.Iter.First()
	return result
}

func (l *LeveldbModel) Value() string {
	result := l.Iter.Value()
	return string(result)
}

func (l *LeveldbModel) Prev() bool {
	result := l.Iter.Prev()
	return result
}

func (l *LeveldbModel) Close() bool {
	return true
}

func (l *LeveldbModel) Seek(key string) bool {
	result := l.Iter.Seek([]byte(key))
	return result
}
