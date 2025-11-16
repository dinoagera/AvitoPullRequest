package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type UserHandler struct {
	log     *slog.Logger
	service UserService
}

func NewUserHandler(log *slog.Logger, service UserService) *UserHandler {
	return &UserHandler{
		log:     log,
		service: service,
	}
}
func (us *UserHandler) SetActive(w http.ResponseWriter, r *http.Request) {
	var req SetUserActiveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteErrorResponse(w, "BAD_REQUEST", "user_id is required", http.StatusBadRequest)
		return
	}
	if req.UserID == "" {
		WriteErrorResponse(w, "BAD_REQUEST", "user_id is required", http.StatusBadRequest)
		return
	}
	user, err := us.service.SetActive(r.Context(), req.UserID, req.IsActive)
	if err != nil {
		us.log.Info("failed to set active", "err", err)
		WriteError(w, err)
		return
	}
	res := domainTodtoUser(user)
	WriteJSONResponse(w, map[string]interface{}{"user": res}, http.StatusOK)
}
func (us *UserHandler) GetReview(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		WriteErrorResponse(w, "BAD_REQUEST", "user_id is required", http.StatusBadRequest)
		return
	}
	result, err := us.service.GetReview(r.Context(), userID)
	if err != nil {
		us.log.Info("failed to get review", "err", err)
		WriteError(w, err)
		return
	}
	shortPRs := make([]PullRequestShort, len(result))
	for i, pr := range result {
		shortPRs[i] = PullRequestShort{
			ID:       pr.ID,
			Name:     pr.Name,
			AuthorID: pr.AuthorID,
			Status:   pr.Status,
		}
	}
	WriteJSONResponse(w, map[string]interface{}{"user_id": userID, "pull_requests": shortPRs}, http.StatusOK)
}
