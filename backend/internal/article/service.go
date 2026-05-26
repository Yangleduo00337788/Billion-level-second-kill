package article

import (
	"errors"

	"inference-engine/internal/ai"
	"inference-engine/internal/user"

	"gorm.io/gorm"
)

type Service struct {
	repo      *Repository
	userRepo  *user.Repository
	aiService *ai.Service
	db        *gorm.DB
}

func NewService(repo *Repository, userRepo *user.Repository, aiService *ai.Service, db *gorm.DB) *Service {
	return &Service{repo: repo, userRepo: userRepo, aiService: aiService, db: db}
}

type CreateArticleReq struct {
	Title      string `json:"title" binding:"required,max=200"`
	Content    string `json:"content"`
	Summary    string `json:"summary"`
	Cover      string `json:"cover"`
	Status     string `json:"status"`
	CategoryID uint   `json:"category_id"`
	Tags       string `json:"tags"`
}

type UpdateArticleReq struct {
	Title      string `json:"title"`
	Content    string `json:"content"`
	Summary    string `json:"summary"`
	Cover      string `json:"cover"`
	Status     string `json:"status"`
	CategoryID uint   `json:"category_id"`
	Tags       string `json:"tags"`
}

func (s *Service) Create(userID uint, req *CreateArticleReq) (*Article, error) {
	article := &Article{
		UserID:     userID,
		Title:      req.Title,
		Content:    req.Content,
		Summary:    req.Summary,
		Cover:      req.Cover,
		Status:     req.Status,
		CategoryID: req.CategoryID,
		Tags:       req.Tags,
		IsAI:       false,
	}

	if article.Status == "" {
		article.Status = "draft"
	}

	if article.Summary == "" && article.Content != "" {
		summary := s.aiService.GenerateSummary(article.Content)
		if summary != "" {
			article.Summary = summary
			article.IsAI = true
		}
	}

	if err := s.repo.Create(article); err != nil {
		return nil, err
	}

	s.db.Model(&user.User{}).Where("id = ?", userID).UpdateColumn("article_count", gorm.Expr("article_count + 1"))

	return article, nil
}

func (s *Service) Update(id uint, userID uint, req *UpdateArticleReq) (*Article, error) {
	article, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("article not found")
	}

	if article.UserID != userID {
		return nil, errors.New("permission denied")
	}

	if req.Title != "" {
		article.Title = req.Title
	}
	if req.Content != "" {
		article.Content = req.Content
	}
	if req.Summary != "" {
		article.Summary = req.Summary
	}
	if req.Cover != "" {
		article.Cover = req.Cover
	}
	if req.Status != "" {
		article.Status = req.Status
	}
	if req.CategoryID != 0 {
		article.CategoryID = req.CategoryID
	}
	if req.Tags != "" {
		article.Tags = req.Tags
	}

	if err := s.repo.Update(article); err != nil {
		return nil, err
	}

	return article, nil
}

func (s *Service) GetByID(id uint) (*Article, error) {
	article, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("article not found")
	}

	s.repo.IncrementViewCount(id)
	return article, nil
}

func (s *Service) List(page, pageSize int, status string, categoryID, userID uint) ([]Article, int64, error) {
	conditions := make(map[string]interface{})
	if status != "" {
		conditions["status"] = status
	} else {
		conditions["status"] = "published"
	}
	if categoryID > 0 {
		conditions["category_id"] = categoryID
	}
	if userID > 0 {
		conditions["user_id"] = userID
	}

	return s.repo.List(page, pageSize, conditions)
}

func (s *Service) Delete(id, userID uint) error {
	article, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("article not found")
	}

	if article.UserID != userID {
		return errors.New("permission denied")
	}

	if err := s.repo.Delete(id); err != nil {
		return err
	}

	s.db.Model(&user.User{}).Where("id = ?", userID).UpdateColumn("article_count", gorm.Expr("article_count - 1"))
	return nil
}

func (s *Service) LikeArticle(userID, articleID uint) error {
	article, err := s.repo.FindByID(articleID)
	if err != nil {
		return errors.New("article not found")
	}

	liked, _ := s.repo.IsLiked(userID, "article", articleID)
	if liked {
		s.repo.DeleteLike(userID, "article", articleID)
		s.repo.DecrementLikeCount(articleID)
		return nil
	}

	like := &Like{
		UserID:     userID,
		TargetType: "article",
		TargetID:   articleID,
	}
	if err := s.repo.CreateLike(like); err != nil {
		return err
	}

	s.repo.IncrementLikeCount(articleID)
	_ = article
	return nil
}

func (s *Service) FavoriteArticle(userID, articleID uint) error {
	_, err := s.repo.FindByID(articleID)
	if err != nil {
		return errors.New("article not found")
	}

	favorited, _ := s.repo.IsFavorited(userID, articleID)
	if favorited {
		s.repo.DeleteFavorite(userID, articleID)
		s.repo.DecrementFavoriteCount(articleID)
		return nil
	}

	fav := &Favorite{
		UserID:    userID,
		ArticleID: articleID,
	}
	if err := s.repo.CreateFavorite(fav); err != nil {
		return err
	}

	s.repo.IncrementFavoriteCount(articleID)
	return nil
}

func (s *Service) GetUserFavorites(userID uint, page, pageSize int) ([]Article, int64, error) {
	return s.repo.GetUserFavorites(userID, page, pageSize)
}

func (s *Service) GetHot(limit int) ([]Article, error) {
	if limit <= 0 {
		limit = 10
	}
	return s.repo.GetHot(limit)
}

func (s *Service) GetFeed(userID uint, page, pageSize int) ([]Article, int64, error) {
	return s.repo.GetFeed(userID, page, pageSize)
}

func (s *Service) SearchTitle(keyword string, page, pageSize int) ([]Article, int64, error) {
	return s.repo.SearchTitle(keyword, page, pageSize)
}
