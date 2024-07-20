package server

import (
	"STRIVEBackend/internal/api/http/handlers"
	"STRIVEBackend/internal/repository"
	"STRIVEBackend/internal/service"
	"database/sql"

	"github.com/gin-gonic/gin"
)

func InitActivityRoutes(router *gin.Engine, db *sql.DB) {
	activityRepo := &repository.ActivityRepository{DB: db}
	activityService := &service.ActivityService{Repo: activityRepo}
	activityHandler := &handlers.ActivityHandler{Service: activityService}
	authRepo := &repository.AuthRepository{DB: db}
	authService := &service.AuthService{Repo: authRepo}
	authHandler := &handlers.AuthHandler{Service: authService}

	activityGroup := router.Group("/activities")
	{
		activityGroup.POST("", authHandler.GitHubAuthMiddleware(), activityHandler.LogActivity)
		activityGroup.GET("", authHandler.GitHubAuthMiddleware(), activityHandler.GetActivityTotals)
		activityGroup.GET("/dates", authHandler.GitHubAuthMiddleware(), activityHandler.GetActivityDates)
	}
}
