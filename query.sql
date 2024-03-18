-- name: GetRoles :many
SELECT * FROM role
WHERE deleted_at = null;

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


-- name: GetRateLimiters :many
SELECT * FROM rate_limiter
WHERE deleted_at = null;


-- name: CrateRateLimiter :one
INSERT INTO rate_limiter (
    name
) VALUES (
             ?
         )
    RETURNING *;