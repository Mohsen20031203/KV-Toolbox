package redisdb

import (
	"context"
	"fmt"
	dbpak "testgui/internal/Databaces"
	"testgui/internal/Databaces/itertor"
	iterRedis "testgui/internal/Databaces/itertor/Redis"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisDatabase struct {
	client    *redis.Client
	ctx       context.Context
	Host      string
	Port      int
	DataAlias int
	Username  string
	Password  string
	Timeout   uint8
}

func NewDataBaseRedis(host string, port int, dataAlias int, username string, password string, timeout uint8) dbpak.DBClient {
	return &RedisDatabase{
		Host:      host,
		Port:      port,
		DataAlias: dataAlias,
		Username:  username,
		Password:  password,
		Timeout:   timeout,
	}
}

func (r *RedisDatabase) Open() error {
	addr := fmt.Sprintf("%s:%d", r.Host, r.Port)

	r.client = redis.NewClient(&redis.Options{
		Addr:        addr,
		Username:    r.Username,
		Password:    r.Password,
		DB:          r.DataAlias,
		DialTimeout: time.Duration(r.Timeout),
	})
	r.ctx = context.Background()

	_, err := r.client.Ping(r.ctx).Result()
	return err
}

func (r *RedisDatabase) Add(key string, value string) error {
	return r.client.Set(r.ctx, key, value, 0).Err()
}

func (r *RedisDatabase) Close() {
	r.Close()
}

func (r *RedisDatabase) Delete(key string) error {
	return r.client.Del(r.ctx, key).Err()
}

func (r *RedisDatabase) Get(key string) (string, error) {
	return r.client.Get(r.ctx, key).Result()
}

func (r *RedisDatabase) Iterator(start, end *string) itertor.IterDB {

	r.client.Scan(r.ctx, 0, "*", 10)

	return &iterRedis.RedisIter{}
}

func (r *RedisDatabase) Read(start, end *string, count int) (error, []dbpak.KVData) {
	return nil, nil
}
