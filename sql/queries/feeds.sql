-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2,
    $3
)
RETURNING *;

-- name: GetFeed :one
SELECT * FROM feeds
where name = $1;

-- name: GetFeedIdandUserID :many
select id, user_id
from feeds;

-- name: GetFeedIDandNameFromURL :one
select id, name from feeds
where url = $1;

-- name: GetFeeds :many
SELECT * from feeds;

-- name: MarkFeedFetched :exec
update feeds
set last_fetch_at = NOW(), updated_at = NOW()
where id = $1;

-- name: GetNextFeedToFetch :one
select id, url
from feeds
order by last_fetch_at asc nulls first
limit 1;

-- name: DeleteAllFeeds :exec
Delete from feeds;