package service

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"math"
	"net/url"
	"time"
)

const (
	totpDigits = 6
	totpPeriod = 30 // seconds per time step
	totpWindow = 1  // accept ±1 period to handle clock drift
)

// GenerateTOTPSecret generates a 160-bit random secret encoded as base32 (no padding).
func GenerateTOTPSecret() (string, error) {
	b := make([]byte, 20)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generate totp secret: %w", err)
	}
	return base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(b), nil
}

// TOTPUri returns the otpauth:// URI used to generate a QR code on the frontend.
// Both the label components and the issuer query parameter are percent-encoded
// so that email addresses with '+', '&', or other reserved characters produce a
// well-formed URI that all authenticator apps can parse correctly.
func TOTPUri(secret, email, issuer string) string {
	label := url.PathEscape(issuer) + ":" + url.PathEscape(email)
	params := url.Values{
		"secret":    {secret},
		"issuer":    {issuer},
		"algorithm": {"SHA1"},
		"digits":    {fmt.Sprintf("%d", totpDigits)},
		"period":    {fmt.Sprintf("%d", totpPeriod)},
	}
	return "otpauth://totp/" + label + "?" + params.Encode()
}

// ValidateTOTP checks whether code is valid for the given base32-encoded secret,
// accepting codes from the current time step and ±totpWindow steps to tolerate clock drift.
// Returns an error if the secret is malformed rather than silently returning false.
func ValidateTOTP(secret, code string) (bool, error) {
	counter := time.Now().Unix() / totpPeriod
	for delta := int64(-totpWindow); delta <= int64(totpWindow); delta++ {
		got, err := generateTOTP(secret, counter+delta)
		if err != nil {
			return false, fmt.Errorf("totp validation: %w", err)
		}
		if got == code {
			return true, nil
		}
	}
	return false, nil
}

// generateTOTP computes the TOTP code for a given secret and counter value
// following RFC 6238 (TOTP) and RFC 4226 (HOTP).
func generateTOTP(secret string, counter int64) (string, error) {
	key, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(secret)
	if err != nil {
		return "", fmt.Errorf("decode totp secret: %w", err)
	}

	// Encode counter as 8-byte big-endian unsigned integer.
	msg := make([]byte, 8)
	binary.BigEndian.PutUint64(msg, uint64(counter))

	// Compute HMAC-SHA1.
	mac := hmac.New(sha1.New, key)
	mac.Write(msg)
	h := mac.Sum(nil)

	// Dynamic truncation (RFC 4226 §5.4).
	offset := h[len(h)-1] & 0x0f
	code := (int(h[offset])&0x7f)<<24 |
		(int(h[offset+1])&0xff)<<16 |
		(int(h[offset+2])&0xff)<<8 |
		int(h[offset+3]&0xff)

	otp := code % int(math.Pow10(totpDigits))
	return fmt.Sprintf("%0*d", totpDigits, otp), nil
}
