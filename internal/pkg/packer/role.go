package packer

import (
	"permission-open/internal/pkg/biz/consts"
	"permission-open/internal/pkg/dal/model"
	"permission-open/internal/pkg/dal/rao"
)

func TransRolesModelToRaoRoles(roles []*model.Role) []*rao.Role {
	ret := make([]*rao.Role, 0, len(roles))
	for _, r := range roles {
		roleAttr := &rao.RoleAttr{
			IsUpdatePermission: true,
		}

		if r.Level == consts.RoleLevelSuperManager {
			roleAttr.IsUpdatePermission = false
		}

		// 管理员不可以删除
		if r.Level == consts.RoleLevelManager && r.RoleType == consts.RoleTypeCompany {
			roleAttr.IsUpdatePermission = false
		}

		ret = append(ret, &rao.Role{
			RoleID:    r.RoleID,
			RoleType:  r.RoleType,
			Name:      r.Name,
			Level:     r.Level,
			Attr:      roleAttr,
			IsDefault: r.IsDefault,
		})
	}

	return ret
}

func TransRoleMembersModelToRaoRole(
	role *model.Role,
	roles []*model.Role,
	userRoles []*model.UserRole,
	users []*model.User,
	teams []*model.Team,
	curUserRoleID string,
) []*rao.RoleMember {
	ret := make([]*rao.RoleMember, 0, len(userRoles))

	userMemo := make(map[string]*model.User)
	for _, user := range users {
		userMemo[user.UserID] = user
	}

	rolesMemo := make(map[string]*model.Role)
	for _, r := range roles {
		rolesMemo[r.RoleID] = r
	}

	teamMemo := make(map[string]*model.Team)
	for _, t := range teams {
		teamMemo[t.TeamID] = t
	}

	for _, userRole := range userRoles {
		if user, ok := userMemo[userRole.UserID]; ok {
			roleMember := &rao.RoleMember{}
			roleMember.UserID = user.UserID
			roleMember.Account = user.Account
			roleMember.Mobile = user.Mobile
			roleMember.Avatar = user.Avatar
			roleMember.Email = user.Email
			roleMember.Nickname = user.Nickname

			roleMember.RoleID = userRole.RoleID
			roleMember.RoleName = role.Name
			if !userRole.InviteTime.IsZero() {
				roleMember.InviteTimeSec = userRole.InviteTime.Unix()
			}
			if user, ok := userMemo[userRole.InviteUserID]; ok {
				roleMember.InviteUserID = user.UserID
				roleMember.InviteUserName = userMemo[user.UserID].Nickname
			}
			if team, ok := teamMemo[userRole.TeamID]; ok {
				roleMember.TeamID = team.TeamID
				roleMember.TeamName = team.Name
			}
			if len(curUserRoleID) > 0 {
				// 超管权限最高
				if rolesMemo[curUserRoleID].Level == consts.RoleLevelSuperManager {
					roleMember.IsOperableRole = true
				}
				if r, ok := rolesMemo[userRole.RoleID]; ok {
					if rolesMemo[curUserRoleID].Level <= r.Level {
						if r.Level != consts.RoleLevelManager {
							roleMember.IsOperableRole = true
						}
					}
					// 当前用户是超管，没有权限修改
					if r.Level == consts.RoleLevelSuperManager {
						roleMember.IsOperableRole = false
					}
				}
			}
			ret = append(ret, roleMember)
		}
	}

	return ret
}
