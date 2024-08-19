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

type TokenType string

type tokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

const (
	TokenTypeGithub TokenType = "github"
	TokenTypeJWT    TokenType = "jwt"
)

type AuthService struct {
	Repo *repository.AuthRepository
}

func (authService *AuthService) Login(email string, password string) (*models.User, error) {
	user, err := authService.Repo.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	if user.Password != password {
		return nil, errors.New("invalid credentials")
	}
	return user, nil
}

func (authService *AuthService) SignUp(user *models.User) (int, error) {
	tx, err := authService.Repo.DB.Begin()
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

	userID, err := authService.Repo.CreateUser(tx, user)
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
func (authService *AuthService) sendVerificationEmail(email string, code string) error {
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

func (authService *AuthService) VerifyEmail(email string, code int) error {
	user, err := authService.Repo.GetUserByEmail(email)
	if err != nil {
		return fmt.Errorf("failed to get user by email: %w", err)
	}

	if user.EmailVerified {
		return errors.New("email already verified")
	}
	if user.Code != code {
		return errors.New("invalid verification code")
	}

	return authService.Repo.VerifyUserEmail(user.ID)
}

func (authService *AuthService) AuthenticateUser(authHeader string) (*models.User, error) {
	if authHeader == "" {
		return nil, fmt.Errorf("missing Authorization header")
	}
	token, tokenType, err := authService._parseHeader(authHeader)
	if err != nil {
		return nil, fmt.Errorf("error getting token")
	}

	switch tokenType {
	case "github":
		return authService._authenticateGithubUser(token)
	default:
		return authService._authenticateJWTUser(token)
	}
}

func (authService *AuthService) _parseHeader(authHeader string) (token string, tokenType string, error error) { // modify this so that tokenType is of github or jwt
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || (strings.ToLower(parts[0]) != "bearer" && strings.ToLower(parts[0]) != "github") {
		return "", "", fmt.Errorf("invalid Authorization header")
	}

	token = parts[1]
	if token == "" {
		return "", "", fmt.Errorf("missing token")
	}

	tokenType = strings.ToLower(parts[0])
	return token, tokenType, nil
}

func (authService *AuthService) _authenticateGithubUser(token string) (*models.User, error) {
	githubUser, err := authService.GetGitHubUser(token)
	if err != nil {
		return nil, fmt.Errorf("invalid Github token: %W", err)
	}
	return &models.User{
		ID:       githubUser.ID,
		Username: githubUser.Login,
	}, nil
}

func (authService *AuthService) _authenticateJWTUser(token string) (*models.User, error) {
	user, err := authService.DecodeJWT(token)
	if err != nil {
		return nil, fmt.Errorf("error decoding JWT: %w", err)
	}
	return user, nil
}

func (authService *AuthService) DecodeJWT(token string) (*models.User, error) {
	return authService.Repo.DecodeJWT(token)
}

func (authService *AuthService) GetGitHubUser(token string) (*models.GitHubUser, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

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

func (authService *AuthService) ExchangeGitHubCode(code string) (string, *models.GitHubUser, error) {
	clientID := os.Getenv("GITHUB_CLIENT_ID")
	clientSecret := os.Getenv("GITHUB_CLIENT_SECRET")
	redirectURI := os.Getenv("GITHUB_REDIRECT_URI")
	fmt.Print(redirectURI)
	fmt.Println()

	tokenURL := "https://github.com/login/oauth/access_token"

	requestBody, _ := json.Marshal(map[string]string{
		"client_id":     clientID,
		"client_secret": clientSecret,
		"code":          code,
		"redirect_uri":  redirectURI,
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
	var tokenResponse tokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return "", nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if tokenResponse.AccessToken == "" {
		return "", nil, fmt.Errorf("access token is empty")
	}

	user, err := authService.GetGitHubUser(tokenResponse.AccessToken)
	if err != nil {
		return "", nil, fmt.Errorf("failed to get GitHub user: %w", err)
	}

	return tokenResponse.AccessToken, user, nil
}

func (s *AuthService) GetOrCreateUserIdFromGithub(githubUser *models.GitHubUser) (int, error) {
	userId, err := s.Repo.GetUserIdByGithubId(githubUser.ID)
	if err == sql.ErrNoRows {
		// User doesn't exist, create a new one
		userId, err = s.Repo.CreateUserFromGithub(githubUser)
		if err != nil {
			return 0, err
		}
	} else if err != nil {
		return 0, err
	}
	return userId, nil
}
