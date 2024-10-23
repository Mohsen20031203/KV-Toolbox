package Redisdb

import (
	"context"
	"fmt"
	"log"
	dbpak "testgui/internal/Databaces"

	"github.com/redis/go-redis/v9"
)

type RedisDatabase struct {
	client   *redis.Client
	ctx      context.Context
	Addres   string
	Username string
	Password string
}

var cursor uint64
var Rediscursor []uint64

func NewDataBaseRedis(addres string, username string, password string) dbpak.DBClient {

	return &RedisDatabase{
		Addres:   addres,
		Username: username,
		Password: password,
	}
}

func (r *RedisDatabase) Open() error {

	r.client = redis.NewClient(&redis.Options{
		Addr:     r.Addres,
		Username: r.Username,
		Password: r.Password,
	})
	r.ctx = context.Background()

	_, err := r.client.Ping(r.ctx).Result()
	return err
}

func (r *RedisDatabase) Add(key, value string) error {
	return r.client.Set(r.ctx, key, value, 0).Err()
}

func (r *RedisDatabase) Close() {}

func (r *RedisDatabase) Delete(key string) error {
	return r.client.Del(r.ctx, key).Err()
}

func (r *RedisDatabase) Get(key string) (string, error) {
	return r.client.Get(r.ctx, string(key)).Result()
}

func (r *RedisDatabase) Read(start, end *string, count int) (error, []dbpak.KVData) {
	var Item []dbpak.KVData
	cnt := 0

	if end != nil && start == nil {

		keys, _, err := r.client.Scan(r.ctx, Rediscursor[0], "*", int64(count)).Result()
		if err != nil {
			log.Fatalf("error %v", err)
			return err, nil
		}

		for _, key := range keys {
			value, err := r.client.Get(r.ctx, key).Result()
			if err != nil {
				log.Fatalf(" error %v", err)
			}
			Item = append(Item, dbpak.KVData{Key: key, Value: value})
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

		keys, newCursor, err := r.client.Scan(r.ctx, cursor, "*", int64(count)).Result()
		if err != nil {
			log.Fatalf("error %v", err)
			return err, nil
		}

		for _, key := range keys {
			value, err := r.client.Get(r.ctx, key).Result()
			if err != nil {
				log.Fatalf("error %v", err)
			}
			Item = append(Item, dbpak.KVData{Key: key, Value: value})
		}

		cursor = newCursor
		Rediscursor = append(Rediscursor, cursor)
	}

	return nil, Item
}

func (l *RedisDatabase) Search(valueEntry string) (error, []string) {
	var data []string
	var cursorOP uint64
	for {

		keys, cursor, err := l.client.Scan(l.ctx, cursorOP, fmt.Sprintf("*%s*", valueEntry), 10).Result()
		if err != nil {
			return err, data
		}
		data = append(data, keys...)

		cursorOP = cursor
		if cursorOP == 0 {
			break
		}
	}

	return nil, data
}
