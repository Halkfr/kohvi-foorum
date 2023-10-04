package main

import (
	"golang.org/x/crypto/bcrypt"
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