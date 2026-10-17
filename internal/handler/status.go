package handler

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"

	"opendomain/internal/config"
	"opendomain/internal/models"
	"opendomain/pkg/powerdns"
	"opendomain/pkg/timeutil"
)

const (
	statusOperational = "operational"
	statusDegraded    = "degraded"
	statusOutage      = "outage"

	statusCacheTTL      = 15 * time.Second
	statusCheckTimeout  = 2 * time.Second
	statusSlowThreshold = time.Second
)

type StatusComponent struct {
	Key       string `json:"key"`
	Status    string `json:"status"`
	LatencyMS int64  `json:"latency_ms"`
}

type StatusResponse struct {
	Status      string            `json:"status"`
	CheckedAt   string            `json:"checked_at"`
	Components  []StatusComponent `json:"components"`
	Nameservers []string          `json:"nameservers"`
	RootDomains []string          `json:"root_domains"`
}

type StatusHandler struct {
	db   *gorm.DB
	rdb  *redis.Client
	cfg  *config.Config
	pdns *powerdns.Client

	mu       sync.Mutex
	cached   *StatusResponse
	cachedAt time.Time
}

func NewStatusHandler(db *gorm.DB, rdb *redis.Client, cfg *config.Config) *StatusHandler {
	return &StatusHandler{
		db:   db,
		rdb:  rdb,
		cfg:  cfg,
		pdns: powerdns.NewClient(cfg.PowerDNS.APIURL, cfg.PowerDNS.APIKey),
	}
}

var statusRank = map[string]int{statusOperational: 0, statusDegraded: 1, statusOutage: 2}

func timedCheck(key string, check func(ctx context.Context) error) StatusComponent {
	ctx, cancel := context.WithTimeout(context.Background(), statusCheckTimeout)
	defer cancel()
	start := time.Now()
	err := check(ctx)
	elapsed := time.Since(start)
	comp := StatusComponent{Key: key, Status: statusOperational, LatencyMS: elapsed.Milliseconds()}
	switch {
	case err != nil:
		comp.Status = statusOutage
	case elapsed > statusSlowThreshold:
		comp.Status = statusDegraded
	}
	return comp
}

func (h *StatusHandler) checkDatabase(ctx context.Context) error {
	sqlDB, err := h.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.PingContext(ctx)
}

func (h *StatusHandler) checkCache(ctx context.Context) error {
	if h.rdb == nil {
		return errors.New("redis client not configured")
	}
	return h.rdb.Ping(ctx).Err()
}

func (h *StatusHandler) checkDNSAPI(ctx context.Context) error {
	if strings.TrimSpace(h.cfg.PowerDNS.APIURL) == "" {
		return errors.New("PowerDNS API URL not configured")
	}
	return h.pdns.Ping(ctx)
}

func (h *StatusHandler) defaultNameservers() []string {
	out := make([]string, 0, 2)
	for _, ns := range []string{h.cfg.DNS.DefaultNS1, h.cfg.DNS.DefaultNS2} {
		ns = strings.TrimSuffix(strings.ToLower(strings.TrimSpace(ns)), ".")
		if ns != "" {
			out = append(out, ns)
		}
	}
	return out
}

func (h *StatusHandler) activeRootDomains() []string {
	out := []string{}
	ctx, cancel := context.WithTimeout(context.Background(), statusCheckTimeout)
	defer cancel()
	_ = h.db.WithContext(ctx).Model(&models.RootDomain{}).
		Where("is_active = ?", true).
		Order("priority DESC, id ASC").
		Pluck("domain", &out).Error
	if out == nil {
		out = []string{}
	}
	return out
}

func (h *StatusHandler) compute() *StatusResponse {
	checks := []struct {
		key string
		fn  func(ctx context.Context) error
	}{
		{"database", h.checkDatabase},
		{"cache", h.checkCache},
		{"dns_api", h.checkDNSAPI},
	}

	components := make([]StatusComponent, len(checks)+1)
	components[0] = StatusComponent{Key: "api", Status: statusOperational, LatencyMS: 0}
	var wg sync.WaitGroup
	for i, chk := range checks {
		wg.Add(1)
		go func(i int, key string, fn func(ctx context.Context) error) {
			defer wg.Done()
			components[i+1] = timedCheck(key, fn)
		}(i, chk.key, chk.fn)
	}
	wg.Wait()

	overall := statusOperational
	for _, comp := range components {
		if statusRank[comp.Status] > statusRank[overall] {
			overall = comp.Status
		}
	}

	return &StatusResponse{
		Status:      overall,
		CheckedAt:   timeutil.Now().Format(time.RFC3339),
		Components:  components,
		Nameservers: h.defaultNameservers(),
		RootDomains: h.activeRootDomains(),
	}
}

// GetPublicStatus 公开：系统状态（结果在内存中缓存，避免被用来压测数据库）
func (h *StatusHandler) GetPublicStatus(c *gin.Context) {
	// Holding the lock while refreshing makes concurrent requests wait for one
	// shared check instead of each hitting the backends.
	h.mu.Lock()
	if h.cached == nil || time.Since(h.cachedAt) > statusCacheTTL {
		h.cached = h.compute()
		h.cachedAt = time.Now()
	}
	resp := h.cached
	h.mu.Unlock()

	c.Header("Cache-Control", "public, max-age=15")
	c.JSON(http.StatusOK, resp)
}
