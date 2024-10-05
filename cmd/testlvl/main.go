package main

import (
	"fmt"

	"github.com/syndtr/goleveldb/leveldb"
)

func main() {
	DB, err := leveldb.OpenFile("constant.Address", nil)

	if err != nil {
		return
	}
	for i := 0; i < 200; i++ {

		k := fmt.Sprintf("key--%d", i)
		v := fmt.Sprintf("value--%d", i)
		err := DB.Put([]byte(k), []byte(v), nil)
		if err != nil {
			return
		}
	}

}
