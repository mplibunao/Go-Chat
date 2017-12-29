package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var id int
var clients = make(map[int]*websocket.Conn)

// var clients = make(map[*websocket.Conn]bool) // connected clients
var broadcast = make(chan Message) // broadcast channel
//var register = make(chan *Client)

// Configure the upgrader
var upgrader = websocket.Upgrader{
	// ReadBufferSize:  1024,
	// WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Define our message object
type Message struct {
	Type     string `json:"type"`
	To       int    `json:"to"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Message  string `json:"message"`
}

// type Messages []Message

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
	//clients[ws] = true
	clients[id] = ws

	for {
		var msg Message
		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, id)
			break
		}

		msg.To = id
		log.Printf("test msg %v", msg)
		// Send the newly received message to the broadcast channel
		broadcast <- msg
	}
}

func handleMessages() {
	for {
		// Grab the next message from the broadcast channel
		msg := <-broadcast

		// Send it out to every client that is currently connected
		for clientID, client := range clients {
			log.Printf("msg %v", msg)
			log.Printf("clientID %v", clientID)
			if clientID == msg.To {
				err := client.WriteJSON(msg)
				if err != nil {
					log.Printf("error: %v", err)
					client.Close()
					delete(clients, clientID)
				}
			}
		}
	}
}
