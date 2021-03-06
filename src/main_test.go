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

func sendJSON(payload string, ws *websocket.Conn, t *testing.T, wsName string) {
	strBytes := []byte(payload)
	message := Message{}
	userParseErr := json.Unmarshal(strBytes, &message)
	if userParseErr != nil {
		t.Fatalf("error parsing user info JSON: %v", userParseErr)
	}

	userWriteErr := ws.WriteJSON(message)
	if userWriteErr != nil {
		t.Fatalf("error sending user info JSON: %v", userWriteErr)
	}
	t.Log("Successfully sent JSON:", message, " through websockets originating from", wsName)
}

func receiveJSON(ws *websocket.Conn, t *testing.T, wsName string) {
	message := Message{}
	readUserErr := ws.ReadJSON(&message)
	if readUserErr != nil {
		t.Fatal("error reading user info JSON:", readUserErr)
	}
	t.Log(wsName, "websocket received JSON payload: ", message)
}

func TestSendingAndReceivingJSON(t *testing.T) {

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

	ws2, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatalf("cannot establish websocket connection: %v", err)
	}
	defer ws2.Close()

	// Create and Send User Info
	userOnline1 := `{ "type": "ADD_USER", "email": "test1@gmail.com", "username": "user1" }`
	userOnline2 := `{ "type": "ADD_USER", "email": "test2@gmail.com", "username": "user2" }`
	sendJSON(userOnline1, ws, t, "conn 1")
	sendJSON(userOnline2, ws2, t, "conn 1")
	receiveJSON(ws, t, "conn 1")
	receiveJSON(ws2, t, "conn 2")
	receiveJSON(ws, t, "conn 1")
	receiveJSON(ws2, t, "conn 2")

	// readUserErr := ws.ReadJSON(&message)
	// if readUserErr != nil {
	// 	t.Fatal("error reading user info JSON:", readUserErr)
	// }
	// t.Log("received user info JSON:", message)

	// Create and Send Message
	// strMsg := `{ "type": "ADD_MESSAGE", "to": 1, "email": "mark@gmail.com", "username": "mplibunao", "message": "Hello World" }`
	// strBytes = []byte(strMsg)
	// parseErr := json.Unmarshal(strBytes, &message)
	// if parseErr != nil {
	// 	t.Fatalf("error parsing JSON: %v", parseErr)
	// }

	// writeErr := ws.WriteJSON(message)
	// if writeErr != nil {
	// 	t.Fatalf("error sending JSON: %v", writeErr)
	// }

	// Receive Message
	//receivedMessage := Message{}
	// readErr := ws.ReadJSON(&message)
	// if readErr != nil {
	// 	t.Fatalf("error reading JSON: %v", readErr)
	// }
	// t.Log("received json:", message)
	// if message.Email != "mark@gmail.com" {
	// 	t.Fatalf("incorrect Email, got: %v, want: %v", message.Email, "mark@gmail.com")
	// }
	// if message.Username != "mplibunao" {
	// 	t.Fatalf("incorrect Username, got: %v, want: %v", message.Username, "mplibunao")
	// }
	// if message.Message != "Hello World" {
	// 	t.Fatalf("incorrect Message, got: %v, want: %v", message.Message, "Hello World")
	// }

	// // Create and Send User Info
	// strUser = `{ "type": "ADD_USER", "email": "mark2@gmail.com", "username": "mplibunao2" }`
	// strBytes = []byte(strUser)
	// userParseErr = json.Unmarshal(strBytes, &message)
	// if userParseErr != nil {
	// 	t.Fatalf("error parsing user info JSON: %v", userParseErr)
	// }

	// userWriteErr = ws2.WriteJSON(message)
	// if userWriteErr != nil {
	// 	t.Fatalf("error sending user info JSON: %v", userWriteErr)
	// }

	// // Receive User Info
	// //receivedUserMsg := Message{}
	// readUserErr = ws2.ReadJSON(&message)
	// if readUserErr != nil {
	// 	t.Fatal("error reading user info JSON:", readUserErr)
	// }
	// t.Log("received user info JSON:", message)

	// // Create and Send Message
	// strMsg = `{ "type": "ADD_MESSAGE", "to": 1, "email": "mark@gmail.com", "username": "mplibunao", "message": "Hello World" }`
	// strBytes = []byte(strMsg)
	// parseErr = json.Unmarshal(strBytes, &message)
	// if parseErr != nil {
	// 	t.Fatalf("error parsing JSON: %v", parseErr)
	// }

	// writeErr = ws2.WriteJSON(message)
	// if writeErr != nil {
	// 	t.Fatalf("error sending JSON: %v", writeErr)
	// }

	// // Create and Send Message
	// strMsg = `{ "type": "ADD_MESSAGE", "to": 1, "email": "mark@gmail.com", "username": "mplibunao", "message": "Hello World" }`
	// strBytes = []byte(strMsg)
	// parseErr = json.Unmarshal(strBytes, &message)
	// if parseErr != nil {
	// 	t.Fatalf("error parsing JSON: %v", parseErr)
	// }

	// writeErr = ws2.WriteJSON(message)
	// if writeErr != nil {
	// 	t.Fatalf("error sending JSON: %v", writeErr)
	// }

	// // Receive Message
	// //receivedMessage := Message{}
	// readErr = ws2.ReadJSON(&message)
	// if readErr != nil {
	// 	t.Fatalf("error reading JSON: %v", readErr)
	// }
	// t.Log("received json:", message)
	// if message.Email != "mark@gmail.com" {
	// 	t.Fatalf("incorrect Email, got: %v, want: %v", message.Email, "mark@gmail.com")
	// }
	// if message.Username != "mplibunao" {
	// 	t.Fatalf("incorrect Username, got: %v, want: %v", message.Username, "mplibunao")
	// }
	// if message.Message != "Hello World" {
	// 	t.Fatalf("incorrect Message, got: %v, want: %v", message.Message, "Hello World")
	// }
}
