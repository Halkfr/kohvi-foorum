package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var (
	wsUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	clients = make(map[int]*websocket.Conn)
)

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	sessionCookie, err := r.Cookie("session_token")
	if err != nil {
		fmt.Println("Cannot access cookie")
		return
	}
	wsUpgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	wsConn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("could not upgrade: %s\n", err.Error())
	}

	defer wsConn.Close()

	type wsMessage struct {
		RecipientUsername string `recipientUsername`
		Content           string `content`
	}

	type returnMessage struct {
		Messages   Messages
		SenderName string
		Sender     bool
	}

	senderId := sessions[sessionCookie.Value]

	clients[senderId] = wsConn

	for {
		var msg wsMessage
		var newMsg returnMessage
		err := wsConn.ReadJSON(&msg)
		if err != nil {
			fmt.Println(err)
			break
		}
		recipientId := fetchIdByUsername(database, msg.RecipientUsername)
		if recipientId != 0 {
			addMessage(database, senderId, recipientId, msg.Content)

			newMsg.Messages = fetchChatLastMessage(database, senderId, recipientId)
			newMsg.Sender = true
			newMsg.SenderName = fetchUserById(database, senderId).Username

			err1 := clients[senderId].WriteJSON(newMsg)
			if err1 != nil {
				fmt.Println("Cant send to first client", err1)
				break
			}
			if clients[recipientId] != nil {
				newMsg.Sender = false
				err2 := clients[recipientId].WriteJSON(newMsg)
				if err2 != nil {
					fmt.Println("Cant send to second client", err2)
					break
				}
			}
			fmt.Println("Message sent by :", senderId, " add notification to: ", recipientId)
		}
	}
}
