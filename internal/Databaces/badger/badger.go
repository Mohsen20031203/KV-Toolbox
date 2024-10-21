package badgerDB

import (
	"strings"
	dbpak "testgui/internal/Databaces"
	jsFile "testgui/internal/json"
	jsondata "testgui/internal/json/jsonData"

	"github.com/dgraph-io/badger/v4"
)

type badgerDatabase struct {
	db      *badger.DB
	Address string
}

func NewDataBaseBadger(address string) dbpak.DBClient {
	return &badgerDatabase{
		Address: address,
	}
}

func (b *badgerDatabase) Open() error {
	var err error
	b.db, err = badger.Open(badger.DefaultOptions(b.Address))
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
	var valORG []byte
	err := b.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		val, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}
		valORG = val
		return nil
	})
	return string(valORG), err
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

func (c *badgerDatabase) Read(start, end *string, count int) (error, []dbpak.KVData) {
	var items []dbpak.KVData
	var opts badger.IteratorOptions
	opts.PrefetchSize = count

	if end != nil && start == nil {
		opts.Reverse = true
	}

	err := c.db.View(func(txn *badger.Txn) error {
		iter := txn.NewIterator(opts)
		defer iter.Close()

		cnt := 0

		if end != nil && start == nil {
			iter.Seek([]byte(*end))
			iter.Next()
			item := iter.Item()
			key := item.Key()
			for iter.Seek([]byte(key)); iter.Valid(); iter.Next() {
				cnt++
				if cnt > count {
					break
				}
				item := iter.Item()
				key := item.Key()

				valCopy, err := item.ValueCopy(nil)
				if err != nil {
					return err
				}

				items = append(items, dbpak.KVData{Key: string(key), Value: string(valCopy)})
			}

			for i := 0; i < len(items)/2; i++ {
				j := len(items) - i - 1
				temp := items[i]
				items[i] = items[j]
				items[j] = temp
			}
		} else {

			if start != nil {
				iter.Seek([]byte(*start))
				iter.Next()
			} else {

				iter.Rewind()
			}

			for ; iter.Valid(); iter.Next() {
				cnt++
				if cnt > count {
					break
				}
				item := iter.Item()
				key := item.Key()

				valCopy, err := item.ValueCopy(nil)
				if err != nil {
					return err
				}

				items = append(items, dbpak.KVData{Key: string(key), Value: string(valCopy)})
			}
		}
		return nil
	})
	if err != nil {
		return err, nil
	}

	return nil, items
}

func (l *badgerDatabase) Search(valueEntry string) (error, []string) {
	var data []string
	var opts badger.IteratorOptions

	err := l.db.View(func(txn *badger.Txn) error {
		Iterator := txn.NewIterator(opts)

		Iterator.Rewind()

		for Iterator.Valid() {

			if strings.Contains(string(Iterator.Item().Key()), valueEntry) {

				data = append(data, string(Iterator.Item().Key()))

			}
			Iterator.Next()
		}
		return nil
	})
	if err != nil {
		return err, data
	}
	return nil, data

}

func (l *badgerDatabase) Jsondatabase() jsFile.JsonFile {
	return &jsondata.ConstantJsonFile{}
}
