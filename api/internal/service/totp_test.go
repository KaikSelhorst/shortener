package service

import (
	"strings"
	"testing"
	"time"
)

func TestGenerateTOTPSecret_Format(t *testing.T) {
	t.Parallel()
	secret, err := GenerateTOTPSecret()
	if err != nil {
		t.Fatalf("GenerateTOTPSecret: %v", err)
	}
	// 20 bytes base32 without padding = 32 characters
	if len(secret) != 32 {
		t.Errorf("expected 32 chars, got %d: %q", len(secret), secret)
	}
	if strings.Contains(secret, "=") {
		t.Errorf("unexpected padding in secret: %q", secret)
	}
}

func TestGenerateTOTPSecret_Unique(t *testing.T) {
	t.Parallel()
	a, _ := GenerateTOTPSecret()
	b, _ := GenerateTOTPSecret()
	if a == b {
		t.Error("GenerateTOTPSecret returned the same secret twice")
	}
}

func TestValidateTOTP_CurrentCode(t *testing.T) {
	t.Parallel()
	secret, _ := GenerateTOTPSecret()
	counter := time.Now().Unix() / totpPeriod
	code, err := generateTOTP(secret, counter)
	if err != nil {
		t.Fatalf("generateTOTP: %v", err)
	}
	ok, err := ValidateTOTP(secret, code)
	if err != nil {
		t.Fatalf("ValidateTOTP: %v", err)
	}
	if !ok {
		t.Error("expected current code to be valid")
	}
}

func TestValidateTOTP_PreviousPeriod(t *testing.T) {
	t.Parallel()
	secret, _ := GenerateTOTPSecret()
	// Code from 1 period ago is within the ±1 window.
	counter := time.Now().Unix()/totpPeriod - 1
	code, err := generateTOTP(secret, counter)
	if err != nil {
		t.Fatalf("generateTOTP: %v", err)
	}
	ok, err := ValidateTOTP(secret, code)
	if err != nil {
		t.Fatalf("ValidateTOTP: %v", err)
	}
	if !ok {
		t.Error("expected code from previous period to be valid (±1 window)")
	}
}

func TestValidateTOTP_ExpiredCode(t *testing.T) {
	t.Parallel()
	secret, _ := GenerateTOTPSecret()
	// Code from 2 periods ago is outside the ±1 window.
	counter := time.Now().Unix()/totpPeriod - 2
	code, err := generateTOTP(secret, counter)
	if err != nil {
		t.Fatalf("generateTOTP: %v", err)
	}
	ok, err := ValidateTOTP(secret, code)
	if err != nil {
		t.Fatalf("ValidateTOTP: %v", err)
	}
	if ok {
		t.Error("expected code from 2 periods ago to be invalid")
	}
}

func TestValidateTOTP_MalformedSecret(t *testing.T) {
	t.Parallel()
	_, err := ValidateTOTP("not-valid-base32!!!", "123456")
	if err == nil {
		t.Error("expected error for malformed secret, got nil")
	}
}

func TestTOTPUri_ContainsRequiredFields(t *testing.T) {
	t.Parallel()
	secret := "JBSWY3DPEHPK3PXP"
	email := "user@example.com"
	issuer := "MyApp"
	uri := TOTPUri(secret, email, issuer)

	checks := []string{
		"otpauth://totp/",
		"secret=" + secret,
		"issuer=" + issuer,
		"algorithm=SHA1",
		"digits=6",
		"period=30",
	}
	for _, want := range checks {
		if !strings.Contains(uri, want) {
			t.Errorf("expected URI to contain %q\ngot: %q", want, uri)
		}
	}
}

func TestTOTPUri_LabelFormat(t *testing.T) {
	t.Parallel()
	uri := TOTPUri("SECRET", "user@example.com", "MyApp")
	// Label must appear as "Issuer:Email" in the path portion.
	if !strings.Contains(uri, "MyApp:user%40example.com") &&
		!strings.Contains(uri, "MyApp:user@example.com") {
		t.Errorf("label not found in expected format\ngot: %q", uri)
	}
}
