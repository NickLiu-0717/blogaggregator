-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1
)
RETURNING *;

-- name: GetUserFromName :one
SELECT * FROM users
where name = $1;

-- name: GetUsers :many
SELECT name from users;

-- name: GetUserNameFromID :one
select name from users
where id = $1;

-- name: DeleteAllUsers :exec
Delete from users;