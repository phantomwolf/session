package session

import (
	"errors"
	"github.com/go-redis/redis"
	"log"
)

type redisBackend struct {
	client *redis.Client
}

func NewRedisBackend() Backend {
	backend := &redisBackend{client: RedisClient()}
	return backend
}

func (backend *redisBackend) Load(key string) (map[string]string, error) {
	data, err := backend.client.HGetAll(key).Result()
	if err != nil {
		log.Printf("[session/storage.go] Loading key %s failed: %s\n", key, err.Error())
		return nil, err
	}
	if data == nil || len(data) == 0 {
		log.Printf("[session/storage.go] No such key %s in redis\n", key)
		return nil, errors.New("No such key")
	}
	return data, nil
}

func (backend *redisBackend) Save(key string, data map[string]string) error {
	// map[string]string => map[string]interface{}
	tmp := map[string]interface{}{}
	for k, v := range data {
		tmp[k] = v
	}
	err := backend.client.HMSet(key, tmp).Err()
	if err != nil {
		log.Printf("[session/storage.go] Saving key %s failed: %s\n", key, err.Error())
		return err
	}
	return nil
}

func (backend *redisBackend) Delete(key string) error {
	err := backend.client.Del(key).Err()
	if err != nil {
		log.Printf("[session/storage.go] Deleting key %s failed: %s\n", key, err.Error())
		return err
	}
	return nil
}
