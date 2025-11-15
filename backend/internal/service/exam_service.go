package service

import (
	"encoding/json"
	"errors"
	"time"

	"learn-hub/internal/model"
	"learn-hub/internal/repository"
)

// AnswerItem 答题项
type AnswerItem struct {
	QuestionID int64  `json:"question_id"`
	Answer     string `json:"answer"`
}

// ExamService 考试服务
type ExamService struct {
	examRepo   *repository.ExamRepository
	recordRepo *repository.ExamRecordRepository
	qRepo      *repository.QuestionRepository
}

// NewExamService 创建考试服务
func NewExamService(
	examRepo *repository.ExamRepository,
	recordRepo *repository.ExamRecordRepository,
	qRepo *repository.QuestionRepository,
) *ExamService {
	return &ExamService{
		examRepo:   examRepo,
		recordRepo: recordRepo,
		qRepo:      qRepo,
	}
}

// StartExam 开始考试
func (s *ExamService) StartExam(userID, examID int64) (*model.ExamRecord, error) {
	// 检查是否已有进行中的记录
	record, err := s.recordRepo.GetByUserAndExam(userID, examID)
	if err == nil && record != nil && record.Status == "in_progress" {
		return record, nil
	}

	// 创建新的考试记录
	now := time.Now()
	newRecord := &model.ExamRecord{
		UserID:    userID,
		ExamID:    examID,
		Status:    "in_progress",
		StartTime: &now,
		Answers:   "[]",
	}

	if err := s.recordRepo.Create(newRecord); err != nil {
		return nil, err
	}

	return newRecord, nil
}

// SubmitExam 提交答卷
func (s *ExamService) SubmitExam(recordID int64, answers []AnswerItem) (float64, error) {
	// 获取考试记录
	record, err := s.recordRepo.GetByID(recordID)
	if err != nil {
		return 0, err
	}

	if record.Status != "in_progress" {
		return 0, errors.New("exam not in progress")
	}

	// 获取所有题目
	questions, err := s.qRepo.ListByExam(record.ExamID)
	if err != nil {
		return 0, err
	}

	// 自动评分
	totalScore := 0.0
	correctCount := 0

	for _, answer := range answers {
		// 查找对应的题目
		var question *model.Question
		for i := range questions {
			if questions[i].ID == answer.QuestionID {
				question = &questions[i]
				break
			}
		}

		if question == nil {
			continue
		}

		// 判断答案是否正确
		if s.isAnswerCorrect(question, answer.Answer) {
			totalScore += question.Score
			correctCount++
		}
	}

	// 更新考试记录
	now := time.Now()
	record.Status = "graded"
	record.SubmitTime = &now
	record.Score = &totalScore
	record.Answers = string(mustMarshal(answers))

	if err := s.recordRepo.Update(record); err != nil {
		return 0, err
	}

	return totalScore, nil
}

// isAnswerCorrect 判断答案是否正确
func (s *ExamService) isAnswerCorrect(question *model.Question, userAnswer string) bool {
	switch question.Type {
	case "single_choice", "multiple_choice":
		// 单选和多选直接比较
		return question.Answer == userAnswer
	case "fill_blank":
		// 填空题支持模糊匹配
		return s.fuzzyMatch(question.Answer, userAnswer)
	default:
		return false
	}
}

// fuzzyMatch 填空题模糊匹配
func (s *ExamService) fuzzyMatch(correctAnswer, userAnswer string) bool {
	// 简单的模糊匹配：去除空格后比较
	// 可以根据需要添加更复杂的匹配逻辑
	return correctAnswer == userAnswer
}

// GetExamRecords 获取用户的考试成绩
func (s *ExamService) GetExamRecords(userID, examID int64, page, pageSize int) ([]model.ExamRecord, int64, error) {
	return s.recordRepo.ListByExam(examID, page, pageSize)
}

// GetUserExamRecords 获取用户的所有考试记录
func (s *ExamService) GetUserExamRecords(userID int64, page, pageSize int) ([]model.ExamRecord, int64, error) {
	return s.recordRepo.ListByUser(userID, page, pageSize)
}

// mustMarshal 必须序列化
func mustMarshal(v interface{}) []byte {
	data, _ := json.Marshal(v)
	return data
}
