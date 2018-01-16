package session

import (
	"github.com/go-redis/redis"
	"github.com/satori/go.uuid"
	"log"
	"time"
)

type Manager struct {
	client *redis.Client
}

func NewManager(client *redis.Client) *Manager {
}

// Once we got session id in cookie, load session from redis
func (m *Manager) LoadSession(id string) (*Session, error) {
	values, err := m.client.HGetAll(id).Result()
	if err != nil || len(values) == 0 {
		log.Fatalf("[session] No such session: %s\n", id)
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

func (m *Manager) SaveSession(s *Session) error {
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
		ok, err = client.HSet(id, k, v).Result()
		if err != nil || ok == false {
			log.Fatalf("[session] Failed to set %s to %s for session %s\n", k, v, id)
			return err
		}
	}
	return nil
}
