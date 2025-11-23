-- name: UpsertUser :exec
INSERT INTO users (user_id, username, team_name, is_active)
VALUES ($1, $2, $3, $4)
ON CONFLICT (user_id) DO UPDATE
SET username = excluded.username,
    team_name = excluded.team_name,
    is_active = excluded.is_active;

-- name: SetUserActive :one
UPDATE users
SET is_active = $2
WHERE user_id = $1
RETURNING *;

-- name: GetUser :one
SELECT * FROM users WHERE user_id = $1;

-- name: DeleteUser :exec
DELETE FROM users WHERE user_id = $1;

-- name: GetActiveTeamReviewCandidates :many
SELECT u.user_id
FROM users u
WHERE u.team_name = $1
  AND u.is_active = TRUE
  AND u.user_id <> $2;

-- name: GetActiveReplacementCandidates :many
SELECT u.user_id
FROM users u
WHERE u.team_name = (
    SELECT inner_u.team_name FROM users inner_u WHERE inner_u.user_id = $1
)
AND u.is_active = TRUE
AND u.user_id <> $1;
