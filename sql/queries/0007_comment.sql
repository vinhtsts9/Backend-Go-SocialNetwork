-- name: GetMaxRightComment :one
select comment_right from Comment 
where post_id = ?
order by comment_right desc
limit 1;

-- name: CreateComment :execresult
INSERT INTO Comment (
    post_id, 
    user_id, 
    user_nickname, 
    comment_content, 
    comment_left, 
    comment_right, 
    comment_parent, 
    isDeleted
) 
VALUES (?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetCommentByLastInsertId :one
select * from Comment where id = ?;
-- name: GetCommentByID :one
SELECT * 
FROM Comment
WHERE id = ?;

-- name: UpdateCommentRightCreate :exec
update Comment
set comment_right = comment_right + 2
where post_id = ? and comment_right >= ?;

-- name: UpdateCommentLeftCreate :exec
update Comment
set comment_left = comment_left + 2
where post_id = ? and comment_left > ?;

-- name: GetCommentByParentID :many
select c.* from Comment c 
where c.post_id = ?
and c.comment_left > ( select sub.comment_left from Comment sub where sub.id = ?)
and c.comment_right < (select sub.comment_right from Comment sub where sub.id = ?)
and c.isDeleted = false
order by c.comment_left
limit 10 
offset 0;

-- name: DeleteCommentsInRange :exec
delete from Comment
where post_id = ?
and comment_left >= ?
and comment_right <= ?;

-- name: UpdateCommentLeft :exec
update Comment 
set comment_left = comment_left - ?
where post_id = ?
and comment_left >= ?;

-- name: UpdateCommentRight :exec
update Comment
set comment_right = comment_right - ?
where post_id = ?
and comment_right >= ?;

-- name: GetRootComment :many
select * from Comment 
where comment_parent is null and  post_id = ?;

-- name: AddReplyCommentParent :exec
update Comment
set reply_count = reply_count + 1
where id = ?;