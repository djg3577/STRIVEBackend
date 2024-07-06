package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/rs/cors"

	"STRIVEBackend/internal/api/http/handlers"
	"STRIVEBackend/internal/api/http/server"
	"STRIVEBackend/internal/config"
	"STRIVEBackend/internal/repository"
	"STRIVEBackend/internal/service"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBHost, cfg.DBPort)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	activityRepo := &repository.ActivityRepository{DB: db}
	// userRepo := &repository.UserRepository{DB: db}

	activityService := &service.ActivityService{Repo: activityRepo}
	scoreService := &service.ScoreService{Repo: activityRepo}

	activityHandler := &handlers.ActivityHandler{Service: activityService}
	scoreHandler := &handlers.ScoreHandler{Service: scoreService}

	router := server.NewRouter(activityHandler, scoreHandler)

	// Setup CORS with specific allowed origins
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	// Wrap the router with the CORS middleware
	handler := c.Handler(router)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
