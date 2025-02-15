-- name: CreateFeedFollow :one
with inserted_feed_follow as (
    INSERT INTO feed_follows (id, created_at, updated_at, feed_id, user_id)
    VALUES (
        gen_random_uuid(),
        NOW(),
        NOW(),
        $1,
        $2
    )
    RETURNING *
)
SELECT 
    iff.id, 
    iff.created_at, 
    iff.updated_at, 
    iff.feed_id, 
    iff.user_id, 
    f.name AS feed_name, 
    u.name AS user_name
FROM inserted_feed_follow AS iff
LEFT JOIN feeds AS f ON iff.feed_id = f.id
LEFT JOIN users AS u ON iff.user_id = u.id;

-- name: GetFeedFollowsForUser :many
Select
    f.name as feed_name
from feed_follows as ff
left join feeds as f on ff.feed_id = f.id
where ff.user_id = $1;

-- name: DeleteFollowFromURLandUser :exec
Delete from feed_follows as ff
using feeds as f
where
ff.feed_id = f.id 
and f.url = $1
and ff.user_id = $2;


-- name: DeleteFeedFollow :exec
Delete from feed_follows;