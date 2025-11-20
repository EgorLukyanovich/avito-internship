-- +goose Up
CREATE TABLE users (
    user_id UUID PRIMARY KEY,
    username VARCHAR NOT NULL UNIQUE,
    team_name VARCHAR NOT NULL REFERENCES teams(team_name) ON DELETE CASCADE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE
);

-- +goose Down
DROP TABLE users;
