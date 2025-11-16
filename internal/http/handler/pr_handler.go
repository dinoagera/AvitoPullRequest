package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type PRHandler struct {
	log     *slog.Logger
	service PRService
}

func NewPRHandler(log *slog.Logger, service PRService) *PRHandler {
	return &PRHandler{
		log:     log,
		service: service,
	}
}
func (pr *PRHandler) CreatePR(w http.ResponseWriter, r *http.Request) {
	var req CreateRPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteErrorResponse(w, "BAD_REQUEST", "Invalid JSON", http.StatusBadRequest)
		return
	}
	if req.ID == "" || req.Name == "" || req.AuthorID == "" {
		WriteErrorResponse(w, "BAD_REQUEST", "pull_request_id, pull_request_name and author_id are required", http.StatusBadRequest)
		return
	}
	pullRequest := dtoToDomainPR(req)
	result, err := pr.service.CreatePR(r.Context(), pullRequest)
	if err != nil {
		pr.log.Info("failed to create pr", "err", err)
		WriteError(w, err)
		return
	}
	WriteJSONResponse(w, map[string]interface{}{"pr": result}, http.StatusCreated)
}
func (pr *PRHandler) MergePR(w http.ResponseWriter, r *http.Request) {
	var req MergePRRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteErrorResponse(w, "BAD_REQUEST", "Invalid JSON", http.StatusBadRequest)
		return
	}
	if req.ID == "" {
		WriteErrorResponse(w, "BAD_REQUEST", "pull_request_id is required", http.StatusBadRequest)
		return
	}
	result, err := pr.service.MergePR(r.Context(), req.ID)
	if err != nil {
		pr.log.Info("failed to merge pr", "err", err)
		WriteError(w, err)
		return
	}
	WriteJSONResponse(w, map[string]interface{}{"pr": result}, http.StatusOK)
}
func (pr *PRHandler) ReassignReviewer(w http.ResponseWriter, r *http.Request) {
	var req ReassignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteErrorResponse(w, "BAD_REQUEST", "Invalid JSON", http.StatusBadRequest)
		return
	}
	if req.PRID == "" || req.OldUserID == "" {
		WriteErrorResponse(w, "BAD_REQUEST", "pull_request_id and old_user_id are required", http.StatusBadRequest)
		return
	}
	result, newUserID, err := pr.service.ReassignReviewer(r.Context(), req.PRID, req.OldUserID)
	if err != nil {
		pr.log.Info("failed to reassign reviewer", "err", err)
		WriteError(w, err)
		return
	}
	WriteJSONResponse(w, map[string]interface{}{
		"pr":          result,
		"replaced_by": newUserID,
	}, http.StatusOK)
}
