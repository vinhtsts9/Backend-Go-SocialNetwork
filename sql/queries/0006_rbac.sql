-- file: rbac_queries.sql

-- AddRoleToUser
-- name: AddRoleToUser :exec
INSERT INTO user_roles (user_id, role_id)
VALUES (?, ?)
ON DUPLICATE KEY UPDATE
    role_id = VALUES(role_id), 
    updated_at = NOW();

-- RemoveRoleFromUser
-- name: RemoveRoleFromUser :exec
DELETE FROM user_roles
WHERE user_id = ? AND role_id = ?;

-- GetRolesByUserID
-- name: GetRolesByUserID :many
SELECT r.id, r.role_name, r.description
FROM roles r
JOIN user_roles ur ON r.id = ur.role_id
WHERE ur.user_id = ?;

-- GetPermissionsByRoleID
-- name: GetPermissionsByRoleID :many
SELECT p.id, p.permission_name, p.resource, p.action
FROM permissions p
JOIN role_permissions rp ON p.id = rp.permission_id
WHERE rp.role_id = ?;

-- AddPermissionToRole
-- name: AddPermissionToRole :exec
INSERT INTO role_permissions (role_id, permission_id)
VALUES (?, ?)
ON DUPLICATE KEY UPDATE
    permission_id = VALUES(permission_id),
    updated_at = NOW();

-- RemovePermissionFromRole
-- name: RemovePermissionFromRole :exec
DELETE FROM role_permissions
WHERE role_id = ? AND permission_id = ?;

-- GetPermissionsByUserID
-- name: GetPermissionsByUserID :many
SELECT p.id, p.permission_name, p.resource, p.action
FROM permissions p
JOIN role_permissions rp ON p.id = rp.permission_id
JOIN user_roles ur ON rp.role_id = ur.role_id
WHERE ur.user_id = ?;

-- GetRoleByID
-- name: GetRoleByID :one
SELECT id, role_name, description
FROM roles
WHERE id = ?;

-- CreatePermission
-- name: CreatePermission :exec
INSERT INTO permissions (permission_name, resource, action, description)
VALUES (?, ?, ?, ?);
SELECT LAST_INSERT_ID() AS id, ?, ?, ?, ?;

-- CreateRole
-- name: CreateRole :exec
INSERT INTO roles (role_name, description)
VALUES (?, ?);
SELECT LAST_INSERT_ID() AS id, ?, ?;

-- CheckUserPermission
-- name: CheckUserPermission :one
select distinct p.*
from permissions p,
(select rp.permission_id
from user_roles ur,
     role_permissions rp
where ur.role_id = rp.role_id 
and ur.user_id = ?   
) t
where p.permission_id = t.permission_id;