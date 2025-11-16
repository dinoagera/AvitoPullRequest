package service

import (
	"context"
	"errors"
	"log/slog"
	"math/rand"
	"time"

	"github.com/dinoagera/AvitoPullRequest/internal/domain"
)

type PRService struct {
	log      *slog.Logger
	userRepo UserRepository
	teamRepo TeamRepository
	prRepo   PRRepository
}

func NewPRService(log *slog.Logger, userRepo UserRepository, teamRepo TeamRepository, prRepo PRRepository) *PRService {
	return &PRService{
		log:      log,
		userRepo: userRepo,
		teamRepo: teamRepo,
		prRepo:   prRepo,
	}
}
func (pr *PRService) CreatePR(ctx context.Context, pullRequest *domain.PullRequest) (*domain.PullRequest, error) {
	author, err := pr.userRepo.GetUser(ctx, pullRequest.AuthorID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, domain.ErrAuthorTeamNotFound
		}
		pr.log.Info("failed to get author", "err", err)
		return nil, err
	}
	team, err := pr.teamRepo.GetTeamByUser(ctx, pullRequest.AuthorID)
	if err != nil {
		if errors.Is(err, domain.ErrTeamNotFound) {
			return nil, domain.ErrAuthorTeamNotFound
		}
		pr.log.Info("failed to get author team", "err", err)
		return nil, err
	}
	exists, err := pr.prRepo.PRExists(ctx, pullRequest.ID)
	if err != nil {
		pr.log.Info("failed to check pr exists", "err", err)
		return nil, err
	}
	if exists {
		return nil, domain.ErrPRExists
	}
	reviewers := pr.selectReviewers(team, author.ID)
	pullRequestChanged := &domain.PullRequest{
		ID:                pullRequest.ID,
		Name:              pullRequest.Name,
		AuthorID:          pullRequest.AuthorID,
		Status:            "OPEN",
		AssignedReviewers: reviewers,
		CreatedAt:         pr.now(),
	}
	if err := pr.prRepo.CreatePR(ctx, pullRequestChanged); err != nil {
		pr.log.Info("failed to create pr", "err", err)
		return nil, err
	}
	return pullRequestChanged, nil
}
func (pr *PRService) MergePR(ctx context.Context, prID string) (*domain.PullRequest, error) {
	pullRequest, err := pr.prRepo.GetPR(ctx, prID)
	if err != nil {
		if errors.Is(err, domain.ErrPRNotFound) {
			return nil, err
		}
		pr.log.Info("failed to get pr", "err", err)
		return nil, err
	}
	if pullRequest.Status == "MERGED" {
		return pullRequest, nil
	}
	pullRequest.Status = "MERGED"
	pullRequest.MergedAt = pr.now()
	if err := pr.prRepo.UpdatePR(ctx, pullRequest); err != nil {
		pr.log.Info("failed to update pr", "err", err)
		return nil, err
	}

	return pullRequest, nil
}
func (ps *PRService) ReassignReviewer(ctx context.Context, prID, oldUserID string) (*domain.PullRequest, string, error) {
	pr, err := ps.prRepo.GetPR(ctx, prID)
	if err != nil {
		if errors.Is(err, domain.ErrPRNotFound) {
			return nil, "", err
		}
		ps.log.Info("failed to get pr", "err", err)
		return nil, "", err
	}
	if pr.Status == "MERGED" {
		return nil, "", domain.ErrPRMerged
	}
	oldUserIndex := -1
	for i, reviewerID := range pr.AssignedReviewers {
		if reviewerID == oldUserID {
			oldUserIndex = i
			break
		}
	}
	if oldUserIndex == -1 {
		return nil, "", domain.ErrNotAssigned
	}
	team, err := ps.teamRepo.GetTeamByUser(ctx, oldUserID)
	if err != nil {
		if errors.Is(err, domain.ErrTeamNotFound) {
			return nil, "", domain.ErrTeamNotFound
		}
		ps.log.Info("failed to get team of old reviewer", "err", err)
		return nil, "", err
	}
	newReviewerID := ps.findNewReviewer(team, pr.AuthorID, oldUserID)
	if newReviewerID == "" {
		return nil, "", domain.ErrNoCandidate
	}
	pr.AssignedReviewers[oldUserIndex] = newReviewerID
	if err := ps.prRepo.UpdatePR(ctx, pr); err != nil {
		ps.log.Info("failed to update pr", "err", err)
		return nil, "", err
	}
	return pr, newReviewerID, nil
}
func (pr *PRService) selectReviewers(team *domain.Team, excludeID string) []string {
	var candidates []string
	for _, user := range team.Members {
		if user.ID != excludeID && user.IsActive {
			candidates = append(candidates, user.ID)
		}
	}
	if len(candidates) > 2 {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(candidates), func(i, j int) {
			candidates[i], candidates[j] = candidates[j], candidates[i]
		})
		return candidates[:2]
	}
	return candidates
}
func (pr *PRService) now() *string {
	now := time.Now().Format(time.RFC3339)
	return &now
}
func (ps *PRService) findNewReviewer(team *domain.Team, authorID, excludeID string) string {
	for _, user := range team.Members {
		if user.ID != authorID && user.ID != excludeID && user.IsActive {
			return user.ID
		}
	}
	return ""
}
