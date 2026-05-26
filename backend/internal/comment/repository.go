package comment

import (
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(comment *Comment) error {
	return r.db.Create(comment).Error
}

func (r *Repository) FindByID(id uint) (*Comment, error) {
	var c Comment
	err := r.db.Preload("User").First(&c, id).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *Repository) GetByArticleID(articleID uint, page, pageSize int) ([]Comment, int64, error) {
	var comments []Comment
	var total int64

	r.db.Model(&Comment{}).Where("article_id = ?", articleID).Count(&total)
	err := r.db.Where("article_id = ?", articleID).
		Preload("User").
		Order("created_at DESC").
		Offset((page - 1) * pageSize).Limit(pageSize).Find(&comments).Error
	if err != nil {
		return nil, 0, err
	}
	return comments, total, nil
}

func (r *Repository) Delete(id uint) error {
	return r.db.Delete(&Comment{}, id).Error
}

func (r *Repository) CountByArticleID(articleID uint) (int64, error) {
	var count int64
	err := r.db.Model(&Comment{}).Where("article_id = ?", articleID).Count(&count).Error
	return count, err
}
