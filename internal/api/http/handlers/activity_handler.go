package handlers

import (
	"STRIVEBackend/internal/service"
	"STRIVEBackend/pkg/models"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type ActivityHandler struct {
	Service *service.ActivityService
}

func (h *ActivityHandler) LogActivity(c *gin.Context) {
	var activity models.Activity
	if err := c.BindJSON(&activity); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	parsedDate, err := time.Parse("2006-01-02", activity.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}
	activity.Date = parsedDate.Format("2006-01-02")
	missingFields := []string{}
	githubUserInterface, exists := c.Get("githubUser")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not authenticated", "githubUser": githubUserInterface, "exists": exists})
		return
	}

	githubUser, ok := githubUserInterface.(*models.GitHubUser)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user data", "githubUser": githubUserInterface, "exists": exists})
		return
	}

	internalUserId, err := h.Service.GetOrCreateUserIdFromGithub(githubUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error in getting user id from github id"})
		return
	}

	activity.UserID = internalUserId

	if activity.UserID == 0 {
		missingFields = append(missingFields, "UserID")
	}
	if activity.ActivityName == "" {
		missingFields = append(missingFields, "ActivityName")
	}
	activity.ActivityName = strings.ToUpper(activity.ActivityName)
	if activity.Duration == 0 {
		missingFields = append(missingFields, "Duration")
	}

	if len(missingFields) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Missing required fields: %s", strings.Join(missingFields, ", "))})
		return
	}

	if err := h.Service.LogActivity(&activity); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Activity logged successfully"})
}

func (h *ActivityHandler) GetActivityTotals(c *gin.Context) {
	userID, _ := c.Get("userID")

	activity_totals, err := h.Service.GetActivityTotals(userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error: ERROR IN GETTING ACTIVITY TOTALS": err.Error()})
		return
	}

	c.JSON(http.StatusOK, activity_totals)
}

func (h *ActivityHandler) GetActivityDates(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not authenticated"})
		return
	}

	activity_dates, err := h.Service.GetActivityDates(userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error in getting activity dates:": err.Error()})
		return
	}

	c.JSON(http.StatusOK, activity_dates)
}
