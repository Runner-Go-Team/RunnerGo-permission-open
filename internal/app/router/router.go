package router

import (
	"permission-open/internal/app/middleware"
	"time"

	"permission-open/internal/pkg/handler"

	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/go-omnibus/proof"
)

func RegisterRouter(r *gin.Engine) {
	// cors
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"POST", "GET", "OPTIONS", "DELETE", "PUT", "PATCH"},
		AllowHeaders:     []string{"Authorization", "Content-Type", "Upgrade", "Origin", "Connection", "Accept-Encoding", "Accept-Language", "Host", "x-requested-with", "CurrentTeamID"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.Use(ginzap.Ginzap(proof.Logger.Z, time.RFC3339, true))

	r.Use(ginzap.RecoveryWithZap(proof.Logger.Z, true))

	// routers
	api := r.Group("/permission/api")

	// 用户鉴权
	auth := api.Group("/v1/auth/")
	auth.POST("user_login", handler.AuthLogin)             // 用户登录
	auth.POST("company_register", handler.CompanyRegister) // 创建企业及超级管理员
	//auth.POST("update_permission", handler.UpdatePermission) // 更新企业权限（调试专用）
	//auth.GET("clear", handler.Clear) // 清洗历史数据

	// 无需 jwt
	open := api.Group("/open")
	open.GET("/v1/team/company/members", handler.TeamCompanyMembers)          // 当前团队和企业成员关系
	open.POST("/v1/team/member/save", handler.SaveTeamMembers)                // 添加团队成员
	open.GET("/v1/role/member/info", handler.RoleMember)                      // 我的角色信息
	open.GET("/v1/permission/user/get_marks", handler.UserAllPermissionMarks) // 获取用户的全部角色对应的mark
	open.GET("/v1/notice/group/list", handler.ListNoticeGroup)                // 三方组通知列表
	open.GET("/v1/notice/get_third_users", handler.GetThirdNoticeUsers)       // 获取三方组织架构

	// 设置静态文件路径
	r.Static("/static", "./static/")

	// 开启接口鉴权
	api.Use(middleware.JWT())
	api.Use(middleware.CheckPermission()) // 权限

	//api.POST("/v1/checkUrl", handler.PermissionCheckUrl)        // 通过 userID url 判断是否有权限

	// 用户
	user := api.Group("/v1/user/")
	user.GET("get", handler.GetUserInfo) // 获取用户信息
	user.POST("update_password", handler.UpdatePassword)
	user.POST("update_nickname", handler.UpdateNickname)
	user.POST("update_avatar", handler.UpdateAvatar)
	user.POST("update_account", handler.UpdateAccount)
	user.POST("update_email", handler.UpdateEmail)
	user.POST("verify_password", handler.VerifyPassword)
	user.GET("verify_usable", handler.VerifyUserUsable) // 验证用户是否是有效用户

	// 用户配置
	setting := api.Group("/v1/setting/")
	setting.GET("get", handler.GetUserSettings)
	setting.POST("set", handler.SetUserSettings)

	// 企业
	company := api.Group("/v1/company/")
	company.GET("info", handler.CompanyInfo)                              // 企业信息
	company.GET("members", handler.CompanyMembers)                        // 企业成员列表
	company.GET("teams", handler.CompanyTeams)                            // 企业团队
	company.POST("member/save", handler.CompanySaveMember)                // 新增企业成员
	company.POST("member/export", handler.ExportMembers)                  // 导入企业成员
	company.POST("member/remove", handler.CompanyRemoveMember)            // 删除企业用户
	company.POST("member/update", handler.CompanyUpdateMember)            // 修改设置企业用户状态  1:正常  2:禁用
	company.POST("member/update_password", handler.CompanyUpdatePassword) // 修改企业成员密码

	// 团队
	team := api.Group("/v1/team/")
	team.POST("save", handler.SaveTeam)                     // 新建团队
	team.POST("update", handler.UpdateTeam)                 // 修改团队
	team.GET("list", handler.ListTeam)                      // 团队列表   order (1:我加入的团队 2:我是团队管理员的团队 3:我收藏的团队)
	team.GET("info", handler.TeamInfo)                      // 团队信息 - 团队成员列表
	team.POST("disband", handler.DisbandTeam)               // 解散团队
	team.POST("collection", handler.CollectionTeam)         // 收藏/取消收藏团队
	team.GET("company/members", handler.TeamCompanyMembers) // 当前团队和企业成员关系

	team.POST("member/save", handler.SaveTeamMembers)               // 添加团队成员
	team.POST("transfer_super_role", handler.TeamTransferSuperRole) // 移交团队管理员
	team.POST("member/remove", handler.RemoveTeamMember)            // 移除用户团队

	// 角色
	role := api.Group("/v1/role/")
	role.GET("list", handler.RoleList)          // 角色列表
	role.POST("save", handler.SaveRole)         // 新建角色
	role.POST("remove", handler.RemoveRole)     // 删除角色
	role.GET("is_remove", handler.IsRemoveRole) // 判断能否删除角色

	role.GET("member/info", handler.RoleMember)            // 我的角色信息
	role.GET("member/list", handler.RoleMembers)           // 角色成员列表
	role.POST("company/set", handler.SetCompanyMemberRole) // 更改企业角色
	role.POST("team/set", handler.SetTeamMemberRole)       // 更改团队角色

	// 权限
	permission := api.Group("/v1/permission/")
	permission.GET("list", handler.PermissionList)         // 权限列表
	permission.POST("role/set", handler.SetRolePermission) // 设置角色权限

	permission.GET("user/get", handler.UserPermissions)              // 用户权限列表
	permission.GET("user/get_marks", handler.UserAllPermissionMarks) // 获取用户的全部角色对应的mark

	// 操作日志
	operation := api.Group("/v1/operation")
	operation.GET("/list", handler.ListOperations)

	// clients manage
	r.POST("/management/api/company/get_newest_stress_plan_list", handler.ManagerGetNewestStressPlanList) // 获取团队最新性能计划列表
	r.POST("/management/api/company/get_newest_auto_plan_list", handler.ManagerGetNewestAutoPlanList)     // 获取团队最新自动计划列表

	// 三方通知
	notice := api.Group("/v1/notice")
	notice.POST("save", handler.SaveNotice)                    // 新建三方通知
	notice.POST("update", handler.UpdateNotice)                // 修改三方通知
	notice.POST("set_status", handler.SetStatusNotice)         // 禁用|启用三方通知
	notice.GET("list", handler.ListNotice)                     // 三方通知列表
	notice.GET("detail", handler.DetailNotice)                 // 三方通知详情
	notice.POST("remove", handler.RemoveNotice)                // 删除
	notice.GET("get_third_users", handler.GetThirdNoticeUsers) // 获取三方组织架构
	// 三方通知组
	noticeGroup := api.Group("/v1/notice/group")
	noticeGroup.POST("save", handler.SaveNoticeGroup)     // 新建三方组通知
	noticeGroup.POST("update", handler.UpdateNoticeGroup) // 修改三方组通知
	noticeGroup.GET("list", handler.ListNoticeGroup)      // 三方组通知列表
	noticeGroup.GET("detail", handler.DetailNoticeGroup)  // 三方通知组详情
	noticeGroup.POST("remove", handler.RemoveNoticeGroup) // 删除三方通知组
}
