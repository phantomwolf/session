package session

import (
	"errors"

	log "github.com/sirupsen/logrus"
)

type Repository interface {
	Add(sess *Session) error
	Update(sess *Session) error
	Remove(id string) error
	Find(id string) (*Session, error)
}

type repository struct {
	storage Storage
}

func NewRepository(storage Storage) Repository {
	return &repository{storage: storage}
}

func (repo *repository) Find(id string) (*Session, error) {
	data, err := repo.storage.Load(id)
	if err != nil {
		log.Debugf("[repository.go:Find] Loading %s failed: %s\n", id, err.Error())
		return nil, err
	}
	return Load(id, data)
}

func (repo *repository) Remove(id string) error {
	return repo.storage.Delete(id)
}

func (repo *repository) Update(sess *Session) error {
	if err := repo.storage.Save(sess.Id(), sess.Uid(), sess.ToMap()); err != nil {
		log.Printf("[repository.go:Update] Updating session %s failed: %s\n", sess.Id(), err.Error())
		return err
	}
	return nil
}

func (repo *repository) Add(sess *Session) error {
	if repo.storage.Exists(sess.Id()) {
		log.Printf("[repository.go] Session %s already exists\n", sess.Id())
		return errors.New("Session already exists")
	}

	err := repo.storage.Save(key, data)
	if err != nil {
		log.Printf("[session/repository.go] Session %s saving failure: %s\n", key, err.Error())
		return err
	}
	return nil
}
