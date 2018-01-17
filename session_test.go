package session

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	s, err := New(time.Now().Add(time.Second * 2))
	if err != nil {
		t.Fatalf("Failed to create session: %s\n", err)
		t.FailNow()
	}
	t.Logf("session id: %s\n", s.ID.String())
	t.Logf("expire time: %s\n", s.Expires.Format(time.RFC1123))
	t.Logf("expired: %v\n", s.IsExpired())
	time.Sleep(time.Second * 3)
	t.Logf("expired: %v\n", s.IsExpired())
}

func TestGetSet(t *testing.T) {
	s, err := New(time.Now().Add(time.Hour * 1))
	if err != nil {
		t.Fatalf("Failed to create session: %s\n", err)
		t.FailNow()
	}
	s.Set("uid", "386")
	if val, ok := s.Get("uid"); ok == false || val != "386" {
		t.Fatal("uid is %s, expected 386", val)
		t.FailNow()
	}
}
