package points

import (
	"log"
	"time"

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
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	db.AutoMigrate(&PointsLog{})
	return &Service{db: db}
}

// AwardPoints awards points to a user for an action
func (s *Service) AwardPoints(userID uint, action string) {
	var rule PointsRule
	if err := s.db.Where("action = ? AND status = 1", action).First(&rule).Error; err != nil {
		return // No rule for this action
	}

	if rule.Points <= 0 {
		return
	}

	// Check limit type
	if !s.canAward(userID, action, rule.LimitType) {
		return
	}

	// Update user points atomically
	result := s.db.Exec("UPDATE users SET points = points + ? WHERE id = ? AND deleted_at IS NULL", rule.Points, userID)
	if result.Error != nil {
		log.Printf("Failed to award points to user %d: %v", userID, result.Error)
		return
	}

	// Get new balance
	var balance int
	s.db.Raw("SELECT points FROM users WHERE id = ?", userID).Scan(&balance)

	// Log the points transaction
	s.db.Create(&PointsLog{
		UserID:  userID,
		Action:  action,
		Points:  rule.Points,
		Balance: balance,
	})
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
