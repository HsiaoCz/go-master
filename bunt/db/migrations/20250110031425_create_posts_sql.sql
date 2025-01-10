-- +goose Up
CREATE TABLE IF NOT EXISTS posts (
    id integer primary key autoincrement,
    post_id text unique not null,
    user_id text not null,
    title text not null,
    caption text not null,
    content text not null,
    image_url text,
    video_url text,
    location text not null,
    likes integer default 0,
    comments integer default 0,
    created_at datetime not null,
    updated_at datetime not null
);
-- +goose StatementBegin
-- +goose StatementEnd
-- +goose Down
DROP TABLE IF EXISTS posts;
-- +goose StatementBegin
-- +goose StatementEnd