package i18n

import "testing"

func TestEmbeddedLocalesUsedWhenPathMissing(t *testing.T) {
	tr := &Translator{locales: make(map[string]map[string]string)}
	if err := tr.LoadLocales("/nonexistent"); err != nil {
		t.Fatalf("LoadLocales: %v", err)
	}
	if got := tr.T("zh-CN", "error.validation"); got != "验证失败" {
		t.Fatalf("zh-CN error.validation = %q", got)
	}
	if got := tr.T("en-US", "auth.invalid_token"); got == "auth.invalid_token" {
		t.Fatalf("en-US auth.invalid_token not translated")
	}
}
