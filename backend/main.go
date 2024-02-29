package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Question struct {
	ID           string   `json:"id"`
	QuestionText string   `json:"questionText"`
	Options      []string `json:"options"`
	CorrectIndex int      `json:"correctIndex"`
}

type PlayerSession struct {
	Score int
}

type SessionStore struct {
	sync.Mutex
	Sessions map[string]*PlayerSession
}

type ClientQuestion struct {
    ID           string   `json:"id"`
    QuestionText string   `json:"questionText"`
    Options      []string `json:"options"`
}

func (store *SessionStore) CreateSession() string {
	store.Lock()
	defer store.Unlock()

	uniqueSessionID := generateSessionID()
	store.Sessions[uniqueSessionID] = &PlayerSession{Score: 0}

	return uniqueSessionID
}

func (store *SessionStore) GetSession(sessionID string) (*PlayerSession, bool) {
	store.Lock()
	defer store.Unlock()

	session, exists := store.Sessions[sessionID]
	return session, exists
}

func generateSessionID() string {
	randBytes := make([]byte, 16)
	_, err := rand.Read(randBytes)
	if err != nil {
		log.Fatalf("Failed to generate session ID: %v", err)
	}
	return fmt.Sprintf("%x", randBytes)
}

type GameServer struct {
	QuestionsMap map[string]Question
	Sessions     *SessionStore
}

func main() {
	rand.Seed(time.Now().UnixNano()) // Initialize random seed once as a global var

	// Setup the server
	router, err := setupServer()
	if err != nil {
		log.Fatalf("Server setup failed: %v", err)
	}

	// set port to PORT or 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start the server
	log.Println("Server starting on port " + port)
	log.Fatal(router.Run(":" + port))
}

// setupServer configures and returns a new Gin instance with all routes.
// It also returns an error if there is a failure in setting up the server, e.g. loading questions.
func setupServer() (*gin.Engine, error) {
	questions, err := loadQuestions()
	if err != nil {
		return nil, err
	}

	sessions := &SessionStore{Sessions: make(map[string]*PlayerSession)}
	server := NewGameServer(questions, sessions)

	// Create Gin router and setup routes
	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	config := cors.DefaultConfig()
	// allow all origins
	config.AllowAllOrigins = true
	router.Use(cors.New(config))

	router.POST("/game/start", server.StartGameHandler)
	router.GET("/questions", server.QuestionsHandler)
	router.POST("/answer", server.AnswerHandler)
	router.POST("/game/end", server.EndGameHandler)

	return router, nil
}

func NewGameServer(questions []Question, store *SessionStore) *GameServer {
	questionsMap := make(map[string]Question)
	for _, q := range questions {
		questionsMap[q.ID] = q
	}
	return &GameServer{
		QuestionsMap: questionsMap,
		Sessions:     store,
	}
}

func (gs *GameServer) StartGameHandler(c *gin.Context) {
	sessionID := gs.Sessions.CreateSession()
	c.JSON(http.StatusOK, gin.H{"sessionId": sessionID})
}

func (gs *GameServer) QuestionsHandler(c *gin.Context) {
    shuffledQuestions := shuffleQuestions(gs.QuestionsMap)

    // Do not send correct answer to client
    clientQuestions := make([]ClientQuestion, len(shuffledQuestions[:10]))
    for i, q := range shuffledQuestions[:10] {
        clientQuestions[i] = ClientQuestion{
            ID:           q.ID,
            QuestionText: q.QuestionText,
            Options:      q.Options,
        }
    }

    c.JSON(http.StatusOK, clientQuestions)
}

func (gs *GameServer) AnswerHandler(c *gin.Context) {
	var submittedAnswer struct {
		SessionID  string `json:"sessionId"`
		QuestionID string `json:"questionId"`
		Answer     int    `json:"answer"`
	}
	if err := c.ShouldBindJSON(&submittedAnswer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	session, exists := gs.Sessions.GetSession(submittedAnswer.SessionID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}

	correct, err := gs.checkAnswer(submittedAnswer.QuestionID, submittedAnswer.Answer)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
		return
	}

	if correct {
		session.Score += 10 // Increment score for correct answer
	}

	c.JSON(http.StatusOK, gin.H{
		"correct":      correct,
		"currentScore": session.Score, // Return the current score
	})
}

func (gs *GameServer) EndGameHandler(c *gin.Context) {
	var request struct {
		SessionID string `json:"sessionId"`
	}
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	session, exists := gs.Sessions.GetSession(request.SessionID)
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"finalScore": session.Score})
}

func (gs *GameServer) checkAnswer(questionID string, submittedAnswer int) (bool, error) {
	question, exists := gs.QuestionsMap[questionID]
	if !exists {
		return false, errors.New("question not found")
	}
	return question.CorrectIndex == submittedAnswer, nil
}

func shuffleQuestions(questionsMap map[string]Question) []Question {
	qs := make([]Question, 0, len(questionsMap))
	for _, q := range questionsMap {
		qs = append(qs, q)
	}
	rand.Shuffle(len(qs), func(i, j int) { qs[i], qs[j] = qs[j], qs[i] })
	return qs
}

func loadQuestions() ([]Question, error) {
	fileBytes, err := ioutil.ReadFile("questions.json")
	if err != nil {
		return nil, err
	}

	var questions []Question
	if err := json.Unmarshal(fileBytes, &questions); err != nil {
		return nil, err
	}

	return questions, nil
}
