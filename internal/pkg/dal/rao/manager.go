package rao

type SendManagerApiResp struct {
	Code int         `json:"code"`
	Em   string      `json:"em"`
	Et   string      `json:"et"`
	Data interface{} `json:"data"`
}

type TressPlanListResp struct {
	TeamID string `json:"team_id" binding:"required"`
	Page   int    `json:"page" binding:"required"`
	Size   int    `json:"size" binding:"required"`
}

type AutoPlanListResp struct {
	TeamID string `json:"team_id" binding:"required"`
	Page   int    `json:"page" binding:"required"`
	Size   int    `json:"size" binding:"required"`
}
