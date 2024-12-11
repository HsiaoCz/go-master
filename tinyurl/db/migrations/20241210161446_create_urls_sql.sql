-- +goose Up
CREATE TABLE IF NOT EXISTS urls(
    id integer primary key,
    original_url text not null,
    short_code text unique not null,
    is_custom boolean not null ,
    expired_at datetime not null,
    created_at datetime not null,
);
-- +goose StatementBegin
-- +goose StatementEnd
-- +goose Down
DROP TABLE IF EXISTS urls;
-- +goose StatementBegin
-- +goose StatementEnd
