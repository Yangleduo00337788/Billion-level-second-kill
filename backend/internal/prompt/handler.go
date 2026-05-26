package prompt

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

	var req CreatePromptReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.ErrBadRequest, err.Error())
		return
	}

	prompt, err := h.svc.Create(userID, &req)
	if err != nil {
		response.Error(c, response.ErrBadRequest, err.Error())
		return
	}

	response.Success(c, prompt)
}

func (h *Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	category := c.Query("category")
	tag := c.Query("tag")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 10
	}

	prompts, total, err := h.svc.List(page, pageSize, category, tag)
	if err != nil {
		response.Error(c, response.ErrInternal, err.Error())
		return
	}

	response.Page(c, prompts, total, page, pageSize)
}

func (h *Handler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, response.ErrBadRequest, "invalid prompt id")
		return
	}

	prompt, err := h.svc.GetByID(uint(id))
	if err != nil {
		response.Error(c, response.ErrNotFound, err.Error())
		return
	}

	response.Success(c, prompt)
}

func (h *Handler) Update(c *gin.Context) {
	userID := middleware.GetUserID(c)
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, response.ErrBadRequest, "invalid prompt id")
		return
	}

	var req UpdatePromptReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.ErrBadRequest, err.Error())
		return
	}

	prompt, err := h.svc.Update(uint(id), userID, &req)
	if err != nil {
		response.Error(c, response.ErrBadRequest, err.Error())
		return
	}

	response.Success(c, prompt)
}

func (h *Handler) Delete(c *gin.Context) {
	userID := middleware.GetUserID(c)
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, response.ErrBadRequest, "invalid prompt id")
		return
	}

	if err := h.svc.Delete(uint(id), userID); err != nil {
		response.Error(c, response.ErrBadRequest, err.Error())
		return
	}

	response.Success(c, nil)
}

func RegisterRoutes(r *gin.RouterGroup, handler *Handler) {
	prompts := r.Group("/prompts")
	{
		prompts.GET("", handler.List)
		prompts.GET("/:id", handler.GetByID)
		prompts.POST("", middleware.Auth(), handler.Create)
		prompts.PUT("/:id", middleware.Auth(), handler.Update)
		prompts.DELETE("/:id", middleware.Auth(), handler.Delete)
	}
}
