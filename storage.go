package session

import (
	"github.com/go-redis/redis"
	"github.com/satori/go.uuid"
	"log"
	"time"
)

type Storage interface {
	Load(key string) (map[string]string, error)
	Save(key string, data map[string]string) error
	Delete(key string) error
}

type redisStorage struct {
	client *redis.Client
}

func (store *Storage) Load(key string) (map[string]string, error) {
	client := RedisClient()
	data, err := store.client.HGetAll(key).Result()
	if err != nil {
		log.Printf("[session/storage.go] Loading key %s failed: %s\n", key, err.Error())
		return nil, err
	}
	return data, nil
}

func (store *Storage) Save(key string, data map[string]string) error {
	// map[string]string => map[string]interface{}
	tmp := map[string]interface{}
	for k, v := range data {
		tmp[k] = v
	}
	err := store.client.HMSet(key, tmp).Err()
	if err != nil {
		log.Printf("[session/storage.go] Saving key %s failed: %s\n", key, err.Error())
		return err
	}
	return nil
}

func (store *Storage) Delete(key string) error {
	err := store.client.Del(key).Err()
	if err != nil {
		log.Printf("[session/storage.go] Deleting key %s failed: %s\n", key, err.Error())
		return err
	}
	return nil
}
