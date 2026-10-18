package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"opendomain/internal/models"
)

func TestListDomainsByRootDomainPaginates(t *testing.T) {
	db := newTestDB(t)
	if err := db.AutoMigrate(&models.RootDomain{}, &models.Domain{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	owner := createTestUser(t, db, "owner@example.com", true)
	root := &models.RootDomain{Domain: "loc.cc", IsActive: true}
	if err := db.Create(root).Error; err != nil {
		t.Fatalf("create root: %v", err)
	}
	base := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	for i := 0; i < 5; i++ {
		expires := base.AddDate(1, 0, i)
		d := &models.Domain{
			UserID:       owner.ID,
			RootDomainID: root.ID,
			Subdomain:    fmt.Sprintf("d%d", i),
			FullDomain:   fmt.Sprintf("d%d.loc.cc", i),
			Status:       "active",
			RegisteredAt: base,
			ExpiresAt:    expires,
		}
		if err := db.Create(d).Error; err != nil {
			t.Fatalf("create domain: %v", err)
		}
	}

	h := &DomainHandler{db: db, cfg: testConfig()}
	r := gin.New()
	r.GET("/root-domains/:id/domains", h.ListDomainsByRootDomain)

	get := func(query string) map[string]any {
		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, fmt.Sprintf("/root-domains/%d/domains%s", root.ID, query), nil))
		if w.Code != http.StatusOK {
			t.Fatalf("status %d: %s", w.Code, w.Body.String())
		}
		var body map[string]any
		if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
			t.Fatalf("decode %q: %v", w.Body.String(), err)
		}
		return body
	}

	first := get("?page=1&page_size=2")
	if n := len(first["domains"].([]any)); n != 2 {
		t.Fatalf("page 1 has %d domains, want 2", n)
	}
	p := first["pagination"].(map[string]any)
	if p["total"].(float64) != 5 || p["total_pages"].(float64) != 3 || p["page"].(float64) != 1 {
		t.Fatalf("pagination = %v", p)
	}

	last := get("?page=3&page_size=2")
	if n := len(last["domains"].([]any)); n != 1 {
		t.Fatalf("page 3 has %d domains, want 1", n)
	}

	if n := len(get("?page_size=500")["domains"].([]any)); n != 5 {
		t.Fatalf("oversized page_size returned %d domains, want default page with all 5", n)
	}
}
