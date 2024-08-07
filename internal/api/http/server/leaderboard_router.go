package server

import (
	"STRIVEBackend/internal/api/http/handlers"
	"STRIVEBackend/internal/repository"
	"STRIVEBackend/internal/service"
	"database/sql"
	"github.com/gin-gonic/gin"
)

func InitWebLeaderboardRoutes(api *gin.RouterGroup, db *sql.DB) {
	leaderboardRepo := &repository.LeaderBoardRepository{DB: db}
	leaderboardService := &service.LeaderboardService{Repo: leaderboardRepo}
	leaderboardHandler := &handlers.LeaderboardHandler{Service: leaderboardService}

	api.GET("/ws", leaderboardHandler.HandleWebSocket)

	leaderboardHandler.InitWebSocketHandler()
}
