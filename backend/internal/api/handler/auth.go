package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"learn-hub/config"
	"learn-hub/internal/model"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	db  *gorm.DB
	cfg *config.Config
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler(db *gorm.DB, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		db:  db,
		cfg: cfg,
	}
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token string      `json:"token"`
	User  UserInfo    `json:"user"`
}

// UserInfo 用户信息
type UserInfo struct {
	ID        int64    `json:"id"`
	Username  string   `json:"username"`
	Nickname  string   `json:"nickname"`
	Roles     []string `json:"roles"`
	Permissions []string `json:"permissions"`
}

// Login 登录
// @Summary 用户登录
// @Description 使用账号密码登录
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body LoginRequest true "登录请求"
// @Success 200 {object} Response{data=LoginResponse}
// @Failure 400 {object} Response
// @Failure 401 {object} Response
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 查询用户
	var user model.User
	if err := h.db.Preload("Roles").First(&user, "username = ?", req.Username).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "用户不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "密码错误"})
		return
	}

	// 获取用户权限
	permissions, err := h.getUserPermissions(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 生成 JWT token
	token, err := h.generateToken(user, permissions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 获取角色名称
	roleNames := make([]string, len(user.Roles))
	for i, role := range user.Roles {
		roleNames[i] = role.Name
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": LoginResponse{
			Token: token,
			User: UserInfo{
				ID:          user.ID,
				Username:    user.Username,
				Nickname:    user.Nickname,
				Roles:       roleNames,
				Permissions: permissions,
			},
		},
	})
}

// Register 注册
// @Summary 用户注册
// @Description 创建新用户账号
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body LoginRequest true "注册请求"
// @Success 201 {object} Response{data=UserInfo}
// @Failure 400 {object} Response
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查用户是否已存在
	var count int64
	h.db.Model(&model.User{}).Where("username = ?", req.Username).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户已存在"})
		return
	}

	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 创建用户
	user := model.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Nickname: req.Username,
		Status:   "active",
	}

	if err := h.db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 分配默认角色（user）
	var userRole model.Role
	h.db.First(&userRole, "name = ?", "user")
	if userRole.ID > 0 {
		h.db.Create(&model.UserRole{
			UserID: user.ID,
			RoleID: userRole.ID,
		})
	}

	c.JSON(http.StatusCreated, gin.H{
		"code": 0,
		"message": "success",
		"data": UserInfo{
			ID:       user.ID,
			Username: user.Username,
			Nickname: user.Nickname,
		},
	})
}

// RefreshToken 刷新 token
// @Summary 刷新 JWT Token
// @Description 使用旧 token 获取新 token
// @Tags 认证
// @Security Bearer
// @Produce json
// @Success 200 {object} Response{data=map[string]string}
// @Failure 401 {object} Response
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	// TODO: 实现 token 刷新逻辑
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"token": "new-token",
		},
	})
}

// generateToken 生成 JWT token
func (h *AuthHandler) generateToken(user model.User, permissions []string) (string, error) {
	roleNames := make([]string, len(user.Roles))
	for i, role := range user.Roles {
		roleNames[i] = role.Name
	}

	claims := jwt.MapClaims{
		"sub":         user.ID,
		"username":    user.Username,
		"roles":       roleNames,
		"permissions": permissions,
		"iat":         time.Now().Unix(),
		"exp":         time.Now().Add(time.Duration(h.cfg.JWT.ExpireHours) * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.cfg.JWT.Secret))
}

// getUserPermissions 获取用户权限
func (h *AuthHandler) getUserPermissions(userID int64) ([]string, error) {
	var permissions []string
	err := h.db.
		Table("permissions p").
		Joins("JOIN role_permissions rp ON p.id = rp.permission_id").
		Joins("JOIN roles r ON rp.role_id = r.id").
		Joins("JOIN user_roles ur ON r.id = ur.role_id").
		Where("ur.user_id = ?", userID).
		Distinct("p.name").
		Pluck("p.name", &permissions).Error

	return permissions, err
}
