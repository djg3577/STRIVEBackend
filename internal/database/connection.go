package database

import (
    "database/sql"
    "log"

    "github.com/lib/pq" // Postgres driver
)

var db *sql.DB

func Connect() {
    var err error
    db, err = sql.Open("postgres", "user=postgres password=password dbname=strive sslmode=disable")
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }

    if err := db.Ping(); err != nil {
        log.Fatalf("Failed to ping database: %v", err)
    }
}

func GetDB() *sql.DB {
    return db
}
