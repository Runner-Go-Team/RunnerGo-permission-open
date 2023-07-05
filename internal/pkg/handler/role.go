package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"permission-open/internal/pkg/biz/consts"
	"permission-open/internal/pkg/biz/errno"
	"permission-open/internal/pkg/biz/jwt"
	"permission-open/internal/pkg/biz/response"
	"permission-open/internal/pkg/dal/rao"
	"permission-open/internal/pkg/logic/errmsg"
	"permission-open/internal/pkg/logic/role"
)

// RoleList 获取角色
func RoleList(ctx *gin.Context) {
	var req rao.RoleListReq
	if err := ctx.ShouldBind(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	roles, err := role.GetRoleList(ctx, req)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	response.SuccessWithData(ctx, rao.RoleListResp{
		RoleList: roles,
	})
	return
}

func SaveRole(ctx *gin.Context) {
	var req rao.SaveRoleReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	err := role.SaveRole(ctx, jwt.GetUserIDByCtx(ctx), req)
	if err != nil {
		if errors.Is(err, errmsg.ErrRoleExists) {
			response.ErrorWithMsg(ctx, errno.ErrRoleExists, err.Error())
			return
		}
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	response.Success(ctx)
	return
}

// SetTeamMemberRole 修改用户角色
func SetTeamMemberRole(ctx *gin.Context) {
	var req rao.SetTeamMemberRoleReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	err := role.SetTeamMemberRole(ctx, jwt.GetUserIDByCtx(ctx), req)
	if err != nil {
		if errors.Is(err, errmsg.ErrRoleNotChange) {
			response.ErrorWithMsg(ctx, errno.ErrRoleNotChange, err.Error())
			return
		}
		if errors.Is(err, errmsg.ErrUserForbidden) {
			response.ErrorWithMsg(ctx, errno.ErrUserForbidden, err.Error())
			return
		}
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	response.Success(ctx)
	return
}

// SetCompanyMemberRole 修改用户角色
func SetCompanyMemberRole(ctx *gin.Context) {
	var req rao.SetCompanyMemberRoleReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	err := role.SetCompanyMemberRole(ctx, jwt.GetUserIDByCtx(ctx), req)
	if err != nil {
		if errors.Is(err, errmsg.ErrRoleNotChange) {
			response.ErrorWithMsg(ctx, errno.ErrRoleNotChange, err.Error())
			return
		}
		if errors.Is(err, errmsg.ErrUserForbidden) {
			response.ErrorWithMsg(ctx, errno.ErrUserForbidden, err.Error())
			return
		}
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	response.Success(ctx)
	return
}

// RoleMember 获取用户角色信息
func RoleMember(ctx *gin.Context) {
	var req rao.RoleMemberReq
	if err := ctx.ShouldBind(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	userID := jwt.GetUserIDByCtx(ctx)
	if len(req.UserID) > 0 {
		userID = req.UserID
	}
	info, err := role.GetRoleMember(ctx, userID, &req)
	if err != nil {
		if errors.Is(err, errmsg.ErrUserNotRole) {
			response.ErrorWithMsg(ctx, errno.ErrUserNotRole, "")
			return
		}
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	level := info.Level
	// 超管/团队管理员只能有一个，所以可见从 level2 开始
	// 企业角色管理员不能管理管理员
	if level == consts.RoleLevelSuperManager {
		level = consts.RoleLevelManager
	} else if info.RoleType == consts.RoleTypeCompany && level == consts.RoleLevelManager {
		level = consts.RoleLevelGeneral
	}

	roleLists, err := role.GetRoleList(ctx, rao.RoleListReq{
		RoleType:  info.RoleType,
		Level:     level,
		CompanyId: req.CompanyID,
	})
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	response.SuccessWithData(ctx, rao.RoleMemberResp{
		Role:        info,
		UsableRoles: roleLists,
	})
	return
}

func RoleMembers(ctx *gin.Context) {
	var req rao.RoleMembersReq
	if err := ctx.ShouldBind(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	members, total, err := role.GetRoleMembers(ctx, &req)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	response.SuccessWithData(ctx, rao.RoleMembersResp{
		Members: members,
		Total:   total,
	})
	return
}

// RemoveRole 删除角色，关联角色更新为新的角色
func RemoveRole(ctx *gin.Context) {
	var req rao.RemoveRoleReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	err := role.RemoveRole(ctx, jwt.GetUserIDByCtx(ctx), req.RoleID, req.ChangeRoleID, req.CompanyID)
	if err != nil {
		if errors.Is(err, errmsg.ErrRoleNotDel) {
			response.ErrorWithMsg(ctx, errno.ErrRoleNotDel, "")
			return
		}
		if errors.Is(err, errmsg.ErrRoleNotChange) {
			response.ErrorWithMsg(ctx, errno.ErrRoleNotChange, "")
			return
		}
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	response.Success(ctx)
	return
}

func IsRemoveRole(ctx *gin.Context) {
	var req rao.IsRemoveRoleReq
	if err := ctx.ShouldBind(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	isAllowRemove, err := role.IsAllowRemove(ctx, req.RoleID)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	response.SuccessWithData(ctx, rao.IsRemoveRoleResp{
		IsAllowRemove: isAllowRemove,
	})
	return
}
