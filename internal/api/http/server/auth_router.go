package server

import (
	"STRIVEBackend/internal/api/http/handlers"
	"STRIVEBackend/internal/repository"
	"STRIVEBackend/internal/service"
	"database/sql"

	"github.com/gin-gonic/gin"
)

func InitAuthRoutes(api *gin.RouterGroup, db *sql.DB) {
	authRepo := &repository.AuthRepository{DB: db}
	authService := &service.AuthService{Repo: authRepo}
	authHandler := &handlers.AuthHandler{Service: authService}

	authGroup := api.Group("/auth")
	{
		authGroup.POST("/decode-jwt", authHandler.DecodeJWT)
		authGroup.POST("/login", authHandler.Login)
		authGroup.POST("/sign-up", authHandler.SignUp)
		authGroup.POST("/verify-email", authHandler.VerifyEmail)
		authGroup.POST("/github", authHandler.GitHubLogin)
		authGroup.GET("/testingRedis", authHandler.TestingRedis)
	}
}
