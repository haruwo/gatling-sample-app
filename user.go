package main

import (
	"crypto/sha256"
	"fmt"
	"sync"
)

type User struct {
	LoginId  string `json:"loginId"`
	Name     string `json:"name"`
	Password string `json:"-"`
}

func NewUser() *User {
	return &User{}
}

type UserRepository interface {
	Get(string) *User
	Auth(loginId, password string) *User
	List() []*User
	Put(*User) *User
}

func NewUserRepository() UserRepository {
	return &userRepository{
		entities: map[string]*User{},
	}
}

type userRepository struct {
	sync.Mutex

	entities map[string]*User
}

func (r *userRepository) Get(loginId string) *User {
	r.Lock()
	defer r.Unlock()

	return r.entities[loginId]
}

func (r *userRepository) List() []*User {
	r.Lock()
	defer r.Unlock()

	a := []*User{}
	for _, e := range r.entities {
		a = append(a, e)
	}
	return a
}

func (r *userRepository) Put(e *User) *User {
	r.Lock()
	defer r.Unlock()

	exists := r.entities[e.LoginId]
	r.entities[e.LoginId] = e
	return exists
}

func (r *userRepository) Auth(loginId, password string) *User {
	if loginId == "" || password == "" {
		return nil
	} else if exists := r.Get(loginId); exists == nil {
		user := &User{
			LoginId:  loginId,
			Name:     loginId,
			Password: hashOf(password),
		}
		r.Put(user)
		return user

	} else if exists.Password == hashOf(password) {
		return exists

	} else {
		return nil

	}
}

func hashOf(password string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(password)))
}
