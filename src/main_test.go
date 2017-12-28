package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
)

var message string

func TestExample(t *testing.T) {
	// Create server with the handleConnections handler
	s := httptest.NewServer(http.HandlerFunc(handleConnections))
	go handleMessages()
	defer s.Close()

	// Convert http://127.0.0.1 to ws://127.0.0.1
	u := "ws" + strings.TrimPrefix(s.URL+"/ws", "http")

	// Connect to the server
	c, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatalf("cannot establish websocket connection: %v", err)
	}
	defer c.Close()

	for i := 0; i < 10; i++ {
		strMsg := `{ "email": "mark@gmail.com", "username": "mplibunao", "message": "Hello World" }`
		textBytes := []byte(strMsg)
		sentMessage := Message{}
		err := json.Unmarshal(textBytes, &sentMessage)
		if err != nil {
			t.Fatalf("error parsing json %v", err)
		}
		t.Log("eh")
		c.WriteJSON(sentMessage)

		// receivedMessage := Message{}
		// ws.ReadJSON(&receivedMessage)
		t.Log("woo")
		_, p, err := c.ReadMessage()
		t.Log(string(p))
		if err != nil {
			t.Fatalf("%v", err)
		}

		if string(p) != "hello" {
			t.Fatalf("bad message ")
		}
		// if string(receivedMessage) != "hello" {
		// 	t.Fatalf("bad message ")
		// }

		// receivedMessage := <-broadcast
		// t.Log("Message received %v", receivedMessage)
	}
	t.Log("woy")

	// Grab the next message from the broadcast channel
	//msg := <-broadcast
	// Send it out to every client that is currently connected
	//t.Log("message received %v", msg)

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
