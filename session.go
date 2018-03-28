package session

import (
	"crypto/rand"
	"errors"
	"math/big"
	"time"

	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

const (
	characters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	length     = 32
)

var (
	ErrSessionMissingUID = errors.New("User id missing in session")
	ErrSessionExists     = errors.New("Session already exists")
)

type Session struct {
	id     []byte
	values map[string]string
}

// New returns a new Session object
func New() (*Session, error) {
	max := big.NewInt(len(characters))
	id := make([]byte, length)
	for i := 0; i < length; i++ {
		id[i] = characters[rand.Int(rand.Reader, max)]
	}
	sess := &Session{
		id:     id,
		values: map[string]string{},
	}
	return sess, nil
}

// ID returns the string representation of session id(uuid)
func (sess *Session) ID() string {
	return string(sess.id)
}

// Get retrieves the value associated with <key>
func (sess *Session) Get(key string) (string, bool) {
	val, ok := sess.values[key]
	return val, ok
}

// Set sets the value associated with <key> to <val>
func (sess *Session) Set(key string, val string) {
	sess.values[key] = val
}

// DelVal deletes value associated with <key>
func (sess *Session) Del(key string) {
	delete(sess.values, key)
}
