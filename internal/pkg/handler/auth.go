package handler

import (
	"errors"
	"gorm.io/gorm"
	"permission-open/internal/pkg/public"
	"time"

	"permission-open/internal/pkg/biz/errno"
	"permission-open/internal/pkg/biz/jwt"
	"permission-open/internal/pkg/biz/log"
	"permission-open/internal/pkg/biz/response"
	"permission-open/internal/pkg/conf"
	"permission-open/internal/pkg/dal"
	"permission-open/internal/pkg/dal/rao"
	"permission-open/internal/pkg/logic/auth"
	"permission-open/internal/pkg/logic/errmsg"

	"github.com/gin-gonic/gin"
	"github.com/go-omnibus/omnibus"
)

// CompanyRegister 企业注册
func CompanyRegister(ctx *gin.Context) {
	var req rao.CompanyRegisterReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	tx := dal.GetQuery().User
	cnt, err := tx.WithContext(ctx).Where(tx.Account.Eq(req.Account)).Count()
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}
	if cnt > 0 {
		response.ErrorWithMsg(ctx, errno.ErrYetAccountRegister, "")
		return
	}

	_, err = auth.CompanyRegister(ctx, req.Account, req.Password)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	response.Success(ctx)
	return
}

func UpdatePermission(ctx *gin.Context) {
	var req rao.UpdatePermissionReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	err := auth.UpdatePermission(ctx, req.CompanyID)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	response.Success(ctx)
	return
}

// AuthLogin 登录
func AuthLogin(ctx *gin.Context) {
	var req rao.AuthLoginReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	userInfo, err := auth.Login(ctx, req)
	if err != nil {
		if errors.Is(err, errmsg.ErrAccountNotFound) {
			response.ErrorWithMsg(ctx, errno.ErrAccountNotFound, "")
			return
		}
		if errors.Is(err, errmsg.ErrPasswordFailed) {
			response.ErrorWithMsg(ctx, errno.ErrPasswordFailed, "")
			return
		}
		if errors.Is(err, errmsg.ErrUserDisable) {
			response.ErrorWithMsg(ctx, errno.ErrUserDisable, "")
			return
		}
		if errors.Is(err, errmsg.ErrUserDisable) {
			response.ErrorWithMsg(ctx, errno.ErrUserDisable, "")
			return
		}
		response.ErrorWithMsg(ctx, errno.ErrAuthFailed, "")
		return
	}

	uc := dal.GetQuery().UserCompany
	userCompany, err := uc.WithContext(ctx).Where(uc.UserID.Eq(userInfo.UserID)).First()
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrAuthFailed, "")
		return
	}

	// 开始生成token
	expireTime := conf.Conf.DefaultTokenExpireTime
	d := expireTime * time.Hour
	if req.IsAutoLogin {
		d = 30 * 24 * time.Hour
	}

	token, exp, err := jwt.GenerateTokenByTime(userInfo.UserID, d)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrInvalidToken, err.Error())
		return
	}

	if err := auth.UpdateLoginTime(ctx, userInfo.UserID); err != nil {
		log.Logger.Errorf("update login time err %s", err)
	}

	defaultTeamID := ""
	if req.InviteVerifyCode == "" {
		defaultTeamID, _ = auth.GetAvailTeamID(ctx, userInfo.UserID)
		userSettings := rao.UserSettings{
			CurrentTeamID: defaultTeamID,
		}
		_ = auth.SetUserSettings(ctx, userInfo.UserID, &userSettings)
	}

	if err := auth.RemoveResetLoginUser(ctx, userInfo.UserID); err != nil {
		log.Logger.Errorf("remove remove login err %s", err)
	}

	response.SuccessWithData(ctx, rao.AuthLoginResp{
		Token:         token,
		ExpireTimeSec: exp.Unix(),
		TeamID:        defaultTeamID,
		IsRegister:    true,
		CompanyID:     userCompany.CompanyID,
	})
	return
}

// RefreshToken 续期
func RefreshToken(ctx *gin.Context) {
	tokenString := ctx.GetHeader("Authorization")

	token, exp, err := jwt.RefreshToken(tokenString)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	response.SuccessWithData(ctx, rao.AuthLoginResp{
		Token:         token,
		ExpireTimeSec: exp.Unix(),
	})
	return
}

func VerifyPassword(ctx *gin.Context) {
	var req rao.VerifyPasswordReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	tx := dal.GetQuery().User
	u, err := tx.WithContext(ctx).Where(tx.UserID.Eq(jwt.GetUserIDByCtx(ctx))).First()
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	err = omnibus.CompareBcryptHashAndPassword(u.Password, req.Password)

	response.SuccessWithData(ctx, rao.VerifyPasswordResp{IsMatch: err == nil})
	return
}

func UpdatePassword(ctx *gin.Context) {
	var req rao.UpdatePasswordReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	if req.NewPassword != req.RepeatPassword {
		response.ErrorWithMsg(ctx, errno.ErrParam, "两次密码输入不一致")
		return
	}

	userID := jwt.GetUserIDByCtx(ctx)

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

	response.Success(ctx)
	return
}

func UpdateNickname(ctx *gin.Context) {
	var req rao.UpdateNicknameReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	u := dal.GetQuery().User
	_, err := u.WithContext(ctx).Where(u.Nickname.Eq(req.Nickname)).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		if _, err := u.WithContext(ctx).Where(u.UserID.Eq(jwt.GetUserIDByCtx(ctx))).UpdateColumn(u.Nickname, req.Nickname); err != nil {
			response.ErrorWithMsg(ctx, errno.ErrServer, "update account failed")
			return
		}

		response.Success(ctx)
		return
	}

	response.ErrorWithMsg(ctx, errno.ErrYetNicknameRegister, "")
	return
}

func UpdateAvatar(ctx *gin.Context) {
	var req rao.UpdateAvatarReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	tx := dal.GetQuery().User
	if _, err := tx.WithContext(ctx).Where(tx.UserID.Eq(jwt.GetUserIDByCtx(ctx))).UpdateColumn(tx.Avatar, req.AvatarURL); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, "password = new password")
		return
	}

	response.Success(ctx)
	return
}

func UpdateAccount(ctx *gin.Context) {
	var req rao.UpdateAccountReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	tx := dal.GetQuery().User
	_, err := tx.WithContext(ctx).Where(tx.Account.Eq(req.Account)).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		if _, err := tx.WithContext(ctx).Where(tx.UserID.Eq(jwt.GetUserIDByCtx(ctx))).UpdateColumn(tx.Account, req.Account); err != nil {
			response.ErrorWithMsg(ctx, errno.ErrServer, "update account failed")
			return
		}

		response.Success(ctx)
		return
	}

	response.ErrorWithMsg(ctx, errno.ErrYetAccountRegister, "update account failed")
	return
}

func UpdateEmail(ctx *gin.Context) {
	var req rao.UpdateEmailReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	if !public.IsEmailValid(req.Email) {
		response.ErrorWithMsg(ctx, errno.ErrYetEmailValid, "")
		return
	}

	u := dal.GetQuery().User
	_, err := u.WithContext(ctx).Where(u.Email.Eq(req.Email)).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		if _, err := u.WithContext(ctx).Where(u.UserID.Eq(jwt.GetUserIDByCtx(ctx))).UpdateColumn(u.Email, req.Email); err != nil {
			response.ErrorWithMsg(ctx, errno.ErrServer, "update account failed")
			return
		}

		response.Success(ctx)
		return
	}

	response.ErrorWithMsg(ctx, errno.ErrYetEmailRegister, "")
	return
}

func GetUserInfo(ctx *gin.Context) {
	var req rao.GetUserInfoReq
	if err := ctx.ShouldBind(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}
	res, err := auth.GetUserInfo(ctx, jwt.GetUserIDByCtx(ctx), req.TargetUserID)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	response.SuccessWithData(ctx, res)
	return
}

// GetUserSettings 获取用户配置
func GetUserSettings(ctx *gin.Context) {
	res, err := auth.GetUserSettings(ctx, jwt.GetUserIDByCtx(ctx))
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	response.SuccessWithData(ctx, res)
	return
}

// SetUserSettings 设置用户配置
func SetUserSettings(ctx *gin.Context) {
	var req rao.SetUserSettingsReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	if err := auth.SetUserSettings(ctx, jwt.GetUserIDByCtx(ctx), &req.UserSettings); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, "")
		return
	}

	response.Success(ctx)
	return
}

// VerifyUserUsable 验证用户是否是有效用户
func VerifyUserUsable(ctx *gin.Context) {
	var req rao.VerifyUserUsableReq
	if err := ctx.ShouldBind(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	isUsable := true
	userID := jwt.GetUserIDByCtx(ctx)
	u := dal.GetQuery().User
	if _, err := u.WithContext(ctx).Where(u.UserID.Eq(jwt.GetUserIDByCtx(ctx))).First(); err != nil {
		isUsable = false
	}
	if len(req.TeamID) > 0 {
		t := dal.GetQuery().Team
		if _, err := t.WithContext(ctx).Where(t.TeamID.Eq(req.TeamID)).First(); err != nil {
			isUsable = false
		}
		ut := dal.GetQuery().UserTeam
		if _, err := ut.WithContext(ctx).Where(ut.UserID.Eq(userID), ut.TeamID.Eq(req.TeamID)).First(); err != nil {
			isUsable = false
		}
	}

	response.SuccessWithData(ctx, rao.VerifyUserUsableResp{IsUsable: isUsable})
	return
}
