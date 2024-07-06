package handlers

import (
	"STRIVEBackend/internal/service"
	"net/http"
	"strconv"
)

type ScoreHandler struct {
	Service *service.ScoreService
}

func (h *ScoreHandler) GetDailyScore(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	score, err := h.Service.CalculateDailyScore(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(strconv.Itoa(score)))
}
