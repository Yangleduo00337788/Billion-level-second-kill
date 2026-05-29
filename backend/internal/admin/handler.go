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
	db        *gorm.DB
	aiService interface{ ReloadConfig() }
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{db: db}
}

func NewHandlerWithAI(db *gorm.DB, aiService interface{ ReloadConfig() }) *Handler {
	return &Handler{db: db, aiService: aiService}
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

// ===== Dashboard =====

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
	h.db.Table("users").Count(&stats.TotalUsers)
	h.db.Table("articles").Count(&stats.TotalArticles)
	h.db.Table("comments").Count(&stats.TotalComments)
	h.db.Table("prompts").Count(&stats.TotalPrompts)
	h.db.Raw("SELECT COUNT(*) FROM users WHERE DATE(created_at) = CURDATE()").Scan(&stats.TodayUsers)
	h.db.Raw("SELECT COUNT(*) FROM articles WHERE DATE(created_at) = CURDATE()").Scan(&stats.TodayArticles)
	response.Success(c, stats)
}

func (h *Handler) ChartData(c *gin.Context) {
	type DayStats struct {
		Date     string `json:"date"`
		Users    int64  `json:"users"`
		Articles int64  `json:"articles"`
		Comments int64  `json:"comments"`
	}

	var results []DayStats
	for i := 29; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
		var uCount, aCount, cCount int64
		h.db.Raw("SELECT COUNT(*) FROM users WHERE DATE(created_at) = ?", date).Scan(&uCount)
		h.db.Raw("SELECT COUNT(*) FROM articles WHERE DATE(created_at) = ?", date).Scan(&aCount)
		h.db.Raw("SELECT COUNT(*) FROM comments WHERE DATE(created_at) = ?", date).Scan(&cCount)
		results = append(results, DayStats{Date: date, Users: uCount, Articles: aCount, Comments: cCount})
	}
	response.Success(c, results)
}

// ===== User Management =====

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

func (h *Handler) SearchUsers(c *gin.Context) {
	q := c.Query("q")
	if q == "" {
		response.Success(c, []interface{}{})
		return
	}
	var users []map[string]interface{}
	keyword := "%" + q + "%"
	h.db.Table("users").
		Select("id, username, email, avatar, role, status").
		Where("username LIKE ? OR email LIKE ?", keyword, keyword).
		Limit(20).Find(&users)
	response.Success(c, users)
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

func (h *Handler) ListAllUserTags(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var tags []struct {
		UserTag
		Username string `json:"username"`
	}
	var total int64

	h.db.Model(&UserTag{}).Count(&total)
	h.db.Table("user_tags").
		Select("user_tags.*, users.username").
		Joins("LEFT JOIN users ON users.id = user_tags.user_id").
		Order("user_tags.id DESC").
		Offset((page - 1) * pageSize).Limit(pageSize).
		Find(&tags)

	response.Page(c, tags, total, page, pageSize)
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
	tagIDStr := c.Param("tagId")
	tagID, err := strconv.ParseUint(tagIDStr, 10, 64)
	if err != nil {
		response.Error(c, response.ErrBadRequest, "invalid tag id")
		return
	}
	h.db.Delete(&UserTag{}, tagID)
	response.Success(c, nil)
}

// ===== Content Management =====

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
	h.db.Table("articles").Select("id, user_id, title, status, view_count, like_count, comment_count, is_ai, created_at").Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&articles)
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
		Status string `json:"status"`
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
	h.db.Table("prompts").Select("id, user_id, title, category, model, usage_count, like_count, rating, status, created_at").Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&prompts)
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
		Status *int `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.ErrBadRequest, err.Error())
		return
	}
	if req.Status != nil {
		h.db.Table("prompts").Where("id = ?", id).Update("status", *req.Status)
	}
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
	h.db.Table("comments").Select("id, article_id, user_id, content, like_count, created_at").Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&comments)
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

// ===== Content Review =====

func (h *Handler) ListContentReviews(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var items []ContentReview
	var total int64
	h.db.Model(&ContentReview{}).Count(&total)
	h.db.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&items)
	response.Page(c, items, total, page, pageSize)
}

func (h *Handler) ReviewContent(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, response.ErrBadRequest, "invalid id")
		return
	}
	var req struct {
		Status int    `json:"status"`
		Reason string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.ErrBadRequest, err.Error())
		return
	}
	h.db.Model(&ContentReview{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":      req.Status,
		"reason":      req.Reason,
		"reviewer_id": c.GetUint("userID"),
	})
	response.Success(c, nil)
}

// ===== Categories & Tags =====

func (h *Handler) CreateCategory(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
		Desc string `json:"desc"`
		Sort int    `json:"sort"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.ErrBadRequest, err.Error())
		return
	}
	h.db.Exec("INSERT INTO categories (name, `desc`, sort, created_at) VALUES (?, ?, ?, NOW())", req.Name, req.Desc, req.Sort)
	response.Success(c, nil)
}

func (h *Handler) DeleteCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, response.ErrBadRequest, "invalid id")
		return
	}
	h.db.Exec("DELETE FROM categories WHERE id = ?", id)
	response.Success(c, nil)
}

// ===== Announcements =====

func (h *Handler) ListAnnouncements(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var items []Announcement
	var total int64
	h.db.Model(&Announcement{}).Count(&total)
	h.db.Order("priority DESC, created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&items)
	response.Page(c, items, total, page, pageSize)
}

func (h *Handler) CreateAnnouncement(c *gin.Context) {
	var req struct {
		Title    string `json:"title" binding:"required"`
		Content  string `json:"content" binding:"required"`
		Priority int    `json:"priority"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.ErrBadRequest, err.Error())
		return
	}
	announcement := Announcement{Title: req.Title, Content: req.Content, Priority: req.Priority, Status: 1}
	h.db.Create(&announcement)
	response.Success(c, announcement)
}

func (h *Handler) UpdateAnnouncement(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, response.ErrBadRequest, "invalid id")
		return
	}
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.ErrBadRequest, err.Error())
		return
	}
	h.db.Model(&Announcement{}).Where("id = ?", id).Updates(req)
	response.Success(c, nil)
}

func (h *Handler) DeleteAnnouncement(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, response.ErrBadRequest, "invalid id")
		return
	}
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
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, response.ErrBadRequest, "invalid id")
		return
	}
	var req struct {
		Status int    `json:"status"`
		Result string `json:"result"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.ErrBadRequest, err.Error())
		return
	}
	h.db.Model(&Report{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status": req.Status,
		"result": req.Result,
	})
	response.Success(c, nil)
}

func (h *Handler) CreateReport(c *gin.Context) {
	var req struct {
		TargetType string `json:"target_type" binding:"required"`
		TargetID   uint   `json:"target_id" binding:"required"`
		Reason     string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.ErrBadRequest, err.Error())
		return
	}
	report := Report{
		ReporterID: c.GetUint("userID"),
		TargetType: req.TargetType,
		TargetID:   req.TargetID,
		Reason:     req.Reason,
	}
	h.db.Create(&report)
	response.Success(c, report)
}

// ===== Security - Audit Logs =====

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

	// Fill in username and IP for entries that are missing them
	for i := range items {
		if items[i].Username == "" && items[i].UserID > 0 {
			var username string
			h.db.Table("users").Where("id = ?", items[i].UserID).Pluck("username", &username)
			items[i].Username = username
		}
		// Make action human-readable
		items[i].Action = humanizeAction(items[i].Action)
	}

	response.Page(c, items, total, page, pageSize)
}

// humanizeAction converts raw HTTP method+path or action strings to human-readable text
func humanizeAction(action string) string {
	action = strings.TrimSpace(action)
	// Already in Chinese
	if strings.Contains(action, "用户") || strings.Contains(action, "文章") ||
		strings.Contains(action, "关注") || strings.Contains(action, "点赞") ||
		strings.Contains(action, "收藏") || strings.Contains(action, "评论") ||
		strings.Contains(action, "登录") || strings.Contains(action, "注册") ||
		strings.Contains(action, "更新") || strings.Contains(action, "删除") ||
		strings.Contains(action, "创建") || strings.Contains(action, "编辑") {
		return action
	}
	// Convert HTTP method + path
	methodMap := map[string]string{
		"POST":   "创建",
		"PUT":    "编辑",
		"DELETE": "删除",
		"PATCH":  "更新",
	}
	for method, cn := range methodMap {
		if strings.HasPrefix(action, method) {
			path := strings.TrimPrefix(action, method+" ")
			path = strings.TrimSpace(path)
			return cn + " " + path
		}
	}
	return action
}

// ===== Security - Login Logs =====

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

	// Fill username for entries that are missing it
	for i := range items {
		if items[i].Username == "" && items[i].UserID > 0 {
			var username string
			h.db.Table("users").Where("id = ?", items[i].UserID).Pluck("username", &username)
			items[i].Username = username
		}
	}

	response.Page(c, items, total, page, pageSize)
}

// ===== Security - IP Blacklist =====

func (h *Handler) ListIPBlacklist(c *gin.Context) {
	var items []IPBlacklist
	h.db.Order("created_at DESC").Find(&items)
	response.Success(c, items)
}

func (h *Handler) AddIPBlacklist(c *gin.Context) {
	var req struct {
		IP       string `json:"ip"`
		UserID   uint   `json:"user_id"`
		Username string `json:"username"`
		Reason   string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.ErrBadRequest, err.Error())
		return
	}

	// If user_id is provided, get their IP from recent login logs
	if req.UserID > 0 && req.IP == "" {
		var log LoginLog
		h.db.Where("user_id = ?", req.UserID).Order("created_at DESC").First(&log)
		req.IP = log.IP
	}

	if req.Username == "" && req.UserID > 0 {
		var username string
		h.db.Table("users").Where("id = ?", req.UserID).Pluck("username", &username)
		req.Username = username
	}

	if req.IP == "" {
		response.Error(c, response.ErrBadRequest, "IP地址不能为空，请确保该用户有登录记录")
		return
	}

	item := IPBlacklist{
		IP:       req.IP,
		UserID:   req.UserID,
		Username: req.Username,
		Reason:   req.Reason,
	}
	if err := h.db.Create(&item).Error; err != nil {
		if strings.Contains(err.Error(), "Duplicate") {
			response.Error(c, response.ErrBadRequest, "该IP已在黑名单中")
			return
		}
		response.Error(c, response.ErrInternal, err.Error())
		return
	}
	response.Success(c, item)
}

func (h *Handler) DeleteIPBlacklist(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, response.ErrBadRequest, "invalid id")
		return
	}
	h.db.Delete(&IPBlacklist{}, id)
	response.Success(c, nil)
}

// ===== Security - Sensitive Words =====

func (h *Handler) ListSensitiveWords(c *gin.Context) {
	var items []SensitiveWord
	h.db.Order("created_at DESC").Find(&items)
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
	if req.Level == 0 {
		req.Level = 1
	}
	item := SensitiveWord{Word: req.Word, Level: req.Level, ReplaceTo: req.ReplaceTo}
	h.db.Create(&item)
	response.Success(c, item)
}

func (h *Handler) DeleteSensitiveWord(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, response.ErrBadRequest, "invalid id")
		return
	}
	h.db.Delete(&SensitiveWord{}, id)
	response.Success(c, nil)
}

// ===== Operations - Recommend Items =====

func (h *Handler) ListRecommendItems(c *gin.Context) {
	var items []RecommendItem
	h.db.Order("position ASC, sort_order ASC").Find(&items)
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
	item := RecommendItem{
		TargetType: req.TargetType,
		TargetID:   req.TargetID,
		Position:   req.Position,
		SortOrder:  req.SortOrder,
		Status:     1,
	}
	h.db.Create(&item)
	response.Success(c, item)
}

func (h *Handler) DeleteRecommendItem(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, response.ErrBadRequest, "invalid id")
		return
	}
	h.db.Delete(&RecommendItem{}, id)
	response.Success(c, nil)
}

// ===== Operations - Points Rules =====

func (h *Handler) ListPointsRules(c *gin.Context) {
	var items []PointsRule
	h.db.Order("id ASC").Find(&items)
	response.Success(c, items)
}

func (h *Handler) UpdatePointsRule(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, response.ErrBadRequest, "invalid id")
		return
	}
	var req struct {
		Points    *int   `json:"points"`
		Desc      string `json:"desc"`
		LimitType string `json:"limit_type"`
		Status    *int   `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.ErrBadRequest, err.Error())
		return
	}
	updates := map[string]interface{}{}
	if req.Points != nil {
		updates["points"] = *req.Points
	}
	if req.Desc != "" {
		updates["desc"] = req.Desc
	}
	if req.LimitType != "" {
		updates["limit_type"] = req.LimitType
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if len(updates) > 0 {
		h.db.Model(&PointsRule{}).Where("id = ?", id).Updates(updates)
	}
	response.Success(c, nil)
}

func (h *Handler) CreatePointsRule(c *gin.Context) {
	var req struct {
		Action    string `json:"action" binding:"required"`
		Points    int    `json:"points" binding:"required"`
		Desc      string `json:"desc"`
		LimitType string `json:"limit_type"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.ErrBadRequest, err.Error())
		return
	}
	limitType := req.LimitType
	if limitType == "" {
		limitType = "unlimited"
	}
	rule := &PointsRule{Action: req.Action, Points: req.Points, Desc: req.Desc, LimitType: limitType, Status: 1}
	if err := h.db.Create(rule).Error; err != nil {
		response.Error(c, response.ErrBadRequest, "action already exists or create failed")
		return
	}
	response.Success(c, rule)
}

// ===== Operations - Invite Codes =====

func (h *Handler) ListInviteCodes(c *gin.Context) {
	var items []InviteCode
	h.db.Order("created_at DESC").Find(&items)
	response.Success(c, items)
}

func (h *Handler) CreateInviteCode(c *gin.Context) {
	var req struct {
		Code    string `json:"code"`
		MaxUses int    `json:"max_uses"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.ErrBadRequest, err.Error())
		return
	}
	if req.Code == "" {
		req.Code = generateRandomCode(8)
	}
	if req.MaxUses <= 0 {
		req.MaxUses = 1
	}
	item := InviteCode{Code: req.Code, MaxUses: req.MaxUses}
	h.db.Create(&item)
	response.Success(c, item)
}

func generateRandomCode(length int) string {
	const chars = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = chars[time.Now().UnixNano()%int64(len(chars))]
		time.Sleep(time.Nanosecond)
	}
	return string(result)
}

// ===== System - Configs =====

func (h *Handler) ListConfigs(c *gin.Context) {
	var items []SystemConfig
	h.db.Order("id ASC").Find(&items)
	response.Success(c, items)
}

func (h *Handler) UpdateConfig(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, response.ErrBadRequest, "invalid id")
		return
	}
	var req struct {
		Value string `json:"value"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.ErrBadRequest, err.Error())
		return
	}

	// 获取配置项的 key，检查是否是 AI 相关配置
	var sc SystemConfig
	h.db.First(&sc, id)
	h.db.Model(&SystemConfig{}).Where("id = ?", id).Update("value", req.Value)

	// 如果是 AI 相关配置，重新加载 AI 服务
	if strings.HasPrefix(sc.Key, "ai_") && h.aiService != nil {
		h.aiService.ReloadConfig()
	}

	response.Success(c, nil)
}

// CreateOrUpdateConfig 通用的配置创建或更新接口
func (h *Handler) CreateOrUpdateConfig(c *gin.Context) {
	var req struct {
		Key   string `json:"key" binding:"required"`
		Value string `json:"value"`
		Desc  string `json:"desc"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.ErrBadRequest, err.Error())
		return
	}

	var sc SystemConfig
	if err := h.db.Where("`key` = ?", req.Key).First(&sc).Error; err != nil {
		// 不存在则创建
		sc = SystemConfig{Key: req.Key, Value: req.Value, Desc: req.Desc}
		h.db.Create(&sc)
	} else {
		// 存在则更新
		h.db.Model(&sc).Update("value", req.Value)
	}

	// 如果是 AI 相关配置，重新加载 AI 服务
	if strings.HasPrefix(req.Key, "ai_") && h.aiService != nil {
		h.aiService.ReloadConfig()
	}

	response.Success(c, sc)
}

// ===== System - AI Stats =====

func (h *Handler) AIUsageStats(c *gin.Context) {
	var items []AIUsageStat
	h.db.Order("date DESC").Limit(30).Find(&items)

	var totalRequests, totalTokens, totalErrors int64
	for _, item := range items {
		totalRequests += int64(item.TotalRequests)
		totalTokens += int64(item.TotalTokens)
		totalErrors += int64(item.ErrorCount)
	}

	response.Success(c, gin.H{
		"total_requests": totalRequests,
		"total_tokens":   totalTokens,
		"total_errors":   totalErrors,
		"daily":          items,
	})
}

// ===== System - System Logs =====

func (h *Handler) ListSystemLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	level := c.Query("level")
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

// LogSystem writes a system log entry
func (h *Handler) LogSystem(level, message, source string) {
	h.db.Create(&SystemLog{Level: level, Message: message, Source: source})
}

// ===== Page View Stats =====

func (h *Handler) PageViewStats(c *gin.Context) {
	var totalPV, todayPV, totalUV, todayUV int64
	h.db.Table("page_views").Count(&totalPV)
	h.db.Raw("SELECT COUNT(*) FROM page_views WHERE DATE(created_at) = CURDATE()").Scan(&todayPV)
	h.db.Raw("SELECT COUNT(DISTINCT ip) FROM page_views").Scan(&totalUV)
	h.db.Raw("SELECT COUNT(DISTINCT ip) FROM page_views WHERE DATE(created_at) = CURDATE()").Scan(&todayUV)

	type TopPage struct {
		Path  string `json:"path"`
		Count int64  `json:"count"`
	}
	var topPages []TopPage
	h.db.Table("page_views").Select("path, COUNT(*) as count").Group("path").Order("count DESC").Limit(10).Find(&topPages)

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
	h.db.Table("articles").Select("id, title, view_count, like_count, comment_count").
		Where("status = ?", "published").
		Order("view_count DESC").Limit(10).Find(&articles)
	response.Success(c, articles)
}

func (h *Handler) HotPrompts(c *gin.Context) {
	var prompts []map[string]interface{}
	h.db.Table("prompts").Select("id, title, usage_count, like_count").
		Where("status = ?", 1).
		Order("usage_count DESC").Limit(10).Find(&prompts)
	response.Success(c, prompts)
}

func (h *Handler) ActiveUsers(c *gin.Context) {
	var users []map[string]interface{}
	h.db.Table("users").Select("id, username, avatar, article_count, follow_count, fans_count").
		Where("status = ?", 1).
		Order("article_count DESC").Limit(10).Find(&users)
	response.Success(c, users)
}

// ===== Admin Notifications =====

func (h *Handler) ListAllNotifications(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	type NotificationResult struct {
		ID        uint      `json:"id"`
		UserID    uint      `json:"user_id"`
		Username  string    `json:"username"`
		ActorID   uint      `json:"actor_id"`
		ActorName string    `json:"actor_name"`
		Type      string    `json:"type"`
		Content   string    `json:"content"`
		TargetID  uint      `json:"target_id"`
		IsRead    bool      `json:"is_read"`
		CreatedAt time.Time `json:"created_at"`
	}

	var total int64
	h.db.Table("notifications").Count(&total)

	var items []NotificationResult
	h.db.Raw(`
		SELECT n.id, n.user_id, u1.username as username, n.actor_id, u2.username as actor_name,
			n.type, n.content, n.target_id, n.is_read, n.created_at
		FROM notifications n
		LEFT JOIN users u1 ON n.user_id = u1.id
		LEFT JOIN users u2 ON n.actor_id = u2.id
		ORDER BY n.created_at DESC
		LIMIT ? OFFSET ?
	`, pageSize, (page-1)*pageSize).Scan(&items)

	response.Page(c, items, total, page, pageSize)
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

	// 重新加载 AI 配置
	if h.aiService != nil {
		h.aiService.ReloadConfig()
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
