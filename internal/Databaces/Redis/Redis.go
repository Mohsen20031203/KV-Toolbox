package Redisdb

import (
	"context"
	"fmt"
	"log"
	dbpak "testgui/internal/Databaces"
	"testgui/internal/Databaces/itertor"
	iterRedis "testgui/internal/Databaces/itertor/Redis"

	"github.com/redis/go-redis/v9"
)

type RedisDatabase struct {
	client   *redis.Client
	ctx      context.Context
	Host     string
	Port     string
	Username string
	Password string
}

func NewDataBaseRedis(host string, port string, username string, password string) dbpak.DBClient {

	return &RedisDatabase{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
	}
}

func (r *RedisDatabase) Open() error {
	addr := fmt.Sprintf("%s:%s", r.Host, r.Port)

	r.client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Username: r.Username,
		Password: r.Password,
	})
	r.ctx = context.Background()

	_, err := r.client.Ping(r.ctx).Result()
	return err
}

func (r *RedisDatabase) Add(key, value []byte) error {
	return r.client.Set(r.ctx, string(key), value, 0).Err()
}

func (r *RedisDatabase) Close() {}

func (r *RedisDatabase) Delete(key []byte) error {
	return r.client.Del(r.ctx, string(key)).Err()
}

func (r *RedisDatabase) Get(key []byte) ([]byte, error) {
	result, err := r.client.Get(r.ctx, string(key)).Result()
	return []byte(result), err
}

func (r *RedisDatabase) Iterator(start, end *[]byte) itertor.IterDB {

	return &iterRedis.RedisIter{
		Ctxx:       r.ctx,
		ClientIter: r.client,
		Cursor:     0,
	}
}

func (r *RedisDatabase) Read(start, end *[]byte, count int) (error, []dbpak.KVData) {
	var Item []dbpak.KVData
	var cursor uint64
	cnt := 0

	if end != nil && start == nil {

		keys, err := r.client.ZRevRange(r.ctx, "your_sorted_set_key", 0, int64(count-1)).Result()
		if err != nil {
			log.Fatalf("خطا در دریافت کلیدها: %v", err)
		}

		for _, key := range keys {
			value, err := r.client.Get(r.ctx, key).Result()
			if err != nil {
				log.Fatalf("خطا در دریافت مقدار: %v", err)
			}
			Item = append(Item, dbpak.KVData{Key: []byte(key), Value: []byte(value)})
			cnt++
			if cnt >= count {
				break
			}
		}

		for i := 0; i < len(Item)/2; i++ {
			j := len(Item) - i - 1
			Item[i], Item[j] = Item[j], Item[i]
		}
	} else {

		for {
			keys, newCursor, err := r.client.Scan(r.ctx, cursor, "*", int64(count)).Result()
			if err != nil {
				log.Fatalf("خطا در SCAN: %v", err)
				return err, nil
			}

			for _, key := range keys {
				value, err := r.client.Get(r.ctx, key).Result()
				if err != nil {
					log.Fatalf("خطا در دریافت مقدار: %v", err)
				}
				Item = append(Item, dbpak.KVData{Key: []byte(key), Value: []byte(value)})
				cnt++
				if cnt >= count {
					break
				}
			}

			cursor = newCursor
			if cursor == 0 || cnt >= count {
				break
			}
		}
	}

	return nil, Item
}
