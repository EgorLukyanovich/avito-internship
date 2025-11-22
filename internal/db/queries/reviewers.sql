-- name: AssignReviewer :exec
INSERT INTO pull_request_reviewers (pull_request_id, reviewer_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING;

-- name: DeleteReviewer :exec
DELETE FROM pull_request_reviewers
WHERE pull_request_id = $1 AND reviewer_id = $2;

-- name: GetReviewers :many
SELECT reviewer_id FROM pull_request_reviewers
WHERE pull_request_id = $1
ORDER BY reviewer_id;

-- name: IsReviewerAssigned :one
SELECT reviewer_id FROM pull_request_reviewers WHERE pull_request_id = $1 AND reviewer_id = $2;