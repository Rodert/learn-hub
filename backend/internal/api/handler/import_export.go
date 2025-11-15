package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"learn-hub/internal/model"
	"learn-hub/pkg/excel"
)

// ImportExportHandler 数据导入导出处理器
type ImportExportHandler struct {
	db *gorm.DB
}

// NewImportExportHandler 创建数据导入导出处理器
func NewImportExportHandler(db *gorm.DB) *ImportExportHandler {
	return &ImportExportHandler{db: db}
}

// ImportQuestionsRequest 导入题目请求
type ImportQuestionsRequest struct {
	ExamID int64 `json:"exam_id"`
}

// ImportQuestions 导入题目
// @Summary 导入题目
// @Description 从 Excel 文件导入题目到试卷
// @Tags 导入导出
// @Security Bearer
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Excel 文件"
// @Param exam_id formData int true "试卷 ID"
// @Success 200 {object} Response{data=map[string]interface{}}
// @Failure 400 {object} Response
// @Failure 401 {object} Response
// @Router /import/questions [post]
func (h *ImportExportHandler) ImportQuestions(c *gin.Context) {
	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file provided"})
		return
	}

	// 获取 exam_id
	examIDStr := c.PostForm("exam_id")
	examID, err := strconv.ParseInt(examIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid exam_id"})
		return
	}

	// 打开文件
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer src.Close()

	// 读取 Excel
	reader, err := excel.NewExcelReader(src)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Excel file"})
		return
	}
	defer reader.Close()

	// 读取第一个工作表
	sheetNames := reader.GetSheetNames()
	if len(sheetNames) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No sheet found in Excel"})
		return
	}

	rows, err := reader.ReadRows(sheetNames[0])
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(rows) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No data found in Excel"})
		return
	}

	// 跳过表头，导入数据
	successCount := 0
	failureCount := 0
	var errors []string

	for i := 1; i < len(rows); i++ {
		row := rows[i]
		if len(row) < 5 {
			failureCount++
			errors = append(errors, fmt.Sprintf("Row %d: Invalid data", i+1))
			continue
		}

		// 解析行数据
		// 格式: 题型 | 题目内容 | 选项 | 答案 | 分数
		questionType := excel.InterfaceToString(row[0])
		content := excel.InterfaceToString(row[1])
		options := excel.InterfaceToString(row[2])
		answer := excel.InterfaceToString(row[3])
		scoreStr := excel.InterfaceToString(row[4])

		score, err := excel.StringToFloat64(scoreStr)
		if err != nil {
			score = 1.0
		}

		// 创建题目
		question := model.Question{
			ExamID:      &examID,
			Type:        questionType,
			Content:     content,
			Options:     options,
			Answer:      answer,
			Score:       score,
			CreatedBy:   1, // TODO: 从 context 获取当前用户 ID
		}

		if err := h.db.Create(&question).Error; err != nil {
			failureCount++
			errors = append(errors, fmt.Sprintf("Row %d: %s", i+1, err.Error()))
		} else {
			successCount++
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"success_count": successCount,
			"failure_count": failureCount,
			"errors":        errors,
		},
	})
}

// ImportUsersRequest 导入用户请求
type ImportUsersRequest struct {
	RoleID int64 `json:"role_id"`
}

// ImportUsers 导入用户
// @Summary 导入用户
// @Description 从 Excel 文件导入用户
// @Tags 导入导出
// @Security Bearer
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Excel 文件"
// @Param role_id formData int false "角色 ID（可选）"
// @Success 200 {object} Response{data=map[string]interface{}}
// @Failure 400 {object} Response
// @Failure 401 {object} Response
// @Router /import/users [post]
func (h *ImportExportHandler) ImportUsers(c *gin.Context) {
	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file provided"})
		return
	}

	// 获取 role_id（可选）
	roleIDStr := c.PostForm("role_id")
	var roleID int64
	if roleIDStr != "" {
		roleID, _ = strconv.ParseInt(roleIDStr, 10, 64)
	}

	// 打开文件
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer src.Close()

	// 读取 Excel
	reader, err := excel.NewExcelReader(src)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Excel file"})
		return
	}
	defer reader.Close()

	// 读取第一个工作表
	sheetNames := reader.GetSheetNames()
	if len(sheetNames) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No sheet found in Excel"})
		return
	}

	rows, err := reader.ReadRows(sheetNames[0])
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(rows) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No data found in Excel"})
		return
	}

	// 跳过表头，导入数据
	successCount := 0
	failureCount := 0
	var errors []string

	for i := 1; i < len(rows); i++ {
		row := rows[i]
		if len(row) < 3 {
			failureCount++
			errors = append(errors, fmt.Sprintf("Row %d: Invalid data", i+1))
			continue
		}

		// 解析行数据
		// 格式: 用户名 | 昵称 | 密码
		username := excel.InterfaceToString(row[0])
		nickname := excel.InterfaceToString(row[1])
		password := excel.InterfaceToString(row[2])

		// 检查用户是否已存在
		var count int64
		h.db.Model(&model.User{}).Where("username = ?", username).Count(&count)
		if count > 0 {
			failureCount++
			errors = append(errors, fmt.Sprintf("Row %d: User already exists", i+1))
			continue
		}

		// 密码加密
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			failureCount++
			errors = append(errors, fmt.Sprintf("Row %d: %s", i+1, err.Error()))
			continue
		}

		// 创建用户
		user := model.User{
			Username: username,
			Password: string(hashedPassword),
			Nickname: nickname,
			Status:   "active",
		}

		if err := h.db.Create(&user).Error; err != nil {
			failureCount++
			errors = append(errors, fmt.Sprintf("Row %d: %s", i+1, err.Error()))
			continue
		}

		// 分配角色
		if roleID > 0 {
			h.db.Create(&model.UserRole{
				UserID: user.ID,
				RoleID: roleID,
			})
		} else {
			// 分配默认角色（user）
			var userRole model.Role
			h.db.First(&userRole, "name = ?", "user")
			if userRole.ID > 0 {
				h.db.Create(&model.UserRole{
					UserID: user.ID,
					RoleID: userRole.ID,
				})
			}
		}

		successCount++
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"success_count": successCount,
			"failure_count": failureCount,
			"errors":        errors,
		},
	})
}

// ExportExamScores 导出考试成绩
// @Summary 导出考试成绩
// @Description 导出试卷的所有考试成绩到 Excel
// @Tags 导入导出
// @Security Bearer
// @Produce octet-stream
// @Param exam_id query int true "试卷 ID"
// @Success 200 {file} application/octet-stream
// @Failure 400 {object} Response
// @Failure 401 {object} Response
// @Router /export/exam-scores [get]
func (h *ImportExportHandler) ExportExamScores(c *gin.Context) {
	// 获取 exam_id
	examIDStr := c.Query("exam_id")
	examID, err := strconv.ParseInt(examIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid exam_id"})
		return
	}

	// 查询考试记录
	var records []model.ExamRecord
	if err := h.db.Where("exam_id = ?", examID).Find(&records).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 创建 Excel 文件
	writer := excel.NewExcelWriter()
	defer writer.Close()

	// 写入表头
	headers := []string{"用户 ID", "用户名", "成绩", "状态", "开始时间", "提交时间"}
	writer.WriteHeaders("成绩", headers)

	// 写入数据
	for i, record := range records {
		// 查询用户信息
		var user model.User
		h.db.First(&user, record.UserID)

		startTime := ""
		if record.StartTime != nil {
			startTime = record.StartTime.Format("2006-01-02 15:04:05")
		}

		submitTime := ""
		if record.SubmitTime != nil {
			submitTime = record.SubmitTime.Format("2006-01-02 15:04:05")
		}

		score := ""
		if record.Score != nil {
			score = fmt.Sprintf("%.2f", *record.Score)
		}

		row := []interface{}{
			record.UserID,
			user.Username,
			score,
			record.Status,
			startTime,
			submitTime,
		}

		writer.WriteRow("成绩", i+2, row)
	}

	// 设置响应头
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=exam_scores_%d_%d.xlsx", examID, time.Now().Unix()))

	// 写入到响应
	if err := writer.WriteTo(c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}
