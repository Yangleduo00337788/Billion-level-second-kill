package admin

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"inference-engine/internal/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	db *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{db: db}
}

func (h *Handler) CreateAuditLog(userID uint, username, action, target, detail, ip string) {
	h.db.Create(&AuditLog{
		UserID:   userID,
		Username: username,
		Action:   action,
		Target:   target,
		Detail:   detail,
		IP:       ip,
	})
}

type DashboardStats struct {
	TotalUsers    int64 `json:"total_users"`
	TotalArticles int64 `json:"total_articles"`
	TotalComments int64 `json:"total_comments"`
	TotalPrompts  int64 `json:"total_prompts"`
	TodayUsers    int64 `json:"today_users"`
	TodayArticles int64 `json:"today_articles"`
}

func (h *Handler) Dashboard(c *gin.Context) {
	var stats DashboardStats
	h.db.Model(&struct{}{}).Table("users").Count(&stats.TotalUsers)
	h.db.Model(&struct{}{}).Table("articles").Count(&stats.TotalArticles)
	h.db.Model(&struct{}{}).Table("comments").Count(&stats.TotalComments)
	h.db.Model(&struct{}{}).Table("prompts").Count(&stats.TotalPrompts)
	h.db.Raw("SELECT COUNT(*) FROM users WHERE DATE(created_at) = CURDATE()").Scan(&stats.TodayUsers)
	h.db.Raw("SELECT COUNT(*) FROM articles WHERE DATE(created_at) = CURDATE()").Scan(&stats.TodayArticles)
	response.Success(c, stats)
}

func (h *Handler) ChartData(c *gin.Context) {
	type DayCount struct {
		Date  string `json:"date"`
		Count int64  `json:"count"`
	}
	var userTrend, articleTrend []DayCount
	h.db.Raw("SELECT DATE(created_at) as date, COUNT(*) as count FROM users WHERE created_at >= DATE_SUB(CURDATE(), INTERVAL 7 DAY) GROUP BY DATE(created_at) ORDER BY date").Scan(&userTrend)
	h.db.Raw("SELECT DATE(created_at) as date, COUNT(*) as count FROM articles WHERE created_at >= DATE_SUB(CURDATE(), INTERVAL 7 DAY) GROUP BY DATE(created_at) ORDER BY date").Scan(&articleTrend)
	response.Success(c, gin.H{"user_trend": userTrend, "article_trend": articleTrend})
}

func (h *Handler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var users []map[string]interface{}
	var total int64
	h.db.Table("users").Count(&total)
	h.db.Table("users").Select("id, username, email, role, level, status, article_count, follow_count, fans_count, created_at, updated_at").Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&users)
	response.Page(c, users, total, page, pageSize)
}

func (h *Handler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, response.ErrBadRequest, "invalid id")
		return
	}

	var req struct {
		Role   string `json:"role"`
		Status *int   `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.ErrBadRequest, err.Error())
		return
	}

	updates := map[string]interface{}{}
	if req.Role != "" {
		updates["role"] = req.Role
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if len(updates) > 0 {
		h.db.Table("users").Where("id = ?", id).Updates(updates)
	}
	response.Success(c, nil)
}

func (h *Handler) BanUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, response.ErrBadRequest, "invalid id")
		return
	}

	var req struct {
		Reason  string `json:"reason" binding:"required"`
		BanType string `json:"ban_type" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.ErrBadRequest, err.Error())
		return
	}

	ban := UserBan{UserID: uint(id), Reason: req.Reason, BanType: req.BanType}
	if req.BanType == "temporary" {
		expires := time.Now().Add(7 * 24 * time.Hour)
		ban.ExpiresAt = &expires
	}
	h.db.Create(&ban)
	h.db.Table("users").Where("id = ?", id).Update("status", 0)
	response.Success(c, nil)
}

func (h *Handler) UnbanUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, response.ErrBadRequest, "invalid id")
		return
	}
	h.db.Table("users").Where("id = ?", id).Update("status", 1)
	h.db.Where("user_id = ?", id).Delete(&UserBan{})
	response.Success(c, nil)
}

func (h *Handler) ListUserTags(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, response.ErrBadRequest, "invalid id")
		return
	}
	var tags []UserTag
	h.db.Where("user_id = ?", id).Find(&tags)
	response.Success(c, tags)
}

func (h *Handler) AddUserTag(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, response.ErrBadRequest, "invalid id")
		return
	}
	var req struct {
		Tag string `json:"tag" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.ErrBadRequest, err.Error())
		return
	}
	tag := UserTag{UserID: uint(id), Tag: req.Tag}
	h.db.Create(&tag)
	response.Success(c, tag)
}

func (h *Handler) DeleteUserTag(c *gin.Context) {
	id := c.Param("tagId")
	h.db.Delete(&UserTag{}, id)
	response.Success(c, nil)
}

func (h *Handler) ListArticles(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var articles []map[string]interface{}
	var total int64
	h.db.Table("articles").Count(&total)
	h.db.Table("articles").Select("id, user_id, title, status, view_count, like_count, comment_count, is_ai, created_at, updated_at").Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&articles)
	response.Page(c, articles, total, page, pageSize)
}

func (h *Handler) DeleteArticle(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, response.ErrBadRequest, "invalid id")
		return
	}
	h.db.Table("articles").Where("id = ?", id).Delete(nil)
	response.Success(c, nil)
}

func (h *Handler) ModerateArticle(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, response.ErrBadRequest, "invalid id")
		return
	}
	var req struct {
		Status int `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.ErrBadRequest, err.Error())
		return
	}
	h.db.Table("articles").Where("id = ?", id).Update("status", req.Status)
	response.Success(c, nil)
}

func (h *Handler) ListPrompts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var prompts []map[string]interface{}
	var total int64
	h.db.Table("prompts").Count(&total)
	h.db.Table("prompts").Select("id, user_id, title, category, model, usage_count, like_count, status, created_at").Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&prompts)
	response.Page(c, prompts, total, page, pageSize)
}

func (h *Handler) DeletePrompt(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, response.ErrBadRequest, "invalid id")
		return
	}
	h.db.Table("prompts").Where("id = ?", id).Delete(nil)
	response.Success(c, nil)
}

func (h *Handler) ModeratePrompt(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, response.ErrBadRequest, "invalid id")
		return
	}
	var req struct {
		Status int `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.ErrBadRequest, err.Error())
		return
	}
	h.db.Table("prompts").Where("id = ?", id).Update("status", req.Status)
	response.Success(c, nil)
}

func (h *Handler) ListComments(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var comments []map[string]interface{}
	var total int64
	h.db.Table("comments").Count(&total)
	h.db.Table("comments").Select("id, user_id, article_id, content, created_at").Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&comments)
	response.Page(c, comments, total, page, pageSize)
}

func (h *Handler) DeleteComment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, response.ErrBadRequest, "invalid id")
		return
	}
	h.db.Table("comments").Where("id = ?", id).Delete(nil)
	response.Success(c, nil)
}

// ===== Announcements =====
func (h *Handler) ListAnnouncements(c *gin.Context) {
	var items []Announcement
	h.db.Order("priority DESC, created_at DESC").Find(&items)
	response.Success(c, items)
}

func (h *Handler) CreateAnnouncement(c *gin.Context) {
	var req struct {
		Title    string `json:"title" binding:"required"`
		Content  string `json:"content" binding:"required"`
		Status   int    `json:"status"`
		Priority int    `json:"priority"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.ErrBadRequest, err.Error())
		return
	}
	item := Announcement{Title: req.Title, Content: req.Content, Status: req.Status, Priority: req.Priority}
	h.db.Create(&item)
	response.Success(c, item)
}

func (h *Handler) UpdateAnnouncement(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Title    string `json:"title"`
		Content  string `json:"content"`
		Status   *int   `json:"status"`
		Priority *int   `json:"priority"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.ErrBadRequest, err.Error())
		return
	}
	updates := map[string]interface{}{}
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Content != "" {
		updates["content"] = req.Content
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.Priority != nil {
		updates["priority"] = *req.Priority
	}
	h.db.Model(&Announcement{}).Where("id = ?", id).Updates(updates)
	response.Success(c, nil)
}

func (h *Handler) DeleteAnnouncement(c *gin.Context) {
	id := c.Param("id")
	h.db.Delete(&Announcement{}, id)
	response.Success(c, nil)
}

// ===== Reports =====
func (h *Handler) ListReports(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var items []Report
	var total int64
	h.db.Model(&Report{}).Count(&total)
	h.db.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&items)
	response.Page(c, items, total, page, pageSize)
}

func (h *Handler) HandleReport(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Status int    `json:"status"`
		Result string `json:"result"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.ErrBadRequest, err.Error())
		return
	}
	h.db.Model(&Report{}).Where("id = ?", id).Updates(map[string]interface{}{"status": req.Status, "result": req.Result})
	response.Success(c, nil)
}

func (h *Handler) CreateReport(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var req struct {
		TargetType string `json:"target_type" binding:"required"`
		TargetID   uint   `json:"target_id" binding:"required"`
		Reason     string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.ErrBadRequest, err.Error())
		return
	}
	report := Report{ReporterID: userID.(uint), TargetType: req.TargetType, TargetID: req.TargetID, Reason: req.Reason}
	h.db.Create(&report)
	response.Success(c, nil)
}

// ===== Audit Logs =====
func (h *Handler) ListAuditLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var items []AuditLog
	var total int64
	h.db.Model(&AuditLog{}).Count(&total)
	h.db.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&items)
	response.Page(c, items, total, page, pageSize)
}

// ===== Content Reviews =====
func (h *Handler) ListContentReviews(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status := c.DefaultQuery("status", "")
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var items []ContentReview
	var total int64
	query := h.db.Model(&ContentReview{})
	if status != "" {
		query = query.Where("status = ?", status)
	}
	query.Count(&total)
	query.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&items)
	response.Page(c, items, total, page, pageSize)
}

func (h *Handler) ReviewContent(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Status     int    `json:"status"`
		ReviewerID uint   `json:"reviewer_id"`
		Reason     string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.ErrBadRequest, err.Error())
		return
	}
	h.db.Model(&ContentReview{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":      req.Status,
		"reviewer_id": req.ReviewerID,
		"reason":      req.Reason,
	})
	response.Success(c, nil)
}

// ===== Categories =====
func (h *Handler) CreateCategory(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.ErrBadRequest, err.Error())
		return
	}
	h.db.Table("categories").Create(map[string]interface{}{"name": req.Name, "description": req.Description})
	response.Success(c, nil)
}

func (h *Handler) DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	h.db.Table("categories").Where("id = ?", id).Delete(nil)
	response.Success(c, nil)
}

// ===== System Config =====
func (h *Handler) ListConfigs(c *gin.Context) {
	var items []SystemConfig
	h.db.Find(&items)
	response.Success(c, items)
}

func (h *Handler) UpdateConfig(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Value string `json:"value"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.ErrBadRequest, err.Error())
		return
	}
	h.db.Model(&SystemConfig{}).Where("id = ?", id).Update("value", req.Value)
	response.Success(c, nil)
}

// ===== AI Usage Stats =====
func (h *Handler) AIUsageStats(c *gin.Context) {
	var stats []AIUsageStat
	h.db.Order("date DESC").Limit(30).Find(&stats)

	var totalRequests, totalTokens int64
	h.db.Model(&AIUsageStat{}).Select("COALESCE(SUM(total_requests),0)").Scan(&totalRequests)
	h.db.Model(&AIUsageStat{}).Select("COALESCE(SUM(total_tokens),0)").Scan(&totalTokens)

	response.Success(c, gin.H{"stats": stats, "total_requests": totalRequests, "total_tokens": totalTokens})
}

// ===== Sensitive Words =====
func (h *Handler) ListSensitiveWords(c *gin.Context) {
	var items []SensitiveWord
	h.db.Find(&items)
	response.Success(c, items)
}

func (h *Handler) AddSensitiveWord(c *gin.Context) {
	var req struct {
		Word      string `json:"word" binding:"required"`
		Level     int    `json:"level"`
		ReplaceTo string `json:"replace_to"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.ErrBadRequest, err.Error())
		return
	}
	item := SensitiveWord{Word: req.Word, Level: req.Level, ReplaceTo: req.ReplaceTo}
	if item.Level == 0 {
		item.Level = 1
	}
	h.db.Create(&item)
	response.Success(c, item)
}

func (h *Handler) DeleteSensitiveWord(c *gin.Context) {
	id := c.Param("id")
	h.db.Delete(&SensitiveWord{}, id)
	response.Success(c, nil)
}

// ===== Recommend Items =====
func (h *Handler) ListRecommendItems(c *gin.Context) {
	var items []RecommendItem
	h.db.Where("status = 1").Order("sort_order ASC").Find(&items)
	response.Success(c, items)
}

func (h *Handler) AddRecommendItem(c *gin.Context) {
	var req struct {
		TargetType string `json:"target_type" binding:"required"`
		TargetID   uint   `json:"target_id" binding:"required"`
		Position   string `json:"position" binding:"required"`
		SortOrder  int    `json:"sort_order"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.ErrBadRequest, err.Error())
		return
	}
	item := RecommendItem{TargetType: req.TargetType, TargetID: req.TargetID, Position: req.Position, SortOrder: req.SortOrder, Status: 1}
	h.db.Create(&item)
	response.Success(c, item)
}

func (h *Handler) DeleteRecommendItem(c *gin.Context) {
	id := c.Param("id")
	h.db.Delete(&RecommendItem{}, id)
	response.Success(c, nil)
}

// ===== Points Rules =====
func (h *Handler) ListPointsRules(c *gin.Context) {
	var rules []PointsRule
	h.db.Find(&rules)
	response.Success(c, rules)
}

func (h *Handler) UpdatePointsRule(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Points int    `json:"points"`
		Desc   string `json:"desc"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.ErrBadRequest, err.Error())
		return
	}
	h.db.Model(&PointsRule{}).Where("id = ?", id).Updates(map[string]interface{}{"points": req.Points, "desc": req.Desc})
	response.Success(c, nil)
}

// ===== Invite Codes =====
func (h *Handler) ListInviteCodes(c *gin.Context) {
	var codes []InviteCode
	h.db.Find(&codes)
	response.Success(c, codes)
}

func (h *Handler) CreateInviteCode(c *gin.Context) {
	var req struct {
		MaxUses   int        `json:"max_uses"`
		ExpiresAt *time.Time `json:"expires_at"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.ErrBadRequest, err.Error())
		return
	}
	code := InviteCode{Code: generateCode(), MaxUses: req.MaxUses, ExpiresAt: req.ExpiresAt}
	h.db.Create(&code)
	response.Success(c, code)
}

func generateCode() string {
	return time.Now().Format("01021504") + strconv.Itoa(int(time.Now().UnixNano()%10000))
}

// ===== Login Logs =====
func (h *Handler) ListLoginLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var items []LoginLog
	var total int64
	h.db.Model(&LoginLog{}).Count(&total)
	h.db.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&items)
	response.Page(c, items, total, page, pageSize)
}

// ===== IP Blacklist =====
func (h *Handler) ListIPBlacklist(c *gin.Context) {
	var items []IPBlacklist
	h.db.Find(&items)
	response.Success(c, items)
}

func (h *Handler) AddIPBlacklist(c *gin.Context) {
	var req struct {
		IP     string `json:"ip" binding:"required"`
		Reason string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.ErrBadRequest, err.Error())
		return
	}
	item := IPBlacklist{IP: req.IP, Reason: req.Reason}
	h.db.Create(&item)
	response.Success(c, item)
}

func (h *Handler) DeleteIPBlacklist(c *gin.Context) {
	id := c.Param("id")
	h.db.Delete(&IPBlacklist{}, id)
	response.Success(c, nil)
}

// ===== System Logs =====
func (h *Handler) ListSystemLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	level := c.DefaultQuery("level", "")
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var items []SystemLog
	var total int64
	query := h.db.Model(&SystemLog{})
	if level != "" {
		query = query.Where("level = ?", level)
	}
	query.Count(&total)
	query.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&items)
	response.Page(c, items, total, page, pageSize)
}

// ===== Page Views =====
func (h *Handler) PageViewStats(c *gin.Context) {
	var totalPV, todayPV, totalUV, todayUV int64
	h.db.Model(&PageView{}).Count(&totalPV)
	h.db.Raw("SELECT COUNT(*) FROM page_views WHERE DATE(created_at) = CURDATE()").Scan(&todayPV)
	h.db.Raw("SELECT COUNT(DISTINCT ip) FROM page_views").Scan(&totalUV)
	h.db.Raw("SELECT COUNT(DISTINCT ip) FROM page_views WHERE DATE(created_at) = CURDATE()").Scan(&todayUV)

	type TopPage struct {
		Path  string `json:"path"`
		Count int64  `json:"count"`
	}
	var topPages []TopPage
	h.db.Model(&PageView{}).Select("path, COUNT(*) as count").Group("path").Order("count DESC").Limit(10).Find(&topPages)

	response.Success(c, gin.H{
		"total_pv":  totalPV,
		"today_pv":  todayPV,
		"total_uv":  totalUV,
		"today_uv":  todayUV,
		"top_pages": topPages,
	})
}

// ===== Rankings =====
func (h *Handler) HotArticles(c *gin.Context) {
	var articles []map[string]interface{}
	h.db.Table("articles").Select("id, title, view_count, like_count, comment_count").Where("status = 1").Order("view_count DESC").Limit(10).Find(&articles)
	response.Success(c, articles)
}

func (h *Handler) HotPrompts(c *gin.Context) {
	var prompts []map[string]interface{}
	h.db.Table("prompts").Select("id, title, usage_count, like_count").Where("status = 1").Order("usage_count DESC").Limit(10).Find(&prompts)
	response.Success(c, prompts)
}

func (h *Handler) ActiveUsers(c *gin.Context) {
	var users []map[string]interface{}
	h.db.Table("users").Select("id, username, article_count, follow_count, fans_count").Where("status = 1").Order("article_count DESC").Limit(10).Find(&users)
	response.Success(c, users)
}

// ===== AI Config =====
func (h *Handler) GetAIConfig(c *gin.Context) {
	getVal := func(key, def string) string {
		var sc SystemConfig
		if err := h.db.Where("`key` = ?", key).First(&sc).Error; err != nil {
			return def
		}
		if sc.Value == "" {
			return def
		}
		return sc.Value
	}

	maskKey := func(key string) string {
		if key == "" {
			return ""
		}
		if len(key) > 8 {
			return key[:4] + "****" + key[len(key)-4:]
		}
		return "****"
	}

	cfg := map[string]string{
		"ai_provider":    getVal("ai_provider", "openai"),
		"ai_api_key":     maskKey(getVal("ai_api_key", "")),
		"ai_base_url":    getVal("ai_base_url", ""),
		"ai_model":       getVal("ai_model", "gpt-4o"),
		"ai_max_tokens":  getVal("ai_max_tokens", "2048"),
		"ai_temperature": getVal("ai_temperature", "0.7"),
		"ai_rate_limit":  getVal("ai_rate_limit", "60"),
	}
	response.Success(c, cfg)
}

func (h *Handler) UpdateAIConfig(c *gin.Context) {
	var req map[string]string
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.ErrBadRequest, err.Error())
		return
	}

	allowedKeys := map[string]string{
		"ai_provider":    "AI 服务商 (openai/deepseek)",
		"ai_api_key":     "API 密钥",
		"ai_base_url":    "API 基础地址",
		"ai_model":       "模型名称",
		"ai_max_tokens":  "最大 Token 数",
		"ai_temperature": "Temperature",
		"ai_rate_limit":  "每分钟请求限制",
	}

	for key, value := range req {
		if _, ok := allowedKeys[key]; !ok {
			continue
		}
		// 跳过掩码值
		if key == "ai_api_key" && strings.Contains(value, "****") {
			continue
		}
		// 校验
		if key == "ai_max_tokens" {
			if n, err := strconv.Atoi(value); err != nil || n <= 0 {
				response.Error(c, response.ErrBadRequest, "ai_max_tokens 必须为正整数")
				return
			}
		}
		if key == "ai_temperature" {
			if f, err := strconv.ParseFloat(value, 64); err != nil || f < 0 || f > 2 {
				response.Error(c, response.ErrBadRequest, "ai_temperature 范围 0-2")
				return
			}
		}
		if key == "ai_rate_limit" {
			if n, err := strconv.Atoi(value); err != nil || n <= 0 {
				response.Error(c, response.ErrBadRequest, "ai_rate_limit 必须为正整数")
				return
			}
		}

		var sc SystemConfig
		if err := h.db.Where("`key` = ?", key).First(&sc).Error; err != nil {
			desc := allowedKeys[key]
			h.db.Create(&SystemConfig{Key: key, Value: value, Desc: desc})
		} else {
			h.db.Model(&sc).Update("value", value)
		}
	}

	response.Success(c, nil)
}

func (h *Handler) TestAIConnection(c *gin.Context) {
	getVal := func(key string) string {
		var sc SystemConfig
		if err := h.db.Where("`key` = ?", key).First(&sc).Error; err != nil {
			return ""
		}
		return sc.Value
	}

	apiKey := getVal("ai_api_key")
	provider := getVal("ai_provider")

	if apiKey == "" {
		response.Success(c, gin.H{"connected": false, "message": "未配置 API Key"})
		return
	}

	response.Success(c, gin.H{
		"connected": true,
		"provider":  provider,
		"message":   fmt.Sprintf("已配置 %s 服务", provider),
	})
}
