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

func SaveNoticeGroup(ctx *gin.Context) {
	var req rao.SaveNoticeGroupReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	if len(req.NoticeRelates) <= 0 {
		response.ErrorWithMsg(ctx, errno.ErrNoticeRelateNotNull, "")
		return
	}

	groupID, err := notice.SaveNoticeGroup(ctx, jwt.GetUserIDByCtx(ctx), req)
	if err != nil {
		if errors.Is(err, errmsg.ErrNoticeGroupNameRepeat) {
			response.ErrorWithMsg(ctx, errno.ErrNoticeGroupNameRepeat, "")
			return
		}
		if errors.Is(err, errmsg.ErrYetEmailValid) {
			response.ErrorWithMsg(ctx, errno.ErrYetEmailValid, "")
			return
		}
		if errors.Is(err, errmsg.ErrYetEmailNotFound) {
			response.ErrorWithMsg(ctx, errno.ErrYetEmailNotFound, "")
			return
		}
		if errors.Is(err, errmsg.ErrYetUserNotFound) {
			response.ErrorWithMsg(ctx, errno.ErrYetUserNotFound, "")
			return
		}
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	response.SuccessWithData(ctx, rao.SaveNoticeGroupResp{
		GroupID: groupID,
	})

	return
}

func UpdateNoticeGroup(ctx *gin.Context) {
	var req rao.UpdateNoticeGroupReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	if len(req.NoticeRelates) <= 0 {
		response.ErrorWithMsg(ctx, errno.ErrNoticeRelateNotNull, "")
		return
	}

	groupID, err := notice.UpdateNoticeGroup(ctx, jwt.GetUserIDByCtx(ctx), req)
	if err != nil {
		if errors.Is(err, errmsg.ErrNoticeGroupNameRepeat) {
			response.ErrorWithMsg(ctx, errno.ErrNoticeGroupNameRepeat, "")
			return
		}
		if errors.Is(err, errmsg.ErrYetEmailValid) {
			response.ErrorWithMsg(ctx, errno.ErrYetEmailValid, "")
			return
		}
		if errors.Is(err, errmsg.ErrYetEmailNotFound) {
			response.ErrorWithMsg(ctx, errno.ErrYetEmailNotFound, "")
			return
		}
		if errors.Is(err, errmsg.ErrYetUserNotFound) {
			response.ErrorWithMsg(ctx, errno.ErrYetUserNotFound, "")
			return
		}
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	response.SuccessWithData(ctx, rao.UpdateNoticeGroupResp{
		GroupID: groupID,
	})

	return
}

func ListNoticeGroup(ctx *gin.Context) {
	var req rao.ListNoticeGroupReq
	if err := ctx.ShouldBind(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	noticeList, err := notice.ListGroup(ctx, &req)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	response.SuccessWithData(ctx, rao.ListNoticeGroupResp{
		List: noticeList,
	})

	return
}

func DetailNoticeGroup(ctx *gin.Context) {
	var req rao.DetailNoticeGroupReq
	if err := ctx.ShouldBind(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}
	groupInfo, err := notice.DetailGroup(ctx, req.GroupID)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	response.SuccessWithData(ctx, rao.DetailNoticeGroupResp{
		Group: groupInfo,
	})

	return
}

func RemoveNoticeGroup(ctx *gin.Context) {
	var req rao.RemoveNoticeGroupReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	err := notice.RemoveNoticeGroup(ctx, jwt.GetUserIDByCtx(ctx), req.GroupID)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	response.Success(ctx)
	return
}
