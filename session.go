package main

import (
	"sync"

	"github.com/rs/xid"
)

type Session struct {
	Id   string `json:"id"`
	User *User  `json:"user"`
}

func NewSession(user *User) *Session {
	return &Session{
		Id:   xid.New().String(),
		User: user,
	}
}

type SessionRepository interface {
	Get(string) *Session
	List() []*Session
	Put(*Session) *Session
}

func NewSessionRepository() SessionRepository {
	return &sessionRepository{
		entities: map[string]*Session{},
	}
}

type sessionRepository struct {
	sync.Mutex

	entities map[string]*Session
}

func (r *sessionRepository) Get(id string) *Session {
	r.Lock()
	defer r.Unlock()

	return r.entities[id]
}

func (r *sessionRepository) List() []*Session {
	r.Lock()
	defer r.Unlock()

	a := []*Session{}
	for _, e := range r.entities {
		a = append(a, e)
	}
	return a
}

func (r *sessionRepository) Put(e *Session) *Session {
	r.Lock()
	defer r.Unlock()

	exists := r.entities[e.Id]
	r.entities[e.Id] = e
	return exists
}
