package service

import (
	"gorm.io/gorm"
	"learn-hub/internal/model"
)

// PermissionService 权限服务
type PermissionService struct {
	db *gorm.DB
}

// NewPermissionService 创建权限服务
func NewPermissionService(db *gorm.DB) *PermissionService {
	return &PermissionService{db: db}
}

// GetUserPermissions 获取用户的所有权限
func (s *PermissionService) GetUserPermissions(userID int64) ([]string, error) {
	var permissions []string
	err := s.db.
		Table("permissions p").
		Joins("JOIN role_permissions rp ON p.id = rp.permission_id").
		Joins("JOIN roles r ON rp.role_id = r.id").
		Joins("JOIN user_roles ur ON r.id = ur.role_id").
		Where("ur.user_id = ? AND ur.deleted_at IS NULL", userID).
		Distinct("p.name").
		Pluck("p.name", &permissions).Error

	return permissions, err
}

// GetRolePermissions 获取角色的所有权限
func (s *PermissionService) GetRolePermissions(roleID int64) ([]model.Permission, error) {
	var permissions []model.Permission
	err := s.db.
		Table("permissions p").
		Joins("JOIN role_permissions rp ON p.id = rp.permission_id").
		Where("rp.role_id = ? AND rp.deleted_at IS NULL", roleID).
		Find(&permissions).Error

	return permissions, err
}

// AssignPermissionToRole 为角色分配权限
func (s *PermissionService) AssignPermissionToRole(roleID, permissionID int64) error {
	return s.db.Create(&model.RolePermission{
		RoleID:       roleID,
		PermissionID: permissionID,
	}).Error
}

// RemovePermissionFromRole 移除角色的权限
func (s *PermissionService) RemovePermissionFromRole(roleID, permissionID int64) error {
	return s.db.Delete(&model.RolePermission{}, "role_id = ? AND permission_id = ?", roleID, permissionID).Error
}

// HasPermission 检查用户是否有某个权限
func (s *PermissionService) HasPermission(userID int64, permission string) (bool, error) {
	var count int64
	err := s.db.
		Table("permissions p").
		Joins("JOIN role_permissions rp ON p.id = rp.permission_id").
		Joins("JOIN roles r ON rp.role_id = r.id").
		Joins("JOIN user_roles ur ON r.id = ur.role_id").
		Where("ur.user_id = ? AND p.name = ? AND ur.deleted_at IS NULL", userID, permission).
		Count(&count).Error

	return count > 0, err
}

// GetUserRoles 获取用户的所有角色
func (s *PermissionService) GetUserRoles(userID int64) ([]model.Role, error) {
	var roles []model.Role
	err := s.db.
		Table("roles r").
		Joins("JOIN user_roles ur ON r.id = ur.role_id").
		Where("ur.user_id = ? AND ur.deleted_at IS NULL", userID).
		Find(&roles).Error

	return roles, err
}

// GetUserMenus 获取用户可访问的菜单
func (s *PermissionService) GetUserMenus(userID int64) ([]model.Menu, error) {
	var menus []model.Menu
	err := s.db.
		Table("menus m").
		Joins("JOIN role_menus rm ON m.id = rm.menu_id").
		Joins("JOIN roles r ON rm.role_id = r.id").
		Joins("JOIN user_roles ur ON r.id = ur.role_id").
		Where("ur.user_id = ? AND m.visible = 1 AND ur.deleted_at IS NULL", userID).
		Distinct("m.id", "m.name", "m.path", "m.icon", "m.component", "m.parent_id", "m.order_num", "m.type", "m.permission").
		Order("m.order_num ASC").
		Find(&menus).Error

	return menus, err
}

// BuildMenuTree 构建菜单树
func (s *PermissionService) BuildMenuTree(menus []model.Menu) []model.Menu {
	// 创建 ID 到菜单的映射
	menuMap := make(map[int64]*model.Menu)
	for i := range menus {
		menuMap[menus[i].ID] = &menus[i]
	}

	// 构建树结构
	var roots []model.Menu
	for i := range menus {
		menu := menus[i]
		if menu.ParentID == nil {
			roots = append(roots, menu)
		} else if parent, ok := menuMap[*menu.ParentID]; ok {
			parent.Children = append(parent.Children, menu)
		}
	}

	return roots
}
