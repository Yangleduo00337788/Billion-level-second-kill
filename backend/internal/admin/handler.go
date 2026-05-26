package admin

import (
	"strconv"

	"inference-engine/internal/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	db *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{db: db}
}

type DashboardStats struct {
	TotalUsers    int64 `json:"total_users"`
	TotalArticles int64 `json:"total_articles"`
	TotalComments int64 `json:"total_comments"`
	TotalPrompts  int64 `json:"total_prompts"`
	TodayUsers    int64 `json:"today_users"`
	TodayArticles int64 `json:"today_articles"`
}

func (h *Handler) Dashboard(c *gin.Context) {
	var stats DashboardStats
	h.db.Model(&struct{}{}).Table("users").Count(&stats.TotalUsers)
	h.db.Model(&struct{}{}).Table("articles").Count(&stats.TotalArticles)
	h.db.Model(&struct{}{}).Table("comments").Count(&stats.TotalComments)
	h.db.Model(&struct{}{}).Table("prompts").Count(&stats.TotalPrompts)
	h.db.Raw("SELECT COUNT(*) FROM users WHERE DATE(created_at) = CURDATE()").Scan(&stats.TodayUsers)
	h.db.Raw("SELECT COUNT(*) FROM articles WHERE DATE(created_at) = CURDATE()").Scan(&stats.TodayArticles)
	response.Success(c, stats)
}

func (h *Handler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var users []map[string]interface{}
	var total int64
	h.db.Table("users").Count(&total)
	h.db.Table("users").Select("id, username, email, role, level, status, article_count, follow_count, fans_count, created_at, updated_at").Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&users)
	response.Page(c, users, total, page, pageSize)
}

func (h *Handler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, response.ErrBadRequest, "invalid id")
		return
	}

	var req struct {
		Role   string `json:"role"`
		Status *int   `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.ErrBadRequest, err.Error())
		return
	}

	updates := map[string]interface{}{}
	if req.Role != "" {
		updates["role"] = req.Role
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if len(updates) > 0 {
		h.db.Table("users").Where("id = ?", id).Updates(updates)
	}
	response.Success(c, nil)
}

func (h *Handler) ListArticles(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var articles []map[string]interface{}
	var total int64
	h.db.Table("articles").Count(&total)
	h.db.Table("articles").Select("id, user_id, title, status, view_count, like_count, comment_count, is_ai, created_at, updated_at").Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&articles)
	response.Page(c, articles, total, page, pageSize)
}

func (h *Handler) DeleteArticle(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, response.ErrBadRequest, "invalid id")
		return
	}

	h.db.Delete(&struct{}{}, "id = ?", id).Table("articles")
	response.Success(c, nil)
}

func (h *Handler) ListPrompts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var prompts []map[string]interface{}
	var total int64
	h.db.Table("prompts").Count(&total)
	h.db.Table("prompts").Select("id, user_id, title, category, model, usage_count, like_count, status, created_at").Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&prompts)
	response.Page(c, prompts, total, page, pageSize)
}

func (h *Handler) DeletePrompt(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, response.ErrBadRequest, "invalid id")
		return
	}

	h.db.Delete(&struct{}{}, "id = ?", id).Table("prompts")
	response.Success(c, nil)
}
