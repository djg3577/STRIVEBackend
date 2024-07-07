package server

import (
	"database/sql"
	"github.com/gorilla/mux"
)

func NewRouter(db *sql.DB) *mux.Router {
	router := mux.NewRouter()

	// Initialize sub-routers
	InitUserRoutes(router, db)
	InitActivityRoutes(router, db)
	InitScoreRoutes(router, db)

	return router
}
