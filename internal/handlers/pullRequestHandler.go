package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/egor_lukyanovich/avito/internal/db"
	"github.com/egor_lukyanovich/avito/internal/models"
	json_resp "github.com/egor_lukyanovich/avito/pkg/response"
)

type PullRequestHandlers struct {
	q *db.Queries
}

func NewPullRequestHandlers(q *db.Queries) *PullRequestHandlers {
	return &PullRequestHandlers{q: q}
}

func (p *PullRequestHandlers) CreatePullRequest(w http.ResponseWriter, r *http.Request) {
	var req models.PullRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json_resp.RespondError(w, 400, "BAD_REQUEST", "invalid json")
		return
	}

	params := db.CreatePullRequestParams{
		PullRequestID:   req.PullRequestID,
		PullRequestName: req.PullRequestName,
		AuthorID:        req.AuthorID,
	}

	if _, err := p.q.PullRequestExists(r.Context(), req.PullRequestID); err == nil {
		json_resp.RespondError(w,
			409, "PR_EXISTS",
			req.PullRequestID+" already exists",
		)
		return
	}

	author, err := p.q.GetUser(r.Context(), req.AuthorID)
	if err != nil {
		json_resp.RespondError(w, 404, "NOT_FOUND", "author not found")
		return
	}
	if author.TeamName == "" {
		json_resp.RespondError(w, 404, "NOT_FOUND", "author has no team")
		return
	}

	if err := p.q.CreatePullRequest(r.Context(), params); err != nil {
		json_resp.RespondError(w, 500, "INTERNAL_ERROR", err.Error())
	}

	candidates, _ := p.q.GetActiveTeamReviewCandidates(r.Context(), db.GetActiveTeamReviewCandidatesParams{
		TeamName: author.TeamName,
		UserID:   author.UserID,
	})
	assigned := []string{}
	for i := 0; i < len(candidates) && len(assigned) < 2; i++ {
		if candidates[i] == req.AuthorID {
			continue
		}
		err = p.q.AssignReviewer(r.Context(), db.AssignReviewerParams{
			PullRequestID: req.PullRequestID,
			ReviewerID:    candidates[i],
		})
		if err != nil {
			json_resp.RespondError(w, 500, "INTERNAL_ERROR", "failed to assign reviewer: "+err.Error())
			return
		}
		assigned = append(assigned, candidates[i])
	}

	req.Status = string(db.PrStatusOPEN)
	req.AssignedReviewers = assigned

	json_resp.RespondJSON(w, 201, map[string]interface{}{
		"pullRequest": req,
	})
}
