package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"opendomain/internal/config"
	"opendomain/internal/models"
)

type PageHandler struct {
	DB  *gorm.DB
	Cfg *config.Config
}

func NewPageHandler(db *gorm.DB, cfg *config.Config) *PageHandler {
	return &PageHandler{DB: db, Cfg: cfg}
}

// GetPublicPages - Get all published pages
func (h *PageHandler) GetPublicPages(c *gin.Context) {
	category := c.Query("category")

	var pages []models.Page
	query := h.DB.Where("is_published = ?", true)

	if category != "" {
		query = query.Where("category = ?", category)
	}

	if err := query.Order("category, display_order, id").Find(&pages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"pages": pages})
}

// GetPublicPageBySlug - Get a single page by slug
func (h *PageHandler) GetPublicPageBySlug(c *gin.Context) {
	slug := c.Param("slug")

	var page models.Page
	if err := h.DB.Where("slug = ? AND is_published = ?", slug, true).First(&page).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Page not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, page)
}

// Admin handlers

// GetAllPages - Admin: Get all pages
func (h *PageHandler) GetAllPages(c *gin.Context) {
	var pages []models.Page
	if err := h.DB.Order("category, display_order, id").Find(&pages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"pages": pages})
}

// CreatePage - Admin: Create a new page
func (h *PageHandler) CreatePage(c *gin.Context) {
	var page models.Page
	if err := c.ShouldBindJSON(&page); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if page.Title == "" || page.Slug == "" || page.Content == "" || page.Category == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title, slug, content, and category are required"})
		return
	}

	if err := h.DB.Create(&page).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, page)
}

// UpdatePage - Admin: Update a page
func (h *PageHandler) UpdatePage(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page ID"})
		return
	}

	var page models.Page
	if err := h.DB.First(&page, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Page not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var updateData models.Page
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Update fields
	page.Title = updateData.Title
	page.Slug = updateData.Slug
	page.Content = updateData.Content
	page.Category = updateData.Category
	page.IsPublished = updateData.IsPublished
	page.DisplayOrder = updateData.DisplayOrder

	if err := h.DB.Save(&page).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Page updated successfully", "page": page})
}

// DeletePage - Admin: Delete a page
func (h *PageHandler) DeletePage(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page ID"})
		return
	}

	if err := h.DB.Delete(&models.Page{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Page deleted successfully"})
}
