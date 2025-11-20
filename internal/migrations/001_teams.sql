-- +goose Up
CREATE TABLE teams (
    team_name VARCHAR PRIMARY KEY
);

-- +goose Down
DROP TABLE teams;
