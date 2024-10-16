package redisdb

import (
	dbpak "testgui/internal/Databaces"
	"testgui/internal/Databaces/itertor"
	iterRedis "testgui/internal/Databaces/itertor/Redis"
)

type RedisDatabase struct {
	Host     float32
	Port     int16
	Username string
	Password string
	Timeout  uint8
}

func NewDataBaseRedis(host float32, port int16, user string, password string, timeot uint8) dbpak.DBClient {
	return &RedisDatabase{
		Host:     host,
		Port:     port,
		Username: user,
		Password: password,
		Timeout:  timeot,
	}
}

func (r *RedisDatabase) Open() error {
	return nil
}

func (r *RedisDatabase) Add(key string, value string) error {
	return nil
}

func (r *RedisDatabase) Close() {

}

func (r *RedisDatabase) Delete(key string) error {
	return nil
}

func (r *RedisDatabase) Get(key string) (string, error) {
	return "", nil
}

func (r *RedisDatabase) Iterator(start, end *string) itertor.IterDB {
	return &iterRedis.RedisIter{}
}

func (r *RedisDatabase) Read(start, end *string, count int) (error, []dbpak.KVData) {
	return nil, nil
}
