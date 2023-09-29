package main

import (
	"fmt"
	"log"
	"strings"

	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	Id        int
	Email     string `email`
	Username  string `username`
	Password  string `password`
	Birthdate string `birthdate`
	Gender    string `gender`
	Firstname   string `first_name`
	Lastname  string `last_name`
	Timestamp string
}

type Post struct {
	Id        int
	Title     string
	Content   string
	Image     string
	Thread    string
	UserId    int
	Timestamp string
	Comments  []Comment
	User      User
}

type Comment struct {
	Id        int
	Content   string
	PostId    int
	UserId    int
	Timestamp string
	User      User
}

type Notification struct {
	Id       int
	Active   bool
	Object   string
	Title    string
	ObjectId int // object id (to link)

	Action    string
	Sender    string //sender username
	Recipient int    // recipient id

	Timestamp string
}

// users
// -------------------------------------------------------------------------------------

func createUsersTable(db *sql.DB) {
	users_table := `CREATE TABLE users (
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        "Username" TEXT UNIQUE,
        "Email" TEXT UNIQUE,
        "Password" TEXT,
		"Birthdate" TEXT,
        "Gender" TEXT,
        "Firstname" TEXT,
		"Lastname" TEXT,
        timestamp TEXT DEFAULT(strftime('%Y.%m.%d %H:%M', 'now')));`
	query, err := db.Prepare(users_table)
	if err != nil {
		log.Fatal(err)
	}
	query.Exec()
	fmt.Println("Table for users created successfully!")
}

func addUser(db *sql.DB, Username string, Email string, Password string, Birthdate string, Gender string, Firstname string, Lastname string) error {
	records := `INSERT INTO users(Username, Email, Password, Birthdate, Gender, Firstname, Lastname) VALUES (?, ?, ?, ?, ?, ?, ?)`
	query, err := db.Prepare(records)
	if err != nil {
		return err
	}
	_, err = query.Exec(Username, Email, Password, Birthdate, Gender, Firstname, Lastname)
	if err != nil {
		return err
	}
	return nil
}

func fetchUserByEmail(db *sql.DB, email string) User {
	var user User
	db.QueryRow("SELECT * FROM users WHERE email=?", email).Scan(&user.Id, &user.Username, &user.Email, &user.Password, &user.Birthdate, &user.Gender, &user.Firstname, &user.Lastname, &user.Timestamp)
	return user
}

func fetchUserByUsername(db *sql.DB, username string) User {
	var user User
	db.QueryRow("SELECT * FROM users WHERE username=?", username).Scan(&user.Id, &user.Username, &user.Email, &user.Password, &user.Birthdate, &user.Gender, &user.Firstname, &user.Lastname, &user.Timestamp)
	return user
}

func fetchUserById(db *sql.DB, id int) User {
	var user User
	db.QueryRow("SELECT * FROM users WHERE id=?", id).Scan(&user.Id, &user.Username, &user.Email, &user.Password, &user.Birthdate, &user.Gender, &user.Firstname, &user.Lastname, &user.Timestamp)
	return user
}

// threads (categories)
// -------------------------------------------------------------------------------------

func createThreadsTable(db *sql.DB) {
	threads_table := `CREATE TABLE threads (
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        "Subject" TEXT UNIQUE,
        timestamp DATETIME DEFAULT CURRENT_TIMESTAMP);`
	query, err := db.Prepare(threads_table)
	if err != nil {
		log.Fatal(err)
	}
	query.Exec()
	fmt.Println("Table for threads created successfully!")
}

func fetchAllThreads(db *sql.DB) []string {
	record, err := db.Query("SELECT Subject FROM threads")
	if err != nil {
		log.Fatal(err)
	}
	defer record.Close()

	var threads []string
	for record.Next() {
		var thread string
		err = record.Scan(&thread)
		if err != nil {
			log.Println(err)
		}
		threads = append(threads, thread)
	}
	return threads
}

// posts
// -------------------------------------------------------------------------------------

func createPostsTable(db *sql.DB) {
	posts_table := `CREATE TABLE posts (
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"Title" TEXT,
        "Content" TEXT,
        "Subject" TEXT,
        "User_id" INTEGER,
		"Image" TEXT,
        timestamp TEXT DEFAULT(strftime('%Y.%m.%d %H:%M', 'now')));`
	query, err := db.Prepare(posts_table)
	if err != nil {
		log.Fatal(err)
	}
	query.Exec()
	fmt.Println("Table for posts created successfully!")
}

func addPost(db *sql.DB, Title string, Image string, Content string, Subject []string, User_id int) {
	records := `INSERT INTO posts(Title, Image, Content, Subject, User_id) VALUES (?, ?, ?, ?, ?)`
	query, err := db.Prepare(records)
	if err != nil {
		log.Fatal(err)
	}

	_, err = query.Exec(Title, Image, Content, strings.Join(Subject, ", "), User_id)
	if err != nil {
		log.Fatal(err)
	}
}

// func updatePost(db *sql.DB, id int, Title string, Content string, Subject string) {
// 	db.Exec("UPDATE posts SET title = ?, content = ?, subject = ? WHERE id = ?", Title, Content, Subject, id)
// }

func fetchAllPosts(db *sql.DB) []Post {
	record, err := db.Query("SELECT * FROM posts")
	if err != nil {
		log.Fatal(err)
	}
	defer record.Close()

	var posts []Post
	for record.Next() {
		var post Post
		err = record.Scan(&post.Id, &post.Title, &post.Content, &post.Thread, &post.UserId, &post.Image, &post.Timestamp)
		if err != nil {
			log.Println(err)
		}
		posts = append(posts, post)
	}
	return posts
}

func fetchPostsByUser(db *sql.DB, user_id int) []Post {
	record, err := db.Query("SELECT * FROM posts WHERE user_id=?", user_id)
	if err != nil {
		log.Fatal(err)
	}
	defer record.Close()

	var posts []Post
	for record.Next() {
		var post Post
		err = record.Scan(&post.Id, &post.Title, &post.Content, &post.Thread, &post.UserId, &post.Image, &post.Timestamp)
		if err != nil {
			log.Println(err)
		}
		posts = append(posts, post)
	}
	return posts
}

func fetchPostByID(db *sql.DB, id int) Post {
	record, err := db.Query("SELECT * FROM posts WHERE id=?", id)
	if err != nil {
		log.Fatal(err)
	}
	defer record.Close()

	var post Post
	for record.Next() {
		err = record.Scan(&post.Id, &post.Title, &post.Content, &post.Thread, &post.UserId, &post.Image, &post.Timestamp)
		if err != nil {
			log.Println(err)
		}
	}
	return post
}

func fetchPostsByUserComments(db *sql.DB, id int) []Post {
	record, err := db.Query("SELECT p.id, p.Title, p.Content,  p.Subject, p.User_id, p.Image, p.timestamp FROM posts p INNER JOIN comments c ON c.Post_id = p.id WHERE c.User_id=?", id)
	if err != nil {
		log.Fatal(err)
	}
	defer record.Close()

	var posts []Post
	for record.Next() {
		var post Post
		err = record.Scan(&post.Id, &post.Title, &post.Content, &post.Thread, &post.UserId, &post.Image, &post.Timestamp)
		if err != nil {
			log.Println(err)
		}
		if isUnique(post, posts) {
			posts = append(posts, post)
		}
	}
	return posts
}

func updatePostByID(db *sql.DB, id int, title, filepath, content string, subject []string) error {
	_, err := db.Exec("UPDATE Posts SET title = ?, image = ?, content = ?, subject = ? WHERE id = ?", title, filepath, content, strings.Join(subject, ", "), id)
	if err != nil {
		return err
	}
	return nil
}

func updatePostImage(db *sql.DB, id int, image string) error {
	_, err := db.Exec("UPDATE Posts SET image = ? WHERE id = ?", image, id)
	if err != nil {
		return err
	}
	return nil
}

// comments
// -------------------------------------------------------------------------------------

func createCommentsTable(db *sql.DB) {
	posts_table := `CREATE TABLE comments (
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        "Content" TEXT,
        "Post_id" INTEGER,
        "User_id" INTEGER,
        timestamp TEXT DEFAULT(strftime('%Y.%m.%d %H:%M', 'now')));`
	query, err := db.Prepare(posts_table)
	if err != nil {
		log.Fatal(err)
	}
	query.Exec()
	fmt.Println("Table for comments created successfully!")
}

func fetchAllComments(db *sql.DB) []Comment {
	record, err := db.Query("SELECT * FROM comments")
	if err != nil {
		log.Fatal(err)
	}
	defer record.Close()

	var comments []Comment
	for record.Next() {
		var comment Comment
		err = record.Scan(&comment.Id, &comment.Content, &comment.PostId, &comment.UserId, &comment.Timestamp)
		if err != nil {
			log.Println(err)
		}
		comments = append(comments, comment)
	}
	return comments
}

func addComment(db *sql.DB, Content string, Post_id int, User_id int) {
	records := `INSERT INTO comments(Content, Post_id, User_id) VALUES (?, ?, ?)`
	query, err := db.Prepare(records)
	if err != nil {
		log.Fatal(err)
	}
	_, err = query.Exec(Content, Post_id, User_id)
	if err != nil {
		log.Fatal(err)
	}
}

func fetchCommentsByPost(db *sql.DB, post_id int) []Comment {
	record, err := db.Query("SELECT * FROM comments WHERE post_id=?", post_id)
	if err != nil {
		log.Fatal(err)
	}
	defer record.Close()

	var comments []Comment
	for record.Next() {
		var comment Comment
		err = record.Scan(&comment.Id, &comment.Content, &comment.PostId, &comment.UserId, &comment.Timestamp)
		if err != nil {
			log.Println(err)
		}
		comments = append(comments, comment)
	}
	return comments
}

func fetchCommentByID(db *sql.DB, id int) Comment {
	record, err := db.Query("SELECT * FROM comments WHERE id=?", id)
	if err != nil {
		log.Fatal(err)
	}
	defer record.Close()

	var comment Comment
	for record.Next() {
		err = record.Scan(&comment.Id, &comment.Content, &comment.PostId, &comment.UserId, &comment.Timestamp)
		if err != nil {
			log.Println(err)
		}
	}
	return comment
}

func updateCommentByID(db *sql.DB, id int, content string) error {
	_, err := db.Exec("UPDATE Comments SET content = ? WHERE id = ?", content, id)
	if err != nil {
		return err
	}
	return nil
}

func deleteRow(db *sql.DB, table string, id int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = ?", table)
	_, err := db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

//	notification table
//
// -------------------------------------------------------------------------------------

func createNotificationsTable(db *sql.DB) {
	n_table := `CREATE TABLE notifications (
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"Active" BOOLEAN DEFAULT 1,
        "Object" TEXT,
		"Title" TEXT,
		"Object_id" INTEGER,
        "Action" TEXT,
        "Sender" TEXT,
		"Recipient" INTEGER,
		timestamp TEXT DEFAULT(strftime('%Y.%m.%d %H:%M', 'now')));`

	query, err := db.Prepare(n_table)
	if err != nil {
		log.Fatal(err)
	}
	query.Exec()
	fmt.Println("Table for notifications created successfully!")
}

func addNotification(db *sql.DB, Object, Title string, ObjectId int, Action string, Sender string, Recipient int) {
	records := `INSERT INTO notifications(Object, Title, Object_id, Action, Sender, Recipient) VALUES (?, ?, ?, ?, ?, ?)`
	query, err := db.Prepare(records)
	if err != nil {
		log.Fatal(err)
	}

	_, err = query.Exec(Object, Title, ObjectId, Action, Sender, Recipient)
	if err != nil {
		log.Fatal(err)
	}
}

func fetchNotificationsByUserId(db *sql.DB, user_id int) []Notification {
	record, err := db.Query("SELECT * FROM notifications WHERE Recipient=?", user_id)
	if err != nil {
		log.Fatal(err)
	}
	defer record.Close()

	var all []Notification
	for record.Next() {
		var n Notification
		err = record.Scan(&n.Id, &n.Active, &n.Object, &n.Title, &n.ObjectId, &n.Action, &n.Sender, &n.Recipient, &n.Timestamp)
		if err != nil {
			log.Println(err)
		}
		all = append(all, n)
	}
	return all
}

func fetchActiveNotificationsByUserId(db *sql.DB, user_id int) []Notification {
	record, err := db.Query("SELECT * FROM notifications WHERE Recipient=? AND active=true", user_id)
	if err != nil {
		log.Fatal(err)
	}
	defer record.Close()

	var all []Notification
	for record.Next() {
		var n Notification
		err = record.Scan(&n.Id, &n.Active, &n.Object, &n.Title, &n.ObjectId, &n.Action, &n.Sender, &n.Recipient, &n.Timestamp)
		if err != nil {
			log.Println(err)
		}
		all = append(all, n)
	}
	return all
}

func disableNotificationByID(db *sql.DB, id int) error {
	_, err := db.Exec("UPDATE notifications SET active=false WHERE id = ?", id)
	if err != nil {
		return err
	}
	return nil
}
