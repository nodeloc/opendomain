package handler

import (
	"errors"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"opendomain/internal/config"
	"opendomain/internal/models"
	"opendomain/pkg/logger"
)

const (
	sponsorNameMaxLength        = 100
	sponsorURLMaxLength         = 500
	sponsorDescriptionMaxLength = 300
	sponsorDisplayOrderLimit    = 10000
)

var trc20AddressPattern = regexp.MustCompile(`^T[1-9A-HJ-NP-Za-km-z]{33}$`)

// sponsorTierOrder sorts tiers as partner, sponsor, backer.
const sponsorTierOrder = "CASE tier WHEN 'partner' THEN 0 WHEN 'sponsor' THEN 1 WHEN 'backer' THEN 2 ELSE 3 END"

type SponsorHandler struct {
	db  *gorm.DB
	cfg *config.Config
}

func NewSponsorHandler(db *gorm.DB, cfg *config.Config) *SponsorHandler {
	return &SponsorHandler{db: db, cfg: cfg}
}

type sponsorFieldError struct {
	field    string
	key      string
	fallback string
}

func (h *SponsorHandler) loadSettings() models.SponsorSettings {
	return models.SponsorSettings{
		USDTTRC20Address:  models.GetSettingValue(h.db, models.SettingSponsorUSDTTRC20Address, ""),
		GithubSponsorsURL: models.GetSettingValue(h.db, models.SettingSponsorGithubURL, ""),
	}
}

// ListPublicSponsors GET /api/public/sponsors
func (h *SponsorHandler) ListPublicSponsors(c *gin.Context) {
	var sponsors []models.Sponsor
	if err := h.db.Where("is_active = ?", true).
		Order(sponsorTierOrder).Order("display_order ASC").Order("id ASC").
		Find(&sponsors).Error; err != nil {
		logger.Error("list public sponsors: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": tr(c, "error.internal_server", "Internal server error")})
		return
	}
	out := make([]models.PublicSponsor, 0, len(sponsors))
	for i := range sponsors {
		out = append(out, sponsors[i].ToPublic())
	}
	c.Header("Cache-Control", "public, max-age=60")
	c.JSON(http.StatusOK, gin.H{"sponsors": out, "payment": h.loadSettings()})
}

// AdminListSponsors GET /api/admin/sponsors
func (h *SponsorHandler) AdminListSponsors(c *gin.Context) {
	sponsors := []models.Sponsor{}
	if err := h.db.Order("display_order ASC").Order("id ASC").Find(&sponsors).Error; err != nil {
		logger.Error("list sponsors: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": tr(c, "error.internal_server", "Internal server error")})
		return
	}
	c.JSON(http.StatusOK, gin.H{"sponsors": sponsors})
}

// AdminCreateSponsor POST /api/admin/sponsors
func (h *SponsorHandler) AdminCreateSponsor(c *gin.Context) {
	var sponsor models.Sponsor
	if !h.bindSponsor(c, &sponsor) {
		return
	}
	if err := h.db.Create(&sponsor).Error; err != nil {
		logger.Error("create sponsor: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": tr(c, "error.internal_server", "Internal server error")})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"sponsor": sponsor})
}

// AdminUpdateSponsor PUT /api/admin/sponsors/:id
func (h *SponsorHandler) AdminUpdateSponsor(c *gin.Context) {
	sponsor, ok := h.findSponsor(c)
	if !ok {
		return
	}
	if !h.bindSponsor(c, sponsor) {
		return
	}
	// Select("*") writes every column so false / 0 / "" are not skipped.
	if err := h.db.Model(sponsor).Select("*").Omit("id", "created_at").Updates(sponsor).Error; err != nil {
		logger.Error("update sponsor: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": tr(c, "error.internal_server", "Internal server error")})
		return
	}
	c.JSON(http.StatusOK, gin.H{"sponsor": sponsor})
}

// AdminDeleteSponsor DELETE /api/admin/sponsors/:id
func (h *SponsorHandler) AdminDeleteSponsor(c *gin.Context) {
	sponsor, ok := h.findSponsor(c)
	if !ok {
		return
	}
	if err := h.db.Delete(&models.Sponsor{}, sponsor.ID).Error; err != nil {
		logger.Error("delete sponsor: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": tr(c, "error.internal_server", "Internal server error")})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": tr(c, "sponsor.deleted", "Sponsor deleted")})
}

// AdminGetSponsorSettings GET /api/admin/sponsor-settings
func (h *SponsorHandler) AdminGetSponsorSettings(c *gin.Context) {
	c.JSON(http.StatusOK, h.loadSettings())
}

// AdminUpdateSponsorSettings PUT /api/admin/sponsor-settings
func (h *SponsorHandler) AdminUpdateSponsorSettings(c *gin.Context) {
	var req models.SponsorSettings
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": tr(c, "error.validation", "Validation failed")})
		return
	}
	req.USDTTRC20Address = strings.TrimSpace(req.USDTTRC20Address)
	req.GithubSponsorsURL = strings.TrimSpace(req.GithubSponsorsURL)

	if req.USDTTRC20Address != "" && !trc20AddressPattern.MatchString(req.USDTTRC20Address) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": tr(c, "sponsor.usdt_invalid", "Enter a valid USDT (TRC20) address starting with T"),
			"field": "usdt_trc20_address",
		})
		return
	}
	if req.GithubSponsorsURL != "" && !isValidGithubSponsorsURL(req.GithubSponsorsURL) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": tr(c, "sponsor.github_url_invalid", "Enter an https://github.com/... URL"),
			"field": "github_sponsors_url",
		})
		return
	}

	rows := []models.SystemSetting{
		{SettingKey: models.SettingSponsorUSDTTRC20Address, SettingValue: req.USDTTRC20Address, Description: "USDT (TRC20) address shown on the sponsor page; empty hides it"},
		{SettingKey: models.SettingSponsorGithubURL, SettingValue: req.GithubSponsorsURL, Description: "GitHub Sponsors URL shown on the sponsor page; empty hides it"},
	}
	err := h.db.Transaction(func(tx *gorm.DB) error {
		for i := range rows {
			if err := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "setting_key"}},
				DoUpdates: clause.AssignmentColumns([]string{"setting_value", "updated_at"}),
			}).Create(&rows[i]).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		logger.Error("save sponsor settings: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": tr(c, "error.internal_server", "Internal server error")})
		return
	}
	c.JSON(http.StatusOK, h.loadSettings())
}

func (h *SponsorHandler) findSponsor(c *gin.Context) (*models.Sponsor, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": tr(c, "error.invalid_parameter", "Invalid parameter")})
		return nil, false
	}
	var sponsor models.Sponsor
	if err := h.db.First(&sponsor, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": tr(c, "sponsor.not_found", "Sponsor not found")})
		} else {
			logger.Error("find sponsor: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": tr(c, "error.internal_server", "Internal server error")})
		}
		return nil, false
	}
	return &sponsor, true
}

// bindSponsor parses and validates the request body and copies it onto
// sponsor. It writes the error response itself and returns false on failure.
func (h *SponsorHandler) bindSponsor(c *gin.Context, sponsor *models.Sponsor) bool {
	var req models.SponsorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": tr(c, "error.validation", "Validation failed")})
		return false
	}
	if fe := normalizeSponsorRequest(&req); fe != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": tr(c, fe.key, fe.fallback), "field": fe.field})
		return false
	}
	sponsor.Name = req.Name
	sponsor.WebsiteURL = req.WebsiteURL
	sponsor.LogoURL = req.LogoURL
	sponsor.LogoDarkURL = req.LogoDarkURL
	sponsor.DescriptionZh = req.DescriptionZh
	sponsor.DescriptionEn = req.DescriptionEn
	sponsor.Tier = req.Tier
	sponsor.DisplayOrder = req.DisplayOrder
	sponsor.IsActive = req.IsActive == nil || *req.IsActive
	return true
}

func normalizeSponsorRequest(req *models.SponsorRequest) *sponsorFieldError {
	req.Name = strings.TrimSpace(req.Name)
	if n := utf8.RuneCountInString(req.Name); n < 1 || n > sponsorNameMaxLength || hasControlChars(req.Name) {
		return &sponsorFieldError{"name", "sponsor.name_invalid", "Name must be 1 to 100 characters"}
	}

	urls := []struct {
		field string
		value *string
	}{
		{"website_url", &req.WebsiteURL},
		{"logo_url", &req.LogoURL},
		{"logo_dark_url", &req.LogoDarkURL},
	}
	for _, u := range urls {
		*u.value = strings.TrimSpace(*u.value)
		if *u.value != "" && !isValidSponsorURL(*u.value) {
			return &sponsorFieldError{u.field, "sponsor.url_invalid", "Enter a full http:// or https:// URL (up to 500 characters)"}
		}
	}

	req.DescriptionZh = strings.TrimSpace(req.DescriptionZh)
	if utf8.RuneCountInString(req.DescriptionZh) > sponsorDescriptionMaxLength {
		return &sponsorFieldError{"description_zh", "sponsor.description_too_long", "Description must be at most 300 characters"}
	}
	req.DescriptionEn = strings.TrimSpace(req.DescriptionEn)
	if utf8.RuneCountInString(req.DescriptionEn) > sponsorDescriptionMaxLength {
		return &sponsorFieldError{"description_en", "sponsor.description_too_long", "Description must be at most 300 characters"}
	}

	req.Tier = strings.TrimSpace(req.Tier)
	if req.Tier == "" {
		req.Tier = models.SponsorTierSponsor
	}
	if !isValidSponsorTier(req.Tier) {
		return &sponsorFieldError{"tier", "sponsor.tier_invalid", "Tier must be partner, sponsor or backer"}
	}

	if req.DisplayOrder < -sponsorDisplayOrderLimit || req.DisplayOrder > sponsorDisplayOrderLimit {
		return &sponsorFieldError{"display_order", "sponsor.order_invalid", "Display order must be between -10000 and 10000"}
	}
	return nil
}

func isValidSponsorTier(tier string) bool {
	for _, t := range models.SponsorTiers {
		if t == tier {
			return true
		}
	}
	return false
}

func hasControlChars(s string) bool {
	for _, r := range s {
		if unicode.IsControl(r) {
			return true
		}
	}
	return false
}

func parseHTTPURL(raw string) (*url.URL, bool) {
	if len(raw) > sponsorURLMaxLength || hasControlChars(raw) || strings.ContainsAny(raw, " \\") {
		return nil, false
	}
	u, err := url.Parse(raw)
	if err != nil || u.Host == "" || u.Hostname() == "" || u.User != nil {
		return nil, false
	}
	scheme := strings.ToLower(u.Scheme)
	if scheme != "http" && scheme != "https" {
		return nil, false
	}
	return u, true
}

func isValidSponsorURL(raw string) bool {
	_, ok := parseHTTPURL(raw)
	return ok
}

func isValidGithubSponsorsURL(raw string) bool {
	u, ok := parseHTTPURL(raw)
	if !ok || strings.ToLower(u.Scheme) != "https" {
		return false
	}
	return strings.ToLower(u.Hostname()) == "github.com" && u.Port() == ""
}
