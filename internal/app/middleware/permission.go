package middleware

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"permission-open/internal/pkg/biz/errno"
	"permission-open/internal/pkg/dal"
	"permission-open/internal/pkg/dal/global"
	"permission-open/internal/pkg/logic/permission"
	"strings"

	"github.com/gin-gonic/gin"
)

type Params struct {
	TeamID string `json:"team_id"`
}

func CheckPermission() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 注意这时 c.Request.Body 已经读完了，需要重新将读出来的值给放回去，之后的处理就依然可以使用 c.Request.Body 了。
		ByteBody, _ := ioutil.ReadAll(c.Request.Body)
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(ByteBody))

		var req Params
		_ = json.Unmarshal(ByteBody, &req)

		permissionUrls, _ := permission.GetPermissionListUrl(c)
		permissionFunc := strings.ToLower(c.Request.URL.Path)
		isUrlExist := false
		for _, url := range permissionUrls {
			if url == permissionFunc {
				isUrlExist = true
				break
			}
		}

		// 当前路由在权限控制范围内
		if isUrlExist {
			userID := c.GetString("user_id")

			ur := dal.GetQuery().UserRole

			// 查询企业权限
			userRole, err := ur.WithContext(c).Where(ur.UserID.Eq(userID), ur.CompanyID.Eq(global.CompanyID)).First()
			if err != nil {
				c.JSON(http.StatusOK, gin.H{"code": errno.ErrUserForbidden, "em": "forbidden", "et": "没有权限"})
				c.Abort()
				return
			}

			isHave, err := permission.CheckRolePermission(c, userRole.RoleID, permissionFunc)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{"code": errno.ErrServer, "em": "server error", "et": "server error"})
				c.Abort()
				return
			}

			// 查询团队权限
			if !isHave {
				if len(req.TeamID) > 0 {
					userRole, err = ur.WithContext(c).Where(ur.UserID.Eq(userID), ur.TeamID.Eq(req.TeamID)).First()
					if err != nil {
						c.JSON(http.StatusOK, gin.H{"code": errno.ErrUserForbidden, "em": "forbidden", "et": "没有权限"})
						c.Abort()
						return
					}

					isHave, err = permission.CheckRolePermission(c, userRole.RoleID, permissionFunc)
					if err != nil {
						c.JSON(http.StatusOK, gin.H{"code": errno.ErrServer, "em": "server error", "et": "server error"})
						c.Abort()
						return
					}
				}
			}

			if !isHave {
				c.JSON(http.StatusOK, gin.H{"code": errno.ErrUserForbidden, "em": "forbidden", "et": "没有权限"})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
