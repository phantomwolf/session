package session

import (
	"github.com/go-redis/redis"
	"testing"
	"time"
)

func TestSaveDelete(t *testing.T) {
	s, err := New(time.Now().Add(time.Hour * 1))
	if err != nil {
		t.Fatalf("Failed to create session: %s\n", err)
		t.FailNow()
	}

	s.Set("uid", "486")
	s.Set("name", "fool")

	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})
	manager := NewManager(client)
	manager.Save(s)
	// Get data from redis
	data, _ := client.HGetAll(s.ID.String()).Result()
	t.Logf("data in redis: %v\n", data)
	// Delete Session
	manager.Delete(s)
	// Get data from redis
	data, _ = client.HGetAll(s.ID.String()).Result()
	t.Logf("data in redis: %v\n", data)
	client.Close()
}

func TestFind(t *testing.T) {
	s, _ := New(time.Now().Add(time.Hour * 1))
	s.Set("uid", "46")
	s.Set("name", "fool")
	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})
	manager := NewManager(client)
	manager.Save(s)
	// Find Session
	session, err := manager.Find(s.ID.String())
	if err != nil {
		t.Fatalf("Failed to find session %s\n", s.ID.String())
		t.FailNow()
	}
	val, ok := session.Get("uid")
	if ok == false || val != "46" {
	}
}
