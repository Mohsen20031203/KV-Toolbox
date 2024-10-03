package PebbleDB

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

func (constant *PebbleDatabase) Delete(key string) error {
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

func (constant *PebbleDatabase) Get(key string) (string, error) {
	if constant.DB == nil {
		return "", nil
	}

	data, closer, err := constant.DB.Get([]byte(key))
	if err != nil {
		return "", err
	}

	defer closer.Close()

	return string(data), err
}

func (c *PebbleDatabase) Read(start, end *string, count int) (error, []dbpak.KVData) {
	var Item []dbpak.KVData

	err := c.Open()
	if err != nil {
		fmt.Println("err in function Read in databace Pebble")
	}
	defer c.Close()
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
			temp := Item[i]
			Item[i] = Item[j]
			Item[j] = temp
		}
	} else {
		if start != nil {
			iter.SeekGE([]byte(*start))
			iter.Next()
		} else {
			iter.First()
		}

		for iter.Valid() {
			cnt++
			if cnt > count {
				break
			}
			key := string(iter.Key())
			value := string(iter.Value())
			Item = append(Item, dbpak.KVData{Key: key, Value: value})
			iter.Next()
		}
	}

	return nil, Item
}
