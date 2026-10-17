package handler

import (
	"reflect"
	"testing"
)

func TestNormalizeReportedHostname(t *testing.T) {
	cases := map[string]string{
		"  Foo.Loc.CC ":                  "foo.loc.cc",
		"foo.loc.cc.":                    "foo.loc.cc",
		"https://Login.Foo.loc.cc/x?y=1": "login.foo.loc.cc",
		"foo.loc.cc/path":                "foo.loc.cc",
		"http://foo.loc.cc:8080/":        "foo.loc.cc",
	}
	for in, want := range cases {
		got, ok := normalizeReportedHostname(in)
		if !ok || got != want {
			t.Errorf("normalizeReportedHostname(%q) = %q, %v; want %q", in, got, ok, want)
		}
	}
	for _, bad := range []string{"", "foo", "foo bar.com", "-foo.com", "foo..com", "http://", "<script>.com"} {
		if got, ok := normalizeReportedHostname(bad); ok {
			t.Errorf("normalizeReportedHostname(%q) = %q, want invalid", bad, got)
		}
	}
}

func TestHostnameCandidates(t *testing.T) {
	got := hostnameCandidates("www.foo.loc.cc")
	want := []string{"www.foo.loc.cc", "foo.loc.cc", "loc.cc"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("hostnameCandidates = %v, want %v", got, want)
	}
}

func TestIsValidReporterEmail(t *testing.T) {
	for _, ok := range []string{"a@b.co", "first.last+tag@example.com"} {
		if !isValidReporterEmail(ok) {
			t.Errorf("expected %q to be valid", ok)
		}
	}
	for _, bad := range []string{"a", "a@b", "Name <a@b.co>", "a@@b.co", "@b.co"} {
		if isValidReporterEmail(bad) {
			t.Errorf("expected %q to be invalid", bad)
		}
	}
}

func TestEscapeLike(t *testing.T) {
	if got := escapeLike(`50%_\`); got != `50\%\_\\` {
		t.Fatalf("escapeLike = %q", got)
	}
}
