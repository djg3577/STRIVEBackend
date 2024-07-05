package main

import (
    "log"
    "net/http"
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("User service is running"))
    })

    log.Println("Starting User Service on port 8081")
    log.Fatal(http.ListenAndServe(":8081", nil))
}
