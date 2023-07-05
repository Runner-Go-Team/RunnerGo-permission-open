package event

import (
	"context"
	"github.com/gookit/event"
	"permission-open/internal/pkg/biz/consts"
	"permission-open/internal/pkg/logic/operation"
)

// BaseParams 事件所需基础参数
func BaseParams(ctx context.Context, teamID, userID, name string) map[string]interface{} {
	return map[string]interface{}{
		"ctx":     ctx,
		"team_id": teamID,
		"user_id": userID,
		"name":    name,
	}
}

// InsertCompanyOperationLog 添加企业操作日志
func InsertCompanyOperationLog(e event.Event, category int32, action int32) error {
	return operation.InsertCompanyLog(
		e.Get("ctx").(context.Context),
		e.Get("team_id").(string),
		e.Get("user_id").(string),
		e.Get("name").(string),
		category,
		action,
	)
}

// Register 注册事件
func Register() error {
	event.On(consts.ActionCompanySaveMember, event.ListenerFunc(func(e event.Event) error {
		return InsertCompanyOperationLog(e, consts.OperationCategoryCreate, consts.ActionMap[consts.ActionCompanySaveMember])
	}), event.Normal)

	event.On(consts.ActionCompanyExportMember, event.ListenerFunc(func(e event.Event) error {
		return InsertCompanyOperationLog(e, consts.OperationCategoryExport, consts.ActionMap[consts.ActionCompanyExportMember])
	}), event.Normal)

	event.On(consts.ActionCompanyUpdateMember, event.ListenerFunc(func(e event.Event) error {
		return InsertCompanyOperationLog(e, consts.OperationCategoryUpdate, consts.ActionMap[consts.ActionCompanyUpdateMember])
	}), event.Normal)

	event.On(consts.ActionCompanyRemoveMember, event.ListenerFunc(func(e event.Event) error {
		return InsertCompanyOperationLog(e, consts.OperationCategoryDelete, consts.ActionMap[consts.ActionCompanyRemoveMember])
	}), event.Normal)

	event.On(consts.ActionCompanySetRoleMember, event.ListenerFunc(func(e event.Event) error {
		return InsertCompanyOperationLog(e, consts.OperationCategoryUpdate, consts.ActionMap[consts.ActionCompanySetRoleMember])
	}), event.Normal)

	event.On(consts.ActionTeamSave, event.ListenerFunc(func(e event.Event) error {
		return InsertCompanyOperationLog(e, consts.OperationCategoryCreate, consts.ActionMap[consts.ActionTeamSave])
	}), event.Normal)

	event.On(consts.ActionTeamUpdate, event.ListenerFunc(func(e event.Event) error {
		return InsertCompanyOperationLog(e, consts.OperationCategoryUpdate, consts.ActionMap[consts.ActionTeamUpdate])
	}), event.Normal)

	event.On(consts.ActionTeamSaveMember, event.ListenerFunc(func(e event.Event) error {
		return InsertCompanyOperationLog(e, consts.OperationCategoryCreate, consts.ActionMap[consts.ActionTeamSaveMember])
	}), event.Normal)

	event.On(consts.ActionTeamRemoveMember, event.ListenerFunc(func(e event.Event) error {
		return InsertCompanyOperationLog(e, consts.OperationCategoryRemove, consts.ActionMap[consts.ActionTeamRemoveMember])
	}), event.Normal)

	event.On(consts.ActionTeamSetRoleMember, event.ListenerFunc(func(e event.Event) error {
		return InsertCompanyOperationLog(e, consts.OperationCategoryUpdate, consts.ActionMap[consts.ActionTeamSetRoleMember])
	}), event.Normal)

	event.On(consts.ActionTeamDisband, event.ListenerFunc(func(e event.Event) error {
		return InsertCompanyOperationLog(e, consts.OperationCategoryDisband, consts.ActionMap[consts.ActionTeamDisband])
	}), event.Normal)

	event.On(consts.ActionTeamTransfer, event.ListenerFunc(func(e event.Event) error {
		return InsertCompanyOperationLog(e, consts.OperationCategoryTransfer, consts.ActionMap[consts.ActionTeamTransfer])
	}), event.Normal)

	event.On(consts.ActionRoleSave, event.ListenerFunc(func(e event.Event) error {
		return InsertCompanyOperationLog(e, consts.OperationCategoryCreate, consts.ActionMap[consts.ActionRoleSave])
	}), event.Normal)

	event.On(consts.ActionRoleSet, event.ListenerFunc(func(e event.Event) error {
		return InsertCompanyOperationLog(e, consts.OperationCategoryUpdate, consts.ActionMap[consts.ActionRoleSet])
	}), event.Normal)

	event.On(consts.ActionRoleRemove, event.ListenerFunc(func(e event.Event) error {
		return InsertCompanyOperationLog(e, consts.OperationCategoryDelete, consts.ActionMap[consts.ActionRoleRemove])
	}), event.Normal)

	// 三方通知
	event.On(consts.ActionNoticeSave, event.ListenerFunc(func(e event.Event) error {
		return InsertCompanyOperationLog(e, consts.OperationCategoryCreate, consts.ActionMap[consts.ActionNoticeSave])
	}), event.Normal)

	event.On(consts.ActionNoticeUpdate, event.ListenerFunc(func(e event.Event) error {
		return InsertCompanyOperationLog(e, consts.OperationCategoryUpdate, consts.ActionMap[consts.ActionNoticeUpdate])
	}), event.Normal)

	event.On(consts.ActionNoticeSetStatus, event.ListenerFunc(func(e event.Event) error {
		return InsertCompanyOperationLog(e, consts.OperationCategoryUpdate, consts.ActionMap[consts.ActionNoticeSetStatus])
	}), event.Normal)

	event.On(consts.ActionNoticeRemove, event.ListenerFunc(func(e event.Event) error {
		return InsertCompanyOperationLog(e, consts.OperationCategoryDelete, consts.ActionMap[consts.ActionNoticeRemove])
	}), event.Normal)

	// 三方通知组
	event.On(consts.ActionNoticeGroupSave, event.ListenerFunc(func(e event.Event) error {
		return InsertCompanyOperationLog(e, consts.OperationCategoryCreate, consts.ActionMap[consts.ActionNoticeGroupSave])
	}), event.Normal)

	event.On(consts.ActionNoticeGroupUpdate, event.ListenerFunc(func(e event.Event) error {
		return InsertCompanyOperationLog(e, consts.OperationCategoryUpdate, consts.ActionMap[consts.ActionNoticeGroupUpdate])
	}), event.Normal)

	event.On(consts.ActionNoticeGroupRemove, event.ListenerFunc(func(e event.Event) error {
		return InsertCompanyOperationLog(e, consts.OperationCategoryDelete, consts.ActionMap[consts.ActionNoticeGroupRemove])
	}), event.Normal)

	return nil
}
