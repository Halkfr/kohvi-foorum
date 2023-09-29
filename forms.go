package main

import (
	"fmt"
	"net/http"
	// "html/template"
	"encoding/json"
)

type user struct {
	FirstName  string `json:"first_name"`
	SecondName string `json:"last_name"`
	Birthdate  string `json:"birthdate"`
	Gender     string `json:"gender"`
	Nickname   string `json:"nickname"`
	Email      string `json:"email"`
	Password   string `json:"password"`
}

func signup(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var u user
	err := decoder.Decode(&u)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	// Set the status code to 200
	w.WriteHeader(http.StatusOK)

	fmt.Println(u)
}

func signin(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var u user
	err := decoder.Decode(&u)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	// Set the status code to 200
	w.WriteHeader(http.StatusOK)

	fmt.Println(u)
}