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
	err := r.DB.QueryRow("SELECT id, username, email, password FROM users WHERE id = $1", id).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &user, nil
}

func (r *UserRepository) CreateUser(user *models.User) (ID int, err error) {
	err = r.DB.QueryRow("INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id",
		user.Username, user.Email, user.Password).Scan(&ID)
	return
}
