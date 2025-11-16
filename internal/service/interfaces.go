package service

import (
	"context"

	"github.com/dinoagera/AvitoPullRequest/internal/domain"
)

type TeamRepository interface {
	CreateTeam(ctx context.Context, team *domain.Team) error
	GetTeam(ctx context.Context, name string) (*domain.Team, error)
	TeamExists(ctx context.Context, name string) (bool, error)
	GetTeamByUser(ctx context.Context, userID string) (*domain.Team, error)
}
type UserRepository interface {
	SetActive(ctx context.Context, userID string, isActive bool) (*domain.User, error)
	GetUser(ctx context.Context, userID string) (*domain.User, error)
}
type PRRepository interface {
	CreatePR(ctx context.Context, pr *domain.PullRequest) error
	GetPR(ctx context.Context, id string) (*domain.PullRequest, error)
	UpdatePR(ctx context.Context, pr *domain.PullRequest) error
	PRExists(ctx context.Context, id string) (bool, error)
	GetPRsByReviewer(ctx context.Context, reviewerID string) ([]domain.PullRequest, error)
}
type StatsRepository interface {
	GetReviewerStats(ctx context.Context) ([]domain.ReviewerStat, error)
}
