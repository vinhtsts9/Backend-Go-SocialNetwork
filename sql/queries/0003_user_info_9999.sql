-- name: GetUser :one
SELECT
    user_id,
    user_account,
    user_nickname,
    user_avatar,
    user_state,
    user_mobile,
    user_gender,
    user_birthday,
    user_email,
    user_is_authencation,
    created_at,
    updated_at
FROM `user_info`
WHERE user_id = ? LIMIT 1;

-- name: GetUsers :many
SELECT
    user_id,
    user_account,
    user_nickname,
    user_avatar,
    user_state,
    user_mobile,
    user_gender,
    user_birthday,
    user_email,
    user_is_authencation,
    created_at,
    updated_at
FROM `user_info`
WHERE user_id IN (?);

-- name: FindUsers :many
SELECT * FROM user_info WHERE user_account LIKE ? OR user_nickname LIKE ?;

-- name: ListUsers :many
SELECT * FROM user_info LIMIT ? OFFSET ?;

-- name: RemoveUser :exec
DELETE FROM user_info WHERE user_id = ?;

-- -- name: UpdatePassword :exec
-- UPDATE `user_info` SET user_password = ? WHERE user_id = ?;

-- name: AddUserAutoUserId :execresult
INSERT INTO `user_info` (
    user_account,
    user_nickname,
    user_avatar,
    user_state,
    user_mobile,
    user_gender,
    user_birthday,
    user_email,
    user_is_authencation
)
VALUES (?,?,?,?,?,?,?,?,?);

-- name: AddUserHaveUserId :execresult
INSERT INTO `user_info` (
    user_id,
    user_account,
    user_nickname,
    user_avatar,
    user_state,
    user_mobile,
    user_gender,
    user_birthday,
    user_email,
    user_is_authencation
)
VALUES (?,?,?,?,?,?,?,?,?,?);

-- name: EditUserByUserId :execresult
UPDATE `user_info`
SET user_nickname = ?, user_avatar = ?, user_mobile = ?,
user_gender = ?, user_birthday = ?, user_email = ?, updated_at = NOW()
WHERE user_id = ? AND user_is_authencation = 1;


-- name: UpdateAvatar :execresult
UPDATE `user_info`
SET user_avatar = ?
WHERE user_id = ? AND user_is_authencation = 1;