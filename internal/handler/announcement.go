package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"opendomain/internal/config"
	"opendomain/internal/middleware"
	"opendomain/internal/models"
)

type AnnouncementHandler struct {
	db  *gorm.DB
	cfg *config.Config
}

func NewAnnouncementHandler(db *gorm.DB, cfg *config.Config) *AnnouncementHandler {
	return &AnnouncementHandler{db: db, cfg: cfg}
}

// ListPublicAnnouncements 获取公开的公告列表
func (h *AnnouncementHandler) ListPublicAnnouncements(c *gin.Context) {
	var announcements []models.Announcement
	if err := h.db.Preload("Author").
		Where("is_published = ?", true).
		Order("priority DESC, published_at DESC").
		Limit(20).
		Find(&announcements).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch announcements"})
		return
	}

	responses := make([]*models.AnnouncementResponse, len(announcements))
	for i, announcement := range announcements {
		responses[i] = announcement.ToResponse()
	}

	c.JSON(http.StatusOK, gin.H{"announcements": responses})
}

// GetAnnouncement 获取单个公告详情
func (h *AnnouncementHandler) GetAnnouncement(c *gin.Context) {
	announcementID := c.Param("id")

	var announcement models.Announcement
	if err := h.db.Preload("Author").First(&announcement, announcementID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Announcement not found"})
		return
	}

	// 增加浏览量
	h.db.Model(&announcement).UpdateColumn("views", gorm.Expr("views + ?", 1))

	c.JSON(http.StatusOK, announcement.ToResponse())
}

// ListAllAnnouncements 获取所有公告列表 (管理员)
func (h *AnnouncementHandler) ListAllAnnouncements(c *gin.Context) {
	var announcements []models.Announcement
	if err := h.db.Preload("Author").Order("created_at DESC").Find(&announcements).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch announcements"})
		return
	}

	responses := make([]*models.AnnouncementResponse, len(announcements))
	for i, announcement := range announcements {
		responses[i] = announcement.ToResponse()
	}

	c.JSON(http.StatusOK, gin.H{"announcements": responses})
}

// CreateAnnouncement 创建公告 (管理员)
func (h *AnnouncementHandler) CreateAnnouncement(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req models.AnnouncementCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	announcement := &models.Announcement{
		Title:       req.Title,
		Content:     req.Content,
		Type:        req.Type,
		Priority:    req.Priority,
		IsPublished: false,
		AuthorID:    &userID,
	}

	if err := h.db.Create(announcement).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create announcement"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Announcement created successfully",
		"announcement": announcement.ToResponse(),
	})
}

// UpdateAnnouncement 更新公告 (管理员)
func (h *AnnouncementHandler) UpdateAnnouncement(c *gin.Context) {
	announcementID := c.Param("id")

	var announcement models.Announcement
	if err := h.db.First(&announcement, announcementID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Announcement not found"})
		return
	}

	var req models.AnnouncementUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 更新字段
	if req.Title != nil {
		announcement.Title = *req.Title
	}
	if req.Content != nil {
		announcement.Content = *req.Content
	}
	if req.Type != nil {
		announcement.Type = *req.Type
	}
	if req.Priority != nil {
		announcement.Priority = *req.Priority
	}
	if req.IsPublished != nil {
		announcement.IsPublished = *req.IsPublished
		// 如果发布，设置发布时间
		if *req.IsPublished && announcement.PublishedAt == nil {
			now := time.Now()
			announcement.PublishedAt = &now
		}
	}

	if err := h.db.Save(&announcement).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update announcement"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Announcement updated successfully",
		"announcement": announcement.ToResponse(),
	})
}

// DeleteAnnouncement 删除公告 (管理员)
func (h *AnnouncementHandler) DeleteAnnouncement(c *gin.Context) {
	announcementID := c.Param("id")

	if err := h.db.Delete(&models.Announcement{}, announcementID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete announcement"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Announcement deleted successfully"})
}
