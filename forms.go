package main

import (
	"fmt"
	"net/http"
	// "html/template"
	"encoding/json"
)

func signup(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var u User
	err := decoder.Decode(&u)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	// Set the status code to 200
	err = addUser(database, u.Username, u.Email, u.Password, u.Birthdate, u.Gender, u.Firstname, u.Lastname)
	if err == nil {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
	fmt.Println(u)
}

func signin(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var u User
	err := decoder.Decode(&u)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")

	fmt.Println(u)

	potentionalUser := fetchUserByUsername(database, u.Username)
	if potentionalUser == (User{}) {
		potentionalUser = fetchUserByEmail(database, u.Username) // actually it is email in u.Username
		if potentionalUser == (User{}) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	}
	if potentionalUser.Password == u.Password {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}
}
