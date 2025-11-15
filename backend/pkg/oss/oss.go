package oss

import (
	"fmt"
	"io"
	"time"

	"learn-hub/config"
)

// OSSClient OSS 客户端接口
type OSSClient interface {
	// Upload 上传文件
	Upload(key string, reader io.Reader) (string, error)
	// Delete 删除文件
	Delete(key string) error
	// GetPresignedURL 获取预签名 URL
	GetPresignedURL(key string, expiration time.Duration) (string, error)
}

// NewOSSClient 创建 OSS 客户端
func NewOSSClient(cfg config.OSSConfig) (OSSClient, error) {
	switch cfg.Provider {
	case "aliyun":
		return NewAliyunOSSClient(cfg)
	case "tencent":
		return NewTencentOSSClient(cfg)
	case "local":
		// 本地存储（用于开发环境）
		uploadDir := cfg.Endpoint // 使用 endpoint 作为本地上传目录
		if uploadDir == "" {
			uploadDir = "./uploads"
		}
		return NewLocalOSSClient(uploadDir), nil
	default:
		return nil, fmt.Errorf("unsupported OSS provider: %s", cfg.Provider)
	}
}

// AliyunOSSClient 阿里云 OSS 客户端
type AliyunOSSClient struct {
	cfg config.OSSConfig
}

// NewAliyunOSSClient 创建阿里云 OSS 客户端
func NewAliyunOSSClient(cfg config.OSSConfig) (*AliyunOSSClient, error) {
	// TODO: 集成阿里云 OSS SDK
	// 需要添加依赖: github.com/aliyun/aliyun-oss-go-sdk
	return &AliyunOSSClient{cfg: cfg}, nil
}

// Upload 上传文件到阿里云 OSS
func (c *AliyunOSSClient) Upload(key string, reader io.Reader) (string, error) {
	// TODO: 实现上传逻辑
	// 示例: client.PutObject(bucket, key, reader)
	url := fmt.Sprintf("https://%s.oss-%s.aliyuncs.com/%s", c.cfg.Bucket, c.cfg.Region, key)
	return url, nil
}

// Delete 从阿里云 OSS 删除文件
func (c *AliyunOSSClient) Delete(key string) error {
	// TODO: 实现删除逻辑
	// 示例: client.DeleteObject(bucket, key)
	return nil
}

// GetPresignedURL 获取阿里云 OSS 预签名 URL
func (c *AliyunOSSClient) GetPresignedURL(key string, expiration time.Duration) (string, error) {
	// TODO: 实现预签名 URL 生成逻辑
	url := fmt.Sprintf("https://%s.oss-%s.aliyuncs.com/%s", c.cfg.Bucket, c.cfg.Region, key)
	return url, nil
}

// TencentOSSClient 腾讯云 COS 客户端
type TencentOSSClient struct {
	cfg config.OSSConfig
}

// NewTencentOSSClient 创建腾讯云 COS 客户端
func NewTencentOSSClient(cfg config.OSSConfig) (*TencentOSSClient, error) {
	// TODO: 集成腾讯云 COS SDK
	// 需要添加依赖: github.com/tencentyun/cos-go-sdk-v5
	return &TencentOSSClient{cfg: cfg}, nil
}

// Upload 上传文件到腾讯云 COS
func (c *TencentOSSClient) Upload(key string, reader io.Reader) (string, error) {
	// TODO: 实现上传逻辑
	url := fmt.Sprintf("https://%s.cos.%s.myqcloud.com/%s", c.cfg.Bucket, c.cfg.Region, key)
	return url, nil
}

// Delete 从腾讯云 COS 删除文件
func (c *TencentOSSClient) Delete(key string) error {
	// TODO: 实现删除逻辑
	return nil
}

// GetPresignedURL 获取腾讯云 COS 预签名 URL
func (c *TencentOSSClient) GetPresignedURL(key string, expiration time.Duration) (string, error) {
	// TODO: 实现预签名 URL 生成逻辑
	url := fmt.Sprintf("https://%s.cos.%s.myqcloud.com/%s", c.cfg.Bucket, c.cfg.Region, key)
	return url, nil
}
