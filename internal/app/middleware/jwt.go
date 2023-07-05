package middleware

import (
	"errors"
	"gorm.io/gorm"
	"net/http"
	"permission-open/internal/pkg/biz/consts"
	"permission-open/internal/pkg/biz/errno"
	"permission-open/internal/pkg/biz/jwt"
	"permission-open/internal/pkg/dal"
	"permission-open/internal/pkg/dal/global"

	"github.com/gin-gonic/gin"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {

		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusOK, gin.H{"code": errno.ErrMustLogin, "em": "must login", "et": "您当前尚未登录"})
			c.Abort()
			return
		}

		userID, err := jwt.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"code": errno.ErrMustLogin, "em": "must login", "et": "您当前尚未登录"})
			c.Abort()
			return
		} else if userID == "" {
			c.JSON(http.StatusOK, gin.H{"code": errno.ErrMustLogin, "em": "must login", "et": "您当前尚未登录"})
			c.Abort()
			return
		}

		// 查询token里面的用户信息是否存在于数据库
		userTable := dal.GetQuery().User
		_, err = userTable.WithContext(c).Where(userTable.UserID.Eq(userID)).First()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, gin.H{"code": errno.ErrAccountDel, "em": "user not found", "et": "用户不存在或已删除"})
			c.Abort()
			return
		}
		if err != nil {
			// 把token设置为过期
			_, _, err := jwt.GenerateTokenByTime(userID, 0)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{"code": errno.ErrInvalidToken, "em": "invalid token", "et": "登录失效，请重新登录"})
				c.Abort()
				return
			}

			c.JSON(http.StatusOK, gin.H{"code": errno.ErrMustLogin, "em": "must login", "et": "您当前尚未登录"})
			c.Abort()
			return
		}

		// 用户是否已禁用
		uc := dal.GetQuery().UserCompany
		userCompany, err := uc.WithContext(c).Where(uc.UserID.Eq(userID), uc.CompanyID.Eq(global.CompanyID)).First()
		if err != nil || userCompany.Status == consts.CompanyUserStatusDisable {
			c.JSON(http.StatusOK, gin.H{"code": errno.ErrUserDisable, "em": "user disable", "et": "您当前被禁用，请联系超管"})
			c.Abort()
			return
		}

		// 用户是否需要重新登录
		if exists, _ := dal.GetRDB().SIsMember(c, consts.RedisResetLoginUsers, userID).Result(); exists {
			c.JSON(http.StatusOK, gin.H{"code": errno.ErrUserDisable, "em": "reset login", "et": "请重新登录"})
			c.Abort()
			return
		}

		c.Set("user_id", userID)

		c.Next()
	}
}
