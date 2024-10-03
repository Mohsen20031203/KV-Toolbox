package dbpak

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
	Delet(key string) error
}
