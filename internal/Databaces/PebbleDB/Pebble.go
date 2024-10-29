package PebbleDB

import (
	"bytes"
	"fmt"
	dbpak "testgui/internal/Databaces"

	"github.com/cockroachdb/pebble"
)

type PebbleDatabase struct {
	DB      *pebble.DB
	Address string
}

func NewDataBasePebble(address string) dbpak.DBClient {
	return &PebbleDatabase{
		Address: address,
	}
}

func (p *PebbleDatabase) Delete(key []byte) error {
	err := p.DB.Delete(key, nil)
	if err != nil {
		return err
	}
	return nil
}

func (p *PebbleDatabase) Open() error {
	var err error
	p.DB, err = pebble.Open(p.Address, &pebble.Options{})
	return err
}

func (p *PebbleDatabase) Close() {
	p.DB.Close()
}

func (p *PebbleDatabase) Add(key, value []byte) error {
	return p.DB.Set(key, value, nil)
}

func (p *PebbleDatabase) Get(key []byte) ([]byte, error) {
	if p.DB == nil {
		return nil, nil
	}

	data, closer, err := p.DB.Get([]byte(key))
	if err != nil {
		return nil, err
	}

	defer closer.Close()

	return data, err
}

func (p *PebbleDatabase) Read(start, end *[]byte, count int) (error, []dbpak.KVData) {
	var Item []dbpak.KVData

	iterOptions := &pebble.IterOptions{}
	if start != nil {
		iterOptions.LowerBound = *start
	}
	if end != nil {

		iterOptions.UpperBound = *end
	}

	iter, err := p.DB.NewIter(iterOptions)
	if err != nil {
		return err, Item
	}
	defer iter.Close()

	cnt := 0
	if end != nil && start == nil {
		iter.Last()

		key := iter.Key()
		value := iter.Value()

		Item = append(Item, dbpak.KVData{Key: key, Value: value})
		cnt++

		for iter.Prev() {
			cnt++
			if cnt > count {
				break
			}
			key := iter.Key()
			value := iter.Value()

			Item = append(Item, dbpak.KVData{Key: key, Value: value})
		}

		for i := 0; i < len(Item)/2; i++ {
			j := len(Item) - i - 1
			temp := Item[i]
			Item[i] = Item[j]
			Item[j] = temp
		}
	} else {
		if start != nil {
			iter.SeekGE(*start)
			iter.Next()
		} else {
			iter.First()
		}

		for iter.Valid() {
			cnt++
			if cnt > count {
				break
			}
			key := iter.Key()
			value := iter.Value()
			Item = append(Item, dbpak.KVData{Key: key, Value: value})
			iter.Next()
		}
	}

	return nil, Item
}

func (l *PebbleDatabase) Search(valueEntry []byte) (error, [][]byte) {
	var data [][]byte

	Iterator, err := l.DB.NewIter(nil)
	if err != nil {
		return err, data
	}

	defer l.Close()

	if !Iterator.First() {
		return fmt.Errorf("iterator is empty"), data
	}

	for Iterator.Valid() {

		if bytes.Contains(Iterator.Key(), valueEntry) {

			data = append(data, Iterator.Key())

		}
		if !Iterator.Next() {
			break
		}
	}

	return nil, data
}
