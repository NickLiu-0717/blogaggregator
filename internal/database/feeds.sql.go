// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: feeds.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const createFeed = `-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2,
    $3
)
RETURNING id, created_at, updated_at, name, url, user_id, last_fetch_at
`

type CreateFeedParams struct {
	Name   string
	Url    string
	UserID uuid.UUID
}

func (q *Queries) CreateFeed(ctx context.Context, arg CreateFeedParams) (Feed, error) {
	row := q.db.QueryRowContext(ctx, createFeed, arg.Name, arg.Url, arg.UserID)
	var i Feed
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Url,
		&i.UserID,
		&i.LastFetchAt,
	)
	return i, err
}

const deleteAllFeeds = `-- name: DeleteAllFeeds :exec
Delete from feeds
`

func (q *Queries) DeleteAllFeeds(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, deleteAllFeeds)
	return err
}

const getFeed = `-- name: GetFeed :one
SELECT id, created_at, updated_at, name, url, user_id, last_fetch_at FROM feeds
where name = $1
`

func (q *Queries) GetFeed(ctx context.Context, name string) (Feed, error) {
	row := q.db.QueryRowContext(ctx, getFeed, name)
	var i Feed
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Url,
		&i.UserID,
		&i.LastFetchAt,
	)
	return i, err
}

const getFeedIDandNameFromURL = `-- name: GetFeedIDandNameFromURL :one
select id, name from feeds
where url = $1
`

type GetFeedIDandNameFromURLRow struct {
	ID   uuid.UUID
	Name string
}

func (q *Queries) GetFeedIDandNameFromURL(ctx context.Context, url string) (GetFeedIDandNameFromURLRow, error) {
	row := q.db.QueryRowContext(ctx, getFeedIDandNameFromURL, url)
	var i GetFeedIDandNameFromURLRow
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const getFeedIdandUserID = `-- name: GetFeedIdandUserID :many
select id, user_id
from feeds
`

type GetFeedIdandUserIDRow struct {
	ID     uuid.UUID
	UserID uuid.UUID
}

func (q *Queries) GetFeedIdandUserID(ctx context.Context) ([]GetFeedIdandUserIDRow, error) {
	rows, err := q.db.QueryContext(ctx, getFeedIdandUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetFeedIdandUserIDRow
	for rows.Next() {
		var i GetFeedIdandUserIDRow
		if err := rows.Scan(&i.ID, &i.UserID); err != nil {
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

const getFeeds = `-- name: GetFeeds :many
SELECT id, created_at, updated_at, name, url, user_id, last_fetch_at from feeds
`

func (q *Queries) GetFeeds(ctx context.Context) ([]Feed, error) {
	rows, err := q.db.QueryContext(ctx, getFeeds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Feed
	for rows.Next() {
		var i Feed
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Name,
			&i.Url,
			&i.UserID,
			&i.LastFetchAt,
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

const getNextFeedToFetch = `-- name: GetNextFeedToFetch :one
select id, url
from feeds
order by last_fetch_at asc nulls first
limit 1
`

type GetNextFeedToFetchRow struct {
	ID  uuid.UUID
	Url string
}

func (q *Queries) GetNextFeedToFetch(ctx context.Context) (GetNextFeedToFetchRow, error) {
	row := q.db.QueryRowContext(ctx, getNextFeedToFetch)
	var i GetNextFeedToFetchRow
	err := row.Scan(&i.ID, &i.Url)
	return i, err
}

const markFeedFetched = `-- name: MarkFeedFetched :exec
update feeds
set last_fetch_at = NOW(), updated_at = NOW()
where id = $1
`

func (q *Queries) MarkFeedFetched(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, markFeedFetched, id)
	return err
}
