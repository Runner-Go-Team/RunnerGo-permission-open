package packer

import (
	"permission-open/internal/pkg/dal/mao"
	"permission-open/internal/pkg/dal/model"
	"permission-open/internal/pkg/dal/rao"
)

func TransOperationsToRaoOperationList(operations []*mao.CompanyOperationLog, users []*model.User) []*rao.Operation {
	ret := make([]*rao.Operation, 0)

	memo := make(map[string]*model.User)
	for _, user := range users {
		memo[user.UserID] = user
	}

	for _, operationInfo := range operations {
		ret = append(ret, &rao.Operation{
			UserID:         operationInfo.UserID,
			UserName:       memo[operationInfo.UserID].Nickname,
			UserAvatar:     memo[operationInfo.UserID].Avatar,
			Category:       operationInfo.Category,
			Action:         operationInfo.Action,
			Name:           operationInfo.Name,
			CreatedTimeSec: operationInfo.CreatedAt,
		})
	}

	return ret
}
