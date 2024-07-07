package server

import (
	"STRIVEBackend/internal/api/http/handlers"
	"STRIVEBackend/internal/repository"
	"STRIVEBackend/internal/service"
	"database/sql"

	"github.com/gorilla/mux"
)

func InitActivityRoutes(router *mux.Router, db *sql.DB) {
	activityRepo := &repository.ActivityRepository{DB: db}
	activityService := &service.ActivityService{Repo: activityRepo}
	activityHandler := &handlers.ActivityHandler{Service: activityService}

	activityRouter := router.PathPrefix("/activities").Subrouter()
	activityRouter.HandleFunc("", activityHandler.LogActivity).Methods("POST")
}
