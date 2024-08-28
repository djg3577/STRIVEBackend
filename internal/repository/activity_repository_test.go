package repository

import (
	"STRIVEBackend/pkg/models"
	"database/sql"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetActivityTotals(t *testing.T) {
	// Create a new mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Create a new ActivityRepository with the mock db
	repo := &ActivityRepository{DB: db}

	// Set up the expected query
	rows := sqlmock.NewRows([]string{"activity_name", "total_duration"}).
		AddRow("running", 3600).
		AddRow("cycling", 7200)

	mock.ExpectQuery("SELECT activity_name, total_duration FROM activity_summary WHERE user_id = ?").
		WithArgs(1).
		WillReturnRows(rows)

	// Call the method
	totals, err := repo.GetActivityTotals(1)

	// Assert no errors occurred
	assert.NoError(t, err)

	// Assert the returned totals match what we expect
	expectedTotals := &models.ActivityTotals{
		ActivityTotals: map[string]int{
			"running": 3600,
			"cycling": 7200,
		},
	}
	assert.Equal(t, expectedTotals, totals)

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetActivityTotalsError(t *testing.T) {
	// Create a new mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Create a new ActivityRepository with the mock db
	repo := &ActivityRepository{DB: db}

	// Set up the expected query to return an error
	mock.ExpectQuery("SELECT activity_name, total_duration FROM activity_summary WHERE user_id = ?").
		WithArgs(1).
		WillReturnError(sqlmock.ErrCancelled)

	// Call the method
	totals, err := repo.GetActivityTotals(1)

	// Assert an error occurred
	assert.Error(t, err)
	assert.Nil(t, totals)

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetActivityDates(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := &ActivityRepository{DB: db}

	mock.ExpectQuery("SELECT date, COUNT\\(\\*\\) FROM activities WHERE user_id = \\$1 GROUP BY date").
		WithArgs(1).
		WillReturnError(sqlmock.ErrCancelled)

	dates, err := repo.GetActivityDates(1)

	assert.Error(t, err)
	assert.Nil(t, dates)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestGetUserIdByGithubId(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stud database connection", err)
	}
	defer db.Close()
	repo := &ActivityRepository{DB: db}

	t.Run("Successful query", func(t *testing.T) {
		expectedUserId := 42
		githubUserId := 12345

		rows := sqlmock.NewRows([]string{"id"}).AddRow(expectedUserId)
		mock.ExpectQuery("SELECT id FROM users WHERE github_id = \\$1").
		WithArgs(githubUserId).
		WillReturnRows(rows)

		userId, err := repo.GetUserIdByGithubId(githubUserId)

		assert.NoError(t, err)
		assert.Equal(t, expectedUserId, userId)
	})

	t.Run("User not found", func(t *testing.T) {
		githubUserId := 67890

		mock.ExpectQuery("SELECT id FROM users WHERE github_id = \\$1").
		WithArgs(githubUserId).
		WillReturnError(sql.ErrNoRows)

		userId, err := repo.GetUserIdByGithubId(githubUserId)

		assert.Error(t, err)
		assert.Equal(t, sql.ErrNoRows, err)
		assert.Equal(t, 0, userId)
	})

	t.Run("Database error", func(t *testing.T) {
		githubUserId := 12359

		mock.ExpectQuery("SELECT id FROM users WHERE github_id = \\$1").
		WithArgs(githubUserId).
		WillReturnError(fmt.Errorf("database connection lost"))

		userId, err := repo.GetUserIdByGithubId(githubUserId)

		assert.Error(t, err)
		assert.Equal(t, "database connection lost", err.Error())
		assert.Equal(t, 0, userId)
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestCreateUserFromGithub(t *testing.T){	
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stud database connection", err)
	}
	defer db.Close()
	repo := &ActivityRepository{DB: db}


	t.Run("Successful user creation", func(t *testing.T) {
		githubUser := &models.GitHubUser{
			ID: 12345,
			Login: "testuser",
		}
		expectedUserId := 1

		rows := sqlmock.NewRows([]string{"id"}).AddRow(expectedUserId)
		mock.ExpectQuery("INSERT INTO users \\(github_id, username\\) VALUES \\(\\$1, \\$2\\) RETURNING id").
		WithArgs(githubUser.ID,githubUser.Login ).
		WillReturnRows(rows)

		userId, err := repo.CreateUserFromGithub(githubUser)

		assert.NoError(t, err)
		assert.Equal(t, expectedUserId, userId)
	})
}
