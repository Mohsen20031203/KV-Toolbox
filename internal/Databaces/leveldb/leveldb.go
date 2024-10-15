package leveldbb

import (
	dbpak "testgui/internal/Databaces"
	"testgui/internal/Databaces/itertor"
	iterleveldb "testgui/internal/Databaces/itertor/leveldb"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type LeveldbDatabase struct {
	DB      *leveldb.DB
	Address string
}

func NewDataBaseLeveldb(address string) dbpak.DBClient {
	return &LeveldbDatabase{
		Address: address,
	}
}

func (l *LeveldbDatabase) Delete(key string) error {
	err := l.DB.Delete([]byte(key), nil)
	if err != nil {
		return err
	}
	return nil
}

func (l *LeveldbDatabase) Open() error {
	var err error
	l.DB, err = leveldb.OpenFile(l.Address, nil)
	return err
}

func (l *LeveldbDatabase) Close() {
	l.DB.Close()
}

func (l *LeveldbDatabase) Add(key, value string) error {
	return l.DB.Put([]byte(key), []byte(value), nil)
}

func (l *LeveldbDatabase) Get(key string) (string, error) {
	if l.DB == nil {
		return "", nil
	}
	data, err := l.DB.Get([]byte(key), nil)
	if err != nil {
		return "", err
	}
	return string(data), err
}

func (c *LeveldbDatabase) Read(start, end *string, count int) (error, []dbpak.KVData) {
	var Item []dbpak.KVData

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

func (l *LeveldbDatabase) Iterator(start, end *string) itertor.IterDB {
	readRange := &util.Range{}

	if start != nil {
		readRange.Start = []byte(*start)
	}
	if end != nil {
		readRange.Limit = []byte(*end)
	}
	Iter2 := l.DB.NewIterator(readRange, nil)
	return &iterleveldb.LeveldbModel{
		Iter: Iter2,
	}
}
