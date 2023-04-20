package main

import (
	"fmt"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type jsonUser struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
}

func (s *server) HandleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to cheddar lovers !!")
	}
}

func (s *server) handleCreateUser() http.HandlerFunc {
	type request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := request{}
		err := s.decode(w, r, &req)
		if err != nil {
			log.Printf("Cannot parse user body. err=%v", err)
			s.respond(w, r, nil, http.StatusBadRequest)
			return
		}

		hash, err := HashPassword(req.Password)
		if err != nil {
			log.Printf("Cannot generate hash. err=%v", err)
			s.respond(w, r, nil, http.StatusBadRequest)
			return
		}

		userStore := &User{
			ID:       0,
			Username: req.Username,
			Password: hash,
		}

		id, err := s.store.QueryExec(true, "INSERT INTO user (username, password) VALUES ('"+userStore.Username+"', '"+userStore.Password+"')")
		if err != nil {
			log.Printf("Cannot create user in DB. err=%v", err)
			s.respond(w, r, nil, http.StatusInternalServerError)
			return
		}
		resp := mapUserToJson(userStore, id)
		s.respond(w, r, resp, http.StatusOK)
	}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func mapUserToJson(u *User, id int64) jsonUser {
	return jsonUser{
		Id:       id,
		Username: u.Username,
	}
}
