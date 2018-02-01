package session

import (
	"github.com/satori/go.uuid"
	"time"
)

const (
	keyExpires = "_expires"
)

type Session struct {
	// unique session id, used as redis key
	id     uuid.UUID
	values map[string]string
	store  Storage
}

func New(store Storage) (*Session, error) {
	// Generate version 4 uuid
	id, err := uuid.NewV4()
	if err != nil {
		log.Printf("[session/session.go] UUID(version 4) generation failed\n")
		return nil, err
	}
	values := make(map[string]string)
	sess := &Session{
		id:     id,
		values: values,
		store:  store,
	}
	return sess, nil
}

func Load(id string, store Storage) (*Session, error) {
	guid, err := uuid.FromString(id)
	if err != nil {
		log.Printf("[session/session.go] Invalid session id %s: %s\n", id, err.Error())
		return nil, err
	}

	values, err := store.Load(id)
	if err != nil {
		log.Printf("[session/session.go] Loading session %s failed: %s\n", id, err.Error())
		return nil, err
	}
	sess := &Session{
		id:     guid,
		values: values,
		store:  store,
	}
	return sess, nil
}

func (sess *Session) ID() string {
	return sess.id.String()
}

// Get value
func (sess *Session) GetVal(key string) (string, bool) {
	val, ok := sess.values[key]
	return val, ok
}

// Set value
func (sess *Session) SetVal(key string, val string) {
	sess.values[key] = val
}

// Delete value
func (sess *Session) DelVal(key string) {
	delete(sess.values, key)
}

func (sess *Session) Delete() error {
	id := sess.ID()
	err := sess.store.Delete(id)
	if err != nil {
		log.Printf("[session/session.go] Deleting session %s failed: %s\n", id, err.Error())
	}
	return err
}

func (sess *Session) Save() error {
	id := sess.ID()
	err := sess.store.Save(id, sess.values)
	if err != nil {
		log.Printf("[session/session.go] Saving session %s failed: %s\n", id, err.Error())
	}
	return err
}

func (sess *Session) Expired() bool {
	expires, ok := time.Parse(time.RFC1123, sess.Get(keyExpires))
	if ok == false {
		panic("[session/session.go] No expire time in session")
	}
	if expires.After(time.Now()) {
		return false
	}
	return true
}

func (sess *Session) SetExpire(t time.Time) {
	sess.values[keyExpires] = t.Format(time.RFC1123)
}

func (sess *Session) SetExpireAfter(d time.Duration) {
	expire := time.Now().Add(d)
	sess.SetExpire(expire)
}
