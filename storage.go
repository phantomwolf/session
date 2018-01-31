package session

import (
	"github.com/go-redis/redis"
	"github.com/satori/go.uuid"
	"log"
	"time"
)

type Storage interface {
	// Return nil if not found
	Load(key string) map[string]string
	Save(key string, map[string]string) error
}

type redisStorage struct {
	backend *redis.Client
}

func NewStorage(backend string) Storage {
	var store Storage
	switch backend {
	case "redis":
		store = &redisStorage{backend: }
	}
}

func NewManager(client *redis.Client) *Manager {
	return &Manager{
		client: client,
	}
}

// Once we got session id in cookie, load session from redis
func (m *Manager) Find(id string) (*Session, error) {
	values, err := m.client.HGetAll(id).Result()
	if err != nil || len(values) == 0 {
		log.Fatalf("[session] Session %s not found\n", id)
		return nil, err
	}

	guid, err := uuid.FromString(id)
	if err != nil {
		log.Fatalf("[session] Invalid session id: %s\n", id)
		return nil, err
	}

	expires, err := time.Parse(time.RFC1123, values["Expires"])
	if err != nil {
		log.Fatalf("[session] Invalid time: %s\n", values["Expires"])
		return nil, err
	}
	delete(values, "Expires")

	return &Session{
		ID:      guid,
		Values:  values,
		Expires: expires,
	}, nil
}

// Save session to redis
func (m *Manager) Save(s *Session) error {
	id := s.ID.String()
	ok, err := m.client.HSet(
		id,
		"Expires",
		s.Expires.Format(time.RFC1123),
	).Result()
	if err != nil || ok == false {
		log.Fatalf("[session] Failed to set expire time to %s for session %s\n", s.Expires.Format(time.RFC1123), id)
		return err
	}

	for k, v := range s.Values {
		ok, err = m.client.HSet(id, k, v).Result()
		if err != nil || ok == false {
			log.Fatalf("[session] Failed to set %s to %s for session %s\n", k, v, id)
			return err
		}
	}
	return nil
}

// Delete session
func (m *Manager) Delete(s *Session) error {
	id := s.ID.String()
	ret, err := m.client.Del(id).Result()
	if err != nil {
		log.Fatalf("[session] Failed to delete session %s\n", id)
		return err
	}
	if ret == 0 {
		log.Printf("[session] Session %s already deleted\n", id)
	}
	return nil
}
