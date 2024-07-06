package server

import (
	"STRIVEBackend/internal/api/http/handlers"
	"STRIVEBackend/internal/repository"
	"STRIVEBackend/internal/service"
	"database/sql"

	"github.com/gorilla/mux"
)

func InitScoreRoutes(router *mux.Router, db *sql.DB){
	activityRepo := &repository.ActivityRepository{DB: db}
	scoreService := &service.ScoreService{Repo: activityRepo}
	scoreHandler := &handlers.ScoreHandler{Service: scoreService}

	scoreRouter := router.PathPrefix("/scores").Subrouter()
	scoreRouter.HandleFunc("/daily", scoreHandler.GetDailyScore).Methods("GET")
	
}