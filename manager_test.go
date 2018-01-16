package session

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
	s.Save(manager)
	// Get data from redis
	data, _ := client.HGetAll(s.ID.String()).Result()
	t.Logf("data in redis: %v\n", data)
	// Delete Session
	s.Delete(manager)
	// Get data from redis
	data, _ = client.HGetAll(s.ID.String()).Result()
	t.Logf("data in redis: %v\n", data)
	client.Close()
}
