package article

import (
	"strconv"

	"inference-engine/internal/middleware"
	"inference-engine/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Create(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req CreateArticleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.ErrBadRequest, err.Error())
		return
	}

	article, err := h.svc.Create(userID, &req)
	if err != nil {
		response.Error(c, response.ErrBadRequest, err.Error())
		return
	}

	response.Success(c, article)
}

func (h *Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status")
	categoryID, _ := strconv.ParseUint(c.Query("category"), 10, 64)
	userIDFilter, _ := strconv.ParseUint(c.Query("user_id"), 10, 64)

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 10
	}

	articles, total, err := h.svc.List(page, pageSize, status, uint(categoryID), uint(userIDFilter))
	if err != nil {
		response.Error(c, response.ErrInternal, err.Error())
		return
	}

	response.Page(c, articles, total, page, pageSize)
}

func (h *Handler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, response.ErrBadRequest, "invalid article id")
		return
	}

	article, err := h.svc.GetByID(uint(id))
	if err != nil {
		response.Error(c, response.ErrNotFound, err.Error())
		return
	}

	response.Success(c, article)
}

func (h *Handler) Update(c *gin.Context) {
	userID := middleware.GetUserID(c)
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, response.ErrBadRequest, "invalid article id")
		return
	}

	var req UpdateArticleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.ErrBadRequest, err.Error())
		return
	}

	article, err := h.svc.Update(uint(id), userID, &req)
	if err != nil {
		response.Error(c, response.ErrBadRequest, err.Error())
		return
	}

	response.Success(c, article)
}

func (h *Handler) Delete(c *gin.Context) {
	userID := middleware.GetUserID(c)
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, response.ErrBadRequest, "invalid article id")
		return
	}

	if err := h.svc.Delete(uint(id), userID); err != nil {
		response.Error(c, response.ErrBadRequest, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *Handler) Like(c *gin.Context) {
	userID := middleware.GetUserID(c)
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, response.ErrBadRequest, "invalid article id")
		return
	}

	if err := h.svc.LikeArticle(userID, uint(id)); err != nil {
		response.Error(c, response.ErrBadRequest, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *Handler) Favorite(c *gin.Context) {
	userID := middleware.GetUserID(c)
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, response.ErrBadRequest, "invalid article id")
		return
	}

	if err := h.svc.FavoriteArticle(userID, uint(id)); err != nil {
		response.Error(c, response.ErrBadRequest, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *Handler) Hot(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	articles, err := h.svc.GetHot(limit)
	if err != nil {
		response.Error(c, response.ErrInternal, err.Error())
		return
	}

	response.Success(c, articles)
}

func (h *Handler) Feed(c *gin.Context) {
	userID := middleware.GetUserID(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 10
	}

	articles, total, err := h.svc.GetFeed(userID, page, pageSize)
	if err != nil {
		response.Error(c, response.ErrInternal, err.Error())
		return
	}

	_ = total
	response.Success(c, articles)
}

func RegisterRoutes(r *gin.RouterGroup, handler *Handler) {
	articles := r.Group("/articles")
	{
		articles.GET("", handler.List)
		articles.GET("/hot", handler.Hot)
		articles.GET("/feed", middleware.Auth(), handler.Feed)
		articles.GET("/:id", handler.GetByID)
		articles.POST("", middleware.Auth(), handler.Create)
		articles.PUT("/:id", middleware.Auth(), handler.Update)
		articles.DELETE("/:id", middleware.Auth(), handler.Delete)
		articles.POST("/:id/like", middleware.Auth(), handler.Like)
		articles.POST("/:id/favorite", middleware.Auth(), handler.Favorite)
	}
}
