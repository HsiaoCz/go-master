-- +goose Up
CREATE TABLE IF NOT EXISTS stocks(
    id integer primary key autoincrement,
    stock_id text unique not null,
);
-- +goose StatementBegin
-- +goose StatementEnd
-- +goose Down
DROP TABLE IF EXISTS stocks;
-- +goose StatementBegin
-- +goose StatementEnd
