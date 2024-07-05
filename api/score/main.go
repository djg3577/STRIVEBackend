package main

import (
    "log"
    "net/http"
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Score service is running"))
    })

    log.Println("Starting Score Service on port 8083")
    log.Fatal(http.ListenAndServe(":8083", nil))
}
