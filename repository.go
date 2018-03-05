package session

import (
	"errors"
	"log"
	"strconv"
)

type Repository interface {
	Add(sess *Session) error
	Update(sess *Session) error
	Remove(uid uint64) error
	Find(uid uint64) (*Session, error)
}

type repository struct {
	storage Storage
}

func NewRepository(storage Storage) Repository {
	return &repository{storage: storage}
}

func (repo *repository) Find(uid uint64) (*Session, error) {
	key := strconv.FormatUint(uid, 10)
	data, err := repo.storage.Load(key)
	if err != nil {
		return nil, err
	}

	sess, err := FromStorage(key, data)
	return sess, err
}

func (repo *repository) Remove(uid uint64) error {
	key := strconv.FormatUint(uid, 10)
	err := repo.storage.Delete(key)
	return err
}

func (repo *repository) Update(sess *Session) error {
	key, data := sess.ToStorage()
	err := repo.storage.Save(key, data)
	if err != nil {
		log.Printf("[session/repository.go] Session %s saving failure\n", key)
		return err
	}
	return nil
}

func (repo *repository) Add(sess *Session) error {
	key, data := sess.ToStorage()
	if repo.storage.Exists(key) {
		log.Printf("[session/repository.go] Session %s already exists\n", key)
		return errors.New("Session already exists")
	}

	err := repo.storage.Save(key, data)
	if err != nil {
		log.Printf("[session/repository.go] Session %s saving failure: %s\n", key, err.Error())
		return err
	}
	return nil
}
