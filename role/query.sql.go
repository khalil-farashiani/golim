// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: query.sql

package role

import (
	"context"
)

const crateRateLimiter = `-- name: CrateRateLimiter :one
INSERT INTO rate_limiter (
    name
) VALUES (
             ?
         )
    RETURNING id, name, created_at, updated_at, deleted_at
`

func (q *Queries) CrateRateLimiter(ctx context.Context, name string) (RateLimiter, error) {
	row := q.db.QueryRowContext(ctx, crateRateLimiter, name)
	var i RateLimiter
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const createRole = `-- name: CreateRole :one
INSERT INTO role (
        endpoint,
        operation,
        bucket_size,
        add_token_per_min,
        initial_tokens,
        rate_limiter_id
    ) VALUES (
             ?, ?, ?, ?, ?, ?
         )
    RETURNING id, endpoint, operation, bucket_size, created_at, add_token_per_min, initial_tokens, updated_at, deleted_at, rate_limiter_id
`

type CreateRoleParams struct {
	Endpoint       string
	Operation      string
	BucketSize     int64
	AddTokenPerMin int64
	InitialTokens  int64
	RateLimiterID  int64
}

func (q *Queries) CreateRole(ctx context.Context, arg CreateRoleParams) (Role, error) {
	row := q.db.QueryRowContext(ctx, createRole,
		arg.Endpoint,
		arg.Operation,
		arg.BucketSize,
		arg.AddTokenPerMin,
		arg.InitialTokens,
		arg.RateLimiterID,
	)
	var i Role
	err := row.Scan(
		&i.ID,
		&i.Endpoint,
		&i.Operation,
		&i.BucketSize,
		&i.CreatedAt,
		&i.AddTokenPerMin,
		&i.InitialTokens,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.RateLimiterID,
	)
	return i, err
}

const getRateLimiters = `-- name: GetRateLimiters :many
SELECT id, name, created_at, updated_at, deleted_at FROM rate_limiter
WHERE deleted_at = null
`

func (q *Queries) GetRateLimiters(ctx context.Context) ([]RateLimiter, error) {
	rows, err := q.db.QueryContext(ctx, getRateLimiters)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []RateLimiter
	for rows.Next() {
		var i RateLimiter
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
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

const getRoles = `-- name: GetRoles :many
SELECT id, endpoint, operation, bucket_size, created_at, add_token_per_min, initial_tokens, updated_at, deleted_at, rate_limiter_id FROM role
WHERE deleted_at = null
`

func (q *Queries) GetRoles(ctx context.Context) ([]Role, error) {
	rows, err := q.db.QueryContext(ctx, getRoles)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Role
	for rows.Next() {
		var i Role
		if err := rows.Scan(
			&i.ID,
			&i.Endpoint,
			&i.Operation,
			&i.BucketSize,
			&i.CreatedAt,
			&i.AddTokenPerMin,
			&i.InitialTokens,
			&i.UpdatedAt,
			&i.DeletedAt,
			&i.RateLimiterID,
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
