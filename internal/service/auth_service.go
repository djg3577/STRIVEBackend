package service

import (
	"STRIVEBackend/internal/repository"
	"STRIVEBackend/pkg/models"
	"errors"
	"fmt"
	"gopkg.in/mail.v2"
	"math/rand"
	"os"
)

type AuthService struct {
	Repo *repository.AuthRepository
}

func (s *AuthService) Login(email string, password string) (*models.User, error) {
	user, err := s.Repo.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	if user.Password != password {
		return nil, errors.New("invalid credentials")
	}
	return user, nil
}

func (s *AuthService) SignUp(user *models.User) (int, error) {
	tx, err := s.Repo.DB.Begin()
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	verificationCode := generateValidCode()
	user.Code = verificationCode

	userID, err := s.Repo.CreateUser(tx, user)
	if err != nil {
		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// if err = s.sendVerificationEmail(user.Email, verificationCode); err != nil {
	// 	return 0, fmt.Errorf("failed to send verification email: %w", err)
	// }
	fmt.Println("Verification code: ", verificationCode)

	return userID, nil
}

func generateValidCode() int {
	return rand.Intn(1000000)
}

// !!  SET UP AMAZON SES TO SEND EMAILS
func (s *AuthService) sendVerificationEmail(email string, code string) error {
	m := mail.NewMessage()

	from := os.Getenv("EMAIL_FROM")
	m.SetHeader("From", from)

	m.SetHeader("To", email)

	m.SetHeader("Subject", "Verify Your Email")

	body := fmt.Sprintf("Your verification code is: %s", code)
	m.SetBody("text/plain", body)

	d := mail.NewDialer(
		os.Getenv("SMTP_HOST"),
		465,
		os.Getenv("SMTP_USERNAME"),
		os.Getenv("SMTP_PASSWORD"),
	)

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func (s *AuthService) VerifyEmail(email string, code int) error {
	user, err := s.Repo.GetUserByEmail(email)
	fmt.Println("THIS IS THE USER: ", user)
	if err != nil {
		return fmt.Errorf("failed to get user by email: %w", err)
	}

	if user.EmailVerified {
		return errors.New("email already verified")
	}
	fmt.Println("THIS IS THE USER CODE: ", user.Code)
	fmt.Println("THIS IS THE CODE: ", code)
	if user.Code != code {
		return errors.New("invalid verification code")
	}

	return s.Repo.VerifyUserEmail(user.ID)
}

func (s *AuthService) DecodeJWT(token string) (*models.User, error) {
	return s.Repo.DecodeJWT(token)
}
