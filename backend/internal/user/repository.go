package user

import (
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(user *User) error {
	return r.db.Create(user).Error
}

func (r *Repository) FindByID(id uint) (*User, error) {
	var user User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) FindByEmail(email string) (*User, error) {
	var user User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) FindByUsername(username string) (*User, error) {
	var user User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) Update(user *User) error {
	return r.db.Save(user).Error
}

func (r *Repository) List(page, pageSize int) ([]User, int64, error) {
	var users []User
	var total int64

	r.db.Model(&User{}).Count(&total)
	err := r.db.Offset((page - 1) * pageSize).Limit(pageSize).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (r *Repository) FindByIDs(ids []uint) ([]User, error) {
	var users []User
	err := r.db.Where("id IN ?", ids).Find(&users).Error
	return users, err
}

func (r *Repository) CreateFollow(follow *Follow) error {
	return r.db.Create(follow).Error
}

// CreateFollowWithTx 使用事务创建关注关系
func (r *Repository) CreateFollowWithTx(tx *gorm.DB, follow *Follow) error {
	return tx.Create(follow).Error
}

func (r *Repository) DeleteFollow(followerID, followedID uint) error {
	return r.db.Where("follower_id = ? AND followed_id = ?", followerID, followedID).Delete(&Follow{}).Error
}

// DeleteFollowWithTx 使用事务删除关注关系
func (r *Repository) DeleteFollowWithTx(tx *gorm.DB, followerID, followedID uint) error {
	return tx.Where("follower_id = ? AND followed_id = ?", followerID, followedID).Delete(&Follow{}).Error
}

func (r *Repository) IsFollowing(followerID, followedID uint) (bool, error) {
	var count int64
	err := r.db.Model(&Follow{}).Where("follower_id = ? AND followed_id = ?", followerID, followedID).Count(&count).Error
	return count > 0, err
}

func (r *Repository) GetFollowers(userID uint, page, pageSize int) ([]User, int64, error) {
	var follows []Follow
	var total int64

	r.db.Model(&Follow{}).Where("followed_id = ?", userID).Count(&total)
	err := r.db.Where("followed_id = ?", userID).Offset((page - 1) * pageSize).Limit(pageSize).Find(&follows).Error
	if err != nil {
		return nil, 0, err
	}

	var userIDs []uint
	for _, f := range follows {
		userIDs = append(userIDs, f.FollowerID)
	}

	users, err := r.FindByIDs(userIDs)
	return users, total, err
}

func (r *Repository) GetFollowing(userID uint, page, pageSize int) ([]User, int64, error) {
	var follows []Follow
	var total int64

	r.db.Model(&Follow{}).Where("follower_id = ?", userID).Count(&total)
	err := r.db.Where("follower_id = ?", userID).Offset((page - 1) * pageSize).Limit(pageSize).Find(&follows).Error
	if err != nil {
		return nil, 0, err
	}

	var userIDs []uint
	for _, f := range follows {
		userIDs = append(userIDs, f.FollowedID)
	}

	users, err := r.FindByIDs(userIDs)
	return users, total, err
}
