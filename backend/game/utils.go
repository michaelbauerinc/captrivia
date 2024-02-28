package game

import (
    "encoding/json"
    "io/ioutil"
    "log"
    "math/rand"
    "sync"
    "time"
    "fmt"
    "github.com/gorilla/websocket"
)

var (
    rooms   = make(map[string]*Room)
    roomsMu sync.Mutex

    players   = make(map[*websocket.Conn]*Player)
    playersMu sync.Mutex
)

// NewPlayer creates a new Player instance.
func NewPlayer(name string, conn *websocket.Conn) *Player {
    return &Player{
        Name: name,
        Conn: conn,
    }
}

func loadGameQuestions() ([]GameQuestion, error) {
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

func shuffleQuestions(questions []GameQuestion, numQuestions int) []GameQuestion {
    rand.Seed(time.Now().UnixNano())
    qs := make([]GameQuestion, len(questions))
    copy(qs, questions)
    rand.Shuffle(len(qs), func(i, j int) { qs[i], qs[j] = qs[j], qs[i] })

    if numQuestions > len(qs) {
        numQuestions = len(qs)
    }

    return qs[:numQuestions]
}

func sendMessage(conn *websocket.Conn, messageType string, data interface{}) {
    message := OutgoingMessage{Type: messageType, Data: data}
    if err := conn.WriteJSON(message); err != nil {
        log.Printf("Error sending message: %v", err)
    }
}

func gatherScores(room *Room) map[string]int {
    scores := make(map[string]int)
    for player, score := range room.Session.Scores {
        scores[player.Name] = score
    }
    return scores
}

func BroadcastRooms() {
    roomsMu.Lock()
    defer roomsMu.Unlock()

    var roomInfo []map[string]interface{}
    for roomName, room := range rooms {
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

    playersMu.Lock()
    defer playersMu.Unlock()
    for _, player := range players {
        sendMessage(player.Conn, "roomsList", roomInfo)
    }
}

// New utility functions to handle player and room management
func AddPlayer(player *Player) {
    playersMu.Lock()
    defer playersMu.Unlock()
    players[player.Conn] = player
}

func RemovePlayer(player *Player) {
    playersMu.Lock()
    defer playersMu.Unlock()
    delete(players, player.Conn)
    // Optionally handle player removal from rooms here
}

func HandleAction(player *Player, msg []byte) {
    var incomingMessage IncomingMessage
    if err := json.Unmarshal(msg, &incomingMessage); err != nil {
        log.Printf("Error unmarshaling incoming message: %v", err)
        return
    }

    switch incomingMessage.Action {
    case "join":
        handleJoin(player, incomingMessage.RoomName)
    case "create":
        handleCreate(player, incomingMessage.RoomName)
    case "startGame":
        startGame(incomingMessage.RoomName, incomingMessage.NumQuestions)
    case "submitAnswer":
        handleGameInput(player, msg)
    default:
        log.Printf("Unknown action: %s", incomingMessage.Action)
        sendMessage(player.Conn, "error", fmt.Sprintf("Unknown action: %s", incomingMessage.Action))
    }
}

func startGame(roomName string, numQuestions int) {
    roomsMu.Lock()
    room, exists := rooms[roomName]
    if !exists {
        roomsMu.Unlock()
        return
    }
    roomsMu.Unlock()

    // Countdown before starting the game
    for i := 3; i > 0; i-- {
        room.mu.Lock()
        for player := range room.Players {
            // Assuming sendMessage function can handle sending simple string messages
            sendMessage(player.Conn, "countdown", fmt.Sprintf("Game starts in %d...", i))
        }
        room.mu.Unlock()
        time.Sleep(1 * time.Second) // Wait for 1 second before sending the next countdown message
    }
    
    // Load and shuffle questions
    questions, err := loadGameQuestions()
    if err != nil {
        log.Printf("Error loading questions: %v", err)
        return
    }
    shuffledQuestions := shuffleQuestions(questions, numQuestions) // Assume this function exists

    // Initialize the game session
    room.mu.Lock()
    room.Session = &GameSession{
        Questions: shuffledQuestions,
        Scores:    make(map[*Player]int),
    }
    room.mu.Unlock()

    // Send the first question to all players in the room
    sendCurrentQuestion(room)
}

// Removes player from all rooms they are currently in
func RemovePlayerFromAllRooms(player *Player, excludeRoomName string) {
    roomsMu.Lock()
    defer roomsMu.Unlock()

    for roomName, room := range rooms {
        if _, ok := room.Players[player]; ok && roomName != excludeRoomName {
            room.mu.Lock()
            delete(room.Players, player)
            isEmpty := len(room.Players) == 0
            room.mu.Unlock()

            if isEmpty {
                delete(rooms, roomName)
                fmt.Printf("Room %s deleted because it is now empty.\n", roomName)
            }
        }
    }
}

func handleJoin(player *Player, roomName string) {
    // Ensure player is removed from any previous rooms
    RemovePlayerFromAllRooms(player, roomName)

    roomsMu.Lock()
    room, exists := rooms[roomName]
    if !exists {
        sendMessage(player.Conn, "error", "Room does not exist.")
        roomsMu.Unlock()
        return
    }
    roomsMu.Unlock()

    room.mu.Lock()
    room.Players[player] = true
    playerList := make([]string, 0, len(room.Players))
    for p := range room.Players {
        playerList = append(playerList, p.Name)
    }
    room.mu.Unlock()

    // Notify all players in the room, including the new player, about the current list of players
    broadcastPlayersInRoom(room, playerList)

    fmt.Printf("Player %s joined room: %s\n", player.Name, roomName)
    BroadcastRooms()
}

// broadcastPlayersInRoom sends the list of all players in the room to every player in that room
func broadcastPlayersInRoom(room *Room, playerList []string) {
    room.mu.Lock()
    defer room.mu.Unlock()

    for p := range room.Players {
        sendMessage(p.Conn, "playerListUpdate", playerList)
    }
}

// handleCreate creates a new room with the given name
func handleCreate(player *Player, roomName string) {
    roomsMu.Lock()
    _, exists := rooms[roomName]
    if exists {
        sendMessage(player.Conn, "error", "Room already exists.")
    } else {
        rooms[roomName] = &Room{
            Players: make(map[*Player]bool),
            Session: &GameSession{
                Scores: make(map[*Player]int),
            },
        }
        sendMessage(player.Conn, "created", fmt.Sprintf("Room created: %s", roomName))
    }
    roomsMu.Unlock()
    BroadcastRooms()
}

func handleGameInput(player *Player, msg []byte) {
    var inputMsg GameInputMessage
    if err := json.Unmarshal(msg, &inputMsg); err != nil {
        log.Printf("Error unmarshaling game input message: %v", err)
        return
    }

    roomsMu.Lock()
    room, ok := rooms[inputMsg.RoomName]
    roomsMu.Unlock() // Unlock immediately after fetching the room
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
        sendGameOver(room)
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
        sendGameOver(room)
    } else {
        for p := range room.Players {
            // test := string(fmt.Sprintf("%s got the answer %s!", player.Name, string(correct))
            fmt.Println(feedbackString)
            sendFeedback(room, p, feedbackString) // Send feedback for the current question
        }
        
    }

    room.mu.Unlock() // Unlock before potentially sending the next question

    if room.Session.CurrentQuestionIndex < len(room.Session.Questions) {
        sendCurrentQuestion(room) // Proceed to send the next question if the game is not over
    }
}

func sendCurrentQuestion(room *Room) {
    room.mu.Lock()
    defer room.mu.Unlock()

    // Game is not over, send the current question
    question := room.Session.Questions[room.Session.CurrentQuestionIndex]
    fmt.Println(question)
    for player := range room.Players {
        sendMessage(player.Conn, "question", question)
    }
}

func sendFeedback(room *Room, player *Player, correct string) {
    feedback := map[string]interface{}{
        "correct": correct,
        "scores":  gatherScores(room),
    }

    sendMessage(player.Conn, "answerFeedback", feedback)
}

func sendGameOver(room *Room) {
    gameOverMessage := map[string]interface{}{
        "message": "Game over! Thanks for playing.",
        "scores":  gatherScores(room), // Assuming gatherScores is a function that compiles scores
    }
    for p := range room.Players {
        sendMessage(p.Conn, "gameOver", gameOverMessage)
    }
}