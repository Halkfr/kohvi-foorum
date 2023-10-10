package main

import (
	"github.com/gofrs/uuid/v5"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

func isUnique(p Post, posts []Post) bool {
	for _, v := range posts {
		if p.Id == v.Id {
			return false
		}
	}
	return true
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

var sessions = map[string]User{}

func setSessionCookie(w http.ResponseWriter, user User) {
	uuid, _ := uuid.NewV4()
	sessionToken := (uuid).String()

	sessions[sessionToken] = user

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Path:    "/",
		Expires: time.Now().Add(15 * 60 * time.Second),
	})
}
