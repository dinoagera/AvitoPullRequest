package handler

import (
	"context"

	"github.com/dinoagera/AvitoPullRequest/internal/domain"
)

type TeamService interface {
	AddTeam(ctx context.Context, team *domain.Team) (*domain.Team, error)
	GetTeam(ctx context.Context, teamName string) (*domain.Team, error)
}
type UserService interface {
	SetActive(ctx context.Context, userID string, isActive bool) (*domain.User, error)
	GetReview(ctx context.Context, userID string) ([]domain.PullRequest, error)
}
type PRService interface {
	CreatePR(ctx context.Context, pr *domain.PullRequest) (*domain.PullRequest, error)
	MergePR(ctx context.Context, id string) (*domain.PullRequest, error)
	ReassignReviewer(ctx context.Context, prID, OldUserID string) (*domain.PullRequest, string, error)
}
type StatsService interface {
	GetReviewerStats(ctx context.Context) ([]domain.ReviewerStat, error)
}
