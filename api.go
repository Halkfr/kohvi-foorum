package main

import (
	"encoding/json"
	"fmt"
	"github.com/gofrs/uuid"
	"io"
	"net/http"
	"os"
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
	err = addUser(database, u.Username, u.Email, hashedPassword, u.Birthdate, u.Gender, u.Firstname, u.Lastname, u.SessionStatus)
	if err == nil {
		w.WriteHeader(http.StatusOK)
	} else {
		var usernameExist, emailExist int
		var errorMessage string
		database.QueryRow("SELECT COUNT(*) FROM users WHERE Username = ?", u.Username).Scan(&usernameExist)
		database.QueryRow("SELECT COUNT(*) FROM users WHERE Email = ?", u.Email).Scan(&emailExist)

		if usernameExist > 0 {
			errorMessage += "usernameExist"
		}
		if emailExist > 0 {
			errorMessage += "emailExist"
		}
		if errorMessage != "" {
			w.WriteHeader(http.StatusForbidden)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
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
		updateUserStatusById(database, "Offline", sessions[token])
		delete(sessions, token)
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func sessionStatus(w http.ResponseWriter, r *http.Request) {
	sessionCookie, err := r.Cookie("session_token")
	if err == nil {
		token := sessionCookie.Value
		if sessions[token] > 0 {
			fmt.Println("welcome user-id", sessions[token])
			w.WriteHeader(http.StatusOK)
		} else {
			fmt.Println("invalid token", sessions[token])
			w.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func userlist(w http.ResponseWriter, r *http.Request) {
	userCount, _ := getRowCount(database, "USERS")
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	if offset < userCount-1 {
		sessionCookie, _ := r.Cookie("session_token")
		token := sessionCookie.Value
		excludeId := sessions[token]

		w.Header().Set("Content-Type", "application/json")

		users := fetchUserlistOffsetExclude(database, excludeId, limit, offset)

		json, err := json.Marshal(users)

		if err == nil {
			w.WriteHeader(http.StatusOK)
			w.Write(json)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("{\"error\":\"Cannot marshal to json\"}"))
		}
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}

func user(w http.ResponseWriter, r *http.Request) {
	sessionCookie, err := r.Cookie("session_token")
	if err == nil {
		token := sessionCookie.Value
		user := fetchUserById(database, sessions[token])
		json, _ := json.Marshal(user)
		w.Write(json)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func posts(w http.ResponseWriter, r *http.Request) {
	var posts []Post
	var postCount int

	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	thread := r.URL.Query().Get("thread")

	if thread == "Viewall" {
		postCount, _ = getRowCount(database, "POSTS")
	} else {
		postCount, _ = getThreadRowCount(database, "POSTS", thread)
	}
	w.Header().Set("Content-Type", "application/json")
	if offset < postCount {
		if thread == "Viewall" {
			posts = fetchAllPostsOffset(database, limit, offset)
		} else {
			posts = fetchThreadPostsOffset(database, limit, offset, thread)
		}
		json, err := json.Marshal(posts)

		if err == nil {
			w.WriteHeader(http.StatusOK)
			w.Write(json)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("{\"error\":\"Cannot marshal to json\"}"))
		}
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}

func addNewPost(w http.ResponseWriter, r *http.Request) {
	sessionCookie, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, "Failed to get session cookie", http.StatusUnauthorized)
		return
	}
	senderId := sessions[sessionCookie.Value]
	imagePath := ""
	if r.FormValue("image-name") != "" {
		file, _, err := r.FormFile("image")
		if err != nil {
			http.Error(w, "Failed to get the image", http.StatusBadRequest)
			return
		}
		defer file.Close()

		u, _ := uuid.NewV4()
		imagePath = "static/images/" + u.String() + r.FormValue("image-name")

		f, err := os.Create(imagePath)
		if err != nil {
			http.Error(w, "Failed to create image file", http.StatusInternalServerError)
			return
		}
		defer f.Close()
		io.Copy(f, file)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	addPost(database, r.FormValue("title"), imagePath, r.FormValue("content"), []string{r.FormValue("thread")}, senderId)
}

func loadChat(w http.ResponseWriter, r *http.Request) {
	var returnMsgs struct {
		Messages     []Messages
		MsgUsernames []string
	}
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	sessionCookie, err := r.Cookie("session_token")
	if err == nil {
		senderId, _ := strconv.Atoi(r.URL.Query().Get("senderId"))
		recipientId := sessions[sessionCookie.Value]

		w.Header().Set("Content-Type", "application/json")
		returnMsgs.Messages = fetchChatMessages(database, senderId, recipientId, limit, offset)
		for i := 0; i < len(returnMsgs.Messages); i++ {
			returnMsgs.MsgUsernames = append(returnMsgs.MsgUsernames, fetchUserById(database, returnMsgs.Messages[i].SenderId).Username)
		}

		json, err := json.Marshal(returnMsgs)

		if err == nil {
			w.WriteHeader(http.StatusOK)
			w.Write(json)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("{\"error\":\"Cannot marshal to json\"}"))
		}
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("{\"error\":\"Cannot access cookie\"}"))
	}
}

func sendMessage(w http.ResponseWriter, r *http.Request) {
	sessionCookie, err := r.Cookie("session_token")
	if err == nil {
		recipientUsername := r.URL.Query().Get("recipient-username")

		recipientId := fetchIdByUsername(database, recipientUsername)
		senderId := sessions[sessionCookie.Value]

		fmt.Println(recipientId, senderId)
		var requestBody map[string]string

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&requestBody)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		message := requestBody["message"]

		addMessage(database, senderId, recipientId, message)
		w.Header().Set("Content-Type", "application/json")

	} else {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("{\"error\":\"Cannot access cookie\"}"))
	}
}

func getUsername(w http.ResponseWriter, r *http.Request) {
	senderId, _ := strconv.Atoi(r.URL.Query().Get("id"))
	json, err := json.Marshal(fetchUserById(database, senderId).Username)
	if err == nil {
		w.WriteHeader(http.StatusOK)
		w.Write(json)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("{\"error\":\"Cannot marshal to json\"}"))
	}
}

func getPostCreationDate(w http.ResponseWriter, r *http.Request) {
	postId, _ := strconv.Atoi(r.URL.Query().Get("id"))
	json, err := json.Marshal(fetchPostByID(database, postId).Timestamp)
	if err == nil {
		w.WriteHeader(http.StatusOK)
		w.Write(json)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("{\"error\":\"Cannot marshal to json\"}"))
	}
}
