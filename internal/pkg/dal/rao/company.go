package rao

type Company struct {
	CompanyId     string `json:"company_id"`
	Name          string `json:"name"`
	Logo          string `json:"logo"`
	ExpireTimeSec int64  `json:"expire_time_sec"`
}

type CompanyMember struct {
	UserID                  string `json:"user_id"`
	Avatar                  string `json:"avatar"`
	Nickname                string `json:"nickname"`
	Account                 string `json:"account"`
	RoleId                  string `json:"role_id"`
	RoleName                string `json:"role_name"`
	InviteUserID            string `json:"invite_user_id"`
	InviteUserName          string `json:"invite_user_name"`
	RoleLevel               int32  `json:"role_level"`
	Status                  int32  `json:"status"`
	InviteTimeSec           int64  `json:"invite_time_sec"`
	CreatedTimeSec          int64  `json:"created_time_sec"`
	IsOperableCompanyRole   bool   `json:"is_operable_company_role"`   // 当前用户是否允许操作该成员企业角色
	IsOperableRemoveMember  bool   `json:"is_operable_remove_member"`  // 当前用户是否允许移除当前用户
	IsOperableDisableMember bool   `json:"is_operable_disable_member"` // 当前用户是否允许禁用当前用户
	IsOperableUptPassword   bool   `json:"is_operable_upt_password"`   // 当前用户是否允许修改密码
}

type CompanyMemberTeam struct {
	*Team
	RoleId         string `json:"role_id"`
	RoleName       string `json:"role_name"`
	InviteUserID   string `json:"invite_user_id"`
	InviteUserName string `json:"invite_user_name"`
}

type TeamDetail struct {
	TeamId string `json:"team_id"`
	RoleId string `json:"role_id"`
}

type CompanySaveMemberReq struct {
	CompanyId  string        `json:"company_id" binding:"required"`
	Account    string        `json:"account" binding:"required,min=6,max=30"`
	Password   string        `json:"password" binding:"required,min=6,max=20"`
	Nickname   string        `json:"nickname" binding:"required"`
	RoleId     string        `json:"role_id" binding:"required"`
	Email      string        `json:"email"`
	Source     string        `json:"source"`
	TeamDetail []*TeamDetail `json:"team_detail"`
}

type CompanyMembersReq struct {
	CompanyId string `form:"company_id" binding:"required"`
	Keyword   string `form:"keyword"`
	Page      int    `form:"page,default=1"`
	Size      int    `form:"size,default=20"`
}

type CompanyMembersResp struct {
	Members       []*CompanyMember `json:"members"`
	CurrentUserId string           `json:"current_user_id"`
	Total         int64            `json:"total"`
}

type ImportMembersResp struct {
	*ImportDesc
}

type ImportDesc struct {
	SuccessNum    int64       `json:"success_num"`
	FailNum       int64       `json:"fail_num"`
	Path          string      `json:"path"`
	ImportErrDesc []ImportErr `json:"import_err_desc"`
}

type ImportErr struct {
	Nickname string `json:"nickname"`
	Account  string `json:"account"`
	Password string `json:"password"`
	RoleName string `json:"role_name"`
	Email    string `json:"email"`
	ErrMsg   string `json:"err_msg"`
}

type CompanyTeamsReq struct {
	CompanyId string `form:"company_id" binding:"required"`
}

type CompanyTeamsResp struct {
	Teams []*Team `json:"teams"`
}

type CompanyRemoveMemberReq struct {
	CompanyId    string `json:"company_id" binding:"required"`
	TargetUserID string `json:"target_user_id" binding:"required"`
}

type CompanyUpdateMemberReq struct {
	Status       int32  `json:"status" binding:"required"`
	TargetUserID string `json:"target_user_id" binding:"required"`
	CompanyId    string `json:"company_id" binding:"required"`
}

type CompanyUpdatePasswordReq struct {
	NewPassword  string `json:"new_password" binding:"required,min=6,max=20"`
	TargetUserID string `json:"target_user_id" binding:"required"`
	CompanyId    string `json:"company_id" binding:"required"`
}

type CompanyInfoReq struct {
	CompanyId string `form:"company_id" binding:"required"`
}

type CompanyInfoResp struct {
	Company *Company `json:"company"`
}
