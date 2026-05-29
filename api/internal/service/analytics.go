package service

import (
	"net/url"
	"strings"
)

// ParseDeviceType classifies a User-Agent string into one of:
// "bot", "tablet", "mobile", "desktop", or "unknown".
// Order matters: bots are checked before mobile because many crawlers
// include "Mobile" in their UA string (e.g. Googlebot Smartphone).
func ParseDeviceType(ua string) string {
	if ua == "" {
		return "unknown"
	}
	lower := strings.ToLower(ua)

	for _, kw := range []string{"bot", "crawler", "spider", "slurp", "scan", "checker", "archiver", "fetcher"} {
		if strings.Contains(lower, kw) {
			return "bot"
		}
	}

	for _, kw := range []string{"tablet", "ipad"} {
		if strings.Contains(lower, kw) {
			return "tablet"
		}
	}

	for _, kw := range []string{"mobile", "android", "iphone", "ipod", "blackberry", "windows phone", "opera mini", "opera mobi"} {
		if strings.Contains(lower, kw) {
			return "mobile"
		}
	}

	return "desktop"
}

// knownReferrers maps hostname suffixes to their canonical source name.
var knownReferrers = map[string]string{
	"instagram.com": "instagram",
	"facebook.com":  "facebook",
	"fb.com":        "facebook",
	"fb.me":         "facebook",
	"twitter.com":   "twitter",
	"t.co":          "twitter",
	"x.com":         "twitter",
	"tiktok.com":    "tiktok",
	"linkedin.com":  "linkedin",
	"lnkd.in":       "linkedin",
	"whatsapp.com":  "whatsapp",
	"wa.me":         "whatsapp",
	"youtube.com":   "youtube",
	"youtu.be":      "youtube",
	"google.com":    "google",
	"google.co":     "google", // google.co.uk, google.com.br, etc.
	"bing.com":      "google", // group search engines under google label isn't ideal; keep separate
}

// ParseReferrerSource extracts the host from the referer URL and maps it to
// a known social/search source. Returns "direct" for empty referers.
func ParseReferrerSource(referer string) string {
	if referer == "" {
		return "direct"
	}

	u, err := url.Parse(referer)
	if err != nil || u.Host == "" {
		return "other"
	}

	host := strings.ToLower(u.Hostname())
	// Strip www. prefix for matching.
	host = strings.TrimPrefix(host, "www.")

	// Exact match first.
	if src, ok := knownReferrers[host]; ok {
		return src
	}

	// Suffix match for subdomains (e.g. l.instagram.com, lm.facebook.com).
	for suffix, src := range knownReferrers {
		if strings.HasSuffix(host, "."+suffix) {
			return src
		}
	}

	// google.co.* country domains (google.co.uk, google.com.br, etc.)
	if strings.HasPrefix(host, "google.") {
		return "google"
	}

	return "other"
}
