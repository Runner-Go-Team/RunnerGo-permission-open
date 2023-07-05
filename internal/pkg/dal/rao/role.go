package rao

// RoleAttr 角色属性
type RoleAttr struct {
	IsUpdatePermission bool `json:"is_update_permission"`
}

type RoleListReq struct {
	RoleType  int32  `form:"role_type"`
	Level     int32  `form:"level"`
	CompanyId string `form:"company_id" binding:"required"`
}

type RoleListResp struct {
	RoleList []*Role `json:"role_list"`
}

type Role struct {
	RoleType  int32     `json:"role_type"`      // 角色分类（1：企业  2：团队）
	Level     int32     `json:"level"`          // 角色层级
	IsDefault int32     `json:"is_default"`     // 是否是默认角色
	RoleID    string    `json:"role_id"`        // 角色id
	Name      string    `json:"name"`           // 角色名称
	Attr      *RoleAttr `json:"attr,omitempty"` // 角色属性
}

type SaveRoleReq struct {
	RoleType  int32  `json:"role_type" binding:"required"`
	CompanyID string `json:"company_id" binding:"required"`
	Name      string `json:"name" binding:"required"`
}

type SetCompanyMemberRoleReq struct {
	RoleID       string `json:"role_id" binding:"required"`
	TargetUserID string `json:"target_user_id" binding:"required"`
	CompanyID    string `json:"company_id" binding:"required"`
}

type SetTeamMemberRoleReq struct {
	RoleID       string `json:"role_id" binding:"required"`
	TargetUserID string `json:"target_user_id" binding:"required"`
	TeamID       string `json:"team_id" binding:"required"`
}

type RoleMembersReq struct {
	RoleID  string `form:"role_id" binding:"required"`
	TeamID  string `form:"team_id"`
	Keyword string `form:"keyword"`
	Page    int    `form:"page,default=1"`
	Size    int    `form:"size,default=20"`
}

type RoleMembersResp struct {
	Members []*RoleMember `json:"members"`
	Total   int64         `json:"total"`
}

type RoleMemberReq struct {
	UserID    string `form:"user_id"`
	CompanyID string `form:"company_id"`
	TeamID    string `form:"team_id"`
	RoleID    string `form:"role_id"`
}

type RoleMemberResp struct {
	Role        *Role   `json:"role"`
	UsableRoles []*Role `json:"usable_roles"`
}

type RoleMember struct {
	UserID         string `json:"user_id"`
	Account        string `json:"account"`
	Mobile         string `json:"mobile"`
	Avatar         string `json:"avatar"`
	Email          string `json:"email"`
	Nickname       string `json:"nickname"`
	RoleID         string `json:"role_id"`
	RoleName       string `json:"role_name"`
	InviteUserID   string `json:"invite_user_id"`
	InviteUserName string `json:"invite_user_name"`
	InviteTimeSec  int64  `json:"invite_time_sec"`
	IsOperableRole bool   `json:"is_operable_role"` // 当前用户是否允许操作该成员角色
	TeamID         string `json:"team_id,omitempty"`
	TeamName       string `json:"team_name,omitempty"`
}

type RemoveRoleReq struct {
	RoleID       string `json:"role_id" binding:"required"`
	CompanyID    string `json:"company_id" binding:"required"`
	ChangeRoleID string `json:"change_role_id" binding:"required"`
}

type IsRemoveRoleReq struct {
	RoleID string `form:"role_id" binding:"required"`
}

type IsRemoveRoleResp struct {
	IsAllowRemove bool `json:"is_allow_remove"`
}
