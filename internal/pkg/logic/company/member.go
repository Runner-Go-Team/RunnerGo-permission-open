package company

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	gevent "github.com/gookit/event"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
	"io"
	"math/rand"
	"os"
	"permission-open/internal/pkg/biz/consts"
	"permission-open/internal/pkg/biz/jwt"
	"permission-open/internal/pkg/biz/uuid"
	"permission-open/internal/pkg/dal"
	"permission-open/internal/pkg/dal/model"
	"permission-open/internal/pkg/dal/query"
	"permission-open/internal/pkg/dal/rao"
	"permission-open/internal/pkg/event"
	"permission-open/internal/pkg/logic/errmsg"
	"permission-open/internal/pkg/packer"
	"permission-open/internal/pkg/public"
	"strings"
	"time"

	"gorm.io/gen"

	"github.com/go-omnibus/omnibus"
)

// SaveMember 添加企业成员
func SaveMember(ctx *gin.Context, userID string, req rao.CompanySaveMemberReq) (*model.User, error) {
	// 判断角色
	ur := dal.GetQuery().UserRole
	userRole, err := ur.WithContext(ctx).Where(ur.UserID.Eq(userID), ur.CompanyID.Eq(req.CompanyId)).First()
	if err != nil {
		return nil, errmsg.ErrRoleNotExists
	}
	r := dal.GetQuery().Role
	roleInfo, err := r.WithContext(ctx).Where(r.RoleID.Eq(req.RoleId)).First()
	if err != nil {
		return nil, errmsg.ErrRoleNotExists
	}

	userRoleInfo, err := r.WithContext(ctx).Where(r.RoleID.Eq(userRole.RoleID)).First()
	if err != nil {
		return nil, errmsg.ErrRoleNotExists
	}

	// 当前用户角色判断
	if userRoleInfo.Level > roleInfo.Level {
		return nil, errmsg.ErrRoleForbidden
	}
	if userRoleInfo.Level == consts.RoleLevelManager && roleInfo.Level == consts.RoleLevelManager {
		return nil, errmsg.ErrRoleForbidden
	}

	tx := dal.GetQuery().User
	cnt, err := tx.WithContext(ctx).Where(tx.Account.Eq(req.Account)).Count()
	if err != nil {
		return nil, err
	}
	if cnt > 0 {
		return nil, errmsg.ErrYetAccountRegister
	}

	cnt, err = tx.WithContext(ctx).Where(tx.Nickname.Eq(req.Nickname)).Count()
	if err != nil {
		return nil, err
	}
	if cnt > 0 {
		return nil, errmsg.ErrYetNicknameRegister
	}

	hashedPassword, err := omnibus.GenerateBcryptFromPassword(req.Password)
	if err != nil {
		return nil, err
	}
	timeNow := time.Now()
	rand.Seed(timeNow.UnixNano())
	user := model.User{
		UserID:   uuid.GetUUID(),
		Email:    "",
		Account:  req.Account,
		Password: hashedPassword,
		Nickname: req.Nickname,
		Avatar:   consts.DefaultAvatarMemo[rand.Intn(3)],
	}

	// 邮箱
	if len(req.Email) > 0 {
		if !public.IsEmailValid(req.Email) {
			return nil, errmsg.ErrYetEmailValid
		}
		cnt, err = tx.WithContext(ctx).Where(tx.Email.Eq(req.Email)).Count()
		if err != nil {
			return nil, err
		}
		if cnt > 0 {
			return nil, errmsg.ErrYetEmailRegister
		}
		user.Email = req.Email
	}

	// 团队是否重复
	if len(req.TeamDetail) > 0 {
		teamIDs := make([]string, 0, len(req.TeamDetail))
		uniqueTeamIDs := make([]string, 0, len(req.TeamDetail))
		for _, t := range req.TeamDetail {
			teamIDs = append(teamIDs, t.TeamId)
		}
		uniqueTeamIDs = public.SliceUnique(teamIDs)
		if len(teamIDs) != len(uniqueTeamIDs) {
			return nil, errmsg.ErrTeamSaveRepeat
		}
	}

	err = query.Use(dal.DB()).Transaction(func(tx *query.Query) error {
		// step1: 生成用户
		if err = tx.User.WithContext(ctx).Create(&user); err != nil {
			return err
		}

		// step2: 维护企业用户
		userCompany := model.UserCompany{
			UserID:       user.UserID,
			CompanyID:    req.CompanyId,
			InviteUserID: userID,
			InviteTime:   timeNow,
		}
		if err = tx.UserCompany.WithContext(ctx).Create(&userCompany); err != nil {
			return nil
		}

		// step3: 维护企业角色
		userRoles := make([]*model.UserRole, 0)
		userCompanyRole := &model.UserRole{
			RoleID:       req.RoleId,
			UserID:       user.UserID,
			CompanyID:    req.CompanyId,
			InviteUserID: userID,
			InviteTime:   timeNow,
		}
		userRoles = append(userRoles, userCompanyRole)

		t := tx.Team
		// step4: 添加团队   // step5: 维护团队角色
		if len(req.TeamDetail) > 0 {
			userTeams := make([]*model.UserTeam, 0, len(req.TeamDetail))
			for _, td := range req.TeamDetail {
				teamInfo, err := t.WithContext(ctx).Where(t.TeamID.Eq(td.TeamId)).First()
				if err != nil || teamInfo.Type == consts.TeamTypePrivate {
					continue
				}

				userTeam := &model.UserTeam{
					UserID:       user.UserID,
					TeamID:       td.TeamId,
					InviteUserID: userID,
					InviteTime:   timeNow,
				}
				userTeams = append(userTeams, userTeam)

				userTeamRole := &model.UserRole{
					RoleID:       td.RoleId,
					UserID:       user.UserID,
					TeamID:       td.TeamId,
					InviteUserID: userID,
					InviteTime:   timeNow,
				}
				userRoles = append(userRoles, userTeamRole)

			}
			if err = tx.UserTeam.WithContext(ctx).Create(userTeams...); err != nil {
				return err
			}

			// step6: 设置默认团队
			if err = tx.Setting.WithContext(ctx).Create(&model.Setting{
				UserID: user.UserID,
				TeamID: req.TeamDetail[0].TeamId,
			}); err != nil {
				return err
			}
		}

		if err = tx.UserRole.WithContext(ctx).Create(userRoles...); err != nil {
			return nil
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// 触发事件
	if req.Source == "SaveMember" {
		if err, _ = gevent.Trigger(consts.ActionCompanySaveMember, event.BaseParams(ctx, "", userID, req.Nickname)); err != nil {
			return nil, err
		}
	}

	return &user, nil
}

// ImportMembers 导入成员
func ImportMembers(ctx *gin.Context, file io.Reader, userID string, companyID string) (*rao.ImportDesc, error) {
	f, err := excelize.OpenReader(file)
	if err != nil {
		return nil, err
	}

	// 获取 Sheet1 上所有单元格
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		return nil, err
	}

	if len(rows) > 1000 {
		return nil, errmsg.ErrFileMaxLimit
	}

	// 获取企业 role_id
	r := dal.GetQuery().Role
	roles, err := r.WithContext(ctx).Where(r.CompanyID.Eq(companyID), r.RoleType.Eq(consts.RoleTypeCompany)).Find()
	if err != nil {
		return nil, err
	}
	roleMemo := make(map[string]string)
	for _, role := range roles {
		roleMemo[role.Name] = role.RoleID
	}

	var (
		importMembers = make([]rao.ImportErr, 0, len(rows))
		successNum    int64
		failNum       int64
	)

	u := dal.GetQuery().User
	for i, row := range rows {
		if i == 0 || i == 1 || i == 2 {
			continue
		}
		saveMember := rao.CompanySaveMemberReq{CompanyId: companyID}
		importErr := rao.ImportErr{}
		//3.昵称：必填，长度为2-30位，由字母/数字/字符/中文组成，在本企业内不能重复。
		//4.账号：必填，必填，长度为6-20位，由字母/数字/字符组成，在本企业内不能重复。
		//5.密码：必填，长度为6-20位，由字母/数字/字符组成。
		//6.企业角色：必填，请输入企业角色名称。
		for ii, colCell := range row {
			if ii == 0 {
				saveMember.Nickname = colCell
				importErr.Nickname = colCell
				if len(colCell) < 2 || len(colCell) > 30 {
					importErr.ErrMsg += "昵称长度有误 "
				}
			}
			if ii == 1 {
				saveMember.Account = colCell
				importErr.Account = colCell
				if len(colCell) < 6 || len(colCell) > 20 {
					importErr.ErrMsg += "账号长度有误 "
				}
				cnt, _ := u.WithContext(ctx).Where(u.Account.Eq(colCell)).Count()
				if cnt > 0 {
					importErr.ErrMsg = "账号不能重复 "
				}
			}
			if ii == 2 {
				saveMember.Password = colCell
				importErr.Password = colCell
				if len(colCell) < 6 || len(colCell) > 20 {
					importErr.ErrMsg += "密码长度有误 "
				}
			}
			if ii == 3 {
				importErr.RoleName = colCell
				if role, ok := roleMemo[colCell]; ok {
					saveMember.RoleId = role
				} else {
					importErr.ErrMsg += "企业角色名称有误 "
				}
			}
			if ii == 4 {
				if len(colCell) > 0 {
					saveMember.Email = colCell
					importErr.Email = colCell
					if !public.IsEmailValid(colCell) {
						importErr.ErrMsg += "邮箱格式不正确  "
					}
					cnt, _ := u.WithContext(ctx).Where(u.Email.Eq(colCell)).Count()
					if cnt > 0 {
						importErr.ErrMsg = "邮箱不能重复 "
					}
				}
			}
		}
		if importErr.ErrMsg == "" {
			if _, err = SaveMember(ctx, userID, saveMember); err != nil {
				importErr.ErrMsg += err.Error()
				failNum++
			} else {
				successNum++
			}
		} else {
			failNum++
		}

		// 有错误才返回
		if importErr.ErrMsg != "" {
			importMembers = append(importMembers, importErr)
		}

		// 删除 row
		_ = f.RemoveRow("Sheet1", i+1)
	}

	var path string
	if len(importMembers) > 0 {
		err = f.SetCellValue("Sheet1", "F2", "错误信息")

		for key, importErr := range importMembers {
			cell := fmt.Sprintf("A%d", key+4)
			err = f.SetSheetRow("Sheet1", cell, &[]interface{}{
				importErr.Nickname,
				importErr.Account,
				importErr.Password,
				importErr.RoleName,
				importErr.Email,
				importErr.ErrMsg})
			if err != nil {
				return nil, err
			}
		}

		filenameWithSuffix := "RunnerGo创建成员批量导入模板 - 错误报告.xlsx"
		StaticPath := "./static/"
		// 生成新的 Excel
		filepath := StaticPath + userID + "/"

		//创建文件夹
		if !public.CheckFileIsExist(StaticPath) {
			err = os.Mkdir(StaticPath, os.ModePerm)
			if err != nil {
				return nil, err
			}
		}
		if !public.CheckFileIsExist(filepath) {
			err = os.Mkdir(filepath, os.ModePerm)
			if err != nil {
				return nil, err
			}
		}

		if err = f.SaveAs(filepath + filenameWithSuffix); err != nil {
			return nil, err
		}

		path = filepath[1:] + filenameWithSuffix
	}

	if successNum > 0 {
		if err, _ = gevent.Trigger(consts.ActionCompanyExportMember, event.BaseParams(ctx, "", userID, fmt.Sprintf("成功：%d 位成员", successNum))); err != nil {
			return nil, err
		}
	}

	return &rao.ImportDesc{
		SuccessNum:    successNum,
		FailNum:       failNum,
		ImportErrDesc: importMembers,
		Path:          path,
	}, nil
}

// RemoveMember 移除企业成员
func RemoveMember(ctx context.Context, userID string, companyID string, targetUserID string) error {
	u := query.Use(dal.DB()).User
	targetUser, err := u.WithContext(ctx).Where(u.UserID.Eq(targetUserID)).First()
	if err != nil {
		return nil
	}

	r := dal.GetQuery().Role
	roles, err := r.WithContext(ctx).Where(r.CompanyID.Eq(companyID)).Find()
	if err != nil {
		return err
	}

	var (
		companySuperRoleID string
		teamSuperRoleID    string
	)
	for _, role := range roles {
		if role.Level == consts.RoleLevelSuperManager {
			if role.RoleType == consts.RoleTypeCompany {
				companySuperRoleID = role.RoleID
			}
			if role.RoleType == consts.RoleTypeTeam {
				teamSuperRoleID = role.RoleID
			}
		}
	}

	ur := dal.GetQuery().UserRole
	t := dal.GetQuery().Team
	userRoles, err := ur.WithContext(ctx).Where(ur.UserID.Eq(targetUserID)).Find()
	for _, role := range userRoles {
		if role.RoleID == companySuperRoleID {
			return errmsg.ErrCompanyNotRemoveMember
		}
		if len(role.TeamID) > 0 {
			_, err = t.WithContext(ctx).Where(t.TeamID.Eq(role.TeamID)).First()
			if errors.Is(err, gorm.ErrRecordNotFound) {
				continue
			}
			if err != nil {
				return err
			}
			if role.RoleID == teamSuperRoleID {
				return errmsg.ErrTeamNotRemoveCreate
			}
		}
	}

	err = query.Use(dal.DB()).Transaction(func(tx *query.Query) error {
		// step1: 删除用户
		if _, err = tx.User.WithContext(ctx).Where(tx.User.UserID.Eq(targetUserID)).Delete(); err != nil {
			return err
		}

		// step2: 删除企业用户关系
		if _, err = tx.UserCompany.WithContext(ctx).Where(tx.UserCompany.UserID.Eq(targetUserID),
			tx.UserCompany.CompanyID.Eq(companyID)).Delete(); err != nil {
			return err
		}

		// step3: 删除用户角色关系
		if _, err = tx.UserRole.WithContext(ctx).Where(tx.UserRole.UserID.Eq(targetUserID)).Delete(); err != nil {
			return err
		}

		// step4: 用户团队关系
		if _, err = tx.UserTeam.WithContext(ctx).Where(tx.UserTeam.UserID.Eq(targetUserID)).Delete(); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	// 触发事件
	if err, _ = gevent.Trigger(consts.ActionCompanyRemoveMember, event.BaseParams(ctx, "", userID, targetUser.Nickname)); err != nil {
		return err
	}

	return nil
}

// MemberList 企业成员列表
func MemberList(ctx *gin.Context, companyID string, keyword string, page int, size int) ([]*rao.CompanyMember, int64, error) {
	// step1: user_company
	// step2: user
	// step3: role
	// step4: user_role
	// step5: user_team
	// step6: team
	curUserID := jwt.GetUserIDByCtx(ctx)

	limit := size
	offset := (page - 1) * size

	// keyword 搜索昵称/账号
	u := dal.GetQuery().User
	conditions := make([]gen.Condition, 0)
	conditionsAccount := conditions
	keyword = strings.TrimSpace(keyword)
	if len(keyword) > 0 {
		conditions = append(conditions, u.Nickname.Like(fmt.Sprintf("%%%s%%", keyword)))
		conditionsAccount = append(conditionsAccount, u.Account.Like(fmt.Sprintf("%%%s%%", keyword)))
	}
	users, total, err := u.WithContext(ctx).Where(conditions...).Or(conditionsAccount...).Order(u.ID.Desc()).FindByPage(offset, limit)
	if err != nil {
		return nil, 0, err
	}

	userIDs := make([]string, 0, len(users))
	for _, u := range users {
		userIDs = append(userIDs, u.UserID)
	}

	uc := query.Use(dal.DB()).UserCompany
	userCompaniesFilter, err := uc.WithContext(ctx).Where(uc.CompanyID.Eq(companyID)).Where(uc.UserID.In(userIDs...)).Order(uc.ID.Desc()).Find()
	if err != nil {
		return nil, 0, err
	}

	// 查询当前企业下的角色
	r := dal.GetQuery().Role
	roles, err := r.WithContext(ctx).Where(r.CompanyID.Eq(companyID)).Find()
	if err != nil {
		return nil, 0, err
	}

	filterUserIDs := make([]string, 0, len(userCompaniesFilter))
	for _, company := range userCompaniesFilter {
		filterUserIDs = append(filterUserIDs, company.UserID)
		filterUserIDs = append(filterUserIDs, company.InviteUserID)
	}

	// 用户对应的角色
	ur := dal.GetQuery().UserRole
	userRoles, err := ur.WithContext(ctx).Where(ur.CompanyID.Eq(companyID)).Find()
	if err != nil {
		return nil, 0, err
	}

	curUserRole, err := ur.WithContext(ctx).Where(ur.UserID.Eq(curUserID), ur.CompanyID.Eq(companyID)).First()
	if err != nil {
		return nil, 0, err
	}

	userAll, err := u.WithContext(ctx).Where(u.UserID.In(filterUserIDs...)).Find()
	if err != nil {
		return nil, 0, err
	}

	return packer.TransCompaniesModelToRaoCompanyMember(userCompaniesFilter, userAll, roles, userRoles, curUserRole.RoleID), total, nil
}

func UpdateMember(ctx context.Context, userID string, req rao.CompanyUpdateMemberReq) error {
	u := query.Use(dal.DB()).User
	targetUser, err := u.WithContext(ctx).Where(u.UserID.Eq(req.TargetUserID)).First()
	if err != nil {
		return nil
	}

	err = query.Use(dal.DB()).Transaction(func(tx *query.Query) error {
		if req.Status > 0 {
			// step1: 修改企业用户状态
			_, err := tx.UserCompany.WithContext(ctx).Where(tx.UserCompany.UserID.Eq(req.TargetUserID),
				tx.UserCompany.CompanyID.Eq(req.CompanyId)).Update(tx.UserCompany.Status, req.Status)
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
	if err, _ = gevent.Trigger(consts.ActionCompanyUpdateMember, event.BaseParams(ctx, "", userID, targetUser.Nickname)); err != nil {
		return err
	}

	return nil
}
