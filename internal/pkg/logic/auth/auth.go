package auth

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"math/rand"
	"permission-open/internal/pkg/conf"
	"permission-open/internal/pkg/dal/global"
	"permission-open/internal/pkg/logic/errmsg"
	"permission-open/internal/pkg/logic/permission"
	"permission-open/internal/pkg/logic/role"
	"permission-open/internal/pkg/packer"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-omnibus/omnibus"
	"permission-open/internal/pkg/biz/consts"
	"permission-open/internal/pkg/biz/uuid"
	"permission-open/internal/pkg/dal"
	"permission-open/internal/pkg/dal/model"
	"permission-open/internal/pkg/dal/query"
	"permission-open/internal/pkg/dal/rao"
	"permission-open/internal/pkg/logic/team"
)

// CompanyRegister 注册企业及超级管理员
func CompanyRegister(ctx context.Context, account, password string) (*model.User, error) {
	c := query.Use(dal.DB()).Company
	companyInfo, err := c.WithContext(ctx).First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// 查询当前是否已经有企业
	if companyInfo != nil {
		r := query.Use(dal.DB()).Role
		superRole, err := r.WithContext(ctx).Where(r.Level.Eq(consts.RoleTypeCompany),
			r.RoleType.Eq(consts.RoleLevelSuperManager)).First()
		if err != nil {
			return nil, err
		}
		ur := query.Use(dal.DB()).UserRole
		userRole, err := ur.WithContext(ctx).Where(ur.RoleID.Eq(superRole.RoleID)).First()
		if err != nil {
			return nil, err
		}

		global.CompanyID = companyInfo.CompanyID
		global.SuperManageRoleID = superRole.RoleID
		global.SuperManageUserID = userRole.UserID

		return nil, nil
	}

	hashedPassword, err := omnibus.GenerateBcryptFromPassword(password)
	if err != nil {
		return nil, err
	}

	rand.Seed(time.Now().UnixNano())
	user := model.User{
		UserID:   uuid.GetUUID(),
		Email:    "",
		Account:  account,
		Password: hashedPassword,
		Nickname: account,
		Avatar:   consts.DefaultAvatarMemo[rand.Intn(3)],
	}

	cSuperUUID := uuid.GetUUID()
	cManagerUUID := uuid.GetUUID()
	cGeneralUUID := uuid.GetUUID()
	tSuperUUID := uuid.GetUUID()
	tManagerUUID := uuid.GetUUID()

	err = query.Use(dal.DB()).Transaction(func(tx *query.Query) error {
		// step1: 生成用户
		if err = tx.User.WithContext(ctx).Create(&user); err != nil {
			return err
		}

		// step2: 生成企业
		company := model.Company{
			CompanyID: uuid.GetUUID(),
			Name:      conf.Conf.CompanyInitConfig.Name,
		}
		if err = tx.Company.WithContext(ctx).Create(&company); err != nil {
			return err
		}

		// step3: 维护企业用户
		if err = tx.UserCompany.WithContext(ctx).Create(&model.UserCompany{
			UserID:     user.UserID,
			CompanyID:  company.CompanyID,
			InviteTime: time.Now(),
		}); err != nil {
			return nil
		}

		// step4: 创建企业默认角色
		var roles = []*model.Role{
			{
				RoleID:    cSuperUUID,
				RoleType:  consts.RoleTypeCompany,
				Name:      "超管",
				CompanyID: company.CompanyID,
				Level:     consts.RoleLevelSuperManager,
				IsDefault: consts.RoleDefault,
			},
			{
				RoleID:    cManagerUUID,
				RoleType:  consts.RoleTypeCompany,
				Name:      "管理员",
				CompanyID: company.CompanyID,
				Level:     consts.RoleLevelManager,
				IsDefault: consts.RoleDefault,
			},
			{
				RoleID:    cGeneralUUID,
				RoleType:  consts.RoleTypeCompany,
				Name:      "普通成员",
				CompanyID: company.CompanyID,
				Level:     consts.RoleLevelGeneral,
				IsDefault: consts.RoleDefault,
			},
			{
				RoleID:    tSuperUUID,
				RoleType:  consts.RoleTypeTeam,
				Name:      "团队管理员",
				CompanyID: company.CompanyID,
				Level:     consts.RoleLevelSuperManager,
				IsDefault: consts.RoleDefault,
			},
			{
				RoleID:    tManagerUUID,
				RoleType:  consts.RoleTypeTeam,
				Name:      "团队成员",
				CompanyID: company.CompanyID,
				Level:     consts.RoleLevelManager,
				IsDefault: consts.RoleDefault,
			},
		}
		if err = tx.Role.WithContext(ctx).Create(roles...); err != nil {
			return err
		}

		// step5: 关联当前企业角色
		userRole := model.UserRole{
			RoleID:     cSuperUUID,
			UserID:     user.UserID,
			CompanyID:  company.CompanyID,
			InviteTime: time.Now(),
		}
		if err = tx.UserRole.WithContext(ctx).Create(&userRole); err != nil {
			return err
		}

		// step6: 同步角色权限信息
		rolePermissions := make([]*model.RolePermission, 0, 50)

		rtp := tx.RoleTypePermission
		companyRtpList, err := rtp.WithContext(ctx).Where(rtp.RoleType.Eq(consts.RoleTypeCompany)).Find()
		if err != nil {
			return err
		}

		teamRtpList, err := rtp.WithContext(ctx).Where(rtp.RoleType.Eq(consts.RoleTypeTeam)).Find()
		if err != nil {
			return err
		}

		companyPIDs := make([]int64, 0, len(companyRtpList))
		teamPIDs := make([]int64, 0, len(teamRtpList))

		for _, cp := range companyRtpList {
			companyPIDs = append(companyPIDs, cp.PermissionID)
		}

		for _, tp := range teamRtpList {
			teamPIDs = append(teamPIDs, tp.PermissionID)
		}

		rp := tx.RolePermission
		for _, r := range roles {
			if r.RoleType == consts.RoleTypeCompany && r.Level <= consts.RoleLevelManager {
				for _, cp := range companyPIDs {
					rolePermissions = append(rolePermissions, &model.RolePermission{
						RoleID:       r.RoleID,
						PermissionID: cp,
					})
				}
				_, err = rp.WithContext(ctx).Where(rp.RoleID.Eq(r.RoleID)).Delete()
				if err != nil {
					return err
				}
			}

			if r.RoleType == consts.RoleTypeTeam {
				// 团队管理员默认全部权限
				if r.Level == consts.RoleLevelSuperManager {
					for _, tp := range teamPIDs {
						rolePermissions = append(rolePermissions, &model.RolePermission{
							RoleID:       r.RoleID,
							PermissionID: tp,
						})
					}
					_, err = rp.WithContext(ctx).Where(rp.RoleID.Eq(r.RoleID)).Delete()
					if err != nil {
						return err
					}
				}
				// 团队成员默认指定权限
				if r.Level == consts.RoleLevelManager {
					permissionIDs := role.GetRoleTeamLevelManagerPermission()
					for _, pID := range permissionIDs {
						rolePermissions = append(rolePermissions, &model.RolePermission{
							RoleID:       r.RoleID,
							PermissionID: pID,
						})
					}
					_, err = rp.WithContext(ctx).Where(rp.RoleID.Eq(r.RoleID)).Delete()
					if err != nil {
						return err
					}
				}
			}
		}

		err = rp.WithContext(ctx).Create(rolePermissions...)
		if err != nil {
			return err
		}

		global.CompanyID = company.CompanyID
		global.SuperManageRoleID = cSuperUUID
		global.SuperManageUserID = user.UserID

		return nil
	})

	if err != nil {
		return nil, err
	}

	if err = permission.Clear(ctx, user.UserID, global.CompanyID, cManagerUUID, tSuperUUID, tManagerUUID); err != nil {
		return nil, err
	}

	return &user, nil
}

// UpdatePermission 更新企业权限
func UpdatePermission(ctx context.Context, companyID string) error {
	r := query.Use(dal.DB()).Role
	roles, err := r.WithContext(ctx).Where(r.CompanyID.Eq(companyID)).Find()
	if err != nil {
		return err
	}
	rolePermissions := make([]*model.RolePermission, 0, 50)

	rtp := query.Use(dal.DB()).RoleTypePermission
	companyRtpList, err := rtp.WithContext(ctx).Where(rtp.RoleType.Eq(consts.RoleTypeCompany)).Find()
	if err != nil {
		return err
	}

	teamRtpList, err := rtp.WithContext(ctx).Where(rtp.RoleType.Eq(consts.RoleTypeTeam)).Find()
	if err != nil {
		return err
	}

	companyPIDs := make([]int64, 0, len(companyRtpList))
	teamPIDs := make([]int64, 0, len(teamRtpList))

	for _, cp := range companyRtpList {
		companyPIDs = append(companyPIDs, cp.PermissionID)
	}

	for _, tp := range teamRtpList {
		teamPIDs = append(teamPIDs, tp.PermissionID)
	}

	rp := query.Use(dal.DB()).RolePermission
	for _, r := range roles {
		if r.RoleType == consts.RoleTypeCompany && r.Level <= consts.RoleLevelManager {
			for _, cp := range companyPIDs {
				rolePermissions = append(rolePermissions, &model.RolePermission{
					RoleID:       r.RoleID,
					PermissionID: cp,
				})
			}
			_, err = rp.WithContext(ctx).Where(rp.RoleID.Eq(r.RoleID)).Delete()
			if err != nil {
				return err
			}
		}

		if r.RoleType == consts.RoleTypeTeam {
			// 团队管理员默认全部权限
			if r.Level == consts.RoleLevelSuperManager {
				for _, tp := range teamPIDs {
					rolePermissions = append(rolePermissions, &model.RolePermission{
						RoleID:       r.RoleID,
						PermissionID: tp,
					})
				}
				_, err = rp.WithContext(ctx).Where(rp.RoleID.Eq(r.RoleID)).Delete()
				if err != nil {
					return err
				}
			}
			// 团队成员默认指定权限
			if r.Level == consts.RoleLevelManager {
				permissionIDs := role.GetRoleTeamLevelManagerPermission()
				for _, pID := range permissionIDs {
					rolePermissions = append(rolePermissions, &model.RolePermission{
						RoleID:       r.RoleID,
						PermissionID: pID,
					})
				}
				_, err = rp.WithContext(ctx).Where(rp.RoleID.Eq(r.RoleID)).Delete()
				if err != nil {
					return err
				}
			}
		}
	}

	err = rp.WithContext(ctx).Create(rolePermissions...)
	if err != nil {
		return err
	}

	return nil
}

func Login(ctx *gin.Context, req rao.AuthLoginReq) (*model.User, error) {
	tx := query.Use(dal.DB()).User
	user, err := tx.WithContext(ctx).Where(tx.Account.Eq(req.Account)).First()
	if err != nil {
		return nil, errmsg.ErrAccountNotFound
	}

	if err := omnibus.CompareBcryptHashAndPassword(user.Password, req.Password); err != nil {
		return nil, errmsg.ErrPasswordFailed
	}

	uc := query.Use(dal.DB()).UserCompany
	userCompany, err := uc.WithContext(ctx).Where(uc.UserID.Eq(user.UserID)).First()
	if userCompany.Status == consts.CompanyUserStatusDisable {
		return nil, errmsg.ErrUserDisable
	}

	if req.InviteVerifyCode != "" {
		err := team.InviteLogin(ctx, req.InviteVerifyCode, user.UserID)
		if err != nil {
			return nil, err
		}
	}
	return user, nil
}

func UpdateLoginTime(ctx context.Context, userID string) error {
	tx := query.Use(dal.DB()).User
	_, err := tx.WithContext(ctx).Where(tx.UserID.Eq(userID)).UpdateColumn(tx.LastLoginAt, time.Now())
	return err
}

// GetUserInfo 默认查询当前用户信息，如果传 targetUserID 查此用户
func GetUserInfo(ctx context.Context, userID string, targetUserID string) (*rao.GetUserInfoResp, error) {
	if len(targetUserID) > 0 {
		userID = targetUserID
	}
	u := dal.GetQuery().User
	user, err := u.WithContext(ctx).Where(u.UserID.Eq(userID)).First()
	if err != nil {
		return nil, err
	}

	uc := dal.GetQuery().UserCompany
	userCompany, err := uc.WithContext(ctx).Where(uc.UserID.Eq(userID)).First()
	if err != nil {
		return nil, err
	}

	ut := dal.GetQuery().UserTeam
	userTeams, err := ut.WithContext(ctx).Where(ut.UserID.Eq(userID)).Find()
	if err != nil {
		return nil, err
	}

	teamIDs := make([]string, 0, len(userTeams))
	for _, userTeam := range userTeams {
		teamIDs = append(teamIDs, userTeam.TeamID)
	}

	t := dal.GetQuery().Team
	teams, err := t.WithContext(ctx).Where(t.TeamID.In(teamIDs...)).Find()
	if err != nil {
		return nil, err
	}

	ur := dal.GetQuery().UserRole
	userRoles, err := ur.WithContext(ctx).Where(ur.UserID.Eq(userID)).Find()
	if err != nil {
		return nil, err
	}
	roleIDs := make([]string, 0, len(userTeams)+1)
	for _, role := range userRoles {
		roleIDs = append(roleIDs, role.RoleID)
	}

	r := dal.GetQuery().Role
	roles, err := r.WithContext(ctx).Where(r.RoleID.In(roleIDs...)).Find()
	if err != nil {
		return nil, err
	}

	s := dal.GetQuery().Setting
	userSetting, _ := s.WithContext(ctx).Where(s.UserID.Eq(userID)).First()

	return packer.TransUserInfoToRaoUser(user, userCompany, userTeams, teams, userRoles, roles, userSetting), nil
}

// ResetLoginUsers 需要重新登录用户
func ResetLoginUsers(ctx *gin.Context, userID string) error {
	_, err := dal.GetRDB().SAdd(ctx, consts.RedisResetLoginUsers, userID).Result()
	if err != nil {
		return err
	}

	return nil
}

// RemoveResetLoginUser 删除重新登录的用户
func RemoveResetLoginUser(ctx *gin.Context, userID string) error {
	if exists, _ := dal.GetRDB().SIsMember(ctx, consts.RedisResetLoginUsers, userID).Result(); exists {
		_, err := dal.GetRDB().SRem(ctx, consts.RedisResetLoginUsers, userID).Result()
		if err != nil {
			return err
		}
	}

	return nil
}
