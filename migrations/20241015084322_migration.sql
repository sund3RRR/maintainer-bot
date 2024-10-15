-- +goose Up
-- +goose StatementBegin
CREATE TABLE repos(
    id SERIAL PRIMARY KEY,
	chat_id INTEGER,
    host TEXT,
    owner TEXT,
    repo TEXT,
    last_tag TEXT,
    is_release BOOLEAN
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE repos;
-- +goose StatementEnd
