package service

import (
	"STRIVEBackend/internal/repository"
	"STRIVEBackend/pkg/models"
	"database/sql"
)

type ActivityService struct {
	Repo *repository.ActivityRepository
}

func (s *ActivityService) LogActivity(activity *models.Activity) error {
	return s.Repo.CreateActivity(activity)
}

func (s *ActivityService) GetActivityTotals(userID int) (*models.ActivityTotals, error) {
	return s.Repo.GetActivityTotals(userID)
}

func (s *ActivityService) GetActivityDates(userID int) (*models.ActivityDates, error) {
	return s.Repo.GetActivityDates(userID)
}

func (s *ActivityService) GetOrCreateUserIdFromGithub(githubUser *models.GitHubUser) (int, error) {
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
