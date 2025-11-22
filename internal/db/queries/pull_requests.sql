-- name: CreatePullRequest :exec
INSERT INTO pull_requests (pull_request_id, pull_request_name, author_id)
VALUES ($1, $2, $3);

-- name: GetPullRequest :one
SELECT * FROM pull_requests WHERE pull_request_id = $1;

-- name: MergePullRequest :one
UPDATE pull_requests
SET status = 'MERGED', merged_at = NOW()
WHERE pull_request_id = $1 AND status <> 'MERGED'
RETURNING *;

-- name: GetPullRequestShortByReviewer :many
SELECT pr.pull_request_id, pr.pull_request_name, pr.author_id, pr.status
FROM pull_requests pr
JOIN pull_request_reviewers ar USING(pull_request_id)
WHERE ar.reviewer_id = $1;
