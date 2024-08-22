package main

import (
	"STRIVEBackend/internal/api/http/scheduler"
	"STRIVEBackend/internal/api/http/server"
	"STRIVEBackend/internal/config"
	"STRIVEBackend/internal/database"
	"log"
	"net/http"
)

func main() {
	loader := &config.EnvConfigLoader{}
	cfg := config.LoadConfig(loader)

	connector := &database.PostgresConnector{}
	db := database.SetupDatabase(cfg, connector)
	defer db.Close()

	redisClient := database.SetupRedis(cfg)

	jobScheduler := scheduler.NewJobScheduler(redisClient.Client)

	jobScheduler.RunDailyJob(scheduler.TestJob)

	go jobScheduler.Start()

	router := server.SetupRouter(db)
	handler := server.SetupCORS(router)

	startServer(handler)
}

func startServer(handler http.Handler) {
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
