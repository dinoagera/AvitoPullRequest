package postgres

import (
	"context"

	"github.com/dinoagera/AvitoPullRequest/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}
func (ur *UserRepository) GetUser(ctx context.Context, userID string) (*domain.User, error) {
	var user domain.User
	err := ur.pool.QueryRow(ctx, `
		SELECT user_id, username, team_name, is_active 
		FROM users
		WHERE user_id = $1
	`).Scan(&user.ID, &user.Name, &user.TeamName, &user.IsActive)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func (ur *UserRepository) SetActive(ctx context.Context, userID string, isActive bool) (*domain.User, error) {
	result, err := ur.pool.Exec(ctx, `
		UPDATE users
		SET is_active = $1
		WHERE user_id = $2
	`, isActive, userID)
	if err != nil {
		return nil, err
	}
	if result.RowsAffected() == 0 {
		return nil, domain.ErrUserNotFound
	}
	user, err := ur.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}
