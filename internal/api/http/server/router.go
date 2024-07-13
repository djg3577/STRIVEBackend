package server

import (
	"database/sql"
	"github.com/gin-gonic/gin"
)

func NewRouter(db *sql.DB) *gin.Engine {
	router := gin.Default()

	// Initialize sub-routers
	InitUserRoutes(router, db)
	InitActivityRoutes(router, db)
	InitScoreRoutes(router, db)
	InitAuthRoutes(router, db)
	InitWebLeaderboardRoutes(router, db)

	return router
}
