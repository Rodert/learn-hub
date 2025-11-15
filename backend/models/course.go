package models

import (
	"time"

	"gorm.io/gorm"
)

// Course 课程模型
type Course struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Title       string `gorm:"size:200;not null;comment:课程标题" json:"title"`
	Description string `gorm:"type:text;comment:课程描述" json:"description"`
	CoverImage  string `gorm:"size:500;comment:封面图片URL" json:"coverImage"`
	ContentType int    `gorm:"default:1;comment:内容类型:1-视频,2-文本,3-混合" json:"contentType"`
	VideoURL    string `gorm:"size:500;comment:视频URL" json:"videoUrl"`
	TextContent string `gorm:"type:text;comment:文本内容" json:"textContent"`
	Duration    int    `gorm:"default:0;comment:视频时长(秒),文本为0" json:"duration"`
	Status      int    `gorm:"default:0;comment:状态:0-草稿,1-已发布,2-已下架" json:"status"`
	SortOrder   int    `gorm:"default:0;comment:排序" json:"sortOrder"`

	// 关联关系（禁用外键约束，避免迁移顺序问题）
	Records []CourseRecord `gorm:"foreignKey:CourseID;references:ID;constraint:-" json:"records,omitempty"`
}

// TableName 指定表名
func (Course) TableName() string {
	return "sys_course"
}
