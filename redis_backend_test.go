package session

import (
	"testing"
)

func TestLoad(t *testing.T) {
	// Prepare data
	client := RedisClient()
	defer client.Close()
	key := "loadingtest"
	client.HSet(key, "name", "foo")
	client.HSet(key, "age", 20)

	// Load data from redis
	backend := NewRedisBackend()
	data, err := backend.Load(key)
	if err != nil {
		t.Fatalf("redisBackend: Load failed: %s\n", err.Error())
	}
	if data["name"] != "foo" || data["age"] != "20" {
		t.Fatalf("Data loaded is incorrect: %v\n", data)
	}

	// Load non-existing data
	data, err = backend.Load("nosuchkey")
	t.Logf("err: %v\n", err)
	if err == nil {
		t.Fatalf("redisBackend: Loading non-existing key should fail")
	}

	// Close redis connection and load again
	client.Close()
	data, err = backend.Load(key)
	t.Logf("err: %v\n", err)
	if err == nil {
		t.Fatalf("redisBackend: Loading non-existing key should fail")
	}

	// Clean up
	client.Del(key)
}
