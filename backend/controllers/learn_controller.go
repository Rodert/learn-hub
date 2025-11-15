package controllers

import (
	"learn-hub-backend/database"
	"learn-hub-backend/models"
	"learn-hub-backend/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type LearnController struct{}

func NewLearnController() *LearnController {
	return &LearnController{}
}

// GetCourses 获取可学习课程列表（员工端）
func (ctrl *LearnController) GetCourses(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		utils.Unauthorized(c, "未登录")
		return
	}

	current, _ := strconv.Atoi(c.DefaultQuery("current", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	title := c.Query("title")

	if current < 1 {
		current = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	offset := (current - 1) * pageSize

	// 构建查询（只查询已发布的课程）
	query := database.DB.Model(&models.Course{}).Where("status = ?", 1)

	if title != "" {
		query = query.Where("title LIKE ?", "%"+title+"%")
	}

	// 查询总数
	var total int64
	query.Count(&total)

	// 查询列表
	var courses []models.Course
	query.Offset(offset).Limit(pageSize).Order("sort_order ASC, created_at DESC").Find(&courses)

	// 获取用户的学习记录
	userIDUint := userID.(uint)
	var records []models.CourseRecord
	database.DB.Where("user_id = ?", userIDUint).Find(&records)

	// 构建课程ID到记录的映射
	recordMap := make(map[uint]*models.CourseRecord)
	for i := range records {
		recordMap[records[i].CourseID] = &records[i]
	}

	// 转换为前端格式
	var list []map[string]interface{}
	for _, course := range courses {
		record, hasRecord := recordMap[course.ID]

		item := map[string]interface{}{
			"id":          course.ID,
			"title":       course.Title,
			"description": course.Description,
			"coverImage":  course.CoverImage,
			"contentType": course.ContentType,
			"duration":    course.Duration,
		}

		if hasRecord {
			item["progress"] = record.Progress
			item["isCompleted"] = record.IsCompleted
			item["lastStudyAt"] = record.LastStudyAt.Format("2006-01-02 15:04:05")
		} else {
			item["progress"] = 0
			item["isCompleted"] = false
			item["lastStudyAt"] = ""
		}

		list = append(list, item)
	}

	c.JSON(200, gin.H{
		"data":    list,
		"total":   total,
		"success": true,
	})
}

// GetCourseDetail 获取课程详情（员工端）
func (ctrl *LearnController) GetCourseDetail(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		utils.Unauthorized(c, "未登录")
		return
	}

	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if id == 0 {
		utils.BadRequest(c, "课程ID不能为空")
		return
	}

	var course models.Course
	if err := database.DB.Where("status = ?", 1).First(&course, id).Error; err != nil {
		utils.Error(c, 200, "COURSE_NOT_FOUND", "课程不存在或已下架")
		return
	}

	// 获取用户的学习记录
	userIDUint := userID.(uint)
	var record models.CourseRecord
	hasRecord := database.DB.Where("user_id = ? AND course_id = ?", userIDUint, id).First(&record).Error == nil

	result := map[string]interface{}{
		"id":          course.ID,
		"title":       course.Title,
		"description": course.Description,
		"coverImage":  course.CoverImage,
		"contentType": course.ContentType,
		"videoUrl":    course.VideoURL,
		"textContent": course.TextContent,
		"duration":    course.Duration,
	}

	if hasRecord {
		result["progress"] = record.Progress
		result["duration"] = record.Duration
		result["isCompleted"] = record.IsCompleted
		result["lastStudyAt"] = record.LastStudyAt.Format("2006-01-02 15:04:05")
	} else {
		result["progress"] = 0
		result["duration"] = 0
		result["isCompleted"] = false
		result["lastStudyAt"] = ""
	}

	utils.Success(c, result)
}

// UpdateProgress 更新学习进度
func (ctrl *LearnController) UpdateProgress(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		utils.Unauthorized(c, "未登录")
		return
	}

	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if id == 0 {
		utils.BadRequest(c, "课程ID不能为空")
		return
	}

	var req struct {
		Progress int `json:"progress" binding:"required"` // 0-100
		Duration int `json:"duration"`                    // 已学习时长（秒）
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if req.Progress < 0 || req.Progress > 100 {
		utils.BadRequest(c, "进度必须在0-100之间")
		return
	}

	userIDUint := userID.(uint)

	// 查找或创建学习记录
	var record models.CourseRecord
	err := database.DB.Where("user_id = ? AND course_id = ?", userIDUint, id).First(&record).Error

	if err != nil {
		// 创建新记录
		record = models.CourseRecord{
			UserID:      userIDUint,
			CourseID:    uint(id),
			Progress:    req.Progress,
			Duration:    req.Duration,
			IsCompleted: req.Progress >= 100,
			LastStudyAt: time.Now(),
		}
		if record.IsCompleted {
			now := time.Now()
			record.CompletedAt = &now
		}
		if err := database.DB.Create(&record).Error; err != nil {
			utils.InternalError(c, "创建学习记录失败: "+err.Error())
			return
		}
	} else {
		// 更新记录
		record.Progress = req.Progress
		if req.Duration > 0 {
			record.Duration = req.Duration
		}
		record.LastStudyAt = time.Now()

		// 如果进度达到100%，标记为完成
		if req.Progress >= 100 && !record.IsCompleted {
			record.IsCompleted = true
			now := time.Now()
			record.CompletedAt = &now
		}

		if err := database.DB.Save(&record).Error; err != nil {
			utils.InternalError(c, "更新学习记录失败: "+err.Error())
			return
		}
	}

	utils.Success(c, map[string]interface{}{
		"progress":    record.Progress,
		"duration":    record.Duration,
		"isCompleted": record.IsCompleted,
	})
}

// CompleteCourse 标记课程完成
func (ctrl *LearnController) CompleteCourse(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		utils.Unauthorized(c, "未登录")
		return
	}

	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if id == 0 {
		utils.BadRequest(c, "课程ID不能为空")
		return
	}

	var req struct {
		Progress int `json:"progress"` // 可选，默认100
		Duration int `json:"duration"` // 可选
	}

	if err := c.ShouldBindJSON(&req); err == nil {
		// 如果提供了参数，使用提供的值
	} else {
		// 默认值
		req.Progress = 100
	}

	userIDUint := userID.(uint)

	// 查找或创建学习记录
	var record models.CourseRecord
	err := database.DB.Where("user_id = ? AND course_id = ?", userIDUint, id).First(&record).Error

	now := time.Now()
	if err != nil {
		// 创建新记录
		record = models.CourseRecord{
			UserID:      userIDUint,
			CourseID:    uint(id),
			Progress:    req.Progress,
			Duration:    req.Duration,
			IsCompleted: true,
			CompletedAt: &now,
			LastStudyAt: now,
		}
		if err := database.DB.Create(&record).Error; err != nil {
			utils.InternalError(c, "创建学习记录失败: "+err.Error())
			return
		}
	} else {
		// 更新记录
		record.Progress = req.Progress
		if req.Duration > 0 {
			record.Duration = req.Duration
		}
		record.IsCompleted = true
		if record.CompletedAt == nil {
			record.CompletedAt = &now
		}
		record.LastStudyAt = now

		if err := database.DB.Save(&record).Error; err != nil {
			utils.InternalError(c, "更新学习记录失败: "+err.Error())
			return
		}
	}

	utils.Success(c, map[string]interface{}{
		"progress":    record.Progress,
		"duration":    record.Duration,
		"isCompleted": record.IsCompleted,
		"completedAt": record.CompletedAt.Format("2006-01-02 15:04:05"),
	})
}
