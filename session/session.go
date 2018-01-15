package session

import (
	"github.com/satori/go.uuid"
)

type ISession interface {
	ID() uuid.UUID
	Get(field string)
	Set(field string, value string)
}

type Session struct {
	id       uuid.UUID // unique session id
	uid      uint64    // user id
	email    string
	nickname string
}
