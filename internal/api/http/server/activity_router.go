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

	activityGroup := router.Group("/activities")
	{
		activityGroup.POST("", activityHandler.LogActivity)

	}
}
