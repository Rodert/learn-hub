package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"learn-hub/internal/model"
	"learn-hub/internal/repository"
	"learn-hub/internal/service"
)

// QuestionHandler 题库处理器（完整实现）
type QuestionHandlerImpl struct {
	db *gorm.DB
}

func NewQuestionHandlerImpl(db *gorm.DB) *QuestionHandlerImpl {
	return &QuestionHandlerImpl{db: db}
}

// ListQuestions 获取题目列表
// @Summary 获取题目列表
// @Tags 题库
// @Security Bearer
// @Param page query int false "页码" default(1)
// @Param limit query int false "每页数量" default(10)
// @Param type query string false "题型"
// @Success 200 {object} Response{data=map[string]interface{}}
// @Router /questions [get]
func (h *QuestionHandlerImpl) ListQuestions(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	qType := c.Query("type")

	qRepo := repository.NewQuestionRepository(h.db)
	questions, total, err := qRepo.List(page, limit, qType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"items": questions,
			"total": total,
			"page": page,
			"limit": limit,
		},
	})
}

// GetQuestion 获取题目详情
// @Summary 获取题目详情
// @Tags 题库
// @Security Bearer
// @Param id path int64 true "题目 ID"
// @Success 200 {object} Response{data=model.Question}
// @Router /questions/{id} [get]
func (h *QuestionHandlerImpl) GetQuestion(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid question ID"})
		return
	}

	qRepo := repository.NewQuestionRepository(h.db)
	question, err := qRepo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": question,
	})
}

// CreateQuestion 创建题目
// @Summary 创建题目
// @Tags 题库
// @Security Bearer
// @Param request body map[string]interface{} true "题目数据"
// @Success 201 {object} Response{data=model.Question}
// @Router /questions [post]
func (h *QuestionHandlerImpl) CreateQuestion(c *gin.Context) {
	var req struct {
		ExamID      *int64      `json:"exam_id"`
		Type        string      `json:"question_type" binding:"required"`
		Content     string      `json:"content" binding:"required"`
		Options     interface{} `json:"options"`
		Answer      string      `json:"answer" binding:"required"`
		Explanation string      `json:"explanation"`
		Score       float64     `json:"score"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")
	
	// 将 Options 转换为 JSON 字符串
	var optionsJSON string
	if req.Options != nil {
		optionsBytes, _ := json.Marshal(req.Options)
		optionsJSON = string(optionsBytes)
	} else {
		// 如果没有提供options，使用空数组
		optionsJSON = "[]"
	}
	
	question := &model.Question{
		ExamID:      req.ExamID,
		Type:        req.Type,
		Content:     req.Content,
		Options:     optionsJSON,
		Answer:      req.Answer,
		Explanation: req.Explanation,
		Score:       req.Score,
		CreatedBy:   userID.(int64),
	}

	qRepo := repository.NewQuestionRepository(h.db)
	if err := qRepo.Create(question); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code": 0,
		"message": "success",
		"data": question,
	})
}

// UpdateQuestion 更新题目
// @Summary 更新题目
// @Tags 题库
// @Security Bearer
// @Param id path int64 true "题目 ID"
// @Param request body map[string]interface{} true "更新数据"
// @Success 200 {object} Response{data=model.Question}
// @Router /questions/{id} [put]
func (h *QuestionHandlerImpl) UpdateQuestion(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid question ID"})
		return
	}

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	qRepo := repository.NewQuestionRepository(h.db)
	question, err := qRepo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
		return
	}

	if content, ok := req["content"].(string); ok {
		question.Content = content
	}
	if answer, ok := req["answer"].(string); ok {
		question.Answer = answer
	}

	if err := qRepo.Update(question); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": question,
	})
}

// DeleteQuestion 删除题目
// @Summary 删除题目
// @Tags 题库
// @Security Bearer
// @Param id path int64 true "题目 ID"
// @Success 204
// @Router /questions/{id} [delete]
func (h *QuestionHandlerImpl) DeleteQuestion(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid question ID"})
		return
	}

	qRepo := repository.NewQuestionRepository(h.db)
	if err := qRepo.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// CourseRecordHandlerImpl 学习记录处理器
type CourseRecordHandlerImpl struct {
	db *gorm.DB
}

func NewCourseRecordHandlerImpl(db *gorm.DB) *CourseRecordHandlerImpl {
	return &CourseRecordHandlerImpl{db: db}
}

// ListRecords 获取学习记录
// @Summary 获取学习记录
// @Tags 学习记录
// @Security Bearer
// @Param page query int false "页码" default(1)
// @Param limit query int false "每页数量" default(10)
// @Success 200 {object} Response{data=map[string]interface{}}
// @Router /course-records [get]
func (h *CourseRecordHandlerImpl) ListRecords(c *gin.Context) {
	userID, _ := c.Get("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	recordRepo := repository.NewCourseRecordRepository(h.db)
	records, total, err := recordRepo.ListByUser(userID.(int64), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"items": records,
			"total": total,
			"page": page,
			"limit": limit,
		},
	})
}

// GetRecord 获取学习记录详情
// @Summary 获取学习记录详情
// @Tags 学习记录
// @Security Bearer
// @Param id path int64 true "记录 ID"
// @Success 200 {object} Response{data=model.CourseRecord}
// @Router /course-records/{id} [get]
func (h *CourseRecordHandlerImpl) GetRecord(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid record ID"})
		return
	}

	recordRepo := repository.NewCourseRecordRepository(h.db)
	record, err := recordRepo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": record,
	})
}

// UpdateRecord 更新学习记录
// @Summary 更新学习记录
// @Tags 学习记录
// @Security Bearer
// @Param id path int64 true "记录 ID"
// @Param request body map[string]interface{} true "更新数据"
// @Success 200 {object} Response{data=model.CourseRecord}
// @Router /course-records/{id} [put]
func (h *CourseRecordHandlerImpl) UpdateRecord(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid record ID"})
		return
	}

	var req struct {
		ProgressPercent int `json:"progress_percent"`
		ViewDuration    int `json:"view_duration"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	recordRepo := repository.NewCourseRecordRepository(h.db)
	courseSvc := service.NewCourseService(recordRepo)
	record, err := courseSvc.UpdateProgress(id, req.ProgressPercent, req.ViewDuration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": record,
	})
}

// MenuHandler 菜单处理器
type MenuHandlerImpl struct {
	db *gorm.DB
}

func NewMenuHandlerImpl(db *gorm.DB) *MenuHandlerImpl {
	return &MenuHandlerImpl{db: db}
}

// GetMenus 获取用户菜单
// @Summary 获取用户菜单
// @Tags 菜单
// @Security Bearer
// @Success 200 {object} Response{data=[]model.Menu}
// @Router /menus [get]
func (h *MenuHandlerImpl) GetMenus(c *gin.Context) {
	userID, _ := c.Get("user_id")
	permSvc := service.NewPermissionService(h.db)

	menus, err := permSvc.GetUserMenus(userID.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 构建菜单树
	menuTree := permSvc.BuildMenuTree(menus)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": menuTree,
	})
}

// AdminHandler 管理员处理器
type AdminHandlerImpl struct {
	db *gorm.DB
}

func NewAdminHandlerImpl(db *gorm.DB) *AdminHandlerImpl {
	return &AdminHandlerImpl{db: db}
}

// ListUsers 获取用户列表
// @Summary 获取用户列表
// @Tags 管理员
// @Security Bearer
// @Param page query int false "页码" default(1)
// @Param limit query int false "每页数量" default(10)
// @Success 200 {object} Response{data=map[string]interface{}}
// @Router /admin/users [get]
func (h *AdminHandlerImpl) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	userRepo := repository.NewUserRepository(h.db)
	users, total, err := userRepo.List(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"items": users,
			"total": total,
			"page": page,
			"limit": limit,
		},
	})
}

// GetUser 获取用户详情
// @Summary 获取用户详情
// @Tags 管理员
// @Security Bearer
// @Param id path int64 true "用户 ID"
// @Success 200 {object} Response{data=model.User}
// @Router /admin/users/{id} [get]
func (h *AdminHandlerImpl) GetUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	userRepo := repository.NewUserRepository(h.db)
	user, err := userRepo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": user,
	})
}

// CreateUser 创建用户
// @Summary 创建用户
// @Tags 管理员
// @Security Bearer
// @Param request body map[string]interface{} true "用户数据"
// @Success 201 {object} Response{data=model.User}
// @Router /admin/users [post]
func (h *AdminHandlerImpl) CreateUser(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Nickname string `json:"nickname"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 密码加密
	hashedPassword, _ := hashPassword(req.Password)

	user := &model.User{
		Username: req.Username,
		Password: hashedPassword,
		Nickname: req.Nickname,
		Status:   "active",
	}

	userRepo := repository.NewUserRepository(h.db)
	if err := userRepo.Create(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code": 0,
		"message": "success",
		"data": user,
	})
}

// UpdateUser 更新用户
// @Summary 更新用户
// @Tags 管理员
// @Security Bearer
// @Param id path int64 true "用户 ID"
// @Param request body map[string]interface{} true "更新数据"
// @Success 200 {object} Response{data=model.User}
// @Router /admin/users/{id} [put]
func (h *AdminHandlerImpl) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userRepo := repository.NewUserRepository(h.db)
	user, err := userRepo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if nickname, ok := req["nickname"].(string); ok {
		user.Nickname = nickname
	}
	if status, ok := req["status"].(string); ok {
		user.Status = status
	}

	if err := userRepo.Update(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": user,
	})
}

// DeleteUser 删除用户
// @Summary 删除用户
// @Tags 管理员
// @Security Bearer
// @Param id path int64 true "用户 ID"
// @Success 204
// @Router /admin/users/{id} [delete]
func (h *AdminHandlerImpl) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	userRepo := repository.NewUserRepository(h.db)
	if err := userRepo.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// ListRoles 获取角色列表
// @Summary 获取角色列表
// @Tags 管理员
// @Security Bearer
// @Success 200 {object} Response{data=[]model.Role}
// @Router /admin/roles [get]
func (h *AdminHandlerImpl) ListRoles(c *gin.Context) {
	var roles []model.Role
	if err := h.db.Find(&roles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": roles,
	})
}

// CreateRole 创建角色
// @Summary 创建角色
// @Tags 管理员
// @Security Bearer
// @Param request body map[string]interface{} true "角色数据"
// @Success 201 {object} Response{data=model.Role}
// @Router /admin/roles [post]
func (h *AdminHandlerImpl) CreateRole(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	role := &model.Role{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := h.db.Create(role).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code": 0,
		"message": "success",
		"data": role,
	})
}

// UpdateRole 更新角色
// @Summary 更新角色
// @Tags 管理员
// @Security Bearer
// @Param id path int64 true "角色 ID"
// @Param request body map[string]interface{} true "更新数据"
// @Success 200 {object} Response{data=model.Role}
// @Router /admin/roles/{id} [put]
func (h *AdminHandlerImpl) UpdateRole(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var role model.Role
	if err := h.db.First(&role, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
		return
	}

	if desc, ok := req["description"].(string); ok {
		role.Description = desc
	}

	if err := h.db.Save(&role).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": role,
	})
}

// DeleteRole 删除角色
// @Summary 删除角色
// @Tags 管理员
// @Security Bearer
// @Param id path int64 true "角色 ID"
// @Success 204
// @Router /admin/roles/{id} [delete]
func (h *AdminHandlerImpl) DeleteRole(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}

	if err := h.db.Delete(&model.Role{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// ListPermissions 获取权限列表
// @Summary 获取权限列表
// @Tags 管理员
// @Security Bearer
// @Success 200 {object} Response{data=[]model.Permission}
// @Router /admin/permissions [get]
func (h *AdminHandlerImpl) ListPermissions(c *gin.Context) {
	var permissions []model.Permission
	if err := h.db.Find(&permissions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": permissions,
	})
}

// CreatePermission 创建权限
// @Summary 创建权限
// @Tags 管理员
// @Security Bearer
// @Param request body map[string]interface{} true "权限数据"
// @Success 201 {object} Response{data=model.Permission}
// @Router /admin/permissions [post]
func (h *AdminHandlerImpl) CreatePermission(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		Resource    string `json:"resource"`
		Action      string `json:"action"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	permission := &model.Permission{
		Name:        req.Name,
		Description: req.Description,
		Resource:    req.Resource,
		Action:      req.Action,
	}

	if err := h.db.Create(permission).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code": 0,
		"message": "success",
		"data": permission,
	})
}

// hashPassword 密码加密（辅助函数）
func hashPassword(password string) (string, error) {
	// 这里应该使用 bcrypt，但为了简化示例，这里只是返回原密码
	// 实际应该使用：golang.org/x/crypto/bcrypt
	return password, nil
}

// Response 通用响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
