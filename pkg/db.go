package db

type DBClient interface {
	Open() bool
	Close() bool
	Add(key, value string) bool
	Get(key string) string
	Exists(key string) bool
}
