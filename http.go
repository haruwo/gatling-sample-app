package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func Render(w http.ResponseWriter, obj interface{}) error {
	return RenderWithStatus(w, http.StatusOK, obj)
}

func RenderWithStatus(w http.ResponseWriter, statusCode int, obj interface{}) error {
	b, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if _, err := w.Write(b); err != nil {
		return err
	} else {
		return nil
	}
}

func Error(w http.ResponseWriter, message string) {
	http.Error(w, message, http.StatusInternalServerError)
}

func NotFound(w http.ResponseWriter) {
	http.Error(w, "Not Found", http.StatusNotFound)
}

func BadRequest(w http.ResponseWriter) {
	http.Error(w, "Bad Request", http.StatusBadRequest)
}

func Unauthorized(w http.ResponseWriter) {
	http.Error(w, "Unauthorized", http.StatusUnauthorized)
}

func WithUser(r *http.Request, user *User) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), "current.user", user))
}

func Decode(r *http.Request, o interface{}) (interface{}, error) {
	if err := json.NewDecoder(r.Body).Decode(o); err != nil {
		return nil, err
	} else {
		return o, err
	}
}

func CurrentUser(r *http.Request) *User {
	return r.Context().Value("current.user").(*User)
}

type Middleware = func(http.Handler) http.Handler

func RequireAuth(repos SessionRepository) Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if header := r.Header.Get("Authorization"); header == "" {
				log.Print("Can't extract Authorization header")
				reject(w, r)
			} else if token, err := extractToken(header); err != nil {
				log.Printf("Can't extract token from header `%s`", header)
				reject(w, r)
			} else if user, err := authToken(repos, token); err != nil {
				log.Printf("Fail auth token by `%s`", token)
				reject(w, r)
			} else {
				h.ServeHTTP(w, WithUser(r, user))
			}
		})
	}
}

func reject(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Bad Request", http.StatusBadRequest)
}

func extractToken(hv string) (string, error) {
	if a := strings.SplitN(hv, " ", 2); len(a) < 2 {
		return "", fmt.Errorf("Invalid Header Value '%s'", hv)
	} else {
		return a[1], nil
	}
}

func authToken(repos SessionRepository, token string) (*User, error) {
	if session := repos.Get(token); session == nil {
		return nil, fmt.Errorf("Invalid Token Value '%s'", token)
	} else {
		return session.User, nil
	}
}
