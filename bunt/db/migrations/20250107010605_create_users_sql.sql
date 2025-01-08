-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    id integer primary key,
    username text not null,
    user_id text unique not null,
    password_hash text not null,
    avatar_url text not null,
    background_url text not null,
    email text unique not null,
    bio text not null,
    created_at datetime not null,
    updated_at datetime not null
);
-- +goose StatementBegin
-- +goose StatementEnd
-- +goose Down
DROP TABLE IF EXISTS users;
-- +goose StatementBegin
-- +goose StatementEnd