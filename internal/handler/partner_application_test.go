package handler

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"opendomain/internal/config"
	"opendomain/internal/middleware"
	"opendomain/internal/models"
)

func newPartnerRouter(t *testing.T) (*gin.Engine, *PartnerApplicationHandler, *fakeSender) {
	t.Helper()
	db := newTestDB(t)
	if err := db.AutoMigrate(&models.PartnerApplication{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	mailer, fake := syncMailer()
	cfg := &config.Config{SiteName: "Yunamo", FrontendURL: "https://yunamo.com/", SupportEmail: "support@yunamo.com"}
	h := NewPartnerApplicationHandler(db, nil, cfg, mailer)
	r := gin.New()
	r.Use(middleware.I18nMiddleware())
	r.POST("/api/public/partner-applications", h.CreatePartnerApplication)
	r.GET("/api/admin/partner-applications", h.ListPartnerApplications)
	r.PUT("/api/admin/partner-applications/:id", h.UpdatePartnerApplication)
	r.DELETE("/api/admin/partner-applications/:id", h.DeletePartnerApplication)
	return r, h, fake
}

func validPartnerBody() gin.H {
	return gin.H{
		"name":             "Ada <script>",
		"email":            "ada@example.com",
		"organization":     "Example Cloud",
		"website":          "https://example.com",
		"partnership_type": "resources",
		"message":          "We can offer two VPS in Tokyo for the project.\nLet's talk.",
		"contact":          "tg: ada",
	}
}

func TestPartnerApplicationStoresAndEmailsSupport(t *testing.T) {
	r, h, fake := newPartnerRouter(t)

	code, body := doJSON(r, "POST", "/api/public/partner-applications", validPartnerBody(), "Accept-Language", "en-US")
	if code != http.StatusCreated {
		t.Fatalf("create = %d %v", code, body)
	}
	var stored models.PartnerApplication
	if err := h.db.First(&stored).Error; err != nil {
		t.Fatalf("load: %v", err)
	}
	if stored.Status != models.PartnerApplicationStatusNew || stored.Locale != "en-US" || stored.PartnershipType != "resources" {
		t.Fatalf("unexpected stored application: %+v", stored)
	}

	sent := fake.all()
	if len(sent) != 1 {
		t.Fatalf("expected one support email, got %d", len(sent))
	}
	m := sent[0]
	if m.To != "support@yunamo.com" || m.ReplyTo != "ada@example.com" {
		t.Fatalf("To/Reply-To = %q / %q", m.To, m.ReplyTo)
	}
	if !strings.Contains(m.Subject, fmt.Sprintf("#%d", stored.ID)) || !strings.Contains(m.Subject, "Example Cloud") {
		t.Fatalf("subject = %q", m.Subject)
	}
	if strings.Contains(m.HTML, "<script>") || !strings.Contains(m.HTML, "https://yunamo.com/admin/partner-applications") {
		t.Fatal("support email HTML is not escaped or lacks the admin link")
	}
	if !strings.Contains(m.Text, "two VPS in Tokyo") || !strings.Contains(m.Text, "tg: ada") {
		t.Fatalf("support email text is missing the message or contact: %q", m.Text)
	}
}

func TestPartnerApplicationValidation(t *testing.T) {
	r, h, fake := newPartnerRouter(t)
	cases := map[string]func(gin.H){
		"missing name":    func(b gin.H) { b["name"] = "  " },
		"bad email":       func(b gin.H) { b["email"] = "ada@example.com, x@evil.com" },
		"missing email":   func(b gin.H) { delete(b, "email") },
		"bad type":        func(b gin.H) { b["partnership_type"] = "investor" },
		"short message":   func(b gin.H) { b["message"] = "hi" },
		"long message":    func(b gin.H) { b["message"] = strings.Repeat("a", partnerMessageMaxLength+1) },
		"js website":      func(b gin.H) { b["website"] = "javascript:alert(1)" },
		"website no host": func(b gin.H) { b["website"] = "https://" },
		"long contact":    func(b gin.H) { b["contact"] = strings.Repeat("c", partnerContactMaxLength+1) },
	}
	for name, mutate := range cases {
		b := validPartnerBody()
		mutate(b)
		if code, body := doJSON(r, "POST", "/api/public/partner-applications", b); code != http.StatusBadRequest || body["error"] == nil {
			t.Errorf("%s: got %d %v, want 400 with error", name, code, body)
		}
	}
	var n int64
	h.db.Model(&models.PartnerApplication{}).Count(&n)
	if n != 0 || len(fake.all()) != 0 {
		t.Fatalf("invalid submissions stored %d rows and sent %d emails", n, len(fake.all()))
	}
}

func TestPartnerApplicationHoneypotAndRateLimit(t *testing.T) {
	r, h, fake := newPartnerRouter(t)

	b := validPartnerBody()
	b["fax"] = "123"
	if code, _ := doJSON(r, "POST", "/api/public/partner-applications", b); code != http.StatusCreated {
		t.Fatalf("honeypot submission = %d, want a fake 201", code)
	}
	var n int64
	h.db.Model(&models.PartnerApplication{}).Count(&n)
	if n != 0 || len(fake.all()) != 0 {
		t.Fatal("honeypot submission was stored or emailed")
	}

	for i := 0; i < partnerApplicationRateLimit; i++ {
		if code, body := doJSON(r, "POST", "/api/public/partner-applications", validPartnerBody()); code != http.StatusCreated {
			t.Fatalf("submission %d = %d %v", i+1, code, body)
		}
	}
	if code, _ := doJSON(r, "POST", "/api/public/partner-applications", validPartnerBody()); code != http.StatusTooManyRequests {
		t.Fatalf("submission over the limit = %d, want 429", code)
	}
}

func TestPartnerApplicationAdminFlow(t *testing.T) {
	r, h, _ := newPartnerRouter(t)
	seed := []models.PartnerApplication{
		{Name: "Ada", Email: "ada@example.com", Organization: "Example Cloud", PartnershipType: "resources", Message: "servers for you", Status: "new"},
		{Name: "Bob", Email: "bob@example.com", Organization: "100%_Media", PartnershipType: "community", Message: "co-marketing idea", Status: "contacted"},
	}
	for i := range seed {
		h.db.Create(&seed[i])
	}

	code, body := doJSON(r, "GET", "/api/admin/partner-applications?status=new", nil)
	if code != http.StatusOK || body["total"].(float64) != 1 {
		t.Fatalf("list new = %d %v", code, body)
	}
	counts := body["counts"].(map[string]any)
	if counts["new"].(float64) != 1 || counts["contacted"].(float64) != 1 || counts["declined"].(float64) != 0 {
		t.Fatalf("counts = %v", counts)
	}
	if code, body = doJSON(r, "GET", "/api/admin/partner-applications?search=EXAMPLE%20CLOUD", nil); code != http.StatusOK || body["total"].(float64) != 1 {
		t.Fatalf("case-insensitive search = %d %v", code, body)
	}
	if code, body = doJSON(r, "GET", "/api/admin/partner-applications?search=100%25_", nil); code != http.StatusOK || body["total"].(float64) != 1 {
		t.Fatalf("search with LIKE wildcards = %d %v", code, body)
	}
	if code, _ = doJSON(r, "GET", "/api/admin/partner-applications?status=bogus", nil); code != http.StatusBadRequest {
		t.Fatalf("bad status filter = %d, want 400", code)
	}

	path := fmt.Sprintf("/api/admin/partner-applications/%d", seed[0].ID)
	code, body = doJSON(r, "PUT", path, gin.H{"status": "accepted", "admin_note": "  signed  "})
	if code != http.StatusOK {
		t.Fatalf("update = %d %v", code, body)
	}
	app := body["application"].(map[string]any)
	if app["status"] != "accepted" || app["admin_note"] != "signed" {
		t.Fatalf("updated application = %v", app)
	}
	if code, _ = doJSON(r, "PUT", path, gin.H{"status": "maybe"}); code != http.StatusBadRequest {
		t.Fatalf("invalid status update = %d, want 400", code)
	}

	if code, _ = doJSON(r, "DELETE", path, nil); code != http.StatusOK {
		t.Fatalf("delete = %d", code)
	}
	if code, _ = doJSON(r, "DELETE", path, nil); code != http.StatusNotFound {
		t.Fatalf("second delete = %d, want 404", code)
	}
}
