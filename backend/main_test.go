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
    testServer := httptest.NewServer(http.HandlerFunc(server.HandleConnections)) // Adjust to your actual handler
    defer testServer.Close()

    wsURL := "ws" + testServer.URL[len("http"):]
    ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
    if err != nil {
        t.Fatalf("Dial failed: %v", err)
    }
    defer ws.Close()

    playerName := "testPlayer"
    roomName := "testRoom"

    t.Run("Connect", func(t *testing.T) {
        connectPlayer(ws, playerName, t)
    })

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
		"numQuestions": 1,
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
		"answerIdx":  answerIndex,
	}
	if err := ws.WriteJSON(actionPayload); err != nil {
		t.Fatalf("Failed to submit answer: %v", err)
	}

}


func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
