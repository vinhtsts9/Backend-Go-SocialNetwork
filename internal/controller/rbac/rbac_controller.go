package rbac

import (
	"go-ecommerce-backend-api/m/v2/global"
	"go-ecommerce-backend-api/m/v2/internal/service"
	"go-ecommerce-backend-api/m/v2/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// management RBAC roles and permissions
var RBAC = new(cRBAC)

type cRBAC struct{}

// CreatePermission
// @Summary      Create Permission
// @Description  Create a new permission
// @Tags         RBAC management
// @Accept       json
// @Produce      json
// @Param        payload body model.CreatePermissionInput true "payload"
// @Success      200  {object}  response.ResponseData
// @Failure      500  {object}  response.ErrorResponseData
// @Router       /rbac/create_permission [post]
func (c *cRBAC) CreatePermission(ctx *gin.Context) {
	var params struct {
		Name        string `json:"name"`
		Resource    string `json:"resource"`
		Action      string `json:"action"`
		Description string `json:"description"`
	}

	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.ErrorResponse(ctx, response.ErrCodeParamInvalid, err.Error())
		return
	}

	err := service.RbacService().CreatePermission(ctx, params.Name, params.Resource, params.Action, params.Description)
	if err != nil {
		global.Logger.Error("Error creating permission", zap.Error(err))
		response.ErrorResponse(ctx, response.ErrCodeParamInvalid, err.Error())
		return
	}

	response.SuccessResponse(ctx, response.ErrCodeSuccess, nil)
}

// CreateRole
// @Summary      Create Role
// @Description  Create a new role
// @Tags         RBAC management
// @Accept       json
// @Produce      json
// @Param        payload body model.CreateRoleInput true "payload"
// @Success      200  {object}  response.ResponseData
// @Failure      500  {object}  response.ErrorResponseData
// @Router       /rbac/create_role [post]
func (c *cRBAC) CreateRole(ctx *gin.Context) {
	var params struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.ErrorResponse(ctx, response.ErrCodeParamInvalid, err.Error())
		return
	}

	err := service.RbacService().CreateRole(ctx, params.Name, params.Description)
	if err != nil {
		global.Logger.Error("Error creating role", zap.Error(err))
		response.ErrorResponse(ctx, response.ErrCodeParamInvalid, err.Error())
		return
	}

	response.SuccessResponse(ctx, response.ErrCodeSuccess, nil)
}

// AddPermissionToRole
// @Summary      Add Permission to Role
// @Description  Assign permission to a role
// @Tags         RBAC management
// @Accept       json
// @Produce      json
// @Param        payload body model.AddPermissionToRoleInput true "payload"
// @Success      200  {object}  response.ResponseData
// @Failure      500  {object}  response.ErrorResponseData
// @Router       /rbac/add_permission_to_role [post]
func (c *cRBAC) AddPermissionToRole(ctx *gin.Context) {
	var params struct {
		RoleID       int32 `json:"role_id"`
		PermissionID int32 `json:"permission_id"`
	}

	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.ErrorResponse(ctx, response.ErrCodeParamInvalid, err.Error())
		return
	}

	err := service.RbacService().AddPermissionToRole(ctx, params.RoleID, params.PermissionID)
	if err != nil {
		global.Logger.Error("Error adding permission to role", zap.Error(err))
		response.ErrorResponse(ctx, response.ErrCodeParamInvalid, err.Error())
		return
	}

	response.SuccessResponse(ctx, response.ErrCodeSuccess, nil)
}

// RemovePermissionFromRole
// @Summary      Remove Permission from Role
// @Description  Remove permission from a role
// @Tags         RBAC management
// @Accept       json
// @Produce      json
// @Param        payload body model.RemovePermissionFromRoleInput true "payload"
// @Success      200  {object}  response.ResponseData
// @Failure      500  {object}  response.ErrorResponseData
// @Router       /rbac/remove_permission_from_role [post]
func (c *cRBAC) RemovePermissionFromRole(ctx *gin.Context) {
	var params struct {
		RoleID       int32 `json:"role_id"`
		PermissionID int32 `json:"permission_id"`
	}

	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.ErrorResponse(ctx, response.ErrCodeParamInvalid, err.Error())
		return
	}

	err := service.RbacService().RemovePermissionFromRole(ctx, params.RoleID, params.PermissionID)
	if err != nil {
		global.Logger.Error("Error removing permission from role", zap.Error(err))
		response.ErrorResponse(ctx, response.ErrCodeParamInvalid, err.Error())
		return
	}

	response.SuccessResponse(ctx, response.ErrCodeSuccess, nil)
}

// AddRoleToUser
// @Summary      Add Role to User
// @Description  Assign a role to a user
// @Tags         RBAC management
// @Accept       json
// @Produce      json
// @Param        payload body model.AddRoleToUserInput true "payload"
// @Success      200  {object}  response.ResponseData
// @Failure      500  {object}  response.ErrorResponseData
// @Router       /rbac/add_role_to_user [post]
func (c *cRBAC) AddRoleToUser(ctx *gin.Context) {
	var params struct {
		UserID int64 `json:"user_id"`
		RoleID int32 `json:"role_id"`
	}

	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.ErrorResponse(ctx, response.ErrCodeParamInvalid, err.Error())
		return
	}

	err := service.RbacService().AddRoleToUser(ctx, params.UserID, params.RoleID)
	if err != nil {
		global.Logger.Error("Error adding role to user", zap.Error(err))
		response.ErrorResponse(ctx, response.ErrCodeParamInvalid, err.Error())
		return
	}

	response.SuccessResponse(ctx, response.ErrCodeSuccess, nil)
}

// RemoveRoleFromUser
// @Summary      Remove Role from User
// @Description  Remove a role from a user
// @Tags         RBAC management
// @Accept       json
// @Produce      json
// @Param        payload body model.RemoveRoleFromUserInput true "payload"
// @Success      200  {object}  response.ResponseData
// @Failure      500  {object}  response.ErrorResponseData
// @Router       /rbac/remove_role_from_user [post]
func (c *cRBAC) RemoveRoleFromUser(ctx *gin.Context) {
	var params struct {
		UserID int64 `json:"user_id"`
		RoleID int32 `json:"role_id"`
	}

	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.ErrorResponse(ctx, response.ErrCodeParamInvalid, err.Error())
		return
	}

	err := service.RbacService().RemoveRoleFromUser(ctx, params.UserID, params.RoleID)
	if err != nil {
		global.Logger.Error("Error removing role from user", zap.Error(err))
		response.ErrorResponse(ctx, response.ErrCodeParamInvalid, err.Error())
		return
	}

	response.SuccessResponse(ctx, response.ErrCodeSuccess, nil)
}

// CheckUserPermission
// @Summary      Check User Permission
// @Description  Check if a user has a specific permission
// @Tags         RBAC management
// @Accept       json
// @Produce      json
// @Param        payload body model.CheckUserPermissionInput true "payload"
// @Success      200  {object}  response.ResponseData
// @Failure      500  {object}  response.ErrorResponseData
// @Router       /rbac/check_user_permission [post]
func (c *cRBAC) CheckUserPermission(ctx *gin.Context) {
	var params struct {
		UserID int64 `json:"user_id"`
	}

	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.ErrorResponse(ctx, response.ErrCodeParamInvalid, err.Error())
		return
	}

	permission, err := service.RbacService().CheckUserPermission(ctx, params.UserID)
	if err != nil {
		global.Logger.Error("Error checking user permission", zap.Error(err))
		response.ErrorResponse(ctx, response.ErrCodeParamInvalid, err.Error())
		return
	}

	response.SuccessResponse(ctx, response.ErrCodeSuccess, permission)
}
