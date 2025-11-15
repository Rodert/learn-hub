package oss

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// LocalOSSClient 本地文件存储客户端（用于开发环境）
type LocalOSSClient struct {
	uploadDir string
}

// NewLocalOSSClient 创建本地文件存储客户端
func NewLocalOSSClient(uploadDir string) *LocalOSSClient {
	// 确保上传目录存在
	os.MkdirAll(uploadDir, 0755)
	return &LocalOSSClient{uploadDir: uploadDir}
}

// Upload 上传文件到本地存储
func (c *LocalOSSClient) Upload(key string, reader io.Reader) (string, error) {
	// 创建文件路径
	filePath := filepath.Join(c.uploadDir, key)

	// 确保目录存在
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	// 创建文件
	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// 复制内容
	if _, err := io.Copy(file, reader); err != nil {
		return "", err
	}

	// 返回访问 URL（相对路径）
	url := fmt.Sprintf("/uploads/%s", key)
	return url, nil
}

// Delete 从本地存储删除文件
func (c *LocalOSSClient) Delete(key string) error {
	filePath := filepath.Join(c.uploadDir, key)
	return os.Remove(filePath)
}

// GetPresignedURL 获取本地存储的预签名 URL
func (c *LocalOSSClient) GetPresignedURL(key string, expiration time.Duration) (string, error) {
	// 本地存储不需要预签名，直接返回访问 URL
	url := fmt.Sprintf("/uploads/%s", key)
	return url, nil
}
