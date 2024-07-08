package repository

import (
	"STRIVEBackend/internal/util"
	"STRIVEBackend/pkg/models"
	"database/sql"
	"fmt"
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
	fmt.Println("INSIDE OF DECODOING JWT: asdasd")
	userID, err := util.ValidateJWT(token)
	if err != nil {
		return nil, err
	}
	fmt.Println("USERID: ", userID)
	var user models.User
	err = r.DB.QueryRow("SELECT id, username, email FROM users WHERE id = $1", userID).Scan(
		&user.ID, &user.Username, &user.Email)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	fmt.Println("USER: ", user)
	return &user, nil
}
