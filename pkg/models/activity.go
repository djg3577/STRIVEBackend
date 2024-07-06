package models

import "time"

type Activity struct {
	ID       int       `json:"id"`
	UserID   int       `json:"user_id"`
	Type     string    `json:"type"`
	Duration int       `json:"duration"` // Duration in minutes
	Date     time.Time `json:"date"`
}
