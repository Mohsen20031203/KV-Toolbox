package iterbadger

import "github.com/dgraph-io/badger/v4"

type BadgerModel struct {
	Iter *badger.Iterator
}

func (l *BadgerModel) Next() bool {
	return true
}

func (l *BadgerModel) Key() string {
	return string("text")
}

func (l *BadgerModel) First() bool {

	return true
}

func (l *BadgerModel) Value() string {

	return string("result")
}

func (l *BadgerModel) Prev() bool {

	return true
}

func (l *BadgerModel) Close() bool {
	return true
}
