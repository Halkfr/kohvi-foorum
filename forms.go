package main

import (
	// "fmt"
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
	// fmt.Println("bla")
	// if r.Method == "POST" {
	// 	r.ParseForm()
	// 	fmt.Println("bla")
	// 	// logic part of log in
	// 	fmt.Println("First name:", r.Form["first-name"])
	// 	fmt.Println("Last name:", r.Form["last-name"])
	// 	fmt.Println("Date of birth:", r.Form["birthdate"])
	// 	fmt.Println("Gender:", r.Form["gender"])
	// 	fmt.Println("Nickname:", r.Form["nickname"])
	// 	fmt.Println("Email:", r.Form["email"])
	// 	fmt.Println("Password:", r.Form["password"])

	// 	http.Redirect(w, r, "/home", http.StatusSeeOther)
	// }

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
