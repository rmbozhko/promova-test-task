-- name: CreatePost :one
INSERT INTO posts (
    title, content
) VALUES (
    $1, $2
) RETURNING *;

-- name: GetPostById :one
SELECT * FROM posts
WHERE id = $1;

-- name: GetPosts :many
SELECT * FROM posts
ORDER BY id;

-- name: UpdatePostById :one
UPDATE posts
SET title = $1, content = $2
WHERE id = $3
RETURNING *;

-- name: DeletePost :exec
DELETE FROM posts
WHERE id = $1;
