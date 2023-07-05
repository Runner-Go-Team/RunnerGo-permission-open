package packer

import (
	"permission-open/internal/pkg/dal/model"
	"permission-open/internal/pkg/dal/rao"
)

func TransUserInfoToRaoUser(
	user *model.User,
	userCompany *model.UserCompany,
	userTeams []*model.UserTeam,
	teams []*model.Team,
	userRoles []*model.UserRole,
	roles []*model.Role,
	userSettings *model.Setting,
) *rao.GetUserInfoResp {

	teamsMemo := make(map[string]*model.Team)
	for _, m := range teams {
		teamsMemo[m.TeamID] = m
	}

	rolesMemo := make(map[string]*model.Role)
	for _, rm := range roles {
		rolesMemo[rm.RoleID] = rm
	}

	userInfo := &rao.UserInfo{
		ID:       user.ID,
		UserID:   user.UserID,
		Account:  user.Account,
		Email:    user.Email,
		Mobile:   user.Mobile,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
	}
	for _, role := range userRoles {
		if role.CompanyID == userCompany.CompanyID {
			userInfo.RoleID = role.RoleID
			userInfo.RoleName = rolesMemo[role.RoleID].Name
		}
	}

	// team  role  user_role.team_id = team
	teamList := make([]*rao.UserTeam, 0, len(userTeams))
	for _, userRole := range userRoles {
		for _, userTeam := range userTeams {
			if userRole.TeamID == userTeam.TeamID {
				raoUserTeam := &rao.UserTeam{
					TeamID:   userTeam.TeamID,
					TeamName: teamsMemo[userTeam.TeamID].Name,
					RoleID:   userRole.RoleID,
					RoleName: rolesMemo[userRole.RoleID].Name,
				}
				if !userTeam.InviteTime.IsZero() {
					raoUserTeam.JoinTimeSec = userTeam.InviteTime.Unix()
				}
				teamList = append(teamList, raoUserTeam)
			}
		}
	}

	var settingTeamID string
	if userSettings != nil {
		settingTeamID = userSettings.TeamID
	}
	userRelated := &rao.UserRelated{
		SettingTeamID: settingTeamID,
		CompanyID:     userCompany.CompanyID,
	}

	return &rao.GetUserInfoResp{
		UserInfo:    userInfo,
		UserRelated: userRelated,
		TeamList:    teamList,
	}
}
