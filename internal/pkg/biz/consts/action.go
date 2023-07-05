package consts

const (
	ActionCompanySaveMember    = "company_save_member"
	ActionCompanyExportMember  = "company_export_member"
	ActionCompanyUpdateMember  = "company_update_member"
	ActionCompanyRemoveMember  = "company_remove_member"
	ActionCompanySetRoleMember = "company_set_role_member"

	ActionTeamSave          = "team_save"
	ActionTeamUpdate        = "team_update"
	ActionTeamSaveMember    = "team_save_member"
	ActionTeamRemoveMember  = "team_remove_member"
	ActionTeamSetRoleMember = "team_set_role_member"
	ActionTeamDisband       = "team_disband"
	ActionTeamTransfer      = "team_transfer"

	ActionRoleSave   = "role_save"
	ActionRoleSet    = "role_set"
	ActionRoleRemove = "role_remove"

	ActionNoticeSave      = "notice_save"
	ActionNoticeUpdate    = "notice_update"
	ActionNoticeSetStatus = "notice_set_status"
	ActionNoticeRemove    = "notice_remove"

	ActionNoticeGroupSave   = "notice_group_save"
	ActionNoticeGroupUpdate = "notice_group_update"
	ActionNoticeGroupRemove = "notice_group_remove"
)

var ActionMap = map[string]int32{
	ActionCompanySaveMember:    101,
	ActionCompanyExportMember:  102,
	ActionCompanyUpdateMember:  103,
	ActionCompanyRemoveMember:  104,
	ActionCompanySetRoleMember: 105,
	ActionTeamSave:             201,
	ActionTeamUpdate:           202,
	ActionTeamSaveMember:       203,
	ActionTeamRemoveMember:     204,
	ActionTeamSetRoleMember:    205,
	ActionTeamDisband:          206,
	ActionTeamTransfer:         207,
	ActionRoleSave:             301,
	ActionRoleSet:              302,
	ActionRoleRemove:           303,
	ActionNoticeSave:           401,
	ActionNoticeUpdate:         402,
	ActionNoticeSetStatus:      403,
	ActionNoticeRemove:         404,
	ActionNoticeGroupSave:      501,
	ActionNoticeGroupUpdate:    502,
	ActionNoticeGroupRemove:    503,
}

var ActionNameMap = map[int32]string{
	101: "创建成员",
	102: "批量导入成员",
	103: "编辑成员",
	104: "删除成员",
	105: "更改企业角色",
	201: "新建团队",
	202: "编辑团队",
	203: "添加团队成员",
	204: "移除团队成员",
	205: "更改团队角色",
	206: "解散团队",
	207: "移交团队管理员",
	301: "新建角色",
	302: "设置角色权限",
	303: "删除角色",
	401: "新建三方通知",
	402: "修改三方通知",
	403: "禁用|启用三方通知",
	404: "删除三方通知",
	501: "新建三方组通知",
	502: "修改三方组通知",
	503: "删除三方组通知",
}
