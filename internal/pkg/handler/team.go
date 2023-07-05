package handler

import (
	"errors"
	"permission-open/internal/pkg/biz/errno"
	"permission-open/internal/pkg/biz/jwt"
	"permission-open/internal/pkg/biz/response"
	"permission-open/internal/pkg/dal/rao"
	"permission-open/internal/pkg/logic/errmsg"
	"permission-open/internal/pkg/logic/team"

	"github.com/gin-gonic/gin"
)

// SaveTeam 创建或修改团队
func SaveTeam(ctx *gin.Context) {
	var req rao.SaveTeamReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	teamID, err := team.SaveTeam(ctx, &req)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	response.SuccessWithData(ctx, rao.SaveTeamResp{
		TeamID: teamID,
	})
	return
}

// UpdateTeam 修改团队
func UpdateTeam(ctx *gin.Context) {
	var req rao.UpdateTeamReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	err := team.UpdateTeam(ctx, req)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	response.Success(ctx)
	return
}

// ListTeam 团队列表
func ListTeam(ctx *gin.Context) {
	var req rao.ListTeamReq
	if err := ctx.ShouldBind(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}
	teams, err := team.ListTeam(ctx, jwt.GetUserIDByCtx(ctx), req)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	response.SuccessWithData(ctx, rao.ListTeamResp{Teams: teams})
	return
}

// RemoveTeamMember 移除成员
func RemoveTeamMember(ctx *gin.Context) {
	var req rao.RemoveTeamMemberReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	if err := team.RemoveMember(ctx, req.TeamID, jwt.GetUserIDByCtx(ctx), req.TargetUserID); err != nil {
		if errors.Is(err, errmsg.ErrTeamNotRemoveSelf) {
			response.ErrorWithMsg(ctx, errno.ErrTeamNotRemoveSelf, "")
			return
		}
		if errors.Is(err, errmsg.ErrTeamNotRemoveCreate) {
			response.ErrorWithMsg(ctx, errno.ErrTeamNotRemoveCreate, "")
			return
		}
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	response.Success(ctx)
	return
}

// DisbandTeam 解散团队
func DisbandTeam(ctx *gin.Context) {
	var req rao.DisbandTeamReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	err := team.DisbandTeam(ctx, req.TeamID, jwt.GetUserIDByCtx(ctx))
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	response.Success(ctx)
	return
}

// CollectionTeam 收藏|取消收藏
func CollectionTeam(ctx *gin.Context) {
	var req rao.CollectionTeamReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	err := team.CollectionTeam(ctx, jwt.GetUserIDByCtx(ctx), req.TeamID, req.Status)
	if err != nil {
		if errors.Is(err, errmsg.ErrTeamNotFound) {
			response.ErrorWithMsg(ctx, errno.ErrTeamNotFound, "")
			return
		}
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}
	response.Success(ctx)
	return
}

func TeamInfo(ctx *gin.Context) {
	var req rao.TeamInfoReq
	if err := ctx.ShouldBind(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	info, total, err := team.Info(ctx, req.TeamID, req.Keyword, req.Page, req.Size)
	if err != nil {
		if errors.Is(err, errmsg.ErrTeamNotFound) {
			response.ErrorWithMsg(ctx, errno.ErrTeamNotFound, "")
			return
		}
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}
	response.SuccessWithData(ctx, rao.TeamInfoResp{
		Team:  info,
		Total: total,
	})
	return
}

func SaveTeamMembers(ctx *gin.Context) {
	var req rao.SaveTeamMemberReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	userID := jwt.GetUserIDByCtx(ctx)
	if len(req.UserID) > 0 {
		userID = req.UserID
	}
	err := team.SaveMembers(ctx, userID, req.TeamID, req.Members)
	if err != nil {
		if errors.Is(err, errmsg.ErrTeamNotFound) {
			response.ErrorWithMsg(ctx, errno.ErrTeamNotFound, "")
			return
		}
		if errors.Is(err, errmsg.ErrTeamTypePrivate) {
			response.ErrorWithMsg(ctx, errno.ErrTeamTypePrivate, "")
			return
		}
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}
	response.Success(ctx)
	return
}

func TeamCompanyMembers(ctx *gin.Context) {
	var req rao.TeamCompanyMembersReq
	if err := ctx.ShouldBind(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	members, total, createdUserID, err := team.CompanyMembers(ctx, req.TeamID, req.Keyword, req.Page, req.Size)
	if err != nil {
		if errors.Is(err, errmsg.ErrTeamNotFound) {
			response.ErrorWithMsg(ctx, errno.ErrTeamNotFound, "")
			return
		}
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}
	response.SuccessWithData(ctx, rao.TeamCompanyMembersResp{
		Members:      members,
		CreateUserID: createdUserID,
		Total:        total,
	})
	return
}

// TeamTransferSuperRole 团队移交团队管理员
func TeamTransferSuperRole(ctx *gin.Context) {
	var req rao.TeamTransferSuperRoleReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	userID := jwt.GetUserIDByCtx(ctx)
	err := team.TransferSuperRole(ctx, userID, req.TeamID, req.TargetUserID)
	if err != nil {
		if errors.Is(err, errmsg.ErrUserNotRole) {
			response.ErrorWithMsg(ctx, errno.ErrUserNotRole, "")
			return
		}
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}
	response.Success(ctx)
	return
}
