package repository

import (
	"STRIVEBackend/pkg/models"
	"database/sql"
)

type UserRepository struct {
	DB *sql.DB
}

func (r *UserRepository) GetUserByID(id int) (*models.User, error) {
	var user models.User
	err := r.DB.QueryRow("SELECT id, username, email, password FROM users WHERE id = $1", id).Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
