package controllers

import (
	"learn-hub-backend/database"
	"learn-hub-backend/models"
	"learn-hub-backend/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CourseController struct{}

func NewCourseController() *CourseController {
	return &CourseController{}
}

// GetCourseList 获取课程列表（管理员）
func (ctrl *CourseController) GetCourseList(c *gin.Context) {
	current, _ := strconv.Atoi(c.DefaultQuery("current", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	title := c.Query("title")
	status := c.Query("status")

	if current < 1 {
		current = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	offset := (current - 1) * pageSize

	// 构建查询
	query := database.DB.Model(&models.Course{})

	if title != "" {
		query = query.Where("title LIKE ?", "%"+title+"%")
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 查询总数
	var total int64
	query.Count(&total)

	// 查询列表
	var courses []models.Course
	query.Offset(offset).Limit(pageSize).Order("sort_order ASC, created_at DESC").Find(&courses)

	// 转换为前端格式
	var list []map[string]interface{}
	for _, course := range courses {
		list = append(list, map[string]interface{}{
			"id":          course.ID,
			"title":       course.Title,
			"description": course.Description,
			"coverImage":  course.CoverImage,
			"contentType": course.ContentType,
			"videoUrl":    course.VideoURL,
			"textContent": course.TextContent,
			"duration":    course.Duration,
			"status":      course.Status,
			"sortOrder":   course.SortOrder,
			"createdAt":   course.CreatedAt.Format("2006-01-02 15:04:05"),
			"updatedAt":   course.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	c.JSON(200, gin.H{
		"data":    list,
		"total":   total,
		"success": true,
	})
}

// GetCourseDetail 获取课程详情
func (ctrl *CourseController) GetCourseDetail(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if id == 0 {
		utils.BadRequest(c, "课程ID不能为空")
		return
	}

	var course models.Course
	if err := database.DB.First(&course, id).Error; err != nil {
		utils.Error(c, 200, "COURSE_NOT_FOUND", "课程不存在")
		return
	}

	utils.Success(c, map[string]interface{}{
		"id":          course.ID,
		"title":       course.Title,
		"description": course.Description,
		"coverImage":  course.CoverImage,
		"contentType": course.ContentType,
		"videoUrl":    course.VideoURL,
		"textContent": course.TextContent,
		"duration":    course.Duration,
		"status":      course.Status,
		"sortOrder":   course.SortOrder,
		"createdAt":   course.CreatedAt.Format("2006-01-02 15:04:05"),
		"updatedAt":   course.UpdatedAt.Format("2006-01-02 15:04:05"),
	})
}

// CreateCourse 创建课程
func (ctrl *CourseController) CreateCourse(c *gin.Context) {
	var req struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description"`
		CoverImage  string `json:"coverImage"`
		ContentType int    `json:"contentType" binding:"required"` // 1-视频，2-文本，3-混合
		VideoURL    string `json:"videoUrl"`
		TextContent string `json:"textContent"`
		Duration    int    `json:"duration"`
		Status      int    `json:"status"`
		SortOrder   int    `json:"sortOrder"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 验证内容类型
	if req.ContentType < 1 || req.ContentType > 3 {
		utils.BadRequest(c, "内容类型错误，必须是1-视频，2-文本，3-混合")
		return
	}

	// 验证内容
	if req.ContentType == 1 && req.VideoURL == "" {
		utils.BadRequest(c, "视频类型课程必须提供视频URL")
		return
	}
	if req.ContentType == 2 && req.TextContent == "" {
		utils.BadRequest(c, "文本类型课程必须提供文本内容")
		return
	}

	course := models.Course{
		Title:       req.Title,
		Description: req.Description,
		CoverImage:  req.CoverImage,
		ContentType: req.ContentType,
		VideoURL:    req.VideoURL,
		TextContent: req.TextContent,
		Duration:    req.Duration,
		Status:      req.Status,
		SortOrder:   req.SortOrder,
	}

	if err := database.DB.Create(&course).Error; err != nil {
		utils.InternalError(c, "创建失败: "+err.Error())
		return
	}

	utils.Success(c, map[string]interface{}{
		"id":          course.ID,
		"title":       course.Title,
		"description": course.Description,
		"contentType": course.ContentType,
		"status":      course.Status,
	})
}

// UpdateCourse 更新课程
func (ctrl *CourseController) UpdateCourse(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if id == 0 {
		utils.BadRequest(c, "课程ID不能为空")
		return
	}

	var course models.Course
	if err := database.DB.First(&course, id).Error; err != nil {
		utils.Error(c, 200, "COURSE_NOT_FOUND", "课程不存在")
		return
	}

	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		CoverImage  string `json:"coverImage"`
		ContentType *int   `json:"contentType"`
		VideoURL    string `json:"videoUrl"`
		TextContent string `json:"textContent"`
		Duration    *int   `json:"duration"`
		Status      *int   `json:"status"`
		SortOrder   *int   `json:"sortOrder"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 更新字段
	if req.Title != "" {
		course.Title = req.Title
	}
	if req.Description != "" {
		course.Description = req.Description
	}
	if req.CoverImage != "" {
		course.CoverImage = req.CoverImage
	}
	if req.ContentType != nil {
		course.ContentType = *req.ContentType
	}
	if req.VideoURL != "" {
		course.VideoURL = req.VideoURL
	}
	if req.TextContent != "" {
		course.TextContent = req.TextContent
	}
	if req.Duration != nil {
		course.Duration = *req.Duration
	}
	if req.Status != nil {
		course.Status = *req.Status
	}
	if req.SortOrder != nil {
		course.SortOrder = *req.SortOrder
	}

	if err := database.DB.Save(&course).Error; err != nil {
		utils.InternalError(c, "更新失败: "+err.Error())
		return
	}

	utils.Success(c, map[string]interface{}{
		"id":          course.ID,
		"title":       course.Title,
		"description": course.Description,
		"contentType": course.ContentType,
		"status":      course.Status,
	})
}

// DeleteCourse 删除课程
func (ctrl *CourseController) DeleteCourse(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if id == 0 {
		utils.BadRequest(c, "课程ID不能为空")
		return
	}

	if err := database.DB.Delete(&models.Course{}, id).Error; err != nil {
		utils.InternalError(c, "删除失败: "+err.Error())
		return
	}

	utils.Success(c, gin.H{})
}

// PublishCourse 发布/下架课程
func (ctrl *CourseController) PublishCourse(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if id == 0 {
		utils.BadRequest(c, "课程ID不能为空")
		return
	}

	var course models.Course
	if err := database.DB.First(&course, id).Error; err != nil {
		utils.Error(c, 200, "COURSE_NOT_FOUND", "课程不存在")
		return
	}

	var req struct {
		Status int `json:"status" binding:"required"` // 1-发布，2-下架
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if req.Status != 1 && req.Status != 2 {
		utils.BadRequest(c, "状态错误，必须是1-发布或2-下架")
		return
	}

	course.Status = req.Status
	if err := database.DB.Save(&course).Error; err != nil {
		utils.InternalError(c, "操作失败: "+err.Error())
		return
	}

	utils.Success(c, map[string]interface{}{
		"id":     course.ID,
		"status": course.Status,
	})
}
