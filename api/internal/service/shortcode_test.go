package service

import (
	"math"
	"testing"
)

func TestShortcodeService_RoundTrip(t *testing.T) {
	t.Parallel()
	svc, err := NewShortcodeService([]byte("test-shortcode-secret-32-chars!!"))
	if err != nil {
		t.Fatalf("NewShortcodeService: %v", err)
	}

	cases := []uint64{1, 42, 100, 999, 1234567890}
	for _, id := range cases {
		code, err := svc.GenerateShortCode(id)
		if err != nil {
			t.Errorf("GenerateShortCode(%d): %v", id, err)
			continue
		}
		got := svc.DecodeShortCode(code)
		if got != id {
			t.Errorf("round-trip id=%d: encoded %q decoded to %d", id, code, got)
		}
	}
}

func TestShortcodeService_ZeroIDError(t *testing.T) {
	t.Parallel()
	svc, _ := NewShortcodeService([]byte("test-shortcode-secret-32-chars!!"))
	_, err := svc.GenerateShortCode(0)
	if err == nil {
		t.Error("expected error for id=0, got nil")
	}
}

func TestShortcodeService_Deterministic(t *testing.T) {
	t.Parallel()
	svc, _ := NewShortcodeService([]byte("test-shortcode-secret-32-chars!!"))
	id := uint64(42)
	a, _ := svc.GenerateShortCode(id)
	b, _ := svc.GenerateShortCode(id)
	if a != b {
		t.Errorf("non-deterministic: got %q and %q for same id", a, b)
	}
}

func TestShortcodeService_MaxUint64(t *testing.T) {
	t.Parallel()
	svc, _ := NewShortcodeService([]byte("test-shortcode-secret-32-chars!!"))
	id := uint64(math.MaxUint64)
	code, err := svc.GenerateShortCode(id)
	if err != nil {
		t.Fatalf("GenerateShortCode(MaxUint64): %v", err)
	}
	got := svc.DecodeShortCode(code)
	if got != id {
		t.Errorf("MaxUint64 round-trip: expected %d, got %d", id, got)
	}
}

func TestShortcodeService_DifferentIDsDifferentCodes(t *testing.T) {
	t.Parallel()
	svc, _ := NewShortcodeService([]byte("test-shortcode-secret-32-chars!!"))
	a, _ := svc.GenerateShortCode(1)
	b, _ := svc.GenerateShortCode(2)
	if a == b {
		t.Errorf("different IDs produced same code: %q", a)
	}
}
