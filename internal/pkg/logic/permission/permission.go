package permission

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	gevent "github.com/gookit/event"
	"gorm.io/gen"
	"gorm.io/gorm"
	"permission-open/internal/pkg/biz/consts"
	"permission-open/internal/pkg/dal"
	"permission-open/internal/pkg/dal/model"
	"permission-open/internal/pkg/dal/query"
	"permission-open/internal/pkg/dal/rao"
	"permission-open/internal/pkg/event"
	"permission-open/internal/pkg/logic/errmsg"
	"permission-open/internal/pkg/packer"
	"permission-open/internal/pkg/public"
	"sort"
	"strings"
	"time"
)

// GetPermissionList 权限列表
func GetPermissionList(ctx *gin.Context, roleID string) ([]*rao.PermissionGroup, error) {
	r := dal.GetQuery().Role
	role, err := r.WithContext(ctx).Where(r.RoleID.Eq(roleID)).First()
	if err != nil {
		return nil, errmsg.ErrRoleExists
	}

	roleType := role.RoleType
	rtp := dal.GetQuery().RoleTypePermission
	rtpList, err := rtp.WithContext(ctx).Where(rtp.RoleType.Eq(roleType)).Find()
	if err != nil {
		return nil, err
	}

	var permissionIDs = make([]int64, 0, len(rtpList))
	for _, p := range rtpList {
		permissionIDs = append(permissionIDs, p.PermissionID)
	}

	p := dal.GetQuery().Permission
	permissions, err := p.WithContext(ctx).Where(p.PermissionID.In(permissionIDs...)).Find()
	if err != nil {
		return nil, err
	}

	havaPIDs := make([]int64, 0)
	if len(roleID) > 0 {
		rp := dal.GetQuery().RolePermission
		havePermission, err := rp.WithContext(ctx).Where(rp.RoleID.Eq(roleID)).Find()
		if err != nil {
			return nil, err
		}
		for _, permission := range havePermission {
			havaPIDs = append(havaPIDs, permission.PermissionID)
		}
	}

	ret := make([]*rao.PermissionGroup, 0)
	groupMemos := make(map[int32][]*rao.Permission)
	for _, p := range permissions {
		permission := &rao.Permission{
			Title:          p.Title,
			PermissionType: p.Type,
			GroupID:        p.GroupID,
			PermissionID:   p.PermissionID,
			Mark:           p.Mark,
		}
		for _, pid := range havaPIDs {
			if p.PermissionID == pid {
				permission.IsHave = true // 当前角色是否拥有此权限
			}
		}
		groupMemos[p.GroupID] = append(groupMemos[p.GroupID], permission)
	}

	for groupID, groupMemo := range groupMemos {
		ret = append(ret, &rao.PermissionGroup{
			GroupID:     groupID,
			GroupName:   consts.PerGroupNameMap[groupID],
			Mark:        consts.PerGroupMarkMap[groupID],
			Permissions: groupMemo,
		})
	}

	sort.Slice(ret, func(i, j int) bool {
		return ret[i].GroupID < ret[j].GroupID
	})

	return ret, nil
}

// SetRolePermission 设置角色权限
func SetRolePermission(ctx *gin.Context, userID string, roleID string, permissionMarks []string, roleName string) error {
	// step1：查询当前角色是否存在
	r := dal.GetQuery().Role
	role, err := r.WithContext(ctx).Where(r.RoleID.Eq(roleID)).First()
	if err != nil {
		return err
	}
	// 超管、管理员角色权限不可修改
	if role.Level == consts.RoleLevelSuperManager {
		return errmsg.ErrRoleNotChange
	}

	// step2：查询权限
	p := dal.GetQuery().Permission
	permissions, err := p.WithContext(ctx).Where(p.Mark.In(permissionMarks...)).Find()
	if err != nil {
		return err
	}

	// step3：过滤不存在的权限
	newPIDs := make([]int64, 0, len(permissionMarks))
	oldPIDs := make([]int64, 0, len(permissionMarks))
	for _, permission := range permissions {
		newPIDs = append(newPIDs, permission.PermissionID)
	}

	// step4：新的权限和旧的权限取差集   新多余旧新增 旧的不存在新的删除
	rp := dal.GetQuery().RolePermission
	rolePerOrm, err := rp.WithContext(ctx).Where(rp.RoleID.Eq(roleID)).Find()
	if err != nil {
		return err
	}

	for _, rp := range rolePerOrm {
		oldPIDs = append(oldPIDs, rp.PermissionID)
	}

	createPIDs := public.SliceDiff(newPIDs, oldPIDs)
	delPIDs := public.SliceDiff(oldPIDs, newPIDs)

	err = query.Use(dal.DB()).Transaction(func(tx *query.Query) error {
		// 新增
		if len(createPIDs) > 0 {
			createRolePer := make([]*model.RolePermission, 0, len(createPIDs))
			for _, pID := range createPIDs {
				createRolePer = append(createRolePer, &model.RolePermission{
					RoleID:       roleID,
					PermissionID: pID,
				})
			}
			err = tx.RolePermission.WithContext(ctx).Create(createRolePer...)
			if err != nil {
				return err
			}
		}

		// 删除
		if len(delPIDs) > 0 {
			_, err = tx.RolePermission.WithContext(ctx).Where(rp.PermissionID.In(delPIDs...), rp.RoleID.Eq(roleID)).Delete()
			if err != nil {
				return err
			}
		}

		// 修改团队名称
		if len(roleName) > 0 && role.Name != roleName {
			count, err := r.WithContext(ctx).Where(r.Name.Eq(roleName), r.CompanyID.Eq(role.CompanyID)).Count()
			if err != nil {
				return err
			}
			if count > 0 {
				return errmsg.ErrRoleExists
			}
			_, err = tx.Role.WithContext(ctx).Where(tx.Role.RoleID.Eq(roleID)).Update(tx.Role.Name, roleName)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	// 触发事件
	if err, _ = gevent.Trigger(consts.ActionRoleSet, event.BaseParams(ctx, "", userID, role.Name)); err != nil {
		return err
	}

	return nil
}

func GetUserPermissionList(ctx *gin.Context, userID string, companyID string, teamID string) ([]*rao.Permission, error) {
	// 查询企业角色权限   如果传团队ID，也同样包含团队角色权限
	ur := dal.GetQuery().UserRole
	// 企业角色
	userRole, err := ur.WithContext(ctx).Where(ur.UserID.Eq(userID), ur.CompanyID.Eq(companyID)).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errmsg.ErrUserNotRole
	}
	if err != nil {
		return nil, err
	}

	roles := make([]string, 0, 2)
	roles = append(roles, userRole.RoleID)
	if len(teamID) > 0 {
		userTeamRole, err := ur.WithContext(ctx).Where(ur.UserID.Eq(userID), ur.TeamID.Eq(teamID)).First()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errmsg.ErrUserNotRole
		}
		if err != nil {
			return nil, err
		}
		roles = append(roles, userTeamRole.RoleID)
	}

	rp := dal.GetQuery().RolePermission
	rolePermission, err := rp.WithContext(ctx).Where(rp.RoleID.In(roles...)).Find()
	if err != nil {
		return nil, err
	}

	havePIDs := make([]int64, 0, len(rolePermission))
	for _, permission := range rolePermission {
		havePIDs = append(havePIDs, permission.PermissionID)
	}

	p := dal.GetQuery().Permission
	permissions, err := p.WithContext(ctx).Where(p.PermissionID.In(havePIDs...)).Find()
	if err != nil {
		return nil, err
	}

	permissionsRao := make([]*rao.Permission, 0, len(havePIDs))
	for _, p := range permissions {
		permissionsRao = append(permissionsRao, &rao.Permission{
			IsHave:         true,
			PermissionType: p.Type,
			GroupID:        p.GroupID,
			PermissionID:   p.PermissionID,
			Title:          p.Title,
			Mark:           p.Mark,
		})
	}

	return permissionsRao, nil
}

func GetUserAllPermissionMarks(ctx *gin.Context, userID string) (*rao.UserAllPermissionMarksResp, error) {
	// step1 :获取用户所有角色ID
	ur := dal.GetQuery().UserRole
	userRoles, err := ur.WithContext(ctx).Where(ur.UserID.Eq(userID)).Find()
	if err != nil {
		return nil, err
	}

	roleIDs := make([]string, 0, len(userRoles))
	for _, role := range userRoles {
		roleIDs = append(roleIDs, role.RoleID)
	}

	// step2 :用户所有角色权限
	rp := dal.GetQuery().RolePermission
	rolePermission, err := rp.WithContext(ctx).Where(rp.RoleID.In(roleIDs...)).Find()
	if err != nil {
		return nil, err
	}

	havePIDs := make([]int64, 0, len(rolePermission))
	for _, p := range rolePermission {
		havePIDs = append(havePIDs, p.PermissionID)
	}

	// step3: 所有权限
	p := dal.GetQuery().Permission
	permissions, err := p.WithContext(ctx).Where(p.PermissionID.In(havePIDs...)).Find()
	if err != nil {
		return nil, err
	}
	return packer.TransMarksToRaoAllPermission(userRoles, rolePermission, permissions), nil
}

// GetPermissionListUrl 获取权限 URL
func GetPermissionListUrl(c *gin.Context) ([]string, error) {
	ret, err := dal.GetRDB().Get(c, consts.RedisPermissionListUrl).Result()
	if err != nil {
		// 如果返回的错误是key不存在
		if errors.Is(err, redis.Nil) {
			p := dal.GetQuery().Permission
			permissions, err := p.WithContext(c).Where().Find()
			if err != nil {
				return nil, err
			}
			var permissionsUrl = make([]string, 0, len(permissions))
			for _, p := range permissions {
				permissionsUrl = append(permissionsUrl, p.URL)
			}

			ret = strings.Join(permissionsUrl, ",")
			if err := dal.GetRDB().Set(c, consts.RedisPermissionListUrl, ret, time.Second*3600).Err(); err != nil {
				return nil, err
			}
		}
	}

	return strings.Split(ret, ","), nil
}

// CheckRolePermission 角色是否有当前权限
func CheckRolePermission(c *gin.Context, roleID string, permissionFunc string) (bool, error) {
	rp := dal.GetQuery().RolePermission
	rolePermissions, err := rp.WithContext(c).Where(rp.RoleID.Eq(roleID)).Find()
	if err != nil {
		return false, err
	}

	p := dal.GetQuery().Permission
	permissions, err := p.WithContext(c).Where().Find()
	if err != nil {
		return false, err
	}

	// step1: 路径是否存在权限控制
	// step2: 存在 -> 判断当前角色是否有这个权限
	var (
		isUrlExist   = false
		permissionID int64
	)
	for _, permission := range permissions {
		if permission.URL == permissionFunc {
			isUrlExist = true
			permissionID = permission.PermissionID
			break
		}
	}
	// step1: 当前 URL 权限控制范围内
	if !isUrlExist {
		return true, nil
	}

	// step2: 判断当前角色是否有这个权限
	for _, p := range rolePermissions {
		if p.PermissionID == permissionID {
			return true, nil
		}
	}

	return false, nil
}

// CheckUrl 通过 userID url 判断是否有权限
func CheckUrl(ctx *gin.Context, userID string, teamID string, permissionFunc string) (bool, error) {
	ur := dal.GetQuery().UserRole
	conditions := make([]gen.Condition, 0)
	conditions = append(conditions, ur.UserID.Eq(userID))

	// 判断是查团队权限 or 企业权限
	if len(teamID) > 0 {
		conditions = append(conditions, ur.TeamID.Eq(teamID))
	} else {
		uc := dal.GetQuery().UserCompany
		userCompany, err := uc.WithContext(ctx).Where(uc.UserID.Eq(userID)).First()
		if err != nil {
			return false, err
		}
		conditions = append(conditions, ur.CompanyID.Eq(userCompany.CompanyID))
	}
	userRole, err := ur.WithContext(ctx).Where(conditions...).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, errmsg.ErrUserNotRole
	}
	if err != nil {
		return false, err
	}

	return CheckRolePermission(ctx, userRole.RoleID, permissionFunc)
}

// Clear 清理历史数据
func Clear(ctx context.Context, userID, companyID, roleIDC2, roleIDT1, roleIDT2 string) error {
	if err := query.Use(dal.DB()).Transaction(func(tx *query.Query) error {
		u := tx.User
		users, err := u.WithContext(ctx).Where(u.Account.Eq("")).Find()
		if err != nil {
			return err
		}

		for _, user := range users {
			_, err = u.WithContext(ctx).Where(u.UserID.Eq(user.UserID)).Update(u.Account, u.Email)
			if err != nil {
				return err
			}
		}

		users, err = u.WithContext(ctx).Find()
		if err != nil {
			return err
		}
		userCompany := make([]*model.UserCompany, 0, len(users))
		userRole := make([]*model.UserRole, 0, len(users))
		uc := tx.UserCompany
		ur := tx.UserRole
		for _, user := range users {
			if user.UserID != userID {
				_, err = uc.WithContext(ctx).Where(uc.UserID.Eq(user.UserID), uc.CompanyID.Eq(companyID)).First()
				if errors.Is(err, gorm.ErrRecordNotFound) {
					userCompany = append(userCompany, &model.UserCompany{
						UserID:       user.UserID,
						CompanyID:    companyID,
						InviteUserID: userID,
						InviteTime:   user.CreatedAt,
					})
				}

				_, err = ur.WithContext(ctx).Where(ur.UserID.Eq(user.UserID), ur.CompanyID.Eq(companyID)).First()
				if errors.Is(err, gorm.ErrRecordNotFound) {
					userRole = append(userRole, &model.UserRole{
						RoleID:       roleIDC2,
						UserID:       user.UserID,
						CompanyID:    companyID,
						InviteUserID: userID,
						InviteTime:   user.CreatedAt,
					})
				}
			}
		}

		if len(userCompany) > 0 {
			if err = uc.WithContext(ctx).Create(userCompany...); err != nil {
				return err
			}
		}

		if len(userRole) > 0 {
			if err = ur.WithContext(ctx).Create(userRole...); err != nil {
				return err
			}
		}

		t := tx.Team
		if _, err = t.WithContext(ctx).Where(t.ID.Gte(0)).Update(t.CompanyID, companyID); err != nil {
			return err
		}

		ut := tx.UserTeam
		userTeams, err := ut.WithContext(ctx).Find()
		if err != nil {
			return err
		}

		userRoleTeam := make([]*model.UserRole, 0, len(userTeams))
		for _, team := range userTeams {
			roleID := roleIDT1
			if team.RoleID != 1 {
				roleID = roleIDT2
			}

			_, err = ur.WithContext(ctx).Where(ur.UserID.Eq(team.UserID), ur.TeamID.Eq(team.TeamID)).First()
			if errors.Is(err, gorm.ErrRecordNotFound) {
				userRoleTeam = append(userRoleTeam, &model.UserRole{
					RoleID:       roleID,
					UserID:       team.UserID,
					TeamID:       team.TeamID,
					InviteUserID: team.InviteUserID,
					InviteTime:   team.CreatedAt,
				})
			}
		}
		if len(userRoleTeam) > 0 {
			if err = ur.WithContext(ctx).Create(userRoleTeam...); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
