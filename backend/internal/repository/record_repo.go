package repository

import (
	"gorm.io/gorm"
	"learn-hub/internal/model"
)

// ExamRecordRepository 考试记录仓储
type ExamRecordRepository struct {
	db *gorm.DB
}

// NewExamRecordRepository 创建考试记录仓储
func NewExamRecordRepository(db *gorm.DB) *ExamRecordRepository {
	return &ExamRecordRepository{db: db}
}

// GetByID 根据 ID 获取考试记录
func (r *ExamRecordRepository) GetByID(id int64) (*model.ExamRecord, error) {
	var record model.ExamRecord
	if err := r.db.First(&record, id).Error; err != nil {
		return nil, err
	}
	return &record, nil
}

// Create 创建考试记录
func (r *ExamRecordRepository) Create(record *model.ExamRecord) error {
	return r.db.Create(record).Error
}

// Update 更新考试记录
func (r *ExamRecordRepository) Update(record *model.ExamRecord) error {
	return r.db.Save(record).Error
}

// ListByUser 获取用户的考试记录
func (r *ExamRecordRepository) ListByUser(userID int64, page, pageSize int) ([]model.ExamRecord, int64, error) {
	var records []model.ExamRecord
	var total int64

	if err := r.db.Model(&model.ExamRecord{}).
		Where("user_id = ?", userID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&records).Error; err != nil {
		return nil, 0, err
	}

	return records, total, nil
}

// ListByExam 获取试卷的考试记录
func (r *ExamRecordRepository) ListByExam(examID int64, page, pageSize int) ([]model.ExamRecord, int64, error) {
	var records []model.ExamRecord
	var total int64

	if err := r.db.Model(&model.ExamRecord{}).
		Where("exam_id = ?", examID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := r.db.Where("exam_id = ?", examID).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&records).Error; err != nil {
		return nil, 0, err
	}

	return records, total, nil
}

// GetByUserAndExam 获取用户在某个试卷的记录
func (r *ExamRecordRepository) GetByUserAndExam(userID, examID int64) (*model.ExamRecord, error) {
	var record model.ExamRecord
	if err := r.db.Where("user_id = ? AND exam_id = ?", userID, examID).
		Order("created_at DESC").
		First(&record).Error; err != nil {
		return nil, err
	}
	return &record, nil
}

// CourseRecordRepository 学习记录仓储
type CourseRecordRepository struct {
	db *gorm.DB
}

// NewCourseRecordRepository 创建学习记录仓储
func NewCourseRecordRepository(db *gorm.DB) *CourseRecordRepository {
	return &CourseRecordRepository{db: db}
}

// GetByID 根据 ID 获取学习记录
func (r *CourseRecordRepository) GetByID(id int64) (*model.CourseRecord, error) {
	var record model.CourseRecord
	if err := r.db.First(&record, id).Error; err != nil {
		return nil, err
	}
	return &record, nil
}

// Create 创建学习记录
func (r *CourseRecordRepository) Create(record *model.CourseRecord) error {
	return r.db.Create(record).Error
}

// Update 更新学习记录
func (r *CourseRecordRepository) Update(record *model.CourseRecord) error {
	return r.db.Save(record).Error
}

// ListByUser 获取用户的学习记录
func (r *CourseRecordRepository) ListByUser(userID int64, page, pageSize int) ([]model.CourseRecord, int64, error) {
	var records []model.CourseRecord
	var total int64

	if err := r.db.Model(&model.CourseRecord{}).
		Where("user_id = ?", userID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&records).Error; err != nil {
		return nil, 0, err
	}

	return records, total, nil
}

// GetByUserAndMaterial 获取用户在某个资料的学习记录
func (r *CourseRecordRepository) GetByUserAndMaterial(userID, materialID int64) (*model.CourseRecord, error) {
	var record model.CourseRecord
	if err := r.db.Where("user_id = ? AND material_id = ?", userID, materialID).
		First(&record).Error; err != nil {
		return nil, err
	}
	return &record, nil
}

// GetUserProgress 获取用户的学习进度统计
func (r *CourseRecordRepository) GetUserProgress(userID int64) (map[string]interface{}, error) {
	var stats struct {
		Total     int64
		Completed int64
		InProgress int64
	}

	if err := r.db.Model(&model.CourseRecord{}).
		Where("user_id = ?", userID).
		Select("COUNT(*) as total, SUM(CASE WHEN status = 'completed' THEN 1 ELSE 0 END) as completed, SUM(CASE WHEN status = 'in_progress' THEN 1 ELSE 0 END) as in_progress").
		Scan(&stats).Error; err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total":       stats.Total,
		"completed":   stats.Completed,
		"in_progress": stats.InProgress,
		"progress_percent": func() int64 {
			if stats.Total == 0 {
				return 0
			}
			return (stats.Completed * 100) / stats.Total
		}(),
	}, nil
}
