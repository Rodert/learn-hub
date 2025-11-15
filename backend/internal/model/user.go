package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// User 用户表
type User struct {
	ID        int64          `gorm:"primaryKey" json:"id"`
	Username  string         `gorm:"uniqueIndex;size:100" json:"username"`
	Password  string         `gorm:"size:255" json:"-"`
	Nickname  string         `gorm:"size:100" json:"nickname"`
	OpenID    string         `gorm:"size:255" json:"openid"`
	Status    string         `gorm:"type:enum('active','inactive','banned');default:'active'" json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联
	Roles []Role `gorm:"many2many:user_roles;" json:"roles,omitempty"`
}

// Role 角色表
type Role struct {
	ID          int64          `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"uniqueIndex;size:100" json:"name"`
	Description string         `gorm:"size:255" json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
	Menus       []Menu       `gorm:"many2many:role_menus;" json:"menus,omitempty"`
}

// UserRole 用户角色关联表
type UserRole struct {
	ID        int64          `gorm:"primaryKey" json:"id"`
	UserID    int64          `gorm:"index" json:"user_id"`
	RoleID    int64          `gorm:"index" json:"role_id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Permission 权限表
type Permission struct {
	ID          int64          `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"uniqueIndex;size:100" json:"name"`
	Description string         `gorm:"size:255" json:"description"`
	Resource    string         `gorm:"size:100" json:"resource"`
	Action      string         `gorm:"size:50" json:"action"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// RolePermission 角色权限关联表
type RolePermission struct {
	ID           int64          `gorm:"primaryKey" json:"id"`
	RoleID       int64          `gorm:"index" json:"role_id"`
	PermissionID int64          `gorm:"index" json:"permission_id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// Menu 菜单表
type Menu struct {
	ID         int64          `gorm:"primaryKey" json:"id"`
	Name       string         `gorm:"size:100" json:"name"`
	Path       string         `gorm:"size:255" json:"path"`
	Icon       string         `gorm:"size:100" json:"icon"`
	Component  string         `gorm:"size:255" json:"component"`
	ParentID   *int64         `gorm:"index" json:"parent_id"`
	OrderNum   int            `gorm:"default:0" json:"order_num"`
	Visible    int            `gorm:"default:1" json:"visible"`
	Type       string         `gorm:"type:enum('menu','button');default:'menu'" json:"type"`
	Permission string         `gorm:"size:100" json:"permission"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联
	Children []Menu `gorm:"foreignKey:ParentID;constraint:OnDelete:CASCADE" json:"children,omitempty"`
}

// RoleMenu 角色菜单关联表
type RoleMenu struct {
	ID        int64          `gorm:"primaryKey" json:"id"`
	RoleID    int64          `gorm:"index" json:"role_id"`
	MenuID    int64          `gorm:"index" json:"menu_id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// JSONMap 自定义 JSON 类型
type JSONMap map[string]interface{}

func (j JSONMap) Value() (driver.Value, error) {
	return json.Marshal(j)
}

func (j *JSONMap) Scan(value interface{}) error {
	bytes := value.([]byte)
	return json.Unmarshal(bytes, &j)
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

func (Role) TableName() string {
	return "roles"
}

func (UserRole) TableName() string {
	return "user_roles"
}

func (Permission) TableName() string {
	return "permissions"
}

func (RolePermission) TableName() string {
	return "role_permissions"
}

func (Menu) TableName() string {
	return "menus"
}

func (RoleMenu) TableName() string {
	return "role_menus"
}
