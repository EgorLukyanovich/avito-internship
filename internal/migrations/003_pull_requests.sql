-- +goose Up
CREATE TYPE pr_status AS ENUM ('OPEN', 'MERGED');

CREATE TABLE pull_requests (
    pull_request_id VARCHAR PRIMARY KEY,
    pull_request_name VARCHAR NOT NULL,
    author_id VARCHAR NOT NULL REFERENCES users(user_id),
    status pr_status NOT NULL DEFAULT 'OPEN',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    merged_at TIMESTAMP WITH TIME ZONE
);

-- +goose Down
DROP TABLE pull_requests;
DROP TYPE pr_status;
