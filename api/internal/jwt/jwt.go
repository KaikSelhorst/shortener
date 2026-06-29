// Package jwt implements HS256 JWT signing and verification using only stdlib.
package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

var header = base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))

// Sign creates a HS256 JWT token from an arbitrary JSON-marshalable payload.
func Sign(payload any, secret []byte) (string, error) {
	raw, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	body := header + "." + base64.RawURLEncoding.EncodeToString(raw)
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(body))
	return body + "." + base64.RawURLEncoding.EncodeToString(mac.Sum(nil)), nil
}

// Verify validates the HS256 signature and unmarshals the payload into dst.
// Returns an error if the signature is invalid or the payload cannot be decoded.
func Verify(token string, secret []byte, dst any) error {
	parts := strings.SplitN(token, ".", 3)
	if len(parts) != 3 {
		return fmt.Errorf("jwt: invalid token format")
	}
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(parts[0] + "." + parts[1]))
	expected := mac.Sum(nil)

	got, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil || !hmac.Equal(expected, got) {
		return fmt.Errorf("jwt: invalid signature")
	}
	raw, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return fmt.Errorf("jwt: invalid payload encoding")
	}
	return json.Unmarshal(raw, dst)
}

// NumericDate is a Unix timestamp that marshals/unmarshals as a JSON number,
// compatible with the JWT spec (RFC 7519 §2).
type NumericDate int64

func NewNumericDate(t time.Time) NumericDate { return NumericDate(t.Unix()) }
func (n NumericDate) Time() time.Time        { return time.Unix(int64(n), 0) }
func (n NumericDate) IsZero() bool           { return n == 0 }
