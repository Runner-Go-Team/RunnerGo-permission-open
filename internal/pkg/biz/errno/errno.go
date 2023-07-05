// Package errno 定义所有错误码
package errno

const (
	Ok = 0

	ErrParam           = 10001
	ErrServer          = 10002
	ErrNonce           = 10003
	ErrTimeStamp       = 10004
	ErrRPCFailed       = 10005
	ErrInvalidToken    = 10006
	ErrMarshalFailed   = 10007
	ErrUnMarshalFailed = 10008
	ErrMustDID         = 10011
	ErrMustSN          = 10012
	ErrHttpFailed      = 10013
	ErrRedisFailed     = 10100
	ErrMongoFailed     = 10101
	ErrMysqlFailed     = 10102
	ErrRecordNotFound  = 10103

	ErrSignError       = 20001
	ErrRepeatRequest   = 20002
	ErrMustLogin       = 20003
	ErrAuthFailed      = 20004
	ErrAccountDel      = 20005
	ErrAccountNotFound = 20006
	ErrPasswordFailed  = 20007

	ErrCompanyNotRemoveMember   = 20101
	ErrCompanyNotFound          = 20102
	ErrCompanySupUpdatePassword = 20103

	ErrUserDisable         = 20201
	ErrUserNotRole         = 20202
	ErrUserForbidden       = 20203
	ErrYetAccountRegister  = 20204
	ErrYetNicknameRegister = 20205
	ErrYetEmailRegister    = 20206
	ErrYetEmailValid       = 20207
	ErrYetEmailNotFound    = 20208
	ErrYetUserNotFound     = 20209

	ErrTeamNotFound        = 20301
	ErrTeamNotRemoveSelf   = 20302
	ErrTeamNotRemoveCreate = 20303
	ErrTeamSaveRepeat      = 20304
	ErrTeamTypePrivate     = 20305

	ErrRoleNotDel    = 20401
	ErrRoleNotChange = 20402
	ErrRoleExists    = 20403
	ErrRoleNotExists = 20404
	ErrRoleForbidden = 20405

	ErrFileMaxSize  = 20501
	ErrFileMaxLimit = 20502

	ErrNoticeNotFound         = 20602
	ErrNoticeNameRepeat       = 20603
	ErrNoticeGroupNameRepeat  = 20604
	ErrNoticeWebhookURLRepeat = 20605
	ErrNoticeAppIDRepeat      = 20606
	ErrNoticeRelateNotNull    = 20607
)

// CodeAlertMap 错图码映射错误提示，展示给用户
var CodeAlertMap = map[int]string{
	Ok:                          "成功",
	ErrServer:                   "服务器错误",
	ErrParam:                    "参数校验错误",
	ErrSignError:                "签名错误",
	ErrRepeatRequest:            "重放请求",
	ErrNonce:                    "_nonce参数错误",
	ErrTimeStamp:                "_timestamp参数错误",
	ErrRecordNotFound:           "数据库记录不存在",
	ErrRPCFailed:                "请求下游服务失败",
	ErrInvalidToken:             "无效的token",
	ErrMarshalFailed:            "序列化失败",
	ErrUnMarshalFailed:          "反序列化失败",
	ErrRedisFailed:              "redis操作失败",
	ErrMongoFailed:              "mongo操作失败",
	ErrMysqlFailed:              "mysql操作失败",
	ErrMustLogin:                "没有获取到登录态",
	ErrMustDID:                  "缺少设备DID信息",
	ErrMustSN:                   "缺少设备SN信息",
	ErrHttpFailed:               "请求下游Http服务失败",
	ErrAuthFailed:               "认证错误",
	ErrYetAccountRegister:       "用户账户已注册",
	ErrUserDisable:              "用户账户已封禁",
	ErrTeamNotFound:             "团队不存在",
	ErrFileMaxSize:              "文件大小超过2M",
	ErrRoleNotDel:               "该角色不可删除",
	ErrRoleNotChange:            "角色不可修改",
	ErrTeamNotRemoveSelf:        "不能移除自己",
	ErrTeamNotRemoveCreate:      "不能移除团队管理员",
	ErrUserNotRole:              "用户不存在角色",
	ErrUserForbidden:            "无权访问",
	ErrTeamSaveRepeat:           "团队重复",
	ErrRoleExists:               "角色名重复",
	ErrTeamTypePrivate:          "私有团队不能添加成员",
	ErrCompanyNotRemoveMember:   "不能移除超管",
	ErrYetNicknameRegister:      "昵称已存在",
	ErrYetEmailRegister:         "邮箱已注册",
	ErrYetEmailValid:            "邮箱格式不正确",
	ErrRoleNotExists:            "角色不存在",
	ErrRoleForbidden:            "不能添加该角色，角色较高",
	ErrAccountNotFound:          "账号不存在",
	ErrPasswordFailed:           "密码错误",
	ErrFileMaxLimit:             "一次导入的成员数据最多为1000条",
	ErrCompanyNotFound:          "企业不存在",
	ErrNoticeNameRepeat:         "通知名称不可重复",
	ErrNoticeGroupNameRepeat:    "通知组名称不可重复",
	ErrNoticeWebhookURLRepeat:   "通知三方 hook 地址不可重复",
	ErrCompanySupUpdatePassword: "超级管理员才能修改密码",
	ErrNoticeAppIDRepeat:        "app ID不可重复",
	ErrNoticeRelateNotNull:      "已配置的第三方集成不能为空",
	ErrYetEmailNotFound:         "邮箱不能为空",
	ErrYetUserNotFound:          "用户不能为空",
}

// CodeMsgMap 错误码映射错误信息，不展示给用户
var CodeMsgMap = map[int]string{
	Ok:                          "success",
	ErrServer:                   "internal server error",
	ErrParam:                    "param error",
	ErrSignError:                "signature error",
	ErrRepeatRequest:            "repeat request",
	ErrNonce:                    "nonce error",
	ErrTimeStamp:                "timestamp error",
	ErrRecordNotFound:           "record not found",
	ErrRPCFailed:                "rpc failed",
	ErrInvalidToken:             "invalid token",
	ErrMarshalFailed:            "marshal failed",
	ErrUnMarshalFailed:          "unmarshal failed",
	ErrRedisFailed:              "redis operate failed",
	ErrMongoFailed:              "mongo operate failed",
	ErrMysqlFailed:              "mysql operate failed",
	ErrMustLogin:                "must login",
	ErrMustDID:                  "must DID",
	ErrMustSN:                   "must SN",
	ErrHttpFailed:               "http failed",
	ErrAuthFailed:               "auth failed",
	ErrYetAccountRegister:       "ErrYetAccountRegister",
	ErrUserDisable:              "ErrUserDisable",
	ErrTeamNotFound:             "ErrTeamNotFound",
	ErrFileMaxSize:              "ErrFileMaxSize",
	ErrRoleNotDel:               "ErrRoleNotDel",
	ErrRoleNotChange:            "ErrRoleNotChange",
	ErrTeamNotRemoveSelf:        "ErrTeamNotRemoveSelf",
	ErrTeamNotRemoveCreate:      "ErrTeamNotRemoveCreate",
	ErrUserNotRole:              "ErrUserNotRole",
	ErrUserForbidden:            "ErrUserForbidden",
	ErrTeamSaveRepeat:           "ErrTeamSaveRepeat",
	ErrRoleExists:               "ErrRoleExists",
	ErrTeamTypePrivate:          "ErrTeamTypePrivate",
	ErrCompanyNotRemoveMember:   "ErrCompanyNotRemoveMember",
	ErrYetNicknameRegister:      "ErrYetNicknameRegister",
	ErrYetEmailRegister:         "ErrYetEmailRegister",
	ErrYetEmailValid:            "ErrYetEmailValid",
	ErrRoleNotExists:            "ErrRoleNotExists",
	ErrRoleForbidden:            "ErrRoleForbidden",
	ErrAccountNotFound:          "ErrAccountNotFound",
	ErrPasswordFailed:           "ErrPasswordFailed",
	ErrFileMaxLimit:             "ErrFileMaxLimit",
	ErrCompanyNotFound:          "ErrCompanyNotFound",
	ErrNoticeNameRepeat:         "ErrNoticeNameRepeat",
	ErrNoticeGroupNameRepeat:    "ErrNoticeGroupNameRepeat",
	ErrNoticeWebhookURLRepeat:   "ErrNoticeWebhookURLRepeat",
	ErrCompanySupUpdatePassword: "ErrCompanySupUpdatePassword",
	ErrNoticeAppIDRepeat:        "ErrNoticeAppIDRepeat",
	ErrNoticeRelateNotNull:      "ErrNoticeRelateNotNull",
	ErrYetEmailNotFound:         "ErrYetEmailNotFound",
	ErrYetUserNotFound:          "ErrYetUserNotFound",
}
