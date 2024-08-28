package repository

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetTopScores(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening stud database", err)
	}
	defer db.Close()

	repo := &LeaderBoardRepository{DB: db}

	rows := sqlmock.NewRows([]string{"username", "score"}).
	AddRow("user1", 100).
	AddRow("user2", 90).
	AddRow("user3", 80)

	mock.ExpectQuery("SELECT u.username as username , SUM\\(count\\) as score FROM Users as u JOIN activity_summary a ON a.user_id = u.id GROUP BY u.username").
	WillReturnRows(rows)

	scores, err := repo.GetTopScores()

	assert.NoError(t, err)

	expectedScores := []UserScore{
		{Username: "user1", Score: 100},
		{Username: "user2", Score: 90},
		{Username: "user3", Score: 80},
	}

	assert.Equal(t, expectedScores, scores)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}


}