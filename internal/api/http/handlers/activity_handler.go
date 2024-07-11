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

// TODO: IF THE ACTIVITY ALREADY EXISTS WE UPDATE THE TIME FOR THAT ACTIVITY
// ??: TRACK THE TIME A USER SPENDS ON VSCODE? make a vscode extension that logs the time spent on vscode
// ?? HOW CAN WE PREVENT USERS FROM LOGGING ACTIVITIES IN THE FUTURE? OR IN THE PAST?
// TODO: Implement a check in the app that disallows logging activities with a timestamp later than the current system time.
// !! Set a threshold for how far back users can log activities (e.g., no more than 24 hours in the past).
// ?? HOW CAN WE PREVENT ADDING OF HUGE AMOUNTS OF TIME that is not possible?
// !! Set reasonable limits on the maximum amount of time that can be logged for an activity per day (e.g., a maximum of 16 hours).
// !! Validate the time input against these limits before saving the activity.
// ?? HOW CAN WE SOMEHOW BE STRCIKTER WITH THE TIME LOGGED? SHOULD WE ADD NOTES WHEN LOGGING TIME?
// TODO FIGURE OUT WAYS WE CAN VERIFY THE TIME LOGGED, AND VERIFY THE OUTPUT OF THE TIME LOGGED
// !! ADD A NOTES FIELD TO THE ACTIVITY MODEL and ask user what they did and their output

// !! FOR THE SCORE - implement feature where the more time you spent grinding the higher the score and the more intense badge you get
// !! the web app also changes colors based on the score
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
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not authenticated"})
		return
	}

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
