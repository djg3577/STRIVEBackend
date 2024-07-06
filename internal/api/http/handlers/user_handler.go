package handlers

import (
	"STRIVEBackend/internal/service"
	"STRIVEBackend/pkg/models"
	"encoding/json"
	"fmt"
	"net/http"
)

type UserHandler struct {
	Service *service.UserService
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	ID, err := h.Service.CreateUser(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "User created with ID: %d", ID)
}
