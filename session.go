package session

import (
	"github.com/satori/go.uuid"
	"time"
)

type Session struct {
	// unique session id, used as redis key
	ID      uuid.UUID
	Values  map[string]string
	Expires time.Time
	Store   Storage
}

func New(expires time.Time, store Storage) (*Session, error) {
	// Generate version 4 uuid
	id, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	return &Session{
		ID:      id,
		Values:  map[string]string{},
		Expires: expires,
		Store:   store,
	}, nil
}

// Get value
func (s *Session) Get(key string) (string, bool) {
	val, ok := s.Values[key]
	return val, ok
}

// Set value
func (s *Session) Set(key string, val string) {
	s.Values[key] = val
}

// Detect if the session is expired
func (s *Session) IsExpired() bool {
	if s.Expires.Before(time.Now()) {
		return true
	}
	return false
}
