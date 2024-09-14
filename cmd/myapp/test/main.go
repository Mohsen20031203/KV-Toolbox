package main

import (
	"fmt"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type KVData struct {
	Key   string
	Value string
}

func testReadRange(db *leveldb.DB, start, end *string, count int) ([]KVData, error) {
	readRange := &util.Range{}
	if start != nil {
		readRange.Start = []byte(*start)
	}
	if end != nil {
		readRange.Limit = []byte(*end)
	}

	iter := db.NewIterator(readRange, nil)
	defer iter.Release()

	var Item []KVData
	cnt := 0

	if end != nil && start == nil {
		iter.Last()
		key := string(iter.Key())
		value := string(iter.Value())
		cnt++
		Item = append(Item, KVData{Key: key, Value: value})
		for iter.Prev() {
			cnt++
			if cnt > count {
				break
			}
			key := string(iter.Key())
			value := string(iter.Value())
			Item = append(Item, KVData{Key: key, Value: value})
		}
	} else {
		for iter.Next() {
			cnt++
			if cnt > count {
				break
			}
			key := string(iter.Key())
			value := string(iter.Value())
			Item = append(Item, KVData{Key: key, Value: value})
		}
	}

	return Item, nil
}

func main() {
	db, err := leveldb.OpenFile("/Users/macbookpro/Documents/GitHub/Bitcoin-prices/example_db", nil)
	if err != nil {
		fmt.Println("Error opening DB:", err)
		return
	}
	defer db.Close()

	_ = db.Put([]byte("aaa"), []byte("apple"), nil)
	_ = db.Put([]byte("bbb"), []byte("banana"), nil)

	start := "d"

	result, err := testReadRange(db, nil, &start, 5)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	for _, item := range result {
		fmt.Printf("Key: %s, Value: %s\n", item.Key, item.Value)
	}
}
