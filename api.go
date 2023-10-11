package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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
	hashedPassword, err := hashPassword(u.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = addUser(database, u.Username, u.Email, hashedPassword, u.Birthdate, u.Gender, u.Firstname, u.Lastname)
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
	if checkPassword(u.Password, potentionalUser.Password) {
		setSessionCookie(w, potentionalUser)
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}
}

func signout(w http.ResponseWriter, r *http.Request) {
	sessionCookie, err := r.Cookie("session_token")
	// add error handler to check if sessionCookie is nil
	token := sessionCookie.Value
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})
	if err == nil {
		delete(sessions, token)
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func users(w http.ResponseWriter, r *http.Request) {
	first, _ := strconv.Atoi(r.URL.Query().Get("first"))
	last, _ := strconv.Atoi(r.URL.Query().Get("last"))

	w.Header().Set("Content-Type", "application/json")

	users := fetchAllUsers(database)

	/* TODO: Sort by last message & A-Z */
	batch := users[first:last]
	json, err := json.Marshal(batch)

	if err == nil {
		w.WriteHeader(http.StatusOK)
		w.Write(json)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("{\"error\":\"Cannot marshal to json\"}"))
	}
}
