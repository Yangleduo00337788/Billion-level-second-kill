package chat

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Hub struct {
	clients    map[uint]*Client
	broadcast  chan *ChatMessage
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

type Client struct {
	UserID uint
	Conn   *websocket.Conn
	Send   chan []byte
	Hub    *Hub
}

type ChatMessage struct {
	ID         uint   `json:"id"`
	ChatID     uint   `json:"chat_id"`
	SenderID   uint   `json:"sender_id"`
	ReceiverID uint   `json:"receiver_id"`
	Content    string `json:"content"`
	MsgType    string `json:"msg_type"`
	FileURL    string `json:"file_url"`
	Status     string `json:"status"`
	CreatedAt  string `json:"created_at"`
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[uint]*Client),
		broadcast:  make(chan *ChatMessage, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if existing, ok := h.clients[client.UserID]; ok {
				close(existing.Send)
				delete(h.clients, client.UserID)
			}
			h.clients[client.UserID] = client
			h.mu.Unlock()
			log.Printf("client connected: user %d", client.UserID)

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.UserID]; ok {
				delete(h.clients, client.UserID)
				close(client.Send)
			}
			h.mu.Unlock()
			log.Printf("client disconnected: user %d", client.UserID)

		case message := <-h.broadcast:
			h.mu.RLock()
			client, ok := h.clients[message.SenderID]
			h.mu.RUnlock()
			if ok {
				data, err := json.Marshal(message)
				if err != nil {
					log.Printf("failed to marshal message: %v", err)
					continue
				}
				select {
				case client.Send <- data:
				default:
					h.mu.Lock()
					delete(h.clients, client.UserID)
					close(client.Send)
					h.mu.Unlock()
				}
			}
		}
	}
}

func (h *Hub) SendToUser(userID uint, message *ChatMessage) {
	h.mu.RLock()
	client, ok := h.clients[userID]
	h.mu.RUnlock()
	if !ok {
		return
	}

	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("failed to marshal message: %v", err)
		return
	}

	select {
	case client.Send <- data:
	default:
		h.mu.Lock()
		delete(h.clients, client.UserID)
		close(client.Send)
		h.mu.Unlock()
	}
}

func (h *Hub) BroadcastMessage(message *ChatMessage) {
	h.broadcast <- message
}

func (h *Hub) IsOnline(userID uint) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	_, ok := h.clients[userID]
	return ok
}

func (h *Hub) GetOnlineUsers() []uint {
	h.mu.RLock()
	defer h.mu.RUnlock()
	users := make([]uint, 0, len(h.clients))
	for userID := range h.clients {
		users = append(users, userID)
	}
	return users
}

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 4096
)

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("websocket error: %v", err)
			}
			break
		}

		var msg ChatMessage
		if err := json.Unmarshal(message, &msg); err != nil {
			continue
		}
		msg.SenderID = c.UserID
		msg.CreatedAt = time.Now().Format(time.RFC3339)

		msg.Status = "sent"

		if msg.ReceiverID > 0 {
			c.Hub.SendToUser(msg.ReceiverID, &msg)
		}
		c.Hub.SendToUser(msg.SenderID, &msg)
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte("\n"))
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
