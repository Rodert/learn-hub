package controllers

import (
	"fmt"
	"learn-hub-backend/database"
	"learn-hub-backend/models"
	"learn-hub-backend/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserController struct{}

func NewUserController() *UserController {
	return &UserController{}
}

// GetUserList 获取用户列表
func (ctrl *UserController) GetUserList(c *gin.Context) {
	current, _ := strconv.Atoi(c.DefaultQuery("current", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	username := c.Query("username")
	status := c.Query("status")

	if current < 1 {
		current = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	offset := (current - 1) * pageSize

	// 构建查询
	query := database.DB.Model(&models.User{})

	if username != "" {
		query = query.Where("username LIKE ?", "%"+username+"%")
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 查询总数
	var total int64
	query.Count(&total)

	// 查询列表
	var users []models.User
	query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&users)

	// 转换为前端格式
	var list []map[string]interface{}
	for _, user := range users {
		// 获取用户角色
		var roles []string
		database.DB.Model(&user).Association("Roles").Find(&user.Roles)
		for _, role := range user.Roles {
			roles = append(roles, role.Name)
		}

		list = append(list, map[string]interface{}{
			"id":        user.ID,
			"username":  user.Username,
			"name":      user.Name,
			"email":     user.Email,
			"phone":     user.Phone,
			"avatar":    user.Avatar,
			"userid":    user.UserID,
			"access":    user.Access,
			"status":    user.Status,
			"roles":     roles,
			"createdAt": user.CreatedAt.Format("2006-01-02 15:04:05"),
			"updatedAt": user.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	c.JSON(200, gin.H{
		"data":    list,
		"total":   total,
		"success": true,
	})
}

// CreateUser 创建用户
func (ctrl *UserController) CreateUser(c *gin.Context) {
	var req struct {
		Username string   `json:"username" binding:"required"`
		Password string   `json:"password" binding:"required"`
		Name     string   `json:"name"`
		Email    string   `json:"email"`
		Phone    string   `json:"phone"`
		Access   string   `json:"access"`
		Status   int      `json:"status"`
		RoleIds  []uint   `json:"roleIds"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 检查用户名是否已存在
	var existUser models.User
	if err := database.DB.Where("username = ?", req.Username).First(&existUser).Error; err == nil {
		utils.Error(c, 200, "USER_EXISTS", "用户名已存在")
		return
	}

	// 加密密码
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.InternalError(c, "密码加密失败")
		return
	}

	// 生成唯一的 UserID（8位数字，从00000001开始）
	var nextUserID string
	var users []models.User
	// 查找所有有 UserID 的用户
	database.DB.Where("user_id != '' AND user_id IS NOT NULL").
		Select("user_id").
		Find(&users)
	
	var maxID uint = 0
	for _, u := range users {
		if u.UserID != "" {
			var id uint
			if _, err := fmt.Sscanf(u.UserID, "%d", &id); err == nil {
				if id > maxID {
					maxID = id
				}
			}
		}
	}
	
	// 生成下一个 UserID
	nextID := maxID + 1
	nextUserID = fmt.Sprintf("%08d", nextID)

	// 创建用户
	user := models.User{
		Username: req.Username,
		Password: hashedPassword,
		Name:     req.Name,
		Email:    req.Email,
		Phone:    req.Phone,
		Access:   req.Access,
		Status:   req.Status,
		UserID:   nextUserID,
	}
	if user.Status == 0 {
		user.Status = 1
	}
	if user.Access == "" {
		user.Access = "user"
	}

	if err := database.DB.Create(&user).Error; err != nil {
		utils.InternalError(c, "创建失败: "+err.Error())
		return
	}

	// 分配角色
	if len(req.RoleIds) > 0 {
		var roles []models.Role
		database.DB.Where("id IN ?", req.RoleIds).Find(&roles)
		database.DB.Model(&user).Association("Roles").Append(&roles)
		
		// 根据角色自动设置 access 字段
		// 如果用户有 admin 角色，access 设为 "admin"，否则设为 "user"
		hasAdminRole := false
		for _, role := range roles {
			if role.Code == "admin" {
				hasAdminRole = true
				break
			}
		}
		if hasAdminRole {
			user.Access = "admin"
		} else if user.Access == "" {
			user.Access = "user"
		}
		// 保存 access 字段的更新
		database.DB.Model(&user).Update("access", user.Access)
	} else if user.Access == "" {
		// 如果没有分配角色，默认设为 "user"
		user.Access = "user"
		database.DB.Model(&user).Update("access", user.Access)
	}

	utils.Success(c, map[string]interface{}{
		"id":       user.ID,
		"username": user.Username,
		"name":     user.Name,
		"email":    user.Email,
		"status":   user.Status,
		"access":   user.Access,
	})
}

// UpdateUser 更新用户
func (ctrl *UserController) UpdateUser(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if id == 0 {
		utils.BadRequest(c, "用户ID不能为空")
		return
	}

	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		utils.Error(c, 200, "USER_NOT_FOUND", "用户不存在")
		return
	}

	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
		Password string `json:"password"`
		Access   string `json:"access"`
		Status   *int   `json:"status"`
		RoleIds  []uint `json:"roleIds"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 更新字段（注意：access 字段会在更新角色后自动设置）
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	// 只有在没有更新角色的情况下，才允许手动设置 access
	if req.Access != "" && req.RoleIds == nil {
		user.Access = req.Access
	}
	if req.Password != "" {
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			utils.InternalError(c, "密码加密失败")
			return
		}
		user.Password = hashedPassword
	}
	if req.Status != nil {
		user.Status = *req.Status
	}

	// 先保存基本字段的更新
	if err := database.DB.Save(&user).Error; err != nil {
		utils.InternalError(c, "更新失败: "+err.Error())
		return
	}

	// 更新角色（如果提供了角色列表）
	if req.RoleIds != nil {
		var roles []models.Role
		database.DB.Where("id IN ?", req.RoleIds).Find(&roles)
		database.DB.Model(&user).Association("Roles").Replace(&roles)
		
		// 根据角色自动设置 access 字段
		// 如果用户有 admin 角色，access 设为 "admin"，否则设为 "user"
		hasAdminRole := false
		for _, role := range roles {
			if role.Code == "admin" {
				hasAdminRole = true
				break
			}
		}
		if hasAdminRole {
			user.Access = "admin"
		} else {
			user.Access = "user"
		}
		// 保存 access 字段的更新
		database.DB.Model(&user).Update("access", user.Access)
	}

	utils.Success(c, map[string]interface{}{
		"id":       user.ID,
		"username": user.Username,
		"name":     user.Name,
		"email":    user.Email,
		"status":   user.Status,
		"access":   user.Access,
	})
}

// DeleteUser 删除用户
func (ctrl *UserController) DeleteUser(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if id == 0 {
		utils.BadRequest(c, "用户ID不能为空")
		return
	}

	if err := database.DB.Delete(&models.User{}, id).Error; err != nil {
		utils.InternalError(c, "删除失败: "+err.Error())
		return
	}

	utils.Success(c, gin.H{})
}

// GetAllRoles 获取所有角色（用于下拉选择）
func (ctrl *UserController) GetAllRoles(c *gin.Context) {
	var roles []models.Role
	database.DB.Where("status = ?", 1).Find(&roles)

	var list []map[string]interface{}
	for _, role := range roles {
		list = append(list, map[string]interface{}{
			"id":   role.ID,
			"code": role.Code,
			"name": role.Name,
		})
	}

	utils.Success(c, list)
}

