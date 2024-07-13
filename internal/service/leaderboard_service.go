package service

import (
	"STRIVEBackend/internal/repository"
)

type LeaderboardService struct {
	Repo *repository.LeaderBoardRepository
}

func (s *LeaderboardService) GetTopScores() ([]repository.UserScore, error) {
	return s.Repo.GetTopScores()
}
