package user

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           uint           `gorm:"primarykey;comment:用户ID" json:"id"`
	Username     string         `gorm:"uniqueIndex;size:50;not null;comment:用户名" json:"username"`
	Email        string         `gorm:"uniqueIndex;size:100;not null;comment:邮箱" json:"email"`
	Password     string         `gorm:"not null;comment:密码(加密)" json:"-"`
	Avatar       string         `gorm:"size:255;comment:头像URL" json:"avatar"`
	Bio          string         `gorm:"size:500;comment:个人简介" json:"bio"`
	Role         string         `gorm:"size:20;default:user;comment:角色(user/creator/admin)" json:"role"`
	Level        int            `gorm:"default:1;comment:等级" json:"level"`
	Status       int            `gorm:"default:1;comment:状态(1=正常 0=禁用)" json:"status"`
	FollowCount  int            `gorm:"default:0;comment:关注数" json:"follow_count"`
	FansCount    int            `gorm:"default:0;comment:粉丝数" json:"fans_count"`
	ArticleCount int            `gorm:"default:0;comment:文章数" json:"article_count"`
	Points       int            `gorm:"default:0;comment:积分" json:"points"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index;comment:软删除时间" json:"-"`
}

func (User) TableComment() string {
	return "用户表"
}

type Follow struct {
	ID         uint           `gorm:"primarykey;comment:关注ID" json:"id"`
	FollowerID uint           `gorm:"index:idx_follower;not null;comment:关注者ID" json:"follower_id"`
	FollowedID uint           `gorm:"index:idx_followed;not null;comment:被关注者ID" json:"followed_id"`
	CreatedAt  time.Time      `json:"created_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index;comment:软删除时间" json:"-"`
}

func (Follow) TableComment() string {
	return "关注关系表"
}
