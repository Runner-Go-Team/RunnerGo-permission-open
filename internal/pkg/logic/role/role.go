package role

import (
	"context"
	"errors"
	"fmt"
	gevent "github.com/gookit/event"
	"gorm.io/gen"
	"gorm.io/gorm"
	"permission-open/internal/pkg/biz/consts"
	"permission-open/internal/pkg/biz/jwt"
	"permission-open/internal/pkg/biz/uuid"
	"permission-open/internal/pkg/dal"
	"permission-open/internal/pkg/dal/global"
	"permission-open/internal/pkg/dal/model"
	"permission-open/internal/pkg/dal/query"
	"permission-open/internal/pkg/dal/rao"
	"permission-open/internal/pkg/event"
	"permission-open/internal/pkg/logic/errmsg"
	"permission-open/internal/pkg/packer"
	"permission-open/internal/pkg/public"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// GetRoleList 获取角色列表
func GetRoleList(ctx *gin.Context, req rao.RoleListReq) ([]*rao.Role, error) {
	if len(req.CompanyId) == 0 {
		req.CompanyId = global.CompanyID
	}

	r := dal.GetQuery().Role
	conditions := make([]gen.Condition, 0)
	conditions = append(conditions, r.CompanyID.Eq(req.CompanyId))

	// 角色分类（1：企业  2：团队）
	if req.RoleType > 0 {
		conditions = append(conditions, r.RoleType.Eq(req.RoleType))
	}
	// 角色层级(相当于查看权限--可选项根据操作者角色展示（如超管可见“管理员”选项，普通成员就不会看到）)
	if req.Level > 0 {
		conditions = append(conditions, r.Level.Gte(req.Level))
	}
	role, err := r.WithContext(ctx).Where(conditions...).Find()
	if err != nil {
		return nil, err
	}
	return packer.TransRolesModelToRaoRoles(role), nil
}

func SaveRole(ctx *gin.Context, userID string, req rao.SaveRoleReq) error {
	r := dal.GetQuery().Role
	count, err := r.WithContext(ctx).Where(r.Name.Eq(req.Name), r.CompanyID.Eq(req.CompanyID)).Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return errmsg.ErrRoleExists
	}

	level := consts.RoleLevelGeneral
	if req.RoleType == consts.RoleTypeTeam {
		level = consts.RoleLevelManager
	}

	err = r.WithContext(ctx).Create(&model.Role{
		RoleID:    uuid.GetUUID(),
		RoleType:  req.RoleType,
		Name:      req.Name,
		CompanyID: req.CompanyID,
		Level:     int32(level),
		IsDefault: consts.RoleCustom,
	})
	if err != nil {
		return err
	}

	// 触发事件
	if err, _ = gevent.Trigger(consts.ActionRoleSave, event.BaseParams(ctx, "", userID, req.Name)); err != nil {
		return err
	}

	return err
}

func SetCompanyMemberRole(ctx *gin.Context, userID string, req rao.SetCompanyMemberRoleReq) error {
	r := dal.GetQuery().Role
	ur := dal.GetQuery().UserRole
	role, err := r.WithContext(ctx).Where(r.RoleID.Eq(req.RoleID)).First()
	if err != nil {
		return err
	}

	// 角色不是企业角色
	if role.RoleType != consts.RoleTypeCompany {
		return errmsg.ErrRoleNotChange
	}

	// 角色是超管，不能修改
	if role.Level == consts.RoleLevelSuperManager {
		return errmsg.ErrRoleNotChange
	}

	// 角色修改为管理员  只有超管有
	if role.Level == consts.RoleLevelManager {
		if userID != global.SuperManageUserID {
			return errmsg.ErrUserForbidden
		}
	}

	u := dal.GetQuery().User
	targetUser, err := u.WithContext(ctx).Where(u.UserID.Eq(req.TargetUserID)).First()
	if err != nil {
		return err
	}

	userRole, err := ur.WithContext(ctx).Where(ur.UserID.Eq(userID), ur.CompanyID.Eq(req.CompanyID)).First()
	if err != nil {
		return err
	}
	curRole, err := r.WithContext(ctx).Where(r.RoleID.Eq(userRole.RoleID)).First()
	if err != nil {
		return err
	}

	tarUserRole, err := ur.WithContext(ctx).Where(ur.UserID.Eq(req.TargetUserID), ur.CompanyID.Eq(req.CompanyID)).First()
	if err != nil {
		return err
	}
	tarRole, err := r.WithContext(ctx).Where(r.RoleID.Eq(tarUserRole.RoleID)).First()
	if err != nil {
		return err
	}

	if tarRole.Level == consts.RoleLevelSuperManager {
		return errmsg.ErrRoleNotChange
	}

	// 当前用户角色级别 小于 被操作人不允许
	if tarRole.Level < curRole.Level {
		return errmsg.ErrUserForbidden
	}

	_, err = ur.WithContext(ctx).Where(ur.CompanyID.Eq(req.CompanyID), ur.UserID.Eq(req.TargetUserID)).Updates(
		&model.UserRole{
			RoleID:       req.RoleID,
			InviteUserID: userID,
			InviteTime:   time.Now(),
		})
	if err != nil {
		return err
	}

	// 触发事件
	if err, _ = gevent.Trigger(consts.ActionCompanySetRoleMember, event.BaseParams(ctx, "", userID, targetUser.Nickname+" "+role.Name)); err != nil {
		return err
	}

	return nil
}

func SetTeamMemberRole(ctx *gin.Context, userID string, req rao.SetTeamMemberRoleReq) error {
	r := dal.GetQuery().Role
	ur := dal.GetQuery().UserRole
	role, err := r.WithContext(ctx).Where(r.RoleID.Eq(req.RoleID)).First()
	if err != nil {
		return err
	}

	// 角色不是团队角色
	if role.RoleType != consts.RoleTypeTeam {
		return errmsg.ErrRoleNotChange
	}

	// 角色是团队管理员，不能修改
	if role.Level == consts.RoleLevelSuperManager {
		return errmsg.ErrRoleNotChange
	}

	userRole, err := ur.WithContext(ctx).Where(ur.UserID.Eq(userID), ur.TeamID.Eq(req.TeamID)).First()
	if err != nil {
		return err
	}
	curRole, err := r.WithContext(ctx).Where(r.RoleID.Eq(userRole.RoleID)).First()
	if err != nil {
		return err
	}

	tarUserRole, err := ur.WithContext(ctx).Where(ur.UserID.Eq(req.TargetUserID), ur.TeamID.Eq(req.TeamID)).First()
	if err != nil {
		return err
	}
	tarRole, err := r.WithContext(ctx).Where(r.RoleID.Eq(tarUserRole.RoleID)).First()
	if err != nil {
		return err
	}

	u := dal.GetQuery().User
	targetUser, err := u.WithContext(ctx).Where(u.UserID.Eq(req.TargetUserID)).First()
	if err != nil {
		return err
	}

	// 角色是团队管理员，不能修改
	if tarRole.Level == consts.RoleLevelSuperManager {
		return errmsg.ErrRoleNotChange
	}

	// 当前用户角色级别 小于 被操作人不允许
	if tarRole.Level < curRole.Level {
		return errmsg.ErrUserForbidden
	}

	_, err = ur.WithContext(ctx).Where(ur.TeamID.Eq(req.TeamID), ur.UserID.Eq(req.TargetUserID)).Updates(
		&model.UserRole{
			RoleID:       req.RoleID,
			InviteUserID: userID,
			InviteTime:   time.Now(),
		})
	if err != nil {
		return err
	}

	// 触发事件
	if err, _ = gevent.Trigger(consts.ActionTeamSetRoleMember, event.BaseParams(ctx, req.TeamID, userID, targetUser.Nickname+" "+role.Name)); err != nil {
		return err
	}

	return nil
}

// IsAllowRemove 判断当前角色是否允许被删除
func IsAllowRemove(ctx *gin.Context, roleID string) (bool, error) {
	r := dal.GetQuery().Role
	role, err := r.WithContext(ctx).Where(r.RoleID.Eq(roleID)).First()
	if err != nil {
		return false, err
	}

	// 超管、管理员、团队管理员、团队成员不可以删除
	if role.Level == consts.RoleLevelSuperManager {
		return false, nil
	}

	if role.Level == consts.RoleLevelManager && role.RoleType == consts.RoleTypeCompany {
		return false, nil
	}

	// 默认角色不可删除
	if role.IsDefault == consts.RoleDefault {
		return false, nil
	}

	return true, err
}

func RemoveRole(ctx *gin.Context, userID string, roleID string, changeRoleID string, companyID string) error {
	r := dal.GetQuery().Role
	role, err := r.WithContext(ctx).Where(r.RoleID.Eq(roleID)).First()
	if err != nil {
		return err
	}

	// 超管、管理员、团队管理员、团队成员不可以删除
	if role.Level == consts.RoleLevelSuperManager {
		return errmsg.ErrRoleNotDel
	}

	if role.Level == consts.RoleLevelManager && role.RoleType == consts.RoleTypeCompany {
		return errmsg.ErrRoleNotDel
	}

	changeRole, err := r.WithContext(ctx).Where(r.RoleID.Eq(changeRoleID)).First()
	if err != nil {
		return err
	}
	// 超管、管理员全局唯一
	if changeRole.Level == consts.RoleLevelSuperManager {
		return errmsg.ErrRoleNotChange
	}

	// 默认角色不可删除
	if role.IsDefault == consts.RoleDefault {
		return errmsg.ErrRoleNotDel
	}

	err = query.Use(dal.DB()).Transaction(func(tx *query.Query) error {
		_, err = tx.Role.WithContext(ctx).Where(tx.Role.RoleID.Eq(roleID)).Delete()
		if err != nil {
			return err
		}

		// 角色关联成新的角色
		_, err = tx.UserRole.WithContext(ctx).Where(tx.UserRole.RoleID.Eq(roleID)).Updates(&model.UserRole{
			RoleID:       changeRoleID,
			InviteUserID: userID,
			InviteTime:   time.Now(),
		})

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	// 触发事件
	if err, _ = gevent.Trigger(consts.ActionRoleRemove, event.BaseParams(ctx, "", userID, role.Name)); err != nil {
		return err
	}

	return nil
}

func GetRoleMembers(ctx *gin.Context, req *rao.RoleMembersReq) ([]*rao.RoleMember, int64, error) {
	roleID := req.RoleID
	teamID := req.TeamID
	limit := req.Size
	offset := (req.Page - 1) * req.Size
	keyword := req.Keyword

	curUserID := jwt.GetUserIDByCtx(ctx)

	// keyword 搜索昵称/账号
	u := dal.GetQuery().User
	conditions := make([]gen.Condition, 0)
	conditionsAccount := conditions
	keyword = strings.TrimSpace(keyword)
	if len(keyword) > 0 {
		conditions = append(conditions, u.Nickname.Like(fmt.Sprintf("%%%s%%", keyword)))
		conditionsAccount = append(conditionsAccount, u.Account.Like(fmt.Sprintf("%%%s%%", keyword)))
	}
	users, err := u.WithContext(ctx).Where(conditions...).Or(conditionsAccount...).Order(u.ID.Desc()).Find()
	if err != nil {
		return nil, 0, err
	}

	userIDs := make([]string, 0, len(users))
	for _, u := range users {
		userIDs = append(userIDs, u.UserID)
	}

	r := dal.GetQuery().Role
	role, err := r.WithContext(ctx).Where(r.RoleID.Eq(roleID)).First()
	if err != nil {
		return nil, 0, err
	}

	uc := dal.GetQuery().UserRole
	conditions = make([]gen.Condition, 0)
	conditions = append(conditions, uc.RoleID.Eq(roleID))
	conditions = append(conditions, uc.UserID.In(userIDs...))
	if len(teamID) > 0 {
		conditions = append(conditions, uc.TeamID.Eq(teamID))
	}
	userRoles, total, err := uc.WithContext(ctx).Where(conditions...).Order(uc.InviteTime).FindByPage(offset, limit)
	if err != nil {
		return nil, 0, err
	}

	teamUserIDs := make([]string, 0, len(userRoles))
	for _, userRole := range userRoles {
		teamUserIDs = append(teamUserIDs, userRole.UserID)
	}

	// 团队管理员超管默认存在，需要过滤
	filterUserRoles := make([]*model.UserRole, 0, len(userRoles))
	userTeams := make([]*model.UserTeam, 0, len(userRoles))
	if len(teamID) > 0 {
		ut := dal.GetQuery().UserTeam
		userTeams, err = ut.WithContext(ctx).Where(
			ut.UserID.In(teamUserIDs...),
			ut.TeamID.Eq(teamID),
			ut.IsShow.Eq(consts.TeamIsShow),
		).Find()
		if err != nil {
			return nil, 0, err
		}
	}

	for _, r := range userRoles {
		if len(teamID) > 0 {
			for _, team := range userTeams {
				if team.UserID == r.UserID {
					filterUserRoles = append(filterUserRoles, r)
				}
			}
		} else {
			filterUserRoles = append(filterUserRoles, r)
		}
	}

	userIDs = make([]string, 0, len(userRoles))
	teamIDs := make([]string, 0, len(userRoles))
	for _, userRole := range filterUserRoles {
		userIDs = append(userIDs, userRole.UserID)
		userIDs = append(userIDs, userRole.InviteUserID)
		if len(userRole.TeamID) > 0 {
			teamIDs = append(teamIDs, userRole.TeamID)
		}
	}
	userIDs = public.SliceUnique(userIDs)

	t := dal.GetQuery().Team
	teams, err := t.WithContext(ctx).Where(t.TeamID.In(teamIDs...)).Find()
	if err != nil {
		return nil, 0, err
	}

	curUserRoles, err := uc.WithContext(ctx).Where(uc.UserID.Eq(curUserID)).Find()
	if err != nil {
		return nil, 0, err
	}

	var curUserRoleID string
	for _, r := range curUserRoles {
		if len(teamID) > 0 {
			if r.TeamID == teamID {
				curUserRoleID = r.RoleID
				break
			}
		} else {
			if len(r.CompanyID) > 0 {
				curUserRoleID = r.RoleID
				break
			}
		}
	}

	roles, err := r.WithContext(ctx).Find()
	if err != nil {
		return nil, 0, err
	}

	userAll, err := u.WithContext(ctx).Where(u.UserID.In(userIDs...)).Find()
	if err != nil {
		return nil, 0, err
	}

	return packer.TransRoleMembersModelToRaoRole(role, roles, filterUserRoles, userAll, teams, curUserRoleID), total, nil
}

func GetRoleMember(ctx context.Context, userID string, req *rao.RoleMemberReq) (*rao.Role, error) {
	var roleID string
	if len(req.RoleID) > 0 {
		roleID = req.RoleID
	} else {
		ur := dal.GetQuery().UserRole
		conditions := make([]gen.Condition, 0)
		conditions = append(conditions, ur.UserID.Eq(userID))

		// 角色分类（1：企业  2：团队）
		if len(req.TeamID) > 0 {
			conditions = append(conditions, ur.TeamID.Eq(req.TeamID))
		} else {
			conditions = append(conditions, ur.CompanyID.Eq(req.CompanyID))
		}

		userRole, err := ur.WithContext(ctx).Where(conditions...).First()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errmsg.ErrUserNotRole
		}
		if err != nil {
			return nil, err
		}
		roleID = userRole.RoleID
	}

	r := dal.GetQuery().Role
	role, err := r.WithContext(ctx).Where(r.RoleID.Eq(roleID)).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errmsg.ErrUserNotRole
	}
	if err != nil {
		return nil, err
	}

	return &rao.Role{
		RoleType: role.RoleType,
		Level:    role.Level,
		RoleID:   role.RoleID,
		Name:     role.Name,
	}, nil
}

// GetRoleTeamLevelManagerPermission 获取团队成员默认权限
func GetRoleTeamLevelManagerPermission() []int64 {
	return []int64{203}
}
