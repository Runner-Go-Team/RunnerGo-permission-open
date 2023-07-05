package notice

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	gevent "github.com/gookit/event"
	"gorm.io/gen"
	"gorm.io/gorm"
	"permission-open/internal/pkg/biz/consts"
	"permission-open/internal/pkg/biz/uuid"
	"permission-open/internal/pkg/dal"
	"permission-open/internal/pkg/dal/model"
	"permission-open/internal/pkg/dal/query"
	"permission-open/internal/pkg/dal/rao"
	"permission-open/internal/pkg/event"
	"permission-open/internal/pkg/logic/errmsg"
	"permission-open/internal/pkg/packer"
	"permission-open/internal/pkg/public"
	"strings"
)

func SaveNoticeGroup(ctx context.Context, userID string, req rao.SaveNoticeGroupReq) (string, error) {
	var groupID string
	noticeRelates := req.NoticeRelates

	tng := query.Use(dal.DB()).ThirdNoticeGroup
	_, err := tng.WithContext(ctx).Where(tng.Name.Eq(req.Name)).First()
	if err == nil {
		return "", errmsg.ErrNoticeGroupNameRepeat
	}

	noticeGroupData := &model.ThirdNoticeGroup{
		Name: req.Name,
	}

	tn := query.Use(dal.DB()).ThirdNotice
	// 验证数据
	noticeIDs := make([]string, 0, len(noticeRelates))
	for _, nr := range noticeRelates {
		noticeIDs = append(noticeIDs, nr.NoticeID)
	}

	noticeList, err := tn.WithContext(ctx).Where(tn.NoticeID.In(noticeIDs...)).Find()
	if err != nil {
		return "", err
	}
	noticeMemo := make(map[string]*model.ThirdNotice)
	for _, notice := range noticeList {
		noticeMemo[notice.NoticeID] = notice
	}

	for _, nr := range noticeRelates {
		if n, ok := noticeMemo[nr.NoticeID]; ok {
			// 验证邮箱
			if n.ChannelID == consts.NoticeChannelIDEmail {
				if len(nr.Params.Emails) <= 0 {
					return "", errmsg.ErrYetEmailNotFound
				}
				for _, email := range nr.Params.Emails {
					if !public.IsEmailValid(email) {
						return "", errmsg.ErrYetEmailValid
					}
				}
			}
			// 验证用户
			if n.ChannelID == consts.NoticeChannelIDFApp ||
				n.ChannelID == consts.NoticeChannelIDDingApp ||
				n.ChannelID == consts.NoticeChannelIDWxApp {
				if len(nr.Params.UserIDs) <= 0 {
					return "", errmsg.ErrYetUserNotFound
				}
			}
		}
	}

	if err := dal.GetQuery().Transaction(func(tx *query.Query) error {
		// step1: 添加 group
		// step2: 添加 group_relate
		noticeGroupData.GroupID = uuid.GetUUID()
		if err := tx.ThirdNoticeGroup.WithContext(ctx).Create(noticeGroupData); err != nil {
			return err
		}
		groupID = noticeGroupData.GroupID

		noticeGroupRelates := make([]*model.ThirdNoticeGroupRelate, 0, len(noticeRelates))
		for _, nr := range noticeRelates {
			marshal, err := json.Marshal(nr.Params)
			if err != nil {
				return err
			}
			params := string(marshal)
			noticeGroupRelates = append(noticeGroupRelates, &model.ThirdNoticeGroupRelate{
				GroupID:  groupID,
				NoticeID: nr.NoticeID,
				Params:   params,
			})
		}
		if err := tx.ThirdNoticeGroupRelate.WithContext(ctx).Create(noticeGroupRelates...); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return "", err
	}

	// 触发事件
	if err, _ := gevent.Trigger(consts.ActionNoticeGroupSave, event.BaseParams(ctx, "", userID, req.Name)); err != nil {
		return "", err
	}

	return groupID, nil
}

func UpdateNoticeGroup(ctx context.Context, userID string, req rao.UpdateNoticeGroupReq) (string, error) {
	// step1: 修改 group
	// step2： 删除之前的通知，添加新的通知
	groupID := req.GroupID
	noticeRelates := req.NoticeRelates

	tng := query.Use(dal.DB()).ThirdNoticeGroup
	_, err := tng.WithContext(ctx).Where(tng.Name.Eq(req.Name), tng.GroupID.Neq(req.GroupID)).First()
	if err == nil {
		return "", errmsg.ErrNoticeGroupNameRepeat
	}

	noticeGroupData := &model.ThirdNoticeGroup{
		Name: req.Name,
	}

	tn := query.Use(dal.DB()).ThirdNotice
	// 验证数据
	noticeIDs := make([]string, 0, len(noticeRelates))
	for _, nr := range noticeRelates {
		noticeIDs = append(noticeIDs, nr.NoticeID)
	}

	noticeList, err := tn.WithContext(ctx).Where(tn.NoticeID.In(noticeIDs...)).Find()
	if err != nil {
		return "", err
	}
	noticeMemo := make(map[string]*model.ThirdNotice)
	for _, notice := range noticeList {
		noticeMemo[notice.NoticeID] = notice
	}

	for _, nr := range noticeRelates {
		if n, ok := noticeMemo[nr.NoticeID]; ok {
			// 验证邮箱
			if n.ChannelID == consts.NoticeChannelIDEmail {
				if len(nr.Params.Emails) <= 0 {
					return "", errmsg.ErrYetEmailNotFound
				}
				for _, email := range nr.Params.Emails {
					if !public.IsEmailValid(email) {
						return "", errmsg.ErrYetEmailValid
					}
				}
			}
			// 验证用户
			if n.ChannelID == consts.NoticeChannelIDFApp ||
				n.ChannelID == consts.NoticeChannelIDDingApp ||
				n.ChannelID == consts.NoticeChannelIDWxApp {
				if len(nr.Params.UserIDs) <= 0 {
					return "", errmsg.ErrYetUserNotFound
				}
			}
		}
	}

	oldIDs := make([]string, 0)
	ngr := dal.GetQuery().ThirdNoticeGroupRelate
	groupRelateOrm, err := ngr.WithContext(ctx).Where(ngr.GroupID.Eq(groupID)).Find()
	if err != nil {
		return "", err
	}
	for _, nro := range groupRelateOrm {
		oldIDs = append(oldIDs, nro.NoticeID)
	}

	if err := dal.GetQuery().Transaction(func(tx *query.Query) error {
		if len(req.GroupID) > 0 {
			// step1: 修改 group
			if _, err := tx.ThirdNoticeGroup.WithContext(ctx).Where(tx.ThirdNoticeGroup.GroupID.Eq(groupID)).Updates(noticeGroupData); err != nil {
				return err
			}

			// step2: 删除之前的 通知
			if _, err := tx.ThirdNoticeGroupRelate.WithContext(ctx).Where(
				tx.ThirdNoticeGroupRelate.GroupID.Eq(groupID),
				tx.ThirdNoticeGroupRelate.NoticeID.In(oldIDs...)).Delete(); err != nil {
				return err
			}

			// step3: 新增
			noticeGroupRelates := make([]*model.ThirdNoticeGroupRelate, 0, len(noticeRelates))
			for _, nr := range noticeRelates {
				marshal, err := json.Marshal(nr.Params)
				if err != nil {
					return err
				}
				params := string(marshal)
				noticeGroupRelates = append(noticeGroupRelates, &model.ThirdNoticeGroupRelate{
					GroupID:  groupID,
					NoticeID: nr.NoticeID,
					Params:   params,
				})
			}
			if err := tx.ThirdNoticeGroupRelate.WithContext(ctx).Create(noticeGroupRelates...); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return "", err
	}

	// 触发事件
	if err, _ := gevent.Trigger(consts.ActionNoticeGroupUpdate, event.BaseParams(ctx, "", userID, req.Name)); err != nil {
		return "", err
	}

	return groupID, nil
}

func ListGroup(ctx context.Context, req *rao.ListNoticeGroupReq) ([]*rao.NoticeGroup, error) {
	ng := dal.GetQuery().ThirdNoticeGroup
	keyword := strings.TrimSpace(req.Keyword)
	channelID := req.ChannelID
	conditions := make([]gen.Condition, 0)
	if len(keyword) > 0 {
		conditions = append(conditions, ng.Name.Like(fmt.Sprintf("%%%s%%", keyword)))
	}

	tn := dal.GetQuery().ThirdNotice
	ngr := dal.GetQuery().ThirdNoticeGroupRelate
	if channelID > 0 {
		thirdNotice, err := tn.WithContext(ctx).Where(tn.ChannelID.Eq(channelID)).Find()
		if err != nil {
			return nil, err
		}
		noticeID := make([]string, 0, len(thirdNotice))
		for _, n := range thirdNotice {
			noticeID = append(noticeID, n.NoticeID)
		}
		noticeGroupRelate, err := ngr.WithContext(ctx).Where(ngr.NoticeID.In(noticeID...)).Find()
		if err != nil {
			return nil, err
		}

		groupID := make([]string, 0, len(noticeGroupRelate))
		for _, relate := range noticeGroupRelate {
			groupID = append(groupID, relate.GroupID)
		}
		conditions = append(conditions, ng.GroupID.In(groupID...))
	}

	noticeGroup, err := ng.WithContext(ctx).Where(conditions...).Order(ng.CreatedAt.Desc()).Find()
	if err != nil {
		return nil, err
	}

	groupIDs := make([]string, 0, len(noticeGroup))
	for _, group := range noticeGroup {
		groupIDs = append(groupIDs, group.GroupID)
	}

	noticeGroupRelate, err := ngr.WithContext(ctx).Where(ngr.GroupID.In(groupIDs...)).Find()
	if err != nil {
		return nil, err
	}

	noticeIDs := make([]string, 0, len(noticeGroupRelate))
	for _, r := range noticeGroupRelate {
		noticeIDs = append(noticeIDs, r.NoticeID)
	}

	thirdNotice, err := tn.WithContext(ctx).Where(tn.NoticeID.In(noticeIDs...)).Find()
	if err != nil {
		return nil, err
	}

	tnc := query.Use(dal.DB()).ThirdNoticeChannel
	thirdNoticeChannelList, err := tnc.WithContext(ctx).Find()
	if err != nil {
		return nil, err
	}

	return packer.TransNoticeGroupsModelToRaoNoticeGroup(noticeGroup, noticeGroupRelate, thirdNotice, thirdNoticeChannelList), nil
}

func DetailGroup(ctx context.Context, groupID string) (*rao.NoticeGroup, error) {
	ng := dal.GetQuery().ThirdNoticeGroup
	noticeGroup, err := ng.WithContext(ctx).Where(ng.GroupID.Eq(groupID)).First()
	if err != nil {
		return nil, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errmsg.ErrNoticeGroupNotFound
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) && err != nil {
		return nil, err
	}

	ngr := dal.GetQuery().ThirdNoticeGroupRelate
	noticeGroupRelate, err := ngr.WithContext(ctx).Where(ngr.GroupID.Eq(groupID)).Find()
	if err != nil {
		return nil, err
	}

	return packer.TransNoticeGroupModelToRaoNoticeGroup(noticeGroup, noticeGroupRelate), nil
}

func RemoveNoticeGroup(ctx context.Context, userID string, groupID string) error {
	tng := dal.GetQuery().ThirdNoticeGroup
	thirdNoticeGroup, err := tng.WithContext(ctx).Where(tng.GroupID.Eq(groupID)).First()
	if err != nil {
		return err
	}

	if err := query.Use(dal.DB()).Transaction(func(tx *query.Query) error {
		if _, err := tx.ThirdNoticeGroup.WithContext(ctx).Where(tx.ThirdNoticeGroup.GroupID.Eq(groupID)).Delete(); err != nil {
			return err
		}

		if _, err := tx.ThirdNoticeGroupRelate.WithContext(ctx).Where(tx.ThirdNoticeGroupRelate.GroupID.Eq(groupID)).Delete(); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	if err, _ := gevent.Trigger(consts.ActionNoticeGroupRemove, event.BaseParams(ctx, "", userID, thirdNoticeGroup.Name)); err != nil {
		return err
	}

	return nil
}
