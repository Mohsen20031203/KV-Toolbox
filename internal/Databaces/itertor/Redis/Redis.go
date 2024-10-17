package iterRedis

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

type RedisIter struct {
	Ctxx       context.Context
	ClientIter *redis.Client
	Cursor     uint64
}

func (r *RedisIter) Close() bool {
	return true
}
func (r *RedisIter) Next() bool {

	key, cursor, err := r.ClientIter.Scan(r.Ctxx, r.Cursor, "*", 5).Result()
	if err != nil {
		log.Fatalf("error %v", err)
	}
	if cursor == 0 {
		return false
	}
	_ = key
	r.Cursor = cursor

	return true
}

func (r *RedisIter) Key() string {
	keys, _, err := r.ClientIter.Scan(r.Ctxx, r.Cursor, "*", 1).Result()
	if err != nil {
		log.Fatalf("خطا در اسکن: %v", err)
	}
	return keys[0]

}

func (r *RedisIter) First() bool {
	r.Cursor = 0
	return true
}

func (r *RedisIter) Value() string {
	keys, _, err := r.ClientIter.Scan(r.Ctxx, r.Cursor, "*", 1).Result()
	if err != nil {
		log.Fatalf("error Scan in func value %v", err)
	}
	result, err := r.ClientIter.Get(r.Ctxx, keys[0]).Result()
	if err != nil {
		return string(err.Error())
	}
	return result
}

func (r *RedisIter) Seek(startKey string) bool {
	var cursor uint64
	for {
		keys, newCursor, err := r.ClientIter.Scan(r.Ctxx, cursor, "*", 10).Result()
		if err != nil {
			return false
		}

		for _, key := range keys {
			if key == startKey {

				r.Cursor = newCursor
				return true
			}
		}

		if newCursor == 0 {
			break
		}

		cursor = newCursor
	}

	return false
}

func (r *RedisIter) Valid() bool {

	keys, newCursor, err := r.ClientIter.Scan(r.Ctxx, r.Cursor, "*", 1).Result()
	if err != nil {
		return false
	}
	if len(keys) == 0 && newCursor == 0 {
		return false
	}
	return true
}

func (r *RedisIter) Prev() bool {
	var cursornow uint64
	for {
		_, newCursor, err := r.ClientIter.Scan(r.Ctxx, cursornow, "*", 1).Result()
		if err != nil {
			return false
		}

		if r.Cursor == newCursor {
			r.Cursor = cursornow
			break
		}

		if newCursor == 0 {
			break
		}

		cursornow = newCursor
	}

	return false
}
