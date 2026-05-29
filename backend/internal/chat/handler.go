package chat

import (
	"net/http"
	"strings"

	"inference-engine/internal/middleware"
	"inference-engine/internal/pkg/jwt"
	"inference-engine/internal/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// 校验 Origin，防止跨站 WebSocket 劫持
		origin := r.Header.Get("Origin")
		if origin == "" {
			return false
		}
		// 允许本地开发和生产环境
		allowedOrigins := []string{
			"http://localhost", "http://localhost:5173", "http://localhost:3000",
			"http://127.0.0.1", "http://127.0.0.1:5173", "http://127.0.0.1:3000",
		}
		for _, allowed := range allowedOrigins {
			if strings.HasPrefix(origin, allowed) {
				return true
			}
		}
		// 生产环境应该配置实际域名
		return false
	},
}

type Handler struct {
	hub *Hub
}

func NewHandler(hub *Hub) *Handler {
	return &Handler{hub: hub}
}

func (h *Handler) HandleWebSocket(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		token = c.GetHeader("Authorization")
		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}
	}

	if token == "" {
		response.ErrorWithStatus(c, 401, response.ErrUnauthorized, "missing token")
		return
	}

	claims, err := jwt.ParseToken(token)
	if err != nil {
		response.ErrorWithStatus(c, 401, response.ErrUnauthorized, "invalid token")
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	client := &Client{
		UserID: claims.UserID,
		Conn:   conn,
		Send:   make(chan []byte, 256),
		Hub:    h.hub,
	}

	h.hub.register <- client

	go client.WritePump()
	go client.ReadPump()
}

func (h *Handler) GetChatList(c *gin.Context) {
	userID := middleware.GetUserID(c)
	_ = userID

	response.Success(c, gin.H{
		"online_users": h.hub.GetOnlineUsers(),
		"is_online":    h.hub.IsOnline(userID),
	})
}

func (h *Handler) GetMessages(c *gin.Context) {
	response.Success(c, []interface{}{})
}

func (h *Handler) SendMessage(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req struct {
		ChatID  uint   `json:"chat_id"`
		Content string `json:"content" binding:"required"`
		MsgType string `json:"msg_type"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.ErrBadRequest, err.Error())
		return
	}

	if req.MsgType == "" {
		req.MsgType = "text"
	}

	msg := &ChatMessage{
		SenderID:  userID,
		Content:   req.Content,
		MsgType:   req.MsgType,
		Status:    "sent",
		CreatedAt: "",
	}

	h.hub.BroadcastMessage(msg)
	response.Success(c, msg)
}

func RegisterRoutes(r *gin.RouterGroup, handler *Handler) {
	r.GET("/ws", handler.HandleWebSocket)

	chat := r.Group("/chat")
	chat.Use(middleware.Auth())
	{
		chat.GET("/list", handler.GetChatList)
		chat.GET("/:id/messages", handler.GetMessages)
		chat.POST("/send", handler.SendMessage)
	}
}
