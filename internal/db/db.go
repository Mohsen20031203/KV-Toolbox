package dbpak

type Database struct {
	Key   string
	Value string
}

type DBClient interface {
	Open() error
	Close()
	Add(key, value string) error
	Get(key string) string
	Read() (error, []Database)
	Delet(key string) error
}
