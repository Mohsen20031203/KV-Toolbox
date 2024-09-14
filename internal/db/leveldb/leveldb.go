package leveldbb

import (
	"fmt"
	dbpak "testgui/internal/db"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type ConstantDatabase struct {
	Address string
	DB      *leveldb.DB
}

func NewDataBase(address string) dbpak.DBClient {
	return &ConstantDatabase{
		Address: address,
	}
}

func (constant *ConstantDatabase) Delet(key string) error {
	err := constant.DB.Delete([]byte(key), nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *ConstantDatabase) Open() error {

	var err error
	opts := &opt.Options{
		ReadOnly: true,
	}
	c.DB, err = leveldb.OpenFile(c.Address, opts)
	if err != nil {
		return err
	}

	return nil
}

func (constant *ConstantDatabase) Close() {
	if constant.DB != nil {
		constant.DB.Close()
	}
}

func (constant *ConstantDatabase) Add(key, value string) error {
	if constant.DB == nil {
		return fmt.Errorf("database not opened")
	}
	return constant.DB.Put([]byte(key), []byte(value), nil)
}

func (constant *ConstantDatabase) Get(key string) string {
	if constant.DB == nil {
		return ""
	}
	data, err := constant.DB.Get([]byte(key), nil)
	if err != nil {
		return ""
	}
	return string(data)
}

func (c *ConstantDatabase) Read(start, end *string, count int) (error, []dbpak.KVData) {
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
		key := string(iter.Key())                                 // new
		value := string(iter.Value())                             // new
		Item = append(Item, dbpak.KVData{Key: key, Value: value}) // new
		cnt++                                                     // new
		for iter.Prev() {
			cnt++
			if cnt > count {
				break
			}
			key := string(iter.Key())
			value := string(iter.Value())
			Item = append(Item, dbpak.KVData{Key: key, Value: value})
		}
	} else {
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
