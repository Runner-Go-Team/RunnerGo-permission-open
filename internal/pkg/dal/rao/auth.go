package rao

type UserRelated struct {
	SettingTeamID string `json:"setting_team_id"`
	CompanyID     string `json:"company_id"`
}

type CompanyRegisterReq struct {
	Account  string `json:"account" binding:"required,min=6,max=30"`
	Password string `json:"password" binding:"required,min=6,max=20"`
}

type UpdatePermissionReq struct {
	CompanyID string `json:"company_id" binding:"required"`
}

type AuthLoginReq struct {
	Account          string `json:"account" binding:"required,min=6,max=30"`
	Password         string `json:"password" binding:"required,min=6,max=20"`
	IsAutoLogin      bool   `json:"is_auto_login"`
	InviteVerifyCode string `json:"invite_verify_code"`
}

type AuthLoginResp struct {
	IsRegister    bool   `json:"is_register"`
	ExpireTimeSec int64  `json:"expire_time_sec"`
	Token         string `json:"token"`
	TeamID        string `json:"team_id"`
	CompanyID     string `json:"company_id"`
}

type SetUserSettingsReq struct {
	UserSettings UserSettings `json:"settings"`
}

type GetUserSettingsResp struct {
	UserSettings *UserSettings `json:"settings"`
}

type UserSettings struct {
	CurrentTeamID string `json:"current_team_id" binding:"required,gt=0"`
}

type UserInfo struct {
	ID       int64  `json:"id"`
	UserID   string `json:"user_id"`
	Account  string `json:"account"`
	Email    string `json:"email"`
	Mobile   string `json:"mobile"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	RoleID   string `json:"role_id"`
	RoleName string `json:"role_name"`
}

type UserTeam struct {
	TeamID      string `json:"team_id"`
	TeamName    string `json:"team_name"`
	RoleID      string `json:"role_id"`
	RoleName    string `json:"role_name"`
	JoinTimeSec int64  `json:"join_time_sec"`
}

type GetUserInfoResp struct {
	UserInfo    *UserInfo    `json:"user_info"`
	UserRelated *UserRelated `json:"user_related"`
	TeamList    []*UserTeam  `json:"team_list"`
}

type VerifyPasswordReq struct {
	Password string `json:"password"`
}

type VerifyPasswordResp struct {
	IsMatch bool `json:"is_match"`
}

type UpdatePasswordReq struct {
	NewPassword    string `json:"new_password" binding:"required,min=6,max=20,eqfield=RepeatPassword"`
	RepeatPassword string `json:"repeat_password" binding:"required,min=6,max=20"`
}

type UpdateNicknameReq struct {
	Nickname string `json:"nickname" binding:"required,min=2"`
}

type UpdateAvatarReq struct {
	AvatarURL string `json:"avatar_url" binding:"required"`
}

type UpdateAccountReq struct {
	Account string `json:"account" binding:"required,min=6,max=30"`
}

type UpdateEmailReq struct {
	Email string `json:"email" binding:"required"`
}

type GetUserInfoReq struct {
	TargetUserID string `form:"target_user_id"`
}

type VerifyUserUsableReq struct {
	TeamID string `form:"team_id"`
}

type VerifyUserUsableResp struct {
	IsUsable bool `json:"is_usable"`
}
