-- Get all post
-- name: GetAllpost :many
SELECT *
FROM post p
left join user_follows uf 
on p.user_id = uf.following_id
WHERE uf.follower_id =?
order by p.created_at desc;



-- Get post by ID
-- name: GetPostById :one
SELECT 
    id, 
    title, 
    COALESCE(image_paths, '[]') AS image_paths,  -- tránh NULL trả về
    user_nickname, 
    created_at, 
    updated_at
FROM post
WHERE id = ?;

-- Create a new post
-- name: CreatePost :execresult
INSERT INTO post (title, image_paths, user_id, created_at, updated_at, user_nickname)
VALUES (?, ?, ?, NOW(), NOW(), ?);

-- Update a post
-- name: UpdatePost :execresult
UPDATE post
SET title = ?, image_paths = ?, updated_at = NOW()
WHERE id = ?;

-- Delete a post
-- name: DeletePost :exec
DELETE FROM post
WHERE id = ?;

-- Get post by user ID
-- name: GetpostByUserId :many
SELECT 
    id, 
    title, 
    COALESCE(image_paths, '[]') AS image_paths,  -- Tránh NULL trả về
    user_id, 
    created_at, 
    updated_at
FROM post
WHERE user_id = ?;
