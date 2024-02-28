package game

import (
    "github.com/gorilla/websocket"
    "sync"
    "net/http"
)

type GameQuestion struct {
    ID           string   `json:"id"`
    QuestionText string   `json:"questionText"`
    Options      []string `json:"options"`
    CorrectIndex int      `json:"correctIndex"`
}

type GameSession struct {
    Questions []GameQuestion
    CurrentQuestionIndex int
    Scores map[*Player]int
}

type IncomingMessage struct {
    Action       string `json:"action"`
    RoomName     string `json:"roomName"`
    NumQuestions int    `json:"numQuestions"`
}


type OutgoingMessage struct {
    Type string      `json:"type"`
    Data interface{} `json:"data"`
}

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

type Player struct {
    Name string
    Conn *websocket.Conn
}

type Room struct {
    mu          sync.Mutex
    Players     map[*Player]bool
    PlayerNames []string
    Session     *GameSession
}

type GameInputMessage struct {
    Action     string `json:"action"`
    RoomName   string `json:"roomName"`
    QuestionID string `json:"questionId"`
    AnswerIdx  int    `json:"answerIdx"`
}