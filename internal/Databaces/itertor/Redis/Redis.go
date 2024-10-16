package iterRedis

import (
	"github.com/redis/go-redis/v9"
)

type RedisIter struct {
	Iter *redis.ScanIterator
}

func (r *RedisIter) Close() bool {
	return true
}
func (r *RedisIter) Next() bool {
	return true
}
func (r *RedisIter) Key() string {
	return ""
}
func (r *RedisIter) First() bool {
	return true
}
func (r *RedisIter) Value() string {
	return "true"
}
func (r *RedisIter) Seek(key string) bool {
	return true
}
func (r *RedisIter) Valid() bool {
	return true
}

func (r *RedisIter) Prev() bool {
	return true
}
