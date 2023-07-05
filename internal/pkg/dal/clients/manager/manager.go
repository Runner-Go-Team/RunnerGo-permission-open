package manager

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"permission-open/internal/pkg/biz/log"
	"permission-open/internal/pkg/conf"
	"permission-open/internal/pkg/dal/rao"
)

var (
	// NewestStressPlanList 获取团队最新性能计划列表
	NewestStressPlanList = "/management/api/company/get_newest_stress_plan_list"

	// NewestAutoPlanList 获取团队最新自动化计划列表
	NewestAutoPlanList = "/management/api/company/get_newest_auto_plan_list"
)

func GetNewestStressPlanList(body *rao.TressPlanListResp) (interface{}, error) {
	bodyByte, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	log.Logger.Info("body:", string(bodyByte))

	url := conf.Conf.Clients.Manager.Domain + NewestStressPlanList

	response, err := resty.New().R().
		SetHeader("Content-Type", "application/json").
		SetBody(bodyByte).
		Post(url)
	if err != nil {
		return nil, err
	}

	resp := rao.SendManagerApiResp{}
	err = json.Unmarshal(response.Body(), &resp)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

func GetNewestAutoPlanList(body *rao.AutoPlanListResp) (interface{}, error) {
	bodyByte, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	log.Logger.Info("body:", string(bodyByte))

	url := conf.Conf.Clients.Manager.Domain + NewestAutoPlanList

	response, err := resty.New().R().
		SetHeader("Content-Type", "application/json").
		SetBody(bodyByte).
		Post(url)
	if err != nil {
		return nil, err
	}

	resp := rao.SendManagerApiResp{}
	err = json.Unmarshal(response.Body(), &resp)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}
