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
