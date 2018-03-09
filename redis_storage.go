package session

import (
	"errors"

	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
)

var (
	ErrStorageKeyNotFound = errors.New("No such key")
)

type redisStorage struct {
	limit  int
	client *redis.Client
}

// NewRedisStorage returns a new redisStorage.
// A user can have at most <limit> sessions.
func NewRedisStorage(limit int, client *redis.Client) Storage {
	return &redisStorage{limit: limit, client: client}
}

// Load retrieves a hash associated with <key> from redis
func (rs *redisStorage) Load(key string) (map[string]string, error) {
	if !rs.Exists(key) {
		log.Debugf("[redis_storage.go:Load] No such key %s\n", key)
		return nil, ErrStorageKeyNotFound
	}
	data, err := rs.client.HGetAll(key).Result()
	if err != nil {
		log.Debugf("[redis_storage.go:Load] Loading key %s failed: %s\n", key, err.Error())
		return nil, err
	}
	return data, nil
}

// Save saves a hash to redis. If user already has MaxSessionNum sessions,
// remove the oldest one
func (rs *redisStorage) Save(key string, uid string data map[string]interface{}) error {
	if len, err := rs.client.LLen(uid).Result(); err != nil {
		log.Debugf("[redis_storage.go:Save] LLEN %s failed: %s\n", uid, err.Error())
		return err
	} else if len >= int64(rs.limit) {
		id, err := rs.client.RPop(uid).Result()
		if err != nil {
			log.Debugf("[redis_storage.go:Save] RPOP %s failed: %s\n", uid, err.Error())
			return err
		}
		// The session might already be removed.
		// No need to check result.
		rs.Delete(id)
	}
	// Add session id to user's session list
	if err := rs.client.LPush(uid, key).Err(); err != nil {
		log.Debugf("[redis_storage.go:Save] LPUSH %s failed: %s\n", uid, err.Error())
		return err
	}
	// Save session to redis
	if err := rs.client.HMSet(key, data).Err(); err != nil {
		log.Debugf("[redis_storage.go:Save] Saving key %s failed: %s\n", key, err.Error())
		return err
	}
	return nil
}

func (rs *redisStorage) Delete(key string) error {
	if err := rs.client.Del(key).Err(); err != nil {
		log.Debugf("[redis_storage.go:Delete] Deleting key %s failed: %s\n", key, err.Error())
		return err
	}
	return nil
}

func (rs *redisStorage) Exists(key string) bool {
	res := rs.client.Exists(key).Val()
	return (res != 0)
}
