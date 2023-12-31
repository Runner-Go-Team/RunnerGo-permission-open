// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"

	"gorm.io/gorm"
)

const TableNamePreinstallConf = "preinstall_conf"

// PreinstallConf mapped from table <preinstall_conf>
type PreinstallConf struct {
	ID            int32          `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`                      // 主键id
	ConfName      string         `gorm:"column:conf_name;not null" json:"conf_name"`                             // 配置名称
	TeamID        string         `gorm:"column:team_id;not null" json:"team_id"`                                 // 团队ID
	UserID        string         `gorm:"column:user_id;not null;default:0" json:"user_id"`                       // 用户ID
	UserName      string         `gorm:"column:user_name;not null" json:"user_name"`                             // 用户名称
	TaskType      int32          `gorm:"column:task_type;not null" json:"task_type"`                             // 任务类型
	TaskMode      int32          `gorm:"column:task_mode;not null" json:"task_mode"`                             // 压测模式
	ControlMode   int32          `gorm:"column:control_mode;not null" json:"control_mode"`                       // 控制模式：0-集中模式，1-单独模式
	DebugMode     string         `gorm:"column:debug_mode;not null;default:stop" json:"debug_mode"`              // debug模式：stop-关闭，all-开启全部日志，only_success-开启仅成功日志，only_error-开启仅错误日志
	ModeConf      string         `gorm:"column:mode_conf;not null" json:"mode_conf"`                             // 压测配置详情
	TimedTaskConf string         `gorm:"column:timed_task_conf;not null" json:"timed_task_conf"`                 // 定时任务相关配置
	CreatedAt     time.Time      `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"created_at"` // 创建时间
	UpdatedAt     time.Time      `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP" json:"updated_at"` // 更新时间
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`                                    // 删除时间
}

// TableName PreinstallConf's table name
func (*PreinstallConf) TableName() string {
	return TableNamePreinstallConf
}
