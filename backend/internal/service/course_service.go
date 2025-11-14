package service

import (
	"time"

	"learn-hub/internal/model"
	"learn-hub/internal/repository"
)

// CourseService 学习记录服务
type CourseService struct {
	recordRepo *repository.CourseRecordRepository
}

// NewCourseService 创建学习记录服务
func NewCourseService(recordRepo *repository.CourseRecordRepository) *CourseService {
	return &CourseService{
		recordRepo: recordRepo,
	}
}

// StartLearning 开始学习
func (s *CourseService) StartLearning(userID, materialID int64) (*model.CourseRecord, error) {
	// 检查是否已有记录
	record, err := s.recordRepo.GetByUserAndMaterial(userID, materialID)
	if err == nil && record != nil {
		// 如果已有记录，更新状态为 in_progress
		if record.Status == "not_started" {
			record.Status = "in_progress"
			if err := s.recordRepo.Update(record); err != nil {
				return nil, err
			}
		}
		return record, nil
	}

	// 创建新的学习记录
	newRecord := &model.CourseRecord{
		UserID:     userID,
		MaterialID: materialID,
		Status:     "in_progress",
	}

	if err := s.recordRepo.Create(newRecord); err != nil {
		return nil, err
	}

	return newRecord, nil
}

// UpdateProgress 更新学习进度
func (s *CourseService) UpdateProgress(recordID int64, progressPercent int, viewDuration int) (*model.CourseRecord, error) {
	record, err := s.recordRepo.GetByID(recordID)
	if err != nil {
		return nil, err
	}

	record.ProgressPercent = progressPercent
	record.ViewDuration = viewDuration

	// 如果进度达到 100%，标记为完成
	if progressPercent >= 100 {
		record.Status = "completed"
		now := time.Now()
		record.CompletedAt = &now
	}

	if err := s.recordRepo.Update(record); err != nil {
		return nil, err
	}

	return record, nil
}

// CompleteLearning 完成学习
func (s *CourseService) CompleteLearning(recordID int64) (*model.CourseRecord, error) {
	record, err := s.recordRepo.GetByID(recordID)
	if err != nil {
		return nil, err
	}

	record.Status = "completed"
	record.ProgressPercent = 100
	now := time.Now()
	record.CompletedAt = &now

	if err := s.recordRepo.Update(record); err != nil {
		return nil, err
	}

	return record, nil
}

// GetUserProgress 获取用户学习进度统计
func (s *CourseService) GetUserProgress(userID int64) (map[string]interface{}, error) {
	return s.recordRepo.GetUserProgress(userID)
}

// GetUserRecords 获取用户的学习记录
func (s *CourseService) GetUserRecords(userID int64, page, pageSize int) ([]model.CourseRecord, int64, error) {
	return s.recordRepo.ListByUser(userID, page, pageSize)
}
