package handlers

import (
	"STRIVEBackend/internal/service"
	"STRIVEBackend/pkg/models"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
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

	parsedDate, err := time.Parse("2006-01-02", activity.Date)
	if err != nil {
		http.Error(w, "Invalid date format. Use YYYY-MM-DD", http.StatusBadRequest)
		return
	}
	activity.Date = parsedDate.Format("2006-01-02")
	missingFields := []string{}

	// !! FIX THESE IF STATEMENTS THEY LOOK BAD LOL
	if activity.UserID == 0 {
		missingFields = append(missingFields, "UserID")
	}
	if activity.Type == "" {
		missingFields = append(missingFields, "Type")
	}
	if activity.Duration == 0 {
		missingFields = append(missingFields, "Duration")
	}

	if len(missingFields) > 0 {
		http.Error(w, fmt.Sprintf("Missing required fields: %s", strings.Join(missingFields, ", ")), http.StatusBadRequest)
		return
	}

	if err := h.Service.LogActivity(&activity); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
