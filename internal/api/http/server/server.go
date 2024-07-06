package server

import (
	"net/http"
)

func StartServer(router http.Handler) error {
	return http.ListenAndServe(":8080", router)
}
