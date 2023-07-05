package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"permission-open/internal/pkg/biz/errno"
	"permission-open/internal/pkg/biz/jwt"
	"permission-open/internal/pkg/biz/response"
	"permission-open/internal/pkg/dal/rao"
	"permission-open/internal/pkg/logic/errmsg"
	"permission-open/internal/pkg/logic/notice"
)

func SaveNotice(ctx *gin.Context) {
	var req rao.SaveNoticeReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	noticeID, err := notice.SaveNotice(ctx, jwt.GetUserIDByCtx(ctx), req)
	if err != nil {
		if errors.Is(err, errmsg.ErrNoticeNameRepeat) {
			response.ErrorWithMsg(ctx, errno.ErrNoticeNameRepeat, "")
			return
		}
		if errors.Is(err, errmsg.ErrNoticeWebhookURLRepeat) {
			response.ErrorWithMsg(ctx, errno.ErrNoticeWebhookURLRepeat, "")
			return
		}
		if errors.Is(err, errmsg.ErrNoticeAppIDRepeat) {
			response.ErrorWithMsg(ctx, errno.ErrNoticeAppIDRepeat, "")
			return
		}
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	response.SuccessWithData(ctx, rao.SaveNoticeResp{
		NoticeId: noticeID,
	})

	return
}

func UpdateNotice(ctx *gin.Context) {
	var req rao.SaveNoticeReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	noticeID, err := notice.SaveNotice(ctx, jwt.GetUserIDByCtx(ctx), req)
	if err != nil {
		if errors.Is(err, errmsg.ErrNoticeNameRepeat) {
			response.ErrorWithMsg(ctx, errno.ErrNoticeNameRepeat, "")
			return
		}
		if errors.Is(err, errmsg.ErrNoticeWebhookURLRepeat) {
			response.ErrorWithMsg(ctx, errno.ErrNoticeWebhookURLRepeat, "")
			return
		}
		if errors.Is(err, errmsg.ErrNoticeAppIDRepeat) {
			response.ErrorWithMsg(ctx, errno.ErrNoticeAppIDRepeat, "")
			return
		}
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	response.SuccessWithData(ctx, rao.SaveNoticeResp{
		NoticeId: noticeID,
	})

	return
}

func SetStatusNotice(ctx *gin.Context) {
	var req rao.SetStatusNoticeReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	err := notice.SetStatus(ctx, jwt.GetUserIDByCtx(ctx), &req)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	response.Success(ctx)
	return
}

func ListNotice(ctx *gin.Context) {
	noticeList, err := notice.List(ctx)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	response.SuccessWithData(ctx, rao.ListNoticeResp{
		List: noticeList,
	})

	return
}

func DetailNotice(ctx *gin.Context) {
	var req rao.DetailNoticeReq
	if err := ctx.ShouldBind(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}
	noticeInfo, err := notice.Detail(ctx, req.NoticeID)
	if err != nil {
		if errors.Is(err, errmsg.ErrNoticeNotFound) {
			response.ErrorWithMsg(ctx, errno.ErrNoticeNotFound, err.Error())
			return
		}
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	response.SuccessWithData(ctx, rao.DetailNoticeResp{
		Notice: noticeInfo,
	})

	return
}

func RemoveNotice(ctx *gin.Context) {
	var req rao.RemoveNoticeReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	err := notice.RemoveNotice(ctx, jwt.GetUserIDByCtx(ctx), req.NoticeID)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	response.Success(ctx)
	return
}

func GetThirdNoticeUsers(ctx *gin.Context) {
	var req rao.GetNoticeUsersReq
	if err := ctx.ShouldBind(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}
	list, err := notice.GetThirdUsers(ctx, req.NoticeID)
	if err != nil {
		if errors.Is(err, errmsg.ErrNoticeNotFound) {
			response.ErrorWithMsg(ctx, errno.ErrNoticeNotFound, err.Error())
			return
		}
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	response.SuccessWithData(ctx, list)
	return
}
