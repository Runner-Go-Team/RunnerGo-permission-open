package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"permission-open/internal/pkg/biz/errno"
	"permission-open/internal/pkg/biz/jwt"
	"permission-open/internal/pkg/biz/response"
	"permission-open/internal/pkg/dal/rao"
	"permission-open/internal/pkg/logic/errmsg"
	"permission-open/internal/pkg/logic/permission"
)

func PermissionList(ctx *gin.Context) {
	var req rao.PermissionListReq
	if err := ctx.ShouldBind(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	permissions, err := permission.GetPermissionList(ctx, req.RoleID)
	if err != nil {
		if errors.Is(err, errmsg.ErrRoleExists) {
			response.ErrorWithMsg(ctx, errno.ErrRoleExists, err.Error())
			return
		}
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	response.SuccessWithData(ctx, rao.PermissionListResp{
		List: permissions,
	})
	return
}

// SetRolePermission 设置角色权限
func SetRolePermission(ctx *gin.Context) {
	var req rao.SetRolePermissionReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	err := permission.SetRolePermission(ctx, jwt.GetUserIDByCtx(ctx), req.RoleID, req.PermissionMarks, req.RoleName)
	if err != nil {
		if errors.Is(err, errmsg.ErrRoleNotChange) {
			response.ErrorWithMsg(ctx, errno.ErrRoleNotChange, "")
			return
		}
		if errors.Is(err, errmsg.ErrRoleExists) {
			response.ErrorWithMsg(ctx, errno.ErrRoleExists, "")
			return
		}
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	response.Success(ctx)
	return
}

// UserPermissions 获取用户权限
func UserPermissions(ctx *gin.Context) {
	var req rao.UserPermissionsReq
	if err := ctx.ShouldBind(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	// 如果传 TargetUserID 查目标用户权限信息
	userID := jwt.GetUserIDByCtx(ctx)
	if len(req.TargetUserID) > 0 {
		userID = req.TargetUserID
	}

	permissions, err := permission.GetUserPermissionList(ctx, userID, req.CompanyId, req.TeamID)
	if err != nil {
		if errors.Is(err, errmsg.ErrUserNotRole) {
			response.ErrorWithMsg(ctx, errno.ErrUserNotRole, err.Error())
			return
		}
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	response.SuccessWithData(ctx, rao.UserPermissionsResp{
		Permissions: permissions,
	})
	return
}

// UserAllPermissionMarks 获取用户的全部角色对应的mark
func UserAllPermissionMarks(ctx *gin.Context) {
	var req rao.OpenUserPermissionMarksReq
	if err := ctx.ShouldBind(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	// 如果传 TargetUserID 查目标用户权限信息
	userID := jwt.GetUserIDByCtx(ctx)
	if len(req.UserId) > 0 {
		userID = req.UserId
	}
	permissions, err := permission.GetUserAllPermissionMarks(ctx, userID)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	response.SuccessWithData(ctx, permissions)
	return
}

// PermissionCheckUrl 通过URL判断是否有权限
func PermissionCheckUrl(ctx *gin.Context) {
	var req rao.PermissionCheckUrlReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	isHave, err := permission.CheckUrl(ctx, req.UserID, req.TeamID, req.Url)
	if err != nil {
		if errors.Is(err, errmsg.ErrUserNotRole) {
			response.ErrorWithMsg(ctx, errno.ErrUserNotRole, err.Error())
			return
		}
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	response.SuccessWithData(ctx, rao.PermissionCheckUrlResp{
		IsHave: isHave,
	})
}

func Clear(ctx *gin.Context) {
	// 清空 company  user_company    team-company_id   role   user_role    role_permission
	companyID := "5fb4d194-238b-47f6-86bc-0c6010d4953f"
	userID := "bb31bfd3-22fe-42af-95d4-f8a37ef46d52"
	roleIDC2 := "04a14961-7951-4606-a5a2-d2ddc1f959cb" // 管理员
	roleIDT1 := "45ee95d9-2a0c-4256-8638-3e962a1de71e" // 团队管理员
	roleIDT2 := "951d2c4b-c1ef-4e8c-8ca5-a5256974fc88" // 团队成员
	err := permission.Clear(ctx, companyID, userID, roleIDC2, roleIDT1, roleIDT2)

	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	response.Success(ctx)
}
