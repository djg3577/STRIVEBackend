package handlers

import (
	"log"
	"net/http"
	"time"

	"STRIVEBackend/internal/repository"
	"STRIVEBackend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type LeaderboardHandler struct {
	Service *service.LeaderboardService
}

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan []repository.UserScore) // Change to slice of UserScore

// Configure the upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *LeaderboardHandler) HandleWebSocket(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Failed to upgrade to WebSocket:", err)
		return
	}
	defer ws.Close()

	clients[ws] = true
	log.Println("Client connected")

	for {
		var msg repository.UserScore
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error reading JSON: %v", err)
			delete(clients, ws)
			break
		}
		log.Printf("Received message: %v", msg)
	}
	log.Println("Client disconnected")
}

func (h *LeaderboardHandler) handleMessages() {
	for {
		msg := <-broadcast
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error writing JSON: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func (h *LeaderboardHandler) fetchAndBroadcastTopScores() {
	for {
		topScores, err := h.Service.GetTopScores()
		if err != nil {
			log.Printf("error fetching top scores: %v", err)
			continue
		}
		broadcast <- topScores
		time.Sleep(10 * time.Second)
	}
}

func (h *LeaderboardHandler) InitWebSocketHandler() {
	go h.handleMessages()
	go h.fetchAndBroadcastTopScores()
}
