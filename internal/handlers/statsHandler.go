package handlers

import (
	"net/http"

	"github.com/egor_lukyanovich/avito/internal/db"
	json_resp "github.com/egor_lukyanovich/avito/pkg/response"
)

type StatsResponse struct {
	ReviewsByUser map[string]int64 `json:"reviews_by_user"`
	NoReviewCount int64            `json:"no_review_count"`
}

type StatsHandler struct {
	q *db.Queries
}

func NewStatsHandler(q *db.Queries) *StatsHandler {
	return &StatsHandler{q: q}
}

func (h *StatsHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	users, err := h.q.GetReviewStatsByUser(ctx)
	if err != nil {
		json_resp.RespondError(w, 500, "INTERNAL_ERROR", err.Error())
		return
	}

	noReview, err := h.q.GetPRsWithoutReview(ctx)
	if err != nil {
		json_resp.RespondError(w, 500, "INTERNAL_ERROR", err.Error())
		return
	}

	stats := StatsResponse{
		ReviewsByUser: make(map[string]int64),
		NoReviewCount: noReview,
	}

	for _, u := range users {
		stats.ReviewsByUser[u.ReviewerID] = u.ReviewsCount
	}

	json_resp.RespondJSON(w, 200, map[string]interface{}{
		"stats": stats,
	})
}
