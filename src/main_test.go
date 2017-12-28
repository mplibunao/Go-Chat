package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
)

// var testClients = make(map[*websocket.Conn]bool)
// var testBroadcast = make(chan Message)

var upgrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// func testHandleConnections(w http.ResponseWriter, r *http.Request) {
// 	// Upgrade initial GET request to a websocket
// 	ws, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Make sure we close the connection when the function returns
// 	defer ws.Close()

// 	// Register our new client
// 	testClients[ws] = true

// 	for {
// 		var msg Message
// 		// Read in a new message as JSON and map it to a Message object
// 		err := ws.ReadJSON(&msg)
// 		if err != nil {
// 			log.Printf("error: %v", err)
// 			delete(testClients, ws)
// 			break
// 		}
// 		// Send the newly received message to the broadcast channel
// 		testBroadcast <- msg
// 	}
// }

func TestExample(t *testing.T) {
	// Create server with the handleConnections handler
	s := httptest.NewServer(http.HandlerFunc(handleConnections))
	defer s.Close()

	// Convert http://127.0.0.1 to ws://127.0.0.1
	u := "ws" + strings.TrimPrefix(s.URL+"/ws", "http")
	t.Log("url is:", u)

	// Connect to the server
	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer ws.Close()

	for i := 0; i < 10; i++ {

		strMsg := `{"email": "mark@gmail.com", "username": "mplibunao", "message": "Test Message" }`
		textBytes := []byte(strMsg)
		msg := Message{}
		err := json.Unmarshal(textBytes, &msg)
		if err != nil {
			t.Fatalf("error parsing json %v", err)
		}

		ws.WriteJSON(msg)

		// for client := range clients {
		// 	err := client.WriteJSON(msg)

		// }
	}

	// for i := 0; i < 10; i++ {
	// 	if err := ws.WriteMessage(websocket.TextMessage, []byte("hello")); err != nil {
	// 		t.Fatalf("%v", err)
	// 	}

	// 	_, p, err := ws.ReadMessage()
	// 	if err != nil {
	// 		t.Fatalf("%v", err)
	// 	}

	// 	if string(p) != "hello" {
	// 		t.Fatalf("bad message ")
	// 	}

	// 	t.Log("Message received:", string(p))
	// }
}

// func echo(w http.ResponseWriter, r *http.Request) {
// 	ws, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		return
// 	}

// 	defer ws.Close()
// 	for {
// 		mt, message, err := ws.ReadMessage()
// 		if err != nil {
// 			break
// 		}

// 		err = ws.WriteMessage(mt, message)
// 		if err != nil {
// 			break
// 		}
// 	}
// }
