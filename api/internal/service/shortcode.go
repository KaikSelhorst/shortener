package service

import (
	"fmt"
	"strings"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

const base = uint64(len(alphabet)) // 62

// ShortcodeService encodes and decodes short codes using bijective base-62
// over a fixed alphabet. The encoding is a stable contract with the database —
// changing it would invalidate all existing short_code values.
type ShortcodeService struct{}

func NewShortcodeService() (*ShortcodeService, error) {
	return &ShortcodeService{}, nil
}

func (svc *ShortcodeService) GenerateShortCode(id uint64) (string, error) {
	if id == 0 {
		return "", fmt.Errorf("shortcode: id must be greater than zero")
	}
	// MaxUint64 in base 62 needs at most 11 chars; 12 is a safe ceiling.
	var buf [12]byte
	n := 0
	for id > 0 {
		buf[n] = alphabet[id%base]
		id /= base
		n++
	}
	for i, j := 0, n-1; i < j; i, j = i+1, j-1 {
		buf[i], buf[j] = buf[j], buf[i]
	}
	return string(buf[:n]), nil
}

func (svc *ShortcodeService) DecodeShortCode(code string) uint64 {
	if code == "" {
		return 0
	}
	var id uint64
	for _, c := range code {
		idx := strings.IndexRune(alphabet, c)
		if idx < 0 {
			return 0
		}
		id = id*base + uint64(idx)
	}
	return id
}
