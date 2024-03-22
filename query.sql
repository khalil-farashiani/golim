-- name: GetRoles :many
SELECT id, endpoint, operation,  bucket_size, add_token_per_min, initial_tokens, rate_limiter_id FROM role
WHERE deleted_at is NULL and rate_limiter_id = ?;

-- name: CreateRole :one
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
    RETURNING *;


-- name: DeleteRole :exec
UPDATE role
SET deleted_at = CURRENT_TIMESTAMP
WHERE id = ?;


-- name: GetRateLimiters :many
SELECT * FROM rate_limiter
WHERE deleted_at = null;


-- name: CrateRateLimiter :one
INSERT INTO rate_limiter (
    name,
    destination
) VALUES (
             ?,
             ?
         )
    RETURNING *;

-- name: DeleteRateLimiter :exec
UPDATE rate_limiter
    SET deleted_at = CURRENT_TIMESTAMP
    WHERE id = ?;


-- name: GetRole :one
SELECT r.endpoint, r.operation, r.bucket_size, r.add_token_per_min, r.initial_tokens, ra.destination FROM role As r
       LEFT JOIN rate_limiter as ra ON ra.id = r.rate_limiter_id
WHERE r.deleted_at is null and r.endpoint = ? and r.operation = ?;

