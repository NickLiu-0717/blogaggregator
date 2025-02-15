// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: posts.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createPost = `-- name: CreatePost :one
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
RETURNING id, created_at, updated_at, title, url, description, published_at, feed_id
`

type CreatePostParams struct {
	Title       sql.NullString
	Url         string
	Description string
	PublishedAt time.Time
	FeedID      uuid.UUID
}

func (q *Queries) CreatePost(ctx context.Context, arg CreatePostParams) (Post, error) {
	row := q.db.QueryRowContext(ctx, createPost,
		arg.Title,
		arg.Url,
		arg.Description,
		arg.PublishedAt,
		arg.FeedID,
	)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Title,
		&i.Url,
		&i.Description,
		&i.PublishedAt,
		&i.FeedID,
	)
	return i, err
}

const getPostsforUser = `-- name: GetPostsforUser :many
SELECT p.title, p.url, p.description, p.published_at, p.feed_id
from posts as p
join feed_follows as f
on p.feed_id = f.feed_id
where f.user_id = $1
order by p.published_at desc
limit $2
`

type GetPostsforUserParams struct {
	UserID uuid.UUID
	Limit  int32
}

type GetPostsforUserRow struct {
	Title       sql.NullString
	Url         string
	Description string
	PublishedAt time.Time
	FeedID      uuid.UUID
}

func (q *Queries) GetPostsforUser(ctx context.Context, arg GetPostsforUserParams) ([]GetPostsforUserRow, error) {
	rows, err := q.db.QueryContext(ctx, getPostsforUser, arg.UserID, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetPostsforUserRow
	for rows.Next() {
		var i GetPostsforUserRow
		if err := rows.Scan(
			&i.Title,
			&i.Url,
			&i.Description,
			&i.PublishedAt,
			&i.FeedID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
