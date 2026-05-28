package comment

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
	articleIDStr := c.Param("id")
	articleID, err := strconv.ParseUint(articleIDStr, 10, 64)
	if err != nil {
		response.Error(c, response.ErrBadRequest, "invalid article id")
		return
	}

	var req CreateCommentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.ErrBadRequest, err.Error())
		return
	}

	ip := c.ClientIP()
	comment, err := h.svc.Create(userID, uint(articleID), &req, ip)
	if err != nil {
		response.Error(c, response.ErrBadRequest, err.Error())
		return
	}

	response.Success(c, comment)
}

func (h *Handler) List(c *gin.Context) {
	articleIDStr := c.Param("id")
	articleID, err := strconv.ParseUint(articleIDStr, 10, 64)
	if err != nil {
		response.Error(c, response.ErrBadRequest, "invalid article id")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 10
	}

	comments, total, err := h.svc.GetByArticleID(uint(articleID), page, pageSize)
	if err != nil {
		response.Error(c, response.ErrInternal, err.Error())
		return
	}

	response.Page(c, comments, total, page, pageSize)
}

func (h *Handler) Delete(c *gin.Context) {
	userID := middleware.GetUserID(c)
	commentIDStr := c.Param("id")
	commentID, err := strconv.ParseUint(commentIDStr, 10, 64)
	if err != nil {
		response.Error(c, response.ErrBadRequest, "invalid comment id")
		return
	}

	ip := c.ClientIP()
	if err := h.svc.Delete(uint(commentID), userID, ip); err != nil {
		response.Error(c, response.ErrBadRequest, err.Error())
		return
	}

	response.Success(c, nil)
}

func RegisterRoutes(r *gin.RouterGroup, articleGroup *gin.RouterGroup, handler *Handler) {
	articles := articleGroup
	articles.POST("/:id/comments", middleware.Auth(), handler.Create)
	articles.GET("/:id/comments", handler.List)

	comments := r.Group("/comments")
	comments.Use(middleware.Auth())
	{
		comments.DELETE("/:id", handler.Delete)
	}
}
