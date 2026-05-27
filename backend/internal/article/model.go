package article

import (
	"encoding/json"
	"time"

	"inference-engine/internal/user"

	"gorm.io/gorm"
)

type Article struct {
	ID            uint           `gorm:"primarykey" json:"id"`
	UserID        uint           `gorm:"index;not null" json:"user_id"`
	Title         string         `gorm:"size:200;not null" json:"title"`
	Content       string         `gorm:"type:longtext" json:"content"`
	Summary       string         `gorm:"size:500" json:"summary"`
	Cover         string         `gorm:"size:255" json:"cover"`
	Status        string         `gorm:"size:20;default:draft" json:"status"`
	CategoryID    *uint          `gorm:"index" json:"category_id"`
	Tags          string         `gorm:"type:json" json:"tags"`
	ViewCount     int            `gorm:"default:0" json:"view_count"`
	LikeCount     int            `gorm:"default:0" json:"like_count"`
	CommentCount  int            `gorm:"default:0" json:"comment_count"`
	FavoriteCount int            `gorm:"default:0" json:"favorite_count"`
	IsAI          bool           `gorm:"default:false" json:"is_ai"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	User          user.User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Category      *Category      `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Liked         bool           `gorm:"-" json:"liked"`
	Favorited     bool           `gorm:"-" json:"favorited"`
}

func (a *Article) MarshalJSON() ([]byte, error) {
	type Alias Article
	return json.Marshal(&struct {
		*Alias
		Tags interface{} `json:"tags"`
	}{
		Alias: (*Alias)(a),
		Tags:  parseTags(a.Tags),
	})
}

func parseTags(tags string) []string {
	if tags == "" {
		return nil
	}
	var result []string
	if err := json.Unmarshal([]byte(tags), &result); err != nil {
		return []string{tags}
	}
	return result
}

type Category struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Name      string    `gorm:"size:50;not null" json:"name"`
	Desc      string    `gorm:"size:200" json:"desc"`
	Sort      int       `gorm:"default:0" json:"sort"`
	CreatedAt time.Time `json:"created_at"`
}

type Tag struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Name      string    `gorm:"uniqueIndex;size:50;not null" json:"name"`
	Color     string    `gorm:"size:7" json:"color"`
	Count     int       `gorm:"default:0" json:"count"`
	CreatedAt time.Time `json:"created_at"`
}

type Like struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	UserID     uint      `gorm:"index:idx_user_target,unique" json:"user_id"`
	TargetType string    `gorm:"size:20" json:"target_type"`
	TargetID   uint      `gorm:"index:idx_user_target,unique" json:"target_id"`
	CreatedAt  time.Time `json:"created_at"`
}

type Favorite struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	UserID    uint      `gorm:"index:idx_user_article,unique" json:"user_id"`
	ArticleID uint      `gorm:"index:idx_user_article,unique" json:"article_id"`
	CreatedAt time.Time `json:"created_at"`
}
