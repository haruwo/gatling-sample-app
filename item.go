package main

import (
	"sync"

	"github.com/rs/xid"
)

type Item struct {
	Id   string `json:id`
	Name string `json:name`
}

func NewItem() *Item {
	return &Item{
		Id: xid.New().String(),
	}
}

type ItemRepository interface {
	Get(string) *Item
	List() []*Item
	Put(*Item) *Item
}

func NewItemRepository() ItemRepository {
	return &itemRepository{
		entities: map[string]*Item{},
	}
}

type itemRepository struct {
	sync.Mutex

	entities map[string]*Item
}

func (r *itemRepository) Get(id string) *Item {
	r.Lock()
	defer r.Unlock()

	return r.entities[id]
}

func (r *itemRepository) List() []*Item {
	r.Lock()
	defer r.Unlock()

	a := []*Item{}
	for _, e := range r.entities {
		a = append(a, e)
	}
	return a
}

func (r *itemRepository) Put(e *Item) *Item {
	r.Lock()
	defer r.Unlock()

	exists := r.entities[e.Id]
	r.entities[e.Id] = e
	return exists
}
