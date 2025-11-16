package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/dinoagera/AvitoPullRequest/internal/domain"
)

type TeamService struct {
	log  *slog.Logger
	repo TeamRepository
}

func NewTeamService(log *slog.Logger, repo TeamRepository) *TeamService {
	return &TeamService{
		log:  log,
		repo: repo,
	}
}

func (ts *TeamService) AddTeam(ctx context.Context, team *domain.Team) (*domain.Team, error) {
	exists, err := ts.repo.TeamExists(ctx, team.Name)
	if err != nil {
		ts.log.Info("failed to check exists team", "err", err)
		return nil, err
	}
	if exists {
		return nil, domain.ErrTeamExists
	}
	if err := ts.repo.CreateTeam(ctx, team); err != nil {
		ts.log.Info("failed to create team", "err", err)
		return nil, err
	}
	return team, nil
}

func (ts *TeamService) GetTeam(ctx context.Context, teamName string) (*domain.Team, error) {
	exists, err := ts.repo.TeamExists(ctx, teamName)
	if err != nil {
		ts.log.Info("failed to check exists team", "err", err)
		return nil, err
	}
	if !exists {
		return nil, domain.ErrTeamNotFound
	}
	team, err := ts.repo.GetTeam(ctx, teamName)
	if err != nil {
		if errors.Is(err, domain.ErrTeamNotFound) {
			return nil, err
		}
		ts.log.Info("failed to get team", "err", err)
		return nil, err
	}
	return team, nil
}
