package main

import (
	"log"
	"net/http"

	_ "github.com/lib/pq"

	"STRIVEBackend/internal/api/http/server"
	"STRIVEBackend/internal/config"
	"STRIVEBackend/internal/database"
)

func main() {
	loader := &config.EnvConfigLoader{}
	cfg := config.LoadConfig(loader)

	connector := &database.PostgresConnector{}
	db := database.SetupDatabase(cfg, connector)
	defer db.Close()

	router := server.SetupRouter(db)

	handler := server.SetupCORS(router)

	startServer(handler)
}

func startServer(handler http.Handler) {
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}