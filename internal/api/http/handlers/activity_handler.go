package handlers

import (
	"STRIVEBackend/internal/service"
	"STRIVEBackend/pkg/models"
	"encoding/json"
	"net/http"
)

type ActivityHandler struct {
	Service *service.ActivityService
}

func (h *ActivityHandler) LogActivity(w http.ResponseWriter, r *http.Request) {
	var activity models.Activity
	if err := json.NewDecoder(r.Body).Decode(&activity); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.Service.LogActivity(&activity); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
