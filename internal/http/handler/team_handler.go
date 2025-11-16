package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type TeamHandler struct {
	log     *slog.Logger
	service TeamService
}

func NewTeamHandler(log *slog.Logger, service TeamService) *TeamHandler {
	return &TeamHandler{
		log:     log,
		service: service,
	}
}
func (th *TeamHandler) AddTeam(w http.ResponseWriter, r *http.Request) {
	var req AddTeamRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteErrorResponse(w, "BAD_REQUEST", "invalid JSON", http.StatusBadRequest)
		return
	}
	if req.TeamName == "" {
		WriteErrorResponse(w, "BAD_REQUEST", "team_name is required", http.StatusBadRequest)
		return
	}
	if len(req.Members) == 0 {
		WriteErrorResponse(w, "BAD_REQUEST", "members are required", http.StatusBadRequest)
		return
	}
	team := dtoToDomainTeam(req)
	result, err := th.service.AddTeam(r.Context(), team)
	if err != nil {
		th.log.Info("failed to add team", "err", err)
		WriteError(w, err)
		return
	}
	res := domainTodtoTeam(result)
	WriteJSONResponse(w, map[string]interface{}{"team": res}, http.StatusCreated)
}

func (th *TeamHandler) GetTeam(w http.ResponseWriter, r *http.Request) {
	teamName := r.URL.Query().Get("team_name")
	if teamName == "" {
		WriteErrorResponse(w, "BAD_REQUEST", "team_name is required", http.StatusBadRequest)
		return
	}
	team, err := th.service.GetTeam(r.Context(), teamName)
	if err != nil {
		th.log.Info("failed to get team", "err", err)
		WriteError(w, err)
		return
	}
	res := domainTodtoTeam(team)
	WriteJSONResponse(w, res, http.StatusOK)
}
