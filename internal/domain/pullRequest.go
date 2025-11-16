package domain

type PullRequest struct {
	ID                string
	Name              string
	AuthorID          string
	Status            string // "OPEN", "MERGED"
	AssignedReviewers []string
	CreatedAt         *string // nullable
	MergedAt          *string // nullable
}
