package service

import (
	"context"
	"go-ecommerce-backend-api/m/v2/internal/database"
)

// RBACService định nghĩa các phương thức cho tầng service.
type RBACService interface {
	// Quản lý permissions
	CreatePermission(ctx context.Context, name, resource, action, description string) error
	GetPermissionsByRoleID(ctx context.Context, roleID int32) ([]database.GetPermissionsByRoleIDRow, error)
	GetPermissionsByUserID(ctx context.Context, userID int64) ([]database.GetPermissionsByUserIDRow, error)

	// Quản lý roles
	CreateRole(ctx context.Context, name, description string) error
	GetRolesByUserID(ctx context.Context, userID int64) ([]database.Role, error)
	AddPermissionToRole(ctx context.Context, roleID, permissionID int32) error
	RemovePermissionFromRole(ctx context.Context, roleID, permissionID int32) error

	// Quản lý user-roles
	AddRoleToUser(ctx context.Context, userID int64, roleID int32) error
	RemoveRoleFromUser(ctx context.Context, userID int64, roleID int32) error

	// Kiểm tra quyền
	CheckUserPermission(ctx context.Context, userID int64) (database.Permission, error)
}

var localRBACService RBACService

func InitRBACService(i RBACService) {
	localRBACService = i
}

func RbacService() RBACService {
	if localRBACService == nil {
		panic("implement rbacService failed")
	}
	return localRBACService
}
