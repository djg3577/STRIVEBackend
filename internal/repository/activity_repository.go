package repository

import (
	"STRIVEBackend/pkg/models"
	"database/sql"
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
