package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/xid"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.DefaultLogger)
	r.Mount("/auth", auth())
	r.Mount("/api", api())

	bind := fmt.Sprintf(":%d", port())
	fmt.Printf("Listening on %s\n", bind)
	http.ListenAndServe(bind, r)
}

func port() int {
	if len(os.Args) < 2 {
		return 8000
	} else if port, err := strconv.Atoi(os.Args[1]); err != nil {
		panic(err)
	} else {
		return port
	}
}

func auth() chi.Router {
	r := chi.NewRouter()
	r.Post("/login", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			BadRequest(w)
		} else if user := userRepos.Auth(r.Form.Get("loginId"), r.Form.Get("password")); user == nil {
			Unauthorized(w)
		} else {
			session := NewSession(user)
			sessionRepos.Put(session)
			Render(w, session)
		}
	})
	return r
}

func api() chi.Router {
	r := chi.NewRouter()
	r.Use(RequireAuth(sessionRepos))

	r.Get("/item", func(w http.ResponseWriter, r *http.Request) {
		if err := Render(w, itemRepos.List()); err != nil {
			Error(w, "Not found")
		}
	})
	r.Get("/item/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if e := itemRepos.Get(id); e == nil {
			NotFound(w)
		} else if err := Render(w, e); err != nil {
			Error(w, "Not found")
		}
	})
	r.Post("/item", func(w http.ResponseWriter, r *http.Request) {
		o, err := Decode(r, &Item{})
		if err != nil {
			log.Printf("Decode Error: %v", err)
			BadRequest(w)
			return
		}

		e, ok := o.(*Item)
		if !ok {
			log.Printf("Cast Error: %v", o)
			BadRequest(w)
			return
		}

		e.Id = xid.New().String()
		itemRepos.Put(e)
		if err := Render(w, e); err != nil {
			Error(w, err.Error())
		}
	})

	r.Get("/user", func(w http.ResponseWriter, r *http.Request) {
		if err := Render(w, userRepos.List()); err != nil {
			NotFound(w)
		}
	})
	return r
}

var sessionRepos = NewSessionRepository()
var userRepos = NewUserRepository()
var itemRepos = NewItemRepository()
