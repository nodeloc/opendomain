package handler

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"

	"opendomain/internal/config"
	"opendomain/internal/middleware"
	"opendomain/internal/models"
	"opendomain/internal/services"
	"opendomain/pkg/logger"
	"opendomain/pkg/timeutil"
)

const (
	partnerApplicationRateLimit     = 5
	partnerApplicationGlobalLimit   = 100
	partnerApplicationRateWindow    = time.Hour
	partnerApplicationMaxBodyBytes  = 64 << 10
	partnerNameMaxLength            = 100
	partnerOrganizationMaxLength    = 200
	partnerWebsiteMaxLength         = 500
	partnerContactMaxLength         = 200
	partnerMessageMinLength         = 10
	partnerMessageMaxLength         = 5000
	partnerApplicationNoteMaxLength = 5000
)

type PartnerApplicationHandler struct {
	db     *gorm.DB
	rdb    *redis.Client
	cfg    *config.Config
	mailer *services.Mailer
}

func NewPartnerApplicationHandler(db *gorm.DB, rdb *redis.Client, cfg *config.Config, mailer *services.Mailer) *PartnerApplicationHandler {
	return &PartnerApplicationHandler{db: db, rdb: rdb, cfg: cfg, mailer: mailer}
}

func containsString(list []string, v string) bool {
	for _, s := range list {
		if s == v {
			return true
		}
	}
	return false
}

// cleanText trims, strips NUL and collapses invalid UTF-8 in user input.
func cleanText(s string) string {
	return strings.TrimSpace(stripNUL(strings.ToValidUTF8(s, "")))
}

func isValidPartnerWebsite(raw string) bool {
	u, err := url.Parse(raw)
	if err != nil || u.Host == "" || u.User != nil {
		return false
	}
	return u.Scheme == "http" || u.Scheme == "https"
}

// allowApplication mirrors the abuse report limiter: a per-IP limit plus a
// site-wide cap, with Redis first and the database as the fallback counter.
func (h *PartnerApplicationHandler) allowApplication(ctx context.Context, ip string) bool {
	if h.rdb != nil {
		rctx, cancel := context.WithTimeout(ctx, 2*time.Second)
		defer cancel()
		global, gerr := incrRateWindow(rctx, h.rdb, "partner_application:rate:global", partnerApplicationRateWindow)
		perIP, err := incrRateWindow(rctx, h.rdb, "partner_application:rate:"+ip, partnerApplicationRateWindow)
		if err == nil && gerr == nil {
			return perIP <= partnerApplicationRateLimit && global <= partnerApplicationGlobalLimit
		}
		logger.Warnf("partner application rate limit: redis unavailable, falling back to database: %v %v", err, gerr)
	}

	since := timeutil.Now().Add(-partnerApplicationRateWindow)
	var count, total int64
	if err := h.db.Model(&models.PartnerApplication{}).
		Where("ip = ? AND created_at > ?", ip, since).
		Count(&count).Error; err != nil {
		logger.Errorf("partner application rate limit: database count failed: %v", err)
		return true
	}
	if err := h.db.Model(&models.PartnerApplication{}).
		Where("created_at > ?", since).
		Count(&total).Error; err != nil {
		logger.Errorf("partner application rate limit: database count failed: %v", err)
		return true
	}
	return count < partnerApplicationRateLimit && total < partnerApplicationGlobalLimit
}

func (h *PartnerApplicationHandler) received(c *gin.Context, id uint) {
	c.JSON(http.StatusCreated, gin.H{
		"id":      id,
		"message": tr(c, "partner.application_received", "Thanks, we've received your application and will reply by email."),
	})
}

// CreatePartnerApplication 公开：提交合作申请
func (h *PartnerApplicationHandler) CreatePartnerApplication(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, partnerApplicationMaxBodyBytes)

	var req models.PartnerApplicationCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": tr(c, "error.validation", "Validation failed")})
		return
	}

	// Answer bots that filled the honeypot as if it worked, so they don't retry.
	if strings.TrimSpace(req.Fax) != "" {
		h.received(c, 0)
		return
	}

	bad := func(key, fallback string) {
		c.JSON(http.StatusBadRequest, gin.H{"error": tr(c, key, fallback)})
	}

	name := cleanText(req.Name)
	if n := utf8.RuneCountInString(name); n == 0 || n > partnerNameMaxLength {
		bad("partner.invalid_name", "Enter your name (100 characters max)")
		return
	}

	email := strings.TrimSpace(req.Email)
	if email == "" || !isValidReporterEmail(email) {
		bad("validation.email_invalid", "Invalid email format")
		return
	}

	organization := cleanText(req.Organization)
	if utf8.RuneCountInString(organization) > partnerOrganizationMaxLength {
		bad("partner.invalid_organization", "Organization must be at most 200 characters")
		return
	}

	website := cleanText(req.Website)
	if website != "" && (len(website) > partnerWebsiteMaxLength || !isValidPartnerWebsite(website)) {
		bad("partner.invalid_website", "Enter a website starting with http:// or https://")
		return
	}

	partnershipType := strings.ToLower(strings.TrimSpace(req.PartnershipType))
	if !containsString(models.PartnershipTypes, partnershipType) {
		bad("partner.invalid_type", "Choose a partnership type")
		return
	}

	message := cleanText(req.Message)
	if n := utf8.RuneCountInString(message); n < partnerMessageMinLength || n > partnerMessageMaxLength {
		bad("partner.invalid_message", "Message must be between 10 and 5000 characters")
		return
	}

	contact := cleanText(req.Contact)
	if utf8.RuneCountInString(contact) > partnerContactMaxLength {
		bad("partner.invalid_contact", "Other contact must be at most 200 characters")
		return
	}

	ip := reporterIP(c)
	if !h.allowApplication(c.Request.Context(), ip) {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": tr(c, "partner.rate_limited", "Too many applications from your network. Please try again later.")})
		return
	}

	userAgent := c.Request.UserAgent()
	if len(userAgent) > 512 {
		userAgent = userAgent[:512]
	}

	app := &models.PartnerApplication{
		Name:            name,
		Email:           email,
		Organization:    organization,
		Website:         website,
		PartnershipType: partnershipType,
		Message:         message,
		Contact:         contact,
		Status:          models.PartnerApplicationStatusNew,
		IP:              ip,
		UserAgent:       cleanText(userAgent),
		Locale:          services.NormalizeLocale(middleware.GetLocale(c)),
	}
	if err := h.db.Create(app).Error; err != nil {
		logger.Errorf("failed to create partner application: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": tr(c, "error.internal_server", "Internal server error")})
		return
	}

	h.notifySupport(*app)
	h.received(c, app.ID)
}

func (h *PartnerApplicationHandler) notifySupport(app models.PartnerApplication) {
	to := strings.TrimSpace(h.cfg.SupportEmail)
	if to == "" {
		return
	}
	m, err := services.PartnerApplicationMail(h.cfg.SiteName, to, services.PartnerApplicationAlert{
		ID:              app.ID,
		Name:            app.Name,
		Email:           app.Email,
		Organization:    app.Organization,
		Website:         app.Website,
		PartnershipType: app.PartnershipType,
		Contact:         app.Contact,
		Message:         app.Message,
		Locale:          app.Locale,
		CreatedAt:       app.CreatedAt,
		AdminURL:        strings.TrimRight(h.cfg.FrontendURL, "/") + "/admin/partner-applications",
	})
	if err != nil {
		logger.Errorf("partner application %d: failed to build support email: %v", app.ID, err)
		return
	}
	h.mailer.SendAsync(m)
}

// ListPartnerApplications 管理员：合作申请列表
func (h *PartnerApplicationHandler) ListPartnerApplications(c *gin.Context) {
	status := strings.TrimSpace(c.Query("status"))
	search := strings.TrimSpace(c.Query("search"))
	page := 1
	pageSize := 20
	if p, err := strconv.Atoi(c.Query("page")); err == nil && p > 0 {
		page = p
	}
	if ps, err := strconv.Atoi(c.Query("page_size")); err == nil && ps > 0 && ps <= 100 {
		pageSize = ps
	}
	if status != "" && !containsString(models.PartnerApplicationStatuses, status) {
		c.JSON(http.StatusBadRequest, gin.H{"error": tr(c, "partner.invalid_status", "Invalid status")})
		return
	}
	fetchFailed := func() {
		c.JSON(http.StatusInternalServerError, gin.H{"error": tr(c, "partner.fetch_failed", "Failed to fetch partner applications")})
	}

	base := func() *gorm.DB {
		q := h.db.Model(&models.PartnerApplication{})
		if search != "" {
			like := "%" + strings.ToLower(escapeLike(search)) + "%"
			q = q.Where(`LOWER(name) LIKE ? ESCAPE '\' OR LOWER(email) LIKE ? ESCAPE '\' OR LOWER(organization) LIKE ? ESCAPE '\' OR LOWER(message) LIKE ? ESCAPE '\'`, like, like, like, like)
		}
		return q
	}

	counts := gin.H{}
	for _, s := range models.PartnerApplicationStatuses {
		counts[s] = int64(0)
	}
	var grouped []struct {
		Status string
		Count  int64
	}
	if err := base().Select("status, COUNT(*) AS count").Group("status").Scan(&grouped).Error; err != nil {
		fetchFailed()
		return
	}
	for _, g := range grouped {
		if _, ok := counts[g.Status]; ok {
			counts[g.Status] = g.Count
		}
	}

	query := base()
	if status != "" {
		query = query.Where("status = ?", status)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		fetchFailed()
		return
	}
	apps := []models.PartnerApplication{}
	if err := query.Order("created_at DESC, id DESC").
		Offset((page - 1) * pageSize).Limit(pageSize).
		Find(&apps).Error; err != nil {
		fetchFailed()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"applications": apps,
		"total":        total,
		"page":         page,
		"page_size":    pageSize,
		"counts":       counts,
	})
}

func (h *PartnerApplicationHandler) findApplication(c *gin.Context) (*models.PartnerApplication, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": tr(c, "error.invalid_parameter", "Invalid parameter")})
		return nil, false
	}
	var app models.PartnerApplication
	if err := h.db.First(&app, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": tr(c, "partner.not_found", "Partner application not found")})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": tr(c, "error.internal_server", "Internal server error")})
		}
		return nil, false
	}
	return &app, true
}

// UpdatePartnerApplication 管理员：更新合作申请状态/备注
func (h *PartnerApplicationHandler) UpdatePartnerApplication(c *gin.Context) {
	app, ok := h.findApplication(c)
	if !ok {
		return
	}
	var req models.PartnerApplicationUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": tr(c, "error.validation", "Validation failed")})
		return
	}
	if req.Status != nil {
		s := strings.TrimSpace(*req.Status)
		if !containsString(models.PartnerApplicationStatuses, s) {
			c.JSON(http.StatusBadRequest, gin.H{"error": tr(c, "partner.invalid_status", "Invalid status")})
			return
		}
		app.Status = s
	}
	if req.AdminNote != nil {
		note := cleanText(*req.AdminNote)
		if utf8.RuneCountInString(note) > partnerApplicationNoteMaxLength {
			c.JSON(http.StatusBadRequest, gin.H{"error": tr(c, "partner.note_too_long", "Admin note must be at most 5000 characters")})
			return
		}
		app.AdminNote = note
	}
	if err := h.db.Save(app).Error; err != nil {
		logger.Errorf("failed to update partner application %d: %v", app.ID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": tr(c, "partner.update_failed", "Failed to update partner application")})
		return
	}
	c.JSON(http.StatusOK, gin.H{"application": app})
}

// DeletePartnerApplication 管理员：删除合作申请
func (h *PartnerApplicationHandler) DeletePartnerApplication(c *gin.Context) {
	app, ok := h.findApplication(c)
	if !ok {
		return
	}
	if err := h.db.Delete(&models.PartnerApplication{}, app.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": tr(c, "partner.delete_failed", "Failed to delete partner application")})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": tr(c, "partner.deleted", "Partner application deleted")})
}
