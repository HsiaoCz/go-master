-- +goose Up
CREATE TABLE IF NOT EXISTS message_status(
    id integer primary key autoincrement,
    message_id text not null,
    user_id text not null,
    is_read boolean not null,
    read_at datetime not null,
    created_at datetime not null,
    updated_at datetime not null
);
-- +goose StatementBegin
-- +goose StatementEnd
-- +goose Down
DROP TABLE IF EXISTS message_status;
-- +goose StatementBegin
-- +goose StatementEnd