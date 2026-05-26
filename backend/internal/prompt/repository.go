package prompt

import (
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(prompt *Prompt) error {
	return r.db.Create(prompt).Error
}

func (r *Repository) FindByID(id uint) (*Prompt, error) {
	var p Prompt
	err := r.db.Preload("User").First(&p, id).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *Repository) Update(prompt *Prompt) error {
	return r.db.Save(prompt).Error
}

func (r *Repository) Delete(id uint) error {
	return r.db.Delete(&Prompt{}, id).Error
}

func (r *Repository) List(page, pageSize int, category, tag string) ([]Prompt, int64, error) {
	var prompts []Prompt
	var total int64

	query := r.db.Model(&Prompt{}).Preload("User").Where("status = ?", 1)

	if category != "" {
		query = query.Where("category = ?", category)
	}
	if tag != "" {
		query = query.Where("tags LIKE ?", "%"+tag+"%")
	}

	query.Count(&total)
	err := query.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&prompts).Error
	if err != nil {
		return nil, 0, err
	}
	return prompts, total, nil
}

func (r *Repository) GetHot(limit int) ([]Prompt, error) {
	var prompts []Prompt
	err := r.db.Where("status = ?", 1).
		Preload("User").
		Order("usage_count DESC, like_count DESC").
		Limit(limit).Find(&prompts).Error
	return prompts, err
}

func (r *Repository) IncrementUsageCount(id uint) error {
	return r.db.Model(&Prompt{}).Where("id = ?", id).UpdateColumn("usage_count", gorm.Expr("usage_count + 1")).Error
}

func (r *Repository) SearchTitle(keyword string, page, pageSize int) ([]Prompt, int64, error) {
	var prompts []Prompt
	var total int64

	query := r.db.Model(&Prompt{}).Preload("User").
		Where("title LIKE ? AND status = ?", "%"+keyword+"%", 1)

	query.Count(&total)
	err := query.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&prompts).Error
	return prompts, total, err
}
