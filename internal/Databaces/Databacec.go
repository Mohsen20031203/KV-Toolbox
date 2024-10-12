package dbpak

import (
	"testgui/internal/Databaces/itertor"

	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type KVData struct {
	Key   string
	Value string
}

type DBClient interface {
	Open() error
	Close()
	Add(key, value string) error
	Get(key string) (string, error)
	Read(start, end *string, count int) (error, []KVData)
	Delete(key string) error
	Iterator(slice *util.Range, ro *opt.ReadOptions) itertor.IterDB
}
