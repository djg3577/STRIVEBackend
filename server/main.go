package server

import (
    "net/http"

    "github.com/gorilla/mux"
)

func NewServer() *mux.Router {
    router := mux.NewRouter()
    // Add your routes here
    return router
}
