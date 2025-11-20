-- +goose Up
CREATE TABLE pull_request_reviewers (
    pull_request_id UUID REFERENCES pull_requests(pull_request_id) ON DELETE CASCADE,
    reviewer_id UUID REFERENCES users(user_id),
    PRIMARY KEY (pull_request_id, reviewer_id)
);

-- +goose Down
DROP TABLE pull_request_reviewers;
