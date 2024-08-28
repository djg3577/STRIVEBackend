package repository

import (
	"STRIVEBackend/pkg/models"
	"database/sql"
	"fmt"
	"time"
)

type ActivityRepository struct {
	DB *sql.DB
}

func (r *ActivityRepository) CreateActivity(activity *models.Activity) error {
	_, err := r.DB.Exec("INSERT INTO activities (user_id, activity_name, duration, date) VALUES ($1, $2, $3, $4)",
		activity.UserID, activity.ActivityName, activity.Duration, activity.Date)
	return err
}

func (r *ActivityRepository) GetActivityTotals(userID int) (*models.ActivityTotals, error) {
	rows, err := r.DB.Query("SELECT activity_name, total_duration FROM activity_summary WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	activityTotals := models.ActivityTotals{ActivityTotals: make(map[string]int)}
	for rows.Next() {
		var activityName string
		var totalDuration int

		if err := rows.Scan(&activityName, &totalDuration); err != nil {
			return nil, err
		}
		activityTotals.ActivityTotals[activityName] = totalDuration
	}

	return &activityTotals, nil
}

func (r *ActivityRepository) GetActivityDates(userID int) (*models.ActivityDates, error) {
	rows, err := r.DB.Query("SELECT date, COUNT(*) FROM activities WHERE user_id = $1 GROUP BY date", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	activityDates := models.ActivityDates{ActivityDates: make([]models.ActivityDate, 0)}
	for rows.Next() {
		var date string
		var count int

		if err := rows.Scan(&date, &count); err != nil {
			return nil, err
		}

		parsedDate, err := time.Parse(time.RFC3339, date)
		if err != nil {
			return nil, fmt.Errorf("error parsing date: %v", err)
		}
		formattedDate := parsedDate.Format("2006-01-02")
		activityDates.ActivityDates = append(activityDates.ActivityDates, models.ActivityDate{Date: formattedDate, Count: count})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &activityDates, nil
}

func (r *ActivityRepository) GetUserIdByGithubId(githubUserId int) (int, error) {
	var userId int
	err := r.DB.QueryRow("SELECT id FROM users WHERE github_id = $1", githubUserId).Scan(&userId)
	return userId, err
}

func (r *ActivityRepository) CreateUserFromGithub(githubUser *models.GitHubUser) (int, error) {
	var userId int
	err := r.DB.QueryRow("INSERT INTO users (github_id, username) VALUES ($1, $2) RETURNING id", githubUser.ID, githubUser.Login).Scan(&userId)
	return userId, err
}
