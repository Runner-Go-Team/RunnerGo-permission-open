package packer

import (
	"permission-open/internal/pkg/biz/consts"
	"permission-open/internal/pkg/dal/model"
	"permission-open/internal/pkg/dal/rao"
)

func TransCompaniesModelToRaoCompanyMember(
	userCompanies []*model.UserCompany,
	users []*model.User,
	roles []*model.Role,
	userRoles []*model.UserRole,
	curUserCompanyRoleID string,
) []*rao.CompanyMember {
	// step1: user_company
	// step2: user
	// step3: role
	// step4: user_role
	ret := make([]*rao.CompanyMember, 0, len(userCompanies))

	userCompanyMemo := make(map[string]*model.UserCompany)
	for _, userCompany := range userCompanies {
		userCompanyMemo[userCompany.UserID] = userCompany
	}

	userMemo := make(map[string]*model.User)
	for _, user := range users {
		userMemo[user.UserID] = user
	}

	roleMemo := make(map[string]*model.Role)
	for _, role := range roles {
		roleMemo[role.RoleID] = role
	}

	userRoleMemo := make(map[string]string)
	for _, r := range userRoles {
		userRoleMemo[r.UserID] = r.RoleID
	}

	for _, i := range userCompanies {
		user, ok := userMemo[i.UserID]
		if !ok {
			continue
		}

		raoCompanyMember := &rao.CompanyMember{}
		raoCompanyMember.UserID = user.UserID
		raoCompanyMember.Avatar = user.Avatar
		raoCompanyMember.Nickname = user.Nickname
		raoCompanyMember.Account = user.Account
		raoCompanyMember.Status = i.Status

		// 和当前角色判断权限
		if len(curUserCompanyRoleID) > 0 {
			// 超管权限最高
			if roleMemo[curUserCompanyRoleID].Level == consts.RoleLevelSuperManager {
				raoCompanyMember.IsOperableCompanyRole = true
				raoCompanyMember.IsOperableDisableMember = true
				raoCompanyMember.IsOperableRemoveMember = true
				raoCompanyMember.IsOperableUptPassword = true
			}
			if roleMemo[curUserCompanyRoleID].Level <= roleMemo[userRoleMemo[user.UserID]].Level {
				if roleMemo[userRoleMemo[user.UserID]].Level != consts.RoleLevelManager {
					raoCompanyMember.IsOperableCompanyRole = true
					raoCompanyMember.IsOperableDisableMember = true
					raoCompanyMember.IsOperableRemoveMember = true
				}
			}
		}
		// 当前用户是超管，没有权限修改
		if roleMemo[userRoleMemo[user.UserID]].Level == consts.RoleLevelSuperManager {
			raoCompanyMember.IsOperableCompanyRole = false
			raoCompanyMember.IsOperableDisableMember = false
			raoCompanyMember.IsOperableRemoveMember = false
		}

		if inviteUser, ok := userMemo[i.InviteUserID]; ok {
			raoCompanyMember.InviteUserID = i.InviteUserID
			raoCompanyMember.InviteUserName = inviteUser.Nickname
		}
		if role, ok := roleMemo[userRoleMemo[i.UserID]]; ok {
			raoCompanyMember.RoleId = role.RoleID
			raoCompanyMember.RoleName = role.Name
			raoCompanyMember.RoleLevel = role.Level
		}
		if !i.InviteTime.IsZero() {
			raoCompanyMember.InviteTimeSec = i.InviteTime.Unix()
		}
		raoCompanyMember.CreatedTimeSec = i.CreatedAt.Unix()

		ret = append(ret, raoCompanyMember)
	}

	return ret
}
