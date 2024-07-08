package server

import (
	"STRIVEBackend/internal/api/http/handlers"
	"STRIVEBackend/internal/repository"
	"STRIVEBackend/internal/service"
	"database/sql"

	"github.com/gin-gonic/gin"
)

func InitScoreRoutes(router *gin.Engine, db *sql.DB) {
	activityRepo := &repository.ActivityRepository{DB: db}
	scoreService := &service.ScoreService{Repo: activityRepo}
	scoreHandler := &handlers.ScoreHandler{Service: scoreService}

	scoreGroup := router.Group("/scores")
	{
		scoreGroup.GET("/daily", scoreHandler.GetDailyScore)
	}

}
