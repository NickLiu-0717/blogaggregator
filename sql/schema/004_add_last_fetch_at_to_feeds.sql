-- +goose Up
alter table feeds
add column last_fetch_at timestamp;

-- +goose Down
drop table feeds;