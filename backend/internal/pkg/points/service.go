package points

import (
	"fmt"
	"log"
	"time"

	"inference-engine/internal/notify"

	"gorm.io/gorm"
)

type PointsRule struct {
	ID        uint   `gorm:"primarykey"`
	Action    string `gorm:"size:50;uniqueIndex;comment:积分行为"`
	Points    int    `gorm:"not null;comment:积分数量"`
	Desc      string `gorm:"size:200;comment:描述"`
	LimitType string `gorm:"size:20;default:unlimited;comment:限制类型: unlimited=无限制, once=仅首次, daily=每日一次, weekly=每周一次, monthly=每月一次"`
	Status    int    `gorm:"default:1;comment:状态 1=启用 0=禁用"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (PointsRule) TableName() string {
	return "points_rules"
}

type PointsLog struct {
	ID        uint      `gorm:"primarykey"`
	UserID    uint      `gorm:"index;not null;comment:用户ID"`
	Action    string    `gorm:"size:50;not null;comment:积分行为"`
	Points    int       `gorm:"not null;comment:变动积分"`
	Balance   int       `gorm:"not null;comment:变动后余额"`
	CreatedAt time.Time `gorm:"autoCreateTime;comment:创建时间"`
}

func (PointsLog) TableName() string {
	return "points_logs"
}

type Service struct {
	db        *gorm.DB
	notifySvc *notify.Service
}

func NewService(db *gorm.DB) *Service {
	db.AutoMigrate(&PointsLog{})
	return &Service{
		db:        db,
		notifySvc: notify.NewService(db),
	}
}

// 积分行为的中文描述映射
var actionDescMap = map[string]string{
	"register":        "用户注册",
	"publish_article": "发布文章",
	"publish_prompt":  "发布 Prompt",
	"like":            "点赞",
	"comment":         "评论",
	"follow":          "关注",
	"daily_login":     "每日登录",
}

// AwardPoints awards points to a user for an action
func (s *Service) AwardPoints(userID uint, action string) {
	log.Printf("[Points] AwardPoints called: userID=%d, action=%s", userID, action)

	var rule PointsRule
	if err := s.db.Where("action = ? AND status = 1", action).First(&rule).Error; err != nil {
		log.Printf("[Points] No rule found for action=%s, error=%v", action, err)
		return // No rule for this action
	}

	log.Printf("[Points] Rule found: action=%s, points=%d, limitType=%s", rule.Action, rule.Points, rule.LimitType)

	if rule.Points <= 0 {
		log.Printf("[Points] Points <= 0, skipping")
		return
	}

	// Use transaction for atomicity
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Check limit type with transaction
	if !s.canAwardWithTx(tx, userID, action, rule.LimitType) {
		tx.Rollback()
		log.Printf("[Points] Already awarded today, skipping")
		return
	}

	// Update user points atomically
	result := tx.Exec("UPDATE users SET points = points + ? WHERE id = ? AND deleted_at IS NULL", rule.Points, userID)
	if result.Error != nil {
		tx.Rollback()
		log.Printf("[Points] Failed to award points to user %d: %v", userID, result.Error)
		return
	}

	if result.RowsAffected == 0 {
		tx.Rollback()
		log.Printf("[Points] No rows affected for user %d, user may not exist", userID)
		return
	}

	// Get new balance
	var balance int
	tx.Raw("SELECT points FROM users WHERE id = ?", userID).Scan(&balance)

	// Log the points transaction
	if err := tx.Create(&PointsLog{
		UserID:  userID,
		Action:  action,
		Points:  rule.Points,
		Balance: balance,
	}).Error; err != nil {
		tx.Rollback()
		log.Printf("[Points] Failed to create points log: %v", err)
		return
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		log.Printf("[Points] Failed to commit transaction: %v", err)
		return
	}

	// 发送积分变动通知 (outside transaction)
	actionDesc := action
	if desc, ok := actionDescMap[action]; ok {
		actionDesc = desc
	}
	notifyContent := fmt.Sprintf("恭喜！您通过%s获得 %d 积分，当前积分：%d", actionDesc, rule.Points, balance)
	log.Printf("[Points] Sending notification to user %d: %s", userID, notifyContent)
	err := s.notifySvc.Create(&notify.CreateNotifyReq{
		UserID:   userID,
		ActorID:  userID,
		Type:     "points",
		Content:  notifyContent,
		TargetID: userID,
	})
	if err != nil {
		log.Printf("[Points] Failed to send notification: %v", err)
	} else {
		log.Printf("[Points] Notification sent successfully")
	}

	log.Printf("[Points] Awarded %d points to user %d for action=%s, new balance=%d", rule.Points, userID, action, balance)
}

// canAward checks if points can be awarded based on limit type
func (s *Service) canAward(userID uint, action, limitType string) bool {
	switch limitType {
	case "unlimited":
		return true

	case "once":
		// Check if user has ever received this action
		var count int64
		s.db.Model(&PointsLog{}).Where("user_id = ? AND action = ?", userID, action).Count(&count)
		return count == 0

	case "daily":
		// Check if user has received this action today
		today := time.Now().Format("2006-01-02")
		var count int64
		s.db.Model(&PointsLog{}).Where("user_id = ? AND action = ? AND DATE(created_at) = ?", userID, action, today).Count(&count)
		return count == 0

	case "weekly":
		// Check if user has received this action this week
		weekStart := getWeekStart()
		var count int64
		s.db.Model(&PointsLog{}).Where("user_id = ? AND action = ? AND created_at >= ?", userID, action, weekStart).Count(&count)
		return count == 0

	case "monthly":
		// Check if user has received this action this month
		monthStart := time.Now().Format("2006-01-") + "01"
		var count int64
		s.db.Model(&PointsLog{}).Where("user_id = ? AND action = ? AND created_at >= ?", userID, action, monthStart).Count(&count)
		return count == 0

	default:
		return true
	}
}

// canAwardWithTx checks if points can be awarded based on limit type (with transaction)
func (s *Service) canAwardWithTx(tx *gorm.DB, userID uint, action, limitType string) bool {
	switch limitType {
	case "unlimited":
		return true

	case "once":
		var count int64
		tx.Model(&PointsLog{}).Where("user_id = ? AND action = ?", userID, action).Count(&count)
		return count == 0

	case "daily":
		today := time.Now().Format("2006-01-02")
		var count int64
		tx.Model(&PointsLog{}).Where("user_id = ? AND action = ? AND DATE(created_at) = ?", userID, action, today).Count(&count)
		return count == 0

	case "weekly":
		weekStart := getWeekStart()
		var count int64
		tx.Model(&PointsLog{}).Where("user_id = ? AND action = ? AND created_at >= ?", userID, action, weekStart).Count(&count)
		return count == 0

	case "monthly":
		monthStart := time.Now().Format("2006-01-") + "01"
		var count int64
		tx.Model(&PointsLog{}).Where("user_id = ? AND action = ? AND created_at >= ?", userID, action, monthStart).Count(&count)
		return count == 0

	default:
		return true
	}
}

func getWeekStart() time.Time {
	now := time.Now()
	weekday := now.Weekday()
	if weekday == 0 {
		weekday = 7
	}
	return now.AddDate(0, 0, -int(weekday-1)).Truncate(24 * time.Hour)
}

// GetUserPoints returns the user's current points balance
func (s *Service) GetUserPoints(userID uint) int {
	var points int
	s.db.Raw("SELECT points FROM users WHERE id = ?", userID).Scan(&points)
	return points
}

// GetPointsHistory returns the user's points transaction history
func (s *Service) GetPointsHistory(userID uint, page, pageSize int) ([]PointsLog, int64) {
	var logs []PointsLog
	var total int64
	s.db.Model(&PointsLog{}).Where("user_id = ?", userID).Count(&total)
	s.db.Where("user_id = ?", userID).Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&logs)
	return logs, total
}
