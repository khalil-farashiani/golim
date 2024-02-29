-- name: GetRoles :many
SELECT * FROM role
WHERE deleted_at = null;

-- name: CreateRole :one
INSERT INTO role (
        endpoint,
        operation,
        bucket_size,
        add_token_per_sec,
        initial_tokens
    ) VALUES (
             ?, ?, ?, ?, ?
         )
    RETURNING *;
