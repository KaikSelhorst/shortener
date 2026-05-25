package service

import (
	"log"
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var nonAlphanumeric = regexp.MustCompile(`[^a-z0-9]+`)

func GenerateSlug(name string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	normalized, _, err := transform.String(t, name)
	if err != nil {
		log.Printf("slug: unicode normalization failed for %q, using raw name: %v", name, err)
		normalized = name
	}
	slug := strings.ToLower(normalized)
	slug = nonAlphanumeric.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	return slug
}
