package rao

type PermissionListReq struct {
	RoleID string `form:"role_id" binding:"required"` // 角色id
}

type PermissionListResp struct {
	List []*PermissionGroup `json:"list"`
}

type Permission struct {
	IsHave         bool   `json:"is_have"`
	PermissionType int32  `json:"permission_type"`
	GroupID        int32  `json:"group_id"`
	PermissionID   int64  `json:"permission_id"`
	Title          string `json:"title"`
	Mark           string `json:"mark"`
}

type PermissionGroup struct {
	GroupID     int32         `json:"group_id"`
	GroupName   string        `json:"group_name"`
	Mark        string        `json:"mark"`
	Permissions []*Permission `json:"permissions"`
}

type SetRolePermissionReq struct {
	RoleID          string   `json:"role_id" binding:"required"` // 角色id
	RoleName        string   `json:"role_name"`
	PermissionMarks []string `json:"permission_marks" binding:"required"`
}

type UserPermissionsReq struct {
	CompanyId    string `form:"company_id" binding:"required"`
	TargetUserID string `form:"target_user_id"` // 用户id
	TeamID       string `form:"team_id"`
}

type UserPermissionsResp struct {
	Permissions []*Permission `json:"permissions"`
}

type UserAllPermissionMarksResp struct {
	Teams   map[string][]string `json:"teams"`
	Company map[string][]string `json:"company"`
}

type PermissionCheckUrlReq struct {
	UserID string `form:"user_id" binding:"required"` // 用户id
	TeamID string `form:"team_id"`
	Url    string `form:"url" binding:"required"`
}

type PermissionCheckUrlResp struct {
	IsHave bool `json:"is_have"`
}

type OpenUserPermissionMarksReq struct {
	UserId string `form:"user_id"`
}
