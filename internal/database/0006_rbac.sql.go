// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: 0006_rbac.sql

package database

import (
	"context"
	"database/sql"
)

const addPermissionToRole = `-- name: AddPermissionToRole :exec
INSERT INTO role_permissions (role_id, permission_id)
VALUES (?, ?)
ON DUPLICATE KEY UPDATE
    permission_id = VALUES(permission_id),
    updated_at = NOW()
`

type AddPermissionToRoleParams struct {
	RoleID       sql.NullInt32
	PermissionID sql.NullInt32
}

// AddPermissionToRole
func (q *Queries) AddPermissionToRole(ctx context.Context, arg AddPermissionToRoleParams) error {
	_, err := q.db.ExecContext(ctx, addPermissionToRole, arg.RoleID, arg.PermissionID)
	return err
}

const addRoleToUser = `-- name: AddRoleToUser :exec

INSERT INTO user_roles (user_id, role_id)
VALUES (?, ?)
ON DUPLICATE KEY UPDATE
    role_id = VALUES(role_id), 
    updated_at = NOW()
`

type AddRoleToUserParams struct {
	UserID sql.NullInt64
	RoleID sql.NullInt32
}

// file: rbac_queries.sql
// AddRoleToUser
func (q *Queries) AddRoleToUser(ctx context.Context, arg AddRoleToUserParams) error {
	_, err := q.db.ExecContext(ctx, addRoleToUser, arg.UserID, arg.RoleID)
	return err
}

const checkUserPermission = `-- name: CheckUserPermission :one
select distinct p.id, p.permission_name, p.description, p.resource, p.action
from permissions p,
(select rp.permission_id
from user_roles ur,
     role_permissions rp
where ur.role_id = rp.role_id 
and ur.user_id = ?   
) t
where p.id = t.permission_id
`

// CheckUserPermission
func (q *Queries) CheckUserPermission(ctx context.Context, userID sql.NullInt64) (Permission, error) {
	row := q.db.QueryRowContext(ctx, checkUserPermission, userID)
	var i Permission
	err := row.Scan(
		&i.ID,
		&i.PermissionName,
		&i.Description,
		&i.Resource,
		&i.Action,
	)
	return i, err
}

const createPermission = `-- name: CreatePermission :exec
INSERT INTO permissions (permission_name, resource, action, description)
VALUES (?, ?, ?, ?)
`

type CreatePermissionParams struct {
	PermissionName string
	Resource       string
	Action         string
	Description    sql.NullString
}

// CreatePermission
func (q *Queries) CreatePermission(ctx context.Context, arg CreatePermissionParams) error {
	_, err := q.db.ExecContext(ctx, createPermission,
		arg.PermissionName,
		arg.Resource,
		arg.Action,
		arg.Description,
	)
	return err
}

const createRole = `-- name: CreateRole :exec
INSERT INTO roles (role_name, description)
VALUES (?, ?)
`

type CreateRoleParams struct {
	RoleName    string
	Description sql.NullString
}

// CreateRole
func (q *Queries) CreateRole(ctx context.Context, arg CreateRoleParams) error {
	_, err := q.db.ExecContext(ctx, createRole, arg.RoleName, arg.Description)
	return err
}

const getPermissionsByRoleID = `-- name: GetPermissionsByRoleID :many
SELECT p.id, p.permission_name, p.resource, p.action
FROM permissions p
JOIN role_permissions rp ON p.id = rp.permission_id
WHERE rp.role_id = ?
`

type GetPermissionsByRoleIDRow struct {
	ID             int32
	PermissionName string
	Resource       string
	Action         string
}

// GetPermissionsByRoleID
func (q *Queries) GetPermissionsByRoleID(ctx context.Context, roleID sql.NullInt32) ([]GetPermissionsByRoleIDRow, error) {
	rows, err := q.db.QueryContext(ctx, getPermissionsByRoleID, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetPermissionsByRoleIDRow
	for rows.Next() {
		var i GetPermissionsByRoleIDRow
		if err := rows.Scan(
			&i.ID,
			&i.PermissionName,
			&i.Resource,
			&i.Action,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getPermissionsByUserID = `-- name: GetPermissionsByUserID :many
SELECT distinct p.id, p.permission_name, p.resource, p.action
FROM permissions p
JOIN role_permissions rp ON p.id = rp.permission_id
JOIN user_roles ur ON rp.role_id = ur.role_id
WHERE ur.user_id = ?
`

type GetPermissionsByUserIDRow struct {
	ID             int32
	PermissionName string
	Resource       string
	Action         string
}

// GetPermissionsByUserID
func (q *Queries) GetPermissionsByUserID(ctx context.Context, userID sql.NullInt64) ([]GetPermissionsByUserIDRow, error) {
	rows, err := q.db.QueryContext(ctx, getPermissionsByUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetPermissionsByUserIDRow
	for rows.Next() {
		var i GetPermissionsByUserIDRow
		if err := rows.Scan(
			&i.ID,
			&i.PermissionName,
			&i.Resource,
			&i.Action,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getRoleByID = `-- name: GetRoleByID :one
SELECT id, role_name, description
FROM roles
WHERE id = ?
`

// GetRoleByID
func (q *Queries) GetRoleByID(ctx context.Context, id int32) (Role, error) {
	row := q.db.QueryRowContext(ctx, getRoleByID, id)
	var i Role
	err := row.Scan(&i.ID, &i.RoleName, &i.Description)
	return i, err
}

const getRolesByUserID = `-- name: GetRolesByUserID :many
SELECT r.id, r.role_name, r.description
FROM roles r
JOIN user_roles ur ON r.id = ur.role_id
WHERE ur.user_id = ?
`

// GetRolesByUserID
func (q *Queries) GetRolesByUserID(ctx context.Context, userID sql.NullInt64) ([]Role, error) {
	rows, err := q.db.QueryContext(ctx, getRolesByUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Role
	for rows.Next() {
		var i Role
		if err := rows.Scan(&i.ID, &i.RoleName, &i.Description); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const removePermissionFromRole = `-- name: RemovePermissionFromRole :exec
DELETE FROM role_permissions
WHERE role_id = ? AND permission_id = ?
`

type RemovePermissionFromRoleParams struct {
	RoleID       sql.NullInt32
	PermissionID sql.NullInt32
}

// RemovePermissionFromRole
func (q *Queries) RemovePermissionFromRole(ctx context.Context, arg RemovePermissionFromRoleParams) error {
	_, err := q.db.ExecContext(ctx, removePermissionFromRole, arg.RoleID, arg.PermissionID)
	return err
}

const removeRoleFromUser = `-- name: RemoveRoleFromUser :exec
DELETE FROM user_roles
WHERE user_id = ? AND role_id = ?
`

type RemoveRoleFromUserParams struct {
	UserID sql.NullInt64
	RoleID sql.NullInt32
}

// RemoveRoleFromUser
func (q *Queries) RemoveRoleFromUser(ctx context.Context, arg RemoveRoleFromUserParams) error {
	_, err := q.db.ExecContext(ctx, removeRoleFromUser, arg.UserID, arg.RoleID)
	return err
}
