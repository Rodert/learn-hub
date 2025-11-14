package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"learn-hub/internal/model"
	"learn-hub/internal/repository"
	"learn-hub/internal/service"
)

// UserHandler 用户处理器
type UserHandler struct {
	db *gorm.DB
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{db: db}
}

// GetUser 获取用户信息
// @Summary 获取用户信息
// @Tags 用户
// @Security Bearer
// @Param id path int64 true "用户 ID"
// @Success 200 {object} Response{data=model.User}
// @Router /users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	userRepo := repository.NewUserRepository(h.db)
	user, err := userRepo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": user,
	})
}

// UpdateUser 更新用户信息
// @Summary 更新用户信息
// @Tags 用户
// @Security Bearer
// @Param id path int64 true "用户 ID"
// @Param request body map[string]interface{} true "更新数据"
// @Success 200 {object} Response{data=model.User}
// @Router /users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
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

	// 更新允许的字段
	if nickname, ok := req["nickname"].(string); ok {
		user.Nickname = nickname
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

// GetProfile 获取当前用户信息
// @Summary 获取当前用户信息
// @Tags 用户
// @Security Bearer
// @Success 200 {object} Response{data=model.User}
// @Router /users/profile/me [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userRepo := repository.NewUserRepository(h.db)
	user, err := userRepo.GetByID(userID.(int64))
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

// MaterialHandler 资料处理器
type MaterialHandler struct {
	db *gorm.DB
}

func NewMaterialHandler(db *gorm.DB) *MaterialHandler {
	return &MaterialHandler{db: db}
}

// ListMaterials 获取资料列表
// @Summary 获取资料列表
// @Tags 资料
// @Security Bearer
// @Param page query int false "页码" default(1)
// @Param limit query int false "每页数量" default(10)
// @Param status query string false "状态"
// @Success 200 {object} Response{data=map[string]interface{}}
// @Router /materials [get]
func (h *MaterialHandler) ListMaterials(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	status := c.Query("status")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	materialRepo := repository.NewMaterialRepository(h.db)
	materials, total, err := materialRepo.List(page, limit, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"items": materials,
			"total": total,
			"page": page,
			"limit": limit,
		},
	})
}

// GetMaterial 获取资料详情
// @Summary 获取资料详情
// @Tags 资料
// @Security Bearer
// @Param id path int64 true "资料 ID"
// @Success 200 {object} Response{data=model.Material}
// @Router /materials/{id} [get]
func (h *MaterialHandler) GetMaterial(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid material ID"})
		return
	}

	materialRepo := repository.NewMaterialRepository(h.db)
	material, err := materialRepo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Material not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": material,
	})
}

// CreateMaterial 创建资料
// @Summary 创建资料
// @Tags 资料
// @Security Bearer
// @Param request body map[string]interface{} true "资料数据"
// @Success 201 {object} Response{data=model.Material}
// @Router /materials [post]
func (h *MaterialHandler) CreateMaterial(c *gin.Context) {
	var req struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description"`
		ContentType string `json:"content_type"`
		Content     string `json:"content"`
		FileURL     string `json:"file_url"`
		FileSize    int64  `json:"file_size"`
		CoverURL    string `json:"cover_url"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")
	material := &model.Material{
		Title:     req.Title,
		Desc:      req.Description,
		Type:      req.ContentType,
		Content:   req.Content,
		FileURL:   req.FileURL,
		FileSize:  req.FileSize,
		CoverURL:  req.CoverURL,
		Status:    "draft",
		CreatedBy: userID.(int64),
	}

	materialRepo := repository.NewMaterialRepository(h.db)
	if err := materialRepo.Create(material); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code": 0,
		"message": "success",
		"data": material,
	})
}

// UpdateMaterial 更新资料
// @Summary 更新资料
// @Tags 资料
// @Security Bearer
// @Param id path int64 true "资料 ID"
// @Param request body map[string]interface{} true "更新数据"
// @Success 200 {object} Response{data=model.Material}
// @Router /materials/{id} [put]
func (h *MaterialHandler) UpdateMaterial(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid material ID"})
		return
	}

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	materialRepo := repository.NewMaterialRepository(h.db)
	material, err := materialRepo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Material not found"})
		return
	}

	// 更新字段
	if title, ok := req["title"].(string); ok {
		material.Title = title
	}
	if desc, ok := req["description"].(string); ok {
		material.Desc = desc
	}
	if status, ok := req["status"].(string); ok {
		material.Status = status
	}

	if err := materialRepo.Update(material); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": material,
	})
}

// DeleteMaterial 删除资料
// @Summary 删除资料
// @Tags 资料
// @Security Bearer
// @Param id path int64 true "资料 ID"
// @Success 204
// @Router /materials/{id} [delete]
func (h *MaterialHandler) DeleteMaterial(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid material ID"})
		return
	}

	materialRepo := repository.NewMaterialRepository(h.db)
	if err := materialRepo.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// ExamHandler 考试处理器
type ExamHandler struct {
	db *gorm.DB
}

func NewExamHandler(db *gorm.DB) *ExamHandler {
	return &ExamHandler{db: db}
}

// ListExams 获取试卷列表
// @Summary 获取试卷列表
// @Tags 考试
// @Security Bearer
// @Param page query int false "页码" default(1)
// @Param limit query int false "每页数量" default(10)
// @Success 200 {object} Response{data=map[string]interface{}}
// @Router /exams [get]
func (h *ExamHandler) ListExams(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	examRepo := repository.NewExamRepository(h.db)
	exams, total, err := examRepo.List(page, limit, "published")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"items": exams,
			"total": total,
			"page": page,
			"limit": limit,
		},
	})
}

// GetExam 获取试卷详情
// @Summary 获取试卷详情
// @Tags 考试
// @Security Bearer
// @Param id path int64 true "试卷 ID"
// @Success 200 {object} Response{data=model.Exam}
// @Router /exams/{id} [get]
func (h *ExamHandler) GetExam(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid exam ID"})
		return
	}

	examRepo := repository.NewExamRepository(h.db)
	exam, err := examRepo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Exam not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": exam,
	})
}

// CreateExam 创建试卷
// @Summary 创建试卷
// @Tags 考试
// @Security Bearer
// @Param request body map[string]interface{} true "试卷数据"
// @Success 201 {object} Response{data=model.Exam}
// @Router /exams [post]
func (h *ExamHandler) CreateExam(c *gin.Context) {
	var req struct {
		Title      string  `json:"title" binding:"required"`
		Description string `json:"description"`
		TotalScore float64 `json:"total_score"`
		PassScore  float64 `json:"pass_score"`
		TimeLimit  int     `json:"time_limit"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")
	exam := &model.Exam{
		Title:      req.Title,
		Desc:       req.Description,
		TotalScore: req.TotalScore,
		PassScore:  req.PassScore,
		TimeLimit:  req.TimeLimit,
		Status:     "draft",
		CreatedBy:  userID.(int64),
	}

	examRepo := repository.NewExamRepository(h.db)
	if err := examRepo.Create(exam); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code": 0,
		"message": "success",
		"data": exam,
	})
}

// UpdateExam 更新试卷
// @Summary 更新试卷
// @Tags 考试
// @Security Bearer
// @Param id path int64 true "试卷 ID"
// @Param request body map[string]interface{} true "更新数据"
// @Success 200 {object} Response{data=model.Exam}
// @Router /exams/{id} [put]
func (h *ExamHandler) UpdateExam(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid exam ID"})
		return
	}

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	examRepo := repository.NewExamRepository(h.db)
	exam, err := examRepo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Exam not found"})
		return
	}

	if title, ok := req["title"].(string); ok {
		exam.Title = title
	}
	if status, ok := req["status"].(string); ok {
		exam.Status = status
	}

	if err := examRepo.Update(exam); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": exam,
	})
}

// DeleteExam 删除试卷
// @Summary 删除试卷
// @Tags 考试
// @Security Bearer
// @Param id path int64 true "试卷 ID"
// @Success 204
// @Router /exams/{id} [delete]
func (h *ExamHandler) DeleteExam(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid exam ID"})
		return
	}

	examRepo := repository.NewExamRepository(h.db)
	if err := examRepo.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// StartExam 开始考试
// @Summary 开始考试
// @Tags 考试
// @Security Bearer
// @Param id path int64 true "试卷 ID"
// @Success 200 {object} Response{data=map[string]interface{}}
// @Router /exams/{id}/start [post]
func (h *ExamHandler) StartExam(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid exam ID"})
		return
	}

	userID, _ := c.Get("user_id")
	examRepo := repository.NewExamRepository(h.db)
	qRepo := repository.NewQuestionRepository(h.db)
	recordRepo := repository.NewExamRecordRepository(h.db)

	// 获取试卷信息
	exam, err := examRepo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Exam not found"})
		return
	}

	// 获取题目
	questions, err := qRepo.ListByExam(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 创建考试记录
	examSvc := service.NewExamService(examRepo, recordRepo, qRepo)
	record, err := examSvc.StartExam(userID.(int64), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"exam_record_id": record.ID,
			"exam": exam,
			"questions": questions,
		},
	})
}

// SubmitExam 提交答卷
// @Summary 提交答卷
// @Tags 考试
// @Security Bearer
// @Param id path int64 true "试卷 ID"
// @Param request body map[string]interface{} true "答题数据"
// @Success 200 {object} Response{data=map[string]interface{}}
// @Router /exams/{id}/submit [post]
func (h *ExamHandler) SubmitExam(c *gin.Context) {
	var req struct {
		ExamRecordID int64 `json:"exam_record_id" binding:"required"`
		Answers []struct {
			QuestionID int64  `json:"question_id"`
			Answer     string `json:"answer"`
		} `json:"answers"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	examRepo := repository.NewExamRepository(h.db)
	qRepo := repository.NewQuestionRepository(h.db)
	recordRepo := repository.NewExamRecordRepository(h.db)

	svc := service.NewExamService(examRepo, recordRepo, qRepo)
	
	// 转换答题数据
	answers := make([]service.AnswerItem, len(req.Answers))
	for i, a := range req.Answers {
		answers[i] = service.AnswerItem{
			QuestionID: a.QuestionID,
			Answer:     a.Answer,
		}
	}

	score, err := svc.SubmitExam(req.ExamRecordID, answers)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"score": score,
		},
	})
}

// GetExamRecords 获取考试成绩
// @Summary 获取考试成绩
// @Tags 考试
// @Security Bearer
// @Param id path int64 true "试卷 ID"
// @Success 200 {object} Response{data=map[string]interface{}}
// @Router /exams/{id}/records [get]
func (h *ExamHandler) GetExamRecords(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid exam ID"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	recordRepo := repository.NewExamRecordRepository(h.db)
	records, total, err := recordRepo.ListByExam(id, page, limit)
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

// QuestionHandler 题库处理器
type QuestionHandler struct {
	db *gorm.DB
}

func NewQuestionHandler(db *gorm.DB) *QuestionHandler {
	return &QuestionHandler{db: db}
}

func (h *QuestionHandler) ListQuestions(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "ListQuestions"})
}

func (h *QuestionHandler) GetQuestion(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "GetQuestion"})
}

func (h *QuestionHandler) CreateQuestion(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"message": "CreateQuestion"})
}

func (h *QuestionHandler) UpdateQuestion(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "UpdateQuestion"})
}

func (h *QuestionHandler) DeleteQuestion(c *gin.Context) {
	c.JSON(http.StatusNoContent, gin.H{})
}

// CourseRecordHandler 学习记录处理器
type CourseRecordHandler struct {
	db *gorm.DB
}

func NewCourseRecordHandler(db *gorm.DB) *CourseRecordHandler {
	return &CourseRecordHandler{db: db}
}

func (h *CourseRecordHandler) ListRecords(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "ListRecords"})
}

func (h *CourseRecordHandler) GetRecord(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "GetRecord"})
}

func (h *CourseRecordHandler) UpdateRecord(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "UpdateRecord"})
}

// MenuHandler 菜单处理器
type MenuHandler struct {
	db *gorm.DB
}

func NewMenuHandler(db *gorm.DB) *MenuHandler {
	return &MenuHandler{db: db}
}

func (h *MenuHandler) GetMenus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "GetMenus"})
}

// AdminHandler 管理员处理器
type AdminHandler struct {
	db *gorm.DB
}

func NewAdminHandler(db *gorm.DB) *AdminHandler {
	return &AdminHandler{db: db}
}

func (h *AdminHandler) ListUsers(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "ListUsers"})
}

func (h *AdminHandler) GetUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "GetUser"})
}

func (h *AdminHandler) CreateUser(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"message": "CreateUser"})
}

func (h *AdminHandler) UpdateUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "UpdateUser"})
}

func (h *AdminHandler) DeleteUser(c *gin.Context) {
	c.JSON(http.StatusNoContent, gin.H{})
}

func (h *AdminHandler) ListRoles(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "ListRoles"})
}

func (h *AdminHandler) CreateRole(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"message": "CreateRole"})
}

func (h *AdminHandler) UpdateRole(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "UpdateRole"})
}

func (h *AdminHandler) DeleteRole(c *gin.Context) {
	c.JSON(http.StatusNoContent, gin.H{})
}

func (h *AdminHandler) ListPermissions(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "ListPermissions"})
}

func (h *AdminHandler) CreatePermission(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"message": "CreatePermission"})
}
