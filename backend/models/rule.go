package models

import (
	"time"

	"gorm.io/gorm"
)

// Rule 规则模型
type Rule struct {
	ID        uint           `gorm:"primarykey" json:"key"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Name     string `gorm:"size:200;not null" json:"name"`
	Desc     string `gorm:"size:500" json:"desc"`
	Status   int    `gorm:"default:1;comment:0-关闭,1-运行中,2-已上线,3-异常" json:"status"`
	CallNo   int    `gorm:"default:0;comment:服务调用次数" json:"callNo"`
	Owner    string `gorm:"size:100" json:"owner"`
	Avatar   string `gorm:"size:500" json:"avatar"`
	Href     string `gorm:"size:500" json:"href"`
	Disabled bool   `gorm:"default:0" json:"disabled"`
	Progress int    `gorm:"default:0" json:"progress"`
}

// TableName 指定表名
func (Rule) TableName() string {
	return "sys_rule"
}

