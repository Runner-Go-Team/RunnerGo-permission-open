package notice

import (
	"context"
	"encoding/json"
	"permission-open/internal/pkg/biz/consts"
	"permission-open/internal/pkg/biz/log"
	"permission-open/internal/pkg/dal"
	"permission-open/internal/pkg/dal/rao"
	"time"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcontact "github.com/larksuite/oapi-sdk-go/v3/service/contact/v3"
)

// GetFerShuUsers 获取飞书组织架构（官方未支持全部）
func GetFerShuUsers(ctx context.Context, appId, appSecret string) (*rao.ThirdCompanyUsers, error) {
	// 先从缓存读取数据，缓存不存在，读取三方接口
	redisKey := consts.RedisNoticeFerShuUsersPrefix + appId
	ret, err := dal.GetRDB().Get(ctx, redisKey).Result()
	if err == nil { // 从 Redis 获取数据
		// 将 JSON 字符串解码为切片结构体数据
		var result *rao.ThirdCompanyUsers
		err = json.Unmarshal([]byte(ret), &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	}

	// 创建 Client
	client := lark.NewClient(appId, appSecret)
	listScopeRespData, err := ListScopeRespData(ctx, client)
	if err != nil {
		return nil, err
	}

	departments, err := GetFerShuDepartments(ctx, appId, appSecret, "0")

	parentUsers := make([]rao.ThirdUserInfo, 0, len(listScopeRespData.UserIds))
	for _, userID := range listScopeRespData.UserIds {
		u, _ := GetUserRespData(ctx, client, userID)
		thirdUser := rao.ThirdUserInfo{
			OpenID: *u.User.OpenId,
			Name:   *u.User.Name,
			Avatar: *u.User.Avatar.AvatarOrigin,
		}
		parentUsers = append(parentUsers, thirdUser)
	}

	thirdCompanyUsers := &rao.ThirdCompanyUsers{
		DepartmentList: departments,
		UserList:       parentUsers,
	}

	marshal, err := json.Marshal(thirdCompanyUsers)
	if err == nil {
		if err := dal.GetRDB().Set(ctx, redisKey, marshal, time.Second*3600).Err(); err != nil {
			return nil, err
		}
	}

	return thirdCompanyUsers, nil
}

func ListScopeRespData(ctx context.Context, client *lark.Client) (*larkcontact.ListScopeRespData, error) {
	ret := &larkcontact.ListScopeRespData{}

	// 创建请求对象
	req := larkcontact.NewListScopeReqBuilder().
		PageSize(100).
		Build()
	// 发起请求
	resp, err := client.Contact.Scope.List(ctx, req)
	if err != nil {
		log.Logger.Info("FeiShuGetScope- Scope.List -err:", err)
		return ret, err
	}

	// 服务端错误处理
	if !resp.Success() {
		log.Logger.Info("FeiShuGetScope- resp.Success - err:", err)
		return ret, err
	}

	return resp.Data, nil
}

func GetFerShuDepartments(ctx context.Context, appId, appSecret string, departmentId string) ([]rao.ThirdDepartmentInfo, error) {
	// 创建 Client
	client := lark.NewClient(appId, appSecret)
	departments := make([]rao.ThirdDepartmentInfo, 0)

	// 获取子部门列表
	subDepartmentList, err := GetSonDepartmentRespData(ctx, client, departmentId)
	if err != nil {
		return nil, err
	}
	for _, subDepartment := range subDepartmentList.Items {
		// 递归获取子部门
		subDepartments, err := GetFerShuDepartments(ctx, appId, appSecret, *subDepartment.OpenDepartmentId)
		if err != nil {
			return nil, err
		}

		// 获取当前部门下用户数据
		subDepartmentUsers := make([]rao.ThirdUserInfo, 0)
		subUsers, err := GetFindByDepartmentUser(ctx, client, *subDepartment.OpenDepartmentId)
		if err != nil {
			return nil, err
		}
		for _, u := range subUsers.Items {
			thirdUser := rao.ThirdUserInfo{
				OpenID: *u.OpenId,
				Name:   *u.Name,
				Avatar: *u.Avatar.AvatarOrigin,
			}
			subDepartmentUsers = append(subDepartmentUsers, thirdUser)
		}

		// 当前部门信息
		departmentInfo := rao.ThirdDepartmentInfo{
			DepartmentID:   *subDepartment.OpenDepartmentId,
			Name:           *subDepartment.Name,
			MemberCount:    *subDepartment.MemberCount,
			UserList:       subDepartmentUsers,
			DepartmentList: subDepartments,
		}
		departments = append(departments, departmentInfo)
	}

	return departments, nil
}

func GetDepartmentRespData(ctx context.Context, client *lark.Client, departmentId string) (*larkcontact.GetDepartmentRespData, error) {
	ret := &larkcontact.GetDepartmentRespData{}

	req := larkcontact.NewGetDepartmentReqBuilder().
		DepartmentId(departmentId).
		Build()

	// 发起请求
	resp, err := client.Contact.Department.Get(ctx, req)
	if err != nil {
		log.Logger.Info("GetFeiShuDepartment- Department.Get -err:", err)
		return ret, err
	}

	// 服务端错误处理
	if !resp.Success() {
		log.Logger.Info("GetFeiShuDepartment- resp.Success - err:", err)
		return ret, err
	}

	return resp.Data, nil
}

// GetSonDepartmentRespData 获取子部门列表
func GetSonDepartmentRespData(ctx context.Context, client *lark.Client, departmentId string) (*larkcontact.ChildrenDepartmentRespData, error) {
	ret := &larkcontact.ChildrenDepartmentRespData{}

	// 创建请求对象
	req := larkcontact.NewChildrenDepartmentReqBuilder().
		DepartmentId(departmentId).
		PageSize(50).
		Build()
	// 发起请求
	resp, err := client.Contact.Department.Children(ctx, req)
	if err != nil {
		log.Logger.Info("GetSonDepartmentRespData- Department.Get -err:", err)
		return ret, err
	}

	// 服务端错误处理
	if !resp.Success() {
		log.Logger.Info("GetSonDepartmentRespData- resp.Success - err:", err)
		return ret, err
	}

	return resp.Data, nil
}

// GetFindByDepartmentUser 获取部门直属用户列表
func GetFindByDepartmentUser(ctx context.Context, client *lark.Client, departmentId string) (*larkcontact.FindByDepartmentUserRespData, error) {
	ret := &larkcontact.FindByDepartmentUserRespData{}

	// 创建请求对象
	req := larkcontact.NewFindByDepartmentUserReqBuilder().
		UserIdType(`open_id`).
		DepartmentIdType(`open_department_id`).
		DepartmentId(departmentId).
		PageSize(50).
		Build()
	// 发起请求
	resp, err := client.Contact.User.FindByDepartment(ctx, req)

	if err != nil {
		log.Logger.Info("GetFindByDepartmentUser- Department.Get -err:", err)
		return ret, err
	}

	// 服务端错误处理
	if !resp.Success() {
		log.Logger.Info("GetFindByDepartmentUser- resp.Success - err:", err)
		return ret, err
	}

	return resp.Data, nil
}

func GetUserRespData(ctx context.Context, client *lark.Client, userID string) (*larkcontact.GetUserRespData, error) {
	ret := &larkcontact.GetUserRespData{}

	// 创建请求对象
	req := larkcontact.NewGetUserReqBuilder().
		UserId(userID).
		Build()
	// 发起请求
	resp, err := client.Contact.User.Get(ctx, req)

	if err != nil {
		log.Logger.Info("GetUserRespData- User.Get - err:", err)
		return ret, err
	}

	// 服务端错误处理
	if !resp.Success() {
		log.Logger.Info("GetUserRespData- resp.Success - err:", err)
		return ret, err
	}

	return resp.Data, nil
}
