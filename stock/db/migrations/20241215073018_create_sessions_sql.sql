-- +goose Up
CREATE TABLE IF NOT EXISTS sessions(
    id integer primary key,
);
-- +goose StatementBegin
-- +goose StatementEnd
-- +goose Down
DROP TABLE IF EXISTS sessions;
-- +goose StatementBegin
-- +goose StatementEnd
