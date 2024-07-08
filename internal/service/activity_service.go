package service

import (
	"STRIVEBackend/internal/repository"
	"STRIVEBackend/pkg/models"
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
