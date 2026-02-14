package router

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"

	"opendomain/internal/config"
	"opendomain/internal/handler"
	"opendomain/internal/middleware"
)

// Setup 设置路由
func Setup(db *gorm.DB, rdb *redis.Client, cfg *config.Config) *gin.Engine {
	r := gin.New()

	// 全局中间件
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.I18nMiddleware())

	// CORS 配置
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{cfg.FrontendURL},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * 3600,
	}))

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "OpenDomain API is running",
		})
	})

	// API 路由组
	api := r.Group("/api")
	{
		// 初始化处理器
		userHandler := handler.NewUserHandler(db, cfg)
		domainHandler := handler.NewDomainHandler(db, cfg)
		dnsHandler := handler.NewDNSHandler(db, cfg)
		couponHandler := handler.NewCouponHandler(db, cfg)
		invitationHandler := handler.NewInvitationHandler(db, cfg)
		announcementHandler := handler.NewAnnouncementHandler(db, cfg)
		domainScanHandler := handler.NewDomainScanHandler(db, cfg)
		orderHandler := handler.NewOrderHandler(db, cfg)
		paymentHandler := handler.NewPaymentHandler(db, cfg)
		pageHandler := handler.NewPageHandler(db, cfg)
		settingHandler := handler.NewSettingHandlerWithRedis(db, rdb, cfg)
		fossBillingSyncHandler := handler.NewFOSSBillingSyncHandler(db, cfg)

		// 公开路由
		public := api.Group("/public")
		{
			// 站点配置
			public.GET("/site-config", settingHandler.GetPublicSiteConfig)
			// 根域名列表
			public.GET("/root-domains", domainHandler.ListRootDomains)
			// 公告列表
			public.GET("/announcements", announcementHandler.ListPublicAnnouncements)
			public.GET("/announcements/:id", announcementHandler.GetAnnouncement)
			// 域名健康检查
			public.GET("/domain-health", domainScanHandler.ListDomainHealth)
			public.GET("/domain-health/:domainId", domainScanHandler.GetDomainHealth)
			public.GET("/domain-health/:domainId/scans", domainScanHandler.GetDomainScans)
			public.GET("/health-statistics", domainScanHandler.GetHealthStatistics)
			// 页面
			public.GET("/pages", pageHandler.GetPublicPages)
			public.GET("/pages/:slug", pageHandler.GetPublicPageBySlug)
			// 待激活域名列表(公开)
			public.GET("/pending-domains", fossBillingSyncHandler.GetPublicPendingDomains)
		}

		// 支付回调路由（公开，通过签名验证）
		api.GET("/payments/callback", paymentHandler.HandleCallback)
		api.GET("/payments/return", paymentHandler.HandleReturn)

		// 认证路由
		auth := api.Group("/auth")
		{
			auth.POST("/register", userHandler.Register)
			auth.POST("/login", userHandler.Login)

			// OAuth 路由
			oauthHandler := handler.NewOAuthHandler(db, cfg)
			auth.GET("/github", oauthHandler.GithubLogin)
			auth.GET("/github/callback", oauthHandler.GithubCallback)
			auth.GET("/google", oauthHandler.GoogleLogin)
			auth.GET("/google/callback", oauthHandler.GoogleCallback)
			auth.GET("/nodeloc", oauthHandler.NodelocLogin)
			auth.GET("/nodeloc/callback", oauthHandler.NodelocCallback)
		}

		// 需要认证的路由
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(cfg))
		{
			// 用户相关
			user := protected.Group("/user")
			{
				user.GET("/profile", userHandler.GetProfile)
				user.PUT("/profile", userHandler.UpdateProfile)
				user.PUT("/change-password", userHandler.ChangePassword)
				// FOSSBilling 同步
				user.POST("/sync-from-fossbilling", fossBillingSyncHandler.SyncFromFOSSBilling)
				user.GET("/sync-status", fossBillingSyncHandler.GetSyncStatus)
			}

			// 域名相关
			domains := protected.Group("/domains")
			{
				domains.GET("/search", domainHandler.SearchDomain)
				domains.POST("", domainHandler.RegisterDomain)
				domains.GET("", domainHandler.ListMyDomains)
				domains.GET("/:id", domainHandler.GetDomain)
				domains.DELETE("/:id", domainHandler.DeleteDomain)
				domains.PUT("/:id/nameservers", domainHandler.ModifyNameservers)
				domains.POST("/:id/renew", domainHandler.RenewDomain)
				domains.POST("/:id/transfer", domainHandler.TransferDomain)
			}

			// 域名扫描记录
			protected.GET("/domain-scans/:id", domainScanHandler.GetDomainScanRecords)

			// DNS 记录管理
			dns := protected.Group("/dns/:domainId/records")
			{
				dns.GET("", dnsHandler.ListRecords)
				dns.POST("", dnsHandler.CreateRecord)
				dns.POST("/sync-from-powerdns", dnsHandler.SyncFromPowerDNS)
				dns.GET("/:recordId", dnsHandler.GetRecord)
				dns.PUT("/:recordId", dnsHandler.UpdateRecord)
				dns.DELETE("/:recordId", dnsHandler.DeleteRecord)
			}

			// 优惠券
			coupons := protected.Group("/coupons")
			{
				coupons.POST("/apply", couponHandler.ApplyCoupon)
				coupons.GET("/my-usage", couponHandler.GetMyCouponUsage)
			}

			// 邀请
			invitations := protected.Group("/invitations")
			{
				invitations.GET("/my", invitationHandler.GetMyInvitations)
				invitations.GET("/stats", invitationHandler.GetInvitationStats)
			}

			// 订单
			orders := protected.Group("/orders")
			{
				orders.POST("/calculate", orderHandler.CalculatePrice)
				orders.POST("", orderHandler.CreateOrder)
				orders.GET("", orderHandler.ListMyOrders)
				orders.GET("/:id", orderHandler.GetOrder)
				orders.POST("/:id/cancel", orderHandler.CancelOrder)
			}

			// 支付
			payments := protected.Group("/payments")
			{
				payments.POST("/:orderId/initiate", paymentHandler.InitiatePayment)
				payments.POST("/:orderId/complete-free", paymentHandler.CompleteFreeOrder)
				payments.GET("/:orderId/status", paymentHandler.QueryPaymentStatus)
			}
		}

		// 管理员路由
		admin := api.Group("/admin")
		admin.Use(middleware.AuthMiddleware(cfg))
		admin.Use(middleware.AdminMiddleware())
		{
			// 系统设置
			admin.GET("/settings", settingHandler.GetSettings)
			admin.PUT("/settings/:key", settingHandler.UpdateSetting)
			admin.GET("/system-info", settingHandler.GetSystemInfo)
			admin.GET("/dashboard-stats", settingHandler.GetDashboardStats)
			admin.POST("/clear-cache", settingHandler.ClearCache)

			// 扫描管理
			admin.GET("/api-quota", domainScanHandler.GetAPIQuotaStatus)
			admin.GET("/scan-summaries", domainScanHandler.GetDomainScanSummaries)
			admin.GET("/scan-records", domainScanHandler.ListDomainScans)
			admin.GET("/suspend-history", domainScanHandler.GetSuspendHistory)

			admin.GET("/users", userHandler.ListUsers)
			admin.PUT("/users/:id", userHandler.AdminUpdateUser)
			admin.PUT("/users/:id/status", userHandler.AdminUpdateUserStatus)
			admin.DELETE("/users/:id", userHandler.AdminDeleteUser)
			admin.GET("/domains", domainHandler.ListAllDomains)
			admin.PUT("/domains/:id/status", domainHandler.AdminUpdateDomainStatus)
			admin.DELETE("/domains/:id", domainHandler.AdminDeleteDomain)
			admin.POST("/sync-fossbilling-domains", fossBillingSyncHandler.AdminSyncAllDomains)
			admin.GET("/pending-domains", fossBillingSyncHandler.ListPendingDomains)
			admin.DELETE("/pending-domains/:id", fossBillingSyncHandler.DeletePendingDomain)
			admin.GET("/orders", orderHandler.ListAllOrders)

			// 根域名管理
			admin.GET("/root-domains", domainHandler.ListAllRootDomains)
			admin.POST("/root-domains", domainHandler.CreateRootDomain)
			admin.PUT("/root-domains/:id", domainHandler.UpdateRootDomain)
			admin.DELETE("/root-domains/:id", domainHandler.DeleteRootDomain)
			admin.GET("/root-domains/:id/domains", domainHandler.ListDomainsByRootDomain)

			// 优惠券管理
			admin.GET("/coupons", couponHandler.ListCoupons)
			admin.POST("/coupons", couponHandler.CreateCoupon)
			admin.GET("/coupons/:id", couponHandler.GetCoupon)
			admin.PUT("/coupons/:id", couponHandler.UpdateCoupon)
			admin.DELETE("/coupons/:id", couponHandler.DeleteCoupon)

			// 公告管理
			admin.GET("/announcements", announcementHandler.ListAllAnnouncements)
			admin.POST("/announcements", announcementHandler.CreateAnnouncement)
			admin.GET("/announcements/:id", announcementHandler.GetAnnouncement)
			admin.PUT("/announcements/:id", announcementHandler.UpdateAnnouncement)
			admin.DELETE("/announcements/:id", announcementHandler.DeleteAnnouncement)

			// 页面管理
			admin.GET("/pages", pageHandler.GetAllPages)
			admin.POST("/pages", pageHandler.CreatePage)
			admin.PUT("/pages/:id", pageHandler.UpdatePage)
			admin.DELETE("/pages/:id", pageHandler.DeletePage)
		}
	}

	return r
}
