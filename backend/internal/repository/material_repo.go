package repository

import (
	"gorm.io/gorm"
	"learn-hub/internal/model"
)

// MaterialRepository 资料仓储
type MaterialRepository struct {
	db *gorm.DB
}

// NewMaterialRepository 创建资料仓储
func NewMaterialRepository(db *gorm.DB) *MaterialRepository {
	return &MaterialRepository{db: db}
}

// GetByID 根据 ID 获取资料
func (r *MaterialRepository) GetByID(id int64) (*model.Material, error) {
	var material model.Material
	if err := r.db.First(&material, id).Error; err != nil {
		return nil, err
	}
	return &material, nil
}

// Create 创建资料
func (r *MaterialRepository) Create(material *model.Material) error {
	return r.db.Create(material).Error
}

// Update 更新资料
func (r *MaterialRepository) Update(material *model.Material) error {
	return r.db.Save(material).Error
}

// Delete 删除资料（软删除）
func (r *MaterialRepository) Delete(id int64) error {
	return r.db.Delete(&model.Material{}, id).Error
}

// List 获取资料列表
func (r *MaterialRepository) List(page, pageSize int, status string) ([]model.Material, int64, error) {
	var materials []model.Material
	var total int64

	query := r.db.Model(&model.Material{})
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("order_num ASC, created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&materials).Error; err != nil {
		return nil, 0, err
	}

	return materials, total, nil
}

// UpdateStatus 更新资料状态
func (r *MaterialRepository) UpdateStatus(id int64, status string) error {
	return r.db.Model(&model.Material{}).Where("id = ?", id).Update("status", status).Error
}
