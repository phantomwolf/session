package session

import (
	"errors"
	"github.com/go-redis/redis"
	"log"
)

type RedisStorage struct {
	client *redis.Client
}

func NewRedisStorage(client *redis.Client) *RedisStorage {
	return &RedisStorage{client: client}
}

func (rs *RedisStorage) Load(key string) (map[string]string, error) {
	if !rs.Exists(key) {
		log.Printf("[session/redis_storage.go] No such key %s\n", key)
		return nil, errors.New("No such key")
	}
	data, err := rs.client.HGetAll(key).Result()
	if err != nil {
		log.Printf("[session/redis_storage.go] Loading key %s failed: %s\n", key, err.Error())
		return nil, err
	}
	return data, nil
}

func (rs *RedisStorage) Save(key string, data map[string]interface{}) error {
	err := rs.client.HMSet(key, data).Err()
	if err != nil {
		log.Printf("[session/redis_storage.go] Saving key %s failed: %s\n", key, err.Error())
		return err
	}
	return nil
}

func (rs *RedisStorage) Delete(key string) error {
	err := rs.client.Del(key).Err()
	if err != nil {
		log.Printf("[session/redis_storage.go] Deleting key %s failed: %s\n", key, err.Error())
		return err
	}
	return nil
}

func (rs *RedisStorage) Exists(key string) bool {
	res := rs.client.Exists(key).Val()
	return (res != 0)
}
