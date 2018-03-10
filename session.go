package session

import (
	"errors"
	"time"

	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

const (
	uidKey    = "_uid"
	expireKey = "_expire"
)

var (
	ErrSessionMissingUID = errors.New("User id missing in session")
	ErrSessionExists     = errors.New("Session already exists")
)

type Session struct {
	id     uuid.UUID
	uid    string
	expire time.Time
	values map[string]string
}

// New returns a new Session which belongs to user <uid>
// and expires after <dur>
func New(uid string, dur time.Duration) (*Session, error) {
	id, err := uuid.NewV4()
	if err != nil {
		log.Debugln("[session.go:New] UUID(version 4) generation failed")
		return nil, err
	}

	sess := &Session{
		id:     id,
		uid:    uid,
		expire: time.Now().Add(dur),
		values: map[string]string{},
	}
	return sess, nil
}

// Load constructs an existing Session whose id is <key>
// and other fields in <data>
func Load(key string, data map[string]string) (*Session, error) {
	id, err := uuid.FromString(key)
	if err != nil {
		log.Debugf("[session.go:Load] Invalid session id %s: %s\n", key, err.Error())
		return nil, err
	}

	uid, ok := data[uidKey]
	if !ok {
		log.Debugf("[session.go:Load] User id not found in session %s\n", key)
		return nil, ErrSessionMissingUID
	}
	delete(data, uidKey)

	expire, err := time.Parse(time.RFC1123, data[expireKey])
	if err != nil {
		log.Debugf("[session.go:Load] Invalid date %s: %s\n", data[expireKey], err.Error())
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

// ToMap returns a map containing all the data of a Session except session id.
// It's useful for storing the Session to storage.
func (sess *Session) ToMap() map[string]interface{} {
	data := map[string]interface{}{}
	for k, v := range sess.values {
		data[k] = v
	}
	data[expireKey] = sess.expire.Format(time.RFC1123)
	data[uidKey] = sess.UID()
	return data
}

// ID returns the string representation of session id(uuid)
func (sess *Session) ID() string {
	return sess.id.String()
}

// UID returns the user id which the session belongs to
func (sess *Session) UID() string {
	return sess.uid
}

// SetUID sets the Session's uid to <uid>
func (sess *Session) SetUID(uid string) {
	sess.uid = uid
}

// GetVal retrieves the value associated with <key>
func (sess *Session) GetVal(key string) (string, bool) {
	val, ok := sess.values[key]
	return val, ok
}

// SetVal sets the value associated with <key> to <val>
func (sess *Session) SetVal(key string, val string) {
	if key == expireKey || key == uidKey {
		log.Panicf("[session.go:SetVal] Key can't be %s or %s\n", expireKey, uidKey)
	}
	sess.values[key] = val
}

// DelVal deletes value associated with <key>
func (sess *Session) DelVal(key string) {
	delete(sess.values, key)
}

// Expired returns true if Session is already expired; otherwise, return false.
func (sess *Session) Expired() bool {
	if sess.expire.Before(time.Now()) {
		return true
	}
	return false
}

// SetExpire sets the expire time
func (sess *Session) SetExpire(t time.Time) {
	sess.expire = t
}

// SetExpireAfter sets the session to exipre after <d>
func (sess *Session) SetExpireAfter(d time.Duration) {
	sess.SetExpire(time.Now().Add(d))
}
