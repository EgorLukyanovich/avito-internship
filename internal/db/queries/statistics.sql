-- name: GetReviewStatsByUser :many
SELECT reviewer_id, COUNT(*) AS reviews_count
FROM pull_request_reviewers
GROUP BY reviewer_id;

-- name: GetPRsWithoutReview :one
SELECT COUNT(*) AS no_review_count
FROM pull_requests pr
LEFT JOIN pull_request_reviewers r 
    ON pr.pull_request_id = r.pull_request_id
WHERE r.reviewer_id IS NULL
  AND pr.status = 'OPEN';
