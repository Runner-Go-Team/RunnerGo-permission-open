package consts

const (
	PerGroupCompany     = 1
	PerGroupTeam        = 2
	PerGroupRole        = 3
	PerGroupGlobalParam = 4
	PerGroupEnv         = 5
	PerGroupTestObj     = 6
	PerGroupScene       = 7
	PerGroupPlan        = 8
	PerGroupReport      = 9
	PerGroupAddressee   = 10
	PerGroupMachine     = 11
)

var PerGroupNameMap = map[int32]string{
	PerGroupCompany:     "企业成员管理",
	PerGroupTeam:        "团队管理",
	PerGroupRole:        "角色管理",
	PerGroupGlobalParam: "全局参数",
	PerGroupEnv:         "环境管理",
	PerGroupTestObj:     "测试对象",
	PerGroupScene:       "场景管理",
	PerGroupPlan:        "测试计划",
	PerGroupReport:      "测试报告",
	PerGroupAddressee:   "邮件通知",
	PerGroupMachine:     "机器管理",
}

var PerGroupMarkMap = map[int32]string{
	PerGroupCompany:     "group_company",
	PerGroupTeam:        "group_team",
	PerGroupRole:        "group_role",
	PerGroupGlobalParam: "group_global_param",
	PerGroupEnv:         "group_env",
	PerGroupTestObj:     "group_test_obj",
	PerGroupScene:       "group_scene",
	PerGroupPlan:        "group_plan",
	PerGroupReport:      "group_report",
	PerGroupAddressee:   "group_address",
	PerGroupMachine:     "group_machine",
}

const (
	PermissionTypeAuth = 1 //权限分类（1：权限  2：功能）
	PermissionTypeFunc = 2

	// RedisPermissionListUrl 权限列表 URL
	RedisPermissionListUrl = "PermissionListUrl"
)
