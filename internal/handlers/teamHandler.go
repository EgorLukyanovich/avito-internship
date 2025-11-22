package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/egor_lukyanovich/avito/internal/db"
	models "github.com/egor_lukyanovich/avito/internal/model"
	json_resp "github.com/egor_lukyanovich/avito/pkg/response"
)

type TeamHandlers struct {
	q *db.Queries
}

func NewTeamHandlers(q *db.Queries) *TeamHandlers {
	return &TeamHandlers{q: q}
}

func (h *TeamHandlers) AddTeam(w http.ResponseWriter, r *http.Request) {
	var req models.Team
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json_resp.RespondError(w, 400, "BAD_REQUEST", "invalid json")
		return
	}

	if _, err := h.q.TeamExists(r.Context(), req.TeamName); err == nil {
		json_resp.RespondError(w,
			400, "TEAM_EXISTS",
			req.TeamName+" already exists",
		)
		return
	}

	if err := h.q.CreateTeam(r.Context(), req.TeamName); err != nil {
		json_resp.RespondError(w, 500, "INTERNAL", err.Error())
		return
	}

	for _, m := range req.Members {
		err := h.q.UpsertUser(r.Context(), db.UpsertUserParams{
			UserID:   m.UserID,
			Username: m.Username,
			TeamName: req.TeamName,
			IsActive: m.IsActive,
		})
		if err != nil {
			json_resp.RespondError(w, 500, "INTERNAL", err.Error())
			return
		}
	}

	json_resp.RespondJSON(w, 201, map[string]interface{}{
		"team": req,
	})
}

func (h *TeamHandlers) GetTeam(w http.ResponseWriter, r *http.Request) {
	teamName := r.URL.Query().Get("team_name")
	if teamName == "" {
		json_resp.RespondError(w, 400, "BAD_REQUEST", "team_name required")
		return
	}

	if _, err := h.q.TeamExists(r.Context(), teamName); err != nil {
		json_resp.RespondError(w, 404, "NOT_FOUND", "team not found")
		return
	}

	rows, err := h.q.GetTeamMembers(r.Context(), teamName)
	if err != nil {
		json_resp.RespondError(w, 500, "INTERNAL", err.Error())
		return
	}

	members := make([]models.TeamMember, 0, len(rows))
	for _, u := range rows {
		members = append(members, models.TeamMember{
			UserID:   u.UserID,
			Username: u.Username,
			IsActive: u.IsActive,
		})
	}

	json_resp.RespondJSON(w, 200, models.Team{
		TeamName: teamName,
		Members:  members,
	})
}
