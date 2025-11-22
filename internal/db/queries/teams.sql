-- name: CreateTeam :exec
INSERT INTO teams (team_name) VALUES ($1);

-- name: GetTeamMembers :many
SELECT u.user_id, u.username, u.is_active
FROM users u
WHERE u.team_name = $1
ORDER BY u.user_id;

-- name: TeamExists :one
SELECT team_name FROM teams WHERE team_name = $1;