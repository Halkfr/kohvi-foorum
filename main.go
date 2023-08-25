package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"html/template"
)

func main() {
	setDB()

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", index)

	// Create a custom server with a timeout
	server := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
	}

	fmt.Println("\nStarting server at http://127.0.0.1:8080/")
	fmt.Printf("Quit the server with CONTROL-C.\n\n")

	// Start the server
	log.Fatal(server.ListenAndServe())
}

var database *sql.DB

func setDB() {

	file, err := os.Create("database.db")
	if err != nil {
		log.Fatal(err)
	}
	file.Close()
	database, _ = sql.Open("sqlite3", "database.db")
	createUsersTable(database)
	createThreadsTable(database)
	createPostsTable(database)
	createCommentsTable(database)
	createCommentsReactionsTable(database)
	createPostsReactionsTable(database)

	p, _ := hashPassword("1234")
	addUser(database, "test", "test@gmail.com", p)
	p, _ = hashPassword("blacksheep")
	addUser(database, "Lasso-less Cowboy", "cowboy@gmail.com", p)
	p, _ = hashPassword("ZoomZoomZap")
	addUser(database, "SnapHappy", "snaphappyphotographer@email.com", p)
	p, _ = hashPassword("kingtomyqueen")
	addUser(database, "RodeoQueen", "rodeoqueen@email.com", p)

	addThread(database, "Ranch", 1)
	addThread(database, "Dogs", 1)
	addThread(database, "Other", 1)

	addPost(database, title1, image1, post1, threads1, 2, 2, 1)
	addPost(database, title2, image2, post2, threads2, 2, 3, 2)
	addPost(database, title3, image3, post3, threads3, 3, 2, 4)
	addPost(database, title4, image4, post4, threads4, 4, 7, 1)

	addComment(database, comment1_1, 1, 3, 1, 0)
	addComment(database, comment1_2, 1, 4, 2, 0)
	addComment(database, comment1_3, 1, 1, 0, 0)
	addComment(database, comment2_1, 2, 3, 2, 1)
	addComment(database, comment2_2, 2, 4, 0, 0)
	addComment(database, comment3_1, 3, 4, 3, 0)
	addComment(database, comment4_1, 4, 2, 0, 1)
	addComment(database, comment4_2, 4, 3, 2, 2)

}

func index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		// createError(w, r, http.StatusNotFound)
		return
	}

	// get data for index page
	data := welcome(w, r)

	if r.URL.Query().Get("modal") != "" {
		data.SigninModalOpen = r.URL.Query().Get("modal")
	}

	posts := fetchAllPosts(database)
	reverse(posts)

	data.Filter = r.URL.Query().Get("filter")
	if data.Filter == "" || data.Filter == "All Categories" || contains(data.Threads, data.Filter) {
		if contains(data.Threads, data.Filter) {
			posts = filterByThread(posts, data.Filter)
		}
	} else {
		// createError(w, r, http.StatusBadRequest)
		return
	}

	data.Posts = fillPosts(&data, posts)

	tmpl, err := template.ParseFiles("static/template/index.html", "static/template/base.html")
	if err != nil {
		fmt.Println(err)
		// createError(w, r, http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		// createError(w, r, http.StatusInternalServerError)
		return
	}
}
