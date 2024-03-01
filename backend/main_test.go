package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/ProlificLabs/captrivia/server"
)

func TestWebSocketGameFlow(t *testing.T) {
	// Initialize the new server structure
	srv := server.NewServer()
	testServer := httptest.NewServer(http.HandlerFunc(srv.HandleConnections))
	defer testServer.Close()

	wsURL := "ws" + testServer.URL[len("http"):]
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Dial failed: %v", err)
	}
	defer ws.Close()

	// Initial connection requires sending playerName as JSON
	playerName := "testPlayer"
	if err := ws.WriteJSON(map[string]interface{}{"playerName": playerName}); err != nil {
		t.Fatalf("Failed to initialize player connection: %v", err)
	}

	roomName := "testRoom"

	t.Run("CreateRoom", func(t *testing.T) {
		createRoom(ws, playerName, roomName, t)
	})

	t.Run("JoinRoom", func(t *testing.T) {
		joinRoom(ws, playerName, roomName, t)
	})

	t.Run("StartGame", func(t *testing.T) {
		startGame(ws, roomName, t)
	})

	t.Run("SubmitAnswer", func(t *testing.T) {
		submitAnswer(ws, playerName, roomName, 0, t) // Assuming 0 is a valid answer index
	})
}

func connectPlayer(ws *websocket.Conn, playerName string, t *testing.T) {
	if err := ws.WriteJSON(map[string]string{"type": "connect", "playerName": playerName}); err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
}

func createRoom(ws *websocket.Conn, playerName, roomName string, t *testing.T) {
	actionPayload := map[string]interface{}{
		"action":     "create",
		"roomName":   roomName,
		"playerName": playerName,
	}
	if err := ws.WriteJSON(actionPayload); err != nil {
		t.Fatalf("Failed to create room: %v", err)
	}
}

func joinRoom(ws *websocket.Conn, playerName, roomName string, t *testing.T) {
	actionPayload := map[string]interface{}{
		"action":     "join",
		"roomName":   roomName,
		"playerName": playerName,
	}
	if err := ws.WriteJSON(actionPayload); err != nil {
		t.Fatalf("Failed to join room: %v", err)
	}
}

func startGame(ws *websocket.Conn, roomName string, t *testing.T) {
	actionPayload := map[string]interface{}{
		"action":   "startGame",
		"roomName": roomName,
		"numQuestions": 1, // Assuming starting a game requires specifying the number of questions
	}
	if err := ws.WriteJSON(actionPayload); err != nil {
		t.Fatalf("Failed to start game: %v", err)
	}
}

func submitAnswer(ws *websocket.Conn, playerName, roomName string, answerIndex int, t *testing.T) {
	actionPayload := map[string]interface{}{
		"action":     "submitAnswer",
		"roomName":   roomName,
		"playerName": playerName,
		"answerIndex": answerIndex, // Make sure the field matches the expected payload on the server side
	}
	if err := ws.WriteJSON(actionPayload); err != nil {
		t.Fatalf("Failed to submit answer: %v", err)
	}
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
