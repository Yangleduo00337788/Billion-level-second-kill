package article

import (
	"errors"
	"fmt"

	"inference-engine/internal/admin"
	"inference-engine/internal/ai"
	"inference-engine/internal/notify"
	"inference-engine/internal/pkg/points"
	"inference-engine/internal/pkg/sensitive"
	"inference-engine/internal/user"

	"gorm.io/gorm"
)

type Service struct {
	repo        *Repository
	userRepo    *user.Repository
	aiService   *ai.Service
	notifySvc   *notify.Service
	auditLogSvc *admin.AuditLogService
	filter      *sensitive.Filter
	pointsSvc   *points.Service
	db          *gorm.DB
}

func NewService(repo *Repository, userRepo *user.Repository, aiService *ai.Service, db *gorm.DB) *Service {
	return &Service{
		repo:        repo,
		userRepo:    userRepo,
		aiService:   aiService,
		notifySvc:   notify.NewService(db),
		auditLogSvc: admin.NewAuditLogService(db),
		filter:      sensitive.NewFilter(db),
		pointsSvc:   points.NewService(db),
		db:          db,
	}
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

func (s *Service) Create(userID uint, req *CreateArticleReq, ip string) (*Article, error) {
	// Sensitive word check
	title, forbidden, words := s.filter.CheckContent(req.Title)
	if forbidden {
		return nil, fmt.Errorf("标题包含违禁词: %v", words)
	}
	content, forbidden, words := s.filter.CheckContent(req.Content)
	if forbidden {
		return nil, fmt.Errorf("内容包含违禁词: %v", words)
	}

	tags := req.Tags
	if tags == "" {
		tags = "[]"
	}

	article := &Article{
		UserID:  userID,
		Title:   title,
		Content: content,
		Summary: req.Summary,
		Cover:   req.Cover,
		Status:  req.Status,
		Tags:    tags,
		IsAI:    false,
	}
	if req.CategoryID > 0 {
		catID := req.CategoryID
		article.CategoryID = &catID
	}

	if article.Status == "" {
		article.Status = "published"
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

	// Audit log
	go s.auditLogSvc.LogUserAction(userID, "创建文章", "article", article.Title, ip)
	go s.pointsSvc.AwardPoints(userID, "publish_article")

	return article, nil
}

func (s *Service) Update(id uint, userID uint, req *UpdateArticleReq, ip string) (*Article, error) {
	article, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("文章不存在")
	}

	if article.UserID != userID {
		return nil, errors.New("没有权限执行此操作")
	}

	if req.Title != "" {
		title, forbidden, words := s.filter.CheckContent(req.Title)
		if forbidden {
			return nil, fmt.Errorf("标题包含违禁词: %v", words)
		}
		article.Title = title
	}
	if req.Content != "" {
		content, forbidden, words := s.filter.CheckContent(req.Content)
		if forbidden {
			return nil, fmt.Errorf("内容包含违禁词: %v", words)
		}
		article.Content = content
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
		catID := req.CategoryID
		article.CategoryID = &catID
	}
	if req.Tags != "" {
		article.Tags = req.Tags
	} else if article.Tags == "" {
		article.Tags = "[]"
	}

	if err := s.repo.Update(article); err != nil {
		return nil, err
	}

	// Audit log
	go s.auditLogSvc.LogUserAction(userID, "编辑文章", "article", article.Title, ip)

	return article, nil
}

func (s *Service) GetByID(id uint, userID uint) (*Article, error) {
	article, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("文章不存在")
	}

	s.repo.IncrementViewCount(id)

	if userID > 0 {
		liked, favorited := s.repo.GetUserInteraction(userID, id)
		article.Liked = liked
		article.Favorited = favorited
	}

	return article, nil
}

func (s *Service) List(page, pageSize int, status string, categoryID, userID uint, keyword string) ([]Article, int64, error) {
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

	return s.repo.List(page, pageSize, conditions, keyword)
}

func (s *Service) Delete(id, userID uint, ip string) error {
	article, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("文章不存在")
	}

	if article.UserID != userID {
		return errors.New("没有权限执行此操作")
	}

	if err := s.repo.Delete(id); err != nil {
		return err
	}

	s.db.Model(&user.User{}).Where("id = ?", userID).UpdateColumn("article_count", gorm.Expr("GREATEST(article_count - 1, 0)"))

	// Audit log
	go s.auditLogSvc.LogUserAction(userID, "删除文章", "article", article.Title, ip)

	return nil
}

func (s *Service) LikeArticle(userID, articleID uint, ip string) error {
	article, err := s.repo.FindByID(articleID)
	if err != nil {
		return errors.New("文章不存在")
	}

	// 使用 ToggleLike 原子操作，避免并发问题
	isLiked, err := s.repo.ToggleLike(userID, "article", articleID)
	if err != nil {
		return err
	}

	if isLiked {
		// 点赞
		s.repo.IncrementLikeCount(articleID)

		// Trigger notification to article author
		if article.UserID != userID {
			actor, _ := s.userRepo.FindByID(userID)
			actorName := "用户"
			if actor != nil {
				actorName = actor.Username
			}
			go s.notifySvc.Create(&notify.CreateNotifyReq{
				UserID:   article.UserID,
				ActorID:  userID,
				Type:     "like",
				Content:  fmt.Sprintf("%s 点赞了你的文章《%s》", actorName, article.Title),
				TargetID: articleID,
			})
		}

		// Audit log
		go s.auditLogSvc.LogUserAction(userID, "点赞文章", "article", article.Title, ip)
		go s.pointsSvc.AwardPoints(userID, "like")
	} else {
		// 取消点赞
		s.repo.DecrementLikeCount(articleID)
	}

	return nil
}

func (s *Service) FavoriteArticle(userID, articleID uint, ip string) error {
	article, err := s.repo.FindByID(articleID)
	if err != nil {
		return errors.New("文章不存在")
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

	// Audit log
	go s.auditLogSvc.LogUserAction(userID, "收藏文章", "article", article.Title, ip)

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

func (s *Service) ListCategories() ([]Category, error) {
	return s.repo.ListCategories()
}
