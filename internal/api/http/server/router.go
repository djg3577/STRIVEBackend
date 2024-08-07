package server

import (
	"database/sql"
	"github.com/gin-gonic/gin"
)

func NewRouter(db *sql.DB) *gin.Engine {
	router := gin.Default()

	api := router.Group("/api")

	InitUserRoutes(api, db)
	InitActivityRoutes(api, db)
	InitScoreRoutes(api, db)
	InitAuthRoutes(api, db)
	InitWebLeaderboardRoutes(api, db)

	return router
}
