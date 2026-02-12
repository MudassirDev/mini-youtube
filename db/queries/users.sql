-- name: CreateAuthor :one
INSERT INTO users (
    email, username, password_hash, display_name
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;
