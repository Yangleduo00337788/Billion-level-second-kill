package comment

import (
	"errors"

	"inference-engine/internal/admin"
	"inference-engine/internal/article"
	"inference-engine/internal/notify"

	"gorm.io/gorm"
)

type Service struct {
	repo        *Repository
	articleRepo *article.Repository
	notifySvc   *notify.Service
	auditLogSvc *admin.AuditLogService
	db          *gorm.DB
}

func NewService(repo *Repository, articleRepo *article.Repository, db *gorm.DB) *Service {
	return &Service{
		repo:        repo,
		articleRepo: articleRepo,
		notifySvc:   notify.NewService(db),
		auditLogSvc: admin.NewAuditLogService(db),
		db:          db,
	}
}

type CreateCommentReq struct {
	Content  string `json:"content" binding:"required"`
	ParentID uint   `json:"parent_id"`
}

func (s *Service) Create(userID, articleID uint, req *CreateCommentReq) (*Comment, error) {
	art, err := s.articleRepo.FindByID(articleID)
	if err != nil {
		return nil, errors.New("article not found")
	}

	comment := &Comment{
		ArticleID: articleID,
		UserID:    userID,
		ParentID:  req.ParentID,
		Content:   req.Content,
	}

	if err := s.repo.Create(comment); err != nil {
		return nil, err
	}

	s.articleRepo.IncrementCommentCount(articleID)

	// Trigger notification to article author
	if art.UserID != userID {
		go s.notifySvc.Create(&notify.CreateNotifyReq{
			UserID:   art.UserID,
			ActorID:  userID,
			Type:     "comment",
			Content:  "评论了你的文章",
			TargetID: articleID,
		})
	}

	// Audit log
	go s.auditLogSvc.LogUserAction(userID, "评论文章", "article", art.Title)

	return comment, nil
}

func (s *Service) GetByArticleID(articleID uint, page, pageSize int) ([]Comment, int64, error) {
	return s.repo.GetByArticleID(articleID, page, pageSize)
}

func (s *Service) Delete(commentID, userID uint) error {
	comment, err := s.repo.FindByID(commentID)
	if err != nil {
		return errors.New("comment not found")
	}

	if comment.UserID != userID {
		return errors.New("permission denied")
	}

	if err := s.repo.Delete(commentID); err != nil {
		return err
	}

	s.articleRepo.DecrementCommentCount(comment.ArticleID)

	// Audit log
	go s.auditLogSvc.LogUserAction(userID, "删除评论", "comment", "评论ID: "+string(rune(commentID)))

	return nil
}
