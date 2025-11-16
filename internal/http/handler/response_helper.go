package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/dinoagera/AvitoPullRequest/internal/domain"
)

func WriteJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func WriteErrorResponse(w http.ResponseWriter, code, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{
		Error: ErrorDetail{
			Code:    code,
			Message: message,
		},
	})
}

func WriteError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrPRNotFound):
		WriteErrorResponse(w, "NOT_FOUND", "PR not found", http.StatusNotFound)
	case errors.Is(err, domain.ErrPRExists):
		WriteErrorResponse(w, "PR_EXISTS", "PR id already exists", http.StatusConflict)
	case errors.Is(err, domain.ErrPRMerged):
		WriteErrorResponse(w, "PR_MERGED", "cannot reassign on merged PR", http.StatusConflict)
	case errors.Is(err, domain.ErrNotAssigned):
		WriteErrorResponse(w, "NOT_ASSIGNED", "reviewer is not assigned to this PR", http.StatusConflict)
	case errors.Is(err, domain.ErrNoCandidate):
		WriteErrorResponse(w, "NO_CANDIDATE", "no active replacement candidate in team", http.StatusConflict)
	case errors.Is(err, domain.ErrTeamNotFound):
		WriteErrorResponse(w, "NOT_FOUND", "team not found", http.StatusNotFound)
	case errors.Is(err, domain.ErrUserNotFound):
		WriteErrorResponse(w, "NOT_FOUND", "user not found", http.StatusNotFound)
	case errors.Is(err, domain.ErrTeamExists):
		WriteErrorResponse(w, "TEAM_EXISTS", "team_name already exists", http.StatusBadRequest)
	case errors.Is(err, domain.ErrAuthorTeamNotFound): 
		WriteErrorResponse(w, "NOT_FOUND", "author or team not found", http.StatusNotFound)
	default:
		WriteErrorResponse(w, "INTERNAL_ERROR", "internal server error", http.StatusInternalServerError)
	}
}
