package handlers

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"slices"
	"time"

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
		CreatedAt:       time.Now(),
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

	prCreated, err := p.q.GetPullRequest(r.Context(), req.PullRequestID)
	if err != nil {
		json_resp.RespondError(w, 500, "INTERNAL_ERROR", "failed to fetch created PR")
		return
	}

	candidates, _ := p.q.GetActiveTeamReviewCandidates(r.Context(), db.GetActiveTeamReviewCandidatesParams{
		TeamName: author.TeamName,
		UserID:   author.UserID,
	})

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	rng.Shuffle(len(candidates), func(i, j int) {
		candidates[i], candidates[j] = candidates[j], candidates[i]
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

	resp := models.PullRequest{
		PullRequestID:     prCreated.PullRequestID,
		PullRequestName:   prCreated.PullRequestName,
		AuthorID:          prCreated.AuthorID,
		Status:            string(prCreated.Status),
		CreatedAt:         &prCreated.CreatedAt,
		AssignedReviewers: assigned,
	}

	json_resp.RespondJSON(w, 201, map[string]interface{}{
		"pullRequest": resp,
	})
}

func (p *PullRequestHandlers) MergePullRequest(w http.ResponseWriter, r *http.Request) {
	var req struct {
		PullRequestID string `json:"pull_request_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json_resp.RespondError(w, 400, "BAD_REQUEST", "invalid json")
		return
	}

	pr, err := p.q.GetPullRequest(r.Context(), req.PullRequestID)
	if err != nil {
		json_resp.RespondError(w, 404, "NOT_FOUND", "PR not found")
		return
	}

	if pr.Status != db.PrStatusMERGED {
		pr, err = p.q.MergePullRequest(r.Context(), req.PullRequestID)
		if err != nil {
			json_resp.RespondError(w, 500, "INTERNAL_ERROR", err.Error())
			return
		}
	}

	revs, _ := p.q.GetReviewers(r.Context(), pr.PullRequestID)
	assigned := make([]string, len(revs))
	copy(assigned, revs)

	json_resp.RespondJSON(w, 200, map[string]interface{}{
		"pullRequest": models.PullRequest{
			PullRequestID:     pr.PullRequestID,
			PullRequestName:   pr.PullRequestName,
			AuthorID:          pr.AuthorID,
			Status:            string(pr.Status),
			AssignedReviewers: assigned,
			MergedAt:          &pr.MergedAt.Time,
		},
	})
}

func (p *PullRequestHandlers) ReassignReviewer(w http.ResponseWriter, r *http.Request) {
	var req struct {
		PullRequestID string `json:"pull_request_id"`
		OldUserID     string `json:"old_user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json_resp.RespondError(w, 400, "BAD_REQUEST", "invalid json")
		return
	}

	pr, err := p.q.GetPullRequest(r.Context(), req.PullRequestID)
	if err != nil {
		json_resp.RespondError(w, 404, "NOT_FOUND", "PR not found")
		return
	}

	if pr.Status == db.PrStatusMERGED {
		json_resp.RespondError(w, 409, "PR_MERGED", "cannot reassign on merged PR")
		return
	}

	revs, err := p.q.GetReviewers(r.Context(), pr.PullRequestID)
	if err != nil {
		json_resp.RespondError(w, 500, "INTERNAL_ERROR", "failed to fetch reviewers")
		return
	}

	if !slices.Contains(revs, req.OldUserID) {
		json_resp.RespondError(w, 409, "NOT_ASSIGNED", "reviewer is not assigned to this PR")
		return
	}

	candidates, err := p.q.GetActiveReplacementCandidates(r.Context(), pr.AuthorID)
	if err != nil {
		json_resp.RespondError(w, 500, "INTERNAL_ERROR", "failed to fetch candidates")
		return
	}

	var available []string
	for _, c := range candidates {
		if c != pr.AuthorID && !slices.Contains(revs, c) {
			available = append(available, c)
		}
	}

	if len(available) == 0 {
		json_resp.RespondError(w, 409, "NO_CANDIDATE", "no available reviewers for reassignment")
		return
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	newReviewer := available[rng.Intn(len(available))]

	err = p.q.DeleteReviewer(r.Context(), db.DeleteReviewerParams{
		PullRequestID: req.PullRequestID,
		ReviewerID:    req.OldUserID,
	})
	if err != nil {
		json_resp.RespondError(w, 500, "INTERNAL_ERROR", "failed to delete old reviewer")
		return
	}

	err = p.q.AssignReviewer(r.Context(), db.AssignReviewerParams{
		PullRequestID: req.PullRequestID,
		ReviewerID:    newReviewer,
	})
	if err != nil {
		json_resp.RespondError(w, 500, "INTERNAL_ERROR", "failed to assign reviewer")
		return
	}

	revs, _ = p.q.GetReviewers(r.Context(), pr.PullRequestID)

	json_resp.RespondJSON(w, 200, map[string]interface{}{
		"pr": models.PullRequest{
			PullRequestID:     pr.PullRequestID,
			PullRequestName:   pr.PullRequestName,
			AuthorID:          pr.AuthorID,
			Status:            string(pr.Status),
			AssignedReviewers: revs,
		},
		"replaced_by": newReviewer,
	})
}
