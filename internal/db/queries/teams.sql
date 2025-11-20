-- name: CreateTeam :exec
INSERT INTO teams (team_name) VALUES ($1);

-- name: GetTeamMembers :many
SELECT u.user_id, u.username, u.is_active
FROM users u
WHERE u.team_name = $1;
