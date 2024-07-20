package repository

import (
	"STRIVEBackend/internal/util"
	"STRIVEBackend/pkg/models"
	"database/sql"
)

type AuthRepository struct {
	DB *sql.DB
}

func (r *AuthRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.DB.QueryRow("SELECT id, username, email, password, code FROM users WHERE email = $1", email).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password, &user.Code)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &user, nil
}

func (r *AuthRepository) CreateUser(tx *sql.Tx, user *models.User) (ID int, err error) {
	err = tx.QueryRow("INSERT INTO users (username, email, password, code, email_verified) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		user.Username, user.Email, user.Password, user.Code, false).Scan(&ID)
	return
}

func (r *AuthRepository) VerifyUserEmail(userID int) error {
	_, err := r.DB.Exec("UPDATE users SET email_verified = true WHERE id = $1", userID)
	return err
}

func (r *AuthRepository) DecodeJWT(token string) (*models.User, error) {
	userID, err := util.ValidateJWT(token)
	if err != nil {
		return nil, err
	}
	var user models.User
	err = r.DB.QueryRow("SELECT id, username, email FROM users WHERE id = $1", userID).Scan(
		&user.ID, &user.Username, &user.Email)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &user, nil
}

func (r *AuthRepository) GetUserIdByGithubId(githubUserId int) (int, error) {
	var userId int
	err := r.DB.QueryRow("SELECT id FROM users WHERE github_id = $1", githubUserId).Scan(&userId)
	return userId, err
}

func (r *AuthRepository) CreateUserFromGithub(githubUserId int) (int, error) {
	var userId int
	err := r.DB.QueryRow("INSERT INTO users (github_id) VALUES ($1) RETURNING id", githubUserId).Scan(&userId)
	return userId, err
}
