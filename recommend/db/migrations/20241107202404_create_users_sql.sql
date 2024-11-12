-- +goose Up
CREATE TABLE IF NOT EXISTS users(
    id integer primary key,
    user_id text unique not null,
    username text not null,
    phone text not null,
    hash_password text not null,
    role boolean not null,
    avatar text not null,
    brief text not null,
    birthday text not null,
    background_image text not null,
    gender text not null,
    created_at datetime not null,
    updated_at datetime not null,
    deleted_at datetime
);
-- +goose StatementBegin
-- +goose StatementEnd
-- +goose Down
DROP TABLE IF EXISTS users;
-- +goose StatementBegin
-- +goose StatementEnd