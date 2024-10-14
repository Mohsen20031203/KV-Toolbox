package badgerDB

import (
	dbpak "testgui/internal/Databaces"
	"testgui/internal/Databaces/itertor"
	iterbadger "testgui/internal/Databaces/itertor/badger"

	"github.com/dgraph-io/badger/v4"
)

type badgerDatabase struct {
	db      *badger.DB
	Address string
}

func NewDataBaseLeveldb(address string) dbpak.DBClient {
	return &badgerDatabase{
		Address: address,
	}
}

func (b *badgerDatabase) Open() error {
	var err error
	b.db, err = badger.Open(badger.DefaultOptions("test"))
	return err
}

func (b *badgerDatabase) Add(key, value string) error {
	return b.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(key), []byte(value))
	})
}

func (b *badgerDatabase) Close() {
	b.db.Close()
}

func (b *badgerDatabase) Get(key string) (string, error) {
	var valORG string
	b.db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		val, err := item.ValueCopy(nil)
		valORG = string(val)
		return nil
	})
	return valORG, nil
}

func (b *badgerDatabase) Delete(key string) error {
	b.db.Update(func(txn *badger.Txn) error {
		err := txn.Delete([]byte(key))
		if err != nil {
			return err
		}
		return nil
	})
	return nil
}

func (b *badgerDatabase) Read(start, end *string, count int) (error, []dbpak.KVData) {
	var kvData []dbpak.KVData
	var err error
	return err, kvData
}

func (b *badgerDatabase) Iterator(start, end *string) itertor.IterDB {
	var it *badger.Iterator
	err := b.db.View(func(txn *badger.Txn) error {
		it = txn.NewIterator(badger.DefaultIteratorOptions)
		return nil
	})
	if err != nil {
		return nil
	}
	return &iterbadger.BadgerModel{
		Iter: *it,
	}
}
