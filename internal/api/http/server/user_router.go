package server

import (
	"STRIVEBackend/internal/api/http/handlers"
	"STRIVEBackend/internal/repository"
	"STRIVEBackend/internal/service"
	"database/sql"

	"github.com/gorilla/mux"
)

func InitUserRoutes(router *mux.Router, db *sql.DB) {
	userRepo := &repository.UserRepository{DB: db}
	userService := &service.UserService{Repo: userRepo}
	userHandler := &handlers.UserHandler{Service: userService}

	userRouter := router.PathPrefix("/users").Subrouter()
	userRouter.HandleFunc("", userHandler.CreateUser).Methods("POST")
}
