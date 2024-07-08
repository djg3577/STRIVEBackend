package repository

import (
	"STRIVEBackend/pkg/models"
	"database/sql"
)

type ActivityRepository struct {
	DB *sql.DB
}

func (r *ActivityRepository) CreateActivity(activity *models.Activity) error {
	_, err := r.DB.Exec("INSERT INTO activities (user_id, type, duration, date) VALUES ($1, $2, $3, $4)",
		activity.UserID, activity.Type, activity.Duration, activity.Date)
	return err
}

func (r *ActivityRepository) GetActivityTotals(userID int) (*models.ActivityTotals, error) {
	rows, err := r.DB.Query("SELECT type, SUM(duration) FROM activities WHERE user_id = $1 GROUP BY type", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	activityTotals := models.ActivityTotals{ActivityTotals: make(map[string]int)}
	for rows.Next() {
		var activityType string
		var totalDuration int
		if err := rows.Scan(&activityType, &totalDuration); err != nil {
			return nil, err
		}
		activityTotals.ActivityTotals[activityType] = totalDuration
	}

	return &activityTotals, nil
}
