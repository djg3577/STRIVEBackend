package service

import (
	"STRIVEBackend/internal/repository"
	"STRIVEBackend/pkg/models"
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"

	"gopkg.in/mail.v2"
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
	if err != nil {
		return fmt.Errorf("failed to get user by email: %w", err)
	}

	if user.EmailVerified {
		return errors.New("email already verified")
	}
	if user.Code != code {
		return errors.New("invalid verification code")
	}

	return s.Repo.VerifyUserEmail(user.ID)
}

func (s *AuthService) AuthenticateUser(authHeader string) (*models.User, error) {
	if authHeader == "" {
			return nil, fmt.Errorf("missing Authorization header")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || (strings.ToLower(parts[0]) != "bearer" && strings.ToLower(parts[0]) != "github") {
			return nil, fmt.Errorf("invalid Authorization header")
	}

	token := parts[1]
	if token == "" {
			return nil, fmt.Errorf("missing token")
	}

	if strings.ToLower(parts[0]) == "github" {
			// Handle GitHub token
			githubUser, err := s.GetGitHubUser(token)
			if err != nil {
					return nil, fmt.Errorf("invalid GitHub token: %w", err)
			}
			// Convert GitHub user to your User model
			user := &models.User{
					ID:    githubUser.ID,
					Username:  githubUser.Login,
			}
			return user, nil
	} else {
			// Handle JWT token
			user, err := s.DecodeJWT(token)
			if err != nil {
					return nil, fmt.Errorf("error decoding JWT: %w", err)
			}
			return user, nil
	}
}

func (s *AuthService) DecodeJWT(token string) (*models.User, error) {
	return s.Repo.DecodeJWT(token)
}

func (s *AuthService) GetGitHubUser(token string) (*models.GitHubUser, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+ token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var user models.GitHubUser

	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &user, nil
}

func (s *AuthService) ExchangeGitHubCode(code string) (string, *models.GitHubUser, error) {
	clientID := os.Getenv("GITHUB_CLIENT_ID")
	clientSecret := os.Getenv("GITHUB_CLIENT_SECRET")
	redirectURI := os.Getenv("GITHUB_REDIRECT_URI")
	fmt.Print(redirectURI)
	fmt.Println()

	tokenURL := "https://github.com/login/oauth/access_token"

	requestBody, _ := json.Marshal(map[string]string{
		"client_id": clientID,
		"client_secret": clientSecret,
		"code": code,
		"redirect_uri": redirectURI,
	})

	req, err := http.NewRequest("POST", tokenURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var tokenResponse struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		Scope       string `json:"scope"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return "", nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if tokenResponse.AccessToken == "" {
		return "", nil, fmt.Errorf("access token is empty")
	}

	user, err := s.GetGitHubUser(tokenResponse.AccessToken)
	if err != nil {
		return "", nil, fmt.Errorf("failed to get GitHub user: %w", err)
	}

	return tokenResponse.AccessToken, user, nil
}

func (s *AuthService) GetOrCreateUserIdFromGithub(githubUserId int) (int, error) {
	userId, err := s.Repo.GetUserIdByGithubId(githubUserId)
	if err == sql.ErrNoRows {
			// User doesn't exist, create a new one
			userId, err = s.Repo.CreateUserFromGithub(githubUserId)
			if err != nil {
					return 0, err
			}
	} else if err != nil {
			return 0, err
	}
	return userId, nil
}