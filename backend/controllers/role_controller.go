package controllers

import (
	"learn-hub-backend/database"
	"learn-hub-backend/models"
	"learn-hub-backend/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RoleController struct{}

func NewRoleController() *RoleController {
	return &RoleController{}
}

// GetRoleList 获取角色列表
func (ctrl *RoleController) GetRoleList(c *gin.Context) {
	current, _ := strconv.Atoi(c.DefaultQuery("current", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	code := c.Query("code")
	status := c.Query("status")

	if current < 1 {
		current = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	offset := (current - 1) * pageSize

	// 构建查询
	query := database.DB.Model(&models.Role{})

	if code != "" {
		query = query.Where("code LIKE ?", "%"+code+"%")
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 查询总数
	var total int64
	query.Count(&total)

	// 查询列表
	var roles []models.Role
	query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&roles)

	// 转换为前端格式
	var list []map[string]interface{}
	for _, role := range roles {
		// 统计用户数
		userCount := database.DB.Model(&role).Association("Users").Count()

		// 统计菜单数
		menuCount := database.DB.Model(&role).Association("Menus").Count()

		list = append(list, map[string]interface{}{
			"id":          role.ID,
			"code":        role.Code,
			"name":        role.Name,
			"description": role.Description,
			"status":      role.Status,
			"userCount":   userCount,
			"menuCount":   menuCount,
			"createdAt":  role.CreatedAt.Format("2006-01-02 15:04:05"),
			"updatedAt":  role.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	c.JSON(200, gin.H{
		"data":    list,
		"total":   total,
		"success": true,
	})
}

// CreateRole 创建角色
func (ctrl *RoleController) CreateRole(c *gin.Context) {
	var req struct {
		Code        string `json:"code" binding:"required"`
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		Status      int    `json:"status"`
		MenuIds     []uint `json:"menuIds"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 检查角色代码是否已存在
	var existRole models.Role
	if err := database.DB.Where("code = ?", req.Code).First(&existRole).Error; err == nil {
		utils.Error(c, 200, "ROLE_EXISTS", "角色代码已存在")
		return
	}

	// 创建角色
	role := models.Role{
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		Status:      req.Status,
	}
	if role.Status == 0 {
		role.Status = 1
	}

	if err := database.DB.Create(&role).Error; err != nil {
		utils.InternalError(c, "创建失败: "+err.Error())
		return
	}

	// 分配菜单
	if len(req.MenuIds) > 0 {
		var menus []models.Menu
		database.DB.Where("id IN ?", req.MenuIds).Find(&menus)
		database.DB.Model(&role).Association("Menus").Append(&menus)
	}

	utils.Success(c, map[string]interface{}{
		"id":          role.ID,
		"code":        role.Code,
		"name":        role.Name,
		"description": role.Description,
		"status":      role.Status,
	})
}

// UpdateRole 更新角色
func (ctrl *RoleController) UpdateRole(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if id == 0 {
		utils.BadRequest(c, "角色ID不能为空")
		return
	}

	var role models.Role
	if err := database.DB.First(&role, id).Error; err != nil {
		utils.Error(c, 200, "ROLE_NOT_FOUND", "角色不存在")
		return
	}

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Status      *int   `json:"status"`
		MenuIds     []uint `json:"menuIds"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 更新字段
	if req.Name != "" {
		role.Name = req.Name
	}
	if req.Description != "" {
		role.Description = req.Description
	}
	if req.Status != nil {
		role.Status = *req.Status
	}

	if err := database.DB.Save(&role).Error; err != nil {
		utils.InternalError(c, "更新失败: "+err.Error())
		return
	}

	// 更新菜单
	if req.MenuIds != nil {
		var menus []models.Menu
		database.DB.Where("id IN ?", req.MenuIds).Find(&menus)
		database.DB.Model(&role).Association("Menus").Replace(&menus)
	}

	utils.Success(c, map[string]interface{}{
		"id":          role.ID,
		"code":        role.Code,
		"name":        role.Name,
		"description": role.Description,
		"status":      role.Status,
	})
}

// DeleteRole 删除角色
func (ctrl *RoleController) DeleteRole(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if id == 0 {
		utils.BadRequest(c, "角色ID不能为空")
		return
	}

	if err := database.DB.Delete(&models.Role{}, id).Error; err != nil {
		utils.InternalError(c, "删除失败: "+err.Error())
		return
	}

	utils.Success(c, gin.H{})
}

// GetAllMenus 获取所有菜单（用于角色分配菜单）
func (ctrl *RoleController) GetAllMenus(c *gin.Context) {
	var menus []models.Menu
	database.DB.Where("status = ?", 1).Order("parent_id, sort_order").Find(&menus)

	// 构建菜单树
	var list []map[string]interface{}
	for _, menu := range menus {
		list = append(list, map[string]interface{}{
			"id":        menu.ID,
			"parentId":  menu.ParentID,
			"name":      menu.Name,
			"path":      menu.Path,
			"component": menu.Component,
			"icon":      menu.Icon,
			"access":    menu.Access,
		})
	}

	utils.Success(c, list)
}

