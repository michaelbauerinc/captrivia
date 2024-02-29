// server/server.go

package server

import (
    "log"
    "net/http"
	"encoding/json"

    "github.com/gorilla/websocket"
    "github.com/ProlificLabs/captrivia/game" // Adjust the import path as necessary
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

func HandleConnections(w http.ResponseWriter, r *http.Request) {
    ws, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Fatal(err)
    }
    defer ws.Close()

    var playerName string
    _, playerNameBytes, err := ws.ReadMessage()
    if err != nil {
        log.Printf("Error getting player name: %v", err)
        return
    }

    var nameJson struct {
        PlayerName string `json:"playerName"`
    }
    if err := json.Unmarshal(playerNameBytes, &nameJson); err == nil && nameJson.PlayerName != "" {
        playerName = nameJson.PlayerName
    } else {
        playerName = string(playerNameBytes)
    }

    player := game.NewPlayer(playerName, ws)
    game.AddPlayer(player)

    log.Printf("Player %s connected.", playerName)
    game.BroadcastRooms() // Adjusted to not require passing global state

    // Cleanup on disconnect
    defer func() {
        game.RemovePlayerFromAllRooms(player, "")
        game.BroadcastRooms() // Optionally, update room info after player disconnects
        game.RemovePlayer(player)
    }()

    for {
        _, msg, err := ws.ReadMessage()
        if err != nil {
            log.Printf("Error: %v", err)
            break
        }

        game.HandleAction(player, msg)
    }
}

// StartServer starts the HTTP server
func StartServer() {
    http.HandleFunc("/ws", HandleConnections)

    log.Println("HTTP server started on :8080")
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}
