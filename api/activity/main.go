package main

import (
    "log"
    "net/http"
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Activity service is running"))
    })

    log.Println("Starting Activity Service on port 8082")
    log.Fatal(http.ListenAndServe(":8082", nil))
}
