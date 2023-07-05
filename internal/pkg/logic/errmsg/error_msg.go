package errmsg

import "errors"

var (
	ErrUserDisable         = errors.New("用户已禁用")
	ErrUserNotRole         = errors.New("用户不存在角色")
	ErrUserForbidden       = errors.New("权限不足")
	ErrYetAccountRegister  = errors.New("账号已注册")
	ErrYetNicknameRegister = errors.New("昵称已注册")
	ErrYetEmailRegister    = errors.New("邮箱已注册")
	ErrYetEmailValid       = errors.New("邮箱格式不正确")
	ErrYetEmailNotFound    = errors.New("邮箱不能为空")
	ErrYetUserNotFound     = errors.New("用户不能为空")
	ErrAccountNotFound     = errors.New("账号不存在")
	ErrPasswordFailed      = errors.New("密码错误")

	ErrTeamNotFound        = errors.New("团队不存在或已解散")
	ErrTeamNotRemoveSelf   = errors.New("不能移除自己")
	ErrTeamNotRemoveCreate = errors.New("不能移除团队管理员")
	ErrTeamSaveRepeat      = errors.New("团队重复")
	ErrTeamTypePrivate     = errors.New("私有团队不能添加成员")

	ErrRoleNotDel    = errors.New("该角色不可删除")
	ErrRoleNotChange = errors.New("角色不可修改")
	ErrRoleExists    = errors.New("角色已存在")
	ErrRoleNotExists = errors.New("角色不存在")
	ErrRoleForbidden = errors.New("不能添加该角色，角色较高")

	ErrCompanyNotRemoveMember = errors.New("不能移除超管")

	ErrFileMaxLimit = errors.New("一次导入的成员数据最多为1000条")

	ErrNoticeNotFound         = errors.New("通知不存在或已删除")
	ErrNoticeNameRepeat       = errors.New("通知名称不可重复")
	ErrNoticeWebhookURLRepeat = errors.New("通知三方 hook 地址不可重复")
	ErrNoticeAppIDRepeat      = errors.New("app ID不可重复")
	ErrNoticeGroupNotFound    = errors.New("通知组不存在或已删除")
	ErrNoticeGroupNameRepeat  = errors.New("通知组名称不可重复")
)
