// server/server.go

package server

import (
    "log"
    "net/http"
    "encoding/json"

    "github.com/gorilla/websocket"
    "github.com/ProlificLabs/captrivia/game"
)

type Server struct {
    GameManager *game.GameManager
    Upgrader    websocket.Upgrader
}

func NewServer() *Server {
    return &Server{
        GameManager: game.NewGameManager(),
        Upgrader: websocket.Upgrader{
            CheckOrigin: func(r *http.Request) bool {
                return true
            },
        },
    }
}

func (s *Server) HandleConnections(w http.ResponseWriter, r *http.Request) {
    ws, err := s.Upgrader.Upgrade(w, r, nil)
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

    player := s.GameManager.NewPlayer(playerName, ws)
    s.GameManager.AddPlayer(player)

    log.Printf("Player %s connected.", playerName)
    s.GameManager.BroadcastRooms()

    // Cleanup on disconnect
    defer func() {
        s.GameManager.RemovePlayerFromAllRooms(player, "")
        s.GameManager.BroadcastRooms() // Optionally, update room info after player disconnects
        s.GameManager.RemovePlayer(player)
    }()

    for {
        _, msg, err := ws.ReadMessage()
        if err != nil {
            log.Printf("Error: %v", err)
            break
        }

        s.GameManager.HandleAction(player, msg)
    }
}

func (s *Server) StartServer() {
    http.HandleFunc("/ws", s.HandleConnections)

    log.Println("HTTP server started on :8080")
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}
