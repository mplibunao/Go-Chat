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
	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatalf("cannot establish websocket connection: %v", err)
	}
	defer ws.Close()

	// Create and Send Message
	strMsg := `{ "to": 1, "email": "mark@gmail.com", "username": "mplibunao", "message": "Hello World" }`
	textBytes := []byte(strMsg)
	sentMessage := Message{}
	parseErr := json.Unmarshal(textBytes, &sentMessage)
	if parseErr != nil {
		t.Fatalf("error parsing JSON: %v", parseErr)
	}

	writeErr := ws.WriteJSON(sentMessage)
	if writeErr != nil {
		t.Fatalf("error sending JSON: %v", writeErr)
	}

	// Receive Message
	receivedMessage := Message{}
	readErr := ws.ReadJSON(&receivedMessage)
	if readErr != nil {
		t.Fatalf("error reading JSON: %v", readErr)
	}
	t.Log("received json:", receivedMessage)
	if receivedMessage.Email != "mark@gmail.com" {
		t.Fatalf("incorrect Email, got: %v, want: %v", receivedMessage.Email, "mark@gmail.com")
	}
	if receivedMessage.Username != "mplibunao" {
		t.Fatalf("incorrect Username, got: %v, want: %v", receivedMessage.Username, "mplibunao")
	}
	if receivedMessage.Message != "Hello World" {
		t.Fatalf("incorrect Message, got: %v, want: %v", receivedMessage.Message, "Hello World")
	}
}
