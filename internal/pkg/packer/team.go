package packer

import (
	"permission-open/internal/pkg/biz/consts"
	"permission-open/internal/pkg/dal/global"
	"permission-open/internal/pkg/dal/model"
	"permission-open/internal/pkg/dal/rao"
)

func TransTeamModelToRaoTeam(
	userID string,
	teamInfo *model.Team,
	userTeams []*model.UserTeam,
	users []*model.User,
	roles []*model.Role,
	userRoles []*model.UserRole,
	userCompanies []*model.UserCompany,
) *rao.Team {
	var (
		curUserCompanyRoleID string   // 当前用户的企业角色ID
		curUserTeamRoleID    string   // 当前用户的角色角色ID
		TeamSuperRoleUserIDs []string // 团队的超级管理员 userID
	)

	for _, userRole := range userRoles {
		// 企业角色
		if userID == userRole.UserID && userRole.CompanyID != "" {
			curUserCompanyRoleID = userRole.RoleID
		}
		// 团队角色
		if userID == userRole.UserID && userRole.TeamID == teamInfo.TeamID {
			curUserTeamRoleID = userRole.RoleID

			for _, r := range roles {
				if r.Level == consts.RoleLevelSuperManager && r.RoleType == consts.RoleTypeTeam {
					TeamSuperRoleUserIDs = append(TeamSuperRoleUserIDs, userRole.UserID)
				}
			}
		}
	}

	userMemo := make(map[string]*model.User)
	for _, user := range users {
		userMemo[user.UserID] = user
	}

	rolesMemo := make(map[string]*model.Role)
	for _, role := range roles {
		rolesMemo[role.RoleID] = role
	}

	userCompaniesMemo := make(map[string]*model.UserCompany)
	for _, userCompany := range userCompanies {
		userCompaniesMemo[userCompany.UserID] = userCompany
	}

	teamMemberMemo := make(map[string][]*rao.Member)
	for _, team := range userTeams {
		// 用户ID不存在或已经被移除
		user, ok := userMemo[team.UserID]
		if !ok {
			continue
		}
		member := &rao.Member{}
		member.UserID = user.UserID
		member.Avatar = user.Avatar
		member.Nickname = user.Nickname
		member.Account = user.Account
		member.Mobile = user.Mobile
		member.Email = user.Email
		member.Status = userCompaniesMemo[user.UserID].Status

		for _, userRole := range userRoles {
			// 企业角色
			if team.UserID == userRole.UserID && userRole.CompanyID != "" {
				member.CompanyRoleID = userRole.RoleID
				member.CompanyRoleName = rolesMemo[userRole.RoleID].Name
				if len(curUserCompanyRoleID) > 0 {
					if rolesMemo[curUserCompanyRoleID].Level <= rolesMemo[userRole.RoleID].Level {
						if rolesMemo[userRole.RoleID].Level != consts.RoleLevelManager {
							member.IsOperableCompanyRole = true
						}
					}
				}
			}
			// 团队角色
			if team.UserID == userRole.UserID && team.TeamID == userRole.TeamID {
				// 移交团队管理员权限
				if len(curUserTeamRoleID) > 0 {
					if rolesMemo[curUserTeamRoleID].Level == consts.RoleLevelSuperManager {
						if userRole.UserID != userID {
							member.IsTransferSuperTeam = true
						}
					}
				}
				member.TeamRoleID = userRole.RoleID
				member.TeamRoleName = rolesMemo[userRole.RoleID].Name
				if len(curUserTeamRoleID) > 0 {
					if rolesMemo[curUserTeamRoleID].Level == consts.RoleLevelSuperManager {
						member.IsOperableTeamRole = true
						member.IsOperableRemoveMember = true
					}
					if rolesMemo[curUserTeamRoleID].Level <= rolesMemo[userRole.RoleID].Level {
						member.IsOperableTeamRole = true
						member.IsOperableRemoveMember = true
					}
				}
			}
		}
		if iUser, ok := userMemo[team.InviteUserID]; ok {
			member.InviteUserID = team.InviteUserID
			member.InviteUserName = iUser.Nickname
		}
		if !team.InviteTime.IsZero() {
			member.JoinTimeSec = team.InviteTime.Unix()
		}

		// 超管特殊处理
		if userID == global.SuperManageUserID {
			member.IsOperableCompanyRole = true
			member.IsOperableTeamRole = true
			member.IsOperableRemoveMember = true
			member.IsTransferSuperTeam = true
		}

		// 当前是团队管理员不展示
		if len(TeamSuperRoleUserIDs) > 0 {
			for _, superUserID := range TeamSuperRoleUserIDs {
				if team.UserID == superUserID {
					member.IsTransferSuperTeam = false
					member.IsOperableTeamRole = false
					member.IsOperableRemoveMember = false
				}
			}
		}

		teamMemberMemo[team.TeamID] = append(teamMemberMemo[team.TeamID], member)
	}

	team := &rao.Team{
		TeamID:          teamInfo.TeamID,
		Name:            teamInfo.Name,
		Type:            teamInfo.Type,
		Description:     teamInfo.Description,
		CreatedUserID:   userMemo[teamInfo.CreatedUserID].UserID,
		CreatedUserName: userMemo[teamInfo.CreatedUserID].Nickname,
		CreatedTimeSec:  teamInfo.CreatedAt.Unix(),
		UpdatedTimeSec:  teamInfo.UpdatedAt.Unix(),
		Members:         teamMemberMemo[teamInfo.TeamID],
	}

	return team
}

func TransTeamModelToRelatedMember(
	teamInfo *model.Team,
	userTeams []*model.UserTeam,
	users []*model.User,
	roles []*model.Role,
	userRoles []*model.UserRole,
	userCompanies []*model.UserCompany,
) []*rao.Member {

	inTeam := make(map[string]struct{})
	for _, team := range userTeams {
		inTeam[team.UserID] = struct{}{}
	}

	userMemo := make(map[string]*model.User)
	for _, user := range users {
		userMemo[user.UserID] = user
	}

	rolesMemo := make(map[string]*model.Role)
	for _, role := range roles {
		rolesMemo[role.RoleID] = role
	}

	userCompaniesMemo := make(map[string]*model.UserCompany)
	for _, userCompany := range userCompanies {
		userCompaniesMemo[userCompany.UserID] = userCompany
	}

	members := make([]*rao.Member, 0, len(userCompanies))
	for _, c := range userCompanies {
		// 超管不展示
		if c.UserID == global.SuperManageUserID {
			continue
		}
		if user, ok := userMemo[c.UserID]; ok {
			member := &rao.Member{}
			member.UserID = user.UserID
			member.Avatar = user.Avatar
			member.Nickname = user.Nickname
			member.Account = user.Account
			member.Mobile = user.Mobile
			member.Email = user.Email
			member.Status = userCompaniesMemo[user.UserID].Status

			for _, userRole := range userRoles {
				// 企业角色
				if c.UserID == userRole.UserID && userRole.CompanyID != "" {
					member.CompanyRoleID = userRole.RoleID
					member.CompanyRoleName = rolesMemo[userRole.RoleID].Name
				}
			}

			// 如果在团队中
			if _, ok = inTeam[c.UserID]; ok {
				for _, userRole := range userRoles {
					// 团队角色
					if c.UserID == userRole.UserID && teamInfo.TeamID == userRole.TeamID {
						member.TeamRoleID = userRole.RoleID
						member.TeamRoleName = rolesMemo[userRole.RoleID].Name
					}
				}
			}
			members = append(members, member)
		}
	}

	return members
}

func TransUserTeamsModelToRaoTeam(
	teams []*model.Team,
	userTeams []*model.UserTeam,
	users []*model.User,
	userTeamCollection []*model.UserTeamCollection,
) []*rao.RelatedTeam {
	ret := make([]*rao.RelatedTeam, 0)

	userTeamMemo := make(map[string]*model.UserTeam)
	for _, team := range userTeamMemo {
		userTeamMemo[team.TeamID] = team
	}

	userMemo := make(map[string]*model.User)
	for _, user := range users {
		userMemo[user.UserID] = user
	}

	teamMemo := make(map[string]*model.Team)
	for _, t := range teams {
		teamMemo[t.TeamID] = t
	}

	teamCollectionMemo := make(map[string]*model.UserTeamCollection)
	for _, tc := range userTeamCollection {
		teamCollectionMemo[tc.TeamID] = tc
	}

	teamMemberMemo := make(map[string][]*rao.Member)
	for _, team := range userTeams {
		user, ok := userMemo[team.UserID]
		if !ok {
			continue
		}
		member := &rao.Member{}
		member.UserID = user.UserID
		member.Avatar = user.Avatar
		member.Nickname = user.Nickname
		member.Account = user.Account
		member.Mobile = user.Mobile
		member.Email = user.Email
		if iUser, ok := userMemo[team.InviteUserID]; ok {
			member.InviteUserID = team.InviteUserID
			member.InviteUserName = iUser.Nickname
		}
		if !team.InviteTime.IsZero() {
			member.JoinTimeSec = team.InviteTime.Unix()
		}
		teamMemberMemo[team.TeamID] = append(teamMemberMemo[team.TeamID], member)
	}

	for _, t := range teams {
		team := &rao.Team{
			Name:            t.Name,
			Type:            t.Type,
			TeamID:          t.TeamID,
			CreatedUserID:   t.CreatedUserID,
			Description:     t.Description,
			CreatedUserName: userMemo[t.CreatedUserID].Nickname,
			CreatedTimeSec:  t.CreatedAt.Unix(),
			UpdatedTimeSec:  t.UpdatedAt.Unix(),
			Members:         teamMemberMemo[t.TeamID],
		}
		isCollect := false
		if _, ok := teamCollectionMemo[t.TeamID]; ok {
			isCollect = true
		}
		ret = append(ret, &rao.RelatedTeam{
			IsCollect: isCollect,
			Team:      team,
		})
	}

	return ret
}
