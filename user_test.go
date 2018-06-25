package main

import (
	"fmt"
	"testing"
)

func TestCRUD(t *testing.T) {
	r := NewUserRepository()

	if r.Get("u1") != nil {
		t.Error(`Get "u1" must return nil`)
	}

	u1 := NewUser()
	u1.LoginId = "u1"

	if r.Put(u1) != nil {
		t.Error("Put u1 must return nil")
	}

	if r.Get("u1") != u1 {
		t.Error(`Get "u1" must return u1`)
	}

	if r.Get("invalid") != nil {
		t.Error(`Get "u2" must return nil`)
	}
}

func TestConcurrentAccess(t *testing.T) {
	r := NewUserRepository()

	go func() {
		for i := 1; i <= 100; i++ {
			go test(r, fmt.Sprintf("u1-%04d", i))
		}
	}()
	go func() {
		for i := 1; i <= 100; i++ {
			go test(r, fmt.Sprintf("u2-%04d", i))
		}
	}()
	go func() {
		for i := 1; i <= 100; i++ {
			go test(r, fmt.Sprintf("u3-%04d", i))
		}
	}()
}

func test(r UserRepository, id string) {
	u := NewUser()
	u.LoginId = id
	r.Get(id)
	r.Auth(id, id)
	r.Get(id)
}
