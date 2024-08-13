package server

import (
	"STRIVEBackend/internal/api/http/handlers"
	"STRIVEBackend/internal/repository"
	"STRIVEBackend/internal/service"
	"database/sql"

	"github.com/gin-gonic/gin"
)

func InitScoreRoutes(api *gin.RouterGroup, db *sql.DB) {
	activityRepo := &repository.ActivityRepository{DB: db}
	scoreService := &service.ScoreService{Repo: activityRepo}
	scoreHandler := &handlers.ScoreHandler{Service: scoreService}

	scoreGroup := api.Group("/scores")
	{
		scoreGroup.GET("/daily", scoreHandler.GetDailyScore)
	}

}
