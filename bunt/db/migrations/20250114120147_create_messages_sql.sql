-- +goose Up
CREATE TABLE IF NOT EXISTS messages(
    id integer primary key autoincrement,
    message_id text unique not null,
    sender_id text not null,
    recevier_id text not null,
    content text not null,
    type text not null,
    is_read boolean not null,
    created_at datetime not null,
);
-- +goose StatementBegin
-- +goose StatementEnd
-- +goose Down
DROP TABLE IF EXISTS messages;
-- +goose StatementBegin
-- +goose StatementEnd