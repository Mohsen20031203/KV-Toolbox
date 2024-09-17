package dbpak

import "github.com/syndtr/goleveldb/leveldb"

type KVData struct {
	Key   string
	Value string
}

type DBClient interface {
	Open() error
	Close()
	Add(key, value string) error
	Get(key string) string
	Read(start, end *string, count int) (error, []KVData)
	Delet(key string) error
	GetDB() leveldb.DB
}
