package main

import (
	"database/sql"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func logRequestMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%v] %v", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	}
}

func (s *server) loggedOnly(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := r.Header.Get("username")
		password := r.Header.Get("password")

		RespErr := &jsonErr{
			Message: "authentication failure",
		}

		credStore := &User{}
		err := s.store.QueryRow(&credStore.Password, "SELECT password FROM user WHERE username='"+username+"'")
		if err != nil {
			if err == sql.ErrNoRows {
				log.Printf("Unauthorized username, err=%v", err)
				s.respond(w, r, RespErr, http.StatusUnauthorized)
				return
			}
			log.Printf("Cannot load user in DB. err=%v", err)
			s.respond(w, r, nil, http.StatusInternalServerError)
			return
		}
		if err = bcrypt.CompareHashAndPassword([]byte(credStore.Password), []byte(password)); err != nil {
			log.Printf("Unauthorized password, err=%v", err)
			s.respond(w, r, RespErr, http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
}
