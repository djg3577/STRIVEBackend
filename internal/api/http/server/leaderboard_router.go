package server

import (
	"STRIVEBackend/internal/api/http/handlers"
	"STRIVEBackend/internal/repository"
	"STRIVEBackend/internal/service"
	"database/sql"
	"github.com/gin-gonic/gin"
)
func init(){
	RegisterRoutes(LeaderBoardRoutes)
}

func LeaderBoardRoutes(api *gin.RouterGroup, db *sql.DB) {
	leaderBoardRepo := &repository.LeaderBoardRepository{DB: db}
	leaderBoardService := &service.LeaderBoardService{Repo: leaderBoardRepo}
	leaderBoardHandler := &handlers.LeaderBoardHandler{Service: leaderBoardService}

	api.GET("/ws", leaderBoardHandler.HandleWebSocket)

	leaderBoardHandler.InitWebSocketHandler()
}
