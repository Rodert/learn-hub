package repository

import (
	"gorm.io/gorm"
	"learn-hub/internal/model"
)

// QuestionRepository 题库仓储
type QuestionRepository struct {
	db *gorm.DB
}

// NewQuestionRepository 创建题库仓储
func NewQuestionRepository(db *gorm.DB) *QuestionRepository {
	return &QuestionRepository{db: db}
}

// GetByID 根据 ID 获取题目
func (r *QuestionRepository) GetByID(id int64) (*model.Question, error) {
	var question model.Question
	if err := r.db.First(&question, id).Error; err != nil {
		return nil, err
	}
	return &question, nil
}

// Create 创建题目
func (r *QuestionRepository) Create(question *model.Question) error {
	return r.db.Create(question).Error
}

// Update 更新题目
func (r *QuestionRepository) Update(question *model.Question) error {
	return r.db.Save(question).Error
}

// Delete 删除题目（软删除）
func (r *QuestionRepository) Delete(id int64) error {
	return r.db.Delete(&model.Question{}, id).Error
}

// ListByExam 获取试卷的题目列表
func (r *QuestionRepository) ListByExam(examID int64) ([]model.Question, error) {
	var questions []model.Question
	if err := r.db.Where("exam_id = ?", examID).
		Order("order_num ASC").
		Find(&questions).Error; err != nil {
		return nil, err
	}
	return questions, nil
}

// List 获取题目列表
func (r *QuestionRepository) List(page, pageSize int, questionType string) ([]model.Question, int64, error) {
	var questions []model.Question
	var total int64

	query := r.db.Model(&model.Question{})
	if questionType != "" {
		query = query.Where("type = ?", questionType)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&questions).Error; err != nil {
		return nil, 0, err
	}

	return questions, total, nil
}

// BatchCreate 批量创建题目
func (r *QuestionRepository) BatchCreate(questions []model.Question) error {
	return r.db.CreateInBatches(questions, 100).Error
}
