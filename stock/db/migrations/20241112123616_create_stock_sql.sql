-- +goose Up
CREATE TABLE IF NOT EXISTS stocks(
    id integer primary key,
);
-- +goose StatementBegin
-- +goose StatementEnd
-- +goose Down
DROP TABLE IF EXISTS stocks;
-- +goose StatementBegin
-- +goose StatementEnd
