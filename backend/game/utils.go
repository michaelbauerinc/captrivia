package game

import (
    "encoding/json"
    "fmt"
    "github.com/gorilla/websocket"
    "io/ioutil"
    "log"
    "math/rand"
    "sync"
    "time"
)

// GameManager manages game rooms and players.
type GameManager struct {
    rooms   map[string]*Room
    roomsMu sync.Mutex

    players   map[*websocket.Conn]*Player
    playersMu sync.Mutex
}

// NewGameManager creates a new instance of GameManager.
func NewGameManager() *GameManager {
    return &GameManager{
        rooms:   make(map[string]*Room),
        players: make(map[*websocket.Conn]*Player),
    }
}

// NewPlayer creates a new Player instance.
func (gm *GameManager) NewPlayer(name string, conn *websocket.Conn) *Player {
    return &Player{
        Name: name,
        Conn: conn,
    }
}

func (gm *GameManager) loadGameQuestions() ([]GameQuestion, error) {
    fileBytes, err := ioutil.ReadFile("questions.json")
    if err != nil {
        return nil, err
    }

    var questions []GameQuestion
    if err := json.Unmarshal(fileBytes, &questions); err != nil {
        return nil, err
    }

    return questions, nil
}

func (gm *GameManager) shuffleQuestions(questions []GameQuestion, numQuestions int) []GameQuestion {
    rand.Seed(time.Now().UnixNano())
    qs := make([]GameQuestion, len(questions))
    copy(qs, questions)
    rand.Shuffle(len(qs), func(i, j int) { qs[i], qs[j] = qs[j], qs[i] })

    if numQuestions > len(qs) {
        numQuestions = len(qs)
    }

    return qs[:numQuestions]
}

func (gm *GameManager) sendMessage(conn *websocket.Conn, messageType string, data interface{}) {
    message := OutgoingMessage{Type: messageType, Data: data}
    if err := conn.WriteJSON(message); err != nil {
        log.Printf("Error sending message: %v", err)
    }
}

func (gm *GameManager) gatherScores(room *Room) map[string]int {
    scores := make(map[string]int)
    for player, score := range room.Session.Scores {
        scores[player.Name] = score
    }
    return scores
}

// BroadcastRooms sends the list of rooms to all connected players.
func (gm *GameManager) BroadcastRooms() {
    gm.roomsMu.Lock()
    defer gm.roomsMu.Unlock()

    var roomInfo []map[string]interface{}
    for roomName, room := range gm.rooms {
        players := make([]string, 0, len(room.Players))
        for player := range room.Players {
            players = append(players, player.Name)
        }
        roomData := map[string]interface{}{
            "roomName": roomName,
            "players":  players,
        }
        roomInfo = append(roomInfo, roomData)
    }

    gm.playersMu.Lock()
    defer gm.playersMu.Unlock()
    for _, player := range gm.players {
        gm.sendMessage(player.Conn, "roomsList", roomInfo)
    }
}

// AddPlayer adds a new player to the manager.
func (gm *GameManager) AddPlayer(player *Player) {
    gm.playersMu.Lock()
    defer gm.playersMu.Unlock()
    gm.players[player.Conn] = player
}

// RemovePlayer removes a player from the manager.
func (gm *GameManager) RemovePlayer(player *Player) {
    gm.playersMu.Lock()
    defer gm.playersMu.Unlock()
    delete(gm.players, player.Conn)
    // Optionally handle player removal from rooms here
}

// HandleAction processes incoming actions from players.
func (gm *GameManager) HandleAction(player *Player, msg []byte) {
    var incomingMessage IncomingMessage
    if err := json.Unmarshal(msg, &incomingMessage); err != nil {
        log.Printf("Error unmarshaling incoming message: %v", err)
        return
    }

    switch incomingMessage.Action {
    case "join":
        gm.handleJoin(player, incomingMessage.RoomName)
    case "create":
        gm.handleCreate(player, incomingMessage.RoomName)
    case "startGame":
        gm.startGame(incomingMessage.RoomName, incomingMessage.NumQuestions)
    case "submitAnswer":
        gm.handleGameInput(player, msg)
    case "leave":
        gm.RemovePlayerFromAllRooms(player, "")
        if gm.rooms[incomingMessage.RoomName] != nil {
            gm.broadcastPlayersInRoom(gm.rooms[incomingMessage.RoomName])
        }

    default:
        log.Printf("Unknown action: %s", incomingMessage.Action)
        gm.sendMessage(player.Conn, "error", fmt.Sprintf("Unknown action: %s", incomingMessage.Action))
    }
}

func (gm *GameManager) startGame(roomName string, numQuestions int) {
    gm.roomsMu.Lock()
    room, exists := gm.rooms[roomName]
    if !exists {
        gm.roomsMu.Unlock()
        return
    }
    gm.roomsMu.Unlock()

    // Countdown before starting the game
    for i := 3; i > 0; i-- {
        room.mu.Lock()
        for player := range room.Players {
            // Assuming sendMessage function can handle sending simple string messages
            gm.sendMessage(player.Conn, "countdown", fmt.Sprintf("Game starts in %d...", i))
        }
        room.mu.Unlock()
        time.Sleep(1 * time.Second) // Wait for 1 second before sending the next countdown message
    }

    // Load and shuffle questions
    questions, err := gm.loadGameQuestions()
    if err != nil {
        log.Printf("Error loading questions: %v", err)
        return
    }
    shuffledQuestions := gm.shuffleQuestions(questions, numQuestions) // Assume this function exists

    // Initialize the game session
    room.mu.Lock()
    room.Session = &GameSession{
        Questions: shuffledQuestions,
        Scores:    make(map[*Player]int),
    }
    room.mu.Unlock()

    // Send the first question to all players in the room
    gm.sendCurrentQuestion(room)
}

func (gm *GameManager) RemovePlayerFromAllRooms(player *Player, excludeRoomName string) {
    playerName := player.Name
    gm.roomsMu.Lock()
    defer gm.roomsMu.Unlock()

    for roomName, room := range gm.rooms {
        room.mu.Lock()
        // Need to iterate through players to find by name
        for p := range room.Players {
            if p.Name == playerName && roomName != excludeRoomName {
                delete(room.Players, p)
                isEmpty := len(room.Players) == 0

                if isEmpty {
                    delete(gm.rooms, roomName)
                    fmt.Printf("Room %s deleted because it is now empty.\n", roomName)
                }
                break // Found and deleted the player, can break the loop
            }
        }
        room.mu.Unlock()
    }
}

func (gm *GameManager) handleJoin(player *Player, roomName string) {
    // Ensure player is removed from any previous rooms
    gm.RemovePlayerFromAllRooms(player, roomName)

    gm.roomsMu.Lock()
    room, exists := gm.rooms[roomName]
    if !exists {
        gm.sendMessage(player.Conn, "error", "Room does not exist.")
        gm.roomsMu.Unlock()
        return
    }
    gm.roomsMu.Unlock()

    room.mu.Lock()
    room.Players[player] = true
    room.mu.Unlock()

    // Notify all players in the room, including the new player, about the current list of players
    gm.broadcastPlayersInRoom(room)

    fmt.Printf("Player %s joined room: %s\n", player.Name, roomName)
    gm.BroadcastRooms()
}

// broadcastPlayersInRoom sends the list of all players in the room to every player in that room
func (gm *GameManager) broadcastPlayersInRoom(room *Room) {
    room.mu.Lock()
    playerList := make([]string, 0, len(room.Players))
    for p := range room.Players {
        playerList = append(playerList, p.Name)
    }
    defer room.mu.Unlock()

    for p := range room.Players {
        gm.sendMessage(p.Conn, "playerListUpdate", playerList)
    }
}

// handleCreate creates a new room with the given name
func (gm *GameManager) handleCreate(player *Player, roomName string) {
    gm.roomsMu.Lock()
    _, exists := gm.rooms[roomName]
    if exists {
        gm.sendMessage(player.Conn, "error", "Room already exists.")
    } else {
        gm.rooms[roomName] = &Room{
            Players: make(map[*Player]bool),
            Session: &GameSession{
                Scores: make(map[*Player]int),
            },
        }
        gm.sendMessage(player.Conn, "created", fmt.Sprintf("Room created: %s", roomName))
    }
    gm.roomsMu.Unlock()
    gm.BroadcastRooms()
}

func (gm *GameManager) handleGameInput(player *Player, msg []byte) {
    var inputMsg GameInputMessage
    if err := json.Unmarshal(msg, &inputMsg); err != nil {
        log.Printf("Error unmarshaling game input message: %v", err)
        return
    }

    gm.roomsMu.Lock()
    room, ok := gm.rooms[inputMsg.RoomName]
    gm.roomsMu.Unlock() // Unlock immediately after fetching the room
    if !ok {
        log.Println("Room not found")
        return // Room not found
    }

    room.mu.Lock()
    if room.Session == nil || len(room.Session.Questions) == 0 {
        log.Println("Session not started or no questions available")
        room.mu.Unlock()
        return
    }

    if room.Session.CurrentQuestionIndex >= len(room.Session.Questions) {
        log.Println("No more questions left to answer")
        gm.sendGameOver(room)
        room.mu.Unlock()
        return
    }

    // Fetch the current question based on the index
    currentQuestion := room.Session.Questions[room.Session.CurrentQuestionIndex]
    correct := currentQuestion.CorrectIndex == inputMsg.AnswerIdx
    feedbackString := fmt.Sprintf("%s got the answer incorrect!", player.Name)
    if correct {
        room.Session.Scores[player] += 10
        room.Session.CurrentQuestionIndex++
        feedbackString = fmt.Sprintf("%s got the answer correct!", player.Name)
    } else {
        room.Session.Scores[player] -= 10
    }

    // Check if there are more questions or if the game has ended
    if room.Session.CurrentQuestionIndex >= len(room.Session.Questions) {
        gm.sendGameOver(room)
    } else {
        for p := range room.Players {
            fmt.Println(feedbackString)
            gm.sendFeedback(room, p, feedbackString) // Send feedback for the current question
        }
        
    }

    room.mu.Unlock() // Unlock before potentially sending the next question

    if room.Session.CurrentQuestionIndex < len(room.Session.Questions) {
        gm.sendCurrentQuestion(room) // Proceed to send the next question if the game is not over
    }
}

func (gm *GameManager) sendCurrentQuestion(room *Room) {
    room.mu.Lock()
    defer room.mu.Unlock()

    // Fetch the current question from the session
    question := room.Session.Questions[room.Session.CurrentQuestionIndex]

    // Create a ClientGameQuestion instance, omitting the CorrectIndex
    clientQuestion := ClientGameQuestion{
        ID:           question.ID,
        QuestionText: question.QuestionText,
        Options:      question.Options,
    }

    // Send the client-friendly question to all players in the room
    for player := range room.Players {
        gm.sendMessage(player.Conn, "question", clientQuestion)
    }
}

func (gm *GameManager) sendFeedback(room *Room, player *Player, correct string) {
    feedback := map[string]interface{}{
        "correct": correct,
        "scores":  gm.gatherScores(room),
    }

    gm.sendMessage(player.Conn, "answerFeedback", feedback)
}

func (gm *GameManager) sendGameOver(room *Room) {
    gameOverMessage := map[string]interface{}{
        "message": "Game over! Thanks for playing.",
        "scores":  gm.gatherScores(room),
    }
    for p := range room.Players {
        gm.sendMessage(p.Conn, "gameOver", gameOverMessage)
    }
}
