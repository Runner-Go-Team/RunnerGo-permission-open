package notice

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-resty/resty/v2"
	"permission-open/internal/pkg/biz/consts"
	"permission-open/internal/pkg/biz/log"
	"permission-open/internal/pkg/dal"
	"permission-open/internal/pkg/dal/rao"
	"strconv"
	"time"
)

const (
	baseURL           = "https://oapi.dingtalk.com"
	deptListSubAPI    = "/topapi/v2/department/listsub"
	userListSimpleAPI = "/topapi/user/listsimple"
)

type AccessToken struct {
	AccessToken string `json:"access_token"`
}

type Department struct {
	AutoAddUser     bool   `json:"auto_add_user"`
	CreateDeptGroup bool   `json:"create_dept_group"`
	DeptId          int64  `json:"dept_id"`
	Name            string `json:"name"`
	ParentId        int    `json:"parent_id"`
}

type DeptListSub struct {
	Errcode   int          `json:"errcode"`
	Errmsg    string       `json:"errmsg"`
	Result    []Department `json:"result"`
	RequestId string       `json:"request_id"`
}

type UserListSimpleResponse struct {
	ErrCode   int                  `json:"errcode"`
	ErrMsg    string               `json:"errmsg"`
	Result    UserListSimpleResult `json:"result"`
	RequestID string               `json:"request_id"`
}

type UserListSimpleResult struct {
	HasMore    bool   `json:"has_more"`
	NextCursor int64  `json:"next_cursor"`
	List       []User `json:"list"`
}

type User struct {
	Name   string `json:"name"`
	UserID string `json:"userid"`
}

func GetDingTalkUsers(ctx context.Context, appKey, appSecret string) (*rao.ThirdCompanyUsers, error) {
	// 先从缓存读取数据，缓存不存在，读取三方接口
	redisKey := consts.RedisNoticeDingTalkUsersPrefix + appKey
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

	// 获取 access_token
	accessToken, err := getAccessToken(appKey, appSecret)
	if err != nil {
		return nil, err
	}

	// 获取部门列表
	departments, err := GetDingTalkDepartments(ctx, accessToken, 1)
	if err != nil {
		return nil, err
	}
	log.Logger.Info("GetDingTalkUsers- departments :", departments)

	parentUsers := make([]rao.ThirdUserInfo, 0)
	// 获取部门成员详情
	simpleUser, err := getUserListSimple(accessToken, 1)
	if err != nil {
		return nil, err
	}
	for _, u := range simpleUser {
		thirdUser := rao.ThirdUserInfo{
			OpenID: u.UserID,
			Name:   u.Name,
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

func GetDingTalkDepartments(ctx context.Context, accessToken string, departmentId int64) ([]rao.ThirdDepartmentInfo, error) {
	departments := make([]rao.ThirdDepartmentInfo, 0)

	// 获取子部门列表
	departmentList, err := getDepartmentSubList(accessToken, departmentId)
	if err != nil {
		return nil, err
	}

	for _, department := range departmentList {
		// 递归获取子部门
		subDepartments, err := GetDingTalkDepartments(ctx, accessToken, department.DeptId)
		if err != nil {
			return nil, err
		}

		var subUsers = make([]rao.ThirdUserInfo, 0)
		// 获取部门成员详情
		subSimpleUser, err := getUserListSimple(accessToken, department.DeptId)
		if err != nil {
			return nil, err
		}
		for _, u := range subSimpleUser {
			thirdUser := rao.ThirdUserInfo{
				OpenID: u.UserID,
				Name:   u.Name,
			}
			subUsers = append(subUsers, thirdUser)
		}

		// 当前部门信息
		departmentInfo := rao.ThirdDepartmentInfo{
			DepartmentID:   strconv.FormatInt(department.DeptId, 10),
			Name:           department.Name,
			UserList:       subUsers,
			DepartmentList: subDepartments,
		}

		departments = append(departments, departmentInfo)
	}

	return departments, nil
}

type getAccessTokenResp struct {
	Errcode     int    `json:"errcode"`
	AccessToken string `json:"access_token"`
	Errmsg      string `json:"errmsg"`
	ExpiresIn   int    `json:"expires_in"`
}

func getAccessToken(appKey, appSecret string) (string, error) {
	client := resty.New()

	resp, err := client.R().
		SetQueryParams(map[string]string{
			"appkey":    appKey,
			"appsecret": appSecret,
		}).
		Get(baseURL + "/gettoken")

	if err != nil {
		return "", err
	}

	var getAccessTokenRes *getAccessTokenResp
	if err := json.Unmarshal(resp.Body(), &getAccessTokenRes); err != nil {
		return "", err
	}

	if getAccessTokenRes.Errcode != 0 {
		return "", errors.New(getAccessTokenRes.Errmsg)
	}

	return getAccessTokenRes.AccessToken, nil
}

type UserListSimpleParams struct {
	DeptID int64 `json:"dept_id"`
	Cursor int64 `json:"cursor"`
	Size   int64 `json:"size"`
}

func getUserListSimple(accessToken string, deptID int64) ([]User, error) {
	client := resty.New()

	users := make([]User, 0)
	for {
		cursor := int64(0)
		departmentInfoParams := UserListSimpleParams{
			DeptID: deptID,
			Cursor: cursor,
			Size:   100,
		}
		bodyByte, err := json.Marshal(departmentInfoParams)
		if err != nil {
			return nil, err
		}

		resp, err := client.R().
			SetQueryParams(map[string]string{
				"access_token": accessToken,
			}).
			SetBody(bodyByte).
			Post(baseURL + userListSimpleAPI)

		if err != nil {
			return nil, err
		}

		var userListSimpleResponse UserListSimpleResponse
		if err := json.Unmarshal(resp.Body(), &userListSimpleResponse); err != nil {
			return nil, err
		}
		if userListSimpleResponse.ErrCode != 0 {
			return nil, errors.New(userListSimpleResponse.ErrMsg)
		}
		users = append(users, userListSimpleResponse.Result.List...)

		if !userListSimpleResponse.Result.HasMore {
			break
		}
		cursor = userListSimpleResponse.Result.NextCursor
	}

	return users, nil
}

type getDepartmentSubListParams struct {
	DeptID int64 `json:"dept_id"`
}

func getDepartmentSubList(accessToken string, deptID int64) ([]Department, error) {
	client := resty.New()

	departmentSubListParams := getDepartmentSubListParams{
		DeptID: deptID,
	}
	bodyByte, err := json.Marshal(departmentSubListParams)
	if err != nil {
		return nil, err
	}

	resp, err := client.R().
		SetQueryParams(map[string]string{
			"access_token": accessToken,
		}).
		SetBody(bodyByte).
		Post(baseURL + deptListSubAPI)

	if err != nil {
		return nil, err
	}

	log.Logger.Info("GetDingTalkUsers -  getDepartmentSubList resp:", resp)
	var deptListSub DeptListSub
	if err := json.Unmarshal(resp.Body(), &deptListSub); err != nil {
		return nil, err
	}
	if deptListSub.Errcode != 0 {
		return nil, errors.New(deptListSub.Errmsg)
	}

	return deptListSub.Result, nil
}
