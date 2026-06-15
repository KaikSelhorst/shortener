package service_test

import (
	"context"
	"strings"
	"testing"

	"github.com/KaikSelhorst/shortener/internal/model"
	"github.com/KaikSelhorst/shortener/internal/repository/fakes"
	"github.com/KaikSelhorst/shortener/internal/service"
)

func newWebhookService() *service.WebhookService {
	return service.NewWebhookService(fakes.NewWebhookRepo(), []byte("test-webhook-secret-key"))
}

// --- EncryptSecret / DecryptSecret ---

func TestWebhookService_EncryptDecrypt_RoundTrip(t *testing.T) {
	t.Parallel()
	svc := newWebhookService()
	plain := "my-super-secret-value"

	enc, err := svc.EncryptSecret(plain)
	if err != nil {
		t.Fatalf("EncryptSecret: %v", err)
	}
	got, err := svc.DecryptSecret(enc)
	if err != nil {
		t.Fatalf("DecryptSecret: %v", err)
	}
	if got != plain {
		t.Errorf("expected %q, got %q", plain, got)
	}
}

func TestWebhookService_EncryptSecret_UniqueNonce(t *testing.T) {
	t.Parallel()
	svc := newWebhookService()
	enc1, _ := svc.EncryptSecret("same-value")
	enc2, _ := svc.EncryptSecret("same-value")
	if enc1 == enc2 {
		t.Error("expected different ciphertexts due to random nonce")
	}
}

func TestWebhookService_DecryptSecret_WrongKey(t *testing.T) {
	t.Parallel()
	svc1 := service.NewWebhookService(nil, []byte("key-a"))
	svc2 := service.NewWebhookService(nil, []byte("key-b"))

	enc, _ := svc1.EncryptSecret("secret")
	_, err := svc2.DecryptSecret(enc)
	if err == nil {
		t.Error("expected error when decrypting with wrong key")
	}
}

func TestWebhookService_DecryptSecret_InvalidBase64(t *testing.T) {
	t.Parallel()
	svc := newWebhookService()
	_, err := svc.DecryptSecret("not-valid-base64!!!")
	if err == nil {
		t.Error("expected error for invalid base64 input")
	}
}

// --- Sign ---

func TestWebhookService_Sign_Format(t *testing.T) {
	t.Parallel()
	sig := service.Sign([]byte(`{"event":"link.clicked"}`), "my-secret")
	if !strings.HasPrefix(sig, "sha256=") {
		t.Errorf("signature should start with 'sha256=', got %q", sig)
	}
	// "sha256=" + 64 hex chars
	if len(sig) != 7+64 {
		t.Errorf("unexpected signature length: %d", len(sig))
	}
}

func TestWebhookService_Sign_Deterministic(t *testing.T) {
	t.Parallel()
	payload := []byte(`{"event":"link.clicked"}`)
	if service.Sign(payload, "secret") != service.Sign(payload, "secret") {
		t.Error("Sign must be deterministic for the same payload and secret")
	}
}

func TestWebhookService_Sign_DifferentSecrets(t *testing.T) {
	t.Parallel()
	payload := []byte(`{"event":"link.clicked"}`)
	if service.Sign(payload, "secret-a") == service.Sign(payload, "secret-b") {
		t.Error("different secrets should produce different signatures")
	}
}

// --- GenerateWebhookSecret ---

func TestWebhookService_GenerateSecret_IsHex64(t *testing.T) {
	t.Parallel()
	s, err := service.GenerateWebhookSecret()
	if err != nil {
		t.Fatalf("GenerateWebhookSecret: %v", err)
	}
	if len(s) != 64 {
		t.Errorf("expected 64 hex chars (32 bytes), got %d", len(s))
	}
}

func TestWebhookService_GenerateSecret_Unique(t *testing.T) {
	t.Parallel()
	s1, _ := service.GenerateWebhookSecret()
	s2, _ := service.GenerateWebhookSecret()
	if s1 == s2 {
		t.Error("secrets should be unique across calls")
	}
}

// --- Dispatch ---

func TestWebhookService_Dispatch_CreatesDeliveryForMatchingWebhook(t *testing.T) {
	t.Parallel()
	repo := fakes.NewWebhookRepo()
	svc := service.NewWebhookService(repo, []byte("test-key"))

	_ = repo.Create(context.Background(), &model.Webhook{
		ProjectID: 1,
		URL:       "https://example.com/hook",
		Secret:    "s",
		Events:    []string{"link.clicked"},
		Enabled:   true,
	})
	// different event — should not receive delivery
	_ = repo.Create(context.Background(), &model.Webhook{
		ProjectID: 1,
		URL:       "https://example.com/hook2",
		Secret:    "s",
		Events:    []string{"link.created"},
		Enabled:   true,
	})

	svc.Dispatch(context.Background(), 1, "link.clicked", map[string]any{"event": "link.clicked"})

	if got := repo.DeliveryCount(); got != 1 {
		t.Errorf("expected 1 delivery, got %d", got)
	}
}

func TestWebhookService_Dispatch_SkipsDisabledWebhooks(t *testing.T) {
	t.Parallel()
	repo := fakes.NewWebhookRepo()
	svc := service.NewWebhookService(repo, []byte("test-key"))

	_ = repo.Create(context.Background(), &model.Webhook{
		ProjectID: 1,
		URL:       "https://example.com/hook",
		Secret:    "s",
		Events:    []string{"link.clicked"},
		Enabled:   false,
	})

	svc.Dispatch(context.Background(), 1, "link.clicked", map[string]any{"event": "link.clicked"})

	if got := repo.DeliveryCount(); got != 0 {
		t.Errorf("expected 0 deliveries for disabled webhook, got %d", got)
	}
}

func TestWebhookService_Dispatch_MultipleWebhooksSameEvent(t *testing.T) {
	t.Parallel()
	repo := fakes.NewWebhookRepo()
	svc := service.NewWebhookService(repo, []byte("test-key"))

	for range 3 {
		_ = repo.Create(context.Background(), &model.Webhook{
			ProjectID: 1,
			URL:       "https://example.com/hook",
			Secret:    "s",
			Events:    []string{"link.clicked"},
			Enabled:   true,
		})
	}

	svc.Dispatch(context.Background(), 1, "link.clicked", map[string]any{"event": "link.clicked"})

	if got := repo.DeliveryCount(); got != 3 {
		t.Errorf("expected 3 deliveries (one per webhook), got %d", got)
	}
}
