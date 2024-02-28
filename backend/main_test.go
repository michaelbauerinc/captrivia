package main_test

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestGameIntegration(t *testing.T) {
	// Create a mock WebSocket server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Upgrade HTTP connection to WebSocket
		upgrader := websocket.Upgrader{}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Fatalf("Error upgrading to WebSocket: %v", err)
		}
		defer conn.Close()

		// Simulate player actions
		handleActions(conn)
	}))
	defer server.Close()

	// Connect to the mock WebSocket server
	wsURL := "ws" + server.URL[4:] // Convert http:// to ws://
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Error connecting to WebSocket server: %v", err)
	}
	defer ws.Close()

	// Simulate game flow
	handleActions(ws)

	// Wait for a few seconds to receive messages
	time.Sleep(5 * time.Second)

	// Additional assertions can be added here to verify the behavior
}

// handleActions simulates player actions during the game
func handleActions(conn *websocket.Conn) {
	// Simulate connecting
	playerName := "TestPlayer"
	err := sendMessage(conn, map[string]string{"playerName": playerName})
	if err != nil {
		log.Fatalf("Error sending player name: %v", err)
	}

	// Simulate creating a room
	roomName := "room1"
	err = sendRoomAction(conn, "create", roomName)
	if err != nil {
		log.Fatalf("Error sending room action: %v", err)
	}

	// Simulate joining the room
	err = sendRoomAction(conn, "join", roomName)
	if err != nil {
		log.Fatalf("Error sending room action: %v", err)
	}

	// Simulate starting the game
	err = sendRoomAction(conn, "startGame", roomName)
	if err != nil {
		log.Fatalf("Error sending room action: %v", err)
	}

	// Simulate answering the questions
	for i := 0; i < 5; i++ { // Assuming there are 5 questions
		err = sendAnswer(conn, roomName, i)
		if err != nil {
			log.Fatalf("Error sending answer: %v", err)
		}
	}
}

// sendAnswer simulates sending an answer to the server
func sendAnswer(conn *websocket.Conn, roomName string, answerIdx int) error {
	// Create a message with the submitted answer index
	message := map[string]string{
		"action":     "submitAnswer",
		"roomName":   roomName,
		"answerIdx":  strconv.Itoa(answerIdx), // Convert int to string
	}
	return sendMessage(conn, message)
}

// sendMessage sends a message to the WebSocket server
func sendMessage(conn *websocket.Conn, message map[string]string) error {
	msgBytes, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("error marshaling message: %v", err)
	}
	err = conn.WriteMessage(websocket.TextMessage, msgBytes)
	if err != nil {
		return fmt.Errorf("error sending message: %v", err)
	}
	return nil
}

// sendRoomAction sends a room action to the WebSocket server
func sendRoomAction(conn *websocket.Conn, action, roomName string) error {
	// Create a message with the room action and room name
	message := map[string]string{
		"action":   action,
		"roomName": roomName,
	}
	return sendMessage(conn, message)
}
