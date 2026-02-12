package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"opendomain/internal/config"
	"opendomain/internal/middleware"
	"opendomain/internal/models"
)

type InvitationHandler struct {
	db  *gorm.DB
	cfg *config.Config
}

func NewInvitationHandler(db *gorm.DB, cfg *config.Config) *InvitationHandler {
	return &InvitationHandler{db: db, cfg: cfg}
}

// GetMyInvitations 获取我的邀请列表
func (h *InvitationHandler) GetMyInvitations(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var invitations []models.Invitation
	if err := h.db.Preload("Invitee").Where("inviter_id = ?", userID).Order("created_at DESC").Find(&invitations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch invitations"})
		return
	}

	responses := make([]*models.InvitationResponse, len(invitations))
	for i, inv := range invitations {
		responses[i] = inv.ToResponse()
	}

	c.JSON(http.StatusOK, gin.H{"invitations": responses})
}

// GetInvitationStats 获取邀请统计
func (h *InvitationHandler) GetInvitationStats(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// 计算总奖励
	var totalRewards int64
	h.db.Model(&models.Invitation{}).Where("inviter_id = ? AND reward_given = ?", userID, true).Count(&totalRewards)

	c.JSON(http.StatusOK, gin.H{
		"invite_code":        user.InviteCode,
		"total_invites":      user.TotalInvites,
		"successful_invites": user.SuccessfulInvites,
		"total_rewards":      totalRewards,
		"current_quota":      user.DomainQuota,
	})
}
