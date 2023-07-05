package notice

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	gevent "github.com/gookit/event"
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
)

func SaveNotice(ctx context.Context, userID string, req rao.SaveNoticeReq) (string, error) {
	channelID := req.ChannelId
	noticeID := req.NoticeID

	tn := query.Use(dal.DB()).ThirdNotice
	_, err := tn.WithContext(ctx).Where(tn.Name.Eq(req.Name), tn.NoticeID.Neq(req.NoticeID)).First()
	if err == nil {
		return "", errmsg.ErrNoticeNameRepeat
	}

	var params string
	switch channelID {
	case consts.NoticeChannelIDFRobot:
		if _, err = tn.WithContext(ctx).Where(tn.Params.Like(fmt.Sprintf("%%%s%%", req.FeiShuRobot.WebhookURL)), tn.NoticeID.Neq(req.NoticeID)).First(); err == nil {
			return "", errmsg.ErrNoticeWebhookURLRepeat
		}
		marshal, err := json.Marshal(req.FeiShuRobot)
		if err != nil {
			return "", err
		}
		params = string(marshal)
	case consts.NoticeChannelIDFApp:
		if _, err = tn.WithContext(ctx).Where(
			tn.Params.Like(fmt.Sprintf("%%%s%%", req.FeiShuApp.AppID)),
			tn.NoticeID.Neq(req.NoticeID),
			tn.ChannelID.Eq(consts.NoticeChannelIDFApp),
		).First(); err == nil {
			return "", errmsg.ErrNoticeAppIDRepeat
		}

		marshal, err := json.Marshal(req.FeiShuApp)
		if err != nil {
			return "", err
		}
		params = string(marshal)
	case consts.NoticeChannelIDWxApp:
		marshal, err := json.Marshal(req.WechatApp)
		if err != nil {
			return "", err
		}
		params = string(marshal)
	case consts.NoticeChannelIDWxRobot:
		if _, err = tn.WithContext(ctx).Where(tn.Params.Like(fmt.Sprintf("%%%s%%", req.WechatRobot.WebhookURL)), tn.NoticeID.Neq(req.NoticeID)).First(); err == nil {
			return "", errmsg.ErrNoticeWebhookURLRepeat
		}
		marshal, err := json.Marshal(req.WechatRobot)
		if err != nil {
			return "", err
		}
		params = string(marshal)
	case consts.NoticeChannelIDEmail:
		marshal, err := json.Marshal(req.SMTPEmail)
		if err != nil {
			return "", err
		}
		params = string(marshal)
	case consts.NoticeChannelIDDingRobot:
		if _, err = tn.WithContext(ctx).Where(tn.Params.Like(fmt.Sprintf("%%%s%%", req.DingTalkRobot.WebhookURL)), tn.NoticeID.Neq(req.NoticeID)).First(); err == nil {
			return "", errmsg.ErrNoticeWebhookURLRepeat
		}
		marshal, err := json.Marshal(req.DingTalkRobot)
		if err != nil {
			return "", err
		}
		params = string(marshal)
	case consts.NoticeChannelIDDingApp:
		if _, err = tn.WithContext(ctx).Where(
			tn.Params.Like(fmt.Sprintf("%%%s%%", req.DingTalkApp.AgentId)),
			tn.NoticeID.Neq(req.NoticeID),
			tn.ChannelID.Eq(consts.NoticeChannelIDDingApp),
		).First(); err == nil {
			return "", errmsg.ErrNoticeAppIDRepeat
		}

		marshal, err := json.Marshal(req.DingTalkApp)
		if err != nil {
			return "", err
		}
		params = string(marshal)
	}

	noticeData := &model.ThirdNotice{
		Name:      req.Name,
		ChannelID: req.ChannelId,
		Params:    params,
	}

	if err := dal.GetQuery().Transaction(func(tx *query.Query) error {
		// step1: 添加 notice
		if len(req.NoticeID) <= 0 {
			noticeData.NoticeID = uuid.GetUUID()
			if err := tx.ThirdNotice.WithContext(ctx).Create(noticeData); err != nil {
				return err
			}
			noticeID = noticeData.NoticeID
		} else { // 修改
			if _, err := tx.ThirdNotice.WithContext(ctx).Where(tx.ThirdNotice.NoticeID.Eq(req.NoticeID)).Updates(noticeData); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return "", err
	}

	// 触发事件
	if len(req.NoticeID) > 0 {
		if err, _ := gevent.Trigger(consts.ActionNoticeUpdate, event.BaseParams(ctx, "", userID, req.Name)); err != nil {
			return "", err
		}
	} else {
		if err, _ := gevent.Trigger(consts.ActionNoticeSave, event.BaseParams(ctx, "", userID, req.Name)); err != nil {
			return "", err
		}
	}

	return noticeID, nil
}

func List(ctx context.Context) ([]*rao.Notice, error) {
	tn := query.Use(dal.DB()).ThirdNotice
	thirdNoticeList, err := tn.WithContext(ctx).Find()
	if err != nil {
		return nil, err
	}

	tnc := query.Use(dal.DB()).ThirdNoticeChannel
	thirdNoticeChannelList, err := tnc.WithContext(ctx).Find()
	if err != nil {
		return nil, err
	}

	return packer.TransNoticesModelToRaoNotice(thirdNoticeList, thirdNoticeChannelList), nil
}

func Detail(ctx context.Context, noticeID string) (*rao.Notice, error) {
	tn := query.Use(dal.DB()).ThirdNotice
	noticeInfo, err := tn.WithContext(ctx).Where(tn.NoticeID.Eq(noticeID)).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errmsg.ErrNoticeNotFound
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) && err != nil {
		return nil, err
	}

	tnc := query.Use(dal.DB()).ThirdNoticeChannel
	thirdNoticeChannel, err := tnc.WithContext(ctx).Where(tnc.ID.Eq(noticeInfo.ChannelID)).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errmsg.ErrNoticeNotFound
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) && err != nil {
		return nil, err
	}

	return packer.TransNoticeModelToRaoNotice(noticeInfo, thirdNoticeChannel), nil
}

func SetStatus(ctx context.Context, userID string, req *rao.SetStatusNoticeReq) error {
	tn := query.Use(dal.DB()).ThirdNotice
	thirdNotice, err := tn.WithContext(ctx).Where(tn.NoticeID.Eq(req.NoticeID)).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	if err != nil {
		return err
	}

	err = dal.GetQuery().Transaction(func(tx *query.Query) error {
		if _, err := tx.ThirdNotice.WithContext(ctx).Where(
			tx.ThirdNotice.NoticeID.Eq(req.NoticeID)).UpdateColumn(tx.ThirdNotice.Status, req.Status); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	if err, _ = gevent.Trigger(consts.ActionNoticeSetStatus, event.BaseParams(ctx, "", userID, thirdNotice.Name)); err != nil {
		return err
	}

	return nil
}

func RemoveNotice(ctx context.Context, userID string, noticeID string) error {
	tn := dal.GetQuery().ThirdNotice
	thirdNotice, err := tn.WithContext(ctx).Where(tn.NoticeID.Eq(noticeID)).First()
	if err != nil {
		return err
	}

	if err := query.Use(dal.DB()).Transaction(func(tx *query.Query) error {
		if _, err := tx.ThirdNotice.WithContext(ctx).Where(tx.ThirdNotice.NoticeID.Eq(noticeID)).Delete(); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	if err, _ := gevent.Trigger(consts.ActionNoticeRemove, event.BaseParams(ctx, "", userID, thirdNotice.Name)); err != nil {
		return err
	}

	return nil
}

func GetThirdUsers(ctx context.Context, noticeID string) (*rao.ThirdCompanyUsers, error) {
	tn := dal.GetQuery().ThirdNotice
	thirdNotice, err := tn.WithContext(ctx).Where(tn.NoticeID.Eq(noticeID)).First()
	if err != nil {
		return nil, err
	}

	switch thirdNotice.ChannelID {
	case consts.NoticeChannelIDFApp:
		s := &rao.FeiShuApp{}
		_ = json.Unmarshal([]byte(thirdNotice.Params), &s)
		users, err := GetFerShuUsers(ctx, s.AppID, s.AppSecret)
		if err != nil {
			return nil, err
		}

		return users, nil
	case consts.NoticeChannelIDDingApp:
		s := &rao.DingTalkApp{}
		_ = json.Unmarshal([]byte(thirdNotice.Params), &s)
		users, err := GetDingTalkUsers(ctx, s.AppKey, s.AppSecret)
		if err != nil {
			return nil, err
		}

		return users, nil
	}

	return nil, err
}
