package admin

import (
	"time"

	"gorm.io/gorm"
)

type Announcement struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Title     string    `gorm:"size:200;not null" json:"title"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	Status    int       `gorm:"default:1;comment:1=发布 0=草稿" json:"status"`
	Priority  int       `gorm:"default:0;comment:优先级" json:"priority"`
	Author    string    `gorm:"size:50" json:"author"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Report struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	ReporterID uint      `gorm:"index;not null" json:"reporter_id"`
	TargetType string    `gorm:"size:20;not null;comment:article/comment/prompt/user" json:"target_type"`
	TargetID   uint      `gorm:"index;not null" json:"target_id"`
	Reason     string    `gorm:"size:500;not null" json:"reason"`
	Status     int       `gorm:"default:0;comment:0=待处理 1=已处理 2=已忽略" json:"status"`
	Result     string    `gorm:"size:200;comment:处理结果" json:"result"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type AuditLog struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	UserID    uint      `gorm:"index;not null" json:"user_id"`
	Username  string    `gorm:"size:50" json:"username"`
	Action    string    `gorm:"size:100;not null;comment:操作类型" json:"action"`
	Target    string    `gorm:"size:200;comment:操作对象" json:"target"`
	Detail    string    `gorm:"type:text;comment:操作详情" json:"detail"`
	IP        string    `gorm:"size:50" json:"ip"`
	CreatedAt time.Time `json:"created_at"`
}

type SystemConfig struct {
	ID    uint   `gorm:"primarykey" json:"id"`
	Key   string `gorm:"uniqueIndex;size:100;not null" json:"key"`
	Value string `gorm:"type:text" json:"value"`
	Desc  string `gorm:"size:200" json:"desc"`
}

type AIUsageStat struct {
	ID               uint      `gorm:"primarykey" json:"id"`
	Date             string    `gorm:"size:10;index;not null" json:"date"`
	TotalRequests    int       `gorm:"default:0" json:"total_requests"`
	TotalTokens      int       `gorm:"default:0" json:"total_tokens"`
	PromptTokens     int       `gorm:"default:0" json:"prompt_tokens"`
	CompletionTokens int       `gorm:"default:0" json:"completion_tokens"`
	ErrorCount       int       `gorm:"default:0" json:"error_count"`
	CreatedAt        time.Time `json:"created_at"`
}

// New models for extended features

type UserBan struct {
	ID        uint       `gorm:"primarykey" json:"id"`
	UserID    uint       `gorm:"index;not null" json:"user_id"`
	Reason    string     `gorm:"size:500;not null" json:"reason"`
	BanType   string     `gorm:"size:20;not null;comment:temporary/permanent" json:"ban_type"`
	ExpiresAt *time.Time `json:"expires_at"`
	Operator  string     `gorm:"size:50" json:"operator"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type UserTag struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	UserID    uint      `gorm:"index;not null" json:"user_id"`
	Tag       string    `gorm:"size:50;not null" json:"tag"`
	CreatedAt time.Time `json:"created_at"`
}

type SensitiveWord struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Word      string    `gorm:"uniqueIndex;size:100;not null" json:"word"`
	Level     int       `gorm:"default:1;comment:1=禁止 2=替换" json:"level"`
	ReplaceTo string    `gorm:"size:100" json:"replace_to"`
	CreatedAt time.Time `json:"created_at"`
}

type ContentReview struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	TargetType string    `gorm:"size:20;not null" json:"target_type"`
	TargetID   uint      `gorm:"index;not null" json:"target_id"`
	Status     int       `gorm:"default:0;comment:0=待审核 1=通过 2=拒绝" json:"status"`
	ReviewerID uint      `json:"reviewer_id"`
	Reason     string    `gorm:"size:500" json:"reason"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type RecommendItem struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	TargetType string    `gorm:"size:20;not null" json:"target_type"`
	TargetID   uint      `gorm:"not null" json:"target_id"`
	Position   string    `gorm:"size:50;not null" json:"position"`
	SortOrder  int       `gorm:"default:0" json:"sort_order"`
	Status     int       `gorm:"default:1" json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

type PointsRule struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Action    string    `gorm:"size:50;not null;uniqueIndex" json:"action"`
	Points    int       `gorm:"not null" json:"points"`
	Desc      string    `gorm:"size:200" json:"desc"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type InviteCode struct {
	ID        uint       `gorm:"primarykey" json:"id"`
	Code      string     `gorm:"uniqueIndex;size:20;not null" json:"code"`
	CreatorID uint       `json:"creator_id"`
	UsedBy    uint       `json:"used_by"`
	MaxUses   int        `gorm:"default:1" json:"max_uses"`
	UsedCount int        `gorm:"default:0" json:"used_count"`
	ExpiresAt *time.Time `json:"expires_at"`
	CreatedAt time.Time  `json:"created_at"`
}

type LoginLog struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	UserID    uint      `gorm:"index;not null" json:"user_id"`
	IP        string    `gorm:"size:50" json:"ip"`
	UserAgent string    `gorm:"size:500" json:"user_agent"`
	Status    int       `gorm:"default:1;comment:1=成功 0=失败" json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type IPBlacklist struct {
	ID        uint       `gorm:"primarykey" json:"id"`
	IP        string     `gorm:"uniqueIndex;size:50;not null" json:"ip"`
	Reason    string     `gorm:"size:200" json:"reason"`
	ExpiresAt *time.Time `json:"expires_at"`
	CreatedAt time.Time  `json:"created_at"`
}

type SystemLog struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Level     string    `gorm:"size:20;not null" json:"level"`
	Message   string    `gorm:"type:text;not null" json:"message"`
	Source    string    `gorm:"size:100" json:"source"`
	Stack     string    `gorm:"type:text" json:"stack"`
	CreatedAt time.Time `json:"created_at"`
}

type PageView struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Path      string    `gorm:"size:200;not null" json:"path"`
	UserID    uint      `json:"user_id"`
	IP        string    `gorm:"size:50" json:"ip"`
	UserAgent string    `gorm:"size:500" json:"user_agent"`
	CreatedAt time.Time `json:"created_at"`
}

// AuditLogService provides audit logging for user actions
type AuditLogService struct {
	db *gorm.DB
}

func NewAuditLogService(db *gorm.DB) *AuditLogService {
	return &AuditLogService{db: db}
}

func (s *AuditLogService) LogUserAction(userID uint, action, target, detail string) {
	var username string
	s.db.Table("users").Where("id = ?", userID).Pluck("username", &username)
	s.db.Create(&AuditLog{UserID: userID, Username: username, Action: action, Target: target, Detail: detail})
}

func (s *AuditLogService) LogAdminAction(userID uint, username, action, target, detail, ip string) {
	s.db.Create(&AuditLog{UserID: userID, Username: username, Action: action, Target: target, Detail: detail, IP: ip})
}
