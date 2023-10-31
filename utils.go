package main

import (
	"github.com/gofrs/uuid/v5"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
	"fmt"
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
	fmt.Println(password, string(bytes))
	return string(bytes), err
}

func checkPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

var sessions = map[string]int{}

func setSessionCookie(w http.ResponseWriter, user User) {
	uuid, _ := uuid.NewV4()
	sessionToken := (uuid).String()

	updateUserStatusById(database, "Online", user.Id)
	sessions[sessionToken] = user.Id

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Path:     "/",
		Expires:  time.Now().Add(15 * 60 * time.Second),
		HttpOnly: true,
	})
}
