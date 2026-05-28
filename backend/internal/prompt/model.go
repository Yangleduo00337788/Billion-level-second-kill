package prompt

import (
	"time"

	"inference-engine/internal/user"

	"gorm.io/gorm"
)

type Prompt struct {
	ID            uint           `gorm:"primarykey;comment:Prompt ID" json:"id"`
	UserID        uint           `gorm:"index;not null;comment:创建者ID" json:"user_id"`
	Title         string         `gorm:"size:200;not null;comment:Prompt标题" json:"title"`
	Content       string         `gorm:"type:text;not null;comment:Prompt内容" json:"content"`
	Description   string         `gorm:"size:500;comment:Prompt描述" json:"description"`
	Category      string         `gorm:"size:50;index;comment:分类" json:"category"`
	Model         string         `gorm:"size:50;comment:适用模型" json:"model"`
	Tags          string         `gorm:"type:json;comment:标签列表(JSON)" json:"tags"`
	UsageCount    int            `gorm:"default:0;comment:使用次数" json:"usage_count"`
	LikeCount     int            `gorm:"default:0;comment:点赞数" json:"like_count"`
	FavoriteCount int            `gorm:"default:0;comment:收藏数" json:"favorite_count"`
	Rating        float64        `gorm:"default:0;comment:评分" json:"rating"`
	Status        int            `gorm:"default:1;comment:状态(1=正常 0=禁用)" json:"status"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index;comment:软删除时间" json:"-"`
	User          user.User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

type PromptLike struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	UserID    uint           `gorm:"index:idx_prompt_like,unique;not null" json:"user_id"`
	PromptID  uint           `gorm:"index:idx_prompt_like,unique;not null" json:"prompt_id"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type PromptFavorite struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	UserID    uint           `gorm:"index:idx_prompt_fav,unique;not null" json:"user_id"`
	PromptID  uint           `gorm:"index:idx_prompt_fav,unique;not null" json:"prompt_id"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Prompt) TableComment() string {
	return "Prompt提示词表"
}
