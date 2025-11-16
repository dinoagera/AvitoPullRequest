package service

import (
	"context"
	"log/slog"

	"github.com/dinoagera/AvitoPullRequest/internal/domain"
)

type StatsService struct {
	log  *slog.Logger
	repo StatsRepository
}

func NewStatsService(log *slog.Logger, repo StatsRepository) *StatsService {
	return &StatsService{
		log:  log,
		repo: repo,
	}
}

func (ss *StatsService) GetReviewerStats(ctx context.Context) ([]domain.ReviewerStat, error) {
	return ss.repo.GetReviewerStats(ctx)
}
