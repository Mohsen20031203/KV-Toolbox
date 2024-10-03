package leveldbb

import (
	"fmt"
	dbpak "testgui/internal/Databaces"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type LeveldbDatabase struct {
	Address string
	DB      *leveldb.DB
}

func NewDataBaseLeveldb(address string) dbpak.DBClient {
	return &LeveldbDatabase{
		Address: address,
	}
}

func (constant *LeveldbDatabase) Delete(key string) error {
	err := constant.DB.Delete([]byte(key), nil)
	if err != nil {
		return err
	}
	return nil
}

func (constant *LeveldbDatabase) Open() error {
	var err error
	constant.DB, err = leveldb.OpenFile(constant.Address, nil)
	return err
}

func (constant *LeveldbDatabase) Close() {
	if constant.DB != nil {
		constant.DB.Close()
	}
}

func (constant *LeveldbDatabase) Add(key, value string) error {
	if constant.DB == nil {
		return fmt.Errorf("database not opened")
	}
	return constant.DB.Put([]byte(key), []byte(value), nil)
}

func (constant *LeveldbDatabase) Get(key string) (string, error) {
	if constant.DB == nil {
		return "", nil
	}
	data, err := constant.DB.Get([]byte(key), nil)
	if err != nil {
		return "", err
	}
	return string(data), err
}

func (c *LeveldbDatabase) Read(start, end *string, count int) (error, []dbpak.KVData) {
	var Item []dbpak.KVData

	err := c.Open()
	if err != nil {
		fmt.Print("error : Open leveldb in func Read")
	}
	defer c.Close()
	readRange := &util.Range{}
	if start != nil {
		readRange.Start = []byte(*start)
	}
	if end != nil {
		readRange.Limit = []byte(*end)
	}
	iter := c.DB.NewIterator(readRange, nil)
	defer iter.Release()
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
		//reverse items
		for i := 0; i < len(Item)/2; i++ {
			j := len(Item) - i - 1
			temp := Item[i]
			Item[i] = Item[j]
			Item[j] = temp
		}
	} else {
		if start != nil {

			iter.Next()
		}
		for iter.Next() {
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
