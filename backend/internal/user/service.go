package user

import (
	"errors"

	"inference-engine/internal/admin"
	"inference-engine/internal/notify"
	"inference-engine/internal/pkg/jwt"
	"inference-engine/internal/pkg/points"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Service struct {
	repo        *Repository
	notifySvc   *notify.Service
	auditLogSvc *admin.AuditLogService
	pointsSvc   *points.Service
	db          *gorm.DB
}

func NewService(repo *Repository, db *gorm.DB) *Service {
	return &Service{
		repo:        repo,
		notifySvc:   notify.NewService(db),
		auditLogSvc: admin.NewAuditLogService(db),
		pointsSvc:   points.NewService(db),
		db:          db,
	}
}

type RegisterReq struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=32"`
}

type LoginReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResp struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

func (s *Service) Register(req *RegisterReq, ip string) (*User, error) {
	// Check if registration is allowed
	var allowRegister struct{ Value string }
	s.db.Table("system_configs").Where("`key` = ?", "allow_register").Select("value").Scan(&allowRegister)
	if allowRegister.Value == "false" {
		return nil, errors.New("registration is currently disabled")
	}

	existing, _ := s.repo.FindByEmail(req.Email)
	if existing != nil {
		return nil, errors.New("email already exists")
	}

	existing, _ = s.repo.FindByUsername(req.Username)
	if existing != nil {
		return nil, errors.New("username already exists")
	}

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPwd),
		Role:     "user",
		Level:    1,
		Status:   1,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	go s.auditLogSvc.LogUserAction(user.ID, "用户注册", "user", user.Username, ip)
	go s.pointsSvc.AwardPoints(user.ID, "register")

	return user, nil
}

func (s *Service) Login(req *LoginReq, ip, userAgent string) (*LoginResp, error) {
	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		// Record failed login attempt
		s.recordLoginLog(0, "", ip, userAgent, 0)
		return nil, errors.New("invalid email or password")
	}

	if user.Status == 0 {
		s.recordLoginLog(user.ID, user.Username, ip, userAgent, 0)
		return nil, errors.New("account has been banned")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		s.recordLoginLog(user.ID, user.Username, ip, userAgent, 0)
		return nil, errors.New("invalid email or password")
	}

	token, err := jwt.GenerateToken(user.ID, user.Role)
	if err != nil {
		return nil, err
	}

	// Record successful login
	s.recordLoginLog(user.ID, user.Username, ip, userAgent, 1)

	go s.auditLogSvc.LogUserAction(user.ID, "用户登录", "user", user.Username, ip)
	go s.pointsSvc.AwardPoints(user.ID, "daily_login")

	return &LoginResp{
		Token: token,
		User:  *user,
	}, nil
}

func (s *Service) recordLoginLog(userID uint, username, ip, userAgent string, status int) {
	deviceType, browser, osName := admin.ParseUserAgent(userAgent)
	go func() {
		s.db.Create(&admin.LoginLog{
			UserID:     userID,
			Username:   username,
			IP:         ip,
			UserAgent:  userAgent,
			DeviceType: deviceType,
			Browser:    browser,
			OS:         osName,
			Status:     status,
		})
	}()
}

// RecordLogout updates the latest login log with logout time
func (s *Service) RecordLogout(userID uint) {
	now := s.db.NowFunc()
	go func() {
		s.db.Model(&admin.LoginLog{}).
			Where("user_id = ? AND logout_at IS NULL", userID).
			Order("created_at DESC").
			Limit(1).
			Update("logout_at", now)
	}()
}

func (s *Service) GetProfile(userID uint) (*User, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (s *Service) GetUserByID(userID uint) (*User, error) {
	return s.repo.FindByID(userID)
}

type UpdateProfileReq struct {
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	Bio      string `json:"bio"`
}

func (s *Service) UpdateProfile(userID uint, req *UpdateProfileReq, ip string) (*User, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if req.Username != "" {
		existing, _ := s.repo.FindByUsername(req.Username)
		if existing != nil && existing.ID != userID {
			return nil, errors.New("username already exists")
		}
		user.Username = req.Username
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}
	if req.Bio != "" {
		user.Bio = req.Bio
	}

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}

	go s.auditLogSvc.LogUserAction(userID, "更新资料", "user", user.Username, ip)

	return user, nil
}

func (s *Service) Follow(followerID, followedID uint, ip string) error {
	if followerID == followedID {
		return errors.New("cannot follow yourself")
	}

	followedUser, err := s.repo.FindByID(followedID)
	if err != nil {
		return errors.New("user not found")
	}

	following, _ := s.repo.IsFollowing(followerID, followedID)
	if following {
		return errors.New("already following")
	}

	tx := s.db.Begin()

	follow := &Follow{
		FollowerID: followerID,
		FollowedID: followedID,
	}
	if err := s.repo.CreateFollow(follow); err != nil {
		tx.Rollback()
		return err
	}

	tx.Model(&User{}).Where("id = ?", followerID).UpdateColumn("follow_count", gorm.Expr("follow_count + 1"))
	tx.Model(&User{}).Where("id = ?", followedID).UpdateColumn("fans_count", gorm.Expr("fans_count + 1"))

	tx.Commit()

	go s.notifySvc.Create(&notify.CreateNotifyReq{
		UserID:   followedID,
		ActorID:  followerID,
		Type:     "follow",
		Content:  "关注了你",
		TargetID: followerID,
	})

	go s.auditLogSvc.LogUserAction(followerID, "关注用户", "user", followedUser.Username, ip)
	go s.pointsSvc.AwardPoints(followerID, "follow")

	return nil
}

func (s *Service) Unfollow(followerID, followedID uint, ip string) error {
	following, _ := s.repo.IsFollowing(followerID, followedID)
	if !following {
		return errors.New("not following")
	}

	followedUser, _ := s.repo.FindByID(followedID)

	tx := s.db.Begin()

	if err := s.repo.DeleteFollow(followerID, followedID); err != nil {
		tx.Rollback()
		return err
	}

	tx.Model(&User{}).Where("id = ?", followerID).UpdateColumn("follow_count", gorm.Expr("follow_count - 1"))
	tx.Model(&User{}).Where("id = ?", followedID).UpdateColumn("fans_count", gorm.Expr("fans_count - 1"))

	tx.Commit()

	if followedUser != nil {
		go s.auditLogSvc.LogUserAction(followerID, "取消关注", "user", followedUser.Username, ip)
	}

	return nil
}

func (s *Service) GetFollowers(userID uint, page, pageSize int) ([]User, int64, error) {
	return s.repo.GetFollowers(userID, page, pageSize)
}

func (s *Service) GetFollowing(userID uint, page, pageSize int) ([]User, int64, error) {
	return s.repo.GetFollowing(userID, page, pageSize)
}

func (s *Service) IsFollowing(followerID, followedID uint) (bool, error) {
	return s.repo.IsFollowing(followerID, followedID)
}