package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"opendomain/internal/models"
)

func newSponsorTestRouter(t *testing.T) (*gin.Engine, *gorm.DB) {
	t.Helper()
	dsn := fmt.Sprintf("file:sponsortest%d?mode=memory&cache=shared", atomic.AddInt64(&testDBSeq, 1))
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{Logger: gormlogger.Default.LogMode(gormlogger.Silent)})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(1)
	t.Cleanup(func() { sqlDB.Close() })
	if err := db.AutoMigrate(&models.Sponsor{}, &models.SystemSetting{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	h := NewSponsorHandler(db, testConfig())
	r := gin.New()
	r.GET("/api/public/sponsors", h.ListPublicSponsors)
	r.GET("/api/admin/sponsors", h.AdminListSponsors)
	r.POST("/api/admin/sponsors", h.AdminCreateSponsor)
	r.PUT("/api/admin/sponsors/:id", h.AdminUpdateSponsor)
	r.DELETE("/api/admin/sponsors/:id", h.AdminDeleteSponsor)
	r.GET("/api/admin/sponsor-settings", h.AdminGetSponsorSettings)
	r.PUT("/api/admin/sponsor-settings", h.AdminUpdateSponsorSettings)
	return r, db
}

func createSponsor(t *testing.T, r http.Handler, body map[string]interface{}) map[string]interface{} {
	t.Helper()
	code, out := doJSON(r, http.MethodPost, "/api/admin/sponsors", body)
	if code != http.StatusCreated {
		t.Fatalf("create sponsor %v: status %d, body %v", body, code, out)
	}
	return out["sponsor"].(map[string]interface{})
}

func sponsorNames(list interface{}) []string {
	var names []string
	for _, item := range list.([]interface{}) {
		names = append(names, item.(map[string]interface{})["name"].(string))
	}
	return names
}

func TestSponsorOrderingAndPublicFiltering(t *testing.T) {
	r, _ := newSponsorTestRouter(t)

	createSponsor(t, r, map[string]interface{}{"name": "Backer A", "tier": "backer", "display_order": -5})
	createSponsor(t, r, map[string]interface{}{"name": "Sponsor B", "display_order": 2})
	createSponsor(t, r, map[string]interface{}{"name": "Partner C", "tier": "partner", "display_order": 9})
	createSponsor(t, r, map[string]interface{}{"name": "Sponsor D", "display_order": 1})
	createSponsor(t, r, map[string]interface{}{"name": "Hidden", "tier": "partner", "is_active": false})
	createSponsor(t, r, map[string]interface{}{"name": "Sponsor E", "display_order": 1})

	code, out := doJSON(r, http.MethodGet, "/api/admin/sponsors", nil)
	if code != http.StatusOK {
		t.Fatalf("admin list: %d %v", code, out)
	}
	if got, want := fmt.Sprint(sponsorNames(out["sponsors"])), "[Backer A Hidden Sponsor D Sponsor E Sponsor B Partner C]"; got != want {
		t.Fatalf("admin order = %s, want %s", got, want)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/public/sponsors", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("public list: %d", w.Code)
	}
	if cc := w.Header().Get("Cache-Control"); cc != "public, max-age=60" {
		t.Fatalf("Cache-Control = %q", cc)
	}

	code, out = doJSON(r, http.MethodGet, "/api/public/sponsors", nil)
	if code != http.StatusOK {
		t.Fatalf("public list: %d %v", code, out)
	}
	if got, want := fmt.Sprint(sponsorNames(out["sponsors"])), "[Partner C Sponsor D Sponsor E Sponsor B Backer A]"; got != want {
		t.Fatalf("public order = %s, want %s", got, want)
	}
	first := out["sponsors"].([]interface{})[0].(map[string]interface{})
	for _, hidden := range []string{"is_active", "display_order", "created_at"} {
		if _, ok := first[hidden]; ok {
			t.Errorf("public sponsor exposes %q", hidden)
		}
	}
	payment, ok := out["payment"].(map[string]interface{})
	if !ok || payment["usdt_trc20_address"] != "" || payment["github_sponsors_url"] != "" {
		t.Fatalf("payment block = %v", out["payment"])
	}
}

func TestSponsorZeroValuesPersist(t *testing.T) {
	r, db := newSponsorTestRouter(t)

	s := createSponsor(t, r, map[string]interface{}{
		"name": "Zero", "display_order": 0, "is_active": false, "website_url": "https://example.com",
	})
	id := uint(s["id"].(float64))
	if s["is_active"] != false || s["tier"] != "sponsor" {
		t.Fatalf("create response = %v", s)
	}
	var stored models.Sponsor
	db.First(&stored, id)
	if stored.IsActive || stored.DisplayOrder != 0 || stored.Tier != "sponsor" {
		t.Fatalf("stored after create = %+v", stored)
	}

	path := fmt.Sprintf("/api/admin/sponsors/%d", id)
	code, out := doJSON(r, http.MethodPut, path, map[string]interface{}{
		"name": "Zero", "display_order": 7, "is_active": true, "website_url": "https://example.com",
	})
	if code != http.StatusOK {
		t.Fatalf("update: %d %v", code, out)
	}
	code, out = doJSON(r, http.MethodPut, path, map[string]interface{}{
		"name": "Zero", "display_order": 0, "is_active": false, "website_url": "",
	})
	if code != http.StatusOK {
		t.Fatalf("update: %d %v", code, out)
	}
	stored = models.Sponsor{}
	db.First(&stored, id)
	if stored.IsActive || stored.DisplayOrder != 0 || stored.WebsiteURL != "" {
		t.Fatalf("stored after update = %+v", stored)
	}

	created := createSponsor(t, r, map[string]interface{}{"name": "Default active"})
	if created["is_active"] != true {
		t.Fatalf("is_active should default to true, got %v", created["is_active"])
	}

	if code, _ := doJSON(r, http.MethodPut, "/api/admin/sponsors/9999", map[string]interface{}{"name": "x"}); code != http.StatusNotFound {
		t.Fatalf("update missing: %d", code)
	}
	if code, _ := doJSON(r, http.MethodPut, "/api/admin/sponsors/abc", map[string]interface{}{"name": "x"}); code != http.StatusBadRequest {
		t.Fatalf("update bad id: %d", code)
	}
	if code, out := doJSON(r, http.MethodDelete, path, nil); code != http.StatusOK || out["message"] == nil {
		t.Fatalf("delete: %d %v", code, out)
	}
	if code, _ := doJSON(r, http.MethodDelete, path, nil); code != http.StatusNotFound {
		t.Fatalf("delete again: %d", code)
	}
}

func TestSponsorValidation(t *testing.T) {
	r, _ := newSponsorTestRouter(t)

	cases := []struct {
		body  map[string]interface{}
		field string
	}{
		{map[string]interface{}{"name": "   "}, "name"},
		{map[string]interface{}{"name": "bad\x07name"}, "name"},
		{map[string]interface{}{"name": "ok", "website_url": "javascript:alert(1)"}, "website_url"},
		{map[string]interface{}{"name": "ok", "website_url": "/relative/path"}, "website_url"},
		{map[string]interface{}{"name": "ok", "logo_url": "data:image/png;base64,AAAA"}, "logo_url"},
		{map[string]interface{}{"name": "ok", "logo_dark_url": "//example.com/logo.png"}, "logo_dark_url"},
		{map[string]interface{}{"name": "ok", "website_url": "https://example.com/" + string(make([]byte, 500))}, "website_url"},
		{map[string]interface{}{"name": "ok", "tier": "gold"}, "tier"},
		{map[string]interface{}{"name": "ok", "display_order": 10001}, "display_order"},
		{map[string]interface{}{"name": "ok", "description_en": fmt.Sprintf("%0301d", 0)}, "description_en"},
	}
	for _, tc := range cases {
		code, out := doJSON(r, http.MethodPost, "/api/admin/sponsors", tc.body)
		if code != http.StatusBadRequest || out["field"] != tc.field {
			t.Errorf("body %v: status %d, out %v; want 400 field %s", tc.body, code, out, tc.field)
		}
	}

	createSponsor(t, r, map[string]interface{}{
		"name": "Valid", "website_url": "HTTPS://example.com/path?q=1", "logo_url": "http://cdn.example.com/a.png",
		"description_zh": "这是一段中文简介", "tier": "partner", "display_order": -10000,
	})

	req := httptest.NewRequest(http.MethodPost, "/api/admin/sponsors", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("empty body: %d", w.Code)
	}
}

func TestSponsorSettingsRoundTrip(t *testing.T) {
	r, _ := newSponsorTestRouter(t)

	code, out := doJSON(r, http.MethodGet, "/api/admin/sponsor-settings", nil)
	if code != http.StatusOK || out["usdt_trc20_address"] != "" || out["github_sponsors_url"] != "" {
		t.Fatalf("initial settings: %d %v", code, out)
	}

	addr := "TLa2f6VPqDgRE67v1736s7bJ8Ray5wYjU7"
	code, out = doJSON(r, http.MethodPut, "/api/admin/sponsor-settings", map[string]interface{}{
		"usdt_trc20_address": "  " + addr + " ", "github_sponsors_url": "https://github.com/sponsors/yunamo",
	})
	if code != http.StatusOK || out["usdt_trc20_address"] != addr || out["github_sponsors_url"] != "https://github.com/sponsors/yunamo" {
		t.Fatalf("save settings: %d %v", code, out)
	}

	code, out = doJSON(r, http.MethodGet, "/api/public/sponsors", nil)
	payment := out["payment"].(map[string]interface{})
	if code != http.StatusOK || payment["usdt_trc20_address"] != addr || payment["github_sponsors_url"] != "https://github.com/sponsors/yunamo" {
		t.Fatalf("public payment: %d %v", code, out)
	}

	code, out = doJSON(r, http.MethodPut, "/api/admin/sponsor-settings", map[string]interface{}{
		"usdt_trc20_address": "", "github_sponsors_url": "",
	})
	if code != http.StatusOK || out["usdt_trc20_address"] != "" || out["github_sponsors_url"] != "" {
		t.Fatalf("clear settings: %d %v", code, out)
	}
	code, out = doJSON(r, http.MethodGet, "/api/admin/sponsor-settings", nil)
	if code != http.StatusOK || out["usdt_trc20_address"] != "" || out["github_sponsors_url"] != "" {
		t.Fatalf("settings after clear: %d %v", code, out)
	}

	bad := []struct {
		body  map[string]interface{}
		field string
	}{
		{map[string]interface{}{"usdt_trc20_address": "0x52908400098527886E0F7030069857D2E4169EE7"}, "usdt_trc20_address"},
		{map[string]interface{}{"usdt_trc20_address": "TLa2f6VPqDgRE67v1736s7bJ8Ray5wYjU0"}, "usdt_trc20_address"},
		{map[string]interface{}{"usdt_trc20_address": "TLa2f6VPqDgRE67v1736s7bJ8Ray5wYjU"}, "usdt_trc20_address"},
		{map[string]interface{}{"github_sponsors_url": "http://github.com/sponsors/x"}, "github_sponsors_url"},
		{map[string]interface{}{"github_sponsors_url": "https://evil.com/sponsors/x"}, "github_sponsors_url"},
		{map[string]interface{}{"github_sponsors_url": "https://github.com.evil.com/x"}, "github_sponsors_url"},
		{map[string]interface{}{"github_sponsors_url": "javascript:alert(1)"}, "github_sponsors_url"},
	}
	for _, tc := range bad {
		code, out := doJSON(r, http.MethodPut, "/api/admin/sponsor-settings", tc.body)
		if code != http.StatusBadRequest || out["field"] != tc.field {
			t.Errorf("body %v: status %d, out %v; want 400 field %s", tc.body, code, out, tc.field)
		}
	}
}
