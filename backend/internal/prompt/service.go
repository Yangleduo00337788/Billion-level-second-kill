package prompt

import (
	"errors"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

type CreatePromptReq struct {
	Title       string `json:"title" binding:"required,max=200"`
	Content     string `json:"content" binding:"required"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Model       string `json:"model"`
	Tags        string `json:"tags"`
}

type UpdatePromptReq struct {
	Title       string `json:"title"`
	Content     string `json:"content"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Model       string `json:"model"`
	Tags        string `json:"tags"`
}

func (s *Service) Create(userID uint, req *CreatePromptReq) (*Prompt, error) {
	p := &Prompt{
		UserID:      userID,
		Title:       req.Title,
		Content:     req.Content,
		Description: req.Description,
		Category:    req.Category,
		Model:       req.Model,
		Tags:        req.Tags,
		Status:      1,
	}

	if err := s.repo.Create(p); err != nil {
		return nil, err
	}

	return p, nil
}

func (s *Service) GetByID(id uint) (*Prompt, error) {
	p, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("prompt not found")
	}

	s.repo.IncrementUsageCount(id)
	return p, nil
}

func (s *Service) Update(id, userID uint, req *UpdatePromptReq) (*Prompt, error) {
	p, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("prompt not found")
	}

	if p.UserID != userID {
		return nil, errors.New("permission denied")
	}

	if req.Title != "" {
		p.Title = req.Title
	}
	if req.Content != "" {
		p.Content = req.Content
	}
	if req.Description != "" {
		p.Description = req.Description
	}
	if req.Category != "" {
		p.Category = req.Category
	}
	if req.Model != "" {
		p.Model = req.Model
	}
	if req.Tags != "" {
		p.Tags = req.Tags
	}

	if err := s.repo.Update(p); err != nil {
		return nil, err
	}

	return p, nil
}

func (s *Service) Delete(id, userID uint) error {
	p, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("prompt not found")
	}

	if p.UserID != userID {
		return errors.New("permission denied")
	}

	return s.repo.Delete(id)
}

func (s *Service) List(page, pageSize int, category, tag string) ([]Prompt, int64, error) {
	return s.repo.List(page, pageSize, category, tag)
}

func (s *Service) SearchTitle(keyword string, page, pageSize int) ([]Prompt, int64, error) {
	return s.repo.SearchTitle(keyword, page, pageSize)
}
