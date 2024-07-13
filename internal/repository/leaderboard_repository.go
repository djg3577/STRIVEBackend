package repository

import "database/sql"

type LeaderBoardRepository struct {
	DB *sql.DB
}

type UserScore struct {
	Username string `json:"username"`
	Score    int    `json:"score"`
}

func (r *LeaderBoardRepository) GetTopScores() ([]UserScore, error) {
	rows, err := r.DB.Query("SELECT u.username as username , SUM(count) as score FROM Users as u JOIN activity_summary a ON a.user_id = u.id GROUP BY u.username")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var topScores []UserScore
	for rows.Next() {
		var score UserScore
		if err := rows.Scan(&score.Username, &score.Score); err != nil {
			return nil, err
		}
		topScores = append(topScores, score)
	}

	return topScores, nil
}
