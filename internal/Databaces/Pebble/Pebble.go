package leveldbb

import (
	"fmt"
	dbpak "testgui/internal/Databaces"

	"github.com/cockroachdb/pebble"
)

type PebbleDatabase struct {
	Address string
	DB      *pebble.DB
}

func NewDataBasePebble(address string) dbpak.DBClient {
	return &PebbleDatabase{
		Address: address,
	}
}

func (constant *PebbleDatabase) Delet(key string) error {
	err := constant.DB.Delete([]byte(key), nil)
	if err != nil {
		return err
	}
	return nil
}

func (constant *PebbleDatabase) Open() error {
	var err error
	constant.DB, err = pebble.Open(constant.Address, &pebble.Options{})
	return err
}

func (constant *PebbleDatabase) Close() {
	if constant.DB != nil {
		constant.DB.Close()
	}
}

func (constant *PebbleDatabase) Add(key, value string) error {
	if constant.DB == nil {
		return fmt.Errorf("database not opened")
	}
	return constant.DB.Set([]byte(key), []byte(value), nil)
}

func (constant *PebbleDatabase) Get(key string) string {
	if constant.DB == nil {
		return ""
	}

	data, closer, err := constant.DB.Get([]byte(key))
	if err != nil {
		return ""
	}

	defer closer.Close()

	return string(data)
}

func (c *PebbleDatabase) Read(start, end *string, count int) (error, []dbpak.KVData) {
	var Item []dbpak.KVData

	iterOptions := &pebble.IterOptions{}
	if start != nil {
		iterOptions.LowerBound = []byte(*start)
	}
	if end != nil {
		iterOptions.UpperBound = []byte(*end)
	}

	iter, err := c.DB.NewIter(iterOptions)
	if err != nil {
		return err, Item
	}
	defer iter.Close()

	cnt := 0
	if end != nil && start == nil {
		iter.Last()

		key := string(iter.Key())
		value := string(iter.Value())
		Item = append(Item, dbpak.KVData{Key: key, Value: value})
		cnt++

		for iter.Prev() {
			cnt++
			if cnt > count {
				break
			}
			key := string(iter.Key())
			value := string(iter.Value())
			Item = append(Item, dbpak.KVData{Key: key, Value: value})
		}

		for i := 0; i < len(Item)/2; i++ {
			j := len(Item) - i - 1
			Item[i], Item[j] = Item[j], Item[i]
		}
	} else {
		if start != nil {
			iter.First()
		}
		for iter.First(); iter.Valid(); iter.Next() {
			cnt++
			if cnt > count {
				break
			}
			key := string(iter.Key())
			value := string(iter.Value())
			Item = append(Item, dbpak.KVData{Key: key, Value: value})
		}
	}

	return nil, Item
}
