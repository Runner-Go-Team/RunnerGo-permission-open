package handler

import (
	"github.com/gin-gonic/gin"
	"permission-open/internal/pkg/biz/errno"
	"permission-open/internal/pkg/biz/jwt"
	"permission-open/internal/pkg/biz/response"
	"permission-open/internal/pkg/conf"
	"permission-open/internal/pkg/dal/rao"
	"permission-open/internal/pkg/logic/operation"
)

// ListOperations 操作日志列表
func ListOperations(ctx *gin.Context) {
	var req rao.ListOperationReq
	if err := ctx.ShouldBind(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	operations, total, err := operation.List(ctx, jwt.GetUserIDByCtx(ctx), req.Size, (req.Page-1)*req.Size)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	response.SuccessWithData(ctx, rao.ListOperationResp{
		Operations:    operations,
		Total:         total,
		RetentionDays: conf.Conf.CompanyOperationLogTime,
	})
	return
}
