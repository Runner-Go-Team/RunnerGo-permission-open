package rao

type Team struct {
	Type            int32     `json:"type"`
	TeamID          string    `json:"team_id"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	CreatedUserID   string    `json:"created_user_id"`
	CreatedUserName string    `json:"created_user_name,omitempty"`
	CreatedTimeSec  int64     `json:"created_time_sec"`
	UpdatedTimeSec  int64     `json:"updated_time_sec"`
	Members         []*Member `json:"members,omitempty"`
}

type RelatedTeam struct {
	IsCollect bool `json:"is_collect"`
	*Team
}

type Member struct {
	UserID                 string `json:"user_id"`
	Account                string `json:"account"`
	Mobile                 string `json:"mobile"`
	Avatar                 string `json:"avatar"`
	Email                  string `json:"email"`
	Nickname               string `json:"nickname"`
	CompanyRoleID          string `json:"company_role_id"`
	CompanyRoleName        string `json:"company_role_name"`
	TeamRoleID             string `json:"team_role_id"`
	TeamRoleName           string `json:"team_role_name"`
	InviteUserID           string `json:"invite_user_id"`
	InviteUserName         string `json:"invite_user_name"`
	Status                 int32  `json:"status"`
	JoinTimeSec            int64  `json:"join_time_sec"`
	IsOperableCompanyRole  bool   `json:"is_operable_company_role"`  // 当前用户是否允许操作该成员企业角色
	IsOperableTeamRole     bool   `json:"is_operable_team_role"`     // 当前用户是否允许操作该成员团队角色
	IsOperableRemoveMember bool   `json:"is_operable_remove_member"` // 当前用户是否允许移除当前角色
	IsTransferSuperTeam    bool   `json:"is_transfer_super_team"`    // 当前用户是否有移交团队管理员的权利
}

type SaveMembers struct {
	UserID     string `json:"user_id" binding:"required"`
	TeamRoleID string `json:"team_role_id" binding:"required"`
}

type SaveTeamReq struct {
	TeamType  int32  `json:"team_type" binding:"required"`
	Name      string `json:"name" binding:"required"`
	CompanyId string `json:"company_id" binding:"required"`
}

type SaveTeamResp struct {
	TeamID string `json:"team_id"  binding:"required"`
}

type UpdateTeamReq struct {
	TeamType    int32   `json:"team_type"`
	TeamID      string  `json:"team_id"  binding:"required"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

type ListTeamReq struct {
	Order     int32  `form:"order"`
	CompanyId string `form:"company_id" binding:"required"`
	Keyword   string `form:"keyword"`
}

type ListTeamResp struct {
	Teams []*RelatedTeam `json:"teams"`
}

type RemoveTeamMemberReq struct {
	TeamID       string `json:"team_id" binding:"required,gt=0"`
	TargetUserID string `json:"target_user_id" binding:"required"`
}

type DisbandTeamReq struct {
	TeamID string `json:"team_id"`
}

type CollectionTeamReq struct {
	Status int32  `json:"status"` // 1:收藏   2:取消收藏
	TeamID string `json:"team_id"`
}

type TeamInfoReq struct {
	TeamID  string `form:"team_id" binding:"required"`
	Keyword string `form:"keyword"`
	Page    int    `form:"page,default=1"`
	Size    int    `form:"size,default=20"`
}

type TeamInfoResp struct {
	Team  *Team `json:"team"`
	Total int64 `json:"total"`
}

type SaveTeamMemberReq struct {
	UserID  string        `json:"user_id"`
	TeamID  string        `json:"team_id" binding:"required"`
	Members []SaveMembers `json:"members" binding:"required"`
}

type TeamCompanyMembersReq struct {
	UserID  string `form:"user_id"`
	TeamID  string `form:"team_id" binding:"required"`
	Keyword string `form:"keyword"`
	Page    int    `form:"page,default=1"`
	Size    int    `form:"size,default=20"`
}

type TeamCompanyMembersResp struct {
	Members      []*Member `json:"members"`
	CreateUserID string    `json:"create_user_id"`
	Total        int64     `json:"total"`
}

type TeamTransferSuperRoleReq struct {
	TeamID       string `json:"team_id" binding:"required"`
	TargetUserID string `json:"target_user_id" binding:"required"`
}
