-- +goose Up
CREATE TABLE IF NOT EXISTS users(
    id integer primary key autoincrement,
    user_id text unqiue not null,
    username text not null,
    password text not null,
    email text unqiue not null,
    avatar text not null,
    Bio text not null,
    background_url text not null,
    created_at datetime not null,
    updated_at datetime not null
);
-- +goose StatementBegin
-- +goose StatementEnd
-- +goose Down
DROP TABLE IF EXISTS users;
-- +goose StatementBegin
-- +goose StatementEnd