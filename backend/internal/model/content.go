package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Material 学习资料表
type Material struct {
	ID        int64          `gorm:"primaryKey" json:"id"`
	Title     string         `gorm:"size:255" json:"title"`
	Desc      string         `gorm:"type:text" json:"description"`
	Type      string         `gorm:"type:enum('text','video','file','mixed');default:'text'" json:"content_type"`
	Content   string         `gorm:"type:longtext" json:"content"`
	FileURL   string         `gorm:"size:500" json:"file_url"`
	FileSize  int64          `json:"file_size"`
	CoverURL  string         `gorm:"size:500" json:"cover_url"`
	OrderNum  int            `gorm:"default:0" json:"order_num"`
	Status    string         `gorm:"type:enum('draft','published','archived');default:'draft'" json:"status"`
	CreatedBy int64          `json:"created_by"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Question 题库表
type Question struct {
	ID          int64          `gorm:"primaryKey" json:"id"`
	ExamID      *int64         `gorm:"index" json:"exam_id"`
	Type        string         `gorm:"type:enum('single_choice','multiple_choice','fill_blank')" json:"question_type"`
	Content     string         `gorm:"type:text" json:"content"`
	Options     string         `gorm:"type:json" json:"options"`
	Answer      string         `gorm:"size:500" json:"answer"`
	Explanation string         `gorm:"type:text" json:"explanation"`
	Score       float64        `gorm:"default:1.00" json:"score"`
	OrderNum    int            `gorm:"default:0" json:"order_num"`
	CreatedBy   int64          `json:"created_by"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// Exam 试卷表
type Exam struct {
	ID          int64          `gorm:"primaryKey" json:"id"`
	Title       string         `gorm:"size:255" json:"title"`
	Desc        string         `gorm:"type:text" json:"description"`
	TotalScore  float64        `gorm:"default:100.00" json:"total_score"`
	PassScore   float64        `gorm:"default:60.00" json:"pass_score"`
	TimeLimit   int            `json:"time_limit"`
	Status      string         `gorm:"type:enum('draft','published','archived');default:'draft'" json:"status"`
	CreatedBy   int64          `json:"created_by"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联
	Questions []Question `gorm:"foreignKey:ExamID" json:"questions,omitempty"`
}

// ExamRecord 用户考试记录表
type ExamRecord struct {
	ID         int64          `gorm:"primaryKey" json:"id"`
	UserID     int64          `gorm:"index" json:"user_id"`
	ExamID     int64          `gorm:"index" json:"exam_id"`
	Score      *float64       `json:"score"`
	Status     string         `gorm:"type:enum('in_progress','submitted','graded');default:'in_progress'" json:"status"`
	Answers    string         `gorm:"type:json" json:"answers"`
	StartTime  *time.Time     `json:"start_time"`
	SubmitTime *time.Time     `json:"submit_time"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

// CourseRecord 用户学习记录表
type CourseRecord struct {
	ID              int64          `gorm:"primaryKey" json:"id"`
	UserID          int64          `gorm:"index" json:"user_id"`
	MaterialID      int64          `gorm:"index" json:"material_id"`
	Status          string         `gorm:"type:enum('not_started','in_progress','completed');default:'not_started'" json:"status"`
	ProgressPercent int            `gorm:"default:0" json:"progress_percent"`
	ViewDuration    int            `json:"view_duration"`
	CompletedAt     *time.Time     `json:"completed_at"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

// Topic 学习专题表
type Topic struct {
	ID        int64          `gorm:"primaryKey" json:"id"`
	Title     string         `gorm:"size:255" json:"title"`
	Desc      string         `gorm:"type:text" json:"description"`
	OrderNum  int            `gorm:"default:0" json:"order_num"`
	Status    string         `gorm:"type:enum('draft','published','archived');default:'draft'" json:"status"`
	CreatedBy int64          `json:"created_by"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联
	Materials []Material `gorm:"many2many:topic_materials;" json:"materials,omitempty"`
	Exams     []Exam     `gorm:"many2many:topic_exams;" json:"exams,omitempty"`
}

// TopicMaterial 专题资料关联表
type TopicMaterial struct {
	ID         int64          `gorm:"primaryKey" json:"id"`
	TopicID    int64          `gorm:"index" json:"topic_id"`
	MaterialID int64          `gorm:"index" json:"material_id"`
	OrderNum   int            `gorm:"default:0" json:"order_num"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

// TopicExam 专题考试关联表
type TopicExam struct {
	ID        int64          `gorm:"primaryKey" json:"id"`
	TopicID   int64          `gorm:"index" json:"topic_id"`
	ExamID    int64          `gorm:"index" json:"exam_id"`
	OrderNum  int            `gorm:"default:0" json:"order_num"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (Material) TableName() string {
	return "materials"
}

func (Question) TableName() string {
	return "questions"
}

func (Exam) TableName() string {
	return "exams"
}

func (ExamRecord) TableName() string {
	return "exam_records"
}

func (CourseRecord) TableName() string {
	return "course_records"
}

func (Topic) TableName() string {
	return "topics"
}

func (TopicMaterial) TableName() string {
	return "topic_materials"
}

func (TopicExam) TableName() string {
	return "topic_exams"
}
