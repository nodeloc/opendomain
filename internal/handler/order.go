package handler

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"opendomain/internal/config"
	"opendomain/internal/middleware"
	"opendomain/internal/models"
)

// OrderHandler 订单处理器
type OrderHandler struct {
	db  *gorm.DB
	cfg *config.Config
}

// NewOrderHandler 创建订单处理器
func NewOrderHandler(db *gorm.DB, cfg *config.Config) *OrderHandler {
	return &OrderHandler{
		db:  db,
		cfg: cfg,
	}
}

// CalculatePrice 计算价格
func (h *OrderHandler) CalculatePrice(c *gin.Context) {
	var req models.OrderCalculateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取根域名
	var rootDomain models.RootDomain
	if err := h.db.First(&rootDomain, req.RootDomainID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Root domain not found"})
		return
	}

	// 检查是否为免费域名
	if rootDomain.IsFree {
		c.JSON(http.StatusOK, models.PriceCalculationResponse{
			BasePrice:      0,
			DiscountAmount: 0,
			FinalPrice:     0,
			CouponApplied:  false,
		})
		return
	}

	// 计算基础价格
	var basePrice float64
	if req.IsLifetime {
		if rootDomain.LifetimePrice == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Lifetime pricing not available for this domain"})
			return
		}
		basePrice = *rootDomain.LifetimePrice
	} else {
		if rootDomain.PricePerYear == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Pricing not configured for this domain"})
			return
		}
		basePrice = *rootDomain.PricePerYear * float64(req.Years)
	}

	// 应用优惠券（如果有）
	var discountAmount float64
	var coupon *models.Coupon
	var couponType *string
	var couponError *string

	if req.CouponCode != nil && *req.CouponCode != "" {
		if err := h.db.Where("UPPER(code) = UPPER(?)", *req.CouponCode).First(&coupon).Error; err == nil {
			fmt.Printf("[DEBUG] Coupon found: code=%s, type=%s, is_active=%v, discount_value=%v\n",
				coupon.Code, coupon.DiscountType, coupon.IsActive, coupon.DiscountValue)

			// 验证优惠券
			if err := h.validateCoupon(coupon, 0); err == nil {
				// 只应用 percentage 和 fixed 类型的优惠券
				if coupon.DiscountType == "percentage" && coupon.DiscountValue != nil {
					discountAmount = basePrice * (*coupon.DiscountValue / 100.0)
					fmt.Printf("[DEBUG] Applied percentage discount: %.2f%%\n", *coupon.DiscountValue)
				} else if coupon.DiscountType == "fixed" && coupon.DiscountValue != nil {
					discountAmount = math.Min(*coupon.DiscountValue, basePrice)
					fmt.Printf("[DEBUG] Applied fixed discount: %.2f\n", *coupon.DiscountValue)
				} else {
					errMsg := fmt.Sprintf("Coupon type '%s' cannot be applied in checkout. Only 'percentage' and 'fixed' discount coupons are valid for domain orders. This coupon may be for account quota increase.", coupon.DiscountType)
					couponError = &errMsg
					fmt.Printf("[DEBUG] %s\n", errMsg)
				}
				couponType = &coupon.DiscountType
			} else {
				errMsg := err.Error()
				couponError = &errMsg
				fmt.Printf("[DEBUG] Coupon validation failed: %v\n", err)
			}
		} else {
			errMsg := fmt.Sprintf("Coupon code '%s' not found", *req.CouponCode)
			couponError = &errMsg
			fmt.Printf("[DEBUG] Coupon not found: %s\n", *req.CouponCode)
		}
	}

	finalPrice := math.Max(0, basePrice-discountAmount)

	c.JSON(http.StatusOK, models.PriceCalculationResponse{
		BasePrice:      basePrice,
		DiscountAmount: discountAmount,
		FinalPrice:     finalPrice,
		CouponApplied:  discountAmount > 0,
		CouponCode:     req.CouponCode,
		CouponType:     couponType,
		CouponError:    couponError,
	})
}

// CreateOrder 创建订单
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req models.OrderCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取根域名
	var rootDomain models.RootDomain
	if err := h.db.First(&rootDomain, req.RootDomainID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Root domain not found"})
		return
	}

	// 检查是否为免费域名
	if rootDomain.IsFree {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "This domain is free. Please use the regular registration endpoint.",
		})
		return
	}

	// 构建完整域名
	fullDomain := fmt.Sprintf("%s.%s", req.Subdomain, rootDomain.Domain)

	// 检查域名是否已被注册
	var existingDomain models.Domain
	if err := h.db.Where("full_domain = ?", fullDomain).First(&existingDomain).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Domain already registered"})
		return
	}

	// 计算价格
	var basePrice float64
	if req.IsLifetime {
		if rootDomain.LifetimePrice == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Lifetime pricing not available"})
			return
		}
		basePrice = *rootDomain.LifetimePrice
	} else {
		if rootDomain.PricePerYear == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Pricing not configured"})
			return
		}
		basePrice = *rootDomain.PricePerYear * float64(req.Years)
	}

	// 应用优惠券
	var discountAmount float64
	var couponID *uint
	var couponCode *string

	if req.CouponCode != nil && *req.CouponCode != "" {
		var coupon models.Coupon
		if err := h.db.Where("UPPER(code) = UPPER(?)", *req.CouponCode).First(&coupon).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Coupon not found"})
			return
		}

		// 验证优惠券
		if err := h.validateCoupon(&coupon, userID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 应用折扣
		if coupon.DiscountType == "percentage" && coupon.DiscountValue != nil {
			discountAmount = basePrice * (*coupon.DiscountValue / 100.0)
		} else if coupon.DiscountType == "fixed" && coupon.DiscountValue != nil {
			discountAmount = math.Min(*coupon.DiscountValue, basePrice)
		} else if coupon.DiscountType == "quota_increase" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "This coupon type should be applied through the regular coupon endpoint",
			})
			return
		}

		couponID = &coupon.ID
		couponCode = &coupon.Code
	}

	finalPrice := math.Max(0, basePrice-discountAmount)

	// 生成订单号
	orderNumber := h.generateOrderNumber()

	// 创建订单（15分钟后过期）
	order := &models.Order{
		OrderNumber:    orderNumber,
		UserID:         userID,
		Subdomain:      req.Subdomain,
		RootDomainID:   req.RootDomainID,
		FullDomain:     fullDomain,
		Years:          req.Years,
		IsLifetime:     req.IsLifetime,
		BasePrice:      basePrice,
		DiscountAmount: discountAmount,
		FinalPrice:     finalPrice,
		CouponID:       couponID,
		CouponCode:     couponCode,
		Status:         "pending",
		ExpiresAt:      time.Now().Add(15 * time.Minute),
	}

	if err := h.db.Create(order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	// 预加载关联数据
	h.db.Preload("RootDomain").First(order, order.ID)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Order created successfully",
		"order":   order.ToResponse(),
	})
}

// GetOrder 获取订单详情
func (h *OrderHandler) GetOrder(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	orderID := c.Param("id")

	var order models.Order
	if err := h.db.Preload("RootDomain").Preload("Payment").
		Where("id = ? AND user_id = ?", orderID, userID).
		First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	c.JSON(http.StatusOK, order.ToResponse())
}

// ListMyOrders 列出我的订单
func (h *OrderHandler) ListMyOrders(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// 分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status := c.Query("status")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	// 构建查询
	query := h.db.Model(&models.Order{}).Where("user_id = ?", userID)
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 获取总数
	var total int64
	query.Count(&total)

	// 获取订单列表
	var orders []models.Order
	if err := query.Preload("RootDomain").Preload("Payment").
		Order("created_at DESC").
		Limit(pageSize).Offset(offset).
		Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}

	// 转换为响应格式
	var orderResponses []*models.OrderResponse
	for i := range orders {
		orderResponses = append(orderResponses, orders[i].ToResponse())
	}

	c.JSON(http.StatusOK, gin.H{
		"orders":    orderResponses,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// CancelOrder 取消订单
func (h *OrderHandler) CancelOrder(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	orderID := c.Param("id")

	var order models.Order
	if err := h.db.Where("id = ? AND user_id = ?", orderID, userID).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	// 只能取消待支付或已过期的订单
	if order.Status != "pending" && order.Status != "expired" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only pending or expired orders can be cancelled"})
		return
	}

	// 更新状态
	order.Status = "cancelled"
	if err := h.db.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Order cancelled successfully",
		"order":   order.ToResponse(),
	})
}

// CleanupExpiredOrders 清理过期订单（后台任务）
func (h *OrderHandler) CleanupExpiredOrders() error {
	result := h.db.Model(&models.Order{}).
		Where("status = ? AND expires_at < ?", "pending", time.Now()).
		Update("status", "expired")

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected > 0 {
		fmt.Printf("Cleaned up %d expired orders\n", result.RowsAffected)
	}

	return nil
}

// validateCoupon 验证优惠券
func (h *OrderHandler) validateCoupon(coupon *models.Coupon, userID uint) error {
	now := time.Now()

	// 检查是否激活
	if !coupon.IsActive {
		return fmt.Errorf("Coupon is not active. This coupon has been disabled by administrator")
	}

	// 检查有效期
	if now.Before(coupon.ValidFrom) {
		return fmt.Errorf("Coupon is not yet valid. This coupon will be available after %s", 
			coupon.ValidFrom.Format("2006-01-02 15:04:05"))
	}
	if coupon.ValidUntil != nil && now.After(*coupon.ValidUntil) {
		return fmt.Errorf("Coupon has expired on %s", 
			coupon.ValidUntil.Format("2006-01-02 15:04:05"))
	}

	// 检查使用次数
	if coupon.MaxUses > 0 && coupon.UsedCount >= coupon.MaxUses {
		return fmt.Errorf("Coupon usage limit reached (%d/%d uses)", 
			coupon.UsedCount, coupon.MaxUses)
	}

	// 对于不可重复使用的优惠券，检查用户是否已使用
	if userID > 0 && !coupon.IsReusable {
		var usage models.CouponUsage
		if err := h.db.Where("coupon_id = ? AND user_id = ?", coupon.ID, userID).
			First(&usage).Error; err == nil {
			return fmt.Errorf("You have already used this coupon on %s. This coupon can only be used once per user", 
				usage.UsedAt.Format("2006-01-02 15:04:05"))
		}
	}

	return nil
}

// generateOrderNumber 生成订单号
func (h *OrderHandler) generateOrderNumber() string {
	// 格式: ORD + 时间戳 + 随机字符
	timestamp := time.Now().Unix()
	randomBytes := make([]byte, 4)
	rand.Read(randomBytes)
	randomStr := hex.EncodeToString(randomBytes)
	return fmt.Sprintf("ORD%d%s", timestamp, randomStr[:8])
}

// ListAllOrders 管理员：获取所有订单
func (h *OrderHandler) ListAllOrders(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status := c.Query("status")
	search := c.Query("search")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	query := h.db.Model(&models.Order{}).
		Preload("User").
		Preload("RootDomain")

	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 搜索功能：支持按订单号、用户名、邮箱、域名搜索
	if search != "" {
		query = query.Joins("LEFT JOIN users ON users.id = orders.user_id").
			Where("orders.order_number LIKE ? OR orders.subdomain LIKE ? OR users.username LIKE ? OR users.email LIKE ?",
				"%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	var total int64
	query.Count(&total)

	var orders []models.Order
	if err := query.Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"orders": orders,
		"pagination": gin.H{
			"page":       page,
			"page_size":  pageSize,
			"total":      total,
			"total_page": int(math.Ceil(float64(total) / float64(pageSize))),
		},
	})
}
