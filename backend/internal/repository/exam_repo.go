package repository

import (
	"gorm.io/gorm"
	"learn-hub/internal/model"
)

// ExamRepository 考试仓储
type ExamRepository struct {
	db *gorm.DB
}

// NewExamRepository 创建考试仓储
func NewExamRepository(db *gorm.DB) *ExamRepository {
	return &ExamRepository{db: db}
}

// GetByID 根据 ID 获取试卷
func (r *ExamRepository) GetByID(id int64) (*model.Exam, error) {
	var exam model.Exam
	if err := r.db.Preload("Questions").First(&exam, id).Error; err != nil {
		return nil, err
	}
	return &exam, nil
}

// Create 创建试卷
func (r *ExamRepository) Create(exam *model.Exam) error {
	return r.db.Create(exam).Error
}

// Update 更新试卷
func (r *ExamRepository) Update(exam *model.Exam) error {
	return r.db.Save(exam).Error
}

// Delete 删除试卷（软删除）
func (r *ExamRepository) Delete(id int64) error {
	return r.db.Delete(&model.Exam{}, id).Error
}

// List 获取试卷列表
func (r *ExamRepository) List(page, pageSize int, status string) ([]model.Exam, int64, error) {
	var exams []model.Exam
	var total int64

	query := r.db.Model(&model.Exam{})
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&exams).Error; err != nil {
		return nil, 0, err
	}

	return exams, total, nil
}

// UpdateStatus 更新试卷状态
func (r *ExamRepository) UpdateStatus(id int64, status string) error {
	return r.db.Model(&model.Exam{}).Where("id = ?", id).Update("status", status).Error
}

// GetQuestions 获取试卷的题目
func (r *ExamRepository) GetQuestions(examID int64) ([]model.Question, error) {
	var questions []model.Question
	if err := r.db.Where("exam_id = ?", examID).Order("order_num ASC").Find(&questions).Error; err != nil {
		return nil, err
	}
	return questions, nil
}
