package models

import (
	"time"

	"gorm.io/gorm"
)

// CourseRecord 学习记录模型
type CourseRecord struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	UserID      uint       `gorm:"not null;index:idx_user_course,unique;comment:用户ID" json:"userId"`
	CourseID    uint       `gorm:"not null;index:idx_user_course,unique;comment:课程ID" json:"courseId"`
	Progress    int        `gorm:"default:0;comment:学习进度(0-100)" json:"progress"`
	Duration    int        `gorm:"default:0;comment:已学习时长(秒)" json:"duration"`
	IsCompleted bool       `gorm:"default:0;comment:是否完成" json:"isCompleted"`
	CompletedAt *time.Time `gorm:"comment:完成时间" json:"completedAt,omitempty"`
	LastStudyAt time.Time  `gorm:"comment:最后学习时间" json:"lastStudyAt"`

	// 关联关系（禁用外键约束，避免迁移顺序问题）
	User   User   `gorm:"foreignKey:UserID;references:ID;constraint:-" json:"user,omitempty"`
	Course Course `gorm:"foreignKey:CourseID;references:ID;constraint:-" json:"course,omitempty"`
}

// TableName 指定表名
func (CourseRecord) TableName() string {
	return "sys_course_record"
}
