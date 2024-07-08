package handlers

import (
	"STRIVEBackend/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ScoreHandler struct {
	Service *service.ScoreService
}

func (h *ScoreHandler) GetDailyScore(c *gin.Context) {
	userID, err := strconv.Atoi(c.Query("user_id"))
	if err != nil {
		c.JSON((http.StatusBadRequest), gin.H{"error": "Invalid user ID"})
		return
	}

	score, err := h.Service.CalculateDailyScore(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.String(http.StatusOK, strconv.Itoa(score))
}
