package rao

type FeiShuRobot struct {
	WebhookURL string `json:"webhook_url"`
	Secret     string `json:"secret"`
}

type FeiShuApp struct {
	AppID     string `json:"app_id"`
	AppSecret string `json:"app_secret"`
}

type WechatRobot struct {
	WebhookURL string `json:"webhook_url"`
}

type WechatApp struct {
	CorpID     string `json:"corp_id"`
	CorpSecret string `json:"corp_secret"`
	AgentID    string `json:"agent_id"`
}

type SMTPEmail struct {
	Host     string      `json:"host"`
	Port     interface{} `json:"port"`
	Email    string      `json:"email"`
	Password string      `json:"password"`
}

type DingTalkRobot struct {
	WebhookURL string `json:"webhook_url"`
	Secret     string `json:"secret"`
}

type DingTalkApp struct {
	AgentId   string `json:"agent_id"`
	AppKey    string `json:"app_key"`
	AppSecret string `json:"app_secret"`
}

type Notice struct {
	NoticeID        string         `json:"notice_id"`
	Name            string         `json:"name"`
	Status          int32          `json:"status"`
	ChannelID       int64          `json:"channel_id"`
	ChannelName     string         `json:"channel_name"`
	ChannelType     int32          `json:"channel_type"`
	ChannelTypeName string         `json:"channel_type_name"`
	FeiShuRobot     *FeiShuRobot   `json:"fei_shu_robot,omitempty"`
	FeiShuApp       *FeiShuApp     `json:"fei_shu_app,omitempty"`
	WechatApp       *WechatApp     `json:"wechat_app,omitempty"`
	WechatRobot     *WechatRobot   `json:"wechat_robot,omitempty"`
	SMTPEmail       *SMTPEmail     `json:"smtp_email,omitempty"`
	DingTalkRobot   *DingTalkRobot `json:"ding_talk_robot,omitempty"`
	DingTalkApp     *DingTalkApp   `json:"ding_talk_app,omitempty"`
	CreatedTimeSec  int64          `json:"created_time_sec"`
	UpdatedTimeSec  int64          `json:"updated_time_sec"`
}

type SaveNoticeReq struct {
	NoticeID      string         `json:"notice_id"`
	Name          string         `json:"name" binding:"required"`
	ChannelId     int64          `json:"channel_id" binding:"required"`
	FeiShuRobot   *FeiShuRobot   `json:"fei_shu_robot"`
	FeiShuApp     *FeiShuApp     `json:"fei_shu_app"`
	WechatApp     *WechatApp     `json:"wechat_app"`
	WechatRobot   *WechatRobot   `json:"wechat_robot"`
	SMTPEmail     *SMTPEmail     `json:"smtp_email"`
	DingTalkRobot *DingTalkRobot `json:"ding_talk_robot"`
	DingTalkApp   *DingTalkApp   `json:"ding_talk_app"`
}

type SaveNoticeResp struct {
	NoticeId string `json:"notice_id" binding:"required"`
}

type ListNoticeResp struct {
	List []*Notice `json:"list"`
}

type DetailNoticeReq struct {
	NoticeID string `form:"notice_id" binding:"required"`
}

type DetailNoticeResp struct {
	Notice *Notice `json:"notice"`
}

type SetStatusNoticeReq struct {
	NoticeID string `json:"notice_id" binding:"required"`
	Status   int32  `json:"status" binding:"required"`
}

type RemoveNoticeReq struct {
	NoticeID string `json:"notice_id"`
}

type GetNoticeUsersReq struct {
	NoticeID string `form:"notice_id" binding:"required"`
}

type ThirdUserInfo struct {
	OpenID string `json:"open_id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

type ThirdDepartmentInfo struct {
	DepartmentID   string                `json:"department_id"`
	Name           string                `json:"name"`
	MemberCount    int                   `json:"member_count"`
	UserList       []ThirdUserInfo       `json:"user_list"`
	DepartmentList []ThirdDepartmentInfo `json:"department_list"`
}

type ThirdCompanyUsers struct {
	DepartmentList []ThirdDepartmentInfo `json:"department_list"`
	UserList       []ThirdUserInfo       `json:"user_list"`
}
