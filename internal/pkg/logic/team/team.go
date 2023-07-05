package team

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gen/field"
	"permission-open/internal/pkg/event"
	"strconv"
	"strings"
	"time"

	"gorm.io/gen"
	"gorm.io/gorm"

	"permission-open/internal/pkg/biz/consts"
	"permission-open/internal/pkg/biz/encrypt"
	"permission-open/internal/pkg/biz/jwt"
	"permission-open/internal/pkg/biz/uuid"
	"permission-open/internal/pkg/conf"
	"permission-open/internal/pkg/dal"
	"permission-open/internal/pkg/dal/global"
	"permission-open/internal/pkg/dal/model"
	"permission-open/internal/pkg/dal/query"
	"permission-open/internal/pkg/dal/rao"
	"permission-open/internal/pkg/logic/errmsg"
	"permission-open/internal/pkg/packer"

	"github.com/gin-gonic/gin"
	"github.com/go-omnibus/proof"
	gevent "github.com/gookit/event"
)

// SettingCreateOrUpdate  创建或修改设置  isUpdate(true:存在修改，false:存在不修改)
func SettingCreateOrUpdate(ctx context.Context, userID string, teamID string, isUpdate bool) error {
	s := query.Use(dal.DB()).Setting
	_, err := s.WithContext(ctx).Where(s.UserID.Eq(userID)).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		if err = s.WithContext(ctx).Create(&model.Setting{UserID: userID, TeamID: teamID}); err != nil {
			return err
		}
		return nil
	}
	if err != nil {
		return err
	}

	if isUpdate {
		_, err = s.WithContext(ctx).Where(s.UserID.Eq(userID)).UpdateColumnSimple(s.TeamID.Value(teamID))
		if err != nil {
			return err
		}
	}

	return nil
}

func SaveTeam(ctx *gin.Context, req *rao.SaveTeamReq) (string, error) {
	userID := jwt.GetUserIDByCtx(ctx)
	var teamID string
	err := dal.GetQuery().Transaction(func(tx *query.Query) error {
		var (
			// 是否添加超级管理员
			isSaveSuper = false
			superUserID = global.SuperManageUserID
		)
		if userID != superUserID {
			isSaveSuper = true
		}
		// step1: team 新建团队
		insertData := model.Team{
			TeamID:        uuid.GetUUID(),
			Name:          req.Name,
			Type:          req.TeamType,
			CompanyID:     req.CompanyId,
			CreatedUserID: userID,
		}
		if err := tx.Team.WithContext(ctx).Create(&insertData); err != nil {
			return err
		}
		teamID = insertData.TeamID

		// step2: 新建用户团队关系
		userTeams := make([]*model.UserTeam, 0, 2)
		userTeams = append(userTeams, &model.UserTeam{
			TeamID: teamID,
			UserID: userID,
		})

		// 超级管理员添加进去
		if isSaveSuper {
			userTeams = append(userTeams, &model.UserTeam{
				TeamID: teamID,
				UserID: superUserID,
				IsShow: consts.TeamIsShowFalse,
			})
		}
		if err := tx.UserTeam.WithContext(ctx).Create(userTeams...); err != nil {
			return err
		}

		// step3: 新建用户团队角色关系
		r := dal.GetQuery().Role
		role, err := r.WithContext(ctx).Where(
			r.CompanyID.Eq(req.CompanyId),
			r.RoleType.Eq(consts.RoleTypeTeam),
			r.Level.Eq(consts.RoleLevelSuperManager)).First()
		if err != nil {
			return err
		}
		ur := dal.GetQuery().UserRole
		userRoles := make([]*model.UserRole, 0, 2)
		userRoles = append(userRoles, &model.UserRole{
			RoleID:     role.RoleID,
			UserID:     userID,
			TeamID:     teamID,
			InviteTime: time.Now(),
		})
		if isSaveSuper {
			userRoles = append(userRoles, &model.UserRole{
				RoleID:     role.RoleID,
				UserID:     superUserID,
				TeamID:     teamID,
				InviteTime: time.Now(),
			})
		}

		if err = ur.WithContext(ctx).Create(userRoles...); err != nil {
			return err
		}

		// step4: 把用户的默认团队设置为当前新建的团队
		if err = SettingCreateOrUpdate(ctx, userID, teamID, true); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return "", err
	}

	// 触发事件
	if err, _ = gevent.Trigger(consts.ActionTeamSave, event.BaseParams(ctx, teamID, userID, req.Name)); err != nil {
		return "", err
	}

	return teamID, nil
}

// UpdateTeam 修改团队
func UpdateTeam(ctx *gin.Context, req rao.UpdateTeamReq) error {
	userID := jwt.GetUserIDByCtx(ctx)
	if req.TeamType == 0 && req.Name == "" && req.Description == nil {
		return nil
	}

	operationName := ""

	t := query.Use(dal.DB()).Team
	fields := make([]field.AssignExpr, 0)
	if req.TeamType > 0 {
		fields = append(fields, t.Type.Value(req.TeamType))
		operationName += "类型：公开 "
	}
	if req.Name != "" {
		fields = append(fields, t.Name.Value(req.Name))
		operationName += "名称：" + req.Name + " "
	}
	if req.Description != nil {
		fields = append(fields, t.Description.Value(*req.Description))
		operationName += "描述：" + *req.Description
	}

	if _, err := t.WithContext(ctx).Where(t.TeamID.Eq(req.TeamID)).UpdateColumnSimple(fields...); err != nil {
		return err
	}

	// 触发事件
	if err, _ := gevent.Trigger(consts.ActionTeamUpdate, event.BaseParams(ctx, req.TeamID, userID, operationName)); err != nil {
		return err
	}

	return nil
}

func ListTeam(ctx context.Context, userID string, req rao.ListTeamReq) ([]*rao.RelatedTeam, error) {
	// order  1:我加入的团队  2:我是团队管理员的团队  3:我收藏的团队
	var order int32 = consts.TeamListOrderUserID
	if req.Order > 0 {
		order = req.Order
	}

	var (
		teamIDs []string
		t       = query.Use(dal.DB()).Team
		ut      = query.Use(dal.DB()).UserTeam
		utc     = query.Use(dal.DB()).UserTeamCollection
	)

	// 收藏
	userTeamCollections, err := utc.WithContext(ctx).Where(utc.UserID.Eq(userID)).Find()
	if err != nil {
		return nil, err
	}

	conditions := make([]gen.Condition, 0)
	keyword := strings.TrimSpace(req.Keyword)
	if len(keyword) > 0 {
		conditions = append(conditions, t.Name.Like(fmt.Sprintf("%%%s%%", keyword)))
	}
	if order == consts.TeamListOrderSuper {
		// 查找当前企业的团队管理员的 role_id
		r := query.Use(dal.DB()).Role
		role, err := r.WithContext(ctx).Where(
			r.CompanyID.Eq(req.CompanyId),
			r.RoleType.Eq(consts.RoleTypeTeam),
			r.Level.Eq(consts.RoleLevelSuperManager)).First()
		if err != nil {
			return nil, err
		}

		ur := query.Use(dal.DB()).UserRole
		userRoles, err := ur.WithContext(ctx).Where(ur.RoleID.Eq(role.RoleID), ur.UserID.Eq(userID)).Find()
		if err != nil {
			return nil, err
		}

		for _, userRole := range userRoles {
			teamIDs = append(teamIDs, userRole.TeamID)
		}
		conditions = append(conditions, t.TeamID.In(teamIDs...))
	} else if order == consts.TeamListOrderCollection {
		for _, team := range userTeamCollections {
			teamIDs = append(teamIDs, team.TeamID)
		}
		conditions = append(conditions, t.TeamID.In(teamIDs...))
	} else {
		userTeams, err := ut.WithContext(ctx).Where(ut.UserID.Eq(userID)).Find()
		if err != nil {
			return nil, err
		}

		for _, team := range userTeams {
			teamIDs = append(teamIDs, team.TeamID)
		}
		conditions = append(conditions, t.TeamID.In(teamIDs...))
	}
	teams, err := t.WithContext(ctx).Where(conditions...).Order(t.UpdatedAt.Desc()).Find()
	if err != nil {
		return nil, err
	}

	newTeamIDs := make([]string, 0, len(teams))
	for _, team := range teams {
		newTeamIDs = append(newTeamIDs, team.TeamID)
	}

	// 查询可用的用户团队数据
	userTeamsNew, err := ut.WithContext(ctx).Where(ut.TeamID.In(newTeamIDs...), ut.IsShow.Eq(consts.TeamIsShow)).Order(ut.CreatedAt).Find()
	if err != nil {
		return nil, err
	}

	var userIDs []string
	for _, team := range userTeamsNew {
		userIDs = append(userIDs, team.UserID)
	}
	u := dal.GetQuery().User
	users, err := u.WithContext(ctx).Where(u.UserID.In(userIDs...)).Find()
	if err != nil {
		return nil, err
	}

	return packer.TransUserTeamsModelToRaoTeam(teams, userTeamsNew, users, userTeamCollections), nil
}

func RemoveMember(ctx *gin.Context, teamID string, userID string, targetUserID string) error {
	u := query.Use(dal.DB()).User
	targetUser, err := u.WithContext(ctx).Where(u.UserID.Eq(targetUserID)).First()
	if err != nil {
		return nil
	}

	ut := query.Use(dal.DB()).UserTeam
	_, err = ut.WithContext(ctx).Where(ut.TeamID.Eq(teamID), ut.UserID.Eq(targetUserID)).First()
	if err != nil {
		return nil
	}

	t := query.Use(dal.DB()).Team
	team, err := t.WithContext(ctx).Where(t.TeamID.Eq(teamID)).First()
	if err != nil {
		return nil
	}

	// 不能移除自己
	if userID == targetUserID {
		return errmsg.ErrTeamNotRemoveSelf
	}

	// 不能移除团队管理员
	r := query.Use(dal.DB()).Role
	role, err := r.WithContext(ctx).Where(r.CompanyID.Eq(team.CompanyID), r.RoleType.Eq(consts.RoleTypeTeam), r.Level.Eq(consts.RoleLevelSuperManager)).First()
	if err != nil {
		return err
	}

	ur := query.Use(dal.DB()).UserRole
	userRole, err := ur.WithContext(ctx).Where(ur.UserID.Eq(targetUserID), ur.TeamID.Eq(teamID)).First()
	if err != nil {
		return err
	}

	if userRole.RoleID == role.RoleID {
		return errmsg.ErrTeamNotRemoveCreate
	}

	err = dal.GetQuery().Transaction(func(tx *query.Query) error {
		_, err = tx.UserTeam.WithContext(ctx).Where(tx.UserTeam.TeamID.Eq(teamID), tx.UserTeam.UserID.Eq(targetUserID)).Delete()

		setting, err := tx.Setting.WithContext(ctx).Where(
			tx.Setting.TeamID.Eq(teamID),
			tx.Setting.UserID.Eq(targetUserID)).First()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		if err != nil {
			return err
		}

		//  默认团队不是当前团队不处理
		if setting.TeamID == teamID {
			teamInfo, err := tx.UserTeam.WithContext(ctx).Where(tx.UserTeam.UserID.Eq(targetUserID)).First()
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if _, err = tx.Setting.WithContext(ctx).Where(tx.Setting.ID.Eq(setting.ID)).Delete(); err != nil {
					return err
				}
				return nil
			}
			if err != nil {
				return err
			}

			// setting 有其他团队更新、没有删除
			if _, err = tx.Setting.WithContext(ctx).Where(tx.Setting.ID.Eq(setting.ID)).
				UpdateColumn(tx.Setting.TeamID, teamInfo.TeamID); err != nil {
				return err
			}
		}

		// 删除团队角色
		_, err = tx.UserRole.WithContext(ctx).Where(tx.UserRole.UserID.Eq(targetUserID), tx.UserRole.TeamID.Eq(teamID)).Delete()
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	// 触发事件
	if err, _ = gevent.Trigger(consts.ActionTeamRemoveMember, event.BaseParams(ctx, teamID, userID, targetUser.Nickname)); err != nil {
		return err
	}

	return nil
}

// DisbandTeam 解散团队
func DisbandTeam(ctx context.Context, teamID string, userID string) error {
	// 团队不存在无需解散
	t := query.Use(dal.DB()).Team
	team, err := t.WithContext(ctx).Where(t.TeamID.Eq(teamID)).First()
	if err != nil {
		return nil
	}

	err = dal.GetQuery().Transaction(func(tx *query.Query) error {
		// 删除当前团队
		_, err := tx.Team.WithContext(ctx).Where(tx.Team.TeamID.Eq(teamID)).Delete()
		if err != nil {
			return err
		}

		// 删除所有用户与解散团队之间的关系数据
		_, err = tx.UserTeam.WithContext(ctx).Where(tx.UserTeam.TeamID.Eq(teamID)).Delete()
		if err != nil {
			return err
		}

		// 默认团队修改或删除
		settings, err := tx.Setting.WithContext(ctx).Where(tx.Setting.TeamID.Eq(teamID)).Find()
		if err != nil {
			return err
		}

		for _, s := range settings {
			// setting 有其他团队更新、没有删除
			teamInfo, err := tx.UserTeam.WithContext(ctx).Where(tx.UserTeam.UserID.Eq(s.UserID)).First()
			if errors.Is(err, gorm.ErrRecordNotFound) {
				_, err = tx.Setting.WithContext(ctx).Where(tx.Setting.ID.Eq(s.ID)).Delete()
				continue
			}
			if err != nil {
				return err
			}

			_, err = tx.Setting.WithContext(ctx).Where(tx.Setting.ID.Eq(s.ID)).UpdateColumn(tx.Setting.TeamID, teamInfo.TeamID)
			if err != nil {
				return err
			}
		}

		// 删除团队角色
		_, err = tx.UserRole.WithContext(ctx).Where(tx.UserRole.TeamID.Eq(teamID)).Delete()
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	// 触发事件
	if err, _ = gevent.Trigger(consts.ActionTeamDisband, event.BaseParams(ctx, teamID, userID, team.Name)); err != nil {
		return err
	}

	return err
}

func InviteLogin(ctx *gin.Context, verifyCode string, userID string) error {
	userInfoString := encrypt.AesDecrypt(verifyCode, conf.Conf.InviteData.AesSecretKey)
	userInfoArr := strings.Split(userInfoString, "_")
	if len(userInfoArr) != 4 {
		return fmt.Errorf("验证码解析错误")
	}

	teamID := userInfoArr[0]
	roleID, _ := strconv.ParseInt(userInfoArr[1], 10, 64)
	inviteUserID := userInfoArr[2]

	// 把当前用户的当前团队设置为被邀请团队
	err := query.Use(dal.DB()).Transaction(func(tx *query.Query) error {
		// 1、把用户当前所属团队修改为被邀请的团队
		updateData := make(map[string]interface{}, 1)
		updateData["team_id"] = teamID
		_, err := tx.Setting.WithContext(ctx).Where(tx.Setting.UserID.Eq(userID)).Updates(updateData)
		if err != nil {
			proof.Infof("邀请登录--修改用户当前团队失败，err:", err)
			return err
		}
		// 2、把当前用户放到被邀请的团队里面
		_, err = tx.UserTeam.WithContext(ctx).Where(tx.UserTeam.UserID.Eq(userID)).Where(tx.UserTeam.TeamID.Eq(teamID)).First()
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		if err == gorm.ErrRecordNotFound { // 没查到，就插入
			insertData := &model.UserTeam{
				UserID:       userID,
				TeamID:       teamID,
				RoleID:       roleID,
				InviteUserID: inviteUserID,
			}
			err = tx.UserTeam.WithContext(ctx).Create(insertData)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// CollectionTeam 收藏|取消收藏
func CollectionTeam(ctx *gin.Context, userID string, teamID string, status int32) error {
	// 查询团队信息
	tx := query.Use(dal.DB())
	_, err := tx.Team.WithContext(ctx).Where(tx.Team.TeamID.Eq(teamID)).First()
	if err != nil {
		return errmsg.ErrTeamNotFound
	}
	utc := query.Use(dal.DB()).UserTeamCollection
	if status == consts.TeamCollection { // 收藏
		_, err := utc.WithContext(ctx).Where(utc.TeamID.Eq(teamID), utc.UserID.Eq(userID)).First()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err = utc.WithContext(ctx).Create(&model.UserTeamCollection{
				TeamID: teamID,
				UserID: userID,
			}); err != nil {
				return err
			}
		}
	} else if status == consts.TeamUnCollection { // 取消收藏
		if _, err = utc.WithContext(ctx).Where(utc.TeamID.Eq(teamID), utc.UserID.Eq(userID)).Delete(); err != nil {
			return err
		}
	}

	return nil
}

// Info 团队信息
func Info(ctx *gin.Context, teamID string, keyword string, page int, size int) (*rao.Team, int64, error) {
	userID := jwt.GetUserIDByCtx(ctx)

	limit := size
	offset := (page - 1) * size

	// 查询团队信息
	t := query.Use(dal.DB()).Team
	team, err := t.WithContext(ctx).Where(t.TeamID.Eq(teamID)).First()
	if err != nil {
		return nil, 0, errmsg.ErrTeamNotFound
	}

	// 查询可用的用户团队数据
	ut := query.Use(dal.DB()).UserTeam
	userTeams, err := ut.WithContext(ctx).Where(ut.TeamID.Eq(teamID), ut.IsShow.Eq(consts.TeamIsShow)).Find()
	if err != nil {
		return nil, 0, err
	}

	var kUserIDs []string
	for _, teamInfo := range userTeams {
		kUserIDs = append(kUserIDs, teamInfo.UserID)
	}

	// keyword 搜索昵称/账号
	u := dal.GetQuery().User
	conditions := make([]gen.Condition, 0)
	conditions = append(conditions, u.UserID.In(kUserIDs...))
	conditionsAccount := conditions
	keyword = strings.TrimSpace(keyword)
	if len(keyword) > 0 {
		conditions = append(conditions, u.Nickname.Like(fmt.Sprintf("%%%s%%", keyword)))
		conditionsAccount = append(conditionsAccount, u.Account.Like(fmt.Sprintf("%%%s%%", keyword)))
	}
	users, err := u.WithContext(ctx).Where(conditions...).Or(conditionsAccount...).Find()
	if err != nil {
		return nil, 0, err
	}

	userIDs := make([]string, 0, len(users))
	for _, user := range users {
		userIDs = append(userIDs, user.UserID)
	}

	userTeams, total, err := ut.WithContext(ctx).Where(
		ut.TeamID.Eq(teamID),
		ut.IsShow.Eq(consts.TeamIsShow),
		ut.UserID.In(userIDs...),
	).Order(ut.ID).FindByPage(offset, limit)
	if err != nil {
		return nil, 0, err
	}

	var teamNotShowUserIDs []string
	if err = ut.WithContext(ctx).Where(ut.TeamID.Eq(teamID), ut.IsShow.Eq(consts.TeamIsShowFalse)).Pluck(ut.UserID, &teamNotShowUserIDs); err != nil {
		return nil, 0, err
	}

	var userIDsFilter []string
	for _, teamInfo := range userTeams {
		userIDsFilter = append(userIDsFilter, teamInfo.UserID)
		userIDsFilter = append(userIDsFilter, userID)
		userIDsFilter = append(userIDsFilter, teamInfo.InviteUserID)
		if len(teamNotShowUserIDs) > 0 {
			userIDsFilter = append(userIDsFilter, teamNotShowUserIDs...)
		}
	}

	users, err = u.WithContext(ctx).Where(u.UserID.In(userIDsFilter...)).Find()
	if err != nil {
		return nil, 0, err
	}

	r := dal.GetQuery().Role
	roles, err := r.WithContext(ctx).Where(r.CompanyID.Eq(team.CompanyID)).Find()
	if err != nil {
		return nil, 0, err
	}

	ur := dal.GetQuery().UserRole
	userRoles, err := ur.WithContext(ctx).Where(ur.UserID.In(userIDsFilter...)).Find()
	if err != nil {
		return nil, 0, err
	}

	uc := dal.GetQuery().UserCompany
	userCompanies, err := uc.WithContext(ctx).Where(uc.CompanyID.Eq(team.CompanyID), uc.UserID.In(userIDsFilter...)).Find()
	if err != nil {
		return nil, 0, err
	}

	return packer.TransTeamModelToRaoTeam(
		userID,
		team,
		userTeams,
		users,
		roles,
		userRoles,
		userCompanies), total, nil
}

func SaveMembers(ctx *gin.Context, userID string, teamID string, members []rao.SaveMembers) error {
	// 查询团队信息
	t := query.Use(dal.DB()).Team
	teamInfo, err := t.WithContext(ctx).Where(t.TeamID.Eq(teamID)).First()
	if err != nil {
		return errmsg.ErrTeamNotFound
	}
	if teamInfo.Type == consts.TeamTypePrivate {
		return errmsg.ErrTeamTypePrivate
	}

	// 查询当前的用户团队数据
	ut := query.Use(dal.DB()).UserTeam
	userTeams, err := ut.WithContext(ctx).Where(ut.TeamID.Eq(teamID)).Find()
	if err != nil {
		return err
	}

	// 存在用户的ID
	var existUserIDs []string
	for _, userTeam := range userTeams {
		existUserIDs = append(existUserIDs, userTeam.UserID)
	}

	userMemo := make(map[string]struct{})
	for _, userId := range existUserIDs {
		userMemo[userId] = struct{}{}
	}

	// 用户团队关系
	count := len(members)
	createUserTeams := make([]*model.UserTeam, 0, count)
	// 用户团队角色关系
	createUserRoles := make([]*model.UserRole, 0, count)
	userIDs := make([]string, 0, count)
	for _, member := range members {
		// 只处理新增用户团队
		if _, ok := userMemo[member.UserID]; !ok {
			userTeam := &model.UserTeam{
				UserID:       member.UserID,
				TeamID:       teamID,
				InviteUserID: userID,
				InviteTime:   time.Now(),
			}
			createUserTeams = append(createUserTeams, userTeam)

			userRole := &model.UserRole{
				RoleID:       member.TeamRoleID,
				UserID:       member.UserID,
				TeamID:       teamID,
				InviteUserID: userID,
				InviteTime:   time.Now(),
			}
			createUserRoles = append(createUserRoles, userRole)

			userIDs = append(userIDs, member.UserID)
		}
	}

	var operationName string
	u := query.Use(dal.DB()).User
	userList, err := u.WithContext(ctx).Where(u.UserID.In(userIDs...)).Find()
	if err != nil {
		return err
	}
	userListMemo := make(map[string]*model.User)
	for _, user := range userList {
		userListMemo[user.UserID] = user
	}

	err = query.Use(dal.DB()).Transaction(func(tx *query.Query) error {
		if len(createUserTeams) > 0 {
			if err = tx.UserTeam.WithContext(ctx).Create(createUserTeams...); err != nil {
				return err
			}
		}
		if len(createUserRoles) > 0 {
			if err = tx.UserRole.WithContext(ctx).Create(createUserRoles...); err != nil {
				return err
			}
		}

		// 设置默认团队
		if len(userIDs) > 0 {
			for _, userID := range userIDs {
				_ = SettingCreateOrUpdate(ctx, userID, teamID, true)

				if user, ok := userListMemo[userID]; ok {
					operationName += user.Nickname + ","
				}
			}
		}

		return nil
	})

	// 触发事件
	if err, _ = gevent.Trigger(consts.ActionTeamSaveMember, event.BaseParams(ctx, teamID, userID, strings.Trim(operationName, ","))); err != nil {
		return err
	}

	return nil
}

func CompanyMembers(ctx context.Context, teamID string, keyword string, page int, size int) ([]*rao.Member, int64, string, error) {
	limit := size
	offset := (page - 1) * size

	t := query.Use(dal.DB()).Team
	team, err := t.WithContext(ctx).Where(t.TeamID.Eq(teamID)).First()
	if err != nil {
		return nil, 0, "", errmsg.ErrTeamNotFound
	}

	// 查询可用的用户团队数据
	ut := query.Use(dal.DB()).UserTeam
	userTeams, err := ut.WithContext(ctx).Where(ut.TeamID.Eq(teamID), ut.IsShow.Eq(consts.TeamIsShow)).Find()
	if err != nil {
		return nil, 0, "", err
	}

	// keyword 搜索昵称/账号
	u := dal.GetQuery().User
	conditions := make([]gen.Condition, 0)
	conditionsAccount := conditions
	keyword = strings.TrimSpace(keyword)
	if len(keyword) > 0 {
		conditions = append(conditions, u.Nickname.Like(fmt.Sprintf("%%%s%%", keyword)))
		conditionsAccount = append(conditionsAccount, u.Account.Like(fmt.Sprintf("%%%s%%", keyword)))
	}
	users, total, err := u.WithContext(ctx).Where(conditions...).Or(conditionsAccount...).Order(u.ID.Desc()).FindByPage(offset, limit)
	if err != nil {
		return nil, 0, "", err
	}

	userIDs := make([]string, 0, len(users))
	for _, u := range users {
		userIDs = append(userIDs, u.UserID)
	}

	// 企业成员
	uc := dal.GetQuery().UserCompany
	userCompanies, err := uc.WithContext(ctx).Where(uc.CompanyID.Eq(team.CompanyID), uc.UserID.In(userIDs...)).Order(uc.ID.Desc()).Find()
	if err != nil {
		return nil, 0, "", err
	}

	r := dal.GetQuery().Role
	roles, err := r.WithContext(ctx).Where(r.CompanyID.Eq(team.CompanyID)).Find()
	if err != nil {
		return nil, 0, "", err
	}

	ur := dal.GetQuery().UserRole
	userRoles, err := ur.WithContext(ctx).Where(ur.UserID.In(userIDs...)).Find()
	if err != nil {
		return nil, 0, "", err
	}

	return packer.TransTeamModelToRelatedMember(team, userTeams, users, roles, userRoles, userCompanies), total, team.CreatedUserID, nil

}

func TransferSuperRole(ctx context.Context, userID string, teamID string, targetUserID string) error {
	if userID == global.SuperManageUserID {
		r := dal.GetQuery().Role
		role, err := r.WithContext(ctx).Where(
			r.RoleType.Eq(consts.RoleTypeTeam),
			r.Level.Eq(consts.RoleLevelSuperManager)).First()
		if err != nil {
			return err
		}

		// 超管  超管是团队管理员
		// 超管  超管不是团队管理员
		ur := dal.GetQuery().UserRole
		userRole, err := ur.WithContext(ctx).Where(ur.UserID.Neq(userID), ur.RoleID.Eq(role.RoleID), ur.TeamID.Eq(teamID)).First()
		if err == nil {
			userID = userRole.UserID
		}
	}

	u := dal.GetQuery().User
	targetUser, err := u.WithContext(ctx).Where(u.UserID.Eq(targetUserID)).First()
	if err != nil {
		return err
	}

	ur := dal.GetQuery().UserRole
	userRole, err := ur.WithContext(ctx).Where(ur.UserID.Eq(userID), ur.TeamID.Eq(teamID)).First()
	if err != nil {
		return err
	}

	r := dal.GetQuery().Role
	role, err := r.WithContext(ctx).Where(r.RoleID.Eq(userRole.RoleID)).First()
	if err != nil {
		return err
	}

	if role.Level != consts.RoleLevelSuperManager {
		return errmsg.ErrUserNotRole
	}

	targetUserRole, err := ur.WithContext(ctx).Where(ur.UserID.Eq(targetUserID), ur.TeamID.Eq(teamID)).First()
	if err != nil {
		return err
	}

	err = query.Use(dal.DB()).Transaction(func(tx *query.Query) error {
		// 角色互换
		urt := tx.UserRole
		_, err = urt.WithContext(ctx).Where(urt.UserID.Eq(userID), urt.TeamID.Eq(teamID)).Update(urt.RoleID, targetUserRole.RoleID)
		if err != nil {
			return err
		}

		_, err = urt.WithContext(ctx).Where(urt.UserID.Eq(targetUserID), urt.TeamID.Eq(teamID)).Update(urt.RoleID, userRole.RoleID)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	// 触发事件
	if err, _ = gevent.Trigger(consts.ActionTeamTransfer, event.BaseParams(ctx, teamID, userID, targetUser.Nickname)); err != nil {
		return err
	}

	return nil
}
