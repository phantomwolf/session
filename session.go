package session

import (
	"github.com/satori/go.uuid"
	"log"
	"strconv"
	"time"
)

const (
	expireKey  = "_expire"
	sessionKey = "_session"
)

type Session struct {
	// unique session id, used as redis key
	id     uuid.UUID
	uid    uint64
	expire time.Time
	values map[string]string
}

// expire after mins minutes
func New(uid uint64, mins int64) (*Session, error) {
	id, err := uuid.NewV4()
	if err != nil {
		log.Printf("[session/session.go] UUID(version 4) generation failed\n")
		return nil, err
	}

	sess := &Session{
		id:     id,
		uid:    uid,
		values: map[string]string{},
	}
	sess.SetExpireAfter(time.Minute * time.Duration(mins))
	return sess, nil
}

func FromStorage(key string, data map[string]string) (*Session, error) {
	uid, err := strconv.ParseUint(key, 10, 64)
	if err != nil {
		log.Printf("[session/session.go] Invalid uid %s: %s\n", key, err.Error())
		return nil, err
	}

	id, err := uuid.FromString(data[sessionKey])
	if err != nil {
		log.Printf("[session/session.go] Invalid session id %s: %s\n", data[sessionKey], err.Error())
		return nil, err
	}
	delete(data, sessionKey)

	expire, err := time.Parse(time.RFC1123, data[expireKey])
	if err != nil {
		log.Printf("[session/session.go] Invalid date %s: %s\n", data[expireKey], err.Error())
		return nil, err
	}
	delete(data, expireKey)

	sess := &Session{
		id:     id,
		uid:    uid,
		expire: expire,
		values: data,
	}
	return sess, nil
}

func (sess *Session) ToStorage() (string, map[string]interface{}) {
	data := map[string]interface{}{}
	for k, v := range sess.values {
		data[k] = v
	}
	data[expireKey] = sess.expire.Format(time.RFC1123)
	data[sessionKey] = sess.ID()
	key := strconv.FormatUint(sess.uid, 10)
	return key, data
}

func (sess *Session) ID() string {
	return sess.id.String()
}

func (sess *Session) UID() uint64 {
	return sess.uid
}

func (sess *Session) SetUID(uid uint64) {
	sess.uid = uid
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

func (sess *Session) Expired() bool {
	if sess.expire.Before(time.Now()) {
		return true
	}
	return false
}

func (sess *Session) SetExpire(t time.Time) {
	sess.expire = t
}

func (sess *Session) SetExpireAfter(d time.Duration) {
	sess.SetExpire(time.Now().Add(d))
}
