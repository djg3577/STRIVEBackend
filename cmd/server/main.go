package main

import (
    "log"
    "net/http"

    "github.com/djg3577/STRIVEBackend/server"
)

func main() {
    srv := server.NewServer()
    log.Println("Starting server on port 8080")
    if err := http.ListenAndServe(":8080", srv); err != nil {
        log.Fatalf("Server failed to start: %v", err)
    }
}
