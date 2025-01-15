-- +goose Up
CREATE TABLE IF NOT EXISTS user_relationships(
    id integer primary key autoincrement,
    user_id text not null,
    related_user_id text not null,
    relationship_type text not null,
    created_at datetime not null,
    updated_at datetime not null
);
-- +goose StatementBegin
-- +goose StatementEnd
-- +goose Down
DROP TABLE IF EXISTS user_relationships;
-- +goose StatementBegin
-- +goose StatementEnd