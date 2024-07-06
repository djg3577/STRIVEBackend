package service

import (
	"STRIVEBackend/internal/repository"
	"STRIVEBackend/pkg/models"
	"time"
)

type ActivityService struct {
	Repo *repository.ActivityRepository
}

func (s *ActivityService) LogActivity(activity *models.Activity) error {
	activity.Date = time.Now()
	return s.Repo.CreateActivity(activity)
}
