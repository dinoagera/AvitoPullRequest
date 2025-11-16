package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/dinoagera/AvitoPullRequest/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TeamRepository struct {
	pool *pgxpool.Pool
}

func NewTeamRepository(pool *pgxpool.Pool) *TeamRepository {
	return &TeamRepository{pool: pool}
}
func (tr *TeamRepository) CreateTeam(ctx context.Context, team *domain.Team) error {
	tx, err := tr.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, "INSERT INTO teams (name) VALUES ($1)", team.Name)
	if err != nil {
		return fmt.Errorf("failed to create team: %w", err)
	}
	for _, user := range team.Members {
		_, err = tx.Exec(ctx, `
            INSERT INTO users (user_id, username, team_name, is_active)
            VALUES ($1, $2, $3, $4)
            ON CONFLICT (user_id) DO UPDATE SET
                username = EXCLUDED.username,
                team_name = EXCLUDED.team_name,
                is_active = EXCLUDED.is_active
        `, user.ID, user.Name, user.TeamName, user.IsActive)
		if err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}
	}
	return tx.Commit(ctx)
}
func (tr *TeamRepository) GetTeam(ctx context.Context, name string) (*domain.Team, error) {
	rows, err := tr.pool.Query(ctx, `
        SELECT user_id, username, team_name, is_active
        FROM users
        WHERE team_name = $1
    `, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get team: %w", err)
	}
	defer rows.Close()
	var members []domain.User
	for rows.Next() {
		var user domain.User
		err := rows.Scan(&user.ID, &user.Name, &user.TeamName, &user.IsActive)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		members = append(members, user)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error: %w", err)
	}
	return &domain.Team{
		Name:    name,
		Members: members,
	}, nil
}
func (tr *TeamRepository) TeamExists(ctx context.Context, name string) (bool, error) {
	var exists bool
	err := tr.pool.QueryRow(ctx, `
        SELECT EXISTS(SELECT 1 FROM teams WHERE name = $1)
    `, name).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check team exists: %w", err)
	}
	return exists, nil
}
func (tr *TeamRepository) GetTeamByUser(ctx context.Context, userID string) (*domain.Team, error) {
	var teamName string
	err := tr.pool.QueryRow(ctx, `
        SELECT team_name
        FROM users
        WHERE user_id = $1
    `, userID).Scan(&teamName)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrTeamNotFound
		}
		return nil, fmt.Errorf("failed to get team name: %w", err)
	}
	team, err := tr.GetTeam(ctx, teamName)
	if err != nil {
		if errors.Is(err, domain.ErrTeamNotFound) {
			return nil, err
		}
		return nil, fmt.Errorf("failed to get team by name: %w", err)
	}

	return team, nil
}
