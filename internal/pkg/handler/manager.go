package handler

import (
	"github.com/gin-gonic/gin"
	"permission-open/internal/pkg/biz/errno"
	"permission-open/internal/pkg/biz/response"
	"permission-open/internal/pkg/dal/clients/manager"
	"permission-open/internal/pkg/dal/rao"
)

func ManagerGetNewestStressPlanList(ctx *gin.Context) {
	var req rao.TressPlanListResp
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	list, err := manager.GetNewestStressPlanList(&req)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrServer, err.Error())
		return
	}

	response.SuccessWithData(ctx, list)
	return
}

func ManagerGetNewestAutoPlanList(ctx *gin.Context) {
	var req rao.AutoPlanListResp
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	list, err := manager.GetNewestAutoPlanList(&req)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrServer, err.Error())
		return
	}

	response.SuccessWithData(ctx, list)
	return
}
