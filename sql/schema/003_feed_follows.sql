-- +goose Up
create table feed_follows (
    id           UUID primary key,
    created_at   timestamp not null,
    updated_at   timestamp not null,
    unique(feed_id, user_id),
    feed_id     UUID not null references feeds(id) on delete cascade,
    user_id      UUID not null references users(id) on delete cascade
);

-- +goose Down
drop table feed_follows;