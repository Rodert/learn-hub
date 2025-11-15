package controllers

import (
	"learn-hub-backend/services"
	"learn-hub-backend/utils"

	"github.com/gin-gonic/gin"
)

type MenuController struct {
	userService *services.UserService
}

func NewMenuController() *MenuController {
	return &MenuController{
		userService: &services.UserService{},
	}
}

// GetMenuList 获取用户菜单列表
func (ctrl *MenuController) GetMenuList(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		utils.Unauthorized(c, "未登录")
		return
	}

	menus, err := ctrl.userService.GetUserMenus(userID.(uint))
	if err != nil {
		utils.InternalError(c, "获取菜单失败")
		return
	}

	utils.Success(c, menus)
}

// GetUserPermissions 获取用户权限列表
func (ctrl *MenuController) GetUserPermissions(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		utils.Unauthorized(c, "未登录")
		return
	}

	user, err := ctrl.userService.GetByID(userID.(uint))
	if err != nil {
		utils.Error(c, 200, "USER_NOT_FOUND", "用户不存在")
		return
	}

	// 收集权限
	var accessList []string
	if user.Access != "" {
		accessList = append(accessList, "can"+user.Access) // 转换为前端格式
	}

	var roles []string
	for _, role := range user.Roles {
		roles = append(roles, role.Code)
	}

	utils.Success(c, gin.H{
		"access":      accessList,
		"roles":       roles,
		"permissions": []string{}, // 可以扩展细粒度权限
	})
}

