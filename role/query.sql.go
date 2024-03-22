// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: query.sql

package role

import (
	"context"
	"database/sql"
)

const crateRateLimiter = `-- name: CrateRateLimiter :one
INSERT INTO rate_limiter (
    name,
    destination
) VALUES (
             ?,
             ?
         )
    RETURNING id, name, destination, deleted_at
`

type CrateRateLimiterParams struct {
	Name        string
	Destination string
}

func (q *Queries) CrateRateLimiter(ctx context.Context, arg CrateRateLimiterParams) (RateLimiter, error) {
	row := q.db.QueryRowContext(ctx, crateRateLimiter, arg.Name, arg.Destination)
	var i RateLimiter
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Destination,
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
    RETURNING id, endpoint, operation, bucket_size, add_token_per_min, initial_tokens, deleted_at, rate_limiter_id
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
		&i.AddTokenPerMin,
		&i.InitialTokens,
		&i.DeletedAt,
		&i.RateLimiterID,
	)
	return i, err
}

const deleteRateLimiter = `-- name: DeleteRateLimiter :exec
UPDATE rate_limiter
    SET deleted_at = CURRENT_TIMESTAMP
    WHERE id = ?
`

func (q *Queries) DeleteRateLimiter(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteRateLimiter, id)
	return err
}

const deleteRole = `-- name: DeleteRole :exec
UPDATE role
SET deleted_at = CURRENT_TIMESTAMP
WHERE id = ?
`

func (q *Queries) DeleteRole(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteRole, id)
	return err
}

const getRateLimiters = `-- name: GetRateLimiters :many
SELECT id, name, destination, deleted_at FROM rate_limiter
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
			&i.Destination,
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

const getRole = `-- name: GetRole :one
SELECT r.endpoint, r.operation, r.bucket_size, r.add_token_per_min, r.initial_tokens, ra.destination FROM role As r
       LEFT JOIN rate_limiter as ra ON ra.id = r.rate_limiter_id
WHERE r.deleted_at is null and r.endpoint = ? and r.operation = ?
`

type GetRoleParams struct {
	Endpoint  string
	Operation string
}

type GetRoleRow struct {
	Endpoint       string
	Operation      string
	BucketSize     int64
	AddTokenPerMin int64
	InitialTokens  int64
	Destination    sql.NullString
}

func (q *Queries) GetRole(ctx context.Context, arg GetRoleParams) (GetRoleRow, error) {
	row := q.db.QueryRowContext(ctx, getRole, arg.Endpoint, arg.Operation)
	var i GetRoleRow
	err := row.Scan(
		&i.Endpoint,
		&i.Operation,
		&i.BucketSize,
		&i.AddTokenPerMin,
		&i.InitialTokens,
		&i.Destination,
	)
	return i, err
}

const getRoles = `-- name: GetRoles :many
SELECT id, endpoint, operation,  bucket_size, add_token_per_min, initial_tokens, rate_limiter_id FROM role
WHERE deleted_at is NULL and rate_limiter_id = ?
`

type GetRolesRow struct {
	ID             int64
	Endpoint       string
	Operation      string
	BucketSize     int64
	AddTokenPerMin int64
	InitialTokens  int64
	RateLimiterID  int64
}

func (q *Queries) GetRoles(ctx context.Context, rateLimiterID int64) ([]GetRolesRow, error) {
	rows, err := q.db.QueryContext(ctx, getRoles, rateLimiterID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetRolesRow
	for rows.Next() {
		var i GetRolesRow
		if err := rows.Scan(
			&i.ID,
			&i.Endpoint,
			&i.Operation,
			&i.BucketSize,
			&i.AddTokenPerMin,
			&i.InitialTokens,
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
