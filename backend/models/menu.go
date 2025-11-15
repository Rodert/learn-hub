package models

import (
	"time"

	"gorm.io/gorm"
)

// Menu 菜单模型
type Menu struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	ParentID   uint   `gorm:"default:0;comment:父菜单ID" json:"parentId"`
	Name       string `gorm:"size:100;not null;comment:菜单名称(国际化key)" json:"name"`
	Path       string `gorm:"size:200;comment:路由路径" json:"path"`
	Component  string `gorm:"size:200;comment:组件路径" json:"component"`
	Icon       string `gorm:"size:50;comment:图标名称" json:"icon"`
	SortOrder  int    `gorm:"default:0;comment:排序顺序" json:"sortOrder"`
	Access     string `gorm:"size:100;comment:权限标识" json:"access"`
	Redirect   string `gorm:"size:200;comment:重定向路径" json:"redirect"`
	Layout     int    `gorm:"default:1;comment:是否显示布局" json:"layout"`
	Hidden     int    `gorm:"default:0;comment:是否隐藏" json:"hidden"`
	Status     int    `gorm:"default:1;comment:0-禁用,1-启用" json:"status"`

	// 关联关系
	Roles []Role `gorm:"many2many:sys_role_menu;" json:"roles,omitempty"`
}

// MenuVO 菜单视图对象（用于返回给前端）
type MenuVO struct {
	Path      string   `json:"path,omitempty"`
	Name      string   `json:"name,omitempty"`
	Icon      string   `json:"icon,omitempty"`
	Component string   `json:"component,omitempty"`
	Access    string   `json:"access,omitempty"`
	Redirect  string   `json:"redirect,omitempty"`
	Routes    []MenuVO `json:"routes,omitempty"`
}

// TableName 指定表名
func (Menu) TableName() string {
	return "sys_menu"
}

