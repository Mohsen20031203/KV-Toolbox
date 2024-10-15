package iterbadger

import "github.com/dgraph-io/badger/v4"

type BadgerModel struct {
	Iter *badger.Iterator
	Opts *badger.IteratorOptions
}

func (l *BadgerModel) Next() bool {
	l.Iter.Next()
	return l.Iter.Valid()
}

func (l *BadgerModel) Key() string {
	item := l.Iter.Item()
	key := item.Key()
	return string(key)
}

func (l *BadgerModel) First() bool {
	l.Iter.Rewind()
	return l.Iter.Valid()
}

func (l *BadgerModel) Value() string {
	item := l.Iter.Item()
	valCopy, err := item.ValueCopy(nil)
	if err != nil {
		return ""
	}
	return string(valCopy)
}

func (l *BadgerModel) Prev() bool {
	l.Opts.Reverse = true
	l.Iter.Next()
	return true
}

func (l *BadgerModel) Close() bool {
	l.Iter.Close()
	return true
}

func (l *BadgerModel) Seek(key string) bool {
	l.Iter.Seek([]byte(key))
	return true
}
