package server

import (
	"STRIVEBackend/internal/api/http/handlers"
	"github.com/gorilla/mux"
)

func NewRouter(activityHandler *handlers.ActivityHandler, ScoreHandler *handlers.ScoreHandler) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/activities", activityHandler.LogActivity).Methods("POST")
	router.HandleFunc("/daily-score", ScoreHandler.GetDailyScore).Methods("GET")

	return router
}
