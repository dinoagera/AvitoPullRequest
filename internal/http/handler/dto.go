package handler

type AddTeamRequest struct {
	TeamName string       `json:"team_name"`
	Members  []TeamMember `json:"members"`
}
type TeamMember struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}
type SetUserActiveRequest struct {
	UserID   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}
type CreateRPRequest struct {
	ID       string `json:"pull_request_id"`
	Name     string `json:"pull_request_name"`
	AuthorID string `json:"author_id"`
}
type ReviewerStat struct {
	UserID        string `json:"user_id"`
	Username      string `json:"username"`
	AssignedCount int    `json:"assigned_count"`
}
type MergePRRequest struct {
	ID string `json:"pull_request_id"`
}
type ReassignRequest struct {
	PRID      string `json:"pull_request_id"`
	OldUserID string `json:"old_user_id"`
}
type PullRequestShort struct {
	ID       string `json:"pull_request_id"`
	Name     string `json:"pull_request_name"`
	AuthorID string `json:"author_id"`
	Status   string `json:"status"`
}
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}
type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
