package models

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Username  string `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Password  string `gorm:"size:255;not null" json:"-"`
	Name      string `gorm:"size:100" json:"name"`
	Email     string `gorm:"size:100" json:"email"`
	Phone     string `gorm:"size:20" json:"phone"`
	Avatar    string `gorm:"size:500" json:"avatar"`
	UserID    string `gorm:"uniqueIndex;size:50" json:"userid"`
	Status    int    `gorm:"default:1;comment:0-禁用,1-启用" json:"status"`
	Access    string `gorm:"size:50;comment:访问级别" json:"access"`
	Signature string `gorm:"size:500" json:"signature"`
	Title     string `gorm:"size:100" json:"title"`
	GroupID   uint   `gorm:"comment:所属组织ID" json:"groupId"`
	Country   string `gorm:"size:50" json:"country"`
	Province  string `gorm:"size:50" json:"province"`
	City      string `gorm:"size:50" json:"city"`
	Address   string `gorm:"size:500" json:"address"`

	// 关联关系
	Roles []Role `gorm:"many2many:sys_user_role;" json:"roles,omitempty"`
}

// TableName 指定表名
func (User) TableName() string {
	return "sys_user"
}

