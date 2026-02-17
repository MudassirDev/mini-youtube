-- name: CreateVideo :one
INSERT INTO videos (
    user_id, title, url, description
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;
