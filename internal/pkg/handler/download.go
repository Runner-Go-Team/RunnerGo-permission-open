package handler

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/gin-gonic/gin"

	"permission-open/internal/pkg/biz/errno"
	"permission-open/internal/pkg/biz/response"
)

const (
	CompanyExportExcel = "RunnerGo创建成员批量导入模板.xlsx"
)

func DownloadCompanyExport(ctx *gin.Context) {

	filename := "./static/files/" + CompanyExportExcel
	f, err := os.Open(filename)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrServer, err.Error())
		return
	}
	defer f.Close()

	// 将文件读取出来
	data, err := ioutil.ReadAll(f)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrServer, err.Error())
		return
	}
	// 设置头信息：Content-Disposition ，消息头指示回复的内容该以何种形式展示，
	// 是以内联的形式（即网页或者页面的一部分），还是以附件的形式下载并保存到本地
	// Content-Disposition: inline
	// Content-Disposition: attachment
	// Content-Disposition: attachment; filename="filename.后缀"
	// 第一个参数或者是inline（默认值，表示回复中的消息体会以页面的一部分或者
	// 整个页面的形式展示），或者是attachment（意味着消息体应该被下载到本地；
	// 大多数浏览器会呈现一个“保存为”的对话框，将filename的值预填为下载后的文件名，
	// 假如它存在的话）。
	ctx.Header("Content-Disposition",
		fmt.Sprintf(`attachment; filename="%s"`, CompanyExportExcel))
	ctx.Writer.Write(data)
}
