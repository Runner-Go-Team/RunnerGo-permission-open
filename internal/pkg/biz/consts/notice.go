package consts

const (
	NoticeChannelTypeFS    = 1 // 类型 1:飞书  2:企业微信  3:邮箱  4:钉钉
	NoticeChannelTypeWx    = 2
	NoticeChannelTypeEmail = 3
	NoticeChannelTypeDing  = 4

	NoticeChannelIDFRobot    = 1 // 1:飞书群机器人
	NoticeChannelIDFApp      = 2 // 2:飞书企业应用
	NoticeChannelIDWxApp     = 3 // 3:企业微信应用
	NoticeChannelIDWxRobot   = 4 // 4:企业微信机器人
	NoticeChannelIDEmail     = 5 // 5:邮箱
	NoticeChannelIDDingRobot = 6 // 6:钉钉群机器人
	NoticeChannelIDDingApp   = 7 // 7:钉钉应用

	// RedisNoticeFerShuUsersPrefix 获取飞书用户
	RedisNoticeFerShuUsersPrefix   = "NoticeFerShuUsers:"
	RedisNoticeDingTalkUsersPrefix = "NoticeDingTalkUsers:"
)

var NoticeChannelTypeCn = map[int32]string{
	NoticeChannelTypeFS:    "飞书",
	NoticeChannelTypeWx:    "企业微信",
	NoticeChannelTypeEmail: "邮箱",
	NoticeChannelTypeDing:  "钉钉",
}
