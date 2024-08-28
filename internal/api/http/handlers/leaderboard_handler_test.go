package handlers_test

import (
	"STRIVEBackend/internal/api/http/handlers"
	"STRIVEBackend/internal/repository"
	"STRIVEBackend/internal/service"
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/mock"
)

type MockLeaderBoardRepository struct {
	mock.Mock
}

func (m *MockLeaderBoardRepository) GetTopScores() ([]repository.UserScore, error) {
	args := m.Called()
	return args.Get(0).([]repository.UserScore), args.Error(1)
}

func TestLeaderBoardWebSocket(t *testing.T) {

	gin.SetMode(gin.TestMode)
	mockRepo := new(MockLeaderBoardRepository)
	leaderBoardService := &service.LeaderBoardService{
		Repo: mockRepo,
	}
	handler := &handlers.LeaderBoardHandler{
		Service: leaderBoardService,
	}

	router := gin.New()
	router.GET("/ws", handler.HandleWebSocket)
	s := httptest.NewServer(router)
	defer s.Close()

	u := "ws" + strings.TrimPrefix(s.URL, "http") + "/ws"

	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatalf("Failed to connect to WebSocket: %v", err)
	}
	defer ws.Close()

	testScores := []repository.UserScore{
		{Username: "user1", Score: 100},
		{Username: "user2", Score: 200},
	}
	mockRepo.On("GetTopScores").Return(testScores, nil)

	go handler.InitWebSocketHandler()

	time.Sleep(100 * time.Millisecond)

	_, msg, err := ws.ReadMessage()
	if err != nil {
		t.Fatalf("Failed to read message: %v", err)
	}

	var receivedScores []repository.UserScore
	err = json.Unmarshal(msg, &receivedScores)
	if err != nil {
		t.Fatalf("Failed to unmarshal message: %v", err)
	}

	if len(receivedScores) != len(testScores) {
		t.Fatalf("Expected %d scores, but got %d", len(testScores), len(receivedScores))
	}
	for i, score := range receivedScores {
		if score != testScores[i] {
			t.Fatalf("Score mismatch at index %d. Expected %v, but got %v", i, testScores[i], score)
		}
	}

	mockRepo.AssertExpectations(t)
}
