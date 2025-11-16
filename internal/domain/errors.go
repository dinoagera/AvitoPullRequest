package domain

import "errors"

// Errors for team domain
var (
	ErrTeamExists   = errors.New("team already exists")
	ErrTeamNotFound = errors.New("team not found")
)

// Errors for user domain
var (
	ErrUserNotFound = errors.New("user not found")
)

// Errors for pullRequest domain
var (
	ErrPRNotFound         = errors.New("PR not found")
	ErrPRExists           = errors.New("PR already exists")
	ErrPRMerged           = errors.New("cannot reassign on merged PR")
	ErrNotAssigned        = errors.New("reviewer is not assigned to this PR")
	ErrNoCandidate        = errors.New("no active replacement candidate in team")
	ErrAuthorTeamNotFound = errors.New("author/team not found")
)
