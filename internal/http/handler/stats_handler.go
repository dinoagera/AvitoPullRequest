package handler

import (
	"log/slog"
	"net/http"
)

type StatsHandler struct {
	log     *slog.Logger
	service StatsService
}

func NewStatsHandler(log *slog.Logger, service StatsService) *StatsHandler {
	return &StatsHandler{
		log:     log,
		service: service,
	}
}

func (sh *StatsHandler) GetReviewerStats(w http.ResponseWriter, r *http.Request) {
	stats, err := sh.service.GetReviewerStats(r.Context())
	if err != nil {
		sh.log.Info("failed to get reviewer stats", "err", err)
		WriteError(w, err)
		return
	}
	WriteJSONResponse(w, stats, http.StatusOK)
}
