package service

import (
	"strings"
	"testing"
)

func TestHashIP(t *testing.T) {
	t.Parallel()
	secret := "test-secret"

	h1 := HashIP("192.168.1.1", secret)
	h2 := HashIP("192.168.1.1", secret)
	h3 := HashIP("10.0.0.1", secret)
	h4 := HashIP("192.168.1.1", "other-secret")

	if h1 != h2 {
		t.Error("same IP + same secret must produce same hash")
	}
	if h1 == h3 {
		t.Error("different IPs must produce different hashes")
	}
	if h1 == h4 {
		t.Error("same IP + different secret must produce different hash")
	}
	if len(h1) != 64 {
		t.Errorf("expected 64-char hex hash, got %d: %q", len(h1), h1)
	}
	if strings.Contains(h1, "192.168") {
		t.Error("hash must not contain the original IP")
	}
}

func TestParseDeviceType(t *testing.T) {
	t.Parallel()
	cases := []struct {
		ua   string
		want string
	}{
		// bots
		{"Googlebot/2.1 (+http://www.google.com/bot.html)", "bot"},
		{"Mozilla/5.0 (compatible; bingbot/2.0)", "bot"},
		{"Mozilla/5.0 (compatible; YandexBot/3.0)", "bot"},
		{"python-requests/2.28.0", "desktop"}, // no bot keyword → falls through to desktop
		// tablets
		{"Mozilla/5.0 (iPad; CPU OS 15_0 like Mac OS X)", "tablet"},
		{"Mozilla/5.0 (Linux; Android 11; Tablet)", "tablet"},
		// mobile
		{"Mozilla/5.0 (iPhone; CPU iPhone OS 15_0 like Mac OS X)", "mobile"},
		{"Mozilla/5.0 (Linux; Android 12; Pixel 6) Mobile", "mobile"},
		{"BlackBerry9300/5.0.0.716", "mobile"},
		{"Windows Phone 8.1; ARM", "mobile"},
		// desktop
		{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36", "desktop"},
		{"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)", "desktop"},
		// unknown
		{"", "unknown"},
	}
	for _, tc := range cases {
		got := ParseDeviceType(tc.ua)
		if got != tc.want {
			t.Errorf("ParseDeviceType(%q) = %q, want %q", tc.ua, got, tc.want)
		}
	}
}

func TestParseDeviceType_BotBeforeMobile(t *testing.T) {
	t.Parallel()
	// Googlebot Smartphone includes "Mobile" — must be classified as bot, not mobile.
	ua := "Mozilla/5.0 (Linux; Android 6.0.1; Nexus 5X Build/MMB29P) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/W.X.Y.Z Mobile Safari/537.36 (compatible; Googlebot/2.1)"
	if got := ParseDeviceType(ua); got != "bot" {
		t.Errorf("expected bot, got %q", got)
	}
}

func TestParseReferrerSource(t *testing.T) {
	t.Parallel()
	cases := []struct {
		referer string
		want    string
	}{
		// direct
		{"", "direct"},
		// known socials
		{"https://www.instagram.com/", "instagram"},
		{"https://l.instagram.com/?u=https%3A%2F%2Fexample.com", "instagram"},
		{"https://www.facebook.com/sharer", "facebook"},
		{"https://fb.me/abc123", "facebook"},
		{"https://twitter.com/home", "twitter"},
		{"https://t.co/abc", "twitter"},
		{"https://x.com/user", "twitter"},
		{"https://www.tiktok.com/@user", "tiktok"},
		{"https://www.linkedin.com/feed/", "linkedin"},
		{"https://lnkd.in/abc", "linkedin"},
		{"https://wa.me/share", "whatsapp"},
		{"https://www.whatsapp.com", "whatsapp"},
		{"https://www.youtube.com/watch?v=abc", "youtube"},
		{"https://youtu.be/abc", "youtube"},
		{"https://www.google.com/search?q=test", "google"},
		{"https://google.co.uk/search?q=test", "google"},
		{"https://google.com.br/", "google"},
		// other
		{"https://news.ycombinator.com/item?id=1", "other"},
		{"https://reddit.com/r/golang", "other"},
		// malformed
		{"not-a-url", "other"},
	}
	for _, tc := range cases {
		got := ParseReferrerSource(tc.referer)
		if got != tc.want {
			t.Errorf("ParseReferrerSource(%q) = %q, want %q", tc.referer, got, tc.want)
		}
	}
}
