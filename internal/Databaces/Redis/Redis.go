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
	Port     int
	Username string
	Password string
}

func NewDataBaseRedis(host string, port int, username string, password string) dbpak.DBClient {

	return &RedisDatabase{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
	}
}

func (r *RedisDatabase) Open() error {
	addr := fmt.Sprintf("%s:%d", r.Host, r.Port)

	r.client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Username: r.Username,
		Password: r.Password,
	})
	r.ctx = context.Background()

	_, err := r.client.Ping(r.ctx).Result()
	return err
}

func (r *RedisDatabase) Add(key string, value string) error {
	return r.client.Set(r.ctx, key, value, 0).Err()
}

func (r *RedisDatabase) Close() {}

func (r *RedisDatabase) Delete(key string) error {
	return r.client.Del(r.ctx, key).Err()
}

func (r *RedisDatabase) Get(key string) (string, error) {
	return r.client.Get(r.ctx, key).Result()
}

func (r *RedisDatabase) Iterator(start, end *string) itertor.IterDB {

	return &iterRedis.RedisIter{
		Ctxx:       r.ctx,
		ClientIter: r.client,
		Cursor:     0,
	}
}

func (r *RedisDatabase) Read(start, end *string, count int) (error, []dbpak.KVData) {
	var Item []dbpak.KVData
	var cursor uint64
	cnt := 0

	// حالت پیمایش معکوس (reverse)
	if end != nil && start == nil {
		// پیمایش از انتها به ابتدا (با فرض این که از یک Sorted Set استفاده کنید)
		keys, err := r.client.ZRevRange(r.ctx, "your_sorted_set_key", 0, int64(count-1)).Result()
		if err != nil {
			log.Fatalf("خطا در دریافت کلیدها: %v", err)
		}

		for _, key := range keys {
			value, err := r.client.Get(r.ctx, key).Result()
			if err != nil {
				log.Fatalf("خطا در دریافت مقدار: %v", err)
			}
			Item = append(Item, dbpak.KVData{Key: key, Value: value})
			cnt++
			if cnt >= count {
				break
			}
		}

		// معکوس کردن آرایه
		for i := 0; i < len(Item)/2; i++ {
			j := len(Item) - i - 1
			Item[i], Item[j] = Item[j], Item[i]
		}
	} else {
		// پیمایش عادی با SCAN
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
				Item = append(Item, dbpak.KVData{Key: key, Value: value})
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
