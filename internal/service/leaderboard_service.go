package service

import (
	"STRIVEBackend/internal/repository"
)

type LeaderBoardService struct {
	Repo repository.LeaderBoardRepositoryInterface
}

func (s *LeaderBoardService) GetTopScores() ([]repository.UserScore, error) {
	return s.Repo.GetTopScores()
}
