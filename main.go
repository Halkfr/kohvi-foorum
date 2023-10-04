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
	mux.HandleFunc("/sign-up/data", signup)	
	mux.HandleFunc("/sign-in/data", signin)
	mux.HandleFunc("/user", users)

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
	createUsersTable(database)
	createThreadsTable(database)
	createPostsTable(database)
	createCommentsTable(database)
	createNotificationsTable(database)

	p, _ := hashPassword("BeanBrewer123")
	addUser(database, "CoffeeExplorer", "CoffeeExplorer@gmail.com", p, "1985-06-12", "female", "Fiona", "Anderson")
	p, _ = hashPassword("EspressoMagic!")
	addUser(database, "LatteArtista", "LatteArtista@yahoo.com", p, "1990-09-05", "male", "Elio", "Rodriguez")
	p, _ = hashPassword("LatteLove#2023")
	addUser(database, "EspressoEnthusiast", "EspressoEnthusiast@outlook.com", p, "1988-04-28", "female", "Bianca", "Martelli")
	p, _ = hashPassword("CaffeineCrusader")
	addUser(database, "BeanToCupMaster", "BeanToCupMaster@coffeeforum.com", p, "1982-11-20", "male", "Maxwell", "Davidson")
	p, _ = hashPassword("CaffeineCrusader")
	addUser(database, "CaffeineConnoisseur", "CaffeineConnoisseur@hotmail.com", p, "1993-07-03", "female", "Olivia", "Chang")
	p, _ = hashPassword("avaJunkie$")
	addUser(database, "BrewingWizard", "BrewingWizard@coffeeaddict.net", p, "1989-03-15", "male", "Milo", "Williams")
	p, _ = hashPassword("BrewersBliss42")
	addUser(database, "MugCollector", "MugCollector@brewersguild.org", p, "1997-12-08", "female", "Sophie", "Mitchell")
	p, _ = hashPassword("CuppaChampion!")
	addUser(database, "AromaAficionado", "AromaAficionado@coffeelovers.com", p, "1991-10-02", "male", "Nico", "Santoro")
	p, _ = hashPassword("AromaAdventures")
	addUser(database, "CaffeineChronicles", "CaffeineChronicles@beanbuzzers.net", p, "1984-02-19", "female", "Lena", "Petrova")
	p, _ = hashPassword("BaristaBond007")
	addUser(database, "BaristaBuddy", "BaristaBuddy@coffeeclubhouse.com", p, "1995-05-07", "male", "Oscar", "Nguyen")


	// addThread(database, "Ranch", 1)
	// addThread(database, "Dogs", 1)
	// addThread(database, "Other", 1)

	// addPost(database, title1, image1, post1, threads1, 2, 2, 1)
	// addPost(database, title2, image2, post2, threads2, 2, 3, 2)
	// addPost(database, title3, image3, post3, threads3, 3, 2, 4)
	// addPost(database, title4, image4, post4, threads4, 4, 7, 1)

	// addComment(database, comment1_1, 1, 3, 1, 0)
	// addComment(database, comment1_2, 1, 4, 2, 0)
	// addComment(database, comment1_3, 1, 1, 0, 0)
	// addComment(database, comment2_1, 2, 3, 2, 1)
	// addComment(database, comment2_2, 2, 4, 0, 0)
	// addComment(database, comment3_1, 3, 4, 3, 0)
	// addComment(database, comment4_1, 4, 2, 0, 1)
	// addComment(database, comment4_2, 4, 3, 2, 2)

}
