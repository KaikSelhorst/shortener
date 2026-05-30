package testutil

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"math"
	"time"
)

// GenerateCurrentTOTPCode computes the current 6-digit TOTP code for a base32-encoded secret.
// Mirrors the RFC 6238 / RFC 4226 logic from the service package so handler tests can
// produce valid codes without depending on the service internals.
func GenerateCurrentTOTPCode(secret string) (string, error) {
	key, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(secret)
	if err != nil {
		return "", fmt.Errorf("testutil: decode totp secret: %w", err)
	}

	counter := time.Now().Unix() / 30
	msg := make([]byte, 8)
	binary.BigEndian.PutUint64(msg, uint64(counter))

	mac := hmac.New(sha1.New, key)
	mac.Write(msg)
	h := mac.Sum(nil)

	offset := h[len(h)-1] & 0x0f
	code := (int(h[offset])&0x7f)<<24 |
		(int(h[offset+1])&0xff)<<16 |
		(int(h[offset+2])&0xff)<<8 |
		int(h[offset+3]&0xff)

	otp := code % int(math.Pow10(6))
	return fmt.Sprintf("%06d", otp), nil
}
