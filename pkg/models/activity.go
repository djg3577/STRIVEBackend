package models

type Activity struct {
	ID           int    `json:"id"`
	UserID       int    `json:"user_id"`
	ActivityName string `json:"activity_name"`
	Duration     int    `json:"duration"` // Duration in minutes
	Date         string `json:"date"`
}

type ActivityTotals struct {
	ActivityTotals map[string]int `json:"activity_totals"`
}

type ActivityDates struct {
	ActivityDates []ActivityDate `json:"activity_dates"`
}

type ActivityDate struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}
