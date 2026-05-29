package user

import (
	"net"
	"strconv"

	"inference-engine/internal/middleware"
	"inference-engine/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

func normalizeIP(ip string) string {
	if ip == "::1" {
		return "127.0.0.1"
	}
	parsed := net.ParseIP(ip)
	if parsed != nil {
		if v4 := parsed.To4(); v4 != nil {
			return v4.String()
		}
	}
	return ip
}

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(c *gin.Context) {
	var req RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.ErrBadRequest, err.Error())
		return
	}

	ip := normalizeIP(c.ClientIP())
	user, err := h.svc.Register(&req, ip)
	if err != nil {
		response.Error(c, response.ErrBadRequest, err.Error())
		return
	}

	response.Success(c, user)
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.ErrBadRequest, err.Error())
		return
	}

	ip := normalizeIP(c.ClientIP())
	ua := c.GetHeader("User-Agent")

	resp, err := h.svc.Login(&req, ip, ua)
	if err != nil {
		response.Error(c, response.ErrBadRequest, err.Error())
		return
	}

	response.Success(c, resp)
}

func (h *Handler) Logout(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID > 0 {
		h.svc.RecordLogout(userID)
	}
	response.Success(c, nil)
}

func (h *Handler) GetProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)
	resp, err := h.svc.GetProfile(userID)
	if err != nil {
		response.Error(c, response.ErrBadRequest, err.Error())
		return
	}

	response.Success(c, resp)
}

func (h *Handler) UpdateProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req UpdateProfileReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.ErrBadRequest, err.Error())
		return
	}

	ip := c.ClientIP()
	user, err := h.svc.UpdateProfile(userID, &req, ip)
	if err != nil {
		response.Error(c, response.ErrBadRequest, err.Error())
		return
	}

	response.Success(c, user)
}

func (h *Handler) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, response.ErrBadRequest, "invalid user id")
		return
	}

	user, err := h.svc.GetUserByID(uint(id))
	if err != nil {
		response.Error(c, response.ErrNotFound, "user not found")
		return
	}

	response.Success(c, user)
}

func (h *Handler) Follow(c *gin.Context) {
	userID := middleware.GetUserID(c)
	idStr := c.Param("id")
	followedID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, response.ErrBadRequest, "invalid user id")
		return
	}

	ip := c.ClientIP()
	if err := h.svc.Follow(userID, uint(followedID), ip); err != nil {
		response.Error(c, response.ErrBadRequest, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *Handler) Unfollow(c *gin.Context) {
	userID := middleware.GetUserID(c)
	idStr := c.Param("id")
	followedID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, response.ErrBadRequest, "invalid user id")
		return
	}

	ip := c.ClientIP()
	if err := h.svc.Unfollow(userID, uint(followedID), ip); err != nil {
		response.Error(c, response.ErrBadRequest, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *Handler) GetFollowers(c *gin.Context) {
	idStr := c.Param("id")
	userID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, response.ErrBadRequest, "invalid user id")
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

	users, total, err := h.svc.GetFollowers(uint(userID), page, pageSize)
	if err != nil {
		response.Error(c, response.ErrInternal, err.Error())
		return
	}

	response.Page(c, users, total, page, pageSize)
}

func (h *Handler) GetFollowing(c *gin.Context) {
	idStr := c.Param("id")
	userID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, response.ErrBadRequest, "invalid user id")
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

	users, total, err := h.svc.GetFollowing(uint(userID), page, pageSize)
	if err != nil {
		response.Error(c, response.ErrInternal, err.Error())
		return
	}

	response.Page(c, users, total, page, pageSize)
}

func RegisterRoutes(r *gin.RouterGroup, handler *Handler, oauthHandler *OAuthHandler) {
	auth := r.Group("/auth")
	{
		auth.POST("/register", handler.Register)
		auth.POST("/login", handler.Login)
		auth.POST("/logout", handler.Logout)
		// OAuth 登录已禁用
		// auth.GET("/google", oauthHandler.GoogleLogin)
		// auth.GET("/google/callback", oauthHandler.GoogleCallback)
		// auth.GET("/github", oauthHandler.GitHubLogin)
		// auth.GET("/github/callback", oauthHandler.GitHubCallback)
		// auth.GET("/wechat", oauthHandler.WeChatLogin)
		// auth.GET("/wechat/callback", oauthHandler.WeChatCallback)
	}

	user := r.Group("/user")
	user.Use(middleware.Auth())
	{
		user.GET("/profile", handler.GetProfile)
		user.PUT("/profile", handler.UpdateProfile)
	}

	r.GET("/user/:id", handler.GetUser)
	r.POST("/user/:id/follow", middleware.Auth(), handler.Follow)
	r.DELETE("/user/:id/follow", middleware.Auth(), handler.Unfollow)
	r.GET("/user/:id/followers", handler.GetFollowers)
	r.GET("/user/:id/following", handler.GetFollowing)
}
