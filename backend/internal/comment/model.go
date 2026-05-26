package comment

import (
	"time"

	"inference-engine/internal/user"

	"gorm.io/gorm"
)

type Comment struct {
	ID        uint           `gorm:"primarykey;comment:评论ID" json:"id"`
	ArticleID uint           `gorm:"index;not null;comment:文章ID" json:"article_id"`
	UserID    uint           `gorm:"index;not null;comment:评论者ID" json:"user_id"`
	ParentID  uint           `gorm:"default:0;index;comment:父评论ID(回复)" json:"parent_id"`
	Content   string         `gorm:"type:text;not null;comment:评论内容" json:"content"`
	LikeCount int            `gorm:"default:0;comment:点赞数" json:"like_count"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index;comment:软删除时间" json:"-"`
	User      user.User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (Comment) TableComment() string {
	return "评论表"
}
