package article

import (
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(article *Article) error {
	return r.db.Create(article).Error
}

func (r *Repository) FindByID(id uint) (*Article, error) {
	var a Article
	err := r.db.Preload("User").Preload("Category").First(&a, id).Error
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *Repository) Update(article *Article) error {
	return r.db.Save(article).Error
}

func (r *Repository) Delete(id uint) error {
	return r.db.Delete(&Article{}, id).Error
}

func (r *Repository) List(page, pageSize int, conditions map[string]interface{}) ([]Article, int64, error) {
	var articles []Article
	var total int64

	query := r.db.Model(&Article{}).Preload("User").Preload("Category")

	if status, ok := conditions["status"]; ok {
		query = query.Where("status = ?", status)
	}
	if categoryID, ok := conditions["category_id"]; ok {
		query = query.Where("category_id = ?", categoryID)
	}
	if userID, ok := conditions["user_id"]; ok {
		query = query.Where("user_id = ?", userID)
	}

	query.Count(&total)
	err := query.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&articles).Error
	if err != nil {
		return nil, 0, err
	}
	return articles, total, nil
}

func (r *Repository) SearchTitle(keyword string, page, pageSize int) ([]Article, int64, error) {
	var articles []Article
	var total int64

	query := r.db.Model(&Article{}).Preload("User").Preload("Category").
		Where("title LIKE ? AND status = ?", "%"+keyword+"%", "published")

	query.Count(&total)
	err := query.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&articles).Error
	if err != nil {
		return nil, 0, err
	}
	return articles, total, nil
}

func (r *Repository) IncrementViewCount(id uint) error {
	return r.db.Model(&Article{}).Where("id = ?", id).UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error
}

func (r *Repository) IncrementLikeCount(id uint) error {
	return r.db.Model(&Article{}).Where("id = ?", id).UpdateColumn("like_count", gorm.Expr("like_count + 1")).Error
}

func (r *Repository) DecrementLikeCount(id uint) error {
	return r.db.Model(&Article{}).Where("id = ?", id).UpdateColumn("like_count", gorm.Expr("like_count - 1")).Error
}

func (r *Repository) IncrementCommentCount(id uint) error {
	return r.db.Model(&Article{}).Where("id = ?", id).UpdateColumn("comment_count", gorm.Expr("comment_count + 1")).Error
}

func (r *Repository) DecrementCommentCount(id uint) error {
	return r.db.Model(&Article{}).Where("id = ?", id).UpdateColumn("comment_count", gorm.Expr("comment_count - 1")).Error
}

func (r *Repository) IncrementFavoriteCount(id uint) error {
	return r.db.Model(&Article{}).Where("id = ?", id).UpdateColumn("favorite_count", gorm.Expr("favorite_count + 1")).Error
}

func (r *Repository) DecrementFavoriteCount(id uint) error {
	return r.db.Model(&Article{}).Where("id = ?", id).UpdateColumn("favorite_count", gorm.Expr("favorite_count - 1")).Error
}

func (r *Repository) GetHot(limit int) ([]Article, error) {
	var articles []Article
	err := r.db.Where("status = ?", "published").
		Preload("User").Preload("Category").
		Order("view_count DESC, like_count DESC").
		Limit(limit).Find(&articles).Error
	return articles, err
}

func (r *Repository) GetUserArticles(userID uint, page, pageSize int, status string) ([]Article, int64, error) {
	var articles []Article
	var total int64

	query := r.db.Model(&Article{}).Where("user_id = ?", userID)
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)
	err := query.Preload("Category").Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&articles).Error
	return articles, total, err
}

func (r *Repository) IsLiked(userID uint, targetType string, targetID uint) (bool, error) {
	var count int64
	err := r.db.Model(&Like{}).Where("user_id = ? AND target_type = ? AND target_id = ?", userID, targetType, targetID).Count(&count).Error
	return count > 0, err
}

func (r *Repository) CreateLike(like *Like) error {
	return r.db.Create(like).Error
}

func (r *Repository) DeleteLike(userID uint, targetType string, targetID uint) error {
	return r.db.Where("user_id = ? AND target_type = ? AND target_id = ?", userID, targetType, targetID).Delete(&Like{}).Error
}

func (r *Repository) IsFavorited(userID, articleID uint) (bool, error) {
	var count int64
	err := r.db.Model(&Favorite{}).Where("user_id = ? AND article_id = ?", userID, articleID).Count(&count).Error
	return count > 0, err
}

func (r *Repository) CreateFavorite(fav *Favorite) error {
	return r.db.Create(fav).Error
}

func (r *Repository) DeleteFavorite(userID, articleID uint) error {
	return r.db.Where("user_id = ? AND article_id = ?", userID, articleID).Delete(&Favorite{}).Error
}

func (r *Repository) GetUserFavorites(userID uint, page, pageSize int) ([]Article, int64, error) {
	var articles []Article
	var total int64

	subQuery := r.db.Model(&Favorite{}).Where("user_id = ?", userID).Select("article_id")
	r.db.Model(&Article{}).Where("id IN (?) AND deleted_at IS NULL", subQuery).Count(&total)
	err := r.db.Where("id IN (?) AND deleted_at IS NULL", subQuery).
		Preload("User").Preload("Category").
		Order("created_at DESC").
		Offset((page - 1) * pageSize).Limit(pageSize).Find(&articles).Error
	return articles, total, err
}

func (r *Repository) GetFeed(userID uint, page, pageSize int) ([]Article, int64, error) {
	var articles []Article
	var total int64

	r.db.Model(&Article{}).Where("status = ? AND user_id != ?", "published", userID).Count(&total)
	err := r.db.Where("status = ? AND user_id != ?", "published", userID).
		Preload("User").Preload("Category").
		Order("created_at DESC").
		Offset((page - 1) * pageSize).Limit(pageSize).Find(&articles).Error

	return articles, total, err
}
