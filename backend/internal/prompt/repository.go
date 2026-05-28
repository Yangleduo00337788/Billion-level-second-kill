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

func (r *Repository) List(page, pageSize int, category, tag string, userID uint, keyword string) ([]Prompt, int64, error) {
	var prompts []Prompt
	var total int64

	query := r.db.Model(&Prompt{}).Preload("User").Where("status = ?", 1)

	if category != "" {
		query = query.Where("category = ?", category)
	}
	if tag != "" {
		query = query.Where("tags LIKE ?", "%"+tag+"%")
	}
	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}
	if keyword != "" {
		query = query.Where("title LIKE ?", "%"+keyword+"%")
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

func (r *Repository) IncrementLikeCount(id uint) error {
	return r.db.Model(&Prompt{}).Where("id = ?", id).UpdateColumn("like_count", gorm.Expr("like_count + 1")).Error
}

func (r *Repository) DecrementLikeCount(id uint) error {
	return r.db.Model(&Prompt{}).Where("id = ?", id).UpdateColumn("like_count", gorm.Expr("like_count - 1")).Error
}

func (r *Repository) IncrementFavoriteCount(id uint) error {
	return r.db.Model(&Prompt{}).Where("id = ?", id).UpdateColumn("favorite_count", gorm.Expr("favorite_count + 1")).Error
}

func (r *Repository) DecrementFavoriteCount(id uint) error {
	return r.db.Model(&Prompt{}).Where("id = ?", id).UpdateColumn("favorite_count", gorm.Expr("favorite_count - 1")).Error
}

func (r *Repository) IsLiked(userID, promptID uint) (bool, error) {
	var count int64
	err := r.db.Model(&PromptLike{}).Where("user_id = ? AND prompt_id = ?", userID, promptID).Count(&count).Error
	return count > 0, err
}

func (r *Repository) CreateLike(like *PromptLike) error {
	return r.db.Create(like).Error
}

func (r *Repository) DeleteLike(userID, promptID uint) error {
	return r.db.Unscoped().Where("user_id = ? AND prompt_id = ?", userID, promptID).Delete(&PromptLike{}).Error
}

func (r *Repository) IsFavorited(userID, promptID uint) (bool, error) {
	var count int64
	err := r.db.Model(&PromptFavorite{}).Where("user_id = ? AND prompt_id = ?", userID, promptID).Count(&count).Error
	return count > 0, err
}

func (r *Repository) CreateFavorite(fav *PromptFavorite) error {
	return r.db.Create(fav).Error
}

func (r *Repository) DeleteFavorite(userID, promptID uint) error {
	return r.db.Unscoped().Where("user_id = ? AND prompt_id = ?", userID, promptID).Delete(&PromptFavorite{}).Error
}

