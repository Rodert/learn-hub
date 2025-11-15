package handler

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"learn-hub/config"
	"learn-hub/pkg/oss"
)

// FileHandler 文件处理器
type FileHandler struct {
	db        *gorm.DB
	cfg       *config.Config
	ossClient oss.OSSClient
}

// NewFileHandler 创建文件处理器
func NewFileHandler(db *gorm.DB, cfg *config.Config, ossClient oss.OSSClient) *FileHandler {
	return &FileHandler{
		db:        db,
		cfg:       cfg,
		ossClient: ossClient,
	}
}

// UploadRequest 上传请求
type UploadRequest struct {
	File     string `form:"file" binding:"required"`
	FileType string `form:"file_type"` // 文件类型：material, question, etc.
}

// UploadResponse 上传响应
type UploadResponse struct {
	URL      string `json:"url"`
	FileName string `json:"file_name"`
	FileSize int64  `json:"file_size"`
}

// Upload 上传文件
// @Summary 上传文件
// @Description 上传文件到 OSS
// @Tags 文件
// @Security Bearer
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "文件"
// @Param file_type formData string false "文件类型"
// @Success 200 {object} Response{data=UploadResponse}
// @Failure 400 {object} Response
// @Failure 401 {object} Response
// @Router /files/upload [post]
func (h *FileHandler) Upload(c *gin.Context) {
	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file provided"})
		return
	}

	// 验证文件大小（最大 100MB）
	maxSize := int64(100 * 1024 * 1024)
	if file.Size > maxSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File too large"})
		return
	}

	// 获取文件内容
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer src.Close()

	// 生成文件名
	timestamp := time.Now().Unix()
	ext := filepath.Ext(file.Filename)
	fileName := fmt.Sprintf("%d_%s%s", timestamp, file.Filename[:len(file.Filename)-len(ext)], ext)

	// 上传到 OSS
	url, err := h.ossClient.Upload(fileName, src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": UploadResponse{
			URL:      url,
			FileName: file.Filename,
			FileSize: file.Size,
		},
	})
}

// DeleteRequest 删除请求
type DeleteRequest struct {
	URL string `json:"url" binding:"required"`
}

// Delete 删除文件
// @Summary 删除文件
// @Description 从 OSS 删除文件
// @Tags 文件
// @Security Bearer
// @Accept json
// @Produce json
// @Param request body DeleteRequest true "删除请求"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 401 {object} Response
// @Router /files/delete [post]
func (h *FileHandler) Delete(c *gin.Context) {
	var req DeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 从 URL 中提取文件名
	// 假设 URL 格式为: https://bucket.oss-region.aliyuncs.com/filename
	// 或: https://bucket.cos.region.myqcloud.com/filename
	// 这里简化处理，实际应该根据 OSS 提供商来解析
	fileName := filepath.Base(req.URL)

	// 从 OSS 删除文件
	if err := h.ossClient.Delete(fileName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}

// GetPresignedURLRequest 获取预签名 URL 请求
type GetPresignedURLRequest struct {
	URL        string `json:"url" binding:"required"`
	Expiration int    `json:"expiration"` // 过期时间（秒），默认 3600
}

// GetPresignedURLResponse 获取预签名 URL 响应
type GetPresignedURLResponse struct {
	PresignedURL string `json:"presigned_url"`
	Expiration   int64  `json:"expiration"`
}

// GetPresignedURL 获取预签名 URL
// @Summary 获取预签名 URL
// @Description 获取文件的预签名 URL（用于直接下载或上传）
// @Tags 文件
// @Security Bearer
// @Accept json
// @Produce json
// @Param request body GetPresignedURLRequest true "获取预签名 URL 请求"
// @Success 200 {object} Response{data=GetPresignedURLResponse}
// @Failure 400 {object} Response
// @Failure 401 {object} Response
// @Router /files/presigned-url [post]
func (h *FileHandler) GetPresignedURL(c *gin.Context) {
	var req GetPresignedURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 设置默认过期时间
	if req.Expiration == 0 {
		req.Expiration = 3600
	}

	// 从 URL 中提取文件名
	fileName := filepath.Base(req.URL)

	// 获取预签名 URL
	presignedURL, err := h.ossClient.GetPresignedURL(fileName, time.Duration(req.Expiration)*time.Second)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": GetPresignedURLResponse{
			PresignedURL: presignedURL,
			Expiration:   time.Now().Add(time.Duration(req.Expiration) * time.Second).Unix(),
		},
	})
}
