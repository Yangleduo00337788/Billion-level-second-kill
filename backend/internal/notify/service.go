package notify

import (
	"time"

	"gorm.io/gorm"
)

type Notification struct {
	ID        uint           `gorm:"primarykey;comment:通知ID" json:"id"`
	UserID    uint           `gorm:"index;not null;comment:接收者ID" json:"user_id"`
	ActorID   uint           `gorm:"index;not null;comment:触发者ID" json:"actor_id"`
	Type      string         `gorm:"size:50;not null;comment:通知类型(like/comment/follow)" json:"type"`
	Content   string         `gorm:"size:500;comment:通知内容" json:"content"`
	TargetID  uint           `comment:目标资源ID" json:"target_id"`
	TargetURL string         `gorm:"size:255;comment:目标链接" json:"target_url"`
	IsRead    bool           `gorm:"default:false;comment:是否已读" json:"is_read"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index;comment:软删除时间" json:"-"`
}

func (Notification) TableComment() string {
	return "通知表"
}

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

func (s *Service) AutoMigrate() error {
	return s.db.AutoMigrate(&Notification{})
}

type CreateNotifyReq struct {
	UserID   uint   `json:"user_id"`
	ActorID  uint   `json:"actor_id"`
	Type     string `json:"type"`
	Content  string `json:"content"`
	TargetID uint   `json:"target_id"`
}

func (s *Service) Create(req *CreateNotifyReq) error {
	notify := &Notification{
		UserID:   req.UserID,
		ActorID:  req.ActorID,
		Type:     req.Type,
		Content:  req.Content,
		TargetID: req.TargetID,
		IsRead:   false,
	}

	return s.db.Create(notify).Error
}

func (s *Service) List(userID uint, page, pageSize int) ([]Notification, int64, error) {
	var notifications []Notification
	var total int64

	s.db.Model(&Notification{}).Where("user_id = ?", userID).Count(&total)
	err := s.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset((page - 1) * pageSize).Limit(pageSize).Find(&notifications).Error
	if err != nil {
		return nil, 0, err
	}

	return notifications, total, nil
}

func (s *Service) MarkAsRead(id, userID uint) error {
	return s.db.Model(&Notification{}).Where("id = ? AND user_id = ?", id, userID).Update("is_read", true).Error
}

func (s *Service) MarkAllAsRead(userID uint) error {
	return s.db.Model(&Notification{}).Where("user_id = ?", userID).Update("is_read", true).Error
}

func (s *Service) GetUnreadCount(userID uint) (int64, error) {
	var count int64
	err := s.db.Model(&Notification{}).Where("user_id = ? AND is_read = ?", userID, false).Count(&count).Error
	return count, err
}
