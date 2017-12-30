package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// Define our message object
type Message struct {
	// Identifies the user
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`

	// Identifies the action for the payload
	Type string `json:"type"`

	// Identifies message and its receipients
	To      int    `json:"to"`
	Message string `json:"message"`
	ToAll   bool   `json:"to_all"`
}

type Messages []Message

var messages = Messages{}
var id int
var clients = make(map[int]*websocket.Conn)
var broadcast = make(chan Message) // broadcast channel

// Configure the upgrader
var upgrader = websocket.Upgrader{
	// ReadBufferSize:  1024,
	// WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {

	// Configure websocket route
	http.HandleFunc("/ws", handleConnections)

	// Start listening for incoming chat messages (go-routine)
	go handleMessages()

	// Start the server on localhost port 8000 and log any errors
	log.Println("http server started on :8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	id++

	// Make sure we close the connection when the function returns
	defer ws.Close()

	// Register our new client
	clients[id] = ws

	echoPreviousMessages()

	for {
		var msg Message
		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("read error: %v", err)
			// Find and delete user in messages slice
			for i, message := range messages {
				if message.ID == id {
					messages[i] = Message{}
				}
			}
			delete(clients, id)
			break
		}

		// If new connection/user attach ID then echo back to all clients so they could see this user as online
		if msg.Type == "ADD_USER" {
			msg.ID = id
		}
		log.Printf("test msg %v", msg)
		// Send the newly received message to the broadcast channel
		broadcast <- msg
	}
}

func echoPreviousMessages() {
	fmt.Println("echoing Previous Messages")
	for _, message := range messages {
		for clientID, client := range clients {

			if clientID == message.To && message.Type == "ADD_MESSAGE" {
				err := client.WriteJSON(message)
				if err != nil {
					log.Printf("close error: %v", err)
					client.Close()
					delete(clients, clientID)
				}
			} else if message.ToAll == true && message.Type == "ADD_MESSAGE" {
				err := client.WriteJSON(message)
				if err != nil {
					log.Printf("close error: %v", err)
					client.Close()
					delete(clients, clientID)
				}
			} else if message.Type == "ADD_USER" {
				err := client.WriteJSON(message)
				if err != nil {
					log.Printf("close error: %v", err)
					client.Close()
					delete(clients, clientID)
				}
			}
		}
	}
}

func handleMessages() {
	for {
		// Grab the next message from the broadcast channel
		msg := <-broadcast
		messages = append(messages, msg)

		for clientID, client := range clients {
			// Send it out to clients based on To Property
			if clientID == msg.To && msg.Type == "ADD_MESSAGE" {
				err := client.WriteJSON(msg)
				if err != nil {
					log.Printf("close error: %v", err)
					client.Close()
					delete(clients, clientID)
				}
			} else if msg.ToAll == true && msg.Type == "ADD_MESSAGE" {
				err := client.WriteJSON(msg)
				if err != nil {
					log.Printf("close error: %v", err)
					client.Close()
					delete(clients, clientID)
				}
			} else if msg.Type == "ADD_USER" {
				err := client.WriteJSON(msg)
				if err != nil {
					log.Printf("close error: %v", err)
					client.Close()
					delete(clients, clientID)
				}
			}
		}
	}
}
