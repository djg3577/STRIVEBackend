package repository

import (
	"STRIVEBackend/pkg/models"
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