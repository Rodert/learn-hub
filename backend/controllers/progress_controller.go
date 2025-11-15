package controllers

import (
	"learn-hub-backend/database"
	"learn-hub-backend/models"
	"learn-hub-backend/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProgressController struct{}

func NewProgressController() *ProgressController {
	return &ProgressController{}
}

// GetCourseProgress 查看课程学习进度（管理员）
func (ctrl *ProgressController) GetCourseProgress(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if id == 0 {
		utils.BadRequest(c, "课程ID不能为空")
		return
	}

	current, _ := strconv.Atoi(c.DefaultQuery("current", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	username := c.Query("username")

	if current < 1 {
		current = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	offset := (current - 1) * pageSize

	// 构建查询
	query := database.DB.Model(&models.CourseRecord{}).
		Where("course_id = ?", id)

	if username != "" {
		// 如果有用户名搜索，使用 JOIN
		query = query.Joins("JOIN sys_user ON sys_course_record.user_id = sys_user.id").
			Where("sys_user.username LIKE ?", "%"+username+"%")
	}

	// 预加载用户信息
	query = query.Preload("User")

	// 查询总数
	var total int64
	query.Count(&total)

	// 查询列表
	var records []models.CourseRecord
	query.Offset(offset).Limit(pageSize).Order("last_study_at DESC").Find(&records)

	// 转换为前端格式
	var list []map[string]interface{}
	for _, record := range records {
		item := map[string]interface{}{
			"userId":      record.UserID,
			"username":    record.User.Username,
			"name":        record.User.Name,
			"progress":    record.Progress,
			"duration":    record.Duration,
			"isCompleted": record.IsCompleted,
			"lastStudyAt": record.LastStudyAt.Format("2006-01-02 15:04:05"),
		}

		if record.CompletedAt != nil {
			item["completedAt"] = record.CompletedAt.Format("2006-01-02 15:04:05")
		} else {
			item["completedAt"] = ""
		}

		list = append(list, item)
	}

	c.JSON(200, gin.H{
		"data":    list,
		"total":   total,
		"success": true,
	})
}

// GetUserProgress 查看用户学习进度（管理员）
func (ctrl *ProgressController) GetUserProgress(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if id == 0 {
		utils.BadRequest(c, "用户ID不能为空")
		return
	}

	current, _ := strconv.Atoi(c.DefaultQuery("current", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	if current < 1 {
		current = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	offset := (current - 1) * pageSize

	// 构建查询
	query := database.DB.Model(&models.CourseRecord{}).
		Where("user_id = ?", id).
		Preload("Course")

	// 查询总数
	var total int64
	query.Count(&total)

	// 查询列表
	var records []models.CourseRecord
	query.Offset(offset).Limit(pageSize).Order("last_study_at DESC").Find(&records)

	// 转换为前端格式
	var list []map[string]interface{}
	for _, record := range records {
		item := map[string]interface{}{
			"courseId":    record.CourseID,
			"courseTitle": record.Course.Title,
			"progress":    record.Progress,
			"duration":    record.Duration,
			"isCompleted": record.IsCompleted,
			"lastStudyAt": record.LastStudyAt.Format("2006-01-02 15:04:05"),
		}

		if record.CompletedAt != nil {
			item["completedAt"] = record.CompletedAt.Format("2006-01-02 15:04:05")
		} else {
			item["completedAt"] = ""
		}

		list = append(list, item)
	}

	c.JSON(200, gin.H{
		"data":    list,
		"total":   total,
		"success": true,
	})
}
