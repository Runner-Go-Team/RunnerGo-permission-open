package packer

import (
	"permission-open/internal/pkg/dal/model"
	"permission-open/internal/pkg/dal/rao"
)

func TransMarksToRaoAllPermission(
	userRoles []*model.UserRole,
	rolePermission []*model.RolePermission,
	permissions []*model.Permission,
) *rao.UserAllPermissionMarksResp {

	pMemo := make(map[int64]*model.Permission)
	for _, p := range permissions {
		pMemo[p.PermissionID] = p
	}

	rpMemo := make(map[string][]string)
	for _, rp := range rolePermission {
		rpMemo[rp.RoleID] = append(rpMemo[rp.RoleID], pMemo[rp.PermissionID].Mark)
	}

	cMap := make(map[string][]string, 1)
	tMap := make(map[string][]string, len(rpMemo)-1)
	for _, role := range userRoles {
		if len(role.CompanyID) > 0 {
			for roleID, marks := range rpMemo {
				if role.RoleID == roleID {
					cMap[role.CompanyID] = marks
				}
			}
		}
		if len(role.TeamID) > 0 {
			for roleID, marks := range rpMemo {
				if role.RoleID == roleID {
					tMap[role.TeamID] = marks
				}
			}
		}
	}

	return &rao.UserAllPermissionMarksResp{
		Teams:   tMap,
		Company: cMap,
	}
}
