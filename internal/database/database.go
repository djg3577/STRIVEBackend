package database

import (
	"STRIVEBackend/internal/config"
	"database/sql"
	"log"
)

type DatabaseConnector interface {
	Connect(cfg *config.Config) (*sql.DB, error)
}

type PostgresConnector struct{}

func (p *PostgresConnector) Connect(cfg *config.Config) (*sql.DB, error) {
	connStr := "user=" + cfg.DBUser + " password=" + cfg.DBPassword + " dbname=" + cfg.DBName + " host=" + cfg.DBHost + " port=" + cfg.DBPort + " sslmode=disable"
	return sql.Open("postgres", connStr)
}

func SetupDatabase(cfg *config.Config, connector DatabaseConnector) *sql.DB {
	db, err := connector.Connect(cfg)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
