package main

import (
	"fmt"
	"log"
	"strings"

	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	Id            int
	Email         string `email`
	Username      string `username`
	Password      string `password`
	Birthdate     string `birthdate`
	Gender        string `gender`
	Firstname     string `firstName`
	Lastname      string `lastName`
	SessionStatus string `sessionStatus`

	Timestamp string
}

type Post struct {
	Id        int
	Title     string `title`
	Content   string `content`
	Image     string `image`
	Thread    string `thread`
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
}

type Messages struct {
	Id          int
	SenderId    int
	RecipientId int
	Content     string
	Timestamp   string
}

type Notification struct {
	Id          int
	SenderId    int
	RecipientId int
	Count       int

	Timestamp string
}

// users
// -------------------------------------------------------------------------------------

func createMessagesTable(db *sql.DB) {
	message_table := `CREATE TABLE messages (
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        "SenderId" INTEGER,
        "RecipientId" INTEGER,
		"Content" TEXT,
        timestamp TEXT DEFAULT(strftime('%Y.%m.%d %H:%M', 'now')));`
	query, err := db.Prepare(message_table)
	if err != nil {
		log.Fatal(err)
	}
	query.Exec()
	fmt.Println("Table for messages created successfully!")
}

func addMessage(db *sql.DB, SenderId, RecipientId int, Content string) error {
	records := `INSERT INTO messages (SenderId, RecipientId, Content) VALUES (?, ?, ?)`
	query, err := db.Prepare(records)
	if err != nil {
		log.Fatal(err)
	}
	_, err = query.Exec(SenderId, RecipientId, Content)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func fetchChatMessages(db *sql.DB, SenderId, RecipientId, limit, offset int) []Messages {
	var allMessages []Messages

	query := `
	SELECT * FROM messages
	WHERE (SenderId, RecipientId) IN ((?, ?), (?, ?))
	ORDER BY id DESC
	LIMIT ? OFFSET ?;`

	record, err := db.Query(query, SenderId, RecipientId, RecipientId, SenderId, limit, offset)
	if err != nil {
		log.Fatal(err)
	}
	defer record.Close()
	for record.Next() {
		message := Messages{}
		err := record.Scan(&message.Id, &message.SenderId, &message.RecipientId, &message.Content, &message.Timestamp)

		if err != nil {
			log.Fatal(err)
		}
		allMessages = append(allMessages, message)
	}
	return allMessages
}

func fetchChatLastMessage(db *sql.DB, SenderId, RecipientId int) Messages {
	query := "SELECT * FROM messages WHERE (SenderId, RecipientId) IN ((?, ?), (?, ?)) ORDER BY id DESC LIMIT 1"

	row := db.QueryRow(query, SenderId, RecipientId, RecipientId, SenderId)

	message := Messages{}
	err := row.Scan(&message.Id, &message.SenderId, &message.RecipientId, &message.Content, &message.Timestamp)
	if err != nil {
		return Messages{}
	}
	return message
}

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
		"SessionStatus" TEXT DEFAULT 'Offline',
        timestamp TEXT DEFAULT(strftime('%Y.%m.%d %H:%M', 'now')));`
	query, err := db.Prepare(users_table)
	if err != nil {
		log.Fatal(err)
	}
	query.Exec()
	fmt.Println("Table for users created successfully!")
}

func addUser(db *sql.DB, Username, Email, Password, Birthdate, Gender, Firstname, Lastname, SessionStatus string) error {
	records := `INSERT INTO users(Username, Email, Password, Birthdate, Gender, Firstname, Lastname, SessionStatus) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	query, err := db.Prepare(records)
	if err != nil {
		return err
	}
	_, err = query.Exec(Username, Email, Password, Birthdate, Gender, Firstname, Lastname, SessionStatus)
	if err != nil {
		return err
	}
	return nil
}

func fetchUserByEmail(db *sql.DB, email string) User {
	var user User
	db.QueryRow("SELECT * FROM users WHERE email=?", email).Scan(&user.Id, &user.Username, &user.Email, &user.Password, &user.Birthdate, &user.Gender, &user.Firstname, &user.Lastname, &user.SessionStatus, &user.Timestamp)
	return user
}

func fetchUserByUsername(db *sql.DB, username string) User {
	var user User
	db.QueryRow("SELECT * FROM users WHERE username=?", username).Scan(&user.Id, &user.Username, &user.Email, &user.Password, &user.Birthdate, &user.Gender, &user.Firstname, &user.Lastname, &user.SessionStatus, &user.Timestamp)
	return user
}

func fetchUserById(db *sql.DB, id int) User {
	var user User
	db.QueryRow("SELECT * FROM users WHERE id=?", id).Scan(&user.Id, &user.Username, &user.Email, &user.Password, &user.Birthdate, &user.Gender, &user.Firstname, &user.Lastname, &user.SessionStatus, &user.Timestamp)
	return user
}

func fetchIdByUsername(db *sql.DB, username string) int {
	var id int
	db.QueryRow("SELECT id FROM users WHERE username=?", username).Scan(&id)
	return id
}

func fetchUserlistOffsetExclude(db *sql.DB, signedInUserId, limit, offset int) []User {
	var allUsers []User

	query := `
    SELECT u.*, MAX(m.id) AS max_message_id
    FROM users u
    LEFT JOIN (
      SELECT DISTINCT SenderId AS user_id, MAX(id) AS id FROM messages WHERE RecipientId = ? 
      GROUP BY SenderId
      UNION ALL
      SELECT DISTINCT RecipientId AS user_id, MAX(id) AS id FROM messages WHERE SenderId = ? 
      GROUP BY RecipientId
    ) m
    ON u.id = m.user_id
    WHERE u.id != ? AND m.id IS NOT NULL
    GROUP BY u.id
    UNION ALL
    SELECT u.*, NULL AS max_message_id
    FROM users u
    WHERE u.id != ?
    AND u.id NOT IN (
      SELECT DISTINCT user_id FROM (
        SELECT DISTINCT SenderId AS user_id FROM messages WHERE RecipientId = ? 
        UNION ALL
        SELECT DISTINCT RecipientId AS user_id FROM messages WHERE SenderId = ? 
      ) temp
    )
    ORDER BY max_message_id DESC NULLS LAST, u.Username COLLATE NOCASE ASC
    LIMIT ? OFFSET ?
    `

	record, err := db.Query(query, signedInUserId, signedInUserId, signedInUserId, signedInUserId, signedInUserId, signedInUserId, limit, offset)
	if err != nil {
		log.Fatal(err)
	}
	defer record.Close()
	for record.Next() {
		var user User
		var maxMessageID sql.NullInt64

		err := record.Scan(&user.Id, &user.Username, &user.Email, &user.Password, &user.Birthdate, &user.Gender, &user.Firstname, &user.Lastname, &user.SessionStatus, &user.Timestamp, &maxMessageID)

		if err != nil {
			log.Fatal(err)
		}

		allUsers = append(allUsers, user)
	}
	return allUsers
}

func updateUserStatusById(db *sql.DB, newStatus string, id int) error {
	_, err := db.Exec("UPDATE users SET SessionStatus=? WHERE id=?", newStatus, id)
	return err
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

func fetchAllPostsOffset(db *sql.DB, limit, offset int) []Post {
	record, err := db.Query("SELECT * FROM posts ORDER BY id DESC LIMIT ? OFFSET ?", limit, offset)
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

func fetchThreadPostsOffset(db *sql.DB, limit, offset int, thread string) []Post {
	record, err := db.Query("SELECT * FROM posts WHERE Subject = ? ORDER BY id DESC LIMIT ? OFFSET ?", thread, limit, offset)
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

func getNumberOfCommentsByPost(db *sql.DB, id int) (int, error) {
	var commentsNumber int
    err := database.QueryRow("SELECT COUNT(*) FROM comments WHERE post_id=?", id).Scan(&commentsNumber)

	return commentsNumber, err
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

//	notification table
//
// -------------------------------------------------------------------------------------

func createNotificationsTable(db *sql.DB) {
	n_table := `CREATE TABLE notifications (
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        "SenderId" INTEGER,
		"RecipientId" INTEGER,
		"Count" INTEGER DEFAULT 0,
		timestamp TEXT DEFAULT(strftime('%Y.%m.%d %H:%M', 'now')));`

	query, err := db.Prepare(n_table)
	if err != nil {
		log.Fatal(err)
	}
	query.Exec()
	fmt.Println("Table for notifications created successfully!")
}

func addNotification(db *sql.DB, signedInUserId, RecipientId int) {
	records := `INSERT INTO notifications(SenderId, RecipientId) VALUES (?, ?)`
	query, err := db.Prepare(records)
	if err != nil {
		log.Fatal(err)
	}

	_, err = query.Exec(signedInUserId, RecipientId)
	if err != nil {
		log.Fatal(err)
	}
}

func fetchNotications(db *sql.DB, signedInUserId, RecipientId int) Notification {
	record, err := db.Query("SELECT * FROM notifications WHERE SenderId=? AND RecipientId=?", signedInUserId, RecipientId)
	if err != nil {
		log.Fatal(err)
	}
	defer record.Close()

	var n Notification
	for record.Next() {
		err = record.Scan(&n.Id, &n.SenderId, &n.RecipientId, &n.Count, &n.Timestamp)
		if err != nil {
			log.Println(err)
		}
	}
	return n
}

func incrementNotification(db *sql.DB, signedInUserId, RecipientId int) { // adds notifications to signedInUserId by RecipientId
	_, err := db.Exec("UPDATE notifications SET Count = Count + 1 WHERE SenderId = ? AND RecipientId = ?", signedInUserId, RecipientId)
	if err != nil {
		log.Fatal(err)
	}
}

func clearNotification(db *sql.DB, signedInUserId, RecipientId int) {
	_, err := db.Exec("UPDATE notifications SET Count = 0 WHERE SenderId = ? AND RecipientId = ?", signedInUserId, RecipientId)
	if err != nil {
		log.Fatal(err)
	}
}

func fetchAllUserNotifications(db *sql.DB, signedInUserId int) int {
	var sum int
	err := db.QueryRow("SELECT SUM(Count) FROM notifications WHERE SenderId = ?", signedInUserId).Scan(&sum)
	if err != nil {
		return 0
	}
	return sum
}

func getRowCount(db *sql.DB, tableName string) (int, error) {
	var count int
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName)

	err := db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func getThreadRowCount(db *sql.DB, tableName, thread string) (int, error) {
	var count int
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE Subject = '%s'", tableName, thread)

	err := db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
