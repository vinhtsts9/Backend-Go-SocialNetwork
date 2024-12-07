-- Get all post
-- name: GetAllpost :many
SELECT id, title, content, user_id, created_at, updated_at
FROM post;

-- Get post by ID
-- name: GetPostById :one
SELECT id, title, content, user_id, created_at, updated_at
FROM post
WHERE id = ?;

-- Create a new post
-- name: CreatePost :exec
INSERT INTO post (title, content, user_id, created_at, updated_at)
VALUES (?, ?, ?, NOW(), NOW());

-- Update a post
-- name: UpdatePost :exec
UPDATE post
SET title = ?, content = ?, updated_at = NOW()
WHERE id = ?;

-- Delete a post
-- name: DeletePost :exec
DELETE FROM post
WHERE id = ?;

-- Get post by user ID
-- name: GetpostByUserId :many
SELECT id, title, content, user_id, created_at, updated_at
FROM post
WHERE user_id = ?;
