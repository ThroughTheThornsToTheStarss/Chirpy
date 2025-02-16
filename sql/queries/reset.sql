-- name: Reset :exec
DELETE FROM users;

-- name: ResetRFToken :exec
DELETE FROM refresh_tokens;
