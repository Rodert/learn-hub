package services

import (
	"learn-hub-backend/database"
	"learn-hub-backend/models"
	"learn-hub-backend/utils"
)

// UserService 用户服务
type UserService struct{}

// GetByUsername 根据用户名获取用户
func (s *UserService) GetByUsername(username string) (*models.User, error) {
	var user models.User
	err := database.DB.Where("username = ?", username).Preload("Roles").First(&user).Error
	return &user, err
}

// GetByID 根据ID获取用户
func (s *UserService) GetByID(id uint) (*models.User, error) {
	var user models.User
	err := database.DB.Preload("Roles").First(&user, id).Error
	return &user, err
}

// VerifyPassword 验证密码
func (s *UserService) VerifyPassword(user *models.User, password string) bool {
	return utils.CheckPassword(password, user.Password)
}

// GetUserMenus 获取用户菜单列表
func (s *UserService) GetUserMenus(userID uint) ([]models.MenuVO, error) {
	var user models.User
	if err := database.DB.Preload("Roles.Menus").First(&user, userID).Error; err != nil {
		return nil, err
	}

	// 收集所有菜单ID（去重）
	menuMap := make(map[uint]models.Menu)
	for _, role := range user.Roles {
		for _, menu := range role.Menus {
			if menu.Status == 1 {
				menuMap[menu.ID] = menu
			}
		}
	}

	// 转换为切片
	var menus []models.Menu
	for _, menu := range menuMap {
		menus = append(menus, menu)
	}

	// 构建菜单树
	return buildMenuTree(menus, 0), nil
}

// buildMenuTree 构建菜单树
func buildMenuTree(menus []models.Menu, parentID uint) []models.MenuVO {
	var tree []models.MenuVO

	for _, menu := range menus {
		if menu.ParentID == parentID && menu.Hidden == 0 {
			vo := models.MenuVO{}

			// 只设置非空字段
			if menu.Path != "" {
				vo.Path = menu.Path
			}
			if menu.Name != "" {
				vo.Name = menu.Name
			}
			if menu.Icon != "" {
				vo.Icon = menu.Icon
			}
			if menu.Component != "" {
				vo.Component = menu.Component
			}
			if menu.Access != "" {
				vo.Access = menu.Access
			}
			if menu.Redirect != "" {
				vo.Redirect = menu.Redirect
			}

			// 递归构建子菜单
			children := buildMenuTree(menus, menu.ID)
			if len(children) > 0 {
				vo.Routes = children
			}

			tree = append(tree, vo)
		}
	}

	return tree
}

