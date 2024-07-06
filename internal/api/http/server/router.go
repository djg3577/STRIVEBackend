package server

import (
	"STRIVEBackend/internal/api/http/handlers"
	"STRIVEBackend/internal/repository"
	"STRIVEBackend/internal/service"
	"database/sql"
	"github.com/gorilla/mux"
)

func NewRouter(db *sql.DB) *mux.Router {
	// Initialize repositories
	activityRepo := &repository.ActivityRepository{DB: db}
	userRepo := &repository.UserRepository{DB: db}

	// Initialize services
	activityService := &service.ActivityService{Repo: activityRepo}
	scoreService := &service.ScoreService{Repo: activityRepo}
	userService := &service.UserService{Repo: userRepo}

	// Initialize handlers
	activityHandler := &handlers.ActivityHandler{Service: activityService}
	scoreHandler := &handlers.ScoreHandler{Service: scoreService}
	userHandler := &handlers.UserHandler{Service: userService}

	// Create the router
	router := mux.NewRouter()

	// Define routes
	router.HandleFunc("/activities", activityHandler.LogActivity).Methods("POST")
	router.HandleFunc("/daily-score", scoreHandler.GetDailyScore).Methods("GET")
	router.HandleFunc("/users", userHandler.CreateUser).Methods("POST")

	return router
}
