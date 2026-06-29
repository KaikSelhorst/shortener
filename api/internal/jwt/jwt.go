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

// Expirable may be implemented by claims types to opt in to automatic expiry
// enforcement inside Verify.
type Expirable interface {
	ExpiresAt() time.Time
}

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
	if err := json.Unmarshal(raw, dst); err != nil {
		return err
	}
	if e, ok := dst.(Expirable); ok {
		if time.Now().After(e.ExpiresAt()) {
			return fmt.Errorf("jwt: token expired")
		}
	}
	return nil
}

type NumericDate int64

func NewNumericDate(t time.Time) NumericDate { return NumericDate(t.Unix()) }
func (n NumericDate) Time() time.Time        { return time.Unix(int64(n), 0) }
func (n NumericDate) IsZero() bool           { return n == 0 }
