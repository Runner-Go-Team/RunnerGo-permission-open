package auth

import (
	"context"
	"permission-open/internal/pkg/dal"
	"permission-open/internal/pkg/dal/query"
	"permission-open/internal/pkg/dal/rao"
)

func SetUserSettings(ctx context.Context, userID string, settings *rao.UserSettings) error {
	currentTeamID := settings.CurrentTeamID
	tx := query.Use(dal.DB()).Setting
	_, err := tx.WithContext(ctx).Where(tx.UserID.Eq(userID)).UpdateColumnSimple(tx.TeamID.Value(currentTeamID))
	if err != nil {
		return err
	}

	return nil
}

func GetUserSettings(ctx context.Context, userID string) (*rao.GetUserSettingsResp, error) {
	tx := query.Use(dal.DB()).Setting
	settingInfo, err := tx.WithContext(ctx).Where(tx.UserID.Eq(userID)).First()
	if err != nil {
		return nil, err
	}

	return &rao.GetUserSettingsResp{
		UserSettings: &rao.UserSettings{
			CurrentTeamID: settingInfo.TeamID,
		},
	}, nil
}

// GetAvailTeamID 获取有效的团队ID
func GetAvailTeamID(ctx context.Context, userID string) (string, error) {
	//获取用户最后一次使用的团队
	tx := query.Use(dal.DB()).Setting
	s, err := tx.WithContext(ctx).Where(tx.UserID.Eq(userID)).First()
	if err != nil {
		return "", err
	}
	lastOperationTeamID := s.TeamID
	return lastOperationTeamID, nil
}
