package packer

import (
	"encoding/json"
	"permission-open/internal/pkg/biz/consts"
	"permission-open/internal/pkg/dal/model"
	"permission-open/internal/pkg/dal/rao"
)

func TransNoticeModelToRaoNotice(thirdNotice *model.ThirdNotice, thirdNoticeChannel *model.ThirdNoticeChannel) *rao.Notice {
	notice := &rao.Notice{
		NoticeID:       thirdNotice.NoticeID,
		Name:           thirdNotice.Name,
		Status:         thirdNotice.Status,
		ChannelID:      thirdNotice.ChannelID,
		CreatedTimeSec: thirdNotice.CreatedAt.Unix(),
		UpdatedTimeSec: thirdNotice.UpdatedAt.Unix(),
	}

	notice.ChannelName = thirdNoticeChannel.Name
	notice.ChannelType = thirdNoticeChannel.Type
	notice.ChannelTypeName = consts.NoticeChannelTypeCn[thirdNoticeChannel.Type]

	switch thirdNotice.ChannelID {
	case consts.NoticeChannelIDFRobot:
		s := &rao.FeiShuRobot{}
		_ = json.Unmarshal([]byte(thirdNotice.Params), s)
		notice.FeiShuRobot = s
	case consts.NoticeChannelIDFApp:
		s := &rao.FeiShuApp{}
		_ = json.Unmarshal([]byte(thirdNotice.Params), &s)
		notice.FeiShuApp = s
	case consts.NoticeChannelIDWxApp:
		s := &rao.WechatApp{}
		_ = json.Unmarshal([]byte(thirdNotice.Params), &s)
		notice.WechatApp = s
	case consts.NoticeChannelIDWxRobot:
		s := &rao.WechatRobot{}
		_ = json.Unmarshal([]byte(thirdNotice.Params), &s)
		notice.WechatRobot = s
	case consts.NoticeChannelIDEmail:
		s := &rao.SMTPEmail{}
		_ = json.Unmarshal([]byte(thirdNotice.Params), &s)
		notice.SMTPEmail = s
	case consts.NoticeChannelIDDingRobot:
		s := &rao.DingTalkRobot{}
		_ = json.Unmarshal([]byte(thirdNotice.Params), &s)
		notice.DingTalkRobot = s
	case consts.NoticeChannelIDDingApp:
		s := &rao.DingTalkApp{}
		_ = json.Unmarshal([]byte(thirdNotice.Params), &s)
		notice.DingTalkApp = s
	}

	return notice
}

func TransNoticesModelToRaoNotice(
	thirdNotice []*model.ThirdNotice,
	thirdNoticeChannel []*model.ThirdNoticeChannel) []*rao.Notice {
	ret := make([]*rao.Notice, 0, len(thirdNotice))

	chanelMemo := make(map[int64]*model.ThirdNoticeChannel)
	for _, c := range thirdNoticeChannel {
		chanelMemo[c.ID] = c
	}

	for _, n := range thirdNotice {
		notice := &rao.Notice{
			NoticeID:       n.NoticeID,
			Name:           n.Name,
			Status:         n.Status,
			ChannelID:      n.ChannelID,
			CreatedTimeSec: n.CreatedAt.Unix(),
			UpdatedTimeSec: n.UpdatedAt.Unix(),
		}

		if c, ok := chanelMemo[n.ChannelID]; ok {
			notice.ChannelName = c.Name
			notice.ChannelType = c.Type
			notice.ChannelTypeName = consts.NoticeChannelTypeCn[chanelMemo[n.ChannelID].Type]
		}

		ret = append(ret, notice)
	}

	return ret
}

func TransNoticeGroupModelToRaoNoticeGroup(
	noticeGroup *model.ThirdNoticeGroup,
	noticeGroupRelate []*model.ThirdNoticeGroupRelate,
) *rao.NoticeGroup {
	noticeGroupRes := &rao.NoticeGroup{
		GroupID:        noticeGroup.GroupID,
		Name:           noticeGroup.Name,
		CreatedTimeSec: noticeGroup.CreatedAt.Unix(),
		UpdatedTimeSec: noticeGroup.UpdatedAt.Unix(),
	}

	noticeGroupRelateRes := make([]*rao.NoticeGroupRelate, 0, len(noticeGroupRelate))
	for _, relate := range noticeGroupRelate {
		s := &rao.NoticeGroupRelateParams{}
		_ = json.Unmarshal([]byte(relate.Params), &s)
		ngr := &rao.NoticeGroupRelate{
			NoticeID: relate.NoticeID,
			Params:   s,
		}
		noticeGroupRelateRes = append(noticeGroupRelateRes, ngr)
	}

	noticeGroupRes.NoticeRelates = noticeGroupRelateRes
	return noticeGroupRes
}

func TransNoticeGroupsModelToRaoNoticeGroup(
	noticeGroup []*model.ThirdNoticeGroup,
	noticeGroupRelate []*model.ThirdNoticeGroupRelate,
	thirdNotice []*model.ThirdNotice,
	thirdNoticeChannel []*model.ThirdNoticeChannel) []*rao.NoticeGroup {

	ret := make([]*rao.NoticeGroup, 0, len(noticeGroup))

	noticeMemo := make(map[string]*model.ThirdNotice)
	for _, t := range thirdNotice {
		noticeMemo[t.NoticeID] = t
	}

	chanelMemo := make(map[int64]*model.ThirdNoticeChannel)
	for _, c := range thirdNoticeChannel {
		chanelMemo[c.ID] = c
	}

	for _, g := range noticeGroup {
		noticeGroupRes := &rao.NoticeGroup{
			GroupID:        g.GroupID,
			Name:           g.Name,
			CreatedTimeSec: g.CreatedAt.Unix(),
			UpdatedTimeSec: g.UpdatedAt.Unix(),
		}
		notices := make([]*rao.Notice, 0)
		for _, gr := range noticeGroupRelate {
			if g.GroupID == gr.GroupID {
				if n, ok := noticeMemo[gr.NoticeID]; ok {
					notice := &rao.Notice{
						NoticeID:       n.NoticeID,
						Name:           n.Name,
						Status:         n.Status,
						ChannelID:      n.ChannelID,
						CreatedTimeSec: n.CreatedAt.Unix(),
						UpdatedTimeSec: n.UpdatedAt.Unix(),
					}

					if c, ok := chanelMemo[n.ChannelID]; ok {
						notice.ChannelName = c.Name
						notice.ChannelType = c.Type
						notice.ChannelTypeName = consts.NoticeChannelTypeCn[chanelMemo[n.ChannelID].Type]
					}

					notices = append(notices, notice)
				}
			}
		}
		noticeGroupRes.Notice = notices

		ret = append(ret, noticeGroupRes)
	}

	return ret
}
