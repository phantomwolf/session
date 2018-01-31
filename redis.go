package session

import (
	"github.com/go-redis/redis"
	"sync"
)

var redisClient *redis.Client
var once sync.Once

func RedisClient() *redis.Client {
	once.Do(func() {
		client := redis.NewClient(&redis.Options{
			Addr:     "127.0.0.1:6379",
			Password: "",
			DB:       0,
		})
		if client == nil {
			panic("Redis connection failed")
		}
		redisClient = client
	})
	return redisClient
}
