package models

import (
	"time"

	"gorm.io/gorm"
)

// Role 角色模型
type Role struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	Code        string `gorm:"uniqueIndex;size:50;not null;comment:角色代码" json:"code"`
	Name        string `gorm:"size:100;not null" json:"name"`
	Description string `gorm:"size:500" json:"description"`
	Status      int    `gorm:"default:1;comment:0-禁用,1-启用" json:"status"`

	// 关联关系
	Users []User `gorm:"many2many:sys_user_role;" json:"users,omitempty"`
	Menus []Menu `gorm:"many2many:sys_role_menu;" json:"menus,omitempty"`
}

// TableName 指定表名
func (Role) TableName() string {
	return "sys_role"
}

