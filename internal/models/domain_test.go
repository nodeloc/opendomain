package models

import (
	"encoding/json"
	"testing"
	"time"
)

func TestDomainResponseClampsFarFutureExpiry(t *testing.T) {
	d := &Domain{
		RegisteredAt: time.Date(2026, 2, 13, 0, 0, 0, 0, time.UTC),
		ExpiresAt:    time.Date(10199, 2, 13, 0, 0, 0, 0, time.UTC),
	}
	b, err := json.Marshal(d.ToResponse())
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var got struct {
		ExpiresAt time.Time `json:"expires_at"`
	}
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.ExpiresAt.Year() != 9999 {
		t.Fatalf("expires_at year = %d, want 9999", got.ExpiresAt.Year())
	}
}
