package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/egor_lukyanovich/avito/internal/db"
	"github.com/egor_lukyanovich/avito/internal/models"
	json_resp "github.com/egor_lukyanovich/avito/pkg/response"
)

type UserHandlers struct {
	q *db.Queries
}

func NewUserHandlers(q *db.Queries) *UserHandlers {
	return &UserHandlers{q: q}
}

type SetIsActiveRequest struct {
	UserID   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

func (h *UserHandlers) UpsertUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		json_resp.RespondError(w, 400, "BAD_REQUEST", "invalid request body")
		return
	}

	params := db.UpsertUserParams{
		UserID:   user.UserID,
		Username: user.Username,
		TeamName: user.TeamName,
		IsActive: user.IsActive,
	}

	if err := h.q.UpsertUser(r.Context(), params); err != nil {
		json_resp.RespondError(w, 500, "INTERNAL_ERROR", err.Error())
		return
	}

	json_resp.RespondJSON(w, 201, map[string]interface{}{
		"user": user,
	})
}

func (h *UserHandlers) SetUserActive(w http.ResponseWriter, r *http.Request) {
	var req SetIsActiveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json_resp.RespondError(w, 400, "BAD_REQUEST", "invalid request body")
		return
	}

	params := db.SetUserActiveParams{
		UserID:   req.UserID,
		IsActive: req.IsActive,
	}

	user, err := h.q.SetUserActive(r.Context(), params)
	if err != nil {
		if err == sql.ErrNoRows {
			json_resp.RespondError(w, 404, "NOT_FOUND", "user not found")
			return
		}
		json_resp.RespondError(w, 500, "INTERNAL_ERROR", err.Error())
		return
	}

	json_resp.RespondJSON(w, 200, map[string]interface{}{
		"user": user,
	})
}

func (h *UserHandlers) GetUser(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		json_resp.RespondError(w, 400, "BAD_REQUEST", "user_id query param required")
		return
	}

	user, err := h.q.GetUser(r.Context(), userID)
	if err != nil {
		if err == sql.ErrNoRows {
			json_resp.RespondError(w, 404, "NOT_FOUND", "user not found")
			return
		}
		json_resp.RespondError(w, 500, "INTERNAL_ERROR", err.Error())
		return
	}

	json_resp.RespondJSON(w, 200, map[string]interface{}{
		"user": user,
	})
}

func (h *UserHandlers) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		json_resp.RespondError(w, 400, "BAD_REQUEST", "user_id query param required")
		return
	}

	if err := h.q.DeleteUser(r.Context(), userID); err != nil {
		json_resp.RespondError(w, 500, "INTERNAL_ERROR", err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandlers) GetPRsForReview(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		json_resp.RespondError(w, 400, "BAD_REQUEST", "user_id query param required")
		return
	}

	prs, err := h.q.GetPullRequestShortByReviewer(r.Context(), userID)
	if err != nil {
		json_resp.RespondError(w, 500, "INTERNAL_ERROR", err.Error())
		return
	}

	result := make([]models.PullRequestShort, 0, len(prs))
	for _, pr := range prs {
		result = append(result, models.PullRequestShort{
			PullRequestID:   pr.PullRequestID,
			PullRequestName: pr.PullRequestName,
			AuthorID:        pr.AuthorID,
			Status:          string(pr.Status),
		})
	}

	json_resp.RespondJSON(w, 200, map[string]interface{}{
		"user_id":       userID,
		"pull_requests": result,
	})
}
