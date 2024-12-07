package impl

import (
	"context"
	"database/sql"
	"go-ecommerce-backend-api/m/v2/internal/database"
)

type rbacService struct {
	r *database.Queries
}

func NewRbacImpl(r *database.Queries) *rbacService {
	return &rbacService{
		r: r,
	}
}

// Tạo permission mới
func (s *rbacService) CreatePermission(ctx context.Context, name, resource, action, description string) error {
	arg := database.CreatePermissionParams{
		PermissionName: name,
		Resource:       resource,
		Action:         action,
		Description:    sql.NullString{String: description, Valid: true},
	}
	return s.r.CreatePermission(ctx, arg)
}

// Lấy danh sách permissions theo role
func (s *rbacService) GetPermissionsByRoleID(ctx context.Context, roleID int32) ([]database.GetPermissionsByRoleIDRow, error) {
	return s.r.GetPermissionsByRoleID(ctx, sql.NullInt32{Int32: roleID, Valid: true})
}

// Lấy danh sách permissions theo user
func (s *rbacService) GetPermissionsByUserID(ctx context.Context, userID int64) ([]database.GetPermissionsByUserIDRow, error) {
	return s.r.GetPermissionsByUserID(ctx, sql.NullInt64{Int64: userID, Valid: true})
}

// Tạo role mới
func (s *rbacService) CreateRole(ctx context.Context, name, description string) error {
	arg := database.CreateRoleParams{
		RoleName:    name,
		Description: sql.NullString{String: description, Valid: true},
	}
	return s.r.CreateRole(ctx, arg)
}

// Lấy danh sách roles của user
func (s *rbacService) GetRolesByUserID(ctx context.Context, userID int64) ([]database.Role, error) {
	return s.r.GetRolesByUserID(ctx, sql.NullInt64{Int64: userID, Valid: true})
}

// Thêm permission vào role
func (s *rbacService) AddPermissionToRole(ctx context.Context, roleID, permissionID int32) error {
	arg := database.AddPermissionToRoleParams{
		RoleID:       sql.NullInt32{Int32: roleID, Valid: true},
		PermissionID: sql.NullInt32{Int32: permissionID, Valid: true},
	}
	return s.r.AddPermissionToRole(ctx, arg)
}

// Gỡ permission khỏi role
func (s *rbacService) RemovePermissionFromRole(ctx context.Context, roleID, permissionID int32) error {
	arg := database.RemovePermissionFromRoleParams{
		RoleID:       sql.NullInt32{Int32: roleID, Valid: true},
		PermissionID: sql.NullInt32{Int32: permissionID, Valid: true},
	}
	return s.r.RemovePermissionFromRole(ctx, arg)
}

// Thêm role cho user
func (s *rbacService) AddRoleToUser(ctx context.Context, userID int64, roleID int32) error {
	arg := database.AddRoleToUserParams{
		UserID: sql.NullInt64{Int64: userID, Valid: true},
		RoleID: sql.NullInt32{Int32: roleID, Valid: true},
	}
	return s.r.AddRoleToUser(ctx, arg)
}

// Gỡ role khỏi user
func (s *rbacService) RemoveRoleFromUser(ctx context.Context, userID int64, roleID int32) error {
	arg := database.RemoveRoleFromUserParams{
		UserID: sql.NullInt64{Int64: userID, Valid: true},
		RoleID: sql.NullInt32{Int32: roleID, Valid: true},
	}
	return s.r.RemoveRoleFromUser(ctx, arg)
}

// Kiểm tra quyền của user
func (s *rbacService) CheckUserPermission(ctx context.Context, userID int64) (database.Permission, error) {
	return s.r.CheckUserPermission(ctx, sql.NullInt64{Int64: userID, Valid: true})
}
