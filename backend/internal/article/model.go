package article

import (
	"encoding/json"
	"time"

	"inference-engine/internal/user"

	"gorm.io/gorm"
)

type Article struct {
	ID            uint           `gorm:"primarykey;comment:文章ID" json:"id"`
	UserID        uint           `gorm:"index;not null;comment:作者ID" json:"user_id"`
	Title         string         `gorm:"size:200;not null;comment:文章标题" json:"title"`
	Content       string         `gorm:"type:longtext;comment:文章内容(Markdown)" json:"content"`
	Summary       string         `gorm:"size:500;comment:文章摘要" json:"summary"`
	Cover         string         `gorm:"size:255;comment:封面图URL" json:"cover"`
	Status        string         `gorm:"size:20;default:draft;comment:状态(draft/published)" json:"status"`
	CategoryID    uint           `gorm:"index;comment:分类ID" json:"category_id"`
	Tags          string         `gorm:"type:json;comment:标签列表(JSON)" json:"tags"`
	ViewCount     int            `gorm:"default:0;comment:浏览量" json:"view_count"`
	LikeCount     int            `gorm:"default:0;comment:点赞数" json:"like_count"`
	CommentCount  int            `gorm:"default:0;comment:评论数" json:"comment_count"`
	FavoriteCount int            `gorm:"default:0;comment:收藏数" json:"favorite_count"`
	IsAI          bool           `gorm:"default:false;comment:是否AI生成" json:"is_ai"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index;comment:软删除时间" json:"-"`
	User          user.User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Category      Category       `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
}

func (Article) TableComment() string {
	return "文章表"
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
	ID        uint      `gorm:"primarykey;comment:分类ID" json:"id"`
	Name      string    `gorm:"size:50;not null;comment:分类名称" json:"name"`
	Desc      string    `gorm:"size:200;comment:分类描述" json:"desc"`
	Sort      int       `gorm:"default:0;comment:排序" json:"sort"`
	CreatedAt time.Time `json:"created_at"`
}

func (Category) TableComment() string {
	return "文章分类表"
}

type Tag struct {
	ID        uint      `gorm:"primarykey;comment:标签ID" json:"id"`
	Name      string    `gorm:"uniqueIndex;size:50;not null;comment:标签名称" json:"name"`
	Color     string    `gorm:"size:7;comment:标签颜色(HEX)" json:"color"`
	Count     int       `gorm:"default:0;comment:使用次数" json:"count"`
	CreatedAt time.Time `json:"created_at"`
}

func (Tag) TableComment() string {
	return "标签表"
}

type Like struct {
	ID         uint      `gorm:"primarykey;comment:点赞ID" json:"id"`
	UserID     uint      `gorm:"index:idx_user_target,unique;comment:用户ID" json:"user_id"`
	TargetType string    `gorm:"size:20;comment:目标类型(article/comment)" json:"target_type"`
	TargetID   uint      `gorm:"index:idx_user_target,unique;comment:目标ID" json:"target_id"`
	CreatedAt  time.Time `json:"created_at"`
}

func (Like) TableComment() string {
	return "点赞表"
}

type Favorite struct {
	ID        uint      `gorm:"primarykey;comment:收藏ID" json:"id"`
	UserID    uint      `gorm:"index:idx_user_article,unique;comment:用户ID" json:"user_id"`
	ArticleID uint      `gorm:"index:idx_user_article,unique;comment:文章ID" json:"article_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (Favorite) TableComment() string {
	return "收藏表"
}
