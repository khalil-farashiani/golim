-- Name: GetRoles :many
SELECT id, endpoint, Operation,  bucket_size, add_token_per_min, initial_tokens, rate_limiter_id FROM role
WHERE deleted_at is NULL and rate_limiter_id = ?;

-- Name: CreateRole :one
INSERT INTO role (
        endpoint,
        Operation,
        bucket_size,
        add_token_per_min,
        initial_tokens,
        rate_limiter_id
    ) VALUES (
             ?, ?, ?, ?, ?, ?
         )
    RETURNING *;


-- Name: DeleteRole :exec
UPDATE role
SET deleted_at = CURRENT_TIMESTAMP
WHERE id = ?;


-- Name: GetRateLimiters :many
SELECT * FROM rate_limiter
WHERE deleted_at = null;


-- Name: CrateRateLimiter :one
INSERT INTO rate_limiter (
    Name,
    Destination
) VALUES (
             ?,
             ?
         )
    RETURNING *;

-- Name: DeleteRateLimiter :exec
UPDATE rate_limiter
    SET deleted_at = CURRENT_TIMESTAMP
    WHERE id = ?;


-- Name: GetRole :one
SELECT r.endpoint, r.Operation, r.bucket_size, r.add_token_per_min, r.initial_tokens, ra.Destination FROM role As r
       LEFT JOIN rate_limiter as ra ON ra.id = r.rate_limiter_id
WHERE r.deleted_at is null and r.endpoint = ? and r.Operation = ?;

