-- name: CreatePost :one
INSERT INTO Posts (id, created_at, updated_at, title, url, description, published_at, feed_id)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *;

-- name: GetPostsforUser :many
SELECT p.title, p.url, p.description, p.published_at, p.feed_id
from posts as p
join feed_follows as f
on p.feed_id = f.feed_id
where f.user_id = $1
order by p.published_at desc
limit $2;