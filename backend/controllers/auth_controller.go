package controllers

import (
	"log"
	"learn-hub-backend/services"
	"learn-hub-backend/utils"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	userService *services.UserService
}

func NewAuthController() *AuthController {
	return &AuthController{
		userService: &services.UserService{},
	}
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
	Type      string `json:"type"`
	AutoLogin bool   `json:"autoLogin"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Status          string `json:"status"`
	Type            string `json:"type"`
	CurrentAuthority string `json:"currentAuthority"`
	Token           string `json:"token,omitempty"`
}

// Login 用户登录
func (ctrl *AuthController) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 获取用户
	user, err := ctrl.userService.GetByUsername(req.Username)
	if err != nil {
		c.JSON(200, LoginResponse{
			Status:          "error",
			Type:            req.Type,
			CurrentAuthority: "guest",
		})
		return
	}

	// 验证密码
	if !ctrl.userService.VerifyPassword(user, req.Password) {
		// 调试日志（生产环境应删除）
		log.Printf("密码验证失败: username=%s, status=%d, password_hash_preview=%s", 
			user.Username, user.Status, 
			func() string {
				if len(user.Password) > 30 {
					return user.Password[:30] + "..."
				}
				return user.Password
			}())
		c.JSON(200, LoginResponse{
			Status:          "error",
			Type:            req.Type,
			CurrentAuthority: "guest",
		})
		return
	}

	// 检查用户状态
	if user.Status != 1 {
		c.JSON(200, LoginResponse{
			Status:          "error",
			Type:            req.Type,
			CurrentAuthority: "guest",
		})
		return
	}

	// 生成Token
	token, err := utils.GenerateToken(user.ID, user.Username, user.Access)
	if err != nil {
		utils.InternalError(c, "生成token失败")
		return
	}

	// 返回登录结果
	loginType := req.Type
	if loginType == "" {
		loginType = "account"
	}

	// 前端期望的响应格式：直接返回 LoginResponse，不包装在 data 中
	c.JSON(200, LoginResponse{
		Status:          "ok",
		Type:            loginType,
		CurrentAuthority: user.Access,
		Token:           token,
	})
}

// GetCurrentUser 获取当前用户信息
func (ctrl *AuthController) GetCurrentUser(c *gin.Context) {
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

	// 构建用户信息响应
	// 如果 name 为空，使用 username 作为 name
	displayName := user.Name
	if displayName == "" {
		displayName = user.Username
	}

	userInfo := map[string]interface{}{
		"name":        displayName,
		"avatar":      user.Avatar,
		"userid":      user.UserID,
		"email":       user.Email,
		"signature":   user.Signature,
		"title":       user.Title,
		"access":      user.Access,
		"country":     user.Country,
		"province":    user.Province,
		"city":        user.City,
		"address":     user.Address,
		"phone":       user.Phone,
		"notifyCount": 12, // 示例数据
		"unreadCount": 11, // 示例数据
	}

	utils.Success(c, userInfo)
}

// Logout 退出登录
func (ctrl *AuthController) Logout(c *gin.Context) {
	// 这里可以实现token黑名单机制
	utils.Success(c, gin.H{})
}

// GetCaptcha 获取验证码（手机登录用）
func (ctrl *AuthController) GetCaptcha(c *gin.Context) {
	phone := c.Query("phone")
	if phone == "" {
		utils.BadRequest(c, "手机号不能为空")
		return
	}

	// 这里应该实现真实的验证码发送逻辑
	// 目前返回固定验证码用于开发测试
	utils.Success(c, gin.H{
		"code":   200,
		"status": "ok",
		"captcha": "1234", // 开发环境返回，生产环境通过短信发送
	})
}

