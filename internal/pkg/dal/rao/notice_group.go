package rao

type NoticeGroupRelate struct {
	NoticeID string                   `json:"notice_id"`
	Params   *NoticeGroupRelateParams `json:"params"`
}

type NoticeGroupRelateParams struct {
	UserIDs []string `json:"user_ids,omitempty"`
	Emails  []string `json:"emails,omitempty"`
}

type NoticeGroup struct {
	GroupID        string               `json:"group_id"`
	Name           string               `json:"name"`
	Notice         []*Notice            `json:"notice,omitempty"`
	NoticeRelates  []*NoticeGroupRelate `json:"notice_relates,omitempty"`
	CreatedTimeSec int64                `json:"created_time_sec"`
	UpdatedTimeSec int64                `json:"updated_time_sec"`
}

type SaveNoticeGroupReq struct {
	Name          string               `json:"name" binding:"required"`
	NoticeRelates []*NoticeGroupRelate `json:"notice_relates" binding:"required"`
}

type SaveNoticeGroupResp struct {
	GroupID string `json:"group_id" binding:"required"`
}

type UpdateNoticeGroupReq struct {
	GroupID       string               `json:"group_id" binding:"required"`
	Name          string               `json:"name" binding:"required"`
	NoticeRelates []*NoticeGroupRelate `json:"notice_relates" binding:"required"`
}

type UpdateNoticeGroupResp struct {
	GroupID string `json:"group_id" binding:"required"`
}

type ListNoticeGroupReq struct {
	Keyword   string `form:"keyword"`
	ChannelID int64  `form:"channel_id"`
}

type ListNoticeGroupResp struct {
	List []*NoticeGroup `json:"list"`
}

type DetailNoticeGroupReq struct {
	GroupID string `form:"group_id" binding:"required"`
}

type DetailNoticeGroupResp struct {
	Group *NoticeGroup `json:"group"`
}

type RemoveNoticeGroupReq struct {
	GroupID string `json:"group_id"`
}
