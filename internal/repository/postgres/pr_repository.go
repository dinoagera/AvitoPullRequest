package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/dinoagera/AvitoPullRequest/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PRRepository struct {
	pool *pgxpool.Pool
}

func NewPRRepository(pool *pgxpool.Pool) *PRRepository {
	return &PRRepository{pool: pool}
}
func (pr *PRRepository) CreatePR(ctx context.Context, pullRequest *domain.PullRequest) error {
	_, err := pr.pool.Exec(ctx, `
        INSERT INTO pull_requests (pr_id, title, author_id, status, assigned_reviewers, created_at)
        VALUES ($1, $2, $3, $4, $5, $6)
    `, pullRequest.ID, pullRequest.Name, pullRequest.AuthorID, pullRequest.Status, pullRequest.AssignedReviewers, pullRequest.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create PR: %w", err)
	}
	return nil
}
func (pr *PRRepository) PRExists(ctx context.Context, id string) (bool, error) {
	var exists bool
	err := pr.pool.QueryRow(ctx, `
        SELECT EXISTS(SELECT 1 FROM pull_requests WHERE pr_id = $1)
    `, id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check PR exists: %w", err)
	}
	return exists, nil
}
func (pr *PRRepository) GetPR(ctx context.Context, id string) (*domain.PullRequest, error) {
	var pullRequest domain.PullRequest
	var createdAt, mergedAt *string
	err := pr.pool.QueryRow(ctx, `
        SELECT pr_id, title, author_id, status, assigned_reviewers, created_at, merged_at
        FROM pull_requests
        WHERE pr_id = $1
    `, id).Scan(&pullRequest.ID, &pullRequest.Name, &pullRequest.AuthorID, &pullRequest.Status, &pullRequest.AssignedReviewers, &createdAt, &mergedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrPRNotFound
		}
		return nil, fmt.Errorf("failed to get PR: %w", err)
	}
	pullRequest.CreatedAt = createdAt
	pullRequest.MergedAt = mergedAt
	return &pullRequest, nil
}
func (pr *PRRepository) UpdatePR(ctx context.Context, pullRequest *domain.PullRequest) error {
	_, err := pr.pool.Exec(ctx, `
        UPDATE pull_requests
        SET status = $1, assigned_reviewers = $2, merged_at = $3
        WHERE pr_id = $4
    `, pullRequest.Status, pullRequest.AssignedReviewers, pullRequest.MergedAt, pullRequest.ID)
	if err != nil {
		return fmt.Errorf("failed to update PR: %w", err)
	}
	return nil
}
func (pr *PRRepository) GetPRsByReviewer(ctx context.Context, reviewerID string) ([]domain.PullRequest, error) {
	rows, err := pr.pool.Query(ctx, `
        SELECT pr_id, title, author_id, status, assigned_reviewers, created_at, merged_at
        FROM pull_requests
        WHERE $1 = ANY(assigned_reviewers)
    `, reviewerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get PRs by reviewer: %w", err)
	}
	defer rows.Close()
	var prs []domain.PullRequest
	for rows.Next() {
		var pullRequest domain.PullRequest
		var createdAt, mergedAt *string
		err := rows.Scan(&pullRequest.ID, &pullRequest.Name, &pullRequest.AuthorID, &pullRequest.Status, &pullRequest.AssignedReviewers, &createdAt, &mergedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan PR: %w", err)
		}
		pullRequest.CreatedAt = createdAt
		pullRequest.MergedAt = mergedAt
		prs = append(prs, pullRequest)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return prs, nil
}
