package rao

type ListOperationReq struct {
	Page int `form:"page,default=1"`
	Size int `form:"size,default=20"`
}

type ListOperationResp struct {
	Operations    []*Operation `json:"operations"`
	Total         int64        `json:"total"`
	RetentionDays int          `json:"retention_days"`
}

type Operation struct {
	UserID         string `json:"user_id"`
	UserName       string `json:"user_name"`
	UserAvatar     string `json:"user_avatar"`
	Category       int32  `json:"category"`
	Action         int32  `json:"action"`
	Name           string `json:"name"`
	CreatedTimeSec int64  `json:"created_time_sec"`
}
