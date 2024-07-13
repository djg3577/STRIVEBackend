package server

import (
	"STRIVEBackend/internal/api/http/handlers"
	"STRIVEBackend/internal/repository"
	"STRIVEBackend/internal/service"
	"database/sql"
	"github.com/gin-gonic/gin"
)

func InitWebLeaderboardRoutes(router *gin.Engine, db *sql.DB) {
	leaderboardRepo := &repository.LeaderBoardRepository{DB: db}
	leaderboardService := &service.LeaderboardService{Repo: leaderboardRepo}
	leaderboardHandler := &handlers.LeaderboardHandler{Service: leaderboardService}

	router.GET("/ws", leaderboardHandler.HandleWebSocket)

	leaderboardHandler.InitWebSocketHandler()
}
