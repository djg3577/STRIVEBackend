package service

import (
	"STRIVEBackend/internal/repository"
	"STRIVEBackend/pkg/models"
)

type UserService struct {
	Repo *repository.UserRepository
}

func (s *UserService) CreateUser(user *models.User) (ID int, err error) {
	return s.Repo.CreateUser(user)
}

func (s *UserService) GetUser(ID int) (*models.User, error) {
	return s.Repo.GetUserByID(ID)
}
