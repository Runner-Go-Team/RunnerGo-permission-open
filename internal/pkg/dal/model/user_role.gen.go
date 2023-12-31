// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"

	"gorm.io/gorm"
)

const TableNameUserRole = "user_role"

// UserRole mapped from table <user_role>
type UserRole struct {
	ID           int64          `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`                        // 主键id
	RoleID       string         `gorm:"column:role_id;not null" json:"role_id"`                                   // 角色id
	UserID       string         `gorm:"column:user_id;not null" json:"user_id"`                                   // 用户id
	CompanyID    string         `gorm:"column:company_id;not null" json:"company_id"`                             // 企业id
	TeamID       string         `gorm:"column:team_id;not null" json:"team_id"`                                   // 团队id
	InviteUserID string         `gorm:"column:invite_user_id;not null;default:0" json:"invite_user_id"`           // 邀请人id
	InviteTime   time.Time      `gorm:"column:invite_time;not null;default:CURRENT_TIMESTAMP" json:"invite_time"` // 邀请时间
	CreatedAt    time.Time      `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

// TableName UserRole's table name
func (*UserRole) TableName() string {
	return TableNameUserRole
}
