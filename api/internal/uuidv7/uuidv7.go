// Package uuidv7 generates time-ordered UUID v7 values using only stdlib.
package uuidv7

import (
	"crypto/rand"
	"fmt"
	"time"
)

// New returns a UUID v7 string (xxxxxxxx-xxxx-7xxx-yxxx-xxxxxxxxxxxx).
// The first 48 bits encode the current Unix millisecond timestamp so UUIDs
// sort chronologically, matching the behaviour of github.com/google/uuid.NewV7.
func New() (string, error) {
	var b [16]byte
	ms := time.Now().UnixMilli()
	b[0] = byte(ms >> 40)
	b[1] = byte(ms >> 32)
	b[2] = byte(ms >> 24)
	b[3] = byte(ms >> 16)
	b[4] = byte(ms >> 8)
	b[5] = byte(ms)
	if _, err := rand.Read(b[6:]); err != nil {
		return "", err
	}
	b[6] = (b[6] & 0x0f) | 0x70 // version 7
	b[8] = (b[8] & 0x3f) | 0x80 // variant 10xx
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16]), nil
}

// IsValid reports whether s is a syntactically valid UUID string
// (8-4-4-4-12 lowercase hex, no version/variant enforcement).
func IsValid(s string) bool {
	if len(s) != 36 {
		return false
	}
	for i, c := range s {
		switch i {
		case 8, 13, 18, 23:
			if c != '-' {
				return false
			}
		default:
			if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
				return false
			}
		}
	}
	return true
}
