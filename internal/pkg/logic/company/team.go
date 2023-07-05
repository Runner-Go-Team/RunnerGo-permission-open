package company

import (
	"github.com/gin-gonic/gin"
	"permission-open/internal/pkg/dal"
	"permission-open/internal/pkg/dal/query"
	"permission-open/internal/pkg/dal/rao"
)

func TeamsList(ctx *gin.Context, userID string, companyId string) ([]*rao.Team, error) {
	ut := query.Use(dal.DB()).UserTeam
	userTeams, err := ut.WithContext(ctx).Where(ut.UserID.Eq(userID)).Find()
	if err != nil {
		return nil, err
	}

	teamIDs := make([]string, 0, len(userTeams))
	for _, t := range userTeams {
		teamIDs = append(teamIDs, t.TeamID)
	}

	t := query.Use(dal.DB()).Team
	teams, err := t.WithContext(ctx).Where(t.CompanyID.Eq(companyId), t.TeamID.In(teamIDs...)).Find()
	if err != nil {
		return nil, err
	}

	teamInfos := make([]*rao.Team, 0, len(teams))
	for _, team := range teams {
		teamInfos = append(teamInfos, &rao.Team{
			Type:           team.Type,
			TeamID:         team.TeamID,
			Name:           team.Name,
			Description:    team.Description,
			CreatedUserID:  team.CreatedUserID,
			CreatedTimeSec: team.CreatedAt.Unix(),
			UpdatedTimeSec: team.UpdatedAt.Unix(),
		})
	}

	return teamInfos, err
}
