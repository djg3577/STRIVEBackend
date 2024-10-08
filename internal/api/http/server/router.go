package server

import (
	"database/sql"
	"github.com/gin-gonic/gin"
)

var RouteRegistry []func(*gin.RouterGroup, *sql.DB)

func RegisterRoutes(initFunc func(*gin.RouterGroup, *sql.DB)) {
	RouteRegistry = append(RouteRegistry, initFunc)
}

func SetupRouter(db *sql.DB) *gin.Engine {
	router := gin.Default()

	api := router.Group("/api")

	InitAuthRoutes(api, db)
	InitActivityRoutes(api, db)
	InitLeaderBoardRoutes(api, db)
	InitUserRoutes(api, db)

	return router
}
