package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	setDB()

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", index)
	mux.HandleFunc("/api/sign-up", signup)
	mux.HandleFunc("/api/sign-in", signin)
	mux.HandleFunc("/api/sign-out", signout)
	mux.HandleFunc("/api/session-status", sessionStatus)
	mux.HandleFunc("/api/users", userlist)
	mux.HandleFunc("/api/user", user)
	mux.HandleFunc("/api/posts", posts)
	mux.HandleFunc("/api/post", post)
	mux.HandleFunc("/api/comments", comments)
	mux.HandleFunc("/api/add-post", addNewPost)
	mux.HandleFunc("/api/load-chat", loadChat)
	mux.HandleFunc("/api/send-message", sendMessage)
	mux.HandleFunc("/api/username", getUsername)
	mux.HandleFunc("/api/post-creation-date", getPostCreationDate)

	mux.HandleFunc("/ws", websocketHandler)

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

func index(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("static/templates/index.html")
	tmpl.Execute(w, nil)
}

var database *sql.DB

func setDB() {

	file, err := os.Create("database.db")
	if err != nil {
		log.Fatal(err)
	}
	file.Close()
	database, _ = sql.Open("sqlite3", "database.db")
	createMessagesTable(database)
	createUsersTable(database)
	createThreadsTable(database)
	createPostsTable(database)
	createCommentsTable(database)
	createNotificationsTable(database)

	// p, _ := hashPassword("BeanBrewer123")
	addUser(database, "CoffeeExplorer", "CoffeeExplorer@gmail.com", "$2a$14$xc7w8oWZtn..60gTT7zFr.wA6l7SiuRLhhHo0WD/fUC4XyLxYVfGy", "1985-06-12", "female", "Fiona", "Anderson", "Offline")
	// p, _ = hashPassword("EspressoMagic!")
	addUser(database, "LatteArtista", "LatteArtista@yahoo.com", "$2a$14$GpHjo9xqCfo.w/Eg2dYtWuKcM8hLVZ0v5s6WMfLEljlRqaxAzj3EC", "1990-09-05", "male", "Elio", "Rodriguez", "Offline")
	// p, _ = hashPassword("LatteLove#2023")
	addUser(database, "EspressoEnthusiast", "EspressoEnthusiast@outlook.com", "$2a$14$jDn/JtRmIdokL4HUqEHsQeCxFUkGBp4/936jQ.ir5WWaC2mMz0llW", "1988-04-28", "female", "Bianca", "Martelli", "Offline")
	// p, _ = hashPassword("CaffeineCrusader")
	addUser(database, "BeanToCupMaster", "BeanToCupMaster@coffeeforum.com", "$2a$14$9FqHX88u4MlZx9R94PhQ0udtXSOQFFNkUdZHDBjEMQxAzl0XNVC.i", "1982-11-20", "male", "Maxwell", "Davidson", "Offline")
	// p, _ = hashPassword("CaffeineCrusader")
	addUser(database, "CaffeineConnoisseur", "CaffeineConnoisseur@hotmail.com", "$2a$14$DSqoxefSW.j0wN9CTzU5aewelS4zMENGpYKySQlLdY7noJTJROssK", "1993-07-03", "female", "Olivia", "Chang", "Offline")
	// p, _ = hashPassword("avaJunkie$")
	addUser(database, "BrewingWizard", "BrewingWizard@coffeeaddict.net", "$2a$14$j1jfxEMvXl2vQLwak1n0d.c.xRL1uUcxNkl1LAUIjbv.sJ/XXNNyS", "1989-03-15", "male", "Milo", "Williams", "Offline")
	// p, _ = hashPassword("BrewersBliss42")
	addUser(database, "MugCollector", "MugCollector@brewersguild.org", "$2a$14$.rzFqscwtONFL7tiqvZ2gePsxL84ziyQEsOQdK1cquk6fjlsQ9fza", "1997-12-08", "female", "Sophie", "Mitchell", "Offline")
	// p, _ = hashPassword("CuppaChampion!")
	addUser(database, "AromaAficionado", "AromaAficionado@coffeelovers.com", "$2a$14$v0l5IENYce/rWIxYPEl95./.vB1CtjlkihAPQnfZ3xl1NO7SkKfPm", "1991-10-02", "male", "Nico", "Santoro", "Offline")
	// p, _ = hashPassword("AromaAdventures")
	addUser(database, "CaffeineChronicles", "CaffeineChronicles@beanbuzzers.net", "$2a$14$2.FSLi.IhIzEFHjjiqBgeeq/taK8U9TEM044MTrf/doPUm7KoyfZW", "1984-02-19", "female", "Lena", "Petrova", "Offline")
	// p, _ = hashPassword("BaristaBond007")
	addUser(database, "BaristaBuddy", "BaristaBuddy@coffeeclubhouse.com", "$2a$14$MmJQI8GG.g.egGytFKBGB.inlTfWoDvyjHeCrRq.BNiIx8BTKfv.a", "1995-05-07", "male", "Oscar", "Nguyen", "Offline")

	addUser(database, "CoffeeLover1", "CoffeeLover1@gmail.com", "$2a$14$HashedPassword1", "1990-01-15", "female", "Emily", "Johnson", "Offline")
	addUser(database, "EspressoAddict2", "EspressoAddict2@yahoo.com", "$2a$14$HashedPassword2", "1985-07-20", "male", "Daniel", "Smith", "Offline")
	addUser(database, "BrewingPro3", "BrewingPro3@outlook.com", "$2a$14$HashedPassword3", "1989-03-05", "female", "Sophia", "Brown", "Offline")
	addUser(database, "CaffeineExpert4", "CaffeineExpert4@coffeeforum.com", "$2a$14$HashedPassword4", "1982-12-10", "male", "William", "Davis", "Offline")
	addUser(database, "MugEnthusiast5", "MugEnthusiast5@hotmail.com", "$2a$14$HashedPassword5", "1993-05-25", "female", "Ella", "Williams", "Offline")
	addUser(database, "AromaMaster6", "AromaMaster6@coffeeaddict.net", "$2a$14$HashedPassword6", "1987-09-12", "male", "Grace", "Roberts", "Offline")
	addUser(database, "CaffeineExplorer7", "CaffeineExplorer7@brewersguild.org", "$2a$14$HashedPassword7", "1995-04-08", "female", "Liam", "Martinez", "Offline")
	addUser(database, "BaristaExtraordinaire8", "BaristaExtraordinaire8@coffeelovers.com", "$2a$14$HashedPassword8", "1991-08-02", "male", "Ava", "Turner", "Offline")
	addUser(database, "BeanToCupAficionado9", "BeanToCupAficionado9@beanbuzzers.net", "$2a$14$HashedPassword9", "1986-02-15", "female", "Mason", "Harris", "Offline")
	addUser(database, "LatteLover10", "LatteLover10@coffeeclubhouse.com", "$2a$14$HashedPassword10", "1994-06-23", "male", "Olivia", "Nguyen", "Offline")

	addThread(database, "Ranch")
	addThread(database, "Dogs")
	addThread(database, "Other")
	addMessage(database, 1, 2, "Hello")
	addMessage(database, 1, 2, "How are you")
	addNotification(database, 2, 1)
	incrementNotification(database, 2, 1)
	incrementNotification(database, 2, 1)

	addMessage(database, 2, 1, "Hi, I am fine!")
	addNotification(database, 1, 2)
	incrementNotification(database, 1, 2)

	addPost(database, title1, image1, post1, threads1, 1)
	addPost(database, title2, image2, post2, threads2, 2)
	addPost(database, title3, image3, post3, threads3, 1)


	// addPost(database, title4, image4, post4, threads4, 4, 7, 1)

	addComment(database, comment1_1, 1, 3)
	addComment(database, comment1_2, 1, 4)
	addComment(database, comment1_3, 1, 2)
	// addComment(database, comment2_1, 2, 3, 2, 1)
	// addComment(database, comment2_2, 2, 4, 0, 0)
	// addComment(database, comment3_1, 3, 4, 3, 0)
	// addComment(database, comment4_1, 4, 2, 0, 1)
	// addComment(database, comment4_2, 4, 3, 2, 2)

}
