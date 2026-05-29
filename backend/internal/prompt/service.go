package prompt

import (
	"errors"
	"inference-engine/internal/pkg/points"

	"gorm.io/gorm"
)

type Service struct {
	repo      *Repository
	db        *gorm.DB
	pointsSvc *points.Service
}

func NewService(repo *Repository, db *gorm.DB) *Service {
	return &Service{
		repo:      repo,
		db:        db,
		pointsSvc: points.NewService(db),
	}
}

type CreatePromptReq struct {
	Title       string `json:"title" binding:"required,max=200"`
	Description string `json:"description"`
	Content     string `json:"content" binding:"required"`
	Category    string `json:"category"`
	Model       string `json:"model"`
	Tags        string `json:"tags"`
}

type UpdatePromptReq struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Content     string `json:"content"`
	Category    string `json:"category"`
	Model       string `json:"model"`
	Tags        string `json:"tags"`
}

func (s *Service) Create(userID uint, req *CreatePromptReq) (*Prompt, error) {
	tags := req.Tags
	if tags == "" {
		tags = "[]"
	}

	prompt := &Prompt{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		Content:     req.Content,
		Category:    req.Category,
		Model:       req.Model,
		Tags:        tags,
		Status:      1,
	}

	if err := s.repo.Create(prompt); err != nil {
		return nil, err
	}

	// 奖励积分
	go s.pointsSvc.AwardPoints(userID, "publish_prompt")

	return prompt, nil
}

func (s *Service) Update(id uint, userID uint, req *UpdatePromptReq) (*Prompt, error) {
	prompt, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("Prompt 不存在")
	}

	if prompt.UserID != userID {
		return nil, errors.New("没有权限执行此操作")
	}

	if req.Title != "" {
		prompt.Title = req.Title
	}
	if req.Description != "" {
		prompt.Description = req.Description
	}
	if req.Content != "" {
		prompt.Content = req.Content
	}
	if req.Category != "" {
		prompt.Category = req.Category
	}
	if req.Model != "" {
		prompt.Model = req.Model
	}
	if req.Tags != "" {
		prompt.Tags = req.Tags
	}

	if err := s.repo.Update(prompt); err != nil {
		return nil, err
	}

	return prompt, nil
}

func (s *Service) GetByID(id uint) (*Prompt, error) {
	prompt, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("Prompt 不存在")
	}
	s.repo.IncrementUsageCount(id)
	return prompt, nil
}

func (s *Service) List(page, pageSize int, category, tag string, userID uint, keyword string) ([]Prompt, int64, error) {
	return s.repo.List(page, pageSize, category, tag, userID, keyword)
}

func (s *Service) Delete(id, userID uint) error {
	prompt, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("Prompt 不存在")
	}

	if prompt.UserID != userID {
		return errors.New("没有权限执行此操作")
	}

	return s.repo.Delete(id)
}

func (s *Service) GetHot(limit int) ([]Prompt, error) {
	if limit <= 0 {
		limit = 10
	}
	return s.repo.GetHot(limit)
}

func (s *Service) Like(userID, promptID uint) error {
	prompt, err := s.repo.FindByID(promptID)
	if err != nil {
		return errors.New("Prompt 不存在")
	}

	liked, _ := s.repo.IsLiked(userID, promptID)
	if liked {
		s.repo.DeleteLike(userID, promptID)
		s.repo.DecrementLikeCount(promptID)
		return nil
	}

	like := &PromptLike{UserID: userID, PromptID: promptID}
	if err := s.repo.CreateLike(like); err != nil {
		return err
	}
	s.repo.IncrementLikeCount(promptID)

	_ = prompt // suppress unused warning
	return nil
}

func (s *Service) Favorite(userID, promptID uint) error {
	_, err := s.repo.FindByID(promptID)
	if err != nil {
		return errors.New("Prompt 不存在")
	}

	favorited, _ := s.repo.IsFavorited(userID, promptID)
	if favorited {
		s.repo.DeleteFavorite(userID, promptID)
		s.repo.DecrementFavoriteCount(promptID)
		return nil
	}

	fav := &PromptFavorite{UserID: userID, PromptID: promptID}
	if err := s.repo.CreateFavorite(fav); err != nil {
		return err
	}
	s.repo.IncrementFavoriteCount(promptID)
	return nil
}
