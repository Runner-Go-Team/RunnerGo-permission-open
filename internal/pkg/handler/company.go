package handler

import (
	"errors"
	"github.com/go-omnibus/omnibus"
	"permission-open/internal/pkg/biz/errno"
	"permission-open/internal/pkg/biz/jwt"
	"permission-open/internal/pkg/biz/response"
	"permission-open/internal/pkg/dal"
	"permission-open/internal/pkg/dal/global"
	"permission-open/internal/pkg/dal/query"
	"permission-open/internal/pkg/dal/rao"
	"permission-open/internal/pkg/logic/auth"
	"permission-open/internal/pkg/logic/company"
	"permission-open/internal/pkg/logic/errmsg"

	"github.com/gin-gonic/gin"
)

// CompanyInfo 企业信息
func CompanyInfo(ctx *gin.Context) {
	var req rao.CompanyInfoReq
	if err := ctx.ShouldBind(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	c := query.Use(dal.DB()).Company
	ci, err := c.WithContext(ctx).Where(c.CompanyID.Eq(req.CompanyId)).First()
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrCompanyNotFound, err.Error())
		return
	}

	response.SuccessWithData(ctx, rao.CompanyInfoResp{
		Company: &rao.Company{
			CompanyId:     ci.CompanyID,
			Name:          ci.Name,
			Logo:          ci.Logo,
			ExpireTimeSec: ci.ExpireAt.Unix(),
		},
	})
	return
}

// CompanySaveMember 新增/修改企业成员
func CompanySaveMember(ctx *gin.Context) {
	var req rao.CompanySaveMemberReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	req.Source = "SaveMember"
	_, err := company.SaveMember(ctx, jwt.GetUserIDByCtx(ctx), req)
	if err != nil {
		// todo 待优化
		if errors.Is(err, errmsg.ErrYetAccountRegister) {
			response.ErrorWithMsg(ctx, errno.ErrYetAccountRegister, err.Error())
			return
		}
		if errors.Is(err, errmsg.ErrYetNicknameRegister) {
			response.ErrorWithMsg(ctx, errno.ErrYetNicknameRegister, err.Error())
			return
		}
		if errors.Is(err, errmsg.ErrYetEmailRegister) {
			response.ErrorWithMsg(ctx, errno.ErrYetEmailRegister, err.Error())
			return
		}
		if errors.Is(err, errmsg.ErrYetEmailValid) {
			response.ErrorWithMsg(ctx, errno.ErrYetEmailValid, err.Error())
			return
		}
		if errors.Is(err, errmsg.ErrTeamSaveRepeat) {
			response.ErrorWithMsg(ctx, errno.ErrTeamSaveRepeat, err.Error())
			return
		}
		if errors.Is(err, errmsg.ErrRoleNotExists) {
			response.ErrorWithMsg(ctx, errno.ErrRoleNotExists, err.Error())
			return
		}
		if errors.Is(err, errmsg.ErrRoleForbidden) {
			response.ErrorWithMsg(ctx, errno.ErrRoleForbidden, err.Error())
			return
		}
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	response.Success(ctx)
	return
}

// ExportMembers 导入成员
func ExportMembers(ctx *gin.Context) {
	file, _, err := ctx.Request.FormFile("file")
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}
	defer file.Close()

	companyID := ctx.PostForm("company_id")
	if companyID == "" {
		response.ErrorWithMsg(ctx, errno.ErrParam, "")
		return
	}

	importDesc, err := company.ImportMembers(ctx, file, jwt.GetUserIDByCtx(ctx), companyID)
	if err != nil {
		if errors.Is(err, errmsg.ErrFileMaxLimit) {
			response.ErrorWithMsg(ctx, errno.ErrFileMaxLimit, "")
			return
		}
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	response.SuccessWithData(ctx, rao.ImportMembersResp{
		ImportDesc: importDesc,
	})
	return
}

// CompanyMembers 企业成员列表
func CompanyMembers(ctx *gin.Context) {
	var req rao.CompanyMembersReq
	if err := ctx.ShouldBind(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	members, total, err := company.MemberList(ctx, req.CompanyId, req.Keyword, req.Page, req.Size)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	response.SuccessWithData(ctx, rao.CompanyMembersResp{
		Members:       members,
		CurrentUserId: jwt.GetUserIDByCtx(ctx),
		Total:         total,
	})
	return
}

// CompanyTeams 企业团队
func CompanyTeams(ctx *gin.Context) {
	var req rao.CompanyTeamsReq
	if err := ctx.ShouldBind(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	teams, err := company.TeamsList(ctx, jwt.GetUserIDByCtx(ctx), req.CompanyId)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	response.SuccessWithData(ctx, rao.CompanyTeamsResp{
		Teams: teams,
	})
	return
}

// CompanyRemoveMember 删除企业用户
func CompanyRemoveMember(ctx *gin.Context) {
	var req rao.CompanyRemoveMemberReq
	if err := ctx.ShouldBind(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	if err := company.RemoveMember(ctx, jwt.GetUserIDByCtx(ctx), req.CompanyId, req.TargetUserID); err != nil {
		if errors.Is(err, errmsg.ErrCompanyNotRemoveMember) {
			response.ErrorWithMsg(ctx, errno.ErrCompanyNotRemoveMember, err.Error())
			return
		}
		if errors.Is(err, errmsg.ErrTeamNotRemoveCreate) {
			response.ErrorWithMsg(ctx, errno.ErrTeamNotRemoveCreate, err.Error())
			return
		}
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	response.Success(ctx)
	return
}

// CompanyUpdateMember 修改企业用户
func CompanyUpdateMember(ctx *gin.Context) {
	var req rao.CompanyUpdateMemberReq
	if err := ctx.ShouldBind(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	if err := company.UpdateMember(ctx, jwt.GetUserIDByCtx(ctx), req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	response.Success(ctx)
	return
}

// CompanyUpdatePassword 修改企业用户密码
func CompanyUpdatePassword(ctx *gin.Context) {
	var req rao.CompanyUpdatePasswordReq
	if err := ctx.ShouldBind(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	currentUserID := jwt.GetUserIDByCtx(ctx)
	if currentUserID != global.SuperManageUserID {
		response.ErrorWithMsg(ctx, errno.ErrCompanySupUpdatePassword, "")
		return
	}

	userID := req.TargetUserID
	tx := dal.GetQuery().User
	_, err := tx.WithContext(ctx).Where(tx.UserID.Eq(userID)).First()
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, "查询用户信息失败")
		return
	}

	hashedPassword, err := omnibus.GenerateBcryptFromPassword(req.NewPassword)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, "获取加密密码失败")
		return
	}
	if _, err := tx.WithContext(ctx).Where(tx.UserID.Eq(userID)).UpdateSimple(tx.Password.Value(hashedPassword)); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, "密码修改失败")
		return
	}

	// 用户需要重新登录
	if err := auth.ResetLoginUsers(ctx, userID); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, "重置 token 失败")
		return
	}

	response.Success(ctx)
	return
}
