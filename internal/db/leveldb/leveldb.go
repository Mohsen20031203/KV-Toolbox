package leveldbb

import (
	"fmt"
	dbpak "testgui/internal/db"

	"github.com/syndtr/goleveldb/leveldb"
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

func (constant *ConstantDatabase) Open() error {
	var err error
	constant.DB, err = leveldb.OpenFile(constant.Address, nil)
	return err
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

func (c *ConstantDatabase) Read() (error, []dbpak.Database) {
	var Item []dbpak.Database

	c.Open()
	defer c.Close()

	iter := c.DB.NewIterator(nil, nil)
	for iter.Next() {
		key := string(iter.Key())
		value := string(iter.Value())
		Item = append(Item, dbpak.Database{Key: key, Value: value})
	}
	iter.Release()

	return nil, Item
}
