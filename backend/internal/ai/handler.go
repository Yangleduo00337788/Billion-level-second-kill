package ai

import (
	"io"

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

func (h *Handler) Chat(c *gin.Context) {
	var req struct {
		Messages []Message `json:"messages" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.ErrBadRequest, err.Error())
		return
	}

	resp, err := h.svc.Chat(req.Messages)
	if err != nil {
		response.Error(c, response.ErrInternal, err.Error())
		return
	}

	response.Success(c, resp)
}

func (h *Handler) StreamChat(c *gin.Context) {
	var req struct {
		Messages []Message `json:"messages" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.ErrBadRequest, err.Error())
		return
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	ch, err := h.svc.StreamChat(req.Messages)
	if err != nil {
		response.Error(c, response.ErrInternal, err.Error())
		return
	}

	c.Stream(func(w io.Writer) bool {
		content, ok := <-ch
		if !ok {
			return false
		}
		c.SSEvent("message", content)
		return true
	})
}

func (h *Handler) GenerateTitle(c *gin.Context) {
	var req struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.ErrBadRequest, err.Error())
		return
	}

	title := h.svc.GenerateTitle(req.Content)
	response.Success(c, gin.H{"title": title})
}

func (h *Handler) GenerateSummary(c *gin.Context) {
	var req struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.ErrBadRequest, err.Error())
		return
	}

	summary := h.svc.GenerateSummary(req.Content)
	response.Success(c, gin.H{"summary": summary})
}

func (h *Handler) SuggestTags(c *gin.Context) {
	var req struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.ErrBadRequest, err.Error())
		return
	}

	tags := h.svc.SuggestTags(req.Content)
	response.Success(c, gin.H{"tags": tags})
}

func RegisterRoutes(r *gin.RouterGroup, handler *Handler) {
	ai := r.Group("/ai")
	ai.Use(middleware.Auth())
	{
		ai.POST("/chat", handler.Chat)
		ai.POST("/chat/stream", handler.StreamChat)
		ai.POST("/generate-title", handler.GenerateTitle)
		ai.POST("/generate-summary", handler.GenerateSummary)
		ai.POST("/suggest-tags", handler.SuggestTags)
	}
}
