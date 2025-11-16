package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/dinoagera/AvitoPullRequest/internal/domain"
)

type UserService struct {
	log      *slog.Logger
	userRepo UserRepository
	prRepo   PRRepository
}

func NewUserService(log *slog.Logger, userRepo UserRepository, prRepo PRRepository) *UserService {
	return &UserService{
		log:      log,
		userRepo: userRepo,
		prRepo:   prRepo,
	}
}
func (us *UserService) SetActive(ctx context.Context, userID string, isActive bool) (*domain.User, error) {
	user, err := us.userRepo.SetActive(ctx, userID, isActive)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, err
		}
		us.log.Info("failed to set active", "err", err)
		return nil, err
	}
	return user, nil
}
func (us *UserService) GetReview(ctx context.Context, userID string) ([]domain.PullRequest, error) {
	_, err := us.userRepo.GetUser(ctx, userID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, err
		}
		us.log.Info("failed to get user", "err", err)
		return nil, err
	}
	prs, err := us.prRepo.GetPRsByReviewer(ctx, userID)
	if err != nil {
		us.log.Info("failed to get prs by reviewer", "err", err)
		return nil, err
	}
	return prs, nil
}
