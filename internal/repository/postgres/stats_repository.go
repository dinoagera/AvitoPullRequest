// repository/postgres/stats_repository.go
package postgres

import (
	"context"
	"fmt"

	"github.com/dinoagera/AvitoPullRequest/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StatsRepository struct {
	pool *pgxpool.Pool
}

func NewStatsRepository(pool *pgxpool.Pool) *StatsRepository {
	return &StatsRepository{pool: pool}
}

func (sr *StatsRepository) GetReviewerStats(ctx context.Context) ([]domain.ReviewerStat, error) {
	rows, err := sr.pool.Query(ctx, `
        SELECT u.user_id, u.username, COUNT(p.assigned_reviewers) as assigned_count
        FROM users u
        LEFT JOIN pull_requests p ON u.user_id = ANY(p.assigned_reviewers)
        GROUP BY u.user_id, u.username
    `)
	if err != nil {
		return nil, fmt.Errorf("failed to get reviewer stats: %w", err)
	}
	defer rows.Close()
	var stats []domain.ReviewerStat
	for rows.Next() {
		var stat domain.ReviewerStat
		err := rows.Scan(&stat.UserID, &stat.Username, &stat.AssignedCount)
		if err != nil {
			return nil, fmt.Errorf("failed to scan reviewer stat: %w", err)
		}
		stats = append(stats, stat)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return stats, nil
}
