package models

type User struct {
	ID            int    `json:"user_id"`
	Username      string `json:"username"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	Code          int    `json:"code"`
	EmailVerified bool   `json:"email_verified"`
}

type GitHubUser struct {
	Login string `json:"login"`
	ID int `json:"id"`
}